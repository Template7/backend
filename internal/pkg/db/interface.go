package db

import (
	"fmt"
	"github.com/Template7/common/structs"
	"time"
)

type ClientInterface interface {

	// user
	GetUserById(userId string) (data structs.User, err error)
	GetUserByMobile(mobile string) (data structs.User, err error)
	GetUserBasicInfo(userId string) (data structs.UserInfo, err error)
	GetFbUser(fbUserId string) (data structs.User, err error)

	CreateUser(user structs.User) (err error)
	UpdateUserBasicInfo(userId string, info structs.UserInfo) (err error)
	UpdateUserStatus(userId string, status structs.UserStatus) (err error)

	// wallet
	GetWallet(userId string) (data structs.WalletData, err error)
	Deposit(walletId string, money structs.Money) (err error)
	Withdraw(walletId string, money structs.Money) (err error)
	Transfer(t TransactionData) (err error)
	GetTransactions(userId string) (data []TransactionData, err error)
}

type TransactionData struct {
	TransactionReq `json:",inline" bson:",inline" validate:"dive"`
	TransactionId  string    `json:"transaction_id" bson:"transaction_id" validate:"uuid"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

// from api
type TransactionReq struct {
	FromWalletId  string `json:"from_wallet_id" bson:"from_wallet_id" validate:"uuid"`
	ToWalletId    string `json:"to_wallet_id" bson:"to_wallet_id" validate:"uuid"`
	Note          string `json:"note" bson:"note"`
	structs.Money `json:",inline" bson:",inline" validate:"required,dive"`
}

func (r TransactionReq) String() string {
	return fmt.Sprintf("from %s to %s, %d %s", r.FromWalletId, r.ToWalletId, r.Amount, r.Unit)
}
