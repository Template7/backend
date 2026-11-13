package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/backend/internal/db/dbtest"
	"github.com/Template7/common/models"
	"github.com/shopspring/decimal"
)

var ctx = context.Background()

func TestTransaction_RollsBackOnError(t *testing.T) {
	c, core := dbtest.New(t)
	dbtest.SeedUser(t, core, "u1", "w1", 100)

	boom := errors.New("boom")
	err := c.Transaction(ctx, func(tx db.Client) error {
		if err := tx.AddBalance(ctx, "w1", "usd", decimal.NewFromInt(-30)); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("want boom, got %v", err)
	}
	if b := dbtest.Balance(t, core, "w1", "usd"); !b.Equal(decimal.NewFromInt(100)) {
		t.Fatalf("balance = %s, want 100 after rollback", b)
	}
}

func TestLockBalances(t *testing.T) {
	c, core := dbtest.New(t)
	dbtest.SeedUser(t, core, "u1", "w1", 100)
	dbtest.SeedUser(t, core, "u2", "w2", 5)

	var got map[string]decimal.Decimal
	err := c.Transaction(ctx, func(tx db.Client) (err error) {
		got, err = tx.LockBalances(ctx, []string{"w2", "w1", "missing"}, "usd")
		return
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || !got["w1"].Equal(decimal.NewFromInt(100)) || !got["w2"].Equal(decimal.NewFromInt(5)) {
		t.Fatalf("got %v", got)
	}
}

func TestAddBalance_UnknownWalletIsNotFound(t *testing.T) {
	c, _ := dbtest.New(t)
	if err := c.AddBalance(ctx, "missing", "usd", decimal.NewFromInt(1)); !db.IsNotFound(err) {
		t.Fatalf("want not found, got %v", err)
	}
}

func TestCreateWallet_CreatesZeroBalances(t *testing.T) {
	c, core := dbtest.New(t)
	if err := c.CreateWallet(ctx, models.Wallet{Id: "w1", UserId: "u1"}, dbtest.Currencies); err != nil {
		t.Fatal(err)
	}
	if n := dbtest.Count(t, core, &models.Balance{}); n != int64(len(dbtest.Currencies)) {
		t.Fatalf("balances = %d", n)
	}
}

func TestDeleteUser_RemovesUserWalletAndBalances(t *testing.T) {
	c, core := dbtest.New(t)
	dbtest.SeedUser(t, core, "u1", "w1", 10)
	dbtest.SeedUser(t, core, "u2", "w2", 10)

	if err := c.DeleteUser(ctx, "u1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n := dbtest.Count(t, core, &models.User{}); n != 1 {
		t.Fatalf("users left = %d, want 1", n)
	}
	if n := dbtest.Count(t, core, &models.Wallet{}); n != 1 {
		t.Fatalf("wallets left = %d, want 1", n)
	}
	if n := dbtest.Count(t, core, &models.Balance{}); n != int64(len(dbtest.Currencies)) {
		t.Fatalf("balances left = %d", n)
	}
}

func TestDeleteUser_WithoutWallet(t *testing.T) {
	c, core := dbtest.New(t)
	if err := core.Create(&models.User{Id: "u1", Username: "u1", Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := c.DeleteUser(ctx, "u1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n := dbtest.Count(t, core, &models.User{}); n != 0 {
		t.Fatal("user without wallet must be deleted too")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	c, _ := dbtest.New(t)
	if err := c.DeleteUser(ctx, "nobody"); !db.IsNotFound(err) {
		t.Fatalf("want not found, got %v", err)
	}
}
