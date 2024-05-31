package test

//func Test_user(t *testing.T) {
//	viper.AddConfigPath("../config")
//
//	userId := "115e8d14-df86-4096-a367-4c5fba94621e"
//	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
//	token, err := auth.New().Login(ctx, "allentest", "password")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log(token)
//	info, err := user.New().GetInfo(ctx, userId)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log(info.String())
//	logger.New().WithContext(ctx).With("info", info).Debug("show info")
//
//	uw := user.New().GetUserWallets(ctx, userId)
//	t.Log(uw)
//
//	err = user.New().UpdateInfo(ctx, userId, entity.UserInfo{
//		Nickname: "nickname",
//	})
//	if err != nil {
//		t.Error(err)
//		return
//	}
//}

//func Test_service_CreateUser(t *testing.T) {
//	viper.AddConfigPath("../config")
//
//	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
//
//	req := userV1.CreateUserRequest{
//		Username: "admin",
//		Password: "password",
//		Role:     authV1.Role_admin,
//	}
//	if err := auth.New().CreateUser(ctx, &req); err != nil {
//		t.Error(err)
//	}
//}
