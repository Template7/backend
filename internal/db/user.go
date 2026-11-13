package db

import (
	"context"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (c *client) CreateUser(ctx context.Context, data models.User) (err error) {
	log := c.log.WithContext(ctx).With("username", data.Username)
	log.Debug("create user")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err = tx.Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create user")
		return
	}

	w := models.Wallet{
		Id:     uuid.NewString(),
		UserId: data.Id,
	}
	if err = tx.Create(&w).Error; err != nil {
		log.WithError(err).Error("fail to create default user wallet")
		return
	}

	bls := []models.Balance{
		{
			WalletId: w.Id,
			Money: models.Money{
				Currency: v1.Currency_usd.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: models.Money{
				Currency: v1.Currency_ntd.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: models.Money{
				Currency: v1.Currency_cny.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: models.Money{
				Currency: v1.Currency_jpy.String(),
			},
		},
	}

	if err = tx.Create(&bls).Error; err != nil {
		log.WithError(err).Error("fail to create default user wallet balances")
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.WithError(err).Error("fail to commit tx")
	}
	return
}

func (c *client) GetUser(ctx context.Context, username string) (data models.User, err error) {
	log := c.log.WithContext(ctx).With("username", username)
	log.Debug("get user")

	if err = c.sql.core.WithContext(ctx).Where("username = ?", username).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}

func (c *client) GetUserById(ctx context.Context, userId string) (data models.User, err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user")

	if err = c.sql.core.WithContext(ctx).Where("id = ?", userId).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user")
	}
	return
}

func (c *client) UpdateUserInfo(ctx context.Context, userId string, info models.UserInfo) (err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("update user info")

	if err = c.sql.core.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId).Update("nickname", info.Nickname).Error; err != nil {
		log.WithError(err).Error("fail to update user info")
	}
	return
}

func (c *client) GetUserWallets(ctx context.Context, userId string) (data []models.UserWalletBalance) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get user wallets balances")

	if err := c.sql.core.WithContext(ctx).Raw("select w.user_id, w.id as wallet_id, b.currency, b.amount from wallet w join balance b on w.id = b.wallet_id where w.user_id = ?", userId).Scan(&data).Error; err != nil {
		log.WithError(err).Error("fail to get user wallets balances")
	}
	return
}

func (c *client) DeleteUser(ctx context.Context, userId string) (err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("delete user")

	// separate deletes instead of one join-delete: an inner join skipped users
	// without wallet/balance rows. History is kept for audit.
	err = c.sql.core.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		walletIds := tx.Model(&models.Wallet{}).Select("id").Where("user_id = ?", userId)
		if err := tx.Where("wallet_id IN (?)", walletIds).Delete(&models.Balance{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userId).Delete(&models.Wallet{}).Error; err != nil {
			return err
		}
		res := tx.Where("id = ?", userId).Delete(&models.User{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
	if err != nil {
		log.WithError(err).Error("fail to delete user")
	}
	return
}

func (c *client) SetUserStatus(ctx context.Context, userId string, status authV1.AccountStatus) (err error) {
	log := c.log.WithContext(ctx).With("userId", userId).With("status", status)
	log.Debug("set account status")

	if err = c.sql.core.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId).Update("status", status).Error; err != nil {
		log.WithError(err).Error("fail to update account status")
	}

	return
}
