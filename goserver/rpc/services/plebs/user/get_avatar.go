package user

import (
	"context"
	"dt/models"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение аватара пользователя.
//кодировка - base64. ответственность за кодировку лежит на пользователе.
//zenrpc:uid id запрашиваемого пользователя
//zenrpc: user not found. запрашиваемый пользователь не найден
//zenrpc:3 specified user does not have avatar. возвращается в случае недоступности стандартного аватара на сервере.
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о файле, содержание файла.
func (s *Service) GetAvatar(ctx context.Context, uid uint) (*views.File, *zenrpc.Error) {
	var requestedUser models.User

	err := s.db.First(&requestedUser, uid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.UserNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	fileView, err := views.FileViewFromModel(&requestedUser.Avatar)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return fileView, nil
}
