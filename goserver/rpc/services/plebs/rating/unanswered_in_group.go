package rating

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//zenrpc:return спсок всех неоценненных пользователей группы для заданного события рейтинга. Если рейтинг уже закончился или отправитель не состоит в этой группе, то будет возвращён пустой список
func (s *Service) UnansweredInGroup(ctx context.Context, ratingID, gid uint) ([]*views.User, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var unanswered []*models.User
	if err := s.db.
		Scopes(scopes.UnansweredInGroup(me.ID, ratingID, gid)).
		Find(&unanswered).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	userViews := make([]*views.User, len(unanswered))
	for i := range unanswered {
		userViews[i] = views.UserViewFromModel(unanswered[i])
	}

	return userViews, nil
}
