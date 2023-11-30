package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
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
	Deposit(ctx context.Context, walletId string, money entity.Money) (err error)
	Withdraw(ctx context.Context, walletId string, money entity.Money) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error)

	// history
	CreateDepositHistory(ctx context.Context, data entity.DepositHistory) (err error)
	CreateWithdrawHistory(ctx context.Context, data entity.WithdrawHistory) (err error)
	CreateTransferHistory(ctx context.Context, data entity.TransferHistory) (err error)
}
