package transaction

import "github.com/Template7/common/structs"

// from api
type RequestData struct {
	FromWalletId  string `json:"from_wallet_id" bson:"from_wallet_id"`
	ToWalletId    string `json:"to_wallet_id" bson:"to_wallet_id"`
	Note          string `json:"note" bson:"note"`
	structs.Money `json:",inline" bson:",inline"`
}


