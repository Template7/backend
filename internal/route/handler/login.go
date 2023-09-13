package handler

import (
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

// NativeLogin
// @Summary Native login
// @Tags v1,login
// @version 1.0
// @Param request body v1.LoginRequest true "Request"
// @produce json
// @Success 200 {object} v1.LoginResponse "Response"
// @failure 400 {object} t7Error.Error
// @Router /api/v1/login/native [post]
func NativeLogin(c *gin.Context) {
	log := logger.New().WithContext(c)
	log.Debug("handle native login")

	defer c.Request.Body.Close()
	bd, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("fail to read resp body")
		c.JSON(http.StatusBadRequest, t7Error.InvalidBody.WithDetail(err.Error()))
		return
	}

	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
	var req v1.LoginRequest
	if err := unmarshaler.Unmarshal(bd, &req); err != nil {
		log.WithError(err).With("resp", string(bd)).Error("fail to decode resp data")
		c.JSON(http.StatusBadRequest, t7Error.DecodeFail.WithDetail(err.Error()))
		return
	}

	token, err := auth.New().Login(c, req.Username, req.Password)
	if err != nil {
		defer c.Abort()
		t7Err, ok := t7Error.ToT7Error(err)
		if !ok {
			log.WithError(err).Error("unknown error")
			c.JSON(http.StatusForbidden, t7Error.InvalidToken)
			return
		}
		c.JSON(t7Err.GetStatus(), t7Err)
		return
	}

	c.JSON(http.StatusOK, v1.LoginResponse{
		Token: token,
	})
}
