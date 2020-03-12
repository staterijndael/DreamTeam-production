package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func RatingEventsOfGroupsOfUser(uid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&models.RatingEvent{}).
			Where("start < current_timestamp").
			Where(`"end" > current_timestamp`).
			Where(
				"organization in (?)",
				db.
					Scopes(GroupsOfUser(uid)).
					Select("organization").
					SubQuery(),
			)
	}
}
