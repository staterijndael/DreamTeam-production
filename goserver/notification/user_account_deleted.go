package notification

import (
	"dt/events"
	"dt/models"
	"dt/scopes"
	"dt/utils"
	"dt/views"
	"encoding/json"
	"github.com/golang-collections/collections/set"
	"github.com/jinzhu/gorm"
)

type UserAccountDeleted struct {
	notificationBase
	Deleted            *models.User
	Groups             []*models.Group
	UniqueAdminsOfOrgs []uint
	OrgIDs             []uint
} //notifier

func (notif *UserAccountDeleted) loadReceivers() {
	s := set.New()
	for _, gr := range notif.Groups {
		s.Insert(gr.AdminID)
	}

	s.Do(func(el interface{}) {
		notif.receivers = append(notif.receivers, el.(uint))
	})
}

func (notif *UserAccountDeleted) loadDashReceivers() {
	notif.dashReceivers = notif.UniqueAdminsOfOrgs
}

func (notif *UserAccountDeleted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.useraccountdeleted",
		Data: notif.View(),
	}
}

func (notif *UserAccountDeleted) View() interface{} {
	return &struct {
		ID      uint        `json:"id"`
		Deleted *views.User `json:"deleted"`
		Seen    *bool       `json:"seen,omitempty"`
	}{
		ID:      notif.GetModel().ID,
		Deleted: views.UserViewFromModel(notif.Deleted),
		Seen:    notif.seen,
	}
}

func (notif *UserAccountDeleted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.UserAccountDeleted)
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

	for _, org := range notif.OrgIDs {
		wall, err := saveWallEvent(db, n, org)
		if err != nil {
			return err
		}
		_, err = saveAOWS(db, wall)
		if err != nil {
			return err
		}
	}

	if _, err := saveUNS(db, n, notif.receivers); err != nil {
		return err
	}

	return nil
}

func (notif *UserAccountDeleted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.UserAccountDeleted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *UserAccountDeleted) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.UserAccountDeleted
	var ok bool
	if event, ok = _event.(*events.UserAccountDeleted); !ok {
		return WrongEventErr
	}

	var deleted models.User
	if err := db.Unscoped().First(&deleted, event.User).Error; err != nil {
		return err
	}

	var nickname models.Nickname
	if err := db.Unscoped().First(&nickname, *deleted.NicknameID).Error; err != nil {
		return err
	}

	deleted.Nickname = &nickname
	var groups []*models.Group
	if err := db.
		Set("gorm:auto_preload", false).
		Model(&models.Group{}).
		Where("id in (?)", event.Groups).
		Find(&groups).Error; err != nil {
		return err
	}

	orgs := set.New()
	for _, gr := range groups {
		orgs.Insert(gr.OrganizationID)
	}

	orgsIDs := make([]uint, 0, orgs.Len())
	orgs.Do(func(el interface{}) {
		orgsIDs = append(orgsIDs, el.(uint))
	})

	uniqueAdminsOfOrgs := make([]uint, 0)
	if err := db.Scopes(scopes.UniqueAdminsOfAllOrgs(orgsIDs, &uniqueAdminsOfOrgs)).Error; err != nil {
		return err
	}

	notif.Model = model
	notif.Deleted = &deleted
	notif.Groups = groups
	notif.UniqueAdminsOfOrgs = uniqueAdminsOfOrgs
	notif.OrgIDs = orgsIDs
	notif.loadReceivers()
	notif.loadDashReceivers()
	return nil
}
