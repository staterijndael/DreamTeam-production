package org

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"time"
)

//включает проведение рейтинга в данной орг-ии.
//zenrpc:11 org not found
//zenrpc:78 только админ орги может включать/выключать ее рейтинг
//zenrpc:79 рейтинг в данной орге уже включен
//zenrpc:return возвращает созданные настройки рейтинга в случае успеха
func (s *Service) EnableRating(ctx context.Context, oid uint) (*views.RatingConfig, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if !org.Admins.Contains(me.ID) {
		return nil, errors.New(errors.OnlyAdminOfOrgCanControlRating, nil, nil)
	}

	conf := models.RatingOrgConfig{OrganizationID: oid}
	if err := s.db.Where(&conf).First(&conf).Error; err == nil {
		return nil, errors.New(errors.RatingAlreadyEnabled, nil, nil)
	} else if !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	conf = models.RatingOrgConfig{
		Organization:   &org,
		OrganizationID: oid,
		WeekDay:        time.Friday,
		StartTime:      17,
	}

	if err := s.db.Create(&conf).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.RatingEnabled{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Organization: oid,
		Config:       conf.ID,
	})

	return views.RatingConfigFromModel(&conf), nil
}
