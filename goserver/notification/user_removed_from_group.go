package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type UserRemovedFromGroup struct {
	notificationBase
	Group     *models.Group
	Removed   *models.User
	RemovedBy *models.User
} //notifier

func (urfg *UserRemovedFromGroup) loadReceivers() {
	for _, member := range urfg.Group.Community.MembersIDs() {
		if member == urfg.RemovedBy.ID {
			continue
		}

		urfg.receivers = append(urfg.receivers, member)
	}

	urfg.receivers = append(urfg.receivers, urfg.Removed.ID)
}

func (urfg *UserRemovedFromGroup) loadDashReceivers() {
	for _, member := range urfg.Group.Organization.Admins.MembersIDs() {
		if member == urfg.RemovedBy.ID {
			continue
		}

		urfg.dashReceivers = append(urfg.dashReceivers, member)
	}
}

func (urfg *UserRemovedFromGroup) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.userremovedfromgroup",
		Data: urfg.View(),
	}
}

func (urfg *UserRemovedFromGroup) View() interface{} {
	return &struct {
		ID        uint         `json:"id"`
		Group     *views.Group `json:"group"`
		Removed   *views.User  `json:"removed"`
		RemovedBy *views.User  `json:"removedBy"`
		Seen      *bool        `json:"seen,omitempty"`
	}{
		ID:        urfg.GetModel().ID,
		Group:     views.GroupFromModelShort(urfg.Group),
		Removed:   views.UserViewFromModel(urfg.Removed),
		RemovedBy: views.UserViewFromModel(urfg.RemovedBy),
		Seen:      urfg.seen,
	}
}

func (urfg *UserRemovedFromGroup) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.UserRemovedFromGroup)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = urfg.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, urfg.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, e.RemovedBy)
	if err != nil {
		return err
	}

	for _, member := range urfg.Group.Community.Members {
		if member.UserID == e.RemovedBy {
			continue
		}

		if _, err = saveSingleUNS(db, n, member.UserID); err != nil {
			return err
		}
	}

	if _, err = saveSingleUNS(db, n, urfg.Removed.ID); err != nil {
		return err
	}

	return nil
}

func (urfg *UserRemovedFromGroup) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.UserRemovedFromGroup
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return urfg.LoadWithEvent(db, e, n)
}

func (urfg *UserRemovedFromGroup) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.UserRemovedFromGroup
	var ok bool
	if event, ok = _event.(*events.UserRemovedFromGroup); !ok {
		return WrongEventErr
	}

	var added models.User
	if err := db.First(&added, event.Removed).Error; err != nil {
		return err
	}

	var addedBy models.User
	if err := db.First(&addedBy, event.RemovedBy).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.
		First(&group, event.Group).Error; err != nil {
		return err
	}

	urfg.Model = model
	urfg.Removed = &added
	urfg.RemovedBy = &addedBy
	urfg.Group = &group
	urfg.loadReceivers()
	urfg.loadDashReceivers()

	return nil
}
