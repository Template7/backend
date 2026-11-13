package auth

import (
	"context"
	"time"

	"github.com/Template7/backend/internal/t7Error"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (s *service) ParseToken(_ context.Context, token string) (*UserTokenClaims, error) {
	tk, err := jwt.ParseWithClaims(token, &UserTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		return nil, t7Error.InvalidToken.WithDetail(err.Error())
	}

	claims, ok := tk.Claims.(*UserTokenClaims)
	if !ok {
		return nil, t7Error.TokenAssertionFail
	}
	return claims, nil
}

// genUserToken signs a token for the user. Status is informational only;
// the middleware always checks the current status.
func (s *service) genUserToken(userId string, role v1.Role, status v1.AccountStatus) (string, error) {
	now := time.Now()
	claims := UserTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(userTokenTtl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserId: userId,
		Role:   role.String(),
		Status: int(status),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims).SignedString(s.jwtSecret)
	if err != nil {
		return "", t7Error.TokenSignFail.WithDetail(err.Error())
	}
	return signed, nil
}
