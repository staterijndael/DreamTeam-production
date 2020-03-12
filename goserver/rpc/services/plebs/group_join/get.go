package group_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение полной информации о запросе. только отправитель запроса, админ целевой группы или орг-ии
// имеет право на данную операцию.
//zenrpc:requestID id запроса.
//zenrpc:41 запрос с данным id не найден.
//zenrpc:46 доступ к ресурсу для данного польз-я закрыт.
//zenrpc:return при удачном выполнении запроса возвращается полную инф-ию о запросе.
func (s *Service) Get(ctx context.Context, requestID uint) (*views.GroupJoinRequest, *zenrpc.Error) {
	var req models.GroupJoinRequest
	if err := s.db.First(&req, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupJoinRequestNotFound, err, nil) // 41
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if !isInitiator(me.ID, &req) && !isAdminOfGroupOrOrg(me.ID, &req) {
		return nil, errors.New(errors.CantViewGroupJoinRequest, nil, nil) // 46
	}

	return views.GroupJoinRequestFromModelShort(&req), nil
}
