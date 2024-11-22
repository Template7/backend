package test

import (
	"context"
	"github.com/Template7/common/models"
	"github.com/Template7/common/t7Id"
	walletV1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDbClient_GetWalletBalanceHistory(t *testing.T) {

	dbCore := newTestDbCore()
	wId := uuid.NewString()
	dh := models.DepositHistory{
		Id:            t7Id.New().Generate().Int64(),
		WalletId:      wId,
		Currency:      walletV1.Currency_ntd.String(),
		Amount:        decimal.NewFromFloat(123),
		BalanceBefore: decimal.NewFromFloat(0),
		BalanceAfter:  decimal.NewFromFloat(123),
		Note:          "test",
	}
	if err := dbCore.Create(&dh).Error; err != nil {
		t.Error(err)
	}

	ctx := context.WithValue(context.Background(), "traceId", "TestDbClient_GetWalletBalanceHistory")
	db := newTestDbClient()
	data, err := db.GetWalletBalanceHistory(ctx, wId)
	if err != nil {
		t.Error(err)
	}

	t.Log(data)
}
