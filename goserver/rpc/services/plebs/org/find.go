package org

import (
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/semrush/zenrpc"
	"strings"
)

//поиск организации
//zenrpc:text title|description|nickname
//zenrpc:page=0 пагинация. отсчет идет с 0. размер стр = 50
//zenrpc:return при удачном выполнении запроса возвращает полную информацию о пользователе.
func (s *Service) Find(text string, page uint64) ([]*views.Org, *zenrpc.Error) {
	queries := common.SplitBySpacesRegex.Split(strings.Trim(text, " \t\r\f\n"), -1)
	var organizations []*models.Organization
	if err := s.db.
		Where(new(models.Organization).FuzzyQuery(queries)).
		Or(
			"nickname in (?)",
			s.db.
				Table(new(models.Nickname).TableName()).
				Select("id").
				Where(new(models.Nickname).FuzzyQuery(queries)).
				SubQuery(),
		).
		Offset(page * 50).
		Limit(50).
		Find(&organizations).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	orgViews := make([]*views.Org, len(organizations))
	for i := range organizations {
		orgViews[i] = views.OrgViewFromModelShort(organizations[i])
	}

	return orgViews, nil
}
