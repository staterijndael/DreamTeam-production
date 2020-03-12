package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func UnseenAOWS(uid, notificationID uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		orgWall := models.OrgWall{NotificationID: notificationID}
		return db.
			Model(&models.AdminOrgWallSeen{}).
			Where(`"user" = ?`, uid).
			Where("seen = ?", false).
			Where(
				"org_wall = ?",
				db.
					Model(&orgWall).
					Where(&orgWall).
					Select("id").
					SubQuery(),
			)
	}
}
