package apiBody

type CreateUserReq struct {
	Mobile string `json:"mobile" bson:"mobile" example:"+886987654321"` // +886987654321
	Email  string `json:"email" bson:"email" example:"username@mail.com"`
}

type SmsReq struct {
	Mobile string `json:"mobile" binding:"required" example:"+886987654321"`
}
