package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/structs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetWallet(ctx context.Context, userId string) (data entity.Wallet, err error) {
	log := c.log.WithContext(ctx).With("userId", userId)
	log.Debug("get wallet")

	// TODO: join wallet and balance
	var wallet structs.Wallet
	if err = c.sql.db.Model(&structs.Wallet{}).Where("userId = ?", userId).Take(&wallet).Error; err != nil {
		log.WithError(err).Error("fail to get wallet")
		return
	}
	return
}

func (c *client) Deposit(ctx context.Context, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("deposit")

	return c.deposit(ctx, c.sql.db, walletId, money)
}

func (c *client) Withdraw(ctx context.Context, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("withdraw")

	return c.withdraw(ctx, c.sql.db, walletId, money)
}

func (c *client) Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("fromWalletId", fromWalletId).With("toWalletId", toWalletId).With("money", money)
	log.Debug("transfer money")

	tx := c.sql.db.Begin()
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

	log.Debug("finish transfer money")
	return
}

func (c *client) deposit(ctx context.Context, tx *gorm.DB, walletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("do deposit")

	err = tx.Model(&entity.Balance{}).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "walletId"}, {Name: "currency"}},
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
