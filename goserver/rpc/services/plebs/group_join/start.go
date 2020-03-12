package group_join

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

//открытие запроса на вступление в группу.
//zenrpc:gid id группы.
//zenrpc:3 группа не найдена.
//zenrpc:38 пользователь уже в группе.
//zenrpc:84 данный запрос уже открыт
//zenrpc:return при удачном выполнении запроса возвращается тело запроса на вступление в группу.
func (s *Service) Start(
	ctx context.Context,
	gid uint,
) (*views.GroupJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	oldReq := models.GroupJoinRequest{
		RequestBase: models.RequestBase{
			Status:      models.Pending,
			InitiatorID: me.ID,
		},
		GroupID: gid,
	}

	if err := s.db.Where(&oldReq).First(&oldReq).Error; err == nil {
		return nil, errors.New(errors.RequestAlreadyOpened, nil, nil)
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	// check not already member of group
	if group.Community.Contains(me.ID) {
		return nil, errors.New(errors.AlreadyInGroup, nil, nil) // 38
	}

	// if in organization (in any group of ...)
	var isMember bool
	if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMember, group.OrganizationID, me.ID)).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	} else if !isMember {
		return nil, errors.New(errors.NotInOrg, err, nil) // 39
	}

	req := models.GroupJoinRequest{
		RequestBase: models.RequestBase{
			Status:      models.Pending,
			InitiatorID: me.ID,
			AcceptorID:  nil,
		},
		GroupID: gid,
	}

	if creation := s.db.Create(&req); creation.Error != nil {
		return nil, errors.New(errors.Internal, creation.Error, nil)
	}

	if err := s.db.First(&req, req.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.GroupJoinRequestStarted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: req.ID,
	})

	return views.GroupJoinRequestFromModel(&req, &req.Initiator, req.Acceptor, &req.Group), nil
}
