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

//zenrpc:return все события рейтинга, к которым причастен текущий поль-ль (по всем орг-ям).
func (s *Service) IsRatingNowInMyGroups(ctx context.Context) ([]*views.Rating, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var ratings []*models.RatingEvent
	if err := s.db.Scopes(scopes.RatingEventsOfGroupsOfUser(me.ID)).Find(&ratings).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	ratingViews := make([]*views.Rating, len(ratings))
	for i := range ratings {
		ratingViews[i] = views.RatingEventFromModel(ratings[i])
	}

	return ratingViews, nil
}
