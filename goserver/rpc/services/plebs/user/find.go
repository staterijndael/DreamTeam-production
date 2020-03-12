package user

import (
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
	"strings"
)

//поиск пользователя
//zenrpc:text last_name|first_name|nickname
//zenrpc:page=0 пагинация. отсчет идет с 0. размер стр = 50
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о пользователе.
func (s *Service) Find(text string, page uint64) ([]*views.User, *zenrpc.Error) {
	queries := common.SplitBySpacesRegex.Split(strings.Trim(text, " \t\r\f\n"), -1)
	var users []*models.User
	if err := s.db.
		Scopes(scopes.FindUser(page, queries)).
		Find(&users).Error; err != nil {
		if err != models.EmptyNicknameModelErr {
			return nil, errors.New(errors.Internal, err, nil)
		}
	}

	userViews := make([]*views.User, len(users))
	for i := range users {
		userViews[i] = views.UserViewFromModel(users[i])
	}

	return userViews, nil
}
