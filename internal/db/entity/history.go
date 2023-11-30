package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type DepositHistory struct {
	Id        int64           `gorm:"primaryKey"` // snowflake id
	UserId    string          `gorm:"primaryKey;type:char(36);not_null"`
	WalletId  string          `gorm:"primaryKey;type:char(36);not_null"`
	Currency  string          `gorm:"primaryKey;type:varchar(4);not null"`
	Amount    decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note      string
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *DepositHistory) TableName() string {
	return "deposit_history"
}

type WithdrawHistory struct {
	Id        int64           `gorm:"primaryKey"` // snowflake id
	UserId    string          `gorm:"primaryKey;type:char(36);not_null"`
	WalletId  string          `gorm:"primaryKey;type:char(36);not_null"`
	Currency  string          `gorm:"primaryKey;type:varchar(4);not null"`
	Amount    decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note      string
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *WithdrawHistory) TableName() string {
	return "withdraw_history"
}

type TransferHistory struct {
	Id           int64           `gorm:"primaryKey"` // snowflake id
	UserId       string          `gorm:"primaryKey;type:char(36);not_null"`
	FromWalletId string          `gorm:"primaryKey;type:char(36);not_null"`
	ToWalletId   string          `gorm:"primaryKey;type:char(36);not_null"`
	Currency     string          `gorm:"primaryKey;type:varchar(4);not null"`
	Amount       decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note         string
	CreatedAt    time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *TransferHistory) TableName() string {
	return "transfer_history"
}
