package fns

import (
	"dt/managers/fns"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//получение информации об орагнизации из фнс
//zenrpc:text шаблон поиска.
//zenrpc:3 fns error. ошибка с сервисом фнс.
//zenrpc:return при удачном выполнении запроса возвращает FNSOrganizationInfo.
func (s *Service) Search(text string) (*fns.OrganizationInfo, *zenrpc.Error) {
	res, err := s.fnsMgr.Search(text)
	if err != nil {
		fnsError, ok := err.(*fns.Error)
		if ok {
			return nil, errors.New(errors.FNSError, fnsError, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	return res, nil
}
