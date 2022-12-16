package db

import (
	"context"
	"github.com/Template7/common/structs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (c impl) CreateUser(user structs.User) (err error) {
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
		Id:     uuid.New().String(),
		UserId: user.UserId,
	}).Error
	return
}

func (c impl) GetUserById(userId string) (data structs.User, err error) {
	filter := bson.M{
		"user_id": userId,
	}
	return c.getUser(filter)
}

func (c impl) GetFbUser(fbUserId string) (data structs.User, err error) {
	filter := bson.M{
		"login_info.channel":         structs.LoginChannelFacebook,
		"login_info.channel_user_id": fbUserId,
	}
	return c.getUser(filter)
}

func (c impl) GetUserByMobile(mobile string) (data structs.User, err error) {
	filter := bson.M{
		"mobile": mobile,
	}
	return c.getUser(filter)
}

func (c impl) getUser(filter bson.M) (data structs.User, err error) {
	err = c.mongo.user.FindOne(context.Background(), filter).Decode(&data)
	return
}

func (c impl) GetUserBasicInfo(userId string) (data structs.UserInfo, err error) {
	log.Debug("get user info: ", userId)

	var temp structs.User
	filter := bson.M{
		"user_id": userId,
	}
	opt := options.FindOne().SetProjection(
		bson.M{
			"_id":        0,
			"basic_info": 1,
		},
	)

	err = c.mongo.user.FindOne(nil, filter, opt).Decode(&temp)
	data = temp.BasicInfo
	return
}

func (c impl) UpdateUserBasicInfo(userId string, data structs.UserInfo) (err error) {
	log.Debug("update user basic info: ", userId)

	filter := bson.M{
		"user_id": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"basic_info": data,
			//"status":     structs.UserStatusActivate,
		},
	}
	_, err = c.mongo.user.UpdateOne(context.Background(), filter, update)
	return
}

func (c impl) UpdateUserStatus(userId string, status structs.UserStatus) (err error) {
	log.Debug("update user status: ", userId)

	filter := bson.M{
		"user_id": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	_, err = c.mongo.user.UpdateOne(nil, filter, update)
	return
}
