package group_creation

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение всех открытых запросов, на которые пользователь может повлиять
//zenrpc:37 вы не можете просматривать даный запрос.
//zenrpc:return при удачном выполнении запроса возвращается массив запросов на создание группы.
func (s *Service) GetByOrg(
	ctx context.Context,
	org uint,
) ([]*views.GroupCreationRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var requests []*models.GroupCreationRequest
	if err := s.db.
		Where("status = 'pending'").
		Where("organization = ?", org).
		Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if requests != nil && len(requests) > 0 &&
		(!requests[0].Organization.Admins.Contains(me.ID) && requests[0].InitiatorID != me.ID) {
		return nil, errors.New(errors.NotAssociated, nil, nil) // 37
	}

	result := make([]*views.GroupCreationRequest, 0)
	for _, reqModel := range requests {
		result = append(result, views.GroupCreationRequestFromModelShort(reqModel))
	}

	return result, nil
}
