package notification

import (
	"dt/events"
	"dt/models"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type OrgNewDirector struct {
	notificationBase
	Org      *models.Organization
	OldAdmin *models.User
	NewAdmin *models.User
} //notifier

func (ond *OrgNewDirector) loadReceivers() {}

func (ond *OrgNewDirector) loadDashReceivers() {
	ond.dashReceivers = ond.Org.Admins.MembersIDs()
}

func (ond *OrgNewDirector) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.newdirector",
		Data: ond.View(),
	}
}

func (ond *OrgNewDirector) View() interface{} {
	return &struct {
		ID       uint        `json:"id"`
		Org      *views.Org  `json:"organization"`
		OldAdmin *views.User `json:"oldAdmin"`
		NewAdmin *views.User `json:"newAdmin"`
		Seen     *bool       `json:"seen,omitempty"`
	}{
		ID:       ond.GetModel().ID,
		Org:      views.OrgViewFromModelShort(ond.Org),
		OldAdmin: views.UserViewFromModel(ond.OldAdmin),
		NewAdmin: views.UserViewFromModel(ond.NewAdmin),
		Seen:     ond.seen,
	}
}

func (ond *OrgNewDirector) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.OrgNewDirector)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = ond.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	wall, err := saveWallEvent(db, n, ond.Org.ID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, ond.OldAdmin.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ond *OrgNewDirector) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.OrgNewDirector
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return ond.LoadWithEvent(db, e, n)
}

func (ond *OrgNewDirector) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.OrgNewDirector
	var ok bool
	if event, ok = _event.(*events.OrgNewDirector); !ok {
		return WrongEventErr
	}

	var oldAdmin models.User
	if err := db.First(&oldAdmin, event.OldDirector).Error; err != nil {
		return err
	}

	var newAdmin models.User
	if err := db.First(&newAdmin, event.NewDirector).Error; err != nil {
		return err
	}

	var org models.Organization
	if err := db.First(&org, event.Org).Error; err != nil {
		return err
	}

	ond.Model = model
	ond.OldAdmin = &oldAdmin
	ond.NewAdmin = &newAdmin
	ond.Org = &org
	ond.loadReceivers()
	ond.loadDashReceivers()

	return nil
}
