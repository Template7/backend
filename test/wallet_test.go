package test

import (
	"context"
	"github.com/Template7/backend/internal/wallet"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"testing"
)

func Test_wallet(t *testing.T) {
	viper.AddConfigPath("../config")

	wId := "8250e425-e02c-410b-ad8a-b2d29149bf8a"
	ctx := context.WithValue(context.Background(), "traceId", uuid.NewString())
	w, err := wallet.New().GetWallet(ctx, wId)
	if err != nil {
		t.Error(err)
		return
	}

	if err := wallet.New().Deposit(ctx, wId, v1.Currency_NTD, 123); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Deposit(ctx, wId, v1.Currency_USD, 45); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Deposit(ctx, wId, v1.Currency_JPY, 678); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Deposit(ctx, wId, v1.Currency_CNY, 90); err != nil {
		t.Error(err)
		return
	}

	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_NTD, 23); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_USD, 45); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_JPY, 67); err != nil {
		t.Error(err)
		return
	}
	if err := wallet.New().Withdraw(ctx, wId, v1.Currency_CNY, 89); err != nil {
		t.Error(err)
		return
	}

	tw := "bd159a64-5a20-493b-93a0-8fcc9b0c607d"
	if err := wallet.New().Transfer(ctx, wId, tw, v1.Currency_NTD, 7); err != nil {
		t.Error(err)
		return
	}

	t.Log(w.String())
}
