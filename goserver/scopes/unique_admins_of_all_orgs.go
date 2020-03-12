package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func UniqueAdminsOfAllOrgs(orgsIDs []uint, uniqueAdmins *[]uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		type IDs struct {
			Id uint `gorm:"column:user"`
		}

		var admins []*IDs
		if res := db.
			Table("membership_of_communities").
			Model(&models.MembershipOfCommunity{}).
			Select(`distinct("user")`).
			Where("community in (?)",
				db.
					Model(&models.Organization{}).
					Select("community").
					Where("id in (?)", orgsIDs).
					SubQuery(),
			).Scan(&admins); res.Error != nil {
			return res
		}

		for _, admin := range admins {
			*uniqueAdmins = append(*uniqueAdmins, admin.Id)
		}

		return db
	}
}
