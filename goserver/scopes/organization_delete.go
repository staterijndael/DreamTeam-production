package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func OrganizationDelete(org models.Organization, groups []*models.Group) Scope {
	return func(db *gorm.DB) *gorm.DB {
		var exMems []*models.OrgExMember
		if res := db.
			Model(models.OrgExMember{}).
			Where("organization = ?", org.ID).
			Find(&exMems); res.Error != nil {
			return res
		}

		for _, exMem := range exMems {
			if res := db.Delete(&exMem); res.Error != nil {
				return res
			}
		}

		if res := db.Where("org_id = ?", org.ID).Delete(&models.RatingOrgConfig{}); res.Error != nil {
			return res
		}

		for _, gr := range groups {
			if res := db.
				Where("community = ?", gr.CommunityID).
				Delete(&models.MembershipOfCommunity{}); res.Error != nil {
				return res
			}

			if res := db.Scopes(GroupDelete(*gr)); res.Error != nil {
				return res
			}

			if res := db.Where("id = ?", gr.CommunityID).Delete(&models.Community{}); res.Error != nil {
				return res
			}
		}

		if res := db.Delete(&org); res.Error != nil {
			return res
		}

		if res := db.Where("id = ?", org.NicknameID).Delete(&models.Nickname{}); res.Error != nil {
			return res
		}

		return db
	}
}
