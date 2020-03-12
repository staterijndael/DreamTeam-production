package wall

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//помечает уведомление как прочитанное.
//zenrpc:notifID id уведомления, которое необходимо отметить как прочитанное.
//zenrpc:48 notification not found. непросмотренное уведомление с таким ID не существует.
//zenrpc:return при удачном выполнении запроса возвращает список уведомлений, в котором лежат сами уведомления и их типы.
func (s *Service) Mark(ctx context.Context, notifID uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var notifSeen models.AdminOrgWallSeen
	if err := s.db.Scopes(scopes.UnseenAOWS(me.ID, notifID)).First(&notifSeen).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.NotificationNotFound, nil, nil) //48
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.Model(&models.AdminOrgWallSeen{}).Where("id = ?", notifSeen.ID).Update("seen", true).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
