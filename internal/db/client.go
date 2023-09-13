package db

import (
	"context"
	"github.com/Template7/backend/internal/config"
	"github.com/Template7/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type client struct {
	mongo struct {
		client             *mongo.Client
		user               *mongo.Collection
		transactionHistory *mongo.Collection
		depositHistory     *mongo.Collection
		withdrawHistory    *mongo.Collection
	}
	sql struct {
		db *gorm.DB
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
		// mongo
		c, err := mongo.Connect(nil, options.Client().ApplyURI(config.New().Mongo.ConnectionString))
		if err != nil {
			panic(err)
		}
		if err := c.Ping(nil, nil); err != nil {
			panic(err)
		}

		db := c.Database(config.New().Mongo.Db)

		// mysql
		sqlDb, err := gorm.Open(mysql.Open(config.New().Sql.ConnectionString), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		conn, err := sqlDb.DB()
		if err != nil {
			panic(err)
		}
		if err := conn.Ping(); err != nil {
			panic(err)
		}

		instance = &client{}
		instance.mongo.client = c
		instance.mongo.user = db.Collection("user")
		instance.mongo.transactionHistory = db.Collection("transactionHistory")
		instance.mongo.depositHistory = db.Collection("depositHistory")
		instance.mongo.withdrawHistory = db.Collection("withdrawHistory")
		instance.sql.db = sqlDb
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
