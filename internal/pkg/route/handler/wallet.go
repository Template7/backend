package handler

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/wallet"
	"github.com/Template7/backend/pkg/apiBody"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func GetWallet(c *gin.Context) {

}

type depositResp struct {
	DepositId string `json:"deposit_id"`
}

func Deposit(c *gin.Context) {
	log.Debug("handle deposit")

	var req apiBody.DepositReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody)
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Warn("invalid body: ", err.Error())
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	depositId, err := wallet.Deposit(req)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, depositResp{DepositId: depositId})
}

type withdrawResp struct {
	WithdrawId string `json:"withdraw_id"`
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

	withDrawId, err := wallet.Withdraw(req)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, withdrawResp{WithdrawId: withDrawId})
}
