package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type OrgJoinRequestDenied struct {
	notificationBase
	Request *models.OrgJoinRequest
} //notifier

func (notif *OrgJoinRequestDenied) loadReceivers() {
	notif.receivers = append(notif.receivers, notif.Request.InitiatorID)

	for _, member := range notif.Request.Organization.Admins.Members {
		if member.UserID == *notif.Request.AcceptorID {
			continue
		}

		notif.receivers = append(notif.receivers, member.UserID)
	}
}

func (notif *OrgJoinRequestDenied) loadDashReceivers() {
	for _, member := range notif.Request.Organization.Admins.Members {
		if member.UserID == *notif.Request.AcceptorID {
			continue
		}

		notif.dashReceivers = append(notif.dashReceivers, member.UserID)
	}
}

func (notif *OrgJoinRequestDenied) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.orgjoinrequestdenied",
		Data: notif.View(),
	}
}

func (notif *OrgJoinRequestDenied) View() interface{} {
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

func (notif *OrgJoinRequestDenied) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgJoinRequestDenied)
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

	if _, err = saveAOWSExcept(db, wall, *notif.Request.AcceptorID); err != nil {
		return err
	}

	_, err = saveUNS(db, n, notif.receivers)
	if err != nil {
		return err
	}

	return nil
}

func (notif *OrgJoinRequestDenied) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgJoinRequestDenied
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *OrgJoinRequestDenied) LoadWithEvent(
	db *gorm.DB,
	_event interface{},
	model *models.Notification,
) error {
	var event *events.OrgJoinRequestDenied
	var ok bool
	if event, ok = _event.(*events.OrgJoinRequestDenied); !ok {
		return WrongEventErr
	}

	var request models.OrgJoinRequest
	if err := db.First(&request, event.Request).Error; err != nil {
		return err
	}

	notif.Request = &request
	notif.Model = model
	notif.loadReceivers()
	notif.loadDashReceivers()

	return nil
}
