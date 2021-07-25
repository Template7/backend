package admin

import (
	"backend/internal/pkg/auth"
	"backend/internal/pkg/db"
	"backend/internal/pkg/db/collection"
	"backend/internal/pkg/t7Error"
	"backend/internal/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func SignIn(data collection.Admin) (token collection.Token, err *t7Error.Error) {
	log.Debug("admin sign in")

	adminData, dbErr := db.New().GetAdmin()
	if dbErr != nil {
		log.Error("fail to get admin: ", dbErr.Error())
		err = t7Error.DbOperationFail.WithDetailAndStatus(dbErr.Error(), http.StatusInternalServerError)
		return
	}


	//hashed, hashErr := util.HashedPassword(data.Password)
	//if hashErr != nil {
	//	log.Error("fail to hash password: ", hashErr.Error())
	//	err = t7Error.HashFail.WithDetailAndStatus(hashErr.Error(), http.StatusInternalServerError)
	//	return
	//}

	if util.CheckPasswordHash(data.Password, adminData.Password)  || data.Username != adminData.Username {
		log.Warn("invalid admin username or password")
		err = t7Error.SignInFail.WithStatus(http.StatusForbidden)
		return
	}

	return auth.GenAdminToken()
}



