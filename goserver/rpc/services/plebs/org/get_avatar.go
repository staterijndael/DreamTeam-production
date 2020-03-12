package org

import (
	"context"
	"dt/models"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение аватара организации.
//zenrpc:oid id орг-ии.
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:1 organization does not have avatar. организация не имеет аватара. возвращается в случае недоступности стандартного аватара на сервере.
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о файле, содержание файла.
func (s *Service) GetAvatar(ctx context.Context, oid uint) (*views.File, *zenrpc.Error) {
	var requestedOrg models.Organization
	if err := s.db.First(&requestedOrg, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	fileView, err := views.FileViewFromModel(&requestedOrg.Avatar)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return fileView, nil
}
