package user

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/stores"
	"github.com/semrush/zenrpc"
)

//установливает nickname
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:6 nickname is busy. данный nickname уже занят другим пользователем
//zenrpc:8 invalid nickname. неверный формат nickname. требуется непустая строка
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetNickName(ctx context.Context, nickname string) (*common.CodeAndMessage, *zenrpc.Error) {
	nicknameModel := models.Nickname{Value: nickname}
	if !nicknameModel.IsValid() {
		return nil, errors.New(errors.InvalidNickname, nil, nil)
	}

	nicknameModel.Validate()
	me := requestContext.CurrentUser(ctx)
	if err := s.db.Create(&nicknameModel).Error; err != nil {
		if stores.IsDuplicate(err) {
			return nil, errors.New(errors.NicknameIsBusy, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me.NicknameID = &nicknameModel.ID
	me.Nickname = &nicknameModel
	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("nickname", nicknameModel.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
