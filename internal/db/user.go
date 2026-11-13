package db

import (
	"context"

	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"gorm.io/gorm"
)

func (c *client) CreateUser(ctx context.Context, data models.User) error {
	return c.db(ctx).Create(&data).Error
}

func (c *client) GetUser(ctx context.Context, username string) (data models.User, err error) {
	err = c.db(ctx).Where("username = ?", username).Take(&data).Error
	return
}

func (c *client) GetUserById(ctx context.Context, userId string) (data models.User, err error) {
	err = c.db(ctx).Where("id = ?", userId).Take(&data).Error
	return
}

func (c *client) UpdateUserInfo(ctx context.Context, userId string, info models.UserInfo) error {
	return checkFound(c.db(ctx).Model(&models.User{}).Where("id = ?", userId).Update("nickname", info.Nickname))
}

func (c *client) SetUserStatus(ctx context.Context, userId string, status authV1.AccountStatus) error {
	return checkFound(c.db(ctx).Model(&models.User{}).Where("id = ?", userId).Update("status", status))
}

func (c *client) DeleteUser(ctx context.Context, userId string) error {
	// separate deletes instead of one join-delete: an inner join skips users
	// without wallet/balance rows. History is kept for audit.
	return c.db(ctx).Transaction(func(tx *gorm.DB) error {
		walletIds := tx.Model(&models.Wallet{}).Select("id").Where("user_id = ?", userId)
		if err := tx.Where("wallet_id IN (?)", walletIds).Delete(&models.Balance{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userId).Delete(&models.Wallet{}).Error; err != nil {
			return err
		}
		return checkFound(tx.Where("id = ?", userId).Delete(&models.User{}))
	})
}

// checkFound turns "no row affected" into ErrNotFound.
func checkFound(res *gorm.DB) error {
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
