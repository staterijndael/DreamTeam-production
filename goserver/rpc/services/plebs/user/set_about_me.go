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

//установливает текст "о себе"
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetAboutMe(ctx context.Context, aboutMe string) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	me.AboutMe = sql.NullString{
		String: aboutMe,
		Valid:  true,
	}

	if err := s.db.
		Model(&models.User{}).Where("id = ?", me.ID).
		Update("about_me", me.AboutMe).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
