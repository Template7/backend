package wallet

import (
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/Template7/backend/internal/db/dbtest"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var ctx = context.Background()

// alice owns wallet-alice (100 usd), bob owns wallet-bob (5 usd)
func setup(t *testing.T) (Service, *gorm.DB) {
	t.Helper()
	client, core := dbtest.New(t)
	dbtest.SeedUser(t, core, "alice", "wallet-alice", 100)
	dbtest.SeedUser(t, core, "bob", "wallet-bob", 5)
	return New(client, logger.New("error", "console", "test")), core
}

func wantCode(t *testing.T, err error, want *t7Error.Error) {
	t.Helper()
	got, ok := t7Error.ToT7Error(err)
	if !ok || got.Code != want.Code {
		t.Fatalf("want %s, got %v", want.Message, err)
	}
}

func usd(t *testing.T, core *gorm.DB, walletId string) decimal.Decimal {
	return dbtest.Balance(t, core, walletId, "usd")
}

func TestDepositWithdraw(t *testing.T) {
	s, core := setup(t)

	if err := s.Deposit(ctx, "alice", "wallet-alice", v1.Currency_usd, 50, "in"); err != nil {
		t.Fatal(err)
	}
	if err := s.Withdraw(ctx, "alice", "wallet-alice", v1.Currency_usd, 30, "out"); err != nil {
		t.Fatal(err)
	}
	if b := usd(t, core, "wallet-alice"); !b.Equal(decimal.NewFromInt(120)) {
		t.Fatalf("balance = %s, want 120", b)
	}

	history, err := s.GetBalanceHistory(ctx, "alice", "wallet-alice")
	if err != nil || len(history) != 2 {
		t.Fatalf("history = %v, %v", history, err)
	}
	var withdraw *v1.BalanceHistory
	for _, h := range history {
		if h.Direction == v1.Direction_withdraw {
			withdraw = h
		}
	}
	if withdraw == nil || withdraw.BalanceBefore != "150" || withdraw.BalanceAfter != "120" {
		t.Fatalf("withdraw history = %v, want before/after 150/120", withdraw)
	}
}

func TestWithdraw_InsufficientBalanceChangesNothing(t *testing.T) {
	s, core := setup(t)

	wantCode(t, s.Withdraw(ctx, "alice", "wallet-alice", v1.Currency_usd, 101, ""), t7Error.InsufficientBalance)
	if b := usd(t, core, "wallet-alice"); !b.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("balance changed to %s", b)
	}
	if n := dbtest.Count(t, core, &models.WithdrawHistory{}); n != 0 {
		t.Fatalf("history rows = %d, want 0", n)
	}
}

func TestTransfer(t *testing.T) {
	s, core := setup(t)

	if err := s.Transfer(ctx, "alice", "wallet-alice", "wallet-bob", v1.Currency_usd, 40, ""); err != nil {
		t.Fatal(err)
	}
	if a, b := usd(t, core, "wallet-alice"), usd(t, core, "wallet-bob"); !a.Equal(decimal.NewFromInt(60)) || !b.Equal(decimal.NewFromInt(45)) {
		t.Fatalf("balances = %s / %s, want 60 / 45", a, b)
	}

	// the receiver sees it as transferIn
	history, err := s.GetBalanceHistory(ctx, "bob", "wallet-bob")
	if err != nil || len(history) != 1 || history[0].Direction != v1.Direction_transferIn {
		t.Fatalf("receiver history = %v, %v", history, err)
	}
}

func TestTransfer_Rejections(t *testing.T) {
	s, core := setup(t)

	cases := []struct {
		name           string
		user, from, to string
		amount         uint32
		want           *t7Error.Error
	}{
		{"not the owner", "alice", "wallet-bob", "wallet-alice", 1, t7Error.WalletNotFound},
		{"same wallet", "alice", "wallet-alice", "wallet-alice", 1, t7Error.InvalidBody},
		{"unknown receiver", "alice", "wallet-alice", "nope", 1, t7Error.WalletNotFound},
		{"insufficient", "alice", "wallet-alice", "wallet-bob", 101, t7Error.InsufficientBalance},
		{"zero amount", "alice", "wallet-alice", "wallet-bob", 0, t7Error.InvalidBody},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wantCode(t, s.Transfer(ctx, tc.user, tc.from, tc.to, v1.Currency_usd, tc.amount, ""), tc.want)
		})
	}
	if a, b := usd(t, core, "wallet-alice"), usd(t, core, "wallet-bob"); !a.Equal(decimal.NewFromInt(100)) || !b.Equal(decimal.NewFromInt(5)) {
		t.Fatalf("balances changed: %s / %s", a, b)
	}
}

func TestOtherUsersWalletLooksMissing(t *testing.T) {
	s, _ := setup(t)

	_, err := s.GetWallet(ctx, "alice", "wallet-bob")
	wantCode(t, err, t7Error.WalletNotFound)
	if got, _ := t7Error.ToT7Error(err); got.GetStatus() != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", got.GetStatus())
	}

	wantCode(t, s.Deposit(ctx, "alice", "wallet-bob", v1.Currency_usd, 1, ""), t7Error.WalletNotFound)
	wantCode(t, s.Withdraw(ctx, "alice", "wallet-bob", v1.Currency_usd, 1, ""), t7Error.WalletNotFound)
	_, err = s.GetBalanceHistory(ctx, "alice", "wallet-bob")
	wantCode(t, err, t7Error.WalletNotFound)
	_, err = s.GetBalanceHistoryByCurrency(ctx, "alice", "wallet-bob", v1.Currency_usd)
	wantCode(t, err, t7Error.WalletNotFound)
}

func TestAmountAboveMaxInt32StaysPositive(t *testing.T) {
	s, core := setup(t)

	if err := s.Deposit(ctx, "alice", "wallet-alice", v1.Currency_usd, math.MaxUint32, ""); err != nil {
		t.Fatal(err)
	}
	want := decimal.NewFromInt(100 + math.MaxUint32)
	if b := usd(t, core, "wallet-alice"); !b.Equal(want) {
		t.Fatalf("balance = %s, want %s", b, want)
	}
}

func TestGetUserWallets(t *testing.T) {
	s, _ := setup(t)

	ws, err := s.GetUserWallets(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(ws) != 1 || ws[0].Id != "wallet-alice" || len(ws[0].Balances) != len(dbtest.Currencies) {
		t.Fatalf("wallets = %v", ws)
	}
}
