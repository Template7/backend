package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type DepositHistory struct {
	Id            int64           `gorm:"primaryKey"` // snowflake id
	UserId        string          `gorm:"index:u;type:char(36);not_null"`
	WalletId      string          `gorm:"index:wc;type:char(36);not_null"`
	Currency      string          `gorm:"index:wc;type:varchar(4);not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note          string
	CreatedAt     time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *DepositHistory) TableName() string {
	return "deposit_history"
}

type WithdrawHistory struct {
	Id            int64           `gorm:"primaryKey"` // snowflake id
	UserId        string          `gorm:"index:u;type:char(36);not_null"`
	WalletId      string          `gorm:"index:wc;type:char(36);not_null"`
	Currency      string          `gorm:"index:wc;type:varchar(4);not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note          string
	CreatedAt     time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *WithdrawHistory) TableName() string {
	return "withdraw_history"
}

type TransferHistory struct {
	Id                    int64           `gorm:"primaryKey"` // snowflake id
	UserId                string          `gorm:"index:u;type:char(36);not_null"`
	FromWalletId          string          `gorm:"index:fc;type:char(36);not_null"`
	ToWalletId            string          `gorm:"index:wc;type:char(36);not_null"`
	Currency              string          `gorm:"index:fc;index:wc;type:varchar(4);not null"`
	Amount                decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	SenderBalanceBefore   decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	SenderBalanceAfter    decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	ReceiverBalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	ReceiverBalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note                  string
	CreatedAt             time.Time `gorm:"autoCreateTime:milli;not null"`
}

func (d *TransferHistory) TableName() string {
	return "transfer_history"
}

type WalletHistory struct {
	Currency string
	Balance  WalletBalanceHistory
}

type WalletBalanceHistory struct {
	RecordId     int64
	Io           string // in/out
	Amount       decimal.Decimal
	AmountBefore decimal.Decimal
	AmountAfter  decimal.Decimal
	Timestamp    time.Time // record created_at
	Note         string
}
