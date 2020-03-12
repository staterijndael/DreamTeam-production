package org

import (
	"dt/managers/eventEmitter"
	"dt/models"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

func DenyAllOrgJoinRequests(db *gorm.DB, emitter *eventEmitter.EventEmitter, denier *models.User, oid uint) error {
	var requests []*models.OrgJoinRequest
	if err := db.
		Model(&models.OrgJoinRequest{}).
		Where("status = 'pending'").
		Where("organization = ?", oid).
		Find(&requests).Error; err != nil {
		return err
	}

	for _, req := range requests {
		if err := db.
			Model(&models.OrgJoinRequest{}).Where("id = ?", req.ID).
			Updates(map[string]interface{}{"status": models.Denied, "acceptor": denier.ID}).Error; err != nil {
			return err
		}

		//TODO issue #4
		//emitter.Emit(&events.OrgJoinRequestDenied{
		//	EventBase: events.EventBase{
		//		Context: ctx,
		//	},
		//	Request:   req.ID,
		//})
	}

	return nil
}

func FindOrgByDirector(db *gorm.DB, uid uint) ([]*models.Organization, *zenrpc.Error) {
	answer := make([]*models.Organization, 0)

	if err := db.
		Where(&models.Organization{DirectorID: uid}).
		Find(&answer).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return answer, nil
}

func FindOrgByAssociated(db *gorm.DB, uid uint) (answer []*models.Organization, rpcErr *zenrpc.Error) {
	if err := db.
		Model(&models.Organization{}).
		Where("community in (?)",
			db.
				Model(&models.MembershipOfCommunity{}).
				Where(`"user" = ?`, uid).
				Select("community").
				SubQuery(),
		).
		Where("director <> ?", uid).
		Find(&answer).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return answer, nil
}

func FindOrgByPerson(db *gorm.DB, uid uint) (*AdministratedOrganizations, *zenrpc.Error) {
	director, err := FindOrgByDirector(db, uid)
	if err != nil {
		return nil, err
	}

	associated, err := FindOrgByAssociated(db, uid)
	if err != nil {
		return nil, err
	}

	directorView, associatedView := make([]*views.Org, len(director)), make([]*views.Org, len(associated))

	for i := range director {
		directorView[i] = views.OrgViewFromModelShort(director[i])
	}

	for i := range associated {
		associatedView[i] = views.OrgViewFromModelShort(associated[i])
	}

	return &AdministratedOrganizations{
		Director: directorView,
		Linked:   associatedView,
	}, nil
}

type AdministratedOrganizations struct {
	Director []*views.Org `json:"director"`
	Linked   []*views.Org `json:"linked"`
}
