package handler

import (
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/backend/internal/pkg/user"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetInfo
// @Summary Get user Info
// @Tags V1,User
// @version 1.0
// @Success 200 {object} structs.User
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

	userInfo, err := user.GetInfo(c.Param("userId"))
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}

type createUserResp struct {
	UserId string `json:"user_id"`
}

// CreateUser
// @Summary Create user
// @Tags V1,User,Admin
// @version 1.0
// @produce json
// @Param Authorization header string true "Access token"
// @Param userData body user.CreateUserReq true "User data"
// @Success 200 {object} createUserResp "User object"
// @failure 400 {object} t7Error.Error "Error object"
// @failure 401 {object} t7Error.Error "Error object"
// @Router /admin/v1/user [post]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func CreateUser(c *gin.Context) {
	log.Debug("handle create user")

	var data user.CreateUserReq
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	userId, err := user.CreateUser(data)
	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, createUserResp{
		UserId: userId,
	})
	return
}

// UpdateUser
// @Summary Update user
// @Tags V1,User
// @version 1.0
// @produce json
// @Param Authorization header string true "Access token"
// @Param user body structs.UserInfo true "User basic info"
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

	//userId := c.Param("user-id")
	//
	//var userData structs.UserInfo
	//if err := c.BindJSON(&userData); err != nil {
	//	c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
	//	return
	//}
	//
	//if err := user.UpdateBasicInfo(userId, userData); err != nil {
	//	c.JSON(err.GetStatus(), err)
	//	return
	//}
	//
	//log.Debug("user updated: ", userId)
	//c.JSON(http.StatusNoContent, nil)
	return
}

func UpdateLoginClient(c *gin.Context) {
	log.Error("handle update login client")

	//var loginClient structs.LoginInfo
	//if err := c.BindJSON(&loginClient); err != nil {
	//	c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
	//	return
	//}
	//if err := user.UpdateLoginClient(c.Param("userId"), loginClient); err != nil {
	//	c.JSON(err.GetStatus(), err)
	//	return
	//}
	//
	//c.JSON(http.StatusNoContent, nil)
	return
}
