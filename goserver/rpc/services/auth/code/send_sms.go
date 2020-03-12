package code

import (
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/utils"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	UserID uint   `json:"userID"`
	Token  string `json:"token"`
	Exp    int64  `json:"exp"`
}

//отправляет верификационный код через смс указанному пользователю.
//zenrpc:-1 произошло некое дерьмо
//zenrpc:2 данный пользователь не найден (неверный пароль или телефон)
//zenrpc:7 невалидный телефон
//zenrpc:90 невалидный пароль
//zenrpc: возвращает ошибку, если была, в противном случае сообщение "oёk"
func (s *Service) SendSms(phone, password string, isAuth string) (*common.CodeAndMessage, *zenrpc.Error) {
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

	if isAuth != "auth" {

		_, err := s.sms.Send(u.ID, formattedPhone)
		if err != nil {
			return nil, zenrpc.NewError(-1, err)
		}

		token, exp, err := generateToken(u.ID, s.conf.SigningAlgorithm, s.conf.JWTIdentityKey)
		if err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}
	}

	return common.ResultOK, nil
}
