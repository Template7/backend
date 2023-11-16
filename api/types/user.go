package types

import (
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
)

type HttpUserInfoResp struct {
	HttpRespBase
	Data struct {
		UserId   string               `json:"userId" example:"userId001"`
		Role     authV1.Role          `json:"role" example:"0"`
		Status   authV1.AccountStatus `json:"status" example:"0"`
		Nickname string               `json:"nickname" example:"example"`
		Email    string               `json:"email" example:"example@email.com"`
	} `json:"data"`
}

type HttpCreateUserReq struct {
	Username string      `json:"username" binding:"required" example:"username"`
	Password string      `json:"password" binding:"required" example:"password"`
	Role     authV1.Role `json:"role" binding:"required" example:"0"`
	Nickname string      `json:"nickname" example:"nickname"`
	Email    string      `json:"email" example:"example@email.com"`
}

func (h HttpCreateUserReq) ToProto() *userV1.CreateUserRequest {
	return &userV1.CreateUserRequest{
		Username: h.Username,
		Password: h.Password,
		Role:     h.Role,
		Nickname: h.Nickname,
		Email:    h.Email,
	}
}
