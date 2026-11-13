package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/common/logger"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service auth.Auth
	log     *logger.Logger
}

func NewAuthController(service auth.Auth, log *logger.Logger) *AuthController {
	return &AuthController{
		service: service,
		log:     log.WithService("authController"),
	}
}

// NativeLogin
// @Summary Native login
// @Tags V1,Login
// @version 1.0
// @Param request body types.HttpLoginReq true "Request"
// @produce json
// @Success 200 {object} types.HttpLoginResp "Response"
// @failure 400 {object} types.HttpRespError
// @failure 401 {object} types.HttpRespError "invalid user name or password"
// @failure 429 {object} types.HttpRespError "too many failed attempts, try again later"
// @Router /api/v1/login/native [post]
func (a *AuthController) NativeLogin(c *gin.Context) {
	var req types.HttpLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := a.service.Login(c.Request.Context(), req.Username, req.Password, c.ClientIP())
	if err != nil {
		response.Error(c, err)
		return
	}

	resp := types.HttpLoginResp{HttpRespBase: response.Ok(c)}
	resp.Data.Token = token
	response.JSON(c, resp)
}
