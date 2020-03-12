package org

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//Проверка включено ли проведение рейтинга в организации
//zenrpc:11 организация не найдена
//zenrpc:86 пользователь не состоит в организации
//zenrpc:return Возвращает структуру, содержащую поле iEnabled и config (только если isEnabled == true), который содержит время начала рейтинга (час в сутках), день недели (от 0 до 6) и оргу
func (s *Service) IsRatingEnabled(ctx context.Context, oid uint) (*views.IsRatingEnabled, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	var isMember bool
	if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMember, oid, me.ID)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	} else if !isMember && !org.Admins.Contains(me.ID) {
		return nil, errors.New(errors.UserIsNotInOrg, nil, nil)
	}

	var conf models.RatingOrgConfig
	if err := s.db.Where("org_id = ?", oid).First(&conf).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &views.IsRatingEnabled{
				IsEnabled: false,
			}, nil
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	return &views.IsRatingEnabled{
		IsEnabled: true,
		Config:    views.RatingConfigFromModel(&conf),
	}, nil
}
