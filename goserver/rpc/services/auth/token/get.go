package token

import (
	"dt/config"
	"dt/models"
	"dt/rpc/services/errors"
	"dt/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Token struct {
	UserID uint   `json:"userID"`
	Token  string `json:"token"`
	Exp    int64  `json:"exp"`
}

//получить авторизационный токен по телефону и коду
//zenrpc:code код из смс или уведомления
//zenrpc: данный пользователь не найден
//zenrpc:7 невалидный телефон
//zenrpc:82 неверный код аутентификации
//zenrpc:return токен или ошибка
func (s *Service) Get(phone string, code int, isAuth string, password string) (*Token, *zenrpc.Error) {
	if !utils.IsPhone([]byte(phone)) {
		return nil, errors.New(errors.InvalidPhone, nil, nil)
	}

	u := models.User{Phone: utils.FormatPhone(phone)}
	if err := s.db.Where(&u).First(&u).Error; err != nil && err != models.EmptyNicknameModelErr {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.UserNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if isAuth == "auth" {

		hashErr := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

		if hashErr != nil {
			return nil, errors.New(errors.UserNotFound, hashErr, nil)
		}
	} else {

		if codeFromManager, ok := s.sms.Get(u.ID); !ok || codeFromManager != code {
			return nil, errors.New(errors.InvalidAuthCode, nil, nil)
		}
	}

	token, exp, err := generateToken(u.ID, s.conf.SigningAlgorithm, s.conf.JWTIdentityKey)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.sms.Delete(u.ID)
	return &Token{UserID: u.ID, Token: token, Exp: exp.Unix()}, nil
}

func generateToken(uid uint, alg, idKey string) (tokenString string, exp *time.Time, err error) {
	token := jwt.New(jwt.GetSigningMethod(alg))
	claims := token.Claims.(jwt.MapClaims)
	now := time.Now()
	expire := now.Add(config.JWTDuration)
	claims[idKey] = uid
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = now.Unix()
	tokenString, err = token.SignedString([]byte(config.VerySecretKey))
	return tokenString, &expire, err
}
