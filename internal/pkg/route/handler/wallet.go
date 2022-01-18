package handler

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/wallet"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Deposit(c *gin.Context) {
	log.Debug("handle deposit")

	var req db.DepositReq
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

func Withdraw(c *gin.Context) {
	log.Debug("handle withdraw")

	var req db.WithdrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Warn("invalid body: ", err.Error())
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	if err := wallet.Withdraw(req.WalletId, req.Money); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
}
