package models

import "github.com/jinzhu/gorm"

type OrgWall struct {
	gorm.Model
	NotificationID uint         `gorm:"column:notification"`
	OrganizationID uint         `gorm:"column:organization"`
	Notification   Notification `gorm:"foreignkey:notification"`
	Organization   Organization `gorm:"foreignkey:organization"`
}

func (*OrgWall) TableName() string {
	return "org_walls"
}

type AdminOrgWallSeen struct {
	gorm.Model
	UserID uint    `gorm:"column:user"`
	WallID uint    `gorm:"column:org_wall"`
	Seen   bool    `gorm:"column:seen"`
	User   User    `gorm:"foreignkey:user"`
	Wall   OrgWall `gorm:"foreignkey:org_wall"`
}

func (*AdminOrgWallSeen) TableName() string {
	return "admin_org_wall_seens"
}
