package code

import (
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/utils"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"golang.org/x/crypto/bcrypt"
)

//отправляет верификационный код через уведомление указанному пользователю.
//zenrpc: данный пользователь не найден (неверный пароль или телефон)
//zenrpc:7 невалидный телефон
//zenrpc:90 невалидный пароль
//zenrpc: возвращает ошибку, если была, в противном случае сообщение "ok"
func (s *Service) SendNotification(phone, password string) (*common.CodeAndMessage, *zenrpc.Error) {
	if !utils.IsPhone([]byte(phone)) {
		return nil, errors.New(errors.InvalidPhone, nil, nil)
	}

	if !utils.IsValidPassword(password) {
		return nil, errors.New(errors.InvalidPassword, nil, nil)
	}

	formattedPhone := utils.FormatPhone(phone)
	u := models.User{Phone: formattedPhone}
	if err := s.db.Where(&u).First(&u).Error; err != nil && err != models.EmptyNicknameModelErr {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.UserNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if hashErr != nil {
		return nil, errors.New(errors.UserNotFound, hashErr, nil)
	}

	msg := &views.JSONRPCNotification{
		Method: "new",
		Params: &utils.Container{
			Type: "newverificationcode",
			Data: &struct {
				Code int `json:"code"`
			}{
				Code: s.sms.GenerateCode(u.ID),
			},
		},
	}

	s.cm.Send(u.ID, msg, nil)
	s.dcm.Send(u.ID, msg, nil)
	return common.ResultOK, nil
}
