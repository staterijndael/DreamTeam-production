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

//отклонение запроса на вступление в организацию. только админ организации имеет право на данную операцию.
//zenrpc:requestID id запроса.
//zenrpc:66 запрос с данным id не найден.
//zenrpc:68 запрос уже закрыт.
//zenrpc:69 данный пользователь не имеет права отклонять данный запрос.
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) Deny(
	ctx context.Context,
	requestID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	var request models.OrgJoinRequest
	if err := s.db.First(&request, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgJoinRequestNotFound, err, nil) // 66
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if !request.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.CantViewOrgJoinRequest, nil, nil) // 69
	}

	if request.Status != models.Pending {
		return nil, errors.New(errors.OrgJoinRequestAlreadyClosed, nil, nil) // 68
	}

	if err := s.db.
		Model(&models.OrgJoinRequest{}).Where("id = ?", request.ID).
		Updates(map[string]interface{}{"status": models.Denied, "acceptor": me.ID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgJoinRequestDenied{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: requestID,
	})

	return common.ResultOK, nil
}
