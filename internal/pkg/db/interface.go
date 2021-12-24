package db

import (
	"github.com/Template7/backend/internal/pkg/transaction"
	"github.com/Template7/common/structs"
	"time"
)

type ClientInterface interface {
	// admin
	//GetAdmin() (data structs.Admin, err error)
	//GetUser(filter GetUserFilter) (users []structs.User, err error)
	initDb()

	// user
	GetUserById(userId string) (data structs.User, err error)
	GetUserBasicInfo(userId string) (data structs.UserInfo, err error)

	CreateUser(user structs.User) (err error)
	UpdateUserBasicInfo(userId string, info structs.UserInfo) (err error)

	// wallet
	GetWallet(userId string) (data structs.WalletData, err error)
	Deposit(walletId string, money structs.Money) (err error)
	Withdraw(walletId string, money structs.Money) (err error)
	Transfer(t TransactionData) (err error)
	GetTransactions(userId string) (data []string, err error)

	// transactionHistory
	//GetTransactions(filter string) (data []string, err error)
}

type TransactionData struct {
	transaction.RequestData
	TransactionId string    `json:"transaction_id" bson:"transaction_id"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}
