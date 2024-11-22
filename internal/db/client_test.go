package db

import (
	"context"
	"encoding/json"
	"github.com/Template7/common/db"
	"github.com/Template7/common/logger"
	"github.com/Template7/common/models"
	"github.com/spf13/viper"
	"testing"
)

// ignore for automation test

func init() {
	viper.AddConfigPath("../../config")
}

func TestNew(t *testing.T) {
	if err := db.NewSql().Debug().AutoMigrate(
		&models.User{},
		&models.Wallet{},
		&models.Balance{},
		&models.DepositHistory{},
		&models.WithdrawHistory{},
		&models.TransferHistory{}); err != nil {
		t.Error(err)
	}
}

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
