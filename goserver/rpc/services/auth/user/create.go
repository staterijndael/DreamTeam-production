package user

import (
	"dt/models"
	errors "dt/rpc/services/errors"
	"dt/stores"
	"dt/utils"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/semrush/zenrpc"
	"golang.org/x/crypto/bcrypt"
)

//zenrpc:7 невалидный телефон
//zenrpc:83 пользователь с таким телефоном уже зарегестрирован
//zenrpc:90 невалидный пароль
//zenrpc:return id созданного пользователя
func (s *Service) Create(phone, password string) (*uint, *zenrpc.Error) {
	if !utils.IsPhone([]byte(phone)) {
		return nil, errors.New(errors.InvalidPhone, nil, nil)
	}

	if !utils.IsValidPassword(password) {
		return nil, errors.New(errors.InvalidPassword, nil, nil)
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, errors.New(errors.Internal, hashErr, nil)
	}

	password = string(hash)
	formattedPhone := utils.FormatPhone(phone)
	user := models.User{Phone: formattedPhone}
	err := s.db.Unscoped().Where(&user).First(&user).Error
	if err == nil {
		return nil, errors.New(errors.PhoneAlreadyRegistered, err, nil)
	}

	if !gorm.IsRecordNotFoundError(err) {
		return nil, errors.New(errors.Internal, err, nil)
	}

	user = models.User{
		Phone:    formattedPhone,
		Password: password,
		Score: postgres.Jsonb{
			RawMessage: models.StartedScore,
		},
		AvatarID: stores.DefaultAvatars.User.ID,
	}

	if s.db.Create(&user).Error != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return &user.ID, nil
}
