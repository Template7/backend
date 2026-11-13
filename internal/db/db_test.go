package db

import (
	"context"
	"errors"
	"testing"

	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// newTestClient returns a client on a fresh in-memory SQLite DB.
// One connection only: anything written outside the open transaction would block.
func newTestClient(t *testing.T) *client {
	t.Helper()
	log := logger.New("error", "console", "test")

	core, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: gormLogger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	sqlDb, _ := core.DB()
	sqlDb.SetMaxOpenConns(1)
	if err := core.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Balance{}); err != nil {
		t.Fatal(err)
	}
	// history models reuse index names across tables (e.g. "wcc"), which MySQL allows
	// but SQLite doesn't, so create those tables by hand without indexes
	for _, ddl := range []string{
		`CREATE TABLE deposit_history (id INTEGER PRIMARY KEY, wallet_id TEXT, currency TEXT, amount NUMERIC,
			balance_before NUMERIC, balance_after NUMERIC, note TEXT, created_at DATETIME)`,
		`CREATE TABLE withdraw_history (id INTEGER PRIMARY KEY, wallet_id TEXT, currency TEXT, amount NUMERIC,
			balance_before NUMERIC, balance_after NUMERIC, note TEXT, created_at DATETIME)`,
		`CREATE TABLE transfer_history (id INTEGER PRIMARY KEY, from_wallet_id TEXT, to_wallet_id TEXT, currency TEXT,
			amount NUMERIC, sender_balance_before NUMERIC, sender_balance_after NUMERIC,
			receiver_balance_before NUMERIC, receiver_balance_after NUMERIC, note TEXT, created_at DATETIME)`,
	} {
		if err := core.Exec(ddl).Error; err != nil {
			t.Fatal(err)
		}
	}

	c := &client{log: log}
	c.sql.core = core
	return c
}

func seedUser(t *testing.T, c *client, userId, walletId string, usd int64) {
	t.Helper()
	if err := c.sql.core.Create(&models.User{Id: userId, Username: userId, Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := c.sql.core.Create(&models.Wallet{Id: walletId, UserId: userId}).Error; err != nil {
		t.Fatal(err)
	}
	b := models.Balance{WalletId: walletId, Money: models.Money{Currency: "usd", Amount: decimal.NewFromInt(usd)}}
	if err := c.sql.core.Create(&b).Error; err != nil {
		t.Fatal(err)
	}
}

func usd(n int64) models.Money {
	return models.Money{Currency: "usd", Amount: decimal.NewFromInt(n)}
}

func count(t *testing.T, c *client, model interface{}) int64 {
	t.Helper()
	var n int64
	if err := c.sql.core.Model(model).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}

func TestWithdraw_UpdatesBalanceAndHistory(t *testing.T) {
	c := newTestClient(t)
	seedUser(t, c, "u1", "w1", 100)

	if err := c.Withdraw(context.Background(), "w1", usd(30), "note"); err != nil {
		t.Fatalf("withdraw: %v", err)
	}
	bal, _ := c.GetBalance(context.Background(), "w1", "usd")
	if !bal.Equal(decimal.NewFromInt(70)) {
		t.Fatalf("balance = %s, want 70", bal)
	}
	if n := count(t, c, &models.WithdrawHistory{}); n != 1 {
		t.Fatalf("withdraw history rows = %d, want 1", n)
	}
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	c := newTestClient(t)
	seedUser(t, c, "u1", "w1", 10)

	err := c.Withdraw(context.Background(), "w1", usd(11), "")
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("want ErrInsufficientBalance, got %v", err)
	}
	bal, _ := c.GetBalance(context.Background(), "w1", "usd")
	if !bal.Equal(decimal.NewFromInt(10)) {
		t.Fatalf("balance changed to %s", bal)
	}
	if n := count(t, c, &models.WithdrawHistory{}); n != 0 {
		t.Fatalf("no history expected, got %d", n)
	}
}

func TestTransfer_MovesMoneyAndWritesHistory(t *testing.T) {
	c := newTestClient(t)
	seedUser(t, c, "u1", "w1", 100)
	seedUser(t, c, "u2", "w2", 5)

	if err := c.Transfer(context.Background(), "w1", "w2", usd(40), ""); err != nil {
		t.Fatalf("transfer: %v", err)
	}
	b1, _ := c.GetBalance(context.Background(), "w1", "usd")
	b2, _ := c.GetBalance(context.Background(), "w2", "usd")
	if !b1.Equal(decimal.NewFromInt(60)) || !b2.Equal(decimal.NewFromInt(45)) {
		t.Fatalf("balances = %s / %s, want 60 / 45", b1, b2)
	}
	if n := count(t, c, &models.TransferHistory{}); n != 1 {
		t.Fatalf("transfer history rows = %d, want 1", n)
	}
}

func TestTransfer_InsufficientBalanceChangesNothing(t *testing.T) {
	c := newTestClient(t)
	seedUser(t, c, "u1", "w1", 10)
	seedUser(t, c, "u2", "w2", 0)

	if err := c.Transfer(context.Background(), "w1", "w2", usd(11), ""); !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("want ErrInsufficientBalance, got %v", err)
	}
	b2, _ := c.GetBalance(context.Background(), "w2", "usd")
	if !b2.IsZero() {
		t.Fatalf("receiver balance changed to %s", b2)
	}
}

func TestDeleteUser_RemovesUserWalletAndBalances(t *testing.T) {
	c := newTestClient(t)
	seedUser(t, c, "u1", "w1", 10)
	seedUser(t, c, "u2", "w2", 10)

	if err := c.DeleteUser(context.Background(), "u1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n := count(t, c, &models.User{}); n != 1 {
		t.Fatalf("users left = %d, want 1", n)
	}
	if n := count(t, c, &models.Wallet{}); n != 1 {
		t.Fatalf("wallets left = %d, want 1", n)
	}
	if n := count(t, c, &models.Balance{}); n != 1 {
		t.Fatalf("balances left = %d, want 1", n)
	}
}

func TestDeleteUser_WithoutWallet(t *testing.T) {
	c := newTestClient(t)
	if err := c.sql.core.Create(&models.User{Id: "u1", Username: "u1", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}

	if err := c.DeleteUser(context.Background(), "u1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n := count(t, c, &models.User{}); n != 0 {
		t.Fatal("user without wallet must be deleted too")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	c := newTestClient(t)
	if err := c.DeleteUser(context.Background(), "nobody"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("want ErrRecordNotFound, got %v", err)
	}
}
