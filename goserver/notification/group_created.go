package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type GroupCreated struct {
	notificationBase
	Group   *models.Group
	Creator *models.User
} //notifier

func (gc *GroupCreated) loadReceivers() {
	if gc.Group.Parent != nil && gc.Group.Parent.AdminID != gc.Creator.ID {
		gc.receivers = append(gc.receivers, gc.Group.Parent.AdminID)
	}
}

func (gc *GroupCreated) loadDashReceivers() {
	for _, member := range gc.Group.Organization.Admins.Members {
		if member.UserID == gc.Creator.ID {
			continue
		}

		gc.dashReceivers = append(gc.dashReceivers, member.UserID)
	}
}

func (gc *GroupCreated) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupcreated",
		Data: gc.View(),
	}
}

func (gc *GroupCreated) View() interface{} {
	return &struct {
		ID      uint         `json:"id"`
		Group   *views.Group `json:"group"`
		Creator *views.User  `json:"creator"`
		Seen    *bool        `json:"seen,omitempty"`
	}{
		ID:      gc.GetModel().ID,
		Creator: views.UserViewFromModel(gc.Creator),
		Group:   views.GroupFromModelShort(gc.Group),
		Seen:    gc.seen,
	}
}

func (gc *GroupCreated) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupCreated)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = gc.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, gc.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, e.Creator)
	if err != nil {
		return err
	}

	if gc.Group.ParentID != nil && e.Creator != gc.Group.Parent.AdminID {
		_, err := saveSingleUNS(db, n, gc.Group.Parent.AdminID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gc *GroupCreated) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupCreated
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return gc.LoadWithEvent(db, e, n)
}

func (gc *GroupCreated) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.GroupCreated
	var ok bool
	if event, ok = _event.(*events.GroupCreated); !ok {
		return WrongEventErr
	}

	var creator models.User
	if err := db.First(&creator, event.Creator).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.
		Preload("Parent").
		First(&group, event.Group).Error; err != nil {
		return err
	}

	gc.Model = model
	gc.Group = &group
	gc.Creator = &creator
	gc.loadReceivers()
	gc.loadDashReceivers()

	return nil
}
