package group

import (
	"dt/managers/eventEmitter"
	"dt/models"
	"github.com/jinzhu/gorm"
)

func DenyAllGroupJoinRequests(db *gorm.DB, emitter *eventEmitter.EventEmitter, denier *models.User, gid uint) error {
	var joinRequests []*models.GroupJoinRequest
	if err := db.
		Model(&models.GroupJoinRequest{}).
		Where("status = 'pending'").
		Where(`"group" = ?`, gid).
		Find(&joinRequests).Error; err != nil {
		return err
	}

	for _, req := range joinRequests {
		if err := db.
			Model(&models.GroupJoinRequest{}).Where("id = ?", req.ID).
			Updates(map[string]interface{}{"status": models.Denied, "acceptor": denier.ID}).Error; err != nil {
			return err
		}

		//TODO issue #4
		//s.emitter.Emit(&events.GroupJoinRequestDenied{
		//	EventBase: events.EventBase{
		//		Context: ctx,
		//	},
		//	Request:   req.ID,
		//})
	}

	return nil
}

func DenyAllGroupCreationRequests(db *gorm.DB, emitter *eventEmitter.EventEmitter, denier *models.User, gid uint) error {
	var creationRequests []*models.GroupCreationRequest
	if err := db.
		Model(&models.GroupCreationRequest{}).
		Where("status = 'pending'").
		Where("hm = ?", gid).
		Find(&creationRequests).Error; err != nil {
		return err
	}

	for _, req := range creationRequests {
		if err := db.Unscoped().Delete(&req.Nickname).Error; err != nil {
			return err
		}

		if err := db.
			Model(&models.GroupCreationRequest{}).Where("id = ?", req.ID).
			Updates(map[string]interface{}{"status": models.Denied, "acceptor": denier.ID}).Error; err != nil {
			return err
		}

		//TODO issue #4
		//s.emitter.Emit(&events.GroupCreationRequestDenied{
		//	EventBase: events.EventBase{
		//		Context: ctx,
		//	},
		//	Request:   req.ID,
		//})
	}

	return nil
}
