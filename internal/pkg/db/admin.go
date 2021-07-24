package db

import (
	"backend/internal/pkg/db/collection"
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func(c client) GetAdmin() (data collection.Admin, err error) {
	log.Debug("get admin")
	err = c.admin.FindOne(context.Background(), bson.M{}).Decode(&data)
	return
}
