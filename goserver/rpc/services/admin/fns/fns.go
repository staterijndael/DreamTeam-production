package fns

import (
	"dt/models"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
	"time"
)

//устанавливает ключ API Федеральной Налоговой Службы.
//zenrpc:key ключ API ФНС
//zenrpc:expDate дата окончания срока дейстия ключа (timestamp)
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetKey(key string, expDate int64) (*common.CodeAndMessage, *zenrpc.Error) {
	err := s.fnsMgr.SetKey(models.FNSKey{
		Key:        key,
		ExpireDate: time.Unix(expDate, 0),
	})
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}

//TODO GetKey, GetAllKeys
