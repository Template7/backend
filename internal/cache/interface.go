package cache

import (
	"context"
	"time"
)

const (
	cacheKeyUserActivationCode = "activationCode"
	cacheKeyLoginFailures      = "loginFailures"
	activationCodeTtl          = 24 * time.Hour
)

type Interface interface {
	SetUserActivationCode(ctx context.Context, userId string, code string) error
	GetUserActivationCode(ctx context.Context, userId string) (string, error)
	DeleteUserActivationCode(ctx context.Context, userId string) error

	// IncrLoginFailures counts one failed login for key and returns the count in the
	// current window; the window starts with the first failure.
	IncrLoginFailures(ctx context.Context, key string, window time.Duration) (int64, error)
	// GetLoginFailures returns the failures counted for key in the current window (0 if none).
	GetLoginFailures(ctx context.Context, key string) (int64, error)
	ResetLoginFailures(ctx context.Context, key string) error
}
