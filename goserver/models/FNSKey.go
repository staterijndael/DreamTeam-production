package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type FNSKey struct {
	gorm.Model
	Key        string    `gorm:"column:key"`
	ExpireDate time.Time `gorm:"column:expire_date"`
}

func (*FNSKey) TableName() string {
	return "fns_keys"
}
