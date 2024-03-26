package cache

import "context"

var (
	cacheKeyUserActivationCode = "activationCode"
)

type Interface interface {
	SetUserActivationCode(ctx context.Context, userId string, code string) error
	GetUserActivationCode(ctx context.Context, userId string) (string, error)
}
