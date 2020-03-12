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
//zenrpc:group id группы
//zenrpc:return при удачном выполнении запроса возвращается массив запросов на создание группы.
func (s *Service) GetByGroup(
	ctx context.Context,
	group uint,
) ([]*views.GroupCreationRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var requests []*models.GroupCreationRequest
	if err := s.db.
		Where("status = 'pending'").
		Where("hm = ?", group).
		Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	result := make([]*views.GroupCreationRequest, 0)
	for _, req := range requests {
		if req.Organization.Admins.Contains(me.ID) ||
			(req.Parent != nil && req.Parent.AdminID == me.ID) ||
			req.InitiatorID == me.ID {
			result = append(result, views.GroupCreationRequestFromModelShort(req))
		}
	}

	return result, nil
}
