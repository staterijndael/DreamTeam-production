package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func RemoveFromGroup(uid uint, group *models.Group) Scope {
	return func(db *gorm.DB) *gorm.DB {
		var membership *models.MembershipOfCommunity
		for _, mem := range group.Community.Members {
			if mem.UserID == uid {
				membership = &mem
				break
			}
		}

		return db.Delete(&membership)
	}
}
