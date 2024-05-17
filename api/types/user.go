package types

import (
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
)

type HttpUserInfoResp struct {
	HttpRespBase
	Data HttpUserInfoRespData `json:"data"`
}

type HttpUserInfoRespData struct {
	UserId   string `json:"userId" example:"userId001"`
	Role     string `json:"role" example:"user"`
	Status   string `json:"status" example:"activated"`
	Nickname string `json:"nickname" example:"example"`
	Email    string `json:"email" example:"example@email.com"`
}

type HttpCreateUserReq struct {
	Username string `json:"username" binding:"required" example:"username"`
	Password string `json:"password" binding:"required" example:"password"`
	Role     string `json:"role" binding:"oneof=admin operator user" example:"user"`
	Nickname string `json:"nickname" example:"nickname"`
	Email    string `json:"email" example:"example@email.com"`
}

func (h HttpCreateUserReq) ToProto() *userV1.CreateUserRequest {
	return &userV1.CreateUserRequest{
		Username: h.Username,
		Password: h.Password,
		Role:     authV1.Role(authV1.Role_value[h.Role]),
		Nickname: h.Nickname,
		Email:    h.Email,
	}
}

type HttpCreateUserResp struct {
	HttpRespBase
	Data HttpCreateUserRespData `json:"data"`
}

type HttpCreateUserRespData struct {
	UserId         string `json:"userId" example:"f47ac10b-58cc-0372-8567-0e02b2c3d479"`
	ActivationCode string `json:"activationCode" example:"0b25b3a8-2e85-452b-8d54-527f7e68a875"`
}

type HttpUpdateUserInfoReq struct {
	Nickname string `json:"nickname" binding:"required" example:"nickname"`
}

type HttpGetUserWalletsResp struct {
	HttpRespBase
	Data []HttpGetUserWalletsRespData `json:"data"`
}

type HttpGetUserWalletsRespData struct {
	Id       string                              `json:"id" example:"af68a360-d035-469c-8ae9-a8640c2ffd19"`
	Balances []HttpGetUserWalletsRespDataBalance `json:"balances"`
}

type HttpGetUserWalletsRespDataBalance struct {
	Currency string `json:"currency" example:"usd"`
	Amount   string `json:"amount" example:"100"`
}

type HttpActivateUserReq struct {
	ActivationCode string `json:"activationCode"`
}

type HttpActivateUserResp struct {
	HttpRespBase
	Data HttpActivateUserRespData `json:"data"`
}

type HttpActivateUserRespData struct {
	Success bool `json:"success"`
}
