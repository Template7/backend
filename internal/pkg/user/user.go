package user

import (
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/structs"
	"github.com/google/uuid"
	"net/http"
)

var (
	log = logger.GetLogger()
)

type GetInfoResp struct {
	UserInfo structs.UserInfo `json:",inline"`
	WalletId string           `json:"wallet_id"`
}

func GetInfo(userId string) (data GetInfoResp, err *t7Error.Error) {
	userBasicInfo, dbErr := db.New().GetUserBasicInfo(userId)
	if dbErr != nil {
		err = t7Error.InvalidDocumentId.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}
	wallet, dbErr := db.New().GetWallet(userId)
	if dbErr != nil {
		err = t7Error.InvalidDocumentId.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}
	data = GetInfoResp{
		UserInfo: userBasicInfo,
		WalletId: wallet.Id,
	}
	return
}

type CreateUserReq struct {
	Mobile string `json:"mobile" bson:"mobile" example:"+886987654321"` // +886987654321
	Email  string `json:"email" bson:"email" example:"username@mail.com"`
}

func CreateUser(req CreateUserReq) (userId string, err *t7Error.Error) {
	log.Debug("create user")

	data := structs.User{
		UserId: uuid.New().String(),
		Mobile: req.Mobile,
		Email:  req.Email,
		Status: structs.UserStatusInitialized,
	}
	if err = createUser(data); err != nil {
		return
	}
	return data.UserId, nil
}

func CreateNativeUser(mobile string) (userId string, err *t7Error.Error) {
	log.Debug("create native user")

	data := structs.User{
		UserId: uuid.New().String(),
		Mobile: mobile,
		Status: structs.UserStatusInitialized,
	}
	if err = createUser(data); err != nil {
		return
	}
	return data.UserId, nil
}

func createUser(data structs.User) (err *t7Error.Error) {
	// TODO: check mobile or email used
	dbErr := db.New().CreateUser(data)
	if dbErr != nil {
		log.Error("fail to create user: ", err.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}

	return verifyUser(data.UserId)
}

// TODO: add some verify logic or send user verify signal to other service
// TODO: once user email confirmed, make it active
func verifyUser(userId string) (err *t7Error.Error) {
	log.Debug("verify user: ", userId)

	if dbErr := db.New().UpdateUserStatus(userId, structs.UserStatusActivate); dbErr != nil {
		log.Error("fail to update user status: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}

func UpdateBasicInfo(userId string, data structs.UserInfo) (err *t7Error.Error) {
	if dbErr := db.New().UpdateUserBasicInfo(userId, data); dbErr != nil {
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
	}
	return
}
