package group_creation

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение информации о запросе.
//zenrpc:requestID id запроса. при уведомлении заменяется сущностью.
//zenrpc:9 запрос на создание не найден
//zenrpc:71 вы не можете просматривать данный запрос на создание группы
//zenrpc:return при удачном выполнении запроса возвращает GroupCreationRequest.
func (s *Service) Get(
	ctx context.Context,
	requestID uint,
) (*views.GroupCreationRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var req models.GroupCreationRequest
	var err error
	if err = s.db.First(&req, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.CreationRequestNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if !req.Organization.Admins.Contains(me.ID) ||
		me.ID != req.InitiatorID ||
		(req.AcceptorID != nil && me.ID != *req.AcceptorID) ||
		(req.Parent != nil && req.Parent.AdminID != me.ID) {
		return nil, errors.New(errors.CantViewThisGroupCreationRequest, nil, nil)
	}

	return views.GroupCreationRequestFromModelShort(&req), nil
}
