package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type OrgJoinRequestAccepted struct {
	notificationBase
	Request *models.OrgJoinRequest
} //notifier

func (req *OrgJoinRequestAccepted) loadReceivers() {
	mems := make(map[uint]*models.MembershipOfCommunity)
	for _, member := range req.Request.Group.Community.Members {
		if member.UserID == *req.Request.AcceptorID {
			continue
		}

		mems[member.ID] = &member
	}

	for _, member := range req.Request.Organization.Admins.Members {
		if member.UserID == *req.Request.AcceptorID {
			continue
		}

		mems[member.ID] = &member
	}

	req.receivers = make([]uint, 0, len(mems))
	for _, val := range mems {
		req.receivers = append(req.receivers, val.UserID)
	}
}

func (req *OrgJoinRequestAccepted) loadDashReceivers() {
	for _, member := range req.Request.Organization.Admins.Members {
		if member.UserID == *req.Request.AcceptorID {
			continue
		}

		req.dashReceivers = append(req.dashReceivers, member.UserID)
	}
}

func (req *OrgJoinRequestAccepted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "Notification.orgjoinrequestaccepted",
		Data: req.View(),
	}
}

func (req *OrgJoinRequestAccepted) View() interface{} {
	return &struct {
		ID      uint                  `json:"id"`
		Request *views.OrgJoinRequest `json:"request"`
		Seen    *bool                 `json:"seen,omitempty"`
	}{
		ID:      req.GetModel().ID,
		Request: views.OrgJoinRequestFromModel(req.Request),
		Seen:    req.seen,
	}
}

func (req *OrgJoinRequestAccepted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgJoinRequestAccepted)
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

	if _, err = saveAOWSExcept(db, wall, *req.Request.AcceptorID); err != nil {
		return err
	}

	_, err = saveUNS(db, n, req.receivers)
	if err != nil {
		return err
	}

	return nil
}

func (req *OrgJoinRequestAccepted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgJoinRequestAccepted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return req.LoadWithEvent(db, e, n)
}

func (req *OrgJoinRequestAccepted) LoadWithEvent(
	db *gorm.DB,
	_event interface{},
	model *models.Notification,
) error {
	var event *events.OrgJoinRequestAccepted
	var ok bool
	if event, ok = _event.(*events.OrgJoinRequestAccepted); !ok {
		return WrongEventErr
	}

	var request models.OrgJoinRequest
	if err := db.First(&request, event.Request).Error; err != nil {
		return err
	}

	req.Request = &request
	req.Model = model
	req.loadReceivers()
	req.loadDashReceivers()

	return nil
}
