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

//выключает проведение рейтинга в данной орг-ии.
//zenrpc:11 org not found
//zenrpc:78 только админ орги может включать/выключать ее рейтинг
//zenrpc:return возвращает "ok" в случае успеха
func (s *Service) DisableRating(ctx context.Context, oid uint) (*common.CodeAndMessage, *zenrpc.Error) {
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

	if err := s.db.Unscoped().Where("org_id = ?", oid).Delete(&models.RatingOrgConfig{}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.RatingDisabled{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Organization: oid,
	})

	return common.ResultOK, nil
}
