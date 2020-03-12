package org

import (
	"dt/models"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение полной информации об организации.
//zenrpc:id id орг-ии.
//zenrpc:11 organization not found. организация с данным id не найдена.
//zenrpc:return при удачном выполнении запроса возвращает полную информацию об орг-ии.
func (s *Service) Get(id uint) (*views.Org, *zenrpc.Error) {
	var organization models.Organization
	if err := s.db.First(&organization, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	return views.OrgViewFromModelShort(&organization), nil
}
