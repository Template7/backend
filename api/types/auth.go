package types

type HttpLoginReq struct {
	Username string `json:"username" binding:"required" example:"username"`
	Password string `json:"password" binding:"required" example:"password"`
}
type HttpLoginResp struct {
	HttpRespBase
	Data struct {
		Token string `json:"token" example:"70596484-67d3-46bd-94bf-08f7c9fb7ac1"`
	} `json:"data"`
}
