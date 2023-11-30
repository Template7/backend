package entity

import (
	"time"
)

type Wallet struct {
	Id        string    `gorm:"type:char(36);primary_key;"`
	UserId    string    `gorm:"uniqueIndex:user_id;type:varchar(36);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

func (w *Wallet) TableName() string {
	return "wallet"
}
