package db

import (
	"context"
	"github.com/Template7/common/config"
	"github.com/Template7/common/db"
	"github.com/Template7/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
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

type QueryOption struct {
	Limit  int64
	Offset int64
	SortBy string
	Desc   bool
}

var (
	once     sync.Once
	instance *client
)

func New() Client {
	once.Do(func() {
		// nosql
		mDb := db.NewNoSql().Database(config.New().Db.NoSql.Db)
		instance.mongo.user = mDb.Collection("user")
		instance.mongo.transactionHistory = mDb.Collection("transactionHistory")
		instance.mongo.depositHistory = mDb.Collection("depositHistory")
		instance.mongo.withdrawHistory = mDb.Collection("withdrawHistory")

		// sql
		instance.sql.core = db.NewSql()
		instance.log = logger.New().WithService("db")

		logger.New().Debug("db client initialized")
	})
	return instance
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
