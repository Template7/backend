package db

import (
	"context"
	"github.com/Template7/common/models"
)

func (c *client) CreateDepositHistory(ctx context.Context, data models.DepositHistory) error {
	return c.db(ctx).Create(&data).Error
}

func (c *client) CreateWithdrawHistory(ctx context.Context, data models.WithdrawHistory) error {
	return c.db(ctx).Create(&data).Error
}

func (c *client) CreateTransferHistory(ctx context.Context, data models.TransferHistory) error {
	return c.db(ctx).Create(&data).Error
}

func (c *client) GetWalletBalanceHistory(ctx context.Context, walletId string) ([]models.WalletBalanceHistory, error) {
	depositQuery := c.core.Table("deposit_history").Select(`
        id, wallet_id, 'deposit' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ?", walletId)

	withdrawQuery := c.core.Table("withdraw_history").Select(`
        id, wallet_id, 'withdraw' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ?", walletId)

	transferOutQuery := c.core.Table("transfer_history").Select(`
        id, from_wallet_id AS wallet_id, 'transferOut' AS direction, currency, amount, sender_balance_before AS balance_before, sender_balance_after AS balance_after, note, created_at
    `).Where("from_wallet_id = ?", walletId)

	transferInQuery := c.core.Table("transfer_history").Select(`
        id, to_wallet_id AS wallet_id, 'transferIn' AS direction, currency, amount, receiver_balance_before AS balance_before, receiver_balance_after AS balance_after, note, created_at
    `).Where("to_wallet_id = ?", walletId)

	unionQuery := c.core.Table("(?) AS u", c.core.Raw(`
        ? UNION ALL ? UNION ALL ? UNION ALL ?
    `, depositQuery, withdrawQuery, transferOutQuery, transferInQuery)).Order("created_at")

	var wbh []models.WalletBalanceHistory
	if err := unionQuery.WithContext(ctx).Scan(&wbh).Error; err != nil {
		return nil, err
	}
	return wbh, nil
}

func (c *client) GetWalletBalanceHistoryByCurrency(ctx context.Context, walletId string, currency string) ([]models.WalletBalanceHistory, error) {
	depositQuery := c.core.Table("deposit_history").Select(`
        id, wallet_id, 'deposit' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ? and currency = ?", walletId, currency)

	withdrawQuery := c.core.Table("withdraw_history").Select(`
        id, wallet_id, 'withdraw' AS direction, currency, amount, balance_before, balance_after, note, created_at
    `).Where("wallet_id = ? and currency = ?", walletId, currency)

	transferOutQuery := c.core.Table("transfer_history").Select(`
        id, from_wallet_id AS wallet_id, 'transferOut' AS direction, currency, amount, sender_balance_before AS balance_before, sender_balance_after AS balance_after, note, created_at
    `).Where("from_wallet_id = ? and currency = ?", walletId, currency)

	transferInQuery := c.core.Table("transfer_history").Select(`
        id, to_wallet_id AS wallet_id, 'transferIn' AS direction, currency, amount, receiver_balance_before AS balance_before, receiver_balance_after AS balance_after, note, created_at
    `).Where("to_wallet_id = ? and currency = ?", walletId, currency)

	unionQuery := c.core.Table("(?) AS u", c.core.Raw(`
        ? UNION ALL ? UNION ALL ? UNION ALL ?
    `, depositQuery, withdrawQuery, transferOutQuery, transferInQuery)).Order("created_at")

	var wbh []models.WalletBalanceHistory
	if err := unionQuery.WithContext(ctx).Scan(&wbh).Error; err != nil {
		return nil, err
	}
	return wbh, nil
}
