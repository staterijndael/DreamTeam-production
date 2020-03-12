package group

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//получение полной информации о группах орг-ии.
//zenrpc:oid
//zenrpc:34 поль-ль не имеет права просматривать данный ресурс (не яв-ся ни одной из групп орг-и, не админ и не директор)
//zenrpc:return при удачном выполнении возваращется информация о группе
func (s *Service) GetByOrg(
	ctx context.Context,
	oid uint,
) ([]*views.Group, *zenrpc.Error) {
	var flag bool
	me := requestContext.CurrentUser(ctx)
	if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&flag, oid, me.ID)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	} else if !flag {
		return nil, errors.New(errors.CantViewThisGroup, nil, nil)
	}

	var groups []*models.Group
	if err := s.db.Where(&models.Group{OrganizationID: oid}).Find(&groups).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	groupsView := make([]*views.Group, len(groups))
	for i := range groups {
		groupsView[i] = views.GroupFromModelShort(groups[i])
	}

	return groupsView, nil
}
