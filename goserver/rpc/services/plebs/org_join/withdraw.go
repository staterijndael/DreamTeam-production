package org_join

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

//Отзыв запроса на вступление в организацию. Только отправитель запроса имеет право на данную операцию.
//zenrpc:66 запрос на вступление в организацию не найден
//zenrpc:67 вы не можете отозвать данный запрос
//zenrpc:68 запрос уже закрыт
//zenrpc:requestID id запроса на вступление в организацию
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok"
func (s *Service) Withdraw(ctx context.Context, requestID uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)

	var request models.OrgJoinRequest
	if err := s.db.First(&request, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgJoinRequestNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID != request.InitiatorID {
		return nil, errors.New(errors.CantWithdrawOrgJoinRequest, nil, nil)
	}

	if request.Status != models.Pending {
		return nil, errors.New(errors.OrgJoinRequestAlreadyClosed, nil, nil)
	}

	if err := s.db.
		Model(&models.OrgJoinRequest{}).Where("id = ?", request.ID).
		Updates(map[string]interface{}{"status": models.Canceled, "acceptor": me.ID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgJoinRequestWithdrawn{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: requestID,
	})

	return common.ResultOK, nil
}
