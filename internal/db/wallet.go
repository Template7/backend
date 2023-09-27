package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetWallet(ctx context.Context, walletId string) (data entity.Wallet, err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet")

	if err = c.sql.core.WithContext(ctx).Where("id = ?", walletId).Preload("Balance").Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get wallet")
	}
	return
}

func (c *client) Deposit(ctx context.Context, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("deposit")

	return c.deposit(ctx, c.sql.core.WithContext(ctx), walletId, money)
}

func (c *client) Withdraw(ctx context.Context, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("withdraw")

	return c.withdraw(ctx, c.sql.core.WithContext(ctx), walletId, money)
}

func (c *client) Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("fromWalletId", fromWalletId).With("toWalletId", toWalletId).With("money", money)
	log.Debug("transfer money")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()

	err = c.withdraw(ctx, tx, fromWalletId, money)
	if err != nil {
		log.WithError(err).Error("fail to withdraw money")
		return
	}
	err = c.deposit(ctx, tx, toWalletId, money)
	if err != nil {
		log.WithError(err).Error("fail to deposit money")
		return
	}

	if err = tx.Commit().Error; err != nil {
		log.WithError(err).Error("fail to commit tx")
		return
	}
	log.Debug("finish transfer money")
	return
}

func (c *client) deposit(ctx context.Context, tx *gorm.DB, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("do deposit")

	err = tx.Model(&entity.Balance{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "wallet_id"}, {Name: "currency"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", money.Amount)}),
		}).Create(&entity.Balance{WalletId: walletId, Money: money}).Error

	if err != nil {
		log.WithError(err).Error("fail to deposit")
	}
	return
}

func (c *client) withdraw(ctx context.Context, tx *gorm.DB, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("do withdraw")

	err = tx.Model(&entity.Balance{}).
		Take(&entity.Balance{}, "wallet_id = ? AND currency = ? AND amount >= ?", walletId, money.Currency, money.Amount).
		Update("amount", gorm.Expr("amount - ?", money.Amount)).Error

	if err != nil {
		log.WithError(err).Error("fail to withdraw")
	}
	return
}
