package user

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"github.com/semrush/zenrpc"
)

//Удаление аккаунта юзера
//zenrpc:85 нельзя удалить юзера, пока он является директором какой-либо организации
//zenrpc:88 нельзя удалить юзера, пока он является админом какой-либо группы
//zenrpc:return При успешном выполнении операции возваращает "ok", иначе ошибку.
func (s *Service) Delete(ctx context.Context) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var count int
	if err := s.db.
		Model(&models.Organization{}).
		Where("director = ?", me.ID).
		Count(&count).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if count > 0 {
		return nil, errors.New(errors.UserDeleteWhileIsDirectorOfAnyOrg, nil, nil)
	}

	if err := s.db.Scopes(scopes.GroupsWhereUserIsAdmin(me.ID)).Count(&count).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if count > 0 {
		return nil, errors.New(errors.UserDeleteWhileIsAdminOfAnyGroup, nil, nil)
	}

	groups := make([]uint, 0)
	if err := s.db.
		Scopes(scopes.IDsOfGroupsOfUser(me.ID, &groups)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.Scopes(scopes.UserDelete(*me)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserAccountDeleted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		User:   me.ID,
		Groups: groups,
	})

	return common.ResultOK, nil
}
