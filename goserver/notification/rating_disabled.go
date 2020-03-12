package notification

import (
	"dt/events"
	"dt/models"
	"dt/scopes"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type RatingDisabled struct {
	notificationBase
	groupMembers []uint               `json:"-"`
	Org          *models.Organization `json:"organization"`
} //notifier

func (rd *RatingDisabled) loadReceivers() {
	rd.receivers = append(rd.Org.Admins.MembersIDs(), rd.groupMembers...)
}

func (rd *RatingDisabled) loadDashReceivers() {
	rd.dashReceivers = append(rd.Org.Admins.MembersIDs(), rd.groupMembers...)
}

func (rd *RatingDisabled) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.ratingdisabled",
		Data: rd.View(),
	}
}

func (rd *RatingDisabled) View() interface{} {
	return &struct {
		ID     uint                `json:"id"`
		Org    *views.Org          `json:"organization"`
		Config *views.RatingConfig `json:"config"`
		Seen   *bool               `json:"seen,omitempty"`
	}{
		ID:   rd.GetModel().ID,
		Org:  views.OrgViewFromModelShort(rd.Org),
		Seen: rd.seen,
	}
}

func (rd *RatingDisabled) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.RatingDisabled)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = rd.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	_, err = saveWallEvent(db, n, rd.Org.ID)
	if err != nil {
		return err
	}

	return nil
}

func (rd *RatingDisabled) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.RatingDisabled
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return rd.LoadWithEvent(db, e, n)
}

func (rd *RatingDisabled) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.RatingDisabled
	var ok bool
	if event, ok = _event.(*events.RatingDisabled); !ok {
		return WrongEventErr
	}

	var org models.Organization
	if err := db.First(&org, event.Organization).Error; err != nil {
		return err
	}

	if err := db.Scopes(scopes.GroupMembersIDsOfOrg(event.Organization, &rd.groupMembers)).Error; err != nil {
		return err
	}

	rd.Model = model
	rd.Org = &org
	rd.loadReceivers()
	rd.loadDashReceivers()

	return nil
}
