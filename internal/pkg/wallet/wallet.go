package wallet

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/structs"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Deposit(walletId string, money structs.Money) (err *t7Error.Error) {
	log.Debug("deposit wallet: ", walletId)

	// TODO: deposit id
	if dbErr := db.New().Deposit(walletId, money); dbErr != nil {
		log.Error("fail to deposit: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}
