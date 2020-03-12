//go:generate zenrpc
package user

import (
	"dt/managers/sms"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db  *gorm.DB
	sms *sms.Manager
} //zenrpc

func New(db *gorm.DB, manager *sms.Manager) *Service {
	return &Service{
		db:  db,
		sms: manager,
	}
}
