package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	Id        string `gorm:"primaryKey;column:id;type:varchar(36);not null"`
	UserId    string `gorm:"uniqueIndex:user_id;column:userId;type:varchar(36);not null"`
	CreatedAt int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli"`

	Balance []Balance `gorm:"foreignKey:WalletId;references:Id"`
}

func (w Wallet) TableName() string {
	return "wallet"
}

type Balance struct {
	WalletId  string `gorm:"primaryKey;type:varchar(36);not_null"`
	Money     `gorm:"embedded"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli"`
}

func (b *Balance) TableName() string {
	return "balance"
}

type Money struct {
	Currency string          `gorm:"primaryKey;type:varchar(4);uniqueIndex:wc;not null"`
	Amount   decimal.Decimal `gorm:"type:decimal(16,4),not null;default:0"`
}
