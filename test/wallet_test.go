package test

//func Test_wallet(t *testing.T) {
//	viper.AddConfigPath("../config")
//
//	wId := "ac03df62-f362-4a89-9a2e-e3cc0d129ea4"
//	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
//	w, err := wallet.New().GetWallet(ctx, wId)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	if err := wallet.New().Deposit(ctx, wId, v1.Currency_ntd, 123); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Deposit(ctx, wId, v1.Currency_ntd, 45); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Deposit(ctx, wId, v1.Currency_jpy, 678); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Deposit(ctx, wId, v1.Currency_cny, 90); err != nil {
//		t.Error(err)
//		return
//	}
//
//	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_ntd, 23); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_usd, 45); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_jpy, 67); err != nil {
//		t.Error(err)
//		return
//	}
//	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_cny, 89); err != nil {
//		t.Error(err)
//		return
//	}
//
//	tw := "bd159a64-5a20-493b-93a0-8fcc9b0c607d"
//	if err := wallet.New().Transfer(ctx, wId, tw, v1.Currency_ntd, 7); err != nil {
//		t.Error(err)
//		return
//	}
//
//	t.Log(w.String())
//}

//func Test_transfer(t *testing.T) {
//	viper.AddConfigPath("../config")
//
//	fw := "ac03df62-f362-4a89-9a2e-e3cc0d129ea4"
//	tw := "bee58b8b-a035-423d-93b5-1d668a94f05c"
//	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
//
//	//wallet.New().Deposit(ctx, fw, v1.Currency_NTD, 7)
//
//	if err := wallet.New().Transfer(ctx, fw, tw, v1.Currency_ntd, 7); err != nil {
//		t.Error(err)
//		return
//	}
//}
