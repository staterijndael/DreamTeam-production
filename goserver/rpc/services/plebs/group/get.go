package group

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//получение полной информации о группе.
//zenrpc:gid id группы.
//zenrpc:3 группа не найдена.
//zenrpc:34 пользователь не имеет права просматривать данный ресурс (не является членом ни одной из групп орг-и, не админ и не директор орг-ии).
//zenrpc:return при удачном выполнении запроса возвращается информация о группе.
func (s *Service) Get(
	ctx context.Context,
	gid uint,
) (*views.Group, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	isDirector := group.Organization.DirectorID == me.ID
	isOrgAdmin := group.Organization.Admins.Contains(me.ID)
	var isMemberOfAnyGroupOfOrg bool
	if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMemberOfAnyGroupOfOrg, group.OrganizationID, me.ID)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	} else if !isDirector && !isOrgAdmin && !isMemberOfAnyGroupOfOrg {
		return nil, errors.New(errors.CantViewThisGroup, nil, nil) // 34
	}

	return views.GroupFromModelShort(&group), nil
}
