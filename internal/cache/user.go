package cache

import (
	"context"
	"fmt"
)

func (c *client) SetUserActivationCode(ctx context.Context, userId string, code string) error {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("set user activation code")

	if err := c.core.Set(ctx, fmt.Sprintf("%s:%s", cacheKeyUserActivationCode, userId), code, 0).Err(); err != nil {
		log.WithError(err).Error("fail to set user activation code")
		return err
	}
	return nil
}

func (c *client) GetUserActivationCode(ctx context.Context, userId string) (string, error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user activation code")

	code, err := c.core.Get(ctx, fmt.Sprintf("%s:%s", cacheKeyUserActivationCode, userId)).Result()
	if err != nil {
		log.WithError(err).Error("fail to get user activation code")
		return "", err
	}
	return code, nil
}
