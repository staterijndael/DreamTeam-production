//go:generate zenrpc
package rating

import (
	"dt/config"
	"dt/managers/eventEmitter"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db      *gorm.DB
	emitter *eventEmitter.EventEmitter
	conf    *config.Config
} //zenrpc

func New(db *gorm.DB, ee *eventEmitter.EventEmitter, conf *config.Config) *Service {
	return &Service{db: db, emitter: ee, conf: conf}
}
