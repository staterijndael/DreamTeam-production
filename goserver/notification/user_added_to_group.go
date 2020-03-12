package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type UserAddedToGroup struct {
	notificationBase
	Group   *models.Group
	Added   *models.User
	AddedBy *models.User
} //notifier

func (uatg *UserAddedToGroup) loadReceivers() {
	for _, member := range uatg.Group.Community.MembersIDs() {
		if member == uatg.AddedBy.ID {
			continue
		}

		uatg.receivers = append(uatg.receivers, member)
	}
}

func (uatg *UserAddedToGroup) loadDashReceivers() {
	for _, member := range uatg.Group.Organization.Admins.MembersIDs() {
		if member == uatg.AddedBy.ID {
			continue
		}

		uatg.dashReceivers = append(uatg.dashReceivers, member)
	}
}

func (uatg *UserAddedToGroup) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.useraddedtogroup",
		Data: uatg.View(),
	}
}

func (uatg *UserAddedToGroup) View() interface{} {
	return &struct {
		ID      uint         `json:"id"`
		Group   *views.Group `json:"group"`
		Added   *views.User  `json:"added"`
		AddedBy *views.User  `json:"addedBy"`
		Seen    *bool        `json:"seen,omitempty"`
	}{
		ID:      uatg.GetModel().ID,
		Group:   views.GroupFromModelShort(uatg.Group),
		Added:   views.UserViewFromModel(uatg.Added),
		AddedBy: views.UserViewFromModel(uatg.AddedBy),
		Seen:    uatg.seen,
	}
}

func (uatg *UserAddedToGroup) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.UserAddedToGroup)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = uatg.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, uatg.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, e.AddedBy)
	if err != nil {
		return err
	}

	for _, member := range uatg.Group.Community.Members {
		if member.UserID == e.AddedBy {
			continue
		}

		if _, err = saveSingleUNS(db, n, member.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (uatg *UserAddedToGroup) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.UserAddedToGroup
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return uatg.LoadWithEvent(db, e, n)
}

func (uatg *UserAddedToGroup) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.UserAddedToGroup
	var ok bool
	if event, ok = _event.(*events.UserAddedToGroup); !ok {
		return WrongEventErr
	}

	var added models.User
	if err := db.First(&added, event.Added).Error; err != nil {
		return err
	}

	var addedBy models.User
	if err := db.First(&addedBy, event.AddedBy).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.
		First(&group, event.Group).Error; err != nil {
		return err
	}

	uatg.Model = model
	uatg.Added = &added
	uatg.AddedBy = &addedBy
	uatg.Group = &group
	uatg.loadReceivers()
	uatg.loadDashReceivers()

	return nil
}
