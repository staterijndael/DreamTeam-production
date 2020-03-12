//go:generate zenrpc
package user

import (
	"dt/config"
	"dt/managers/eventEmitter"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db      *gorm.DB
	conf    *config.Config
	emitter *eventEmitter.EventEmitter
} //zenrpc

func New(store *gorm.DB, conf *config.Config, ee *eventEmitter.EventEmitter) *Service {
	return &Service{db: store, conf: conf, emitter: ee}
}
