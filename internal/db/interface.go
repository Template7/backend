package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/structs"
)

type Client interface {
	// user
	CreateUser(ctx context.Context, data entity.User) (err error)
	GetUser(ctx context.Context, username string) (data entity.User, err error)
	GetUserById(ctx context.Context, userId string) (data structs.User, err error)

	// wallet
	GetWallet(ctx context.Context, userId string) (data entity.Wallet, err error)
	Deposit(ctx context.Context, walletId string, money entity.Money) (err error)
	Withdraw(ctx context.Context, walletId string, money entity.Money) (err error)
	Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error)
}
