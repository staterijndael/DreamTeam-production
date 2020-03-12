//go:generate zenrpc
package group_join

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

func New(db *gorm.DB, c *config.Config, ee *eventEmitter.EventEmitter) *Service {
	return &Service{
		db:      db,
		conf:    c,
		emitter: ee,
	}
}
