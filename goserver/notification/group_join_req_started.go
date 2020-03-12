package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type GroupJoinRequestStarted struct {
	notificationBase
	Request *models.GroupJoinRequest
} //notifier

func (req *GroupJoinRequestStarted) loadReceivers() {
	req.receivers = append(req.receivers, req.Request.Group.AdminID)
}

func (req *GroupJoinRequestStarted) loadDashReceivers() {
	req.dashReceivers = req.Request.Group.Organization.Admins.MembersIDs()
}

func (req *GroupJoinRequestStarted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupjoinrequeststarted",
		Data: req.View(),
	}
}

func (req *GroupJoinRequestStarted) View() interface{} {
	return &struct {
		ID      uint                    `json:"id"`
		Request *views.GroupJoinRequest `json:"request"`
		Seen    *bool                   `json:"seen,omitempty"`
	}{
		ID:      req.GetModel().ID,
		Request: views.GroupJoinRequestFromModelShort(req.Request),
		Seen:    req.seen,
	}
}

func (req *GroupJoinRequestStarted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupJoinRequestStarted)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = req.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, req.Request.Group.OrganizationID)
	if err != nil {
		return err
	}

	if _, err = saveAOWS(db, wall); err != nil {
		return err
	}

	if _, err = saveSingleUNS(db, n, req.Request.Group.AdminID); err != nil {
		return err
	}

	return nil
}

func (req *GroupJoinRequestStarted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupJoinRequestStarted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *GroupJoinRequestStarted) LoadWithEvent(db *gorm.DB, event interface{}, n *models.Notification) error {
	var e *events.GroupJoinRequestStarted
	if _event, ok := event.(*events.GroupJoinRequestStarted); !ok {
		return WrongEventErr
	} else {
		e = _event
	}

	var request models.GroupJoinRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	req.Request = &request
	req.Model = n
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
