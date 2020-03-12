package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type OrgJoinRequestStarted struct {
	notificationBase
	Request *models.OrgJoinRequest
} //notifier

func (notif *OrgJoinRequestStarted) loadReceivers() {
	notif.receivers = notif.Request.Organization.Admins.MembersIDs()
}

func (notif *OrgJoinRequestStarted) loadDashReceivers() {
	notif.dashReceivers = notif.Request.Organization.Admins.MembersIDs()
}

func (notif *OrgJoinRequestStarted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.orgjoinrequeststarted",
		Data: notif.View(),
	}
}

func (notif *OrgJoinRequestStarted) View() interface{} {
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

func (notif *OrgJoinRequestStarted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgJoinRequestStarted)
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

func (notif *OrgJoinRequestStarted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgJoinRequestStarted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *OrgJoinRequestStarted) LoadWithEvent(db *gorm.DB, event interface{}, n *models.Notification) error {
	var e *events.OrgJoinRequestStarted
	if _event, ok := event.(*events.OrgJoinRequestStarted); !ok {
		return WrongEventErr
	} else {
		e = _event
	}

	var request models.OrgJoinRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	notif.Request = &request
	notif.Model = n
	notif.loadReceivers()
	notif.loadDashReceivers()

	return nil
}
