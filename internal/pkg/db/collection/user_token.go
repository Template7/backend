package collection

type UserToken struct {
	UserToken    string `json:"user_token" bson:"user_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
}
