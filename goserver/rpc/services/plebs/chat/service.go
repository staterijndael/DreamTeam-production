//go:generate zenrpc
package chat

import (
	"dt/managers/eventEmitter"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db      *gorm.DB
	emitter *eventEmitter.EventEmitter
} //zenrpc

func New(db *gorm.DB, emitter *eventEmitter.EventEmitter) *Service {
	return &Service{
		db:      db,
		emitter: emitter,
	}
}
