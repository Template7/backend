package collection

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ClaimTypeUser  = "user"
	ClaimTypeAdmin = "admin"
)

type claimType string

type Token struct {
	Id           *primitive.ObjectID    `json:"id,omitempty" bson:"id,omitempty"`
	AccessToken  string                 `json:"access_token" bson:"access_token"`
	RefreshToken string                 `json:"refresh_token" bson:"refresh_token"`
	ClaimType    claimType              `json:"-" bson:"claim_type"`
	OtherInfo    map[string]interface{} `json:"-" bson:"other_info"`
}
