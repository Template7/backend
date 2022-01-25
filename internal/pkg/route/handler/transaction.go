package handler

import (
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/transaction"
	"github.com/Template7/backend/pkg/apiBody"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type transactionResp struct {
	TransactionId string `json:"transaction_id"`
}

func MakeTransfer(c *gin.Context) {
	log.WithContext(c).Debug("handle make transfer")

	var data apiBody.TransactionReq
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Warn("invalid data: ", err.Error())
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody)
		return
	}
	if err := validator.New().Struct(data); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}
	transactionId, err := transaction.Transfer(data)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	resp := transactionResp{
		TransactionId: transactionId,
	}
	c.JSON(http.StatusOK, resp)
}
