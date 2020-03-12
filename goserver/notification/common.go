//go:generate notifier
package notification

import (
	"dt/models"
	"dt/utils"
	"errors"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

var (
	WrongEventErr = errors.New("wrong event")
)

type INotification interface {
	GetModel() *models.Notification
	CreateByEvent(db *gorm.DB, event interface{}) error
	Load(db *gorm.DB, model *models.Notification) error
	LoadWithEvent(db *gorm.DB, event interface{}, model *models.Notification) error
	Receivers() []uint
	DashReceivers() []uint
	ContainerizedView() *utils.Container
	View() interface{}
	Seen(state bool)
}

var (
	notifications map[string]func() INotification //notifier
)

func getByTypeName(name string) (INotification, error) {
	if f, ok := notifications[name]; ok {
		return f(), nil
	}

	return nil, errors.New("loader not implemented")
}

func GetNotification(data interface{}) (INotification, error) {
	var typeName string
	switch converted := data.(type) {
	case *models.Notification:
		typeName = converted.DataType
		break
	case models.Notification:
		typeName = converted.DataType
		break
	default:
		value := reflect.ValueOf(converted)
		kind := value.Kind()
		if reflect.Ptr == kind {
			value = value.Elem()
		}

		typeName = strings.ToLower(value.Type().Name())
	}

	return getByTypeName(typeName)
}
