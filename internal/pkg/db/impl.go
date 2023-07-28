package db

import (
	"context"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/common/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"

	"sync"
)

type impl struct {
	mongo struct {
		client             *mongo.Client
		user               *mongo.Collection
		transactionHistory *mongo.Collection
		depositHistory     *mongo.Collection
		withdrawHistory    *mongo.Collection
	}
	mysql struct {
		db *gorm.DB
	}
}

type QueryOption struct {
	Limit  int64
	Offset int64
	SortBy string
	Desc   bool
}

var (
	once     sync.Once
	instance *impl
)

func New() ClientInterface {
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
		instance = &impl{}
		instance.mongo.client = c
		instance.mongo.user = db.Collection("user")
		instance.mongo.transactionHistory = db.Collection("transactionHistory")
		instance.mongo.depositHistory = db.Collection("depositHistory")
		instance.mongo.withdrawHistory = db.Collection("withdrawHistory")

		// mysql
		sqlDb, err := gorm.Open(mysql.Open(config.New().MySql.ConnectionString), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		instance.mysql.db = sqlDb
		logger.New().Debug("db client initialized")
	})
	return instance
}

func (c *impl) initIndex(db *mongo.Database) (err error) {
	ctx := context.Background()
	for col, idx := range CollectionIndexes {
		logger.New().With("collection", col).Debug("create index")
		_, err = db.Collection(col).Indexes().CreateMany(ctx, idx)
		if err != nil {
			log.Error("fail to create index: ", err.Error())
			return
		}
	}
	return
}
