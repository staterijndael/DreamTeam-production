package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func FindUser(page uint64, queries []string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		user := models.User{}
		nickname := models.Nickname{}
		return db.
			Model(&user).
			Where(user.FuzzyQuery(queries)).
			Or(
				"nickname in (?)",
				db.
					Table(nickname.TableName()).
					Select("id").
					Where(nickname.FuzzyQuery(queries)).
					SubQuery(),
			).
			Offset(page * searchPageSize).
			Limit(searchPageSize)
	}
}
