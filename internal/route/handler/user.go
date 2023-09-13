package handler

import (
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

// GetUserInfo
// @Summary Get user Info
// @Tags V1,User
// @version 1.0
// @Success 200 {object} structs.User
// @failure 400 {object} t7Error.Error
func GetUserInfo(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get user info")

	// TODO: implementation
	return
}

// CreateUser
// @Summary Create user
// @Tags V1,User
// @version 1.0
func CreateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle create user")

	// TODO: implementation
	return
}

// UpdateUser
// @Summary Update user
// @Tags V1,User
// @version 1.0
func UpdateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle update user")

	// TODO: implementation
	return
}
