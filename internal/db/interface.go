// Package db is the storage layer: plain reads and writes, no business rules.
// Callers that need several writes to succeed together use Transaction.
package db

import (
	"context"
	"errors"

	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ErrNotFound is returned when the requested row doesn't exist.
var ErrNotFound = gorm.ErrRecordNotFound

// IsNotFound reports whether err means the row doesn't exist.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

type Client interface {
	UserRepo
	WalletRepo
	HistoryRepo

	// Transaction runs fn in one DB transaction. The Client passed to fn uses that
	// transaction; returning an error from fn rolls everything back.
	Transaction(ctx context.Context, fn func(tx Client) error) error
}

type UserRepo interface {
	CreateUser(ctx context.Context, data models.User) error
	GetUser(ctx context.Context, username string) (models.User, error)
	GetUserById(ctx context.Context, userId string) (models.User, error)
	UpdateUserInfo(ctx context.Context, userId string, info models.UserInfo) error
	SetUserStatus(ctx context.Context, userId string, status authV1.AccountStatus) error
	// DeleteUser removes the user with its wallets and balances; history is kept.
	DeleteUser(ctx context.Context, userId string) error
}

type WalletRepo interface {
	// CreateWallet creates the wallet with a zero balance for each currency.
	CreateWallet(ctx context.Context, wallet models.Wallet, currencies []string) error
	GetWallet(ctx context.Context, walletId string) (models.Wallet, error)
	GetUserWallets(ctx context.Context, userId string) ([]models.UserWalletBalance, error)
	GetWalletBalances(ctx context.Context, walletId string) ([]models.WalletBalance, error)
	// LockBalances reads the balances of the given wallets in one currency and locks
	// those rows until the transaction ends. Use inside Transaction.
	LockBalances(ctx context.Context, walletIds []string, currency string) (map[string]decimal.Decimal, error)
	// AddBalance adds delta (may be negative) to one balance.
	AddBalance(ctx context.Context, walletId string, currency string, delta decimal.Decimal) error
}

type HistoryRepo interface {
	CreateDepositHistory(ctx context.Context, data models.DepositHistory) error
	CreateWithdrawHistory(ctx context.Context, data models.WithdrawHistory) error
	CreateTransferHistory(ctx context.Context, data models.TransferHistory) error
	GetWalletBalanceHistory(ctx context.Context, walletId string) ([]models.WalletBalanceHistory, error)
	GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]models.WalletBalanceHistory, error)
}
