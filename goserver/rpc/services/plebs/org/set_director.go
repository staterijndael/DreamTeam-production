package org

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//zenrpc:11 group not found
//zenrpc:76 только админ группы/орги может назначать нового админа
//zenrpc:77 данный пользователь не состоит в группе
//zenrpc:return возвращает "ok" в случае успеха
func (s *Service) SetDirector(ctx context.Context, oid, uid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if org.DirectorID != me.ID {
		return nil, errors.New(errors.OnlyDirectorCanSetNewDirector, nil, nil)
	}

	var user models.User
	for _, m := range org.Admins.Members {
		if m.UserID == uid {
			user = m.User
			break
		}
	}

	if user.ID == 0 {
		return nil, errors.New(errors.UserNotAdminOfOrg, nil, nil)
	}

	oldAdmin := org.DirectorID
	if err := s.db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("director", uid).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgNewDirector{
		EventBase: events.EventBase{
			Context: ctx,
		},
		OldDirector: oldAdmin,
		NewDirector: uid,
		Org:         oid,
	})

	return common.ResultOK, nil
}
