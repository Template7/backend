package db

import (
	"github.com/Template7/common/structs"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c client) SaveToken(token structs.Token) (tokenId *primitive.ObjectID, err error) {
	log.Debug("save token")

	res, err := c.token.InsertOne(context.Background(), token)
	if err != nil {
		log.Error("fail to insert document: ", err.Error())
		return
	}
	oId := res.InsertedID.(primitive.ObjectID)
	tokenId = &oId
	return
}

func (c client) RemoveToken(id *primitive.ObjectID) (err error) {
	log.Debug("remove token: ", id.Hex())

	filter := bson.M{
		"_id": id,
	}
	_, err = c.token.DeleteOne(context.Background(), filter)
	return
}

func (c client) GetToken(id *primitive.ObjectID) (token structs.Token, err error) {
	log.Debug("get token: ", id.Hex())

	filter := bson.M{
		"_id": id,
	}
	err = c.token.FindOne(context.Background(), filter).Decode(&token)
	return
}
