package user

import "github.com/dgrijalva/jwt-go"

type TokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}
