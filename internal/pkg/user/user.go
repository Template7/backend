package user

import (
	"backend/internal/pkg/config"
	"backend/internal/pkg/db"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/t7Redis"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

const (
	tokenTtl = 365 * 24 * time.Hour // 1 year
)

var secret = []byte(config.New().JwtSign)

func Exist(filter db.GetUserFilter) (exist bool, err *t7Error.Error) {
	log.Debug("check user exist")

	users, err := GetUsers(filter, db.QueryOption{})
	exist = len(users) > 0
	log.Debug("matched count: ", len(users))
	return
}

func GetUsers(filter db.GetUserFilter, option db.QueryOption) (users []collection.User, err *t7Error.Error) {
	log.Debug("get user")

	users, dbErr := db.New().GetUser(filter, option)
	if dbErr == nil {
		return
	}

	log.Warn("fail to get user: ", dbErr.Error())
	switch dbErr {
	case mongo.ErrNoDocuments:
		log.Info("no matched user")
		//err = t7Error.UserNotfound.WithStatus(http.StatusNoContent)
	default:
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func GetByChannel(channel collection.LoginChannel, id string) (data collection.User, err *t7Error.Error) {
	log.Debug("get user by channel")

	data, dbErr := db.New().GetUserByChannel(channel, id)
	if dbErr == nil {
		return
	}

	log.Warn("fail to get user: ", dbErr.Error())
	switch dbErr {
	case mongo.ErrNoDocuments:
		log.Info("no matched user")
		err = t7Error.UserNotfound.WithStatus(http.StatusNoContent)
	default:
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func GetByMobile(mobile string) (data collection.User, err *t7Error.Error){
	data, dbErr := db.New().GetUserByMobile(mobile)
	if dbErr != nil {
		log.Error("fail to get user: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func GetInfo(userId string) (data collection.User, err *t7Error.Error) {
	uId, idErr := primitive.ObjectIDFromHex(userId)
	if idErr != nil {
		err = t7Error.InvalidDocumentId.WithDetail(idErr.Error())
		return
	}
	data, dbErr := db.New().GetUserInfo(uId)
	if dbErr != nil {
		err = t7Error.InvalidDocumentId.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func CreateUser(user collection.User) (userId *primitive.ObjectID, err *t7Error.Error) {
	userId, dbErr := db.New().CreateUser(user)

	if dbErr == nil {
		return
	}

	// check user exist
	switch dbErr.(type) {
	case mongo.WriteException:
		err = t7Error.UserAlreadyExist
	default:
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func CreateNativeUser(mobile string) (u collection.User, err *t7Error.Error) {
	log.Debug("create native user")

	u = collection.User{
		Mobile: mobile,
		Status: collection.UserStatusInitialized,
	}
	uId, dbErr := db.New().CreateUser(u)
	if dbErr != nil {
		// TODO: bad request for duplicated mobile
		log.Error("fail to create native user: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}
	u.Id = uId
	return
}

func DeleteUser(userId string) (err *t7Error.Error) {
	uId, idErr := primitive.ObjectIDFromHex(userId)
	if idErr != nil {
		err = t7Error.InvalidDocumentId.WithDetail(idErr.Error())
		return
	}

	if dbErr := db.New().DeleteUser(uId); dbErr != nil {
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func UpdateBasicInfo(userId string, data collection.UserInfo) (err *t7Error.Error) {
	uId, idErr := primitive.ObjectIDFromHex(userId)
	if idErr != nil {
		log.Warn("invalid user id: ", userId)
		err = t7Error.InvalidDocumentId.WithDetail(idErr.Error())
		return
	}

	if dbErr := db.New().UpdateBasicInfo(uId, data); dbErr != nil {
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func UpdateLoginClient(userId string, loginClient collection.LoginInfo) (err *t7Error.Error) {
	uId, idErr := primitive.ObjectIDFromHex(userId)
	if idErr != nil {
		log.Warn("invalid user id: ", userId)
		err = t7Error.InvalidDocumentId.WithDetail(idErr.Error())
		return
	}
	if dbErr := db.New().UpdateLoginClient(uId, loginClient); dbErr != nil {
		log.Error("fail to update login client: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func SignOut(token string, expiresAt int64) (err *t7Error.Error) {
	log.Debug("user logout")

	life := time.Duration(expiresAt-time.Now().Unix()) * time.Second
	r := t7Redis.New().Set(token, nil, life)
	if r.Err() != nil {
		err = t7Error.RedisOperationFail.WithDetailAndStatus(r.Err().Error(), http.StatusInternalServerError)
	}
	return
}
