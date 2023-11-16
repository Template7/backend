package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/auth"
	middleware "github.com/Template7/backend/internal/route/middleWare"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// NativeLogin
// @Summary Native login
// @Tags v1,login
// @version 1.0
// @Param request body types.HttpLoginReq true "Request"
// @produce json
// @Success 200 {object} types.HttpLoginResp "Response"
// @failure 400 {object} types.HttpRespError
// @Router /api/v1/login/native [post]
func NativeLogin(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle native login")

	var req types.HttpLoginReq
	if err := c.ShouldBindJSON(req); err != nil {
		log.WithError(err).Warn("invalid body")
		c.JSON(http.StatusBadRequest, types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      int(t7Error.InvalidBody.Code),
			Message:   t7Error.InvalidBody.Message,
		})
		return
	}

	token, err := auth.New().Login(c, req.Username, req.Password)
	if err != nil {
		defer c.Abort()
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

	lr := types.HttpLoginResp{
		HttpRespBase: types.HttpRespBase{
			RequestId: c.GetHeader(middleware.HeaderRequestId),
			Code:      types.HttpRespCodeOk,
			Message:   types.HttpRespMsgOk,
		},
	}
	lr.Data.Token = token
	c.JSON(http.StatusOK, lr)
}
