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

//установливает имя
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//senderID не заменяется сущностью sender'а.
//zenrpc:sender id отправителя запроса
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetFirstName(ctx context.Context, firstName string) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	oldUser := *me
	me.FirstName = sql.NullString{
		String: firstName,
		Valid:  true,
	}

	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("first_name", me.FirstName).Error; err != nil {
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
