package db

import (
	"context"
	"encoding/json"
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/db"
	"github.com/Template7/common/logger"
	"github.com/spf13/viper"
	"testing"
)

func init() {
	viper.AddConfigPath("../../config")
}

func TestNew(t *testing.T) {
	if err := db.NewSql().Debug().AutoMigrate(
		&entity.User{},
		&entity.Wallet{},
		&entity.Balance{},
		&entity.DepositHistory{},
		&entity.WithdrawHistory{},
		&entity.TransferHistory{}); err != nil {
		t.Error(err)
	}
}

// ignore for automation test
func TestClient_GetWalletBalanceHistory(t *testing.T) {
	c := New(logger.New())

	ctx := context.WithValue(context.Background(), "traceId", "TestClient_GetWalletBalanceHistory")
	data, err := c.GetWalletBalanceHistory(ctx, "9c5d98c7-65a4-4e97-83ef-feb3969ef421")
	if err != nil {
		t.Error(err)
	}

	b, _ := json.MarshalIndent(data, "", "  ")
	t.Log(string(b))
}
