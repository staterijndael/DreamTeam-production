package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type UserLeftGroup struct {
	notificationBase
	Group *models.Group
	User  *models.User
} //notifier

func (ulg *UserLeftGroup) loadReceivers() {
	for _, member := range ulg.Group.Community.MembersIDs() {
		if member == ulg.User.ID {
			continue
		}

		ulg.receivers = append(ulg.receivers, member)
	}

}

func (ulg *UserLeftGroup) loadDashReceivers() {
	for _, member := range ulg.Group.Organization.Admins.Members {
		if member.UserID == ulg.User.ID {
			continue
		}

		ulg.dashReceivers = append(ulg.dashReceivers, member.UserID)
	}
}

func (ulg *UserLeftGroup) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.userleftgroup",
		Data: ulg.View(),
	}
}

func (ulg *UserLeftGroup) View() interface{} {
	return &struct {
		ID    uint         `json:"id"`
		Group *views.Group `json:"group"`
		User  *views.User  `json:"user"`
		Seen  *bool        `json:"seen,omitempty"`
	}{
		ID:    ulg.GetModel().ID,
		Group: views.GroupFromModelShort(ulg.Group),
		User:  views.UserViewFromModel(ulg.User),
		Seen:  ulg.seen,
	}
}

func (ulg *UserLeftGroup) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.UserLeftGroup)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = ulg.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, ulg.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, e.User)
	if err != nil {
		return err
	}

	for _, member := range ulg.Group.Community.Members {
		if member.UserID == e.User {
			continue
		}

		if _, err = saveSingleUNS(db, n, member.UserID); err != nil {
			return err
		}
	}

	return nil
}

func (ulg *UserLeftGroup) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.UserLeftGroup
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return ulg.LoadWithEvent(db, e, n)
}

func (ulg *UserLeftGroup) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.UserLeftGroup
	var ok bool
	if event, ok = _event.(*events.UserLeftGroup); !ok {
		return WrongEventErr
	}

	var user models.User
	if err := db.First(&user, event.User).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.
		First(&group, event.Group).Error; err != nil {
		return err
	}

	ulg.Model = model
	ulg.User = &user
	ulg.Group = &group
	ulg.loadReceivers()
	ulg.loadDashReceivers()

	return nil
}
