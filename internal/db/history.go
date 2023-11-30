package db

import (
	"context"
	"github.com/Template7/backend/internal/db/entity"
)

func (c *client) CreateDepositHistory(ctx context.Context, data entity.DepositHistory) (err error) {
	log := c.log.WithContext(ctx)
	log.Debug("create deposit history")

	if err = c.sql.core.WithContext(ctx).Create(&data).Error; err != nil {
		log.WithError(err).Error("fail to create deposit history")
	}
	return
}

func (c *client) CreateWithdrawHistory(ctx context.Context, data entity.WithdrawHistory) (err error) {
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
