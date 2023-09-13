package handler

import (
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

func GetWallet(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get wallet")

	// TODO: implementation
	return
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
