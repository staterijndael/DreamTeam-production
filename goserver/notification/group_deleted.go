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

type GroupDeleted struct {
	notificationBase
	Group       *models.Group
	ChildGroups []*models.Group
	DeletedBy   *models.User
} //notifier

func (notif *GroupDeleted) loadReceivers() {
	if notif.Group.Parent != nil && notif.DeletedBy.ID != notif.Group.Parent.AdminID {
		notif.receivers = append(notif.receivers, notif.Group.Parent.AdminID)
	}

	for _, member := range notif.Group.Community.Members {
		if member.UserID == notif.DeletedBy.ID {
			continue
		}

		notif.receivers = append(notif.receivers, member.UserID)
	}

	if notif.ChildGroups != nil {
		for _, gr := range notif.ChildGroups {
			if gr.AdminID == notif.DeletedBy.ID {
				continue
			}

			notif.receivers = append(notif.receivers, gr.AdminID)
		}
	}
}

func (notif *GroupDeleted) loadDashReceivers() {
	for _, member := range notif.Group.Organization.Admins.Members {
		if member.UserID == notif.DeletedBy.ID {
			continue
		}

		notif.dashReceivers = append(notif.dashReceivers, member.UserID)
	}
}

func (notif *GroupDeleted) ContainerizedView() *utils.Container {
	return &utils.Container{
		Type: "notification.groupdeleted",
		Data: notif.View(),
	}
}

func (notif *GroupDeleted) View() interface{} {
	return &struct {
		ID        uint         `json:"id"`
		Group     *views.Group `json:"group"`
		DeletedBy *views.User  `json:"deletedBy"`
		Seen      *bool        `json:"seen,omitempty"`
	}{
		ID:        notif.GetModel().ID,
		DeletedBy: views.UserViewFromModel(notif.DeletedBy),
		Group:     views.GroupFromModelShort(notif.Group),
		Seen:      notif.seen,
	}
}

func (notif *GroupDeleted) CreateByEvent(db *gorm.DB, event interface{}) error {
	e, ok := event.(*events.GroupDeleted)
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

	wall, err := saveWallEvent(db, n, notif.Group.OrganizationID)
	if err != nil {
		return err
	}

	_, err = saveAOWSExcept(db, wall, notif.DeletedBy.ID)
	if err != nil {
		return err
	}

	s := set.New()
	if notif.Group.Parent != nil && notif.DeletedBy.ID != notif.Group.Parent.AdminID {
		s.Insert(notif.Group.Parent.AdminID)
	}

	for _, member := range notif.Group.Community.Members {
		if member.UserID == notif.DeletedBy.ID {
			continue
		}

		s.Insert(member.UserID)
	}

	if notif.ChildGroups != nil {
		for _, gr := range notif.ChildGroups {
			if gr.AdminID == notif.DeletedBy.ID {
				continue
			}

			s.Insert(gr.AdminID)
		}
	}

	var anError bool
	s.Do(func(el interface{}) {
		if anError {
			return
		}

		id, _ := el.(uint)
		if _, err = saveSingleUNS(db, n, id); err != nil {
			anError = true
			return
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func (notif *GroupDeleted) Load(db *gorm.DB, n *models.Notification) error {
	var e *events.GroupDeleted
	if err := json.Unmarshal(n.Data.RawMessage, &e); err != nil {
		return err
	}

	return notif.LoadWithEvent(db, e, n)
}

func (notif *GroupDeleted) LoadWithEvent(db *gorm.DB, _event interface{}, model *models.Notification) error {
	var event *events.GroupDeleted
	var ok bool
	if event, ok = _event.(*events.GroupDeleted); !ok {
		return WrongEventErr
	}

	var deletedBy models.User
	if err := db.First(&deletedBy, event.DeletedBy).Error; err != nil {
		return err
	}

	var group models.Group
	if err := db.
		Unscoped().
		First(&group, event.Group).Error; err != nil {
		return err
	}

	var parent models.Group
	if group.ParentID != nil {
		if err := db.First(&parent, *group.ParentID).Error; err != nil {
			return err
		}

		group.Parent = &parent
	}

	var childGroups []*models.Group
	if err := db.Where(`id in (?)`, []int64(group.ChildrenIDs)).Find(&childGroups).Error; err != nil {
		return err
	}

	notif.Model = model
	notif.Group = &group
	notif.ChildGroups = childGroups
	notif.DeletedBy = &deletedBy
	notif.loadReceivers()
	notif.loadDashReceivers()

	return nil
}
