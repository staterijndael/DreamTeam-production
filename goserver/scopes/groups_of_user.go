package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func GroupsOfUser(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		membership := models.MembershipOfCommunity{UserID: uid}
		return db.
			Model(&models.Group{}).
			Where(
				"community in (?)",
				db.
					Model(&membership).
					Where(&membership).
					Select("community").
					SubQuery(),
			)
	}
}
