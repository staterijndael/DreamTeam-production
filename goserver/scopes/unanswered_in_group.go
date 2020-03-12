package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func UnansweredInGroup(me, eid, gid uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		gr := models.Group{
			Model: gorm.Model{
				ID: gid,
			},
		}

		estimate := models.Estimate{
			EstimatorID: me,
			GroupID:     gid,
			EventID:     eid,
		}

		return db.
			Model(&models.User{}).
			Where("id <> ?", me).
			Where(
				"id in (?)",
				db.
					Model(&models.MembershipOfCommunity{}).
					Where(
						"community = ?",
						db.
							Model(&gr).
							Where(&gr).
							Select("community").
							SubQuery(),
					).
					Select("\"user\"").
					SubQuery(),
			).
			Where(
				"id not in (?)",
				db.
					Model(&estimate).
					Where(&estimate).
					Select("estimated").
					SubQuery(),
			).
			Where(
				"(?) > 0",
				db.
					Model(&models.RatingEvent{}).
					Where("start < current_timestamp").
					Where(`"end" > current_timestamp`).
					Where("id = ?", eid).
					Select("count(*)").
					SubQuery(),
			)
	}
}
