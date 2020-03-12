package models

import "github.com/jinzhu/gorm"

type Community struct {
	gorm.Model
	Members []MembershipOfCommunity `gorm:"foreignkey:community"`
}

func (*Community) TableName() string {
	return "communities"
}

func (c *Community) Contains(uid uint) bool {
	for _, m := range c.Members {
		if m.User.ID == uid {
			return true
		}
	}

	return false
}

func (c *Community) MembersIDs() []uint {
	ids := make([]uint, len(c.Members))
	for i := range c.Members {
		ids[i] = c.Members[i].UserID
	}

	return ids
}

type MembershipOfCommunity struct {
	gorm.Model
	CommunityID uint `gorm:"column:community"`
	UserID      uint `gorm:"column:user"`
	User        User `gorm:"foreignkey:user"`
}

func (*MembershipOfCommunity) TableName() string {
	return "membership_of_communities"
}
