package group_creation

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение всех открытых запросов, на которые пользователь может повлиять
//zenrpc:return при удачном выполнении запроса возвращается массив запросов на создание группы.
func (s *Service) GetAll(
	ctx context.Context,
) ([]*views.GroupCreationRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var requests []*models.GroupCreationRequest
	if err := s.db.Scopes(scopes.GroupCreationReqGetAll(me.ID)).Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	result := make([]*views.GroupCreationRequest, 0)
	for _, reqModel := range requests {
		result = append(result, views.GroupCreationRequestFromModelShort(reqModel))
	}

	return result, nil
}
