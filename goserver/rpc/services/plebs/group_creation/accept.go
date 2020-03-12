package group_creation

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//отклонение запроса на создание группы. только админ родительской группы имеет право на данную операцию.
//zenrpc:requestID id запроса. при уведомлении заменяется сущностью.
//zenrpc:8 пользователь не имеет права одобрять данный запрос.
//zenrpc:9 запрос на создание не найден
//zenrpc:30 запрос уже закрыт (отклонен | одобрен | отозван).
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) Accept(
	ctx context.Context,
	requestID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	code, err, grID, _ := changeStatusIfLinkedOrDirector(ctx, models.Confirmed, requestID,
		errors.CantAcceptCreationRequest, s, s.emitter)

	if err == nil {
		s.emitter.Emit(&events.GroupCreationRequestAccepted{
			EventBase: events.EventBase{
				Context: ctx,
			},
			Request: requestID,
			Group:   grID,
		})
	}

	return code, err
}
