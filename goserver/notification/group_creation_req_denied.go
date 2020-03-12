package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
)

type GroupCreationRequestDenied struct {
	notificationBase
	Request *models.GroupCreationRequest
} //notifier

func (req *GroupCreationRequestDenied) loadReceivers() {
	req.receivers = append(req.receivers, req.Request.InitiatorID)

	if req.Request.Parent != nil && *req.Request.AcceptorID != req.Request.Parent.AdminID {
		req.receivers = append(req.receivers, req.Request.Parent.AdminID)
	}
}

func (req *GroupCreationRequestDenied) loadDashReceivers() {
	for _, member := range req.Request.Organization.Admins.Members {
		if member.UserID == *req.Request.AcceptorID {
			continue
		}

		req.dashReceivers = append(req.dashReceivers, member.UserID)
	}
}

func (req *GroupCreationRequestDenied) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupcreationrequestdenied",
		Data: req.View(),
	}
}

func (req *GroupCreationRequestDenied) View() interface{} {
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

func (req *GroupCreationRequestDenied) CreateByEvent(db *gorm.DB, _event interface{}) error {
	e, ok := _event.(*events.GroupCreationRequestDenied)
	if !ok {
		return errors.New("wrong event")
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

	if _, err = saveAOWSExcept(db, wall, *req.Request.AcceptorID); err != nil {
		return err
	}

	if req.Request.Parent != nil && *req.Request.AcceptorID != req.Request.Parent.AdminID {
		_, err = saveUNS(db, n, []uint{
			*req.Request.ParentID,
			req.Request.InitiatorID,
		})
	} else {
		_, err = saveSingleUNS(db, n, req.Request.InitiatorID)
	}

	if err != nil {
		return err
	}

	return nil
}

func (req *GroupCreationRequestDenied) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupCreationRequestDenied
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *GroupCreationRequestDenied) LoadWithEvent(
	db *gorm.DB,
	event interface{},
	n *models.Notification,
) error {
	e, ok := event.(*events.GroupCreationRequestDenied)
	if !ok {
		return nil
	}

	var request models.GroupCreationRequest
	if err := db.First(&request, e.Request).Error; err != nil {
		return err
	}

	req.Model = n
	req.Request = &request
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
