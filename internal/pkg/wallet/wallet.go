package wallet

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/structs"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Deposit(walletId string, money structs.Money) (err *t7Error.Error) {
	log.Debug("deposit wallet: ", walletId)

	data := db.DepositData{
		DepositReq: db.DepositReq{
			WalletId: walletId,
			Money:    money,
		},
		DepositId: uuid.New().String(),
	}
	if dbErr := db.New().Deposit(data); dbErr != nil {
		log.Error("fail to deposit: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func Withdraw(walletId string, money structs.Money) (err *t7Error.Error) {
	log.Debug("withdraw wallet: ", walletId)

	data := db.WithdrawData{
		WithdrawReq: db.WithdrawReq{
			WalletId: walletId,
			Money:    money,
		},
		WithdrawId: uuid.New().String(),
	}
	if dbErr := db.New().Withdraw(data); dbErr != nil {
		log.Error("fail to withdraw: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}
