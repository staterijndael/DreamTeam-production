package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/golang-collections/collections/set"
	"github.com/jinzhu/gorm"
)

type OrgDeleted struct {
	notificationBase
	Organization *models.Organization
	Members      []uint
} //notifier

func (notif *OrgDeleted) loadReceivers() {
	notif.receivers = notif.Members
}

func (notif *OrgDeleted) loadDashReceivers() {
}

func (notif *OrgDeleted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.orgdeleted",
		Data: notif.View(),
	}
}

func (notif *OrgDeleted) View() interface{} {
	return &struct {
		ID   uint       `json:"id"`
		Org  *views.Org `json:"org"`
		Seen *bool      `json:"seen,omitempty"`
	}{
		ID:   notif.GetModel().ID,
		Org:  views.OrgViewFromModelShort(notif.Organization),
		Seen: notif.seen,
	}
}

func (notif *OrgDeleted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgDeleted)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = notif.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	if _, err = saveUNS(db, n, notif.Members); err != nil {
		return nil
	}

	return nil
}

func (notif *OrgDeleted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgDeleted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *OrgDeleted) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.OrgDeleted
	var ok bool
	if event, ok = _event.(*events.OrgDeleted); !ok {
		return WrongEventErr
	}

	var org models.Organization
	if err := db.Unscoped().First(&org, event.Org).Error; err != nil {
		return err
	}

	s := set.New()
	for _, id := range event.GroupMembers {
		if id == org.DirectorID {
			continue
		}

		s.Insert(id)
	}

	for _, id := range org.Admins.MembersIDs() {
		if id == org.DirectorID {
			continue
		}

		s.Insert(id)
	}

	s.Do(func(el interface{}) {
		notif.Members = append(notif.Members, el.(uint))
	})

	notif.Model = model
	notif.Organization = &org
	notif.loadReceivers()
	notif.loadDashReceivers()

	return nil
}
