package transaction

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/pkg/apiBody"
	"github.com/Template7/common/logger"
	"github.com/google/uuid"
	"net/http"
)

var (
	log = logger.GetLogger()
)

func Transfer(data apiBody.TransactionReq) (transactionId string, err *t7Error.Error) {
	log.Debug("handle transfer: ", data.String())

	// TODO: verify from_wallet_id and token

	transferData := db.TransactionData{
		TransactionReq: data,
		TransactionId:  uuid.New().String(),
	}
	if dbErr := db.New().Transfer(transferData); dbErr != nil {
		log.Error("fail to make transfer: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}

	log.Debug("finish transfer: ", data.String())
	transactionId = transferData.TransactionId
	return
}
