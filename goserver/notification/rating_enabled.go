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

type RatingEnabled struct {
	notificationBase
	groupMembers []uint                  `json:"-"`
	Org          *models.Organization    `json:"organization"`
	Config       *models.RatingOrgConfig `json:"config"`
} //notifier

func (re *RatingEnabled) loadReceivers() {
	re.receivers = append(re.Org.Admins.MembersIDs(), re.groupMembers...)
}

func (re *RatingEnabled) loadDashReceivers() {
	re.dashReceivers = append(re.Org.Admins.MembersIDs(), re.groupMembers...)
}

func (re *RatingEnabled) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.ratingenabled",
		Data: re.View(),
	}
}

func (re *RatingEnabled) View() interface{} {
	return &struct {
		ID     uint                `json:"id"`
		Org    *views.Org          `json:"organization"`
		Config *views.RatingConfig `json:"config"`
		Seen   *bool               `json:"seen,omitempty"`
	}{
		ID:     re.GetModel().ID,
		Org:    views.OrgViewFromModelShort(re.Org),
		Config: views.RatingConfigFromModel(re.Config),
		Seen:   re.seen,
	}
}

func (re *RatingEnabled) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.RatingEnabled)
	if !ok {
		return WrongEventErr
	}

	n, err := saveNotification(db, e)
	if err != nil {
		return err
	}

	if err = re.LoadWithEvent(db, e, n); err != nil {
		return err
	}

	_, err = saveWallEvent(db, n, re.Org.ID)
	if err != nil {
		return err
	}

	return nil
}

func (re *RatingEnabled) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.RatingEnabled
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return re.LoadWithEvent(db, e, n)
}

func (re *RatingEnabled) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.RatingEnabled
	var ok bool
	if event, ok = _event.(*events.RatingEnabled); !ok {
		return WrongEventErr
	}

	var org models.Organization
	if err := db.First(&org, event.Organization).Error; err != nil {
		return err
	}

	var c models.RatingOrgConfig
	if err := db.First(&c, event.Config).Error; err != nil {
		return err
	}

	if err := db.Scopes(scopes.GroupMembersIDsOfOrg(event.Organization, &re.groupMembers)).Error; err != nil {
		return err
	}

	re.Model = model
	re.Config = &c
	re.Org = &org
	re.loadReceivers()
	re.loadDashReceivers()

	return nil
}
