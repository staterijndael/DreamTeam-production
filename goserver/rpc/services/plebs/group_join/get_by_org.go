package group_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение всех запросов на присоединение к группе по заданной организации
//zenrpc:oid id организации
//zenrpc:37 отправитель запроса не является админом этой организации
//zenrpc:return при удачном выполнении возвращается список запросов на присоединение к группе.
func (s *Service) GetByOrg(ctx context.Context, oid uint) ([]*views.GroupJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var requests []*models.GroupJoinRequest
	if err := s.db.
		Model(&models.GroupJoinRequest{}).
		Where("status = 'pending'").
		Where("group in (?)",
			s.db.
				Model(&models.Group{}).
				Where("organization = ?", oid).
				Select("id").
				SubQuery(),
		).
		Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if requests != nil && len(requests) > 0 && !requests[0].Group.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.NotAssociated, nil, nil) // 37
	}

	result := make([]*views.GroupJoinRequest, 0)
	for _, reqModel := range requests {
		result = append(result, views.GroupJoinRequestFromModelShort(reqModel))
	}

	return result, nil
}
