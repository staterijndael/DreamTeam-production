package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
	"time"
)

func GetMessagesByChatBeforeTime(cid uint, after *time.Time) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.Message{}).
			Where("chat = ?", cid).
			Where("created_at < ?", after).
			Order("id desc").
			Limit(searchPageSize)
	}
}

func GetMessagesByChat(cid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(&models.Message{ChatID: cid}).Order("id desc").Limit(50)
	}
}
