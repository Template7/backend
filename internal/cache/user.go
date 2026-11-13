package cache

import (
	"context"
	"fmt"
)

func (c *client) SetUserActivationCode(ctx context.Context, userId string, code string) error {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("set user activation code")

	if err := c.core.Set(ctx, activationCodeKey(userId), code, activationCodeTtl).Err(); err != nil {
		log.WithError(err).Error("fail to set user activation code")
		return err
	}
	return nil
}

func (c *client) GetUserActivationCode(ctx context.Context, userId string) (string, error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user activation code")

	code, err := c.core.Get(ctx, activationCodeKey(userId)).Result()
	if err != nil {
		log.WithError(err).Warn("fail to get user activation code")
		return "", err
	}
	return code, nil
}

func (c *client) DeleteUserActivationCode(ctx context.Context, userId string) error {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("delete user activation code")

	if err := c.core.Del(ctx, activationCodeKey(userId)).Err(); err != nil {
		log.WithError(err).Error("fail to delete user activation code")
		return err
	}
	return nil
}

func activationCodeKey(userId string) string {
	return fmt.Sprintf("%s:%s", cacheKeyUserActivationCode, userId)
}
