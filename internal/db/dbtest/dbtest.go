// Package dbtest provides an in-memory SQLite database for tests.
package dbtest

import (
	"testing"

	"github.com/Template7/backend/internal/db"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	"github.com/glebarez/sqlite"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Currencies created for every seeded wallet.
var Currencies = []string{"usd", "ntd", "cny", "jpy"}

// New returns a db.Client on a fresh in-memory SQLite DB, plus the raw handle for
// seeding and assertions. One connection only: a write made outside an open
// transaction blocks, so tests catch it.
func New(t *testing.T) (db.Client, *gorm.DB) {
	t.Helper()
	log := logger.New("error", "console", "test")

	core, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: gormLogger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	sqlDb, err := core.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDb.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDb.Close() })

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

	return db.New(core, log), core
}

// SeedUser creates a user with one wallet holding usd (and zero in the other currencies).
func SeedUser(t *testing.T, core *gorm.DB, userId, walletId string, usd int64) {
	t.Helper()
	if err := core.Create(&models.User{Id: userId, Username: userId, Password: "x"}).Error; err != nil {
		t.Fatal(err)
	}
	if err := core.Create(&models.Wallet{Id: walletId, UserId: userId}).Error; err != nil {
		t.Fatal(err)
	}
	for _, cur := range Currencies {
		amount := decimal.Zero
		if cur == "usd" {
			amount = decimal.NewFromInt(usd)
		}
		b := models.Balance{WalletId: walletId, Money: models.Money{Currency: cur, Amount: amount}}
		if err := core.Create(&b).Error; err != nil {
			t.Fatal(err)
		}
	}
}

// Balance returns the current balance of one wallet in one currency.
func Balance(t *testing.T, core *gorm.DB, walletId, currency string) decimal.Decimal {
	t.Helper()
	var b models.Balance
	if err := core.Where("wallet_id = ? AND currency = ?", walletId, currency).Take(&b).Error; err != nil {
		t.Fatal(err)
	}
	return b.Amount
}

// Count returns the number of rows of model.
func Count(t *testing.T, core *gorm.DB, model interface{}) int64 {
	t.Helper()
	var n int64
	if err := core.Model(model).Count(&n).Error; err != nil {
		t.Fatal(err)
	}
	return n
}
