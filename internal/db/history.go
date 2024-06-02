package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
	"gorm.io/gorm"
)

func (c *client) createDepositHistory(ctx context.Context, tx *gorm.DB, data entity.DepositHistory) (err error) {
	log := c.log.WithContext(ctx)
	log.Debug("create deposit history")

	if err = tx.Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create deposit history")
	}
	return
}

func (c *client) createWithdrawHistory(ctx context.Context, tx *gorm.DB, data entity.WithdrawHistory) (err error) {
	log := c.log.WithContext(ctx)
	log.Debug("create withdraw history")

	if err = c.sql.core.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create withdraw history")
	}
	return
}

func (c *client) createTransferHistory(ctx context.Context, tx *gorm.DB, data entity.TransferHistory) (err error) {
	log := c.log.WithContext(ctx)
	log.Debug("create transfer history")

	if err = c.sql.core.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create transfer history")
	}
	return
}

func (c *client) GetWalletBalanceHistory(ctx context.Context, walletId string) ([]entity.WalletBalanceHistory, error) {
	log := c.log.WithContext(ctx).With("walletId", walletId)
	log.Debug("get wallet balance history")

	depositQuery := c.sql.core.Table("deposit_history").Select(`
        id, wallet_id, 'deposit' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ?", walletId)

	withdrawQuery := c.sql.core.Table("withdraw_history").Select(`
        id, wallet_id, 'withdraw' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ?", walletId)

	transferOutQuery := c.sql.core.Table("transfer_history").Select(`
        id, from_wallet_id AS wallet_id, 'transferOut' AS direction, currency, amount, sender_balance_before AS balance_before, sender_balance_after AS balance_after, note, created_at
    `).Where("from_wallet_id = ?", walletId)

	transferInQuery := c.sql.core.Table("transfer_history").Select(`
        id, to_wallet_id AS wallet_id, 'transferIn' AS direction, currency, amount, receiver_balance_before AS balance_before, receiver_balance_after AS balance_after, note, created_at
    `).Where("to_wallet_id = ?", walletId)

	unionQuery := c.sql.core.Table("(?) AS u", c.sql.core.Raw(`
        ? UNION ALL ? UNION ALL ? UNION ALL ?
    `, depositQuery, withdrawQuery, transferOutQuery, transferInQuery)).Order("created_at")

	var wbh []entity.WalletBalanceHistory
	if err := unionQuery.Scan(&wbh).Error; err != nil {
		log.WithError(err).Error("fail to scan data")
		return nil, err
	}

	return wbh, nil
}

func (c *client) GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]entity.WalletBalanceHistory, error) {
	log := c.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get wallet balance history by currency")

	depositQuery := c.sql.core.Table("deposit_history").Select(`
        id, wallet_id, 'deposit' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ? and currency = ?", walletId, currency)

	withdrawQuery := c.sql.core.Table("withdraw_history").Select(`
        id, wallet_id, 'withdraw' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ? and currency = ?", walletId, currency)

	transferOutQuery := c.sql.core.Table("transfer_history").Select(`
        id, from_wallet_id AS wallet_id, 'transferOut' AS direction, currency, amount, sender_balance_before AS balance_before, sender_balance_after AS balance_after, note, created_at
    `).Where("from_wallet_id = ? and currency = ?", walletId, currency)

	transferInQuery := c.sql.core.Table("transfer_history").Select(`
        id, to_wallet_id AS wallet_id, 'transferIn' AS direction, currency, amount, receiver_balance_before AS balance_before, receiver_balance_after AS balance_after, note, created_at
    `).Where("to_wallet_id = ? and currency = ?", walletId, currency)

	unionQuery := c.sql.core.Table("(?) AS u", c.sql.core.Raw(`
        ? UNION ALL ? UNION ALL ? UNION ALL ?
    `, depositQuery, withdrawQuery, transferOutQuery, transferInQuery)).Order("created_at")

	var wbh []entity.WalletBalanceHistory
	if err := unionQuery.Scan(&wbh).Error; err != nil {
		log.WithError(err).Error("fail to scan data")
		return nil, err
	}

	return wbh, nil
}
