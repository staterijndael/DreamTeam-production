package user

import (
	"context"
	"database/sql"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//установливает фамилию
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetLastName(ctx context.Context, lastName string) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	oldUser := *me
	me.LastName = sql.NullString{
		String: lastName,
		Valid:  true,
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("last_name", me.LastName).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserRenamed{
		EventBase: events.EventBase{
			Context: ctx,
		},
		OldUser: oldUser,
		User:    *me,
	})

	return common.ResultOK, nil
}
