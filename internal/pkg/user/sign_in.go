package user

import (
	"github.com/Template7/backend/internal/pkg/auth"
	"github.com/Template7/backend/internal/pkg/db"
	"github.com/Template7/backend/internal/pkg/t7Error"
	"github.com/Template7/common/structs"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

const (
	SignInPrefix = "signIn"
)

func MobileSignIn(mobile string) (userToken structs.Token, err *t7Error.Error) {
	log.Debug("mobile sign in")

	data, dbErr := db.New().GetUserByMobile(mobile)

	// TODO: decouple from db implementation?
	if dbErr == mongo.ErrNoDocuments {
		log.Error("user not found")
		err = t7Error.UserNotfound.WithStatus(http.StatusNoContent)
		return
	}
	if dbErr != nil {
		log.Error("fail to get user data: ", err.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}
	return auth.GenUserToken(data.UserId)
}