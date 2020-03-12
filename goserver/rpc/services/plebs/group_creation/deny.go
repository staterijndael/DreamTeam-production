package group_creation

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//отклонение запроса на создание группы. только админ родительской группы и админы организации имееют право на данную операцию.
//zenrpc:requestID id запроса. при уведомлении заменяется сущностью.
//zenrpc:7 пользователь не имеет права отклонять данный запрос.
//zenrpc:9 запрос на создание не найден
//zenrpc:30 запрос уже закрыт (отклонен | одобрен | отозван).
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) Deny(
	ctx context.Context,
	requestID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	code, err, _, req := changeStatusIfLinkedOrDirector(ctx, models.Denied, requestID,
		errors.CantDenyCreationRequest, s, s.emitter)

	if err == nil {
		if err := s.db.Unscoped().Delete(&req.Nickname).Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}

		s.emitter.Emit(&events.GroupCreationRequestDenied{
			EventBase: events.EventBase{
				Context: ctx,
			},
			Request: requestID,
		})
	}

	return code, err
}
