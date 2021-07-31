package handler

import (
	"backend/internal/pkg/auth"
	"backend/internal/pkg/db/collection"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// RefreshToken
// @Summary Refresh access token
// @Tags V1,Token
// @version 1.0
// @Param token body collection.Token true "Token object"
// @Success 200 {object} collection.Token "Token object"
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Param UserId path string true "User ID"
// @Param Authorization header string true "Access token"
// @Router /api/v1/users/{UserId}/token [put]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func RefreshToken(c *gin.Context) {
	log.Debug("handle refresh token")

	var token collection.Token
	if err := c.BindJSON(&token); err != nil {
		log.Warn("invalid body: ", err.Error())
		c.JSON(http.StatusBadRequest, err)
		return
	}
	newToken, err := auth.RefreshToken(token)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, newToken)
	return
}
