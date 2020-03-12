package group

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
	"strings"
)

//поиск группы
//zenrpc:oid id организации
//zenrpc:text title|description|nickname
//zenrpc:page=0 пагинация. отсчет идет с 0. размер стр = 50
//zenrpc:18 организация не найдена
//zenrpc:65 вы не можете просматривать группы этой орг-ии
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о пользователе.
func (s *Service) Find(ctx context.Context, oid uint, text string, page uint64) ([]*views.Group, *zenrpc.Error) {
	queries := common.SplitBySpacesRegex.Split(strings.Trim(text, " \t\r\f\n"), -1)
	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrganizationNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if me.ID != org.DirectorID && !org.Admins.Contains(me.ID) {
		var isMember bool
		if err := s.db.Scopes(scopes.IsMemberOfAnyGroupOfOrg(&isMember, oid, me.ID)).Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		} else if !isMember {
			return nil, errors.New(errors.CantViewGroupsOfThisOrg, nil, nil)
		}
	}

	var groups []*models.Group
	if err := s.db.
		Scopes(scopes.FindGroup(page, oid, queries)).
		Find(&groups).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	groupViews := make([]*views.Group, len(groups))
	for i := range groups {
		groupViews[i] = views.GroupFromModelShort(groups[i])
	}

	return groupViews, nil
}
