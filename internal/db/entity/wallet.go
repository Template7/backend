package entity

import (
	"context"
	"github.com/Template7/common/logger"
	v1 "github.com/Template7/protobuf/gen/proto/template7/wallet"
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	Id        string    `gorm:"type:char(36);primary_key;"`
	UserId    string    `gorm:"uniqueIndex:user_id;type:varchar(36);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`

	Balance []Balance `gorm:"foreignKey:WalletId;references:Id"`
}

func (w *Wallet) TableName() string {
	return "wallet"
}

func (w *Wallet) ToProto(ctx context.Context) *v1.Wallet {
	log := logger.New().WithContext(ctx)

	var blc []*v1.Balance
	for _, b := range w.Balance {
		if _, ok := v1.Currency_value[b.Currency]; !ok {
			log.With("currency", b.Currency).Warn("unsupported currency")
			continue
		}
		blc = append(blc, &v1.Balance{
			Currency: v1.Currency(v1.Currency_value[b.Currency]),
			Amount:   b.Amount.String(),
		})
	}
	return &v1.Wallet{
		Id:       w.Id,
		Balances: blc,
	}
}

type Balance struct {
	WalletId  string `gorm:"primaryKey;type:char(36);not_null"`
	Money     `gorm:"embedded"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli;not null"`
}

func (b *Balance) TableName() string {
	return "balance"
}

type Money struct {
	Currency string          `gorm:"primaryKey;type:varchar(4);not null"`
	Amount   decimal.Decimal `gorm:"type:decimal(16,4);not null;default:0"`
}
