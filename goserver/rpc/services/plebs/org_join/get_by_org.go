package org_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//Получение всех запросов на присоединение к организации. Только для админов организации
//zenrpc:oid id организации
//zenrpc:37 отправитель запроса не является админом этой организации
//zenrpc:return при удачном выполнении возвращается список запросов на присоединение к группе.
func (s *Service) GetByOrg(ctx context.Context, oid uint) ([]*views.OrgJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var requests []*models.OrgJoinRequest
	if err := s.db.
		Model(&models.OrgJoinRequest{}).
		Where("status = 'pending'").
		Where("organization = ?", oid).
		Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if requests != nil && len(requests) > 0 &&
		(!requests[0].Organization.Admins.Contains(me.ID) && requests[0].InitiatorID != me.ID) {
		return nil, errors.New(errors.NotAssociated, nil, nil) // 37
	}

	result := make([]*views.OrgJoinRequest, 0, len(requests))
	for _, reqModel := range requests {
		result = append(result, views.OrgJoinRequestFromModel(reqModel))
	}

	return result, nil
}
