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

func (c *client) CreateTransferHistory(ctx context.Context, data entity.TransferHistory) (err error) {
	log := c.log.WithContext(ctx)
	log.Debug("create transfer history")

	if err = c.sql.core.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create transfer history")
	}
	return
}

func (c *client) GetWalletBalanceHistory(ctx context.Context, walletId string, currency string) ([]entity.WalletBalanceHistory, error) {
	log := c.log.WithContext(ctx).With("walletId", walletId).With("currency", currency)
	log.Debug("get wallet balance history")

	tx := c.sql.core.WithContext(ctx)
	defer tx.Rollback()

	var dep []entity.DepositHistory
	if err := tx.Select("id", "currency", "amount", "balance_before", "balance_after", "note", "created_at").Where("wallet_id = ? and currency = ?", walletId, currency).Find(&dep).Error; err != nil {
		log.WithError(err).Error("fail to get deposit history")
		return nil, err
	}

	var wit []entity.WithdrawHistory
	if err := tx.Select("id", "currency", "amount", "balance_before", "balance_after", "note", "created_at").Where("wallet_id = ? and currency = ?", walletId, currency).Find(&wit).Error; err != nil {
		log.WithError(err).Error("fail to get withdraw history")
		return nil, err
	}

	var trf []entity.TransferHistory
	if err := tx.Select("id", "currency", "amount", "sender_balance_before", "sender_balance_after", "note", "created_at").Where("from_wallet_id = ? and currency = ?", walletId, currency).Find(&trf).Error; err != nil {
		log.WithError(err).Error("fail to get withdraw history")
		return nil, err
	}

	var trt []entity.TransferHistory
	if err := tx.Select("id", "currency", "amount", "receiver_balance_before", "receiver_balance_after", "note", "created_at").Where("to_wallet_id = ? and currency = ?", walletId, currency).Find(&trt).Error; err != nil {
		log.WithError(err).Error("fail to get withdraw history")
		return nil, err
	}

	tx.Commit()

	// data merge
	wbh := make([]entity.WalletBalanceHistory, len(dep)+len(wit)+len(trf)+len(trt))
	for i, h := range dep {
		wbh[i] = entity.WalletBalanceHistory{
			RecordId:     h.Id,
			Io:           "in",
			Amount:       h.Amount,
			AmountBefore: h.BalanceBefore,
			AmountAfter:  h.BalanceAfter,
			Timestamp:    h.CreatedAt,
			Note:         h.Note,
		}
	}
	for i, h := range wit {
		wbh[i+len(dep)] = entity.WalletBalanceHistory{
			RecordId:     h.Id,
			Io:           "out",
			Amount:       h.Amount,
			AmountBefore: h.BalanceBefore,
			AmountAfter:  h.BalanceAfter,
			Timestamp:    h.CreatedAt,
			Note:         h.Note,
		}
	}
	for i, h := range trf {
		wbh[i+len(dep)+len(wit)] = entity.WalletBalanceHistory{
			RecordId:     h.Id,
			Io:           "out",
			Amount:       h.Amount,
			AmountBefore: h.SenderBalanceBefore,
			AmountAfter:  h.SenderBalanceAfter,
			Timestamp:    h.CreatedAt,
			Note:         h.Note,
		}
	}
	for i, h := range trt {
		wbh[i+len(dep)+len(wit)+len(trf)] = entity.WalletBalanceHistory{
			RecordId:     h.Id,
			Io:           "in",
			Amount:       h.Amount,
			AmountBefore: h.ReceiverBalanceBefore,
			AmountAfter:  h.ReceiverBalanceAfter,
			Timestamp:    h.CreatedAt,
			Note:         h.Note,
		}
	}
	return wbh, nil
}
