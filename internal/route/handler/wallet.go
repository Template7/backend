package handler

import (
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/wallet"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWallet(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get wallet")

	data, err := wallet.New().GetWallet(c, c.Param("walletId"))
	if err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to update user info")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}

	c.JSON(http.StatusOK, &data)
}

func Deposit(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle deposit")

	// TODO: implementation
	return
}

func Withdraw(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle withdraw")

	// TODO: implementation
	return
}
