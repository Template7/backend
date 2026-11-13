package auth

import (
	"context"
	"github.com/Template7/backend/internal/t7Error"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (s *service) ParseToken(ctx context.Context, token string) (data *UserTokenClaims, err error) {
	log := s.log.WithContext(ctx)
	log.Debug("parse token")

	tk, err := jwt.ParseWithClaims(token, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		log.WithError(err).Error("fail to parse token")
		return nil, t7Error.InvalidToken.WithDetail(err.Error())
	}

	claims, ok := tk.Claims.(*UserTokenClaims)
	if !ok {
		log.Error("token assertion fail")
		err = t7Error.TokenAssertionFail
		return
	}

	return claims, nil
}

func (s *service) genUserToken(ctx context.Context, userId string, role v1.Role) (string, error) {
	log := s.log.WithContext(ctx).With("userId", userId).With("role", role)
	log.Debug("gen user token")

	utc := UserTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(userTokenTtl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserId: userId,
		Role:   role.String(),
		Status: int(s.GetUserStatus(ctx, userId)),
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &utc)
	tokenString, signErr := tokenClaims.SignedString(s.jwtSecret)
	if signErr != nil {
		log.WithError(signErr).Error("fail to sign jwt")
		return "", t7Error.TokenSignFail.WithDetail(signErr.Error())
	}
	return tokenString, nil
}
