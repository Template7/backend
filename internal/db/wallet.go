package db

import (
	"context"

	"github.com/Template7/common/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *client) CreateWallet(ctx context.Context, wallet models.Wallet, currencies []string) error {
	return c.db(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&wallet).Error; err != nil {
			return err
		}
		if len(currencies) == 0 {
			return nil
		}
		bls := make([]models.Balance, len(currencies))
		for i, cur := range currencies {
			bls[i] = models.Balance{WalletId: wallet.Id, Money: models.Money{Currency: cur}}
		}
		return tx.Create(&bls).Error
	})
}

func (c *client) GetWallet(ctx context.Context, walletId string) (data models.Wallet, err error) {
	err = c.db(ctx).Where("id = ?", walletId).Take(&data).Error
	return
}

func (c *client) GetUserWallets(ctx context.Context, userId string) (data []models.UserWalletBalance, err error) {
	err = c.db(ctx).Raw(
		"select w.user_id, w.id as wallet_id, b.currency, b.amount from wallet w join balance b on w.id = b.wallet_id where w.user_id = ? order by w.id, b.currency",
		userId).Scan(&data).Error
	return
}

func (c *client) GetWalletBalances(ctx context.Context, walletId string) (data []models.WalletBalance, err error) {
	err = c.db(ctx).Raw(
		"select w.id as wallet_id, b.currency, b.amount from wallet w join balance b on w.id = b.wallet_id where w.id = ? order by b.currency",
		walletId).Scan(&data).Error
	return
}

func (c *client) LockBalances(ctx context.Context, walletIds []string, currency string) (map[string]decimal.Decimal, error) {
	q := c.db(ctx).
		Select("wallet_id", "amount").
		Where("wallet_id IN ? AND currency = ?", walletIds, currency).
		Order("wallet_id") // same lock order everywhere avoids deadlocks between opposite transfers
	if c.core.Name() != "sqlite" { // SQLite has no row locks (and no FOR UPDATE); it locks the whole DB on write
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	var rows []models.Balance
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	r := make(map[string]decimal.Decimal, len(rows))
	for _, b := range rows {
		r[b.WalletId] = b.Amount
	}
	return r, nil
}

func (c *client) AddBalance(ctx context.Context, walletId string, currency string, delta decimal.Decimal) error {
	return checkFound(c.db(ctx).Model(&models.Balance{}).
		Where("wallet_id = ? AND currency = ?", walletId, currency).
		Update("amount", gorm.Expr("amount + ?", delta)))
}
