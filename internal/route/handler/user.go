package handler

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db/entity"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
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

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}

	info, err := user.New().GetInfo(c, userId)
	if err != nil {
		log.WithError(err).Error("fail to get user info")
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

	c.JSON(http.StatusOK, info)
}

// CreateUser
// @Summary Create user
// @Tags V1,User
// @version 1.0
func CreateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle create user")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req userV1.CreateUserRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	if err := auth.New().CreateUser(c, &req); err != nil {
		log.WithError(err).Error("fail to create user")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateUser
// @Summary Update user
// @Tags V1,User
// @version 1.0
func UpdateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle update user")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req userV1.UpdateUserInfoRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}

	data := entity.UserInfo{
		NickName: req.Nickname,
	}
	if err := user.New().UpdateInfo(c, userId, data); err != nil {
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
	c.JSON(http.StatusNoContent, nil)
}

func GetUserWallets(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get user wallets")

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusForbidden, t7Error.InvalidToken)
		return
	}

	c.JSON(http.StatusOK, userV1.GetUserWalletResponse{
		Wallets: user.New().GetUserWallets(c, userId),
	})
}
