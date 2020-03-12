package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func OrgJoinReqGetAll(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.OrgJoinRequest{}).
			Where("status = 'pending'").
			Where("initiator = ? or organization in (?)",
				uid,
				db.
					Model(&models.Organization{}).
					Where("community in (?)",
						db.
							Model(&models.MembershipOfCommunity{}).
							Where(`"user" = ?`, uid).
							Select("community").
							SubQuery(),
					).
					Select("id").
					SubQuery(),
			)
	}
}
