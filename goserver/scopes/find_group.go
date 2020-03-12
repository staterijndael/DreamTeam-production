package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

const searchPageSize uint64 = 50

func FindGroup(page uint64, oid uint, queries []string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		group := models.Group{}
		nickname := models.Nickname{}
		return db.
			Model(&group).
			Where("organization = ? and "+group.FuzzyQuery(queries), oid).
			Or(
				"nickname in (?)",
				db.
					Model(&nickname).
					Select("id").
					Where(nickname.FuzzyQuery(queries)).
					SubQuery(),
			).
			Offset(page * searchPageSize).
			Limit(searchPageSize)
	}
}
