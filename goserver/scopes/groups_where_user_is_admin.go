package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func GroupsWhereUserIsAdmin(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.Group{}).
			Where("admin = ?", uid)
	}
}
