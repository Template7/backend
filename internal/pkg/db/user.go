package db

import (
	"context"
	"github.com/Template7/common/structs"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (c client) CreateUser(user structs.User) (err error) {
	log.Debug("create user")

	user.Status = structs.UserStatusInitialized
	user.LastUpdate = time.Now().Unix()
	if len(user.BasicInfo.ProfilePictures) == 0 {
		user.BasicInfo.ProfilePictures = []string{}
	}
	if len(user.BasicInfo.Hobbies) == 0 {
		user.BasicInfo.Hobbies = []string{}
	}

	if _, err = c.mongo.user.InsertOne(context.Background(), user); err != nil {
		log.Error("fail to insert user: ", err.Error())
		return
	}

	// create user wallet
	err = c.mysql.db.Model(&structs.Wallet{}).Create(&structs.Wallet{
		Id: uuid.New().String(),
		UserId: user.UserId,
	}).Error
	return
}

func (c client) GetUserById(userId string) (data structs.User, err error) {
	filter := bson.M{
		"user_id": userId,
	}
	err = c.mongo.user.FindOne(context.Background(), filter).Decode(&data)
	return
}

func (c client) GetUserBasicInfo(userId string) (data structs.UserInfo, err error) {
	log.Debug("get user info: ", userId)

	var temp structs.User
	filter := bson.M{
		"user_id": userId,
	}
	opt := options.FindOne().SetProjection(
		bson.M{
			"_id":          0,
			"basic_info":   1,
		},
	)

	err = c.mongo.user.FindOne(nil, filter, opt).Decode(&temp)
	data = temp.BasicInfo
	return
}

func (c client) UpdateUserBasicInfo(userId string, data structs.UserInfo) (err error) {
	log.Debug("update user basic info: ", userId)

	filter := bson.M{
		"userId": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"basic_info": data,
		},
	}
	_, err = c.mongo.user.UpdateOne(nil, filter, update, options.Update().SetUpsert(true))
	return
}
