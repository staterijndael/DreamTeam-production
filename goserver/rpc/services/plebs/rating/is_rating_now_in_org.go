package rating

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//возвращает событие оценивания, если проходит в данной орг-ии.
//zenrpc:49 поль-ль не имеет права просматривать эвенты рейтинга данной орг-ии.
//zenrpc:50 в данной орг-ии на данный момент не проходит ratingEvent.
//zenrpc:return рейтинг орг-ии.
func (s *Service) IsRatingNowInOrg(ctx context.Context, oid uint) (*views.Rating, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	rating := models.RatingEvent{
		OrganizationID: oid,
	}

	if err := s.db.Where(&rating).First(&rating).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.NoRatingsInThisOrgForNow, nil, nil) // 50
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if !rating.Organization.Admins.Contains(me.ID) {
		var isMember bool
		if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMember, oid, me.ID)).Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		} else if !isMember {
			return nil, errors.New(errors.CantViewRatingEventsOfThisGroup, nil, nil) // 49
		}
	}

	return views.RatingEventFromModel(&rating), nil
}
