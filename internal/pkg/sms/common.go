package sms

type Request struct {
	Mobile string `json:"mobile" binding:"required"`
}

type Confirm struct {
	Mobile string `json:"mobile" binding:"required"`
	Code   string `json:"code" binding:"required"`
}
