package group

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//Удаление группы. Только админ группы или организации имеет право на эту операцию
//zenrpc:23 группа не найдена
//zenrpc:80 только админ группы или организации имеет право на эту операцию
//zenrpc:return Возвращает ok при успешном выполнении операции
func (s *Service) Delete(ctx context.Context, gid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID != group.AdminID && !group.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.OnlyOrgAdminOrGroupAdminMayDeleteGroup, nil, nil)
	}

	tx := s.db.Begin()
	if err := DenyAllGroupJoinRequests(tx, s.emitter, me, gid); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := DenyAllGroupCreationRequests(tx, s.emitter, me, gid); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := tx.Scopes(scopes.MoveToParent(group)).Scopes(scopes.GroupDelete(group)).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	exMembers := models.OrgExMember{
		OrganizationID: group.OrganizationID,
		CommunityID:    group.CommunityID,
		Organization:   group.Organization,
		Community:      group.Community,
	}

	if err := tx.Create(&exMembers).Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.GroupDeleted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group:     gid,
		DeletedBy: me.ID,
	})

	return common.ResultOK, nil
}
