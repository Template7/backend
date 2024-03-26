package auth

import (
	"context"
	authV1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	userV1 "github.com/Template7/protobuf/gen/proto/template7/user"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"testing"
)

func Test_service_CreateUser(t *testing.T) {
	viper.AddConfigPath("../../config")

	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())

	req := userV1.CreateUserRequest{
		Username: "admin",
		Password: "password",
		Role:     authV1.Role_admin,
	}
	if _, err := New().CreateUser(ctx, &req); err != nil {
		t.Error(err)
	}
}
