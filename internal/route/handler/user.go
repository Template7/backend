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
	}
	resp.Data.UserId = info.GetUserId()
	resp.Data.Role = authV1.Role_name[int32(info.GetRole())]
	resp.Data.Status = authV1.AccountStatus_name[int32(info.GetStatus())]
	resp.Data.Nickname = info.GetNickname()
	resp.Data.Email = info.GetEmail()
	c.JSON(http.StatusOK, resp)
}

// CreateUser
// @Summary Create user
// @Tags V1,User
// @version 1.0
// @Param request body types.HttpCreateUserReq true "Request"
// @produce json
// @Success 200 {object} types.HttpLoginResp "Response"
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

	if err := auth.New().CreateUser(c, req.ToProto()); err != nil {
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

	c.JSON(http.StatusNoContent, types.HttpLoginResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
	})
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
		Nickname: req.Nickname,
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
