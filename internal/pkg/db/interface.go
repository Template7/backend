package db

import (
	"context"
	"github.com/Template7/backend/internal/pkg/db/entity"
	"github.com/Template7/backend/pkg/apiBody"
	"github.com/Template7/common/structs"
	"time"
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

type TransactionData struct {
	apiBody.TransactionReq `json:",inline" bson:",inline" validate:"dive"`
	TransactionId          string    `json:"transaction_id" bson:"transaction_id" validate:"uuid"`
	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
}

type DepositData struct {
	apiBody.DepositReq `json:",inline" bson:",inline" validate:"dive"`
	DepositId          string    `json:"deposit_id" bson:"deposit_id" validate:"uuid"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
}

type WithdrawData struct {
	WithdrawReq `json:",inline" bson:",inline" validate:"dive"`
	WithdrawId  string    `json:"withdraw_id" bson:"withdraw_id" validate:"uuid"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type WithdrawReq struct {
	Target   string        `json:"target" bson:"target"`
	WalletId string        `json:"wallet_id" bson:"wallet_id" validate:"uuid"`
	Money    structs.Money `json:"money" bson:"money" validate:"dive"`
}
