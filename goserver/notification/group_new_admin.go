package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type NewGroupAdmin struct {
	notificationBase
	Group      *models.Group
	OldAdmin   *models.User
	NewAdmin   *models.User
	AssignedBy *models.User
} //notifier

func (nga *NewGroupAdmin) loadReceivers() {
	if nga.Group.Parent != nil {
		nga.receivers = append(nga.receivers, nga.Group.Parent.AdminID)
	}

	nga.receivers = append(nga.receivers, nga.Group.Community.MembersIDs()...)
}

func (nga *NewGroupAdmin) loadDashReceivers() {
	nga.dashReceivers = nga.Group.Organization.Admins.MembersIDs()
}

func (nga *NewGroupAdmin) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.newgroupadmin",
		Data: nga.View(),
	}
}

func (nga *NewGroupAdmin) View() interface{} {
	return &struct {
		ID         uint         `json:"id"`
		Group      *views.Group `json:"group"`
		OldAdmin   *views.User  `json:"oldAdmin"`
		NewAdmin   *views.User  `json:"newAdmin"`
		AssignedBy *views.User  `json:"assignedBy"`
		Seen       *bool        `json:"seen,omitempty"`
	}{
		ID:         nga.GetModel().ID,
		OldAdmin:   views.UserViewFromModel(nga.OldAdmin),
		NewAdmin:   views.UserViewFromModel(nga.NewAdmin),
		AssignedBy: views.UserViewFromModel(nga.AssignedBy),
		Group:      views.GroupFromModelShort(nga.Group),
		Seen:       nga.seen,
	}
}

func (nga *NewGroupAdmin) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.NewGroupAdmin)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = nga.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, nga.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, nga.AssignedBy.ID)
	if err != nil {
		return err
	}

	_, err = saveUNS(db, n, nga.Group.Community.MembersIDs())

	if nga.Group.ParentID != nil {
		_, err := saveSingleUNS(db, n, nga.Group.Parent.AdminID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (nga *NewGroupAdmin) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.NewGroupAdmin
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return nga.LoadWithEvent(db, e, n)
}

func (nga *NewGroupAdmin) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.NewGroupAdmin
	var ok bool
	if event, ok = _event.(*events.NewGroupAdmin); !ok {
		return WrongEventErr
	}

	var oldAdmin models.User
	if err := db.First(&oldAdmin, event.OldAdmin).Error; err != nil {
		return err
	}

	var newAdmin models.User
	if err := db.First(&newAdmin, event.NewAdmin).Error; err != nil {
		return err
	}

	var assignedBy models.User
	if event.AssignedBy == event.OldAdmin {
		assignedBy = oldAdmin
	} else {
		if err := db.First(&assignedBy, event.AssignedBy).Error; err != nil {
			return err
		}
	}

	var group models.Group
	if err := db.
		Preload("Parent").
		First(&group, event.Group).Error; err != nil {
		return err
	}

	nga.Model = model
	nga.Group = &group
	nga.AssignedBy = &assignedBy
	nga.OldAdmin = &oldAdmin
	nga.NewAdmin = &newAdmin
	nga.loadReceivers()
	nga.loadDashReceivers()

	return nil
}
