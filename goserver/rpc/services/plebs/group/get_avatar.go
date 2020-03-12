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

//zenrpc:23 группа не найдена
//zenrpc:34 вы не имеете права доступа к этому ресурсу.
func (s *Service) GetAvatar(ctx context.Context, gid uint) (*views.File, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if !group.Community.Contains(me.ID) && !group.Organization.Admins.Contains(me.ID) {
		var isMemberOfAnyGroupOfOrg bool
		if err := s.db.
			Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMemberOfAnyGroupOfOrg, group.OrganizationID, me.ID)).Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}

		if !isMemberOfAnyGroupOfOrg {
			return nil, errors.New(errors.CantViewThisGroup, nil, nil)
		}
	}

	ava, err := views.FileViewFromModel(&group.Avatar)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return ava, nil
}
