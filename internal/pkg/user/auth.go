package user

import (
	"backend/internal/pkg/db"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/route/middle_ware"
	"backend/internal/pkg/t7Error"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func GenToken(userId string) (token collection.UserToken, err *t7Error.Error) {
	log.Debug("gen token for user: ", userId)

	utc := middle_ware.UserTokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTtl).Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, utc)
	tokenString, signErr := tokenClaims.SignedString(secret)
	if signErr != nil {
		err = t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, signErr := jwt.New(jwt.SigningMethodHS256).SignedString(secret)
	if signErr != nil {
		err = t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
		return
	}
	token = collection.UserToken{
		UserToken:    tokenString,
		RefreshToken: refreshToken,
	}

	if dbErr := db.New().SaveUserToken(token); dbErr != nil {
		log.Error("fail to save user token: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

// TODO
func RefreshToken(oriToken collection.UserToken) (refreshedToken collection.UserToken, err *t7Error.Error) {
	return
}
