package entity

import (
	v1 "github.com/Template7/protobuf/gen/proto/template7/auth"
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `gorm:"type:char(36);primary_key;"`
	Username  string    `gorm:"type:varchar(64);uniqueIndex;not_null"`
	Password  string    `gorm:"type:varchar(64);not_null"`
	Info      UserInfo  `gorm:"embedded"`
	Email     string    `gorm:"type:varchar(128)"`
	Status    v1.AccountStatus
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli;not null"`
}

// UserInfo
// the fields which configurable by the user
type UserInfo struct {
	NickName string `gorm:"type:varchar(64);not_null;default:''"`
}

func (u User) TableName() string {
	return "user"
}
