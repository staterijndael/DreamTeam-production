package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func GroupMembersIDsOfOrg(oid uint, members *[]uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		type IDs struct {
			ID uint `gorm:"column:id"`
		}

		var memIDs []*IDs
		gr := models.Group{OrganizationID: oid}
		if res := db.
			Table("users").
			Model(&models.User{}).
			Select("id").
			Where(
				"id in (?)",
				db.
					Model(&models.MembershipOfCommunity{}).
					Select("\"user\"").
					Where(
						"community in (?)",
						db.
							Model(&models.Community{}).
							Select("id").
							Where(
								"id in (?)",
								db.
									Model(&gr).
									Select("community").
									Where(&gr).
									SubQuery(),
							).
							SubQuery(),
					).
					SubQuery(),
			).Find(&memIDs); res.Error != nil {
			return res
		}

		for _, id := range memIDs {
			*members = append(*members, id.ID)
		}

		return db
	}
}
