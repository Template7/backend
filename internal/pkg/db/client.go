package db

import (
	"github.com/Template7/backend/internal/pkg/config"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type client struct {
	admin *mongo.Collection
	user  *mongo.Collection
	token *mongo.Collection
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

func New() *client {
	once.Do(func() {
		c, err := mongo.Connect(nil, options.Client().ApplyURI(config.New().Mongo.ConnectionString))
		if err != nil {
			log.Fatal(err)
		}
		db := c.Database(config.New().Mongo.Db)
		instance = &client{
			admin: db.Collection("admin"),
			user:  db.Collection("user"),
			token: db.Collection("token"),
		}
		if err := c.Ping(nil, nil); err != nil {
			log.Fatal(err)
		}
		instance.initIndex(db)
		log.Debug("mongo client initialized")
	})
	return instance
}

func (c client) initIndex(db *mongo.Database) {
	ctx := context.Background()
	for col, idx := range CollectionIndexes {
		log.Debug("create index for collection: ", col)

		//db := c.Database(config.New().Mongo.Db)
		_, err := db.Collection(col).Indexes().CreateMany(ctx, idx)
		if err != nil {
			log.Error("unable to create index for collection: ", col, ". ", err.Error())
			panic(err)
		}
	}
}
