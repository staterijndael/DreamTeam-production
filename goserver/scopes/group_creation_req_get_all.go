package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func GroupCreationReqGetAll(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.GroupCreationRequest{}).
			Where("status = 'pending'").
			Where("organization in (?) OR hm in (?)",
				db.
					Model(&models.Organization{}).
					Where("community in (?)",
						db.
							Model(&models.MembershipOfCommunity{}).
							Where("user = ?", uid).
							Select("community").
							SubQuery(),
					).
					Select("id").
					SubQuery(),
				db.
					Model(&models.Group{}).
					Where("admin = ?", uid).
					Select("id").
					SubQuery(),
			)

	}
}
