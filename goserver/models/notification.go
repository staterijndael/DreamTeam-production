package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
	"strings"
)

type Notification struct {
	gorm.Model
	Data     postgres.Jsonb `gorm:"column:data"`
	DataType string         `gorm:"column:data_type"`
}

type UserNotificationSeen struct {
	gorm.Model
	UserID         uint         `gorm:"column:user"`
	NotificationID uint         `gorm:"column:notification"`
	Seen           bool         `gorm:"column:seen"`
	User           User         `gorm:"foreignkey:user"`
	Notification   Notification `gorm:"foreignkey:notification"`
}

func (*Notification) TableName() string {
	return "notifications"
}

func (*UserNotificationSeen) TableName() string {
	return "user_notification_seens"
}

func NotificationWithEvent(event interface{}) (*Notification, error) {
	n := &Notification{}
	return n, n.SetEventData(event)
}

func (n *Notification) SetEventData(event interface{}) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	n.Data = postgres.Jsonb{RawMessage: bytes}
	value := reflect.ValueOf(event)
	kind := value.Kind()
	if reflect.Ptr == kind {
		value = value.Elem()
	}

	n.DataType = strings.ToLower(value.Type().Name())
	return nil
}
