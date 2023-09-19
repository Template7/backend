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

	// wallet
	GetUserWallets(ctx context.Context, userId string) (data []entity.Wallet)
	Deposit(ctx context.Context, walletId string, money entity.Money) (err error)
	Withdraw(ctx context.Context, walletId string, money entity.Money) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error)
}
