package notification

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//помечает уведомление как прочитанное.
//zenrpc:notifID id уведомления, которое необходимо отметить как прочитанное.
//zenrpc:48 notification not found. непросмотренное уведомление с таким ID не существует.
//zenrpc:return при удачном выполнении запроса возвращает список уведомлений, в котором лежат сами уведомления и их типы.
func (s *Service) Mark(ctx context.Context, notifID uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	notifSeen := models.UserNotificationSeen{UserID: me.ID, NotificationID: notifID}
	if err := s.db.
		Where(&notifSeen).
		Where("seen = ?", false).
		First(&notifSeen).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.NotificationNotFound, nil, nil) //48
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if notifSeen.UserID != me.ID {
		return nil, errors.New(errors.IncorrectNotificationID, nil, nil) //47
	}

	if err := s.db.Model(&models.UserNotificationSeen{}).Where("id = ?", notifSeen.ID).Update("seen", true).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}
