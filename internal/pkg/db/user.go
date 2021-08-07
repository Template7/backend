package db

import (
	"github.com/Template7/backend/internal/pkg/db/collection"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type GetUserFilter struct {
	NickName   string `json:"nick_name" bson:"nick_name,omitempty"`
	Gender     string `json:"gender,omitempty" bson:"gender,omitempty"` // male | female | other
	Mobile     string `json:"mobile" bson:"mobile,omitempty"`           // +886987654321
	Email      string `json:"email,omitempty" bson:"email,omitempty"`
	LoginOs    string `json:"login_mobile" bson:"login_mobile,omitempty"` // ios | android
	SignUpFrom string `json:"sign_up_from" bson:"sign_up_from,omitempty"` // google | facebook | phone
	Status     string `json:"status" bson:"status,omitempty"`             // initialized | activate | block
}

func (c client) CreateUser(user collection.User) (userId *primitive.ObjectID, err error) {
	log.Debug("create user")

	user.Status = collection.UserStatusInitialized
	user.LastUpdate = time.Now().Unix()
	if len(user.BasicInfo.ProfilePictures) == 0 {
		user.BasicInfo.ProfilePictures = []string{}
	}
	if len(user.BasicInfo.Hobbies) == 0 {
		user.BasicInfo.Hobbies = []string{}
	}

	res, err := c.user.InsertOne(nil, user)
	if err != nil {
		log.Error("fail to insert user: ", err.Error())
		return
	}

	oId := res.InsertedID.(primitive.ObjectID)
	userId = &oId
	return
}

func (c client) GetUser(filter GetUserFilter, option QueryOption) (users []collection.User, err error) {
	log.Debug("get user")

	opt := option.ToMongoOption()

	cursor, err := c.user.Find(context.Background(), filter, &opt)
	if err != nil {
		log.Error("fail to get documents")
		return
	}

	if err := cursor.All(context.Background(), &users); err != nil {
		log.Error("fail to decode document: ", err.Error())
	}
	return
}

func (c client) GetUserByChannel(channel collection.LoginChannel, id string) (data collection.User, err error) {
	log.Debug("get user by channel")

	filter := bson.M{
		"login_client.channel":         channel,
		"login_client.channel_user_id": id,
	}
	err = c.user.FindOne(context.Background(), filter).Decode(&data)
	return
}

func (c client) GetUserByMobile(mobile string) (data collection.User, err error) {
	log.Debug("get user by mobile: ", mobile)

	filter := bson.M{
		"mobile": mobile,
	}
	err = c.user.FindOne(context.Background(), filter).Decode(&data)
	return
}

func (c client) GetUserInfo(userId primitive.ObjectID) (data collection.User, err error) {
	log.Debug("get user info: ", userId.Hex())

	filter := bson.M{
		"_id": userId,
	}
	opt := options.FindOne().SetProjection(
		bson.M{
			"_id":          0,
			"login_client": 0,
			"create_at":    0,
			"last_update":  0,
		},
	)

	err = c.user.FindOne(nil, filter, opt).Decode(&data)
	return
}

func (c client) UpdateBasicInfo(userId primitive.ObjectID, data collection.UserInfo) (err error) {
	log.Debug("update user basic info: ", userId.Hex())

	filter := bson.M{
		"_id": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"basic_info": data,
		},
	}
	_, err = c.user.UpdateOne(nil, filter, update, options.Update().SetUpsert(true))
	return
}

func (c client) UpdateLoginClient(userId primitive.ObjectID, loginClient collection.LoginInfo) (err error) {
	log.Debug("update user login client: ", userId.Hex())

	filter := bson.M{
		"_id": userId,
	}

	update := bson.M{
		"$set": bson.M{
			"login_client": loginClient,
		},
	}
	_, err = c.user.UpdateOne(nil, filter, update, options.Update().SetUpsert(true))
	return
}

func (c client) DeleteUser(userId primitive.ObjectID) (err error) {
	log.Debug("delete user: ", userId.Hex())
	filter := bson.M{
		"_id": userId,
	}
	update := bson.M{
		"last_update": time.Now().Unix(),
		"status":      collection.UserStatusBlock,
	}
	_, err = c.user.UpdateOne(nil, filter, update)
	return
}
