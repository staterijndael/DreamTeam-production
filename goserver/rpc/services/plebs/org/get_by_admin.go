package org

import (
	"context"
	"dt/logwrap"
	"dt/requestContext"
	"github.com/semrush/zenrpc"
)

//получение списка организаций, в которых пользователь имеет права администратора.
//zenrpc:return при удачном выполнении запроса возвращает AdministratedOrganizations.
func (s *Service) GetByAdmin(
	ctx context.Context,
) (*AdministratedOrganizations, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	orgs, err := FindOrgByPerson(s.db, me.ID)
	if err != nil {
		logwrap.Debug("err finding user orgs: %v", err)
		return nil, err
	}

	return orgs, nil
}
