package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
)

func (c *client) CreateUser(ctx context.Context, data entity.User) (err error) {
	log := c.log.WithContext(ctx).With("username", data.Username)
	log.Debug("create user")

	if err = c.sql.core.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create user")
	}
	return
}

func (c *client) GetUser(ctx context.Context, username string) (data entity.User, err error) {
	log := c.log.WithContext(ctx).With("username", username)
	log.Debug("get user")

	if err = c.sql.core.WithContext(ctx).Where("username = ?", username).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}

func (c *client) GetUserById(ctx context.Context, userId string) (data entity.User, err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user")

	if err = c.sql.core.WithContext(ctx).Where("id = ?", userId).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}

func (c *client) UpdateUserInfo(ctx context.Context, userId string, info entity.UserInfo) (err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("update user info")

	if err = c.sql.core.WithContext(ctx).Where("user_id = ?", userId).Update("nickname", info.NickName).Error; err != nil {
		log.WithError(err).Error("fail to update user info")
	}
	return
}
