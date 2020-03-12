package group_creation

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//отзыв запроса на создание группы. только создатель запроса имеет право на данную операцию.
//zenrpc:requestID id запроса. при уведомлении заменяется сущностью.
//zenrpc:31 пользователь не имеет права отзывать данный запрос.
//zenrpc:9 запрос на создание не найден
//zenrpc:30 запрос уже закрыт (отклонен | принят | отозван).
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) Withdraw(
	ctx context.Context,
	requestID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var req models.GroupCreationRequest
	var err error
	if err = s.db.First(&req, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.CreationRequestNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if req.Status != models.Pending || req.InitiatorID != me.ID {
		return nil, errors.New(errors.CantWithdrawCreationRequest, nil, nil)
	}

	if err = s.db.
		Model(&models.GroupCreationRequest{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{"status": models.Canceled, "acceptor": me.ID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.Unscoped().Delete(&req.Nickname).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.GroupCreationRequestWithdrawn{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: requestID,
	})

	return common.ResultOK, nil
}
