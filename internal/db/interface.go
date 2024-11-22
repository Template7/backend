package db

import (
	"context"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/shopspring/decimal"
)

type Client interface {
	// user
	CreateUser(ctx context.Context, data models.User) (err error)
	GetUser(ctx context.Context, username string) (data models.User, err error)
	GetUserById(ctx context.Context, userId string) (data models.User, err error)
	UpdateUserInfo(ctx context.Context, userId string, info models.UserInfo) (err error)
	GetUserWallets(ctx context.Context, userId string) (data []models.UserWalletBalance)
	DeleteUser(ctx context.Context, userId string) (err error)
	SetUserStatus(ctx context.Context, userId string, status authV1.AccountStatus) (err error)

	// wallet
	GetWalletBalances(ctx context.Context, walletId string) (data []models.WalletBalance, err error)
	Deposit(ctx context.Context, walletId string, money models.Money, note string) (err error)
	Withdraw(ctx context.Context, walletId string, money models.Money, note string) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money models.Money, note string) (err error)
	GetBalance(ctx context.Context, walletId string, currency string) (decimal.Decimal, error)

	// history
	GetWalletBalanceHistory(ctx context.Context, walletId string) ([]models.WalletBalanceHistory, error)
	GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]models.WalletBalanceHistory, error)
}
