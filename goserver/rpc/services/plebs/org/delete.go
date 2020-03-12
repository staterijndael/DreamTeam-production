package org

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/rpc/services/plebs/group"
	"dt/scopes"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//Удаление организации. Только директор организации имеет право на эту операцию
//zenrpc:oid id организации
//zenrpc:11 организация не найдена
//zenrpc:81 только директор организации имеет право на эту операцию
//zenrpc:return Возвращает ok при успешном выполнении операции
func (s *Service) Delete(ctx context.Context, oid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var org models.Organization

	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID != org.DirectorID {
		return nil, errors.New(errors.OnlyOrgDirectorMayDeleteOrg, nil, nil)
	}

	var groups []*models.Group
	if err := s.db.
		Set("gorm:auto_preload", false).
		Model(&models.Group{}).
		Where("organization = ?", org.ID).
		Find(&groups).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var members []uint
	if err := s.db.Scopes(scopes.GroupMembersIDsOfOrg(oid, &members)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	tx := s.db.Begin()
	for _, gr := range groups {
		if err := group.DenyAllGroupJoinRequests(tx, s.emitter, me, gr.ID); err != nil {
			tx.Rollback()
			return nil, errors.New(errors.Internal, err, nil)
		}

		if err := group.DenyAllGroupCreationRequests(tx, s.emitter, me, gr.ID); err != nil {
			tx.Rollback()
			return nil, errors.New(errors.Internal, err, nil)
		}
	}

	if err := DenyAllOrgJoinRequests(tx, s.emitter, me, oid); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := tx.Scopes(scopes.OrganizationDelete(org, groups)).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgDeleted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Org:          oid,
		GroupMembers: members,
	})

	return common.ResultOK, nil
}
