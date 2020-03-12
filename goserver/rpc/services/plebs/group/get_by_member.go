package group

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение полной информации обо всех группах, где данный пользователь яв-ся участником или админом.
//zenrpc:return при удачном выполнении запроса возвращается информация о группах.
func (s *Service) GetByMember(
	ctx context.Context,
) ([]*views.Group, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var groups []*models.Group
	if err := s.db.Scopes(scopes.GroupsOfUser(me.ID)).Find(&groups).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	answer := make([]*views.Group, 0)
	for _, g := range groups {
		answer = append(answer, views.GroupFromModelShort(g))
	}

	return answer, nil
}
