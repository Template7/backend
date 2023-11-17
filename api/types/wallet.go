package types

type HttpGetWalletResp struct {
	HttpRespBase
	Data HttpGetUserWalletsRespData `json:"data"`
}
