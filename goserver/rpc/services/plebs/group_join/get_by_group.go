package group_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение всех запросов по группе.
//пользователь должен быть админом этой группы, иначе вернётся пустой список
//zenrpc:gid id группы
//zenrpc:return при удачном выполнении запроса возвращается список запросов на присоединение к группе
func (s *Service) GetByGroup(ctx context.Context, gid uint) ([]*views.GroupJoinRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var requests []*models.GroupJoinRequest
	if err := s.db.
		Model(&models.GroupJoinRequest{}).
		Where("status = 'pending'").
		Where("group = ?", gid).
		Find(&requests).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	result := make([]*views.GroupJoinRequest, 0)
	for _, req := range requests {
		if req.Group.Organization.Admins.Contains(me.ID) ||
			req.Group.AdminID == me.ID ||
			req.InitiatorID == me.ID {
			result = append(result, views.GroupJoinRequestFromModelShort(req))
		}
	}

	return result, nil
}
