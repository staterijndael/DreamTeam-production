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

//удаление пользователя от управления организацией.
//только директор компании имеет право на данную операцию.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя,
// а так же слинкованным к орг-ии пользователям.
//zenrpc:oid id орг-ии. при уведомлении сменяется на полную информацию об орг-ии.
//zenrpc:uid id удаляемого человека. при уведомлении сменяется на полную информацию о данном пользователе.
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:1 operation on organization is not permitted. только директор орг-ии имеет права на данную операцию
//zenrpc:16 попытка удалить директора данной орг-ии.
//zenrpc:17 пользователь с id = uid не слинкован к данной орг-ии.
//zenrpc:35 нельзя диссоциировать себя.
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) RemoveUser(ctx context.Context, oid, uid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	if me.ID == uid {
		return nil, errors.New(errors.CantDissociateYourself, nil, nil) // 35
	}

	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil) // 11
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if uid == org.DirectorID {
		return nil, errors.New(errors.RemoveDirector, nil, nil) // 16
	}

	if org.DirectorID != me.ID {
		return nil, errors.New(errors.OrgOperationNotPermitted, nil, nil) // 1
	}

	if !org.Admins.Contains(uid) {
		return nil, errors.New(errors.NotAssociatedUser, nil, nil) // 17
	}

	if err := s.db.
		Where(`community = ?`, org.CommunityID).
		Where(`"user" = ?`, uid).
		Delete(models.MembershipOfCommunity{}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserDissociatedByDirector{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Org:      oid,
		User:     uid,
		Director: me.ID,
	})

	return common.ResultOK, nil
}
