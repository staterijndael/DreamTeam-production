package group_join

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//отклонение запроса на вступление в группу. только админ группы/организации имеет право на данную операцию.
//zenrpc:requestID id запроса.
//zenrpc:41 запрос с данным id не найден.
//zenrpc:4 данный пользователь не имеет права отклонять данный запрос.
//zenrpc:45 запрос уже закрыт.
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) Deny(
	ctx context.Context,

	requestID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	code, err := changeStatusOfGroupJoinRequest(
		ctx,
		s.db,
		models.Denied,
		errors.CantDenyGroupJoinRequest,
		requestID,
		isAdminOfGroupOrOrg)

	if err == nil {
		s.emitter.Emit(&events.GroupJoinRequestDenied{
			EventBase: events.EventBase{
				Context: ctx,
			},
			Request: requestID,
		})
	}

	return code, err
}
