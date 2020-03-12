package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type GroupCreationRequestWithdrawn struct {
	notificationBase
	Request *models.GroupCreationRequest
} //notifier

func (req *GroupCreationRequestWithdrawn) loadReceivers() {
	if req.Request.Parent != nil {
		req.receivers = append(req.receivers, req.Request.Parent.AdminID)
	}
}

func (req *GroupCreationRequestWithdrawn) loadDashReceivers() {
	req.dashReceivers = req.Request.Organization.Admins.MembersIDs()
}

func (req *GroupCreationRequestWithdrawn) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupcreationrequestwithdrawn",
		Data: req.View(),
	}
}

func (req *GroupCreationRequestWithdrawn) View() interface{} {
	return &struct {
		ID      uint                        `json:"id"`
		Request *views.GroupCreationRequest `json:"request"`
		Seen    *bool                       `json:"seen,omitempty"`
	}{
		ID:      req.GetModel().ID,
		Request: views.GroupCreationRequestFromModelShort(req.Request),
		Seen:    req.seen,
	}
}

func (req *GroupCreationRequestWithdrawn) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupCreationRequestWithdrawn)
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

	wall, err := saveWallEvent(db, n, req.Request.OrganizationID)
	if err != nil {
		return err
	}

	if _, err = saveAOWS(db, wall); err != nil {
		return err
	}

	if req.Request.Parent != nil {
		_, err = saveSingleUNS(db, n, req.Request.Parent.AdminID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (req *GroupCreationRequestWithdrawn) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupCreationRequestWithdrawn
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *GroupCreationRequestWithdrawn) LoadWithEvent(db *gorm.DB, event interface{}, n *models.Notification) error {
	var e *events.GroupCreationRequestWithdrawn
	if _event, ok := event.(*events.GroupCreationRequestWithdrawn); !ok {
		return WrongEventErr
	} else {
		e = _event
	}

	var request models.GroupCreationRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	req.Request = &request
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
