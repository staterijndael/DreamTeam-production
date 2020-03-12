package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func GroupJoinReqGetAll(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.GroupJoinRequest{}).
			Where("status = 'pending'").
			Where(`initiator = ? OR "group" in ?`,
				uid,
				db.
					Model(&models.Group{}).
					Where("admin = ? OR organization in ?",
						uid,
						db.
							Model(&models.Organization{}).
							Where("community in ?",
								db.
									Model(&models.MembershipOfCommunity{}).
									Where(`"user" = ?`, uid).
									Select("community").
									SubQuery(),
							).
							Select("id").
							SubQuery(),
					).
					Select("id").
					SubQuery(),
			)
	}
}
