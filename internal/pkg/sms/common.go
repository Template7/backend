package sms

type Request struct {
	Mobile   string `json:"mobile" binding:"required" example:"+886987654321"`
}

type Confirm struct {
	Mobile string `json:"mobile" binding:"required" example:"+886987654321"`
	Code   string `json:"code" binding:"required" example:"1234567"`
}
