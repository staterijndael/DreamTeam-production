package user

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
	"time"
)

//установливает дату рождения
//.jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//zenrpc:bDate дата рождения в формате unix
//zenrpc:return при удачном выполнении запроса возвращает сообщение "ok".
func (s *Service) SetBDate(ctx context.Context, bDate int64) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	newTime := time.Unix(bDate, 0)
	me.BDate = &newTime

	if err := s.db.Model(&models.User{}).Where("id = ?", me.ID).Update("bdate", me.BDate).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
