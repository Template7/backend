package auth

import (
	"context"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"time"
)

const (
	userTokenTtl = 24 * time.Hour // 1 week
)

type Auth interface {
	ParseToken(ctx context.Context, token string) (data *v1.TokenClaims, err error)
	CheckPermission(ctx context.Context, sub, obj, act string) bool
	Login(ctx context.Context, username string, password string) (token string, err error)
	GetUserRole(ctx context.Context, username string) v1.Role
}