package org

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//установка описания организации.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя,
// а так же слинкованным к орг-ии пользователям.
//zenrpc:oid id орг-ии. при уведомлении сменяется на полную информацию об орг-ии.
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:1 operation on organization is not permitted. только директор и слинкованные пользователи орг-ии имеют права на данную операцию
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetDescription(
	ctx context.Context,
	oid uint,
	description string,
) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil) // 11
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID != org.DirectorID && !org.Admins.Contains(me.ID) {
		return nil, errors.New(errors.OrgOperationNotPermitted, nil, nil)
	}

	if errDB := s.db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("description", description).Error; errDB != nil {
		return nil, errors.New(errors.Internal, errDB, nil)
	}
	return common.ResultOK, nil
}
