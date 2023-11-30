package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type WalletBalance struct {
	WalletId string
	Currency string
	Amount   decimal.Decimal
}

type UserWalletBalance struct {
	UserId   string
	WalletId string
	Currency string
	Amount   decimal.Decimal
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
