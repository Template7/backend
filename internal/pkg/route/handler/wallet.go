package handler

import (
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/wallet"
	"github.com/Template7/common/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type depositReq struct {
	WalletId string        `json:"wallet_id" bson:"wallet_id" validate:"uuid"`
	Money    structs.Money `json:"money" bson:"money" validate:"dive"`
}

func Deposit(c *gin.Context) {
	log.Debug("handle deposit")

	var req depositReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Warn("invalid body: ", err.Error())
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	if err := wallet.Deposit(req.WalletId, req.Money); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
}
