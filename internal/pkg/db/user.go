package db

import (
	"context"
	"github.com/Template7/backend/internal/pkg/db/entity"
	"github.com/Template7/common/structs"
)

func (c *client) CreateUser(ctx context.Context, data entity.User) (err error) {
	log := c.log.WithContext(ctx).With("username", data.Username)
	log.Debug("create user")

	if err = c.sql.db.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create user")
	}
	return
}

func (c *client) GetUser(ctx context.Context, username string) (data entity.User, err error) {
	log := c.log.WithContext(ctx).With("username", username)
	log.Debug("get user")

	if err = c.sql.db.WithContext(ctx).Where("username = ?", username).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}

func (c *client) GetUserById(ctx context.Context, userId string) (data structs.User, err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user")

	if err = c.sql.db.WithContext(ctx).Where("id = ?", userId).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}
