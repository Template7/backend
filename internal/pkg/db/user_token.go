package db

import (
	"backend/internal/pkg/db/collection"
	"context"
	log "github.com/sirupsen/logrus"
)

func (c client) SaveUserToken(token collection.UserToken) (err error) {
	log.Debug("save user token")
	_, err = c.userToken.InsertOne(context.Background(), token)
	return
}
