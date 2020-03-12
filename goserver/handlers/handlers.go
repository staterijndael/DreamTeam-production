package handlers

import (
	"dt/config"
	"dt/logwrap"
	"dt/managers/connections"
	"dt/models"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var (
	handlers = []Handler{
		notifier,
		ratingHandler,
	}

	db   *gorm.DB
	cm   *connections.UserConnectionsManager
	dcm  *connections.DashConnectionsManager
	conf *config.Config
)

func InitHandlers(
	store *gorm.DB,
	connManager *connections.UserConnectionsManager,
	dashConnManager *connections.DashConnectionsManager,
	config *config.Config,
) []Handler {
	dcm = dashConnManager
	db = store
	cm = connManager
	conf = config
	return handlers
}

type Handler func(event interface{})

type IConnectionManager interface {
	Add(u *models.User, con *websocket.Conn)
	Send(uid uint, msg interface{}, exceptConn *websocket.Conn) error
}

func sendToAllMembers(
	members []uint,
	msg interface{},
	parent *websocket.Conn,
	manager IConnectionManager,
) []error {
	errors := make([]error, 0)
	sent := make([]uint, 0)
	logwrap.Debug("[handlers.sendToAllMembers]: sending to %v", members)
	for _, user := range members {
		isAlreadySent := false
		for _, id := range sent {
			if id == user {
				isAlreadySent = true
				break
			}
		}

		if isAlreadySent {
			continue
		}

		if err := manager.Send(user, msg, parent); err != nil {
			errors = append(errors, err)
		} else {
			sent = append(sent, user)
		}
	}

	return errors
}
