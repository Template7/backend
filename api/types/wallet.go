package types

type HttpGetWalletResp struct {
	HttpRespBase
	Data HttpGetUserWalletsRespData `json:"data"`
}

type HttpWalletWithdrawReq struct {
	Currency string `json:"currency" binding:"required,oneof=usd ntd cny jpy" example:"usd"`
	Amount   uint32 `json:"amount" binding:"required" example:"100"`
	Note     string `json:"note"`
}

type HttpWalletDepositReq struct {
	Currency string `json:"currency" binding:"required,oneof=usd ntd cny jpy" example:"usd"`
	Amount   uint32 `json:"amount" binding:"required" example:"100"`
	Note     string `json:"note"`
}

type HttpTransferMoneyReq struct {
	FromWalletId string `json:"fromWalletId" binding:"required" example:"4eb4a439-af97-46e4-8a0c-6d568281c43a"`
	ToWalletId   string `json:"toWalletId" binding:"required" example:"d53ce74f-5f74-4c78-b3ca-1e1d2f7fa43d"`
	Currency     string `json:"currency" binding:"required,oneof=usd ntd cny jpy" example:"usd"`
	Amount       uint32 `json:"amount" binding:"required" example:"100"`
	Note         string `json:"note"`
}

type HttpGetWalletBalanceHistoryResp struct {
	HttpRespBase
	Data []HttpGetWalletBalanceHistoryData `json:"data"`
}

type HttpGetWalletBalanceHistoryData struct {
	Direction     string `json:"direction" example:"deposit"` // one of deposit, withdraw, transferIn, transferOut
	Currency      string `json:"currency" example:"ntd"`
	Amount        string `json:"amount" example:"100"`
	BalanceBefore string `json:"balanceBefore" example:"10"`
	BalanceAfter  string `json:"balanceAfter" example:"110"`
	Timestamp     int64  `json:"timestamp"`
	Note          string `json:"note"`
}

type HttpGetWalletBalanceHistoryByCurrencyResp struct {
	HttpRespBase
	Data []HttpGetWalletBalanceHistoryByCurrencyData `json:"data"`
}

type HttpGetWalletBalanceHistoryByCurrencyData struct {
	Direction     string `json:"direction" example:"deposit"` // one of deposit, withdraw, transferIn, transferOut
	Amount        string `json:"amount" example:"100"`
	BalanceBefore string `json:"balanceBefore" example:"10"`
	BalanceAfter  string `json:"balanceAfter" example:"110"`
	Timestamp     int64  `json:"timestamp"`
	Note          string `json:"note"`
}
