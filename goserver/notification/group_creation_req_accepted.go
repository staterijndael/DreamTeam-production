package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type GroupCreationRequestAccepted struct {
	notificationBase
	Group   *models.Group
	Request *models.GroupCreationRequest
} //notifier

func (req *GroupCreationRequestAccepted) loadReceivers() {
	req.receivers = append(req.receivers, req.Request.InitiatorID)

	if req.Request.Parent != nil && *req.Request.AcceptorID != req.Request.Parent.AdminID {
		req.receivers = append(req.receivers, req.Request.Parent.AdminID)
	}
}

func (req *GroupCreationRequestAccepted) loadDashReceivers() {
	for _, member := range req.Group.Organization.Admins.Members {
		if member.UserID == *req.Request.AcceptorID {
			continue
		}
		req.dashReceivers = append(req.dashReceivers, member.UserID)
	}
}

func (req *GroupCreationRequestAccepted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupcreationrequestaccepted",
		Data: req.View(),
	}
}

func (req *GroupCreationRequestAccepted) View() interface{} {
	return &struct {
		ID      uint                        `json:"id"`
		Group   *views.Group                `json:"group"`
		Request *views.GroupCreationRequest `json:"request"`
		Seen    *bool                       `json:"seen,omitempty"`
	}{
		ID:      req.GetModel().ID,
		Group:   views.GroupFromModelShort(req.Group),
		Request: views.GroupCreationRequestFromModelShort(req.Request),
		Seen:    req.seen,
	}
}

func (req *GroupCreationRequestAccepted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupCreationRequestAccepted)
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

	wall, err := saveWallEvent(db, n, req.Group.OrganizationID)
	if err != nil {
		return err
	}

	if _, err = saveAOWSExcept(db, wall, *req.Request.AcceptorID); err != nil {
		return err
	}

	if req.Request.Parent != nil && *req.Request.AcceptorID != req.Request.Parent.AdminID {
		_, err = saveUNS(db, n, []uint{
			req.Request.Parent.AdminID,
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

func (req *GroupCreationRequestAccepted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupCreationRequestAccepted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *GroupCreationRequestAccepted) LoadWithEvent(
	db *gorm.DB,
	_event interface{},
	model *models.Notification,
) error {
	var event *events.GroupCreationRequestAccepted
	var ok bool
	if event, ok = _event.(*events.GroupCreationRequestAccepted); !ok {
		return WrongEventErr
	}

	var request models.GroupCreationRequest
	if err := db.First(&request, event.Request).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.First(&group, *event.Group).Error; err != nil {
		return err
	}

	req.Group = &group
	req.Group.Parent = request.Parent
	req.Request = &request
	req.Model = model
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
