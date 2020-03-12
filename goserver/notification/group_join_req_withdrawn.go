package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type GroupJoinRequestWithdrawn struct {
	notificationBase
	Request *models.GroupJoinRequest
} //notifier

func (req *GroupJoinRequestWithdrawn) loadReceivers() {
	req.receivers = append(req.receivers, req.Request.Group.AdminID)
}

func (req *GroupJoinRequestWithdrawn) loadDashReceivers() {
	req.dashReceivers = req.Request.Group.Organization.Admins.MembersIDs()
}

func (req *GroupJoinRequestWithdrawn) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupjoinrequestwithdrawn",
		Data: req.View(),
	}
}

func (req *GroupJoinRequestWithdrawn) View() interface{} {
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

func (req *GroupJoinRequestWithdrawn) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupJoinRequestWithdrawn)
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

func (req *GroupJoinRequestWithdrawn) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupJoinRequestWithdrawn
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *GroupJoinRequestWithdrawn) LoadWithEvent(db *gorm.DB, event interface{}, n *models.Notification) error {
	var e *events.GroupJoinRequestWithdrawn
	if _event, ok := event.(*events.GroupJoinRequestWithdrawn); !ok {
		return WrongEventErr
	} else {
		e = _event
	}

	var request models.GroupJoinRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	req.Model = n
	req.Request = &request
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
