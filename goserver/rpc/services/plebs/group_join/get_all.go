package group_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение всех запросов на которые может повлиять данный пользователь.
//zenrpc:return список запросов.
func (s *Service) GetAll(ctx context.Context) ([]*views.GroupJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var reqs []*models.GroupJoinRequest
	if err := s.db.
		Scopes(scopes.GroupJoinReqGetAll(me.ID)).
		Find(&reqs).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	reqViews := make([]*views.GroupJoinRequest, len(reqs))
	for i := range reqs {
		reqViews[i] = views.GroupJoinRequestFromModelShort(reqs[i])
	}

	return reqViews, nil
}
