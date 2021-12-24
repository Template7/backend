package db

import (
	"context"
	"github.com/Template7/backend/internal/pkg/config"
	"github.com/Template7/common/structs"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
	}
	mysql struct {
		db      *gorm.DB
		//wallet  *gorm.Model
		//balance *gorm.Model
	}
}

type QueryOption struct {
	Limit  int64
	Offset int64
	SortBy string
	Desc   bool
}

func (q *QueryOption) ToMongoOption() (option options.FindOptions) {
	//opt := options.Find()
	if q.Limit > 0 {
		option.Limit = &q.Limit
	}
	if q.Offset > 0 {
		option.Skip = &q.Offset
	}
	if q.SortBy != "" {
		option.SetSort(bson.M{
			q.SortBy: parseSortOrder(q.Desc),
		})
	} else {
		option.SetSort(bson.M{
			"_id": parseSortOrder(q.Desc),
		})
	}
	return
}

func parseSortOrder(o bool) int {
	if o {
		return 1
	} else {
		return -1
	}
}

var (
	once     sync.Once
	instance *client
)

func New() ClientInterface {
	once.Do(func() {
		// mongo
		c, err := mongo.Connect(nil, options.Client().ApplyURI(config.New().Mongo.ConnectionString))
		if err != nil {
			log.Fatal(err)
		}
		if err := c.Ping(nil, nil); err != nil {
			log.Fatal(err)
		}
		db := c.Database(config.New().Mongo.Db)
		instance = &client{}
		instance.mongo.client = c
		instance.mongo.user = db.Collection("user")
		instance.mongo.transactionHistory = db.Collection("transactionHistory")

		// mysql
		sqlDb, err := gorm.Open(mysql.Open(config.New().MySql.ConnectionString), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		instance.mysql.db = sqlDb
		//instance.mysql.wallet = sqlDb.Model(&structs.Wallet{})
		instance.initDb()
		log.Debug("mongo client initialized")
	})
	return instance
}

func (c client) initIndex(db *mongo.Database) (err error) {
	ctx := context.Background()
	for col, idx := range CollectionIndexes {
		log.Debug("create index for collection: ", col)
		_, err = db.Collection(col).Indexes().CreateMany(ctx, idx)
		if err != nil {
			log.Error("fail to create index: ", err.Error())
			return
		}
	}
	return
}

func (c client) initTable() error {
	return instance.mysql.db.AutoMigrate(&structs.Wallet{}, &structs.Balance{})
}

func (c client) initDb() {
	if err := c.initIndex(c.mongo.client.Database(config.New().Mongo.Db)); err != nil {
		log.Error("fail to init index: ", err.Error())
		panic(err)
	}
	if err := c.initTable(); err != nil {
		log.Error("fail to init table: ", err.Error())
		panic(err)
	}
}
