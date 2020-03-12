package user

import (
	"context"
	"database/sql"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//установливает email
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetEmail(ctx context.Context, email string) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	me.Email = sql.NullString{
		String: email,
		Valid:  true,
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("email", me.Email).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
