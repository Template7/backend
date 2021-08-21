package db

import (
	"github.com/Template7/common/structs"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func(c client) GetAdmin() (data structs.Admin, err error) {
	log.Debug("get admin")
	err = c.admin.FindOne(context.Background(), bson.M{}).Decode(&data)
	return
}
