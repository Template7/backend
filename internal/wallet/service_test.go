package wallet

import (
	"context"
	"math"
	"net/http"
	"testing"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
)

// fakeDb implements only the db.Client methods the wallet service uses;
// calling any other method panics via the nil embedded interface.
type fakeDb struct {
	db.Client
	userWallets map[string][]string // userId -> walletIds

	transferCalls int
	lastMoney     models.Money
}

func (f *fakeDb) GetUserWallets(_ context.Context, userId string) []models.UserWalletBalance {
	var r []models.UserWalletBalance
	for _, wId := range f.userWallets[userId] {
		r = append(r, models.UserWalletBalance{UserId: userId, WalletId: wId, Currency: "usd"})
	}
	return r
}

func (f *fakeDb) Transfer(_ context.Context, _ string, _ string, money models.Money, _ string) error {
	f.transferCalls++
	f.lastMoney = money
	return nil
}

func (f *fakeDb) Deposit(_ context.Context, _ string, money models.Money, _ string) error {
	f.lastMoney = money
	return nil
}

func (f *fakeDb) Withdraw(_ context.Context, _ string, money models.Money, _ string) error {
	f.lastMoney = money
	return nil
}

func newTestService() (*Service, *fakeDb) {
	f := &fakeDb{userWallets: map[string][]string{
		"alice": {"wallet-alice"},
		"bob":   {"wallet-bob"},
	}}
	return New(f, logger.New("error", "console", "test")), f
}

func TestTransfer_RejectsWalletNotOwnedByCaller(t *testing.T) {
	s, f := newTestService()

	err := s.Transfer(context.Background(), "alice", "wallet-bob", "wallet-alice", v1.Currency_usd, 100, "")

	t7Err, ok := t7Error.ToT7Error(err)
	if !ok || t7Err.Code != t7Error.WalletNotFound.Code || t7Err.GetStatus() != http.StatusNotFound {
		t.Fatalf("want WalletNotFound/404, got %v", err)
	}
	if f.transferCalls != 0 {
		t.Fatalf("db.Transfer must not be called, got %d calls", f.transferCalls)
	}
}

func TestTransfer_RejectsSameWallet(t *testing.T) {
	s, f := newTestService()

	err := s.Transfer(context.Background(), "alice", "wallet-alice", "wallet-alice", v1.Currency_usd, 100, "")

	if t7Err, ok := t7Error.ToT7Error(err); !ok || t7Err.Code != t7Error.InvalidBody.Code {
		t.Fatalf("want InvalidBody, got %v", err)
	}
	if f.transferCalls != 0 {
		t.Fatalf("db.Transfer must not be called")
	}
}

func TestTransfer_OwnerSucceeds(t *testing.T) {
	s, f := newTestService()

	if err := s.Transfer(context.Background(), "alice", "wallet-alice", "wallet-bob", v1.Currency_usd, 100, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.transferCalls != 1 || !f.lastMoney.Amount.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("want one transfer of 100, got calls=%d money=%v", f.transferCalls, f.lastMoney)
	}
}

func TestAmountAboveMaxInt32StaysPositive(t *testing.T) {
	s, f := newTestService()
	const amount = uint32(math.MaxUint32) // was int32(-1) before the fix

	if err := s.Deposit(context.Background(), "wallet-alice", v1.Currency_usd, amount, ""); err != nil {
		t.Fatalf("deposit: %v", err)
	}
	if !f.lastMoney.Amount.Equal(decimal.NewFromInt(math.MaxUint32)) {
		t.Fatalf("deposit amount = %s, want %d", f.lastMoney.Amount, uint32(math.MaxUint32))
	}

	if err := s.Withdraw(context.Background(), "wallet-alice", v1.Currency_usd, math.MaxInt32+1, ""); err != nil {
		t.Fatalf("withdraw: %v", err)
	}
	if !f.lastMoney.Amount.IsPositive() {
		t.Fatalf("withdraw amount must be positive, got %s", f.lastMoney.Amount)
	}
}

func TestZeroAmountRejected(t *testing.T) {
	s, _ := newTestService()

	for name, err := range map[string]error{
		"deposit":  s.Deposit(context.Background(), "wallet-alice", v1.Currency_usd, 0, ""),
		"withdraw": s.Withdraw(context.Background(), "wallet-alice", v1.Currency_usd, 0, ""),
		"transfer": s.Transfer(context.Background(), "alice", "wallet-alice", "wallet-bob", v1.Currency_usd, 0, ""),
	} {
		if t7Err, ok := t7Error.ToT7Error(err); !ok || t7Err.Code != t7Error.InvalidBody.Code {
			t.Errorf("%s: want InvalidBody, got %v", name, err)
		}
	}
}
