package notification

import (
	"dt/models"
	"errors"
	"github.com/jinzhu/gorm"
)

var ErrPanic = errors.New("panic in notification.utils")

func saveNotification(db *gorm.DB, event interface{}) (*models.Notification, error) {
	n, err := models.NotificationWithEvent(event)
	if err != nil {
		return nil, err
	}

	creation := db.Create(n)
	if creation.Error != nil {
		return nil, creation.Error
	}

	if err := db.First(n, n.ID).Error; err != nil {
		return nil, err
	}

	return n, nil
}

func saveWallEvent(db *gorm.DB, n *models.Notification, orgID uint) (*models.OrgWall, error) {
	w := models.OrgWall{
		NotificationID: n.ID,
		OrganizationID: orgID,
	}

	creation := db.Create(&w)
	if creation.Error != nil {
		return nil, creation.Error
	}

	if err := db.First(&w, w.ID).Error; err != nil {
		return nil, err
	}

	return &w, nil
}

func saveAOWS(db *gorm.DB, wall *models.OrgWall) (res []*models.AdminOrgWallSeen, err error) {
	aowsIDS := make([]uint, len(wall.Organization.Admins.Members))
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = ErrPanic
			tx.Rollback()
		}
	}()

	for i := range wall.Organization.Admins.Members {
		aows := models.AdminOrgWallSeen{
			UserID: wall.Organization.Admins.Members[i].UserID,
			WallID: wall.ID,
		}

		if err = tx.Create(&aows).Error; err != nil {
			tx.Rollback()
			return
		}

		aowsIDS[i] = aows.ID
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	if err = db.
		Where(`id in (?)`, aowsIDS).
		Find(&res).Error; err != nil {
		return
	}

	return
}

func saveAOWSExcept(db *gorm.DB, wall *models.OrgWall, excludedID uint) (res []*models.AdminOrgWallSeen, err error) {
	aowsIDS := make([]uint, len(wall.Organization.Admins.Members))
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = ErrPanic
			tx.Rollback()
		}
	}()

	for i := range wall.Organization.Admins.Members {
		if excludedID == wall.Organization.Admins.Members[i].UserID {
			continue
		}

		aows := models.AdminOrgWallSeen{
			UserID: wall.Organization.Admins.Members[i].UserID,
			WallID: wall.ID,
		}

		if err = tx.Create(&aows).Error; err != nil {
			tx.Rollback()
			return
		}

		aowsIDS[i] = aows.ID
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	if err = db.
		Where(`id in (?)`, aowsIDS).
		Find(&res).Error; err != nil {
		return
	}

	return
}

func saveSingleUNS(db *gorm.DB, n *models.Notification, uid uint) (*models.UserNotificationSeen, error) {
	uns := models.UserNotificationSeen{
		UserID:         uid,
		NotificationID: n.ID,
	}

	creation := db.Create(&uns)
	if creation.Error != nil {
		return nil, creation.Error
	}

	if err := db.First(&uns, uns.ID).Error; err != nil {
		return nil, err
	}

	return &uns, nil
}

func saveUNS(db *gorm.DB, n *models.Notification, users []uint) (res []*models.UserNotificationSeen, err error) {
	unsIDS := make([]uint, len(users))
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			err = ErrPanic
			tx.Rollback()
		}
	}()

	for i := range users {
		uns := models.UserNotificationSeen{
			UserID:         users[i],
			NotificationID: n.ID,
		}

		if err = tx.Create(&uns).Error; err != nil {
			tx.Rollback()
			return
		}

		unsIDS[i] = uns.ID
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	if err = db.
		Where(`id in (?)`, unsIDS).
		Find(&res).Error; err != nil {
		return
	}

	return
}

func saveTimeline(db *gorm.DB, notification, group uint, rating *uint) (*models.Timeline, error) {
	tl := models.Timeline{
		GroupID:        group,
		NotificationID: notification,
		EventID:        rating,
	}

	creation := db.Create(&tl)
	if creation.Error != nil {
		return nil, creation.Error
	}

	if err := db.First(&tl, tl.ID).Error; err != nil {
		return nil, err
	}

	return &tl, nil
}
