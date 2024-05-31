package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type DepositHistory struct {
	Id            int64           `gorm:"primaryKey"` // snowflake id
	WalletId      string          `gorm:"index:wcc;type:char(36);not_null"`
	Currency      string          `gorm:"index:wcc;type:varchar(4);not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note          string
	CreatedAt     time.Time `gorm:"index:wcc;autoCreateTime:milli;not null"`
}

func (d *DepositHistory) TableName() string {
	return "deposit_history"
}

type WithdrawHistory struct {
	Id            int64           `gorm:"primaryKey"` // snowflake id
	WalletId      string          `gorm:"index:wc;type:char(36);not_null"`
	Currency      string          `gorm:"index:wc;type:varchar(4);not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	BalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note          string
	CreatedAt     time.Time `gorm:"index:wcc;autoCreateTime:milli;not null"`
}

func (d *WithdrawHistory) TableName() string {
	return "withdraw_history"
}

type TransferHistory struct {
	Id                    int64           `gorm:"primaryKey"` // snowflake id
	FromWalletId          string          `gorm:"index:fcc;type:char(36);not_null"`
	ToWalletId            string          `gorm:"index:wcc;type:char(36);not_null"`
	Currency              string          `gorm:"index:fcc;index:wcc;type:varchar(4);not null"`
	Amount                decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	SenderBalanceBefore   decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	SenderBalanceAfter    decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	ReceiverBalanceBefore decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	ReceiverBalanceAfter  decimal.Decimal `gorm:"type:decimal(16,4);not null"`
	Note                  string
	CreatedAt             time.Time `gorm:"index:fcc;index:wcc;autoCreateTime:milli;not null"`
}

func (d *TransferHistory) TableName() string {
	return "transfer_history"
}

type WalletHistory struct {
	Currency string
	Balance  WalletBalanceHistory
}

// WalletBalanceHistory
// data structure only
type WalletBalanceHistory struct {
	Id            uint64
	Direction     string
	Currency      string
	Amount        decimal.Decimal
	BalanceBefore decimal.Decimal
	BalanceAfter  decimal.Decimal
	CreatedAt     time.Time
	Note          string
}
