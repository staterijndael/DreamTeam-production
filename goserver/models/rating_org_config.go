package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type RatingOrgConfig struct {
	gorm.Model
	OrganizationID uint          `gorm:"column:org_id;unique"`
	StartTime      uint8         `gorm:"column:start_time;default:17"`
	WeekDay        time.Weekday  `gorm:"column:week_day;default:0"`
	Organization   *Organization `gorm:"foreignkey:org_id"`
}

func (*RatingOrgConfig) TableName() string {
	return "rating_org_configs"
}
