package types

type HttpGetWalletResp struct {
	HttpRespBase
	Data HttpGetUserWalletsRespData `json:"data"`
}

type HttpWalletWithdrawReq struct {
	Currency string `json:"currency" binding:"required;oneof=usd,utd" example:"usd"`
	Amount   uint32 `json:"amount" binding:"required" example:"100"`
}
