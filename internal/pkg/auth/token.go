package auth

import (
	"backend/internal/pkg/config"
	"backend/internal/pkg/db"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/user"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	adminTokenTtl = 1 * time.Hour
	userTokenTtl  = 7 * 24 * time.Hour // 1 week
)

func GenUserToken(userId string) (token collection.Token, err *t7Error.Error) {
	log.Debug("gen token for user: ", userId)

	utc := user.TokenClaims{
		UserId: userId,
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

func GenAdminToken() (token collection.Token, err *t7Error.Error) {
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

func genToken(accessToken string) (token collection.Token, err *t7Error.Error) {
	log.Debug("gen token")

	refreshToken, signErr := jwt.New(jwt.SigningMethodHS256).SignedString(config.New().JwtSign)
	if signErr != nil {
		err = t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
		return
	}
	token = collection.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClaimType:    collection.ClaimTypeUser,
	}
	tokenId, dbErr := db.New().SaveToken(token)
	if dbErr != nil {
		log.Error("fail to save user token: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	token.Id = tokenId
	return
}

func RefreshToken(oriToken collection.Token) (refreshedToken collection.Token, err *t7Error.Error) {
	log.Debug("refresh token: ", oriToken.Id.Hex())

	ot, dbErr := db.New().GetToken(oriToken.Id)
	if dbErr != nil {
		log.Error("fail to get token: ", oriToken.Id.Hex(), ". ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}

	switch ot.ClaimType {
	case collection.ClaimTypeUser:
		return refreshUserToken(oriToken)

	default:
		log.Warn("unsupported claim type: ", ot.ClaimType)
		err = t7Error.InvalidToken.WithStatus(http.StatusBadRequest)
		return
	}
}

func refreshUserToken(oriToken collection.Token) (refreshedToken collection.Token, err *t7Error.Error) {
	log.Debug("refresh user token: ", oriToken.Id.Hex())

	token, tokenErr := jwt.ParseWithClaims(oriToken.AccessToken, &user.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.New().JwtSign, nil
	})

	if tokenErr != nil {
		log.Error("fail to parse token: ", tokenErr.Error())
		err = t7Error.TokenParseFail.WithDetailAndStatus(tokenErr.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := token.Claims.(*user.TokenClaims)
	if !ok {
		log.Error("token assertion fail")
		err = t7Error.TokenAssertionFail.WithStatus(http.StatusBadRequest)
		return
	}

	refreshedToken, err = GenUserToken(claims.UserId)
	if err != nil {
		log.Error("fail to gen user token: ", err.Error())
		return
	}
	if dbErr := db.New().RemoveToken(oriToken.Id); dbErr != nil {
		log.Error("fail to remove token: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}
