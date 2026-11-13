package db

import (
	"context"

	"github.com/Template7/common/logger"
	"gorm.io/gorm"
)

type client struct {
	core *gorm.DB
	log  *logger.Logger
}

func New(sqlCore *gorm.DB, log *logger.Logger) Client {
	c := &client{
		core: sqlCore,
		log:  log.WithService("db"),
	}
	c.log.Info("db client initialized")
	return c
}

func (c *client) db(ctx context.Context) *gorm.DB {
	return c.core.WithContext(ctx)
}

func (c *client) Transaction(ctx context.Context, fn func(tx Client) error) error {
	return c.db(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&client{core: tx, log: c.log})
	})
}
