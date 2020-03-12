package scopes

import (
	"dt/models"
	"dt/stores"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

func CreateGroup(
	createdGroupAddress *models.Group,
	creator,
	org uint,
	title,
	description string,
	nickname models.Nickname,
	parent *uint,
) Scope {
	return func(db *gorm.DB) *gorm.DB {
		chat := models.Chat{
			Community: models.Community{
				Members: []models.MembershipOfCommunity{
					{UserID: creator},
				},
			},
		}

		if _db := db.Create(&chat); _db.Error != nil {
			return _db
		}

		gr := models.Group{
			ParentID:       parent,
			CreatorID:      creator,
			AdminID:        creator,
			Title:          title,
			Description:    description,
			NicknameID:     nickname.ID,
			Nickname:       nickname,
			OrganizationID: org,
			CommunityID:    chat.CommunityID,
			ChatID:         chat.ID,
			ChildrenIDs:    make(pq.Int64Array, 0),
			AvatarID:       stores.DefaultAvatars.Group.ID,
		}

		if _db := db.Create(&gr); _db.Error != nil {
			return _db
		}

		if _db := db.First(&gr, gr.ID); _db.Error != nil {
			return _db
		}

		_db := db
		if parent != nil {
			_db = db.Exec("UPDATE groups SET children = array_append(children, ?) WHERE id = ?", gr.ID, *parent)
		}

		if _db.Error == nil {
			*createdGroupAddress = gr
		}

		return _db
	}
}
