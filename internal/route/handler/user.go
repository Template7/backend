package handler

import (
	"github.com/Template7/backend/api/types"
	"github.com/Template7/backend/internal/route/middleware"
	"github.com/Template7/backend/internal/route/response"
	"github.com/Template7/backend/internal/t7Error"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService user.Service
	log         *logger.Logger
}

func NewUserController(service user.Service, log *logger.Logger) *UserController {
	return &UserController{
		userService: service,
		log:         log.WithService("userController"),
	}
}

// currentUserId reads the authenticated user id; on failure the response is already written.
func currentUserId(c *gin.Context) (string, bool) {
	userId, ok := middleware.CurrentUserId(c)
	if !ok {
		response.Error(c, t7Error.InvalidToken.WithDetail("no user id in context"))
	}
	return userId, ok
}

// GetUserInfo
// @Summary Get user Info
// @Tags V1,User
// @version 1.0
// @Success 200 {object} types.HttpUserInfoResp "Response"
// @failure 400 {object} types.HttpRespError
// @Security BearerAuth
// @Router /api/v1/user/info [get]
func (u *UserController) GetUserInfo(c *gin.Context) {
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	info, err := u.userService.GetInfo(c, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.JSON(c, types.HttpUserInfoResp{
		HttpRespBase: response.Ok(c),
		Data: types.HttpUserInfoRespData{
			UserId:   info.GetUserId(),
			Role:     authV1.Role_name[int32(info.GetRole())],
			Status:   authV1.AccountStatus_name[int32(info.GetStatus())],
			Nickname: info.GetNickname(),
			Email:    info.GetEmail(),
		},
	})
}

// CreateUser
// @Summary Create user
// @Description Creates the user with a default wallet and returns a one-time activation code (valid 24h).
// @Tags V1,User
// @version 1.0
// @Param request body types.HttpCreateUserReq true "Request"
// @produce json
// @Success 200 {object} types.HttpCreateUserResp "Response"
// @failure 400 {object} types.HttpRespError
// @Security BearerAuth
// @Router /admin/v1/user [post]
func (u *UserController) CreateUser(c *gin.Context) {
	var req types.HttpCreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userId, actCode, err := u.userService.CreateUser(c, req.ToProto())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.JSON(c, types.HttpCreateUserResp{
		HttpRespBase: response.Ok(c),
		Data: types.HttpCreateUserRespData{
			UserId:         userId,
			ActivationCode: actCode,
		},
	})
}

// ActivateUser
// @Summary Activate user
// @Tags V1,User
// @version 1.0
// @Param request body types.HttpActivateUserReq true "Request"
// @Param userId path string true "User ID"
// @produce json
// @Success 200 {object} types.HttpActivateUserResp "Response"
// @failure 400 {object} types.HttpRespError "invalid or expired activation code"
// @Security BearerAuth
// @Router /admin/v1/users/{userId}/activate [post]
func (u *UserController) ActivateUser(c *gin.Context) {
	var req types.HttpActivateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := u.userService.ActivateUser(c, c.Param("userId"), req.ActivationCode); err != nil {
		response.Error(c, err)
		return
	}

	response.JSON(c, types.HttpActivateUserResp{
		HttpRespBase: response.Ok(c),
		Data:         types.HttpActivateUserRespData{Success: true},
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
// @Security BearerAuth
// @Router /api/v1/user/info [put]
func (u *UserController) UpdateUser(c *gin.Context) {
	var req types.HttpUpdateUserInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userId, ok := currentUserId(c)
	if !ok {
		return
	}

	if err := u.userService.UpdateInfo(c, userId, models.UserInfo{Nickname: req.Nickname}); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, response.Ok(c))
}

// DeleteUser
// @Summary Delete user
// @Tags V1,User
// @version 1.0
// @produce json
// @Success 200 {object} types.HttpRespBase "Response"
// @failure 404 {object} types.HttpRespError "user not found"
// @Security BearerAuth
// @Param userId path string true "User ID"
// @Router /admin/v1/users/{userId} [delete]
func (u *UserController) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		response.BadRequest(c, "empty user id")
		return
	}

	if err := u.userService.DeleteUser(c, userId); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, response.Ok(c))
}

func toBalances(bls []*v1.Balance) []types.HttpGetUserWalletsRespDataBalance {
	r := make([]types.HttpGetUserWalletsRespDataBalance, len(bls))
	for i, bl := range bls {
		r[i] = types.HttpGetUserWalletsRespDataBalance{
			Currency: bl.GetCurrency().String(),
			Amount:   bl.GetAmount(),
		}
	}
	return r
}
