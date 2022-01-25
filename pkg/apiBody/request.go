package apiBody

import "github.com/Template7/common/structs"

type CreateUserReq struct {
	Mobile string `json:"mobile" bson:"mobile" example:"+886987654321"` // +886987654321
	Email  string `json:"email" bson:"email" example:"username@mail.com"`
}

type SmsReq struct {
	Mobile string `json:"mobile" binding:"required" example:"+886987654321"`
}

type UserInfoResp struct {
	UserInfo   structs.UserInfo   `json:",inline"`
	WalletData structs.WalletData `json:"wallet_data"`
}
