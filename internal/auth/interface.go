// Package auth handles credentials, tokens, roles and permissions.
// User accounts themselves are managed by the user package.
package auth

import (
	"context"
	"time"

	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
)

const (
	userTokenTtl = 24 * time.Hour
)

type Auth interface {
	// Login checks the credentials and returns a signed token.
	Login(ctx context.Context, username string, password string) (token string, err error)
	ParseToken(ctx context.Context, token string) (*UserTokenClaims, error)
	// CheckPermission reports whether role sub may do act on obj (path, method).
	CheckPermission(ctx context.Context, sub, obj, act string) bool

	GetUserRole(ctx context.Context, userId string) (v1.Role, error)
	AddUserRole(ctx context.Context, userId string, role v1.Role) error
	DeleteUserRoles(ctx context.Context, userId string) error
}
