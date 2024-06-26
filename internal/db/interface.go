package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/shopspring/decimal"
)

type Client interface {
	// user
	CreateUser(ctx context.Context, data entity.User) (err error)
	GetUser(ctx context.Context, username string) (data entity.User, err error)
	GetUserById(ctx context.Context, userId string) (data entity.User, err error)
	UpdateUserInfo(ctx context.Context, userId string, info entity.UserInfo) (err error)
	GetUserWallets(ctx context.Context, userId string) (data []entity.UserWalletBalance)
	DeleteUser(ctx context.Context, userId string) (err error)
	SetUserStatus(ctx context.Context, userId string, status authV1.AccountStatus) (err error)

	// wallet
	GetWalletBalances(ctx context.Context, walletId string) (data []entity.WalletBalance, err error)
	Deposit(ctx context.Context, walletId string, money entity.Money, note string) (err error)
	Withdraw(ctx context.Context, walletId string, money entity.Money, note string) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money, note string) (err error)
	GetBalance(ctx context.Context, walletId string, currency string) (decimal.Decimal, error)

	// history
	GetWalletBalanceHistory(ctx context.Context, walletId string) ([]entity.WalletBalanceHistory, error)
	GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]entity.WalletBalanceHistory, error)
}
