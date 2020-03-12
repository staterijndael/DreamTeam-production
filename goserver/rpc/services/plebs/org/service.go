//go:generate zenrpc
package org

import (
	"dt/config"
	"dt/managers/eventEmitter"
	"dt/managers/fns"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db      *gorm.DB
	conf    *config.Config
	emitter *eventEmitter.EventEmitter
	fnsMgr  *fns.FNSManager
} //zenrpc

func New(
	store *gorm.DB,
	conf *config.Config,
	fns *fns.FNSManager,
	ee *eventEmitter.EventEmitter,
) *Service {
	return &Service{db: store, conf: conf, fnsMgr: fns, emitter: ee}
}
