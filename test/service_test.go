package test

import (
	"context"
	"fmt"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	viper.AddConfigPath("../config")
}

func Test_service_CreateUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), "traceId", "Test_service_CreateUser")

	req := userV1.CreateUserRequest{
		Username: "admin",
		Password: "password",
		Role:     authV1.Role_admin,
	}

	authAvc := auth.New(newTestDbClient(), newTestDbCore(), newTestCache(), logger.GetLogger(), newTestConfig())

	userId, err := authAvc.CreateUser(ctx, &req)
	if err != nil {
		t.Error(err)
	}

	if userId == "" {
		t.Error("empty user id")
	}

	var data []models.User
	testDbCore.Find(&data)
	fmt.Println(data)
}

func Test_service_GenActivationCode(t *testing.T) {
	ctx := context.WithValue(context.Background(), "traceId", "Test_service_GenActivationCode")
	authAvc := auth.New(newTestDbClient(), newTestDbCore(), newTestCache(), logger.GetLogger(), newTestConfig())

	code, err := authAvc.GenActivationCode(ctx, "testUserId")
	if err != nil {
		t.Error(err)
	}

	if code == "" {
		t.Error("empty act code")
	}
}
