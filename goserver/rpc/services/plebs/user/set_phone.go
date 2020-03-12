package user

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/stores"
	"dt/utils"
	"github.com/semrush/zenrpc"
)

//установливает номер телефона
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:7 invalid phone. неверный формат номера телефона
//zenrpc:83 номер телефона занят другим пользователем
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetPhone(ctx context.Context, phone string) (*common.CodeAndMessage, *zenrpc.Error) {
	if !utils.IsPhone([]byte(phone)) {
		return nil, errors.New(errors.InvalidPhone, nil, nil)
	}

	me := requestContext.CurrentUser(ctx)
	me.Phone = utils.FormatPhone(phone)
	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("phone", me.Phone).Error; err != nil {
		if stores.IsDuplicate(err) {
			return nil, errors.New(errors.PhoneAlreadyRegistered, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
