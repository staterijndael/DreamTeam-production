package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func MoveToParent(gr models.Group) Scope {
	return func(db *gorm.DB) *gorm.DB {
		grChildrens := []int64(gr.ChildrenIDs)
		if res := db.
			Exec(
				`UPDATE "groups" SET parent = ? WHERE id in (?)`,
				gr.ParentID,
				grChildrens,
			); res.Error != nil {
			return res
		}

		if gr.ParentID == nil {
			return db
		}

		var parentGr models.Group
		if res := db.
			Set("gorm:auto_preload", false).
			Where("id = ?", gr.ParentID).
			First(&parentGr); res.Error != nil {
			return res
		}

		childs := make(pq.Int64Array, 0)
		for _, id := range parentGr.ChildrenIDs {
			if uint(id) == gr.ID {
				continue
			}

			childs = append(childs, id)
		}

		for _, id := range grChildrens {
			childs = append(childs, id)
		}

		if res := db.
			Exec(
				`UPDATE "groups" SET children = ? WHERE id = ?`,
				childs,
				gr.ParentID,
			); res.Error != nil {
			return res
		}

		return db
	}
}
