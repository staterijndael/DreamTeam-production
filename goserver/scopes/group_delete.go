package scopes

import (
	"dt/models"
	"github.com/jinzhu/gorm"
)

//requires transaction
func GroupDelete(gr models.Group) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if res := db.Delete(&gr); res.Error != nil {
			return res
		}

		if res := db.Where("id = ?", gr.NicknameID).Delete(&models.Nickname{}); res.Error != nil {
			return res
		}

		if res := db.Where("id = ?", gr.ChatID).Delete(&models.Chat{}); res.Error != nil {
			return res
		}

		return db
	}
}
