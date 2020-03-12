package org_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение полной информации о запросе. только инициатор запроса на вступление в организцаю имеет право на данную операцию
//zenrpc:requestID id запроса.
//zenrpc:66 запрос на вступление в организацию не найден
//zenrpc:69 вы не можете просматривать данный запрос
//zenrpc:return при удачном выполнении запроса возвращается полную инф-ию о запросе.
func (s *Service) Get(ctx context.Context, requestID uint) (*views.OrgJoinRequest, *zenrpc.Error) {
	var req models.OrgJoinRequest
	if err := s.db.First(&req, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgJoinRequestNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if me.ID != req.InitiatorID && !req.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.CantViewOrgJoinRequest, nil, nil)
	}

	return views.OrgJoinRequestFromModel(&req), nil
}
