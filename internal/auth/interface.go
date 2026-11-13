package auth

import (
	"context"
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"time"
)

const (
	userTokenTtl = 24 * time.Hour
)

type Auth interface {
	ParseToken(ctx context.Context, token string) (data *UserTokenClaims, err error)
	CheckPermission(ctx context.Context, sub, obj, act string) bool
	Login(ctx context.Context, username string, password string) (token string, err error)
	GetUserRole(ctx context.Context, username string) v1.Role
	// GetUserStatus returns the current account status, or -1 if the user can't be loaded.
	GetUserStatus(ctx context.Context, userId string) v1.AccountStatus
	CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (string, error)
	GenActivationCode(ctx context.Context, userId string) (string, error)
	DeleteUser(ctx context.Context, userId string) error
	ActivateUser(ctx context.Context, userId string, actCode string) bool
}
