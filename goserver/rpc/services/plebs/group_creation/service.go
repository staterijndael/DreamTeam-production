//go:generate zenrpc
package group_creation

import (
	"dt/managers/eventEmitter"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db      *gorm.DB
	emitter *eventEmitter.EventEmitter
} //zenrpc

func New(db *gorm.DB, ee *eventEmitter.EventEmitter) *Service {
	return &Service{db: db, emitter: ee}
}
