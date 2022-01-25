package apiBody

import (
	"fmt"
	"github.com/Template7/common/structs"
)

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

type TransactionReq struct {
	FromWalletId  string `json:"from_wallet_id" bson:"from_wallet_id" validate:"uuid"`
	ToWalletId    string `json:"to_wallet_id" bson:"to_wallet_id" validate:"uuid"`
	Note          string `json:"note" bson:"note"`
	structs.Money `json:",inline" bson:",inline" validate:"required,dive"`
}

func (r TransactionReq) String() string {
	return fmt.Sprintf("from %s to %s, %d %s", r.FromWalletId, r.ToWalletId, r.Amount, r.Unit)
}
