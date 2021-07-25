package db

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var CollectionIndexes = map[string][]mongo.IndexModel{
	"user": user,
}

var (
	user = []mongo.IndexModel{
		//{
		//	Keys: bsonx.Doc{
		//		{
		//			Key:   "login_client.channel",
		//			Value: bsonx.Int32(1),
		//		},
		//		{
		//			Key:   "login_client.channel_user_id",
		//			Value: bsonx.Int32(1),
		//		},
		//	},
		//	Options: options.Index().SetUnique(true),
		//},
		{
			Keys: bsonx.Doc{
				{
					Key:   "mobile",
					Value: bsonx.Int32(1),
				},
			},
		},
		{
			Keys: bsonx.Doc{
				{
					Key:   "email",
					Value: bsonx.Int32(1),
				},
			},
		},
	}
)
