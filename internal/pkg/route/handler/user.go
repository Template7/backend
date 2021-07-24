package handler

import (
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/user"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetInfo
// @Summary Get user Info
// @Tags v1,user
// @version 1.0
// @Success 200 {object} collection.User
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Param UserId path string true "User ID"
// @Param Authorization header string true "Access token"
// @Router /api/v1/users/{UserId} [get]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func GetInfo(c *gin.Context) {
	log.Debug("handle get info")

	userInfo, err := user.GetInfo(c.Param("user-id"))
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}

// CreateUser
// @Summary Create user
// @Tags v1,user,admin
// @version 1.0
// @produce json
// @Param Authorization header string true "Access token"
// @Param userData body collection.User true "User data"
// @Success 204
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Router /admin/v1/user [post]
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
// @Param Authorization header string true "Access token"
// @Param user body collection.UserInfo true "User basic info"
// @Success 200
// @failure 400 {object} t7Error.Error
// @failure 401 {object} t7Error.Error
// @Param UserId path string true "User ID"
// @Router /api/v1/users/{UserId} [put]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func UpdateUser(c *gin.Context) {
	log.Debug("handle update user")

	userId := c.Param("user-id")

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
// @Router /admin/v1/users [delete]
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
