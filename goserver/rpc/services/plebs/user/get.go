package user

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение полной информации о пользователе
//zenrpc:uid id запрашиваемого пользователя
//zenrpc: user not found. запрашиваемый пользователь не найден
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о пользователе.
func (s *Service) Get(ctx context.Context, uid uint) (*views.User, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var u models.User
	err := s.db.Where(&models.User{Model: gorm.Model{ID: uid}}).First(&u).Error
	if err != nil || u.NicknameID == nil {
		if err == gorm.ErrRecordNotFound || u.NicknameID == nil {
			return nil, errors.New(errors.UserNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if uid == me.ID {
		return views.UserSelfViewFromModel(&u), nil
	}

	return views.UserViewFromModel(&u), nil
}
