package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
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

	bls := []entity.Balance{
		{
			WalletId: w.Id,
			Money: entity.Money{
				Currency: v1.Currency_usd.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: entity.Money{
				Currency: v1.Currency_ntd.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: entity.Money{
				Currency: v1.Currency_cny.String(),
			},
		},
		{
			WalletId: w.Id,
			Money: entity.Money{
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

func (c *client) GetUserWallets(ctx context.Context, userId string) (data []entity.UserWalletBalance) {
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

	if err = c.sql.core.WithContext(ctx).Exec("delete u, b, w from user u join wallet w on u.id = w.user_id join balance b on w.id = b.wallet_id where u.id = ?", userId).Error; err != nil {
		log.WithError(err).Error("fail to delete user")
	}

	return
}
