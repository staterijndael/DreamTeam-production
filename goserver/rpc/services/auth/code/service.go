//go:generate zenrpc
package code

import (
	"dt/managers/connections"
	"dt/managers/sms"
	"github.com/jinzhu/gorm"
)

type Service struct {
	db  *gorm.DB
	sms *sms.Manager
	cm  *connections.UserConnectionsManager
	dcm *connections.DashConnectionsManager
} //zenrpc

func New(
	db *gorm.DB,
	smsMgr *sms.Manager,
	cm *connections.UserConnectionsManager,
	dcm *connections.DashConnectionsManager,
) *Service {
	return &Service{
		db:  db,
		sms: smsMgr,
		cm:  cm,
		dcm: dcm,
	}
}
