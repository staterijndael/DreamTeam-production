//go:generate zenrpc
package token

import (
	"dt/config"
	"dt/managers/sms"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db   *gorm.DB
	sms  *sms.Manager
	conf *config.Config
} //zenrpc

func New(db *gorm.DB, manager *sms.Manager, conf *config.Config) *Service {
	return &Service{
		db:   db,
		sms:  manager,
		conf: conf,
	}
}
