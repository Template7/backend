package types

type HttpGetWalletResp struct {
	HttpRespBase
	Data HttpGetUserWalletsRespData `json:"data"`
}

type HttpWalletWithdrawReq struct {
	Currency string `json:"currency" binding:"required;oneof=usd,utd" example:"usd"`
	Amount   uint32 `json:"amount" binding:"required" example:"100"`
}

type HttpWalletDepositReq struct {
	Currency string `json:"currency" binding:"required;oneof=usd,utd" example:"usd"`
	Amount   uint32 `json:"amount" binding:"required" example:"100"`
}

type HttpTransferMoneyReq struct {
	FromWalletId string `json:"fromWalletId" binding:"required" example:"4eb4a439-af97-46e4-8a0c-6d568281c43a"`
	ToWalletId   string `json:"toWalletId" binding:"required" example:"d53ce74f-5f74-4c78-b3ca-1e1d2f7fa43d"`
	Currency     string `json:"currency" binding:"required;oneof=usd,utd" example:"usd"`
	Amount       uint32 `json:"amount" binding:"required" example:"100"`
}
