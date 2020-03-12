package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type RatingEvent struct {
	gorm.Model
	Start          time.Time      `gorm:"column:start"`
	End            time.Time      `gorm:"column:end"`
	OrganizationID uint           `gorm:"column:organization"`
	Questions      postgres.Jsonb `gorm:"column:questions"`
	Organization   Organization   `gorm:"foreignkey:organization"`
}

func (*RatingEvent) TableName() string {
	return "rating_events"
}

type Estimate struct {
	gorm.Model
	Estimate       postgres.Jsonb `gorm:"column:estimate"`
	EstimatorID    uint           `gorm:"column:estimator"`
	EstimatedID    uint           `gorm:"column:estimated"`
	GroupID        uint           `gorm:"column:group"`
	OrganizationID uint           `gorm:"column:organization"`
	EventID        uint           `gorm:"column:event"`
	Group          Group          `gorm:"foreignkey:group"`
	Event          RatingEvent    `gorm:"foreignkey:event"`
	Estimator      User           `gorm:"foreignkey:estimator"`
	Estimated      User           `gorm:"foreignkey:estimated"`
}

func (*Estimate) TableName() string {
	return "estimates"
}

type RatingEventAverageEstimate struct {
	gorm.Model
	UserID          uint           `gorm:"column:user"`
	GroupID         uint           `gorm:"column:group"`
	EventID         uint           `gorm:"column:event"`
	OrganizationID  uint           `gorm:"column:organization"`
	User            User           `gorm:"foreignkey:user"`
	Group           Group          `gorm:"foreignkey:group"`
	Event           RatingEvent    `gorm:"foreignkey:event"`
	AverageEstimate postgres.Jsonb `gorm:"column:average_estimate"`
}

func (*RatingEventAverageEstimate) TableName() string {
	return "rating_event_average_estimates"
}
