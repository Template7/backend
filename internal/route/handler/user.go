package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db/entity"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetUserInfo
// @Summary Get user Info
// @Tags V1,User
// @version 1.0
// @Success 200 {object} types.HttpUserInfoResp "Response"
// @failure 400 {object} types.HttpRespError
// @Router /api/v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get user info")

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}

	info, err := user.New().GetInfo(c, userId)
	if err != nil {
		log.WithError(err).Error("fail to get user info")
		defer c.Abort()
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	resp := types.HttpUserInfoResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: types.HttpUserInfoRespData{
			UserId:   info.GetUserId(),
			Role:     authV1.Role_name[int32(info.GetRole())],
			Status:   authV1.AccountStatus_name[int32(info.GetStatus())],
			Nickname: info.GetNickname(),
			Email:    info.GetEmail(),
		},
	}
	c.JSON(http.StatusOK, resp)
}

// CreateUser
// @Summary Create user
// @Tags V1,User
// @version 1.0
// @Param request body types.HttpCreateUserReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Router /admin/v1/user [post]
func CreateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle create user")

	var req types.HttpCreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	actCode, err := auth.New().CreateUser(c, req.ToProto())
	if err != nil {
		log.WithError(err).Error("fail to create user")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpCreateUserResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: types.HttpCreateUserRespData{
			ActivationCode: actCode,
		},
	})
}

// UpdateUser
// @Summary Update user info
// @Tags V1,User
// @version 1.0
// @Param request body types.HttpUpdateUserInfoReq true "Request"
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Router /api/v1/user/info [put]
func UpdateUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle update user")

	var req types.HttpUpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}

	data := entity.UserInfo{
		Nickname: req.Nickname,
	}
	if err := user.New().UpdateInfo(c, userId, data); err != nil {
		defer c.Abort()
		log.WithError(err).Error("fail to update user info")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpRespBase{
		RequestId: c.GetHeader(middleware.HeaderRequestId),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	})
}

// GetUserWallets
// @Summary Get user wallets
// @Tags V1,User,Wallet
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpGetUserWalletsResp "Response"
// @failure 400 {object} types.HttpRespError
// @Router /api/v1/user/wallets [get]
func GetUserWallets(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle get user wallets")

	uId, ok := c.Get(middleware.UserId)
	if !ok {
		log.Warn("no user id from the previous middleware")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}
	userId, ok := uId.(string)
	if !ok {
		log.With("uId", uId).Error("type assertion fail")
		c.JSON(http.StatusUnauthorized, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidToken.Code),
			Message:   t7Error.InvalidToken.Message,
		})
		return
	}

	uws := user.New().GetUserWallets(c, userId)
	rd := make([]types.HttpGetUserWalletsRespData, len(uws))
	for i, uw := range uws {
		rd[i] = types.HttpGetUserWalletsRespData{
			Id: uw.Id,
			Balances: func(bls []*v1.Balance) []types.HttpGetUserWalletsRespDataBalance {
				r := make([]types.HttpGetUserWalletsRespDataBalance, len(bls))
				for i, bl := range bls {
					r[i] = types.HttpGetUserWalletsRespDataBalance{
						Currency: bl.Currency.String(),
						Amount:   bl.Amount,
					}
				}
				return r
			}(uw.Balances),
		}
	}

	c.JSON(http.StatusOK, types.HttpGetUserWalletsResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
		Data: rd,
	})
}

// DeleteUser
// @Summary Delete user
// @Tags V1,User
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 400 {object} types.HttpRespError
// @Param userId path string true "User ID"
// @Router /admin/v1/users/{userId} [delete]
func DeleteUser(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle delete user")

	userId := c.Param("userId")
	if userId == "" {
		log.Warn("empty user id")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	if err := auth.New().DeleteUser(c, userId); err != nil {
		log.WithError(err).Error("fail to create user")
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusInternalServerError, types.HttpRespBase{
				RequestId: c.GetHeader(middleware.HeaderRequestId),
				Code:      int(t7Error.Unknown.Code),
				Message:   t7Error.Unknown.Message,
			})
			return
		}
		c.JSON(t7Err.GetStatus(), types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Err.Code),
			Message:   t7Err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, types.HttpRespBase{
		RequestId: c.GetHeader(middleware.HeaderRequestId),
		Code:      types.HttpRespCodeOk,
		Message:   types.HttpRespMsgOk,
	})
}
