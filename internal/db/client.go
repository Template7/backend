package db

import (
	"context"
	"github.com/Template7/common/config"
	"github.com/Template7/common/db"
	"github.com/Template7/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type client struct {
	mongo struct {
		user               *mongo.Collection
		transactionHistory *mongo.Collection
		depositHistory     *mongo.Collection
		withdrawHistory    *mongo.Collection
	}
	sql struct {
		core *gorm.DB
	}
	log *logger.Logger
}

func New(log *logger.Logger) Client {
	c := &client{
		log: log.WithService("db"),
	}

	// nosql
	mDb := db.NewNoSql().Database(config.New().Db.NoSql.Db)
	c.mongo.user = mDb.Collection("user")
	c.mongo.transactionHistory = mDb.Collection("transactionHistory")
	c.mongo.depositHistory = mDb.Collection("depositHistory")
	c.mongo.withdrawHistory = mDb.Collection("withdrawHistory")

	// sql
	c.sql.core = db.NewSql()

	c.log.Info("config initialized")
	return c
}

func (c *client) initIndex(db *mongo.Database) (err error) {
	ctx := context.Background()
	for col, idx := range CollectionIndexes {
		logger.New().With("collection", col).Debug("create index")
		_, err = db.Collection(col).Indexes().CreateMany(ctx, idx)
		if err != nil {
			c.log.WithError(err).Error("fail to init mongo index")
			return
		}
	}
	return
}
