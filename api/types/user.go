package types

import (
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
)

type HttpUserInfoResp struct {
	HttpRespBase
	Data struct {
		UserId   string `json:"userId" example:"userId001"`
		Role     string `json:"role" example:"user"`
		Status   string `json:"status" example:"activated"`
		Nickname string `json:"nickname" example:"example"`
		Email    string `json:"email" example:"example@email.com"`
	} `json:"data"`
}

type HttpCreateUserReq struct {
	Username string `json:"username" binding:"required" example:"username"`
	Password string `json:"password" binding:"required" example:"password"`
	Role     string `json:"role" binding:"oneof=admin,operator,user" example:"user"`
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

type HttpUpdateUserInfoReq struct {
	Nickname string `json:"nickname" binding:"required" example:"nickname"`
}
