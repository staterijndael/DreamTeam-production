package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

func IDsOfGroupsOfUser(uid uint, ids *[]uint) Scope {
	return func(db *gorm.DB) *gorm.DB {
		type GroupsIDs struct {
			Id uint `gorm:"column:id"`
		}

		var groups []*GroupsIDs
		if res := db.
			Scopes(GroupsOfUser(uid)).
			Select("id").
			Table("groups").
			Model(&models.Group{}).
			Scan(&groups); res.Error != nil {
			return res
		}

		for _, gr := range groups {
			*ids = append(*ids, gr.Id)
		}

		return db
	}
}
