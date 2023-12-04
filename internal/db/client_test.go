package db

import (
	"github.com/Template7/backend/internal/db/entity"
	"github.com/Template7/common/db"
	"github.com/spf13/viper"
	"testing"
)

func TestNew(t *testing.T) {
	viper.AddConfigPath("../../config")

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
