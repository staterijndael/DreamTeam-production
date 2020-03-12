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

//добавление пользователя к организации, дает полный доступ, не считая удаления организации и управления списком
// слинкованных пользователей.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя,
// а так же слинкованным к орг-ии пользователям.
//zenrpc:oid id орг-ии. при уведомлении сменяется на полную информацию об орг-ии.
//zenrpc:uid id добавляемого человека. при уведомлении сменяется на полную информацию о данном пользователе.
//zenrpc: user not found. пользователь с id равным uid не найден
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:1 operation on organization is not permitted. только директор орг-ии имеет права на данную операцию
//zenrpc:13 пользователь уже добавлен к организации
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) AddUser(ctx context.Context, oid, uid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var user models.User
	if err := s.db.
		First(&user, uid).Error; err != nil || user.NicknameID == nil {
		if gorm.IsRecordNotFoundError(err) || user.NicknameID == nil {
			return nil, errors.New(errors.UserNotFound, err, nil) //
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil) // 11
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if org.DirectorID != me.ID {
		return nil, errors.New(errors.OrgOperationNotPermitted, nil, nil) // 1
	}

	if org.Admins.Contains(uid) {
		return nil, errors.New(errors.UserAlreadyInOrg, nil, nil) // 13
	}

	if err := s.db.Create(&models.MembershipOfCommunity{UserID: uid, CommunityID: org.CommunityID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserAssociated{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Org:          oid,
		Associated:   uid,
		AssociatedBy: me.ID,
	})

	return common.ResultOK, nil
}
