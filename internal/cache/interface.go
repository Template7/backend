package cache

import (
	"context"
	"time"
)

const (
	cacheKeyUserActivationCode = "activationCode"
	activationCodeTtl          = 24 * time.Hour
)

type Interface interface {
	SetUserActivationCode(ctx context.Context, userId string, code string) error
	GetUserActivationCode(ctx context.Context, userId string) (string, error)
	DeleteUserActivationCode(ctx context.Context, userId string) error
}
