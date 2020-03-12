package models

import "github.com/jinzhu/gorm"

type Timeline struct {
	gorm.Model
	GroupID        uint         `gorm:"column:group"`
	NotificationID uint         `gorm:"column:notification"`
	EventID        *uint        `gorm:"column:event"`
	Group          Group        `gorm:"foreignkey:group"`
	Notification   Notification `gorm:"foreignkey:notification"`
	Event          *RatingEvent `gorm:"foreignkey:event"`
}

func (*Timeline) TableName() string {
	return "timelines"
}
