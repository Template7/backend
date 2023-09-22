package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var CollectionIndexes = map[string][]mongo.IndexModel{}
