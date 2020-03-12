package org

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//установка аватара организации.
//кодировка - base64. ответственность за кодировку лежит на пользователе.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя,
// а так же слинкованным к орг-ии пользователям.
//zenrpc:oid id орг-ии. при уведомлении сменяется на полную информацию об орг-ии.
//zenrpc:avatar filtInput. контент файла в кодировке base64. при отправке уведомления заменяется с fileInput на fileMetaInfo, содержащий мета информацию о файле, но не контент файла
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:1 operation on organization is not permitted. только директор и слинкованные пользователи орг-ии имеют права на данную операцию
//zenrpc:return при удачном выполнении запроса возвращает мета информацию о файле.
func (s *Service) SetAvatar(
	ctx context.Context,
	oid uint,
	avatar *common.FileInput,
) (*views.FileMetaInfo, *zenrpc.Error) {
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

	dbFile, errDBFile := common.CreateDBFile(s.db, s.conf, []byte(avatar.Content))
	if errDBFile != nil {
		return nil, errors.New(errors.Internal, errDBFile, nil)
	}

	if err := s.db.Model(&models.Organization{}).Where("id = ?", org.ID).Update("avatar", dbFile.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return views.FileMetaInfoViewFromModel(dbFile), nil
}
