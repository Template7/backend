package wallet

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/logger"
	"github.com/google/uuid"
	"net/http"
)

var (
	log = logger.GetLogger()
)

func Deposit(req db.DepositReq) (depositId string, err *t7Error.Error) {
	log.Debug("deposit wallet: ", req.WalletId)

	data := db.DepositData{
		DepositReq: req,
		DepositId:  uuid.New().String(),
	}
	if dbErr := db.New().Deposit(data); dbErr != nil {
		log.Error("fail to deposit: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	depositId = data.DepositId
	return
}

func Withdraw(req db.WithdrawReq) (withdrawId string, err *t7Error.Error) {
	log.Debug("withdraw wallet: ", req.WalletId)

	data := db.WithdrawData{
		WithdrawReq: req,
		WithdrawId:  uuid.New().String(),
	}
	if dbErr := db.New().Withdraw(data); dbErr != nil {
		log.Error("fail to withdraw: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	withdrawId = data.WithdrawId
	return
}
