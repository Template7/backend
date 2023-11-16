package types

const (
	HttpRespCodeOk = 3000
	HttpRespMsgOk  = "ok"
)

// HttpRespError
// for swagger only
type HttpRespError struct {
	RequestId string `json:"requestId" example:"27c0a70e-59ab-4a94-872c-5f014aaa047f"`
	Code      int    `json:"code" example:"1024"`
	Message   string `json:"message" example:"token unauthorized"`
}

type HttpRespBase struct {
	RequestId string `json:"requestId" example:"b8974256-1f17-477f-8638-c7ebbac656d7"`
	Code      int    `json:"code" example:"3000"`
	Message   string `json:"message" example:"ok"`
}
