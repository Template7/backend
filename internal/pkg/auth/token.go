package auth

import (
	"context"
	"github.com/Template7/backend/internal/pkg/t7Error"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func (s *service) ParseToken(ctx context.Context, token string) (data *v1.TokenClaims, err error) {
	log := s.log.WithContext(ctx)
	log.Debug("parse token")

	tk, err := jwt.ParseWithClaims(token, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSign, nil
	})
	if err != nil {
		log.WithError(err).Error("fail to parse token")
		return nil, t7Error.InvalidToken.WithDetail(err.Error())
	}

	claims, ok := tk.Claims.(*UserTokenClaims)
	if !ok {
		log.Error("token assertion fail")
		err = t7Error.TokenAssertionFail.WithStatus(http.StatusBadRequest)
		return
	}

	return &claims.TokenClaims, nil
}

func (s *service) genUserToken(ctx context.Context, userId string, role v1.Role) (string, error) {
	log := s.log.WithContext(ctx).With("userId", userId).With("role", role)
	log.Debug("gen user token")

	utc := UserTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(userTokenTtl).Unix(),
		},
		TokenClaims: v1.TokenClaims{
			UserId: userId,
			Role:   role,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &utc)
	tokenString, signErr := tokenClaims.SignedString(jwtSign)
	if signErr != nil {
		log.WithError(signErr).Error("fail to sign jwt")
		return "", t7Error.TokenSignFail.WithDetailAndStatus(signErr.Error(), http.StatusInternalServerError)
	}
	return tokenString, nil
}
