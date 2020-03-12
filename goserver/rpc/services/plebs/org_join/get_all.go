package org_join

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
func (s *Service) GetAll(ctx context.Context) ([]*views.OrgJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var reqs []*models.OrgJoinRequest
	if err := s.db.
		Scopes(scopes.OrgJoinReqGetAll(me.ID)).
		Find(&reqs).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	reqViews := make([]*views.OrgJoinRequest, len(reqs))
	for i := range reqs {
		reqViews[i] = views.OrgJoinRequestFromModel(reqs[i])
	}

	return reqViews, nil
}
