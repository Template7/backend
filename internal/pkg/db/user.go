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

//type GetUserFilter struct {
//	NickName   string `json:"nick_name" bson:"nick_name,omitempty"`
//	Gender     string `json:"gender,omitempty" bson:"gender,omitempty"` // male | female | other
//	Mobile     string `json:"mobile" bson:"mobile,omitempty"`           // +886987654321
//	Email      string `json:"email,omitempty" bson:"email,omitempty"`
//	LoginOs    string `json:"login_mobile" bson:"login_mobile,omitempty"` // ios | android
//	SignUpFrom string `json:"sign_up_from" bson:"sign_up_from,omitempty"` // google | facebook | phone
//	Status     string `json:"status" bson:"status,omitempty"`             // initialized | activate | block
//}

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
	err = c.mysql.db.Create(&structs.Wallet{
		Id: uuid.New().String(),
		UserId: user.UserId,
	}).Error
	return
}

//func (c client) GetUser(filter GetUserFilter, option QueryOption) (users []structs.User, err error) {
//	log.Debug("get user")
//
//	opt := option.ToMongoOption()
//
//	cursor, err := c.mongo.user.Find(context.Background(), filter, &opt)
//	if err != nil {
//		log.Error("fail to get documents")
//		return
//	}
//
//	if err := cursor.All(context.Background(), &users); err != nil {
//		log.Error("fail to decode document: ", err.Error())
//	}
//	return
//}

func (c client) GetUserById(userId string) (data structs.User, err error) {
	filter := bson.M{
		"user_id": userId,
	}
	err = c.mongo.user.FindOne(context.Background(), filter).Decode(&data)
	return
}

//func (c client) GetUserByChannel(channel structs.LoginChannel, id string) (data structs.User, err error) {
//	log.Debug("get user by channel")
//
//	filter := bson.M{
//		"login_client.channel":         channel,
//		"login_client.channel_user_id": id,
//	}
//	err = c.mongo.user.FindOne(context.Background(), filter).Decode(&data)
//	return
//}

//func (c client) GetUserByMobile(mobile string) (data structs.User, err error) {
//	log.Debug("get user by mobile: ", mobile)
//
//	filter := bson.M{
//		"mobile": mobile,
//	}
//	err = c.mongo.user.FindOne(context.Background(), filter).Decode(&data)
//	return
//}

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

//func (c client) UpdateLoginClient(userId primitive.ObjectID, loginClient structs.LoginInfo) (err error) {
//	log.Debug("update user login client: ", userId.Hex())
//
//	filter := bson.M{
//		"_id": userId,
//	}
//
//	update := bson.M{
//		"$set": bson.M{
//			"login_client": loginClient,
//		},
//	}
//	_, err = c.mongo.user.UpdateOne(nil, filter, update, options.Update().SetUpsert(true))
//	return
//}

//func (c client) DeleteUser(userId primitive.ObjectID) (err error) {
//	log.Debug("delete user: ", userId.Hex())
//	filter := bson.M{
//		"_id": userId,
//	}
//	update := bson.M{
//		"last_update": time.Now().Unix(),
//		"status":      structs.UserStatusBlock,
//	}
//	_, err = c.mongo.user.UpdateOne(nil, filter, update)
//	return
//}
