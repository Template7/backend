package db

import (
	"context"
	"fmt"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/t7Id"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) GetWalletBalances(ctx context.Context, walletId string) (data []entity.WalletBalance, err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet balances")

	if err = c.sql.core.WithContext(ctx).Raw("select w.id as wallet_id, b.currency, b.amount from wallet w join balance b on w.id = b.wallet_id where w.id = ?", walletId).Scan(&data).Error; err != nil {
		log.WithError(err).Error("fail to get wallet balances")
	}
	return
}

// Deposit
// deposit money and create deposit history
func (c *client) Deposit(ctx context.Context, walletId string, money entity.Money, note string) error {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("deposit")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()

	bb, err := c.getBalance(ctx, tx, walletId, money.Currency)
	if err != nil {
		log.WithError(err).Error("fail to get balance before")
		return err
	}

	if err := c.deposit(ctx, tx, walletId, money); err != nil {
		log.WithError(err).Error("fail to deposit money")
		return err
	}

	ba, err := c.getBalance(ctx, tx, walletId, money.Currency)
	if err != nil {
		log.WithError(err).Error("fail to get balance after")
		return err
	}

	dh := entity.DepositHistory{
		Id:            t7Id.New().Generate().Int64(),
		WalletId:      walletId,
		Currency:      money.Currency,
		Amount:        money.Amount,
		BalanceBefore: bb,
		BalanceAfter:  ba,
		Note:          note,
	}
	if err := c.createDepositHistory(ctx, tx, dh); err != nil {
		log.WithError(err).Error("fail to create deposit history")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("fail to commit transaction")
		return err
	}

	log.Debug("deposit done")
	return nil
}

func (c *client) Withdraw(ctx context.Context, walletId string, money entity.Money, note string) (err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("withdraw")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()

	bb, err := c.getBalance(ctx, tx, walletId, money.Currency)
	if err != nil {
		log.WithError(err).Error("fail to get balance before")
		return err
	}

	if err := c.withdraw(ctx, tx, walletId, money); err != nil {
		log.WithError(err).Error("fail to withdraw money")
		return err
	}

	ba, err := c.getBalance(ctx, tx, walletId, money.Currency)
	if err != nil {
		log.WithError(err).Error("fail to get balance after")
		return err
	}

	wh := entity.WithdrawHistory{
		Id:            t7Id.New().Generate().Int64(),
		WalletId:      walletId,
		Currency:      money.Currency,
		Amount:        money.Amount,
		BalanceBefore: bb,
		BalanceAfter:  ba,
		Note:          note,
	}
	if err := c.createWithdrawHistory(ctx, tx, wh); err != nil {
		log.WithError(err).Error("fail to create withdraw history")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.WithError(err).Error("fail to commit transaction")
		return err
	}

	log.Debug("withdraw done")
	return nil
}

func (c *client) Transfer(ctx context.Context, fromWalletId string, toWalletId string, money entity.Money) (err error) {
	log := c.log.WithContext(ctx).With("fromWalletId", fromWalletId).With("toWalletId", toWalletId).With("money", money)
	log.Debug("transfer money")

	tx := c.sql.core.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := tx.Error; err != nil {
		return err
	}

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

func (c *client) GetBalance(ctx context.Context, walletId string, currency string) (decimal.Decimal, error) {
	log := c.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get balance")

	return c.getBalance(ctx, c.sql.core.WithContext(ctx), walletId, currency)
}

func (c *client) getBalance(ctx context.Context, tx *gorm.DB, walletId string, currency string) (decimal.Decimal, error) {
	log := c.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get balance")

	var data entity.Balance
	if err := tx.Where(&entity.Balance{}).Select("amount").Where("wallet_id = ? and currency = ?", walletId, currency).Take(&data).Error; err != nil {
		log.WithError(err).Error("fail to get balance")
		return decimal.NewFromInt(-1), err
	}

	return data.Amount, nil
}

func (c *client) GetWalletsBalance(ctx context.Context, walletId []string, currency string) (data []entity.Balance, err error) {
	log := c.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get balances")

	if err = c.sql.core.WithContext(ctx).Where(&entity.Balance{}).Select("wallet_id", "amount").Where("wallet_id in ? and currency = ?", walletId, currency).Find(&data).Error; err != nil {
		log.WithError(err).Error("fail to get balance")
	}

	if len(data) != len(walletId) {
		log.With("walletIdLength", len(walletId)).With("dataLength", len(data)).Warn("length not match")
		err = fmt.Errorf("length not match")
	}

	return
}
