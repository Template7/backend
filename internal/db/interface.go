package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Client interface {
	// user
	CreateUser(ctx context.Context, data entity.User) (err error)
	GetUser(ctx context.Context, username string) (data entity.User, err error)
	GetUserById(ctx context.Context, userId string) (data entity.User, err error)
	UpdateUserInfo(ctx context.Context, userId string, info entity.UserInfo) (err error)
	GetUserWallets(ctx context.Context, userId string) (data []entity.UserWalletBalance)
	DeleteUser(ctx context.Context, userId string) (err error)

	// wallet
	GetWalletBalances(ctx context.Context, walletId string) (data []entity.WalletBalance, err error)
	Deposit(ctx context.Context, walletId string, money entity.Money, note string) (err error)
	Withdraw(ctx context.Context, walletId string, money entity.Money, note string) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money, note string) (err error)
	GetBalance(ctx context.Context, walletId string, currency string) (decimal.Decimal, error)
	getWalletsBalance(ctx context.Context, tx *gorm.DB, walletId []string, currency string) (data []entity.Balance, err error)

	// history
	createDepositHistory(ctx context.Context, tx *gorm.DB, data entity.DepositHistory) (err error)
	createWithdrawHistory(ctx context.Context, tx *gorm.DB, data entity.WithdrawHistory) (err error)
	createTransferHistory(ctx context.Context, tx *gorm.DB, data entity.TransferHistory) (err error)
	GetWalletBalanceHistory(ctx context.Context, walletId string, currency string) ([]entity.WalletBalanceHistory, error)
}
