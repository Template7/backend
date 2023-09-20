package handler

import (
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

func MakeTransfer(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.WithContext(c).Debug("handle make transfer")

	// TODO: implementation
	return
}
