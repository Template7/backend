package auth

import (
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/structs"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	adminTokenTtl = 1 * time.Hour
	userTokenTtl  = 7 * 24 * time.Hour // 1 week
)

type UserTokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
	Status structs.UserStatus
}

func GenUserToken(userId string) (token structs.Token, err *t7Error.Error) {
	log.Debug("gen token for user: ", userId)

	userData, dbErr := db.New().GetUserById(userId)
	if dbErr != nil {
		log.Error("fail to get user data: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}

	utc := UserTokenClaims{
		UserId: userId,
		Status: userData.Status,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(userTokenTtl).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, utc)
	tokenString, signErr := tokenClaims.SignedString(config.New().JwtSign)
	if signErr != nil {
		err = t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
		return
	}
	return genToken(tokenString)
}

func GenAdminToken() (token structs.Token, err *t7Error.Error) {
	log.Debug("gen admin token")

	claim := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(adminTokenTtl).Unix(),
	}
	tokenString, jwtErr := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(config.New().JwtSign)
	if jwtErr != nil {
		log.Error("fail to sign token: ", jwtErr.Error())
		err = t7Error.TokenSignFail.WithDetailAndStatus(jwtErr.Error(), http.StatusInternalServerError)
		return
	}
	return genToken(tokenString)
}

func genToken(accessToken string) (token structs.Token, err *t7Error.Error) {
	log.Debug("gen token")

	refreshToken, signErr := jwt.New(jwt.SigningMethodHS256).SignedString(config.New().JwtSign)
	if signErr != nil {
		err = t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
		return
	}
	token = structs.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClaimType:    structs.ClaimTypeUser,
	}
	return
}

func RefreshToken(oriToken structs.Token) (refreshedToken structs.Token, err *t7Error.Error) {
	log.Debug("refresh token: ", oriToken.Id.Hex())

	switch oriToken.ClaimType {
	case structs.ClaimTypeUser:
		return refreshUserToken(oriToken)

	default:
		log.Warn("unsupported claim type: ", oriToken.ClaimType)
		err = t7Error.InvalidToken.WithStatus(http.StatusBadRequest)
		return
	}
}

func refreshUserToken(oriToken structs.Token) (refreshedToken structs.Token, err *t7Error.Error) {
	log.Debug("refresh user token: ", oriToken.Id.Hex())

	token, tokenErr := jwt.ParseWithClaims(oriToken.AccessToken, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.New().JwtSign, nil
	})

	if tokenErr != nil {
		log.Error("fail to parse token: ", tokenErr.Error())
		err = t7Error.InvalidToken.WithDetailAndStatus(tokenErr.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*UserTokenClaims)
	if !ok {
		log.Error("token assertion fail")
		err = t7Error.TokenAssertionFail.WithStatus(http.StatusBadRequest)
		return
	}

	refreshedToken, err = GenUserToken(claims.UserId)
	if err != nil {
		log.Error("fail to gen user token: ", err.Error())
	}
	return
}
