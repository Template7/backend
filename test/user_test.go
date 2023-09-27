package test

import (
	"context"
	"github.com/Template7/backend/internal/auth"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/backend/internal/user"
	"github.com/Template7/common/logger"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"testing"
)

func Test_user(t *testing.T) {
	viper.AddConfigPath("../config")

	userId := "115e8d14-df86-4096-a367-4c5fba94621e"
	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
	token, err := auth.New().Login(ctx, "allentest", "password")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(token)
	info, err := user.New().GetInfo(ctx, userId)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(info.String())
	logger.New().WithContext(ctx).With("info", info).Debug("show info")

	uw := user.New().GetUserWallets(ctx, userId)
	t.Log(uw)

	err = user.New().UpdateInfo(ctx, userId, entity.UserInfo{
		Nickname: "nickname",
	})
	if err != nil {
		t.Error(err)
		return
	}
}
