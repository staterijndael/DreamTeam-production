package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func UserDelete(user models.User) Scope {
	return func(db *gorm.DB) *gorm.DB {
		tx := db.Begin()
		if res := db.Where(`"user" = ?`, user.ID).Delete(&models.MembershipOfCommunity{}); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := db.
			Where("initiator = ?", user.ID).
			Where("status = 'pending'").
			Delete(&models.OrgJoinRequest{}); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := db.
			Where("initiator = ?", user.ID).
			Where("status = 'pending'").
			Delete(&models.GroupJoinRequest{}); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := db.
			Where("initiator = ?", user.ID).
			Where("status = 'pending'").
			Delete(&models.GroupCreationRequest{}); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := db.Delete(&user.Nickname); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := db.Delete(&user); res.Error != nil {
			tx.Rollback()
			return res
		}

		if res := tx.Commit(); res.Error != nil {
			tx.Rollback()
			return res
		}

		return db
	}
}
