package org_join

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//Открытие запроса на вступление в организацию
//zenrpc:11 организация не найдена
//zenrpc:13 пользователь уже состоит в этой организации
//zenrpc:84 данный запрос уже открыт
//zenrpc:return при удачном выполнении запроса возвращается тело запроса на вступление в организацию
func (s *Service) Start(ctx context.Context, oid uint) (*views.OrgJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	oldReq := models.OrgJoinRequest{
		RequestBase: models.RequestBase{
			Status:      models.Pending,
			InitiatorID: me.ID,
		},
		OrganizationID: oid,
	}

	if err := s.db.Where(&oldReq).First(&oldReq).Error; err == nil {
		return nil, errors.New(errors.RequestAlreadyOpened, nil, nil)
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil) // 11
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	var isMember bool
	if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMember, oid, me.ID)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	} else if isMember {
		return nil, errors.New(errors.UserAlreadyInOrg, err, nil) // 13
	}

	req := models.OrgJoinRequest{
		RequestBase: models.RequestBase{
			Status:      models.Pending,
			InitiatorID: me.ID,
			AcceptorID:  nil,
		},
		OrganizationID: oid,
		GroupID:        nil,
	}

	if err := s.db.Create(&req).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.First(&req, req.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgJoinRequestStarted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: req.ID,
	})

	return views.OrgJoinRequestFromModel(&req), nil
}
