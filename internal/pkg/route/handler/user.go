package handler

import (
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/route/middle_ware"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/user"

	"backend/internal/pkg/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetInfo
// @Summary Get user Info
// @Tags user
// @version 1.0
// @Success 200 {object} collection.User
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /app/v1/users [get]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func GetInfo(c *gin.Context) {
	log.Debug("handle get info")

	userInfo, err := user.GetInfo(c.Param("userId"))
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}

func RefreshToken(c *gin.Context) {
	log.Debug("handle refresh token")

	userId := c.Param("user-id")
	token, err := user.GenToken(userId)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, token)
	return
}

func SignOut(c *gin.Context) {
	log.Debug("handle user sign out")

	claims := c.Keys["claims"]
	tokenClaims := claims.(*middle_ware.UserTokenClaims)

	if err := user.SignOut(util.ParseBearerToken(c.GetHeader("Authorization")), tokenClaims.ExpiresAt); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// CreateUser
// @Summary Create user
// @Tags user
// @version 1.0
// @produce json
// @Param userData body collection.User true "The whole user data"
// @Success 200
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /app/v1/users [put]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func CreateUser(c *gin.Context) {
	log.Debug("handle create user")

	var userData collection.User
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	if _, err := user.CreateUser(userData); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// UpdateUser
// @Summary Update user
// @Tags user
// @version 1.0
// @produce json
// @Param user body collection.User true "The whole user data"
// @Success 200
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /app/v1/users [put]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func UpdateUser(c *gin.Context) {
	log.Debug("handle update user")

	userId := c.Param("userId")

	var userData collection.UserInfo
	if err := c.BindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	if err := user.UpdateBasicInfo(userId, userData); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	log.Debug("user updated: ", userId)
	c.JSON(http.StatusNoContent, nil)
	return
}

func UpdateLoginClient(c *gin.Context) {
	log.Error("handle update login client")

	var loginClient collection.LoginInfo
	if err := c.BindJSON(&loginClient); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}
	if err := user.UpdateLoginClient(c.Param("userId"), loginClient); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
	return
}

// DeleteUser
// @Summary Delete user
// @Tags user
// @version 1.0
// @produce json
// @Success 200
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /app/v1/users [delete]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func DeleteUser(c *gin.Context) {
	log.Debug("handle delete user")
	userId := c.Param("userId")

	if err := user.DeleteUser(userId); err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
	return
}
