package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/google/uuid"
)

func (c *client) CreateUser(ctx context.Context, data entity.User) (err error) {
	log := c.log.WithContext(ctx).With("username", data.Username)
	log.Debug("create user")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err = tx.Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create user")
		return
	}

	w := entity.Wallet{
		Id:     uuid.NewString(),
		UserId: data.Id,
	}
	if err = tx.Create(&w).Error; err != nil {
		log.WithError(err).Error("fail to create default user wallet")
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.WithError(err).Error("fail to commit tx")
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

	if err = c.sql.core.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userId).Update("nickname", info.Nickname).Error; err != nil {
		log.WithError(err).Error("fail to update user info")
	}
	return
}

func (c *client) GetUserWallets(ctx context.Context, userId string) (data []entity.Wallet) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get wallet")

	if err := c.sql.core.Model(&entity.Wallet{}).Preload("Balance").Where("user_id = ?", userId).Find(&data).Error; err != nil {
		log.WithError(err).Error("fail to get wallet")
		return
	}
	return
}
