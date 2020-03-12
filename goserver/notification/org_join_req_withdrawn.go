package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type OrgJoinRequestWithdrawn struct {
	notificationBase
	Request *models.OrgJoinRequest
} //notifier

func (notif *OrgJoinRequestWithdrawn) loadReceivers() {
	notif.receivers = notif.Request.Organization.Admins.MembersIDs()
}

func (notif *OrgJoinRequestWithdrawn) loadDashReceivers() {
	notif.receivers = notif.Request.Organization.Admins.MembersIDs()
}

func (notif *OrgJoinRequestWithdrawn) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.orgjoinrequestwithdrawn",
		Data: notif.View(),
	}
}

func (notif *OrgJoinRequestWithdrawn) View() interface{} {
	return &struct {
		ID      uint                  `json:"id"`
		Request *views.OrgJoinRequest `json:"request"`
		Seen    *bool                 `json:"seen,omitempty"`
	}{
		ID:      notif.GetModel().ID,
		Request: views.OrgJoinRequestFromModel(notif.Request),
		Seen:    notif.seen,
	}
}

func (notif *OrgJoinRequestWithdrawn) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgJoinRequestWithdrawn)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = notif.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, notif.Request.OrganizationID)
	if err != nil {
		return err
	}

	if _, err = saveAOWS(db, wall); err != nil {
		return err
	}

	if _, err = saveUNS(db, n, notif.receivers); err != nil {
		return err
	}

	return nil
}

func (notif *OrgJoinRequestWithdrawn) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgJoinRequestWithdrawn
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *OrgJoinRequestWithdrawn) LoadWithEvent(db *gorm.DB, event interface{}, n *models.Notification) error {
	var e *events.OrgJoinRequestWithdrawn
	if _event, ok := event.(*events.OrgJoinRequestWithdrawn); !ok {
		return WrongEventErr
	} else {
		e = _event
	}

	var request models.OrgJoinRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	notif.Model = n
	notif.Request = &request
	notif.loadReceivers()
	notif.loadDashReceivers()

	return nil
}
