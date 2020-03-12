package notification

import (
	"context"
	"dt/models"
	"dt/notification"
	"dt/recovery_wrapper"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/utils"
	"github.com/semrush/zenrpc"
)

//получение списка непрочитанных уведомлений пользователя, отправившего запрос.
//zenrpc:loadAvatars=true при = true отдаются в том числе и аватарки
//zenrpc:return при удачном выполнении запроса возвращает список уведомлений, в котором лежат сами уведомления и их типы.
func (s *Service) Unseen(ctx context.Context, loadAvatars bool) ([]*utils.Container, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var unseenNotifications []*models.UserNotificationSeen
	if err := s.db.
		Where(`"user" = ?`, me.ID).
		Where("seen = ?", false).
		Find(&unseenNotifications).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	containers := make([]*utils.Container, 0)
	wrapper := recoveryWrapper.Wrapper{}
	for _, n := range unseenNotifications {
		notif, err := notification.GetNotification(n.Notification)
		if err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}

		wrapper.
			Clear().
			Do(func() error {
				return notif.Load(s.db, &n.Notification)
			}).
			Do(func() error {
				containers = append(containers, notif.ContainerizedView())
				return nil
			})

		//if err := notif.Load(s.db, &n.Notification); err != nil {
		//	//TODO issue #4
		//	if gorm.IsRecordNotFoundError(err) {
		//		continue
		//	}
		//
		//	return nil, errors.New(errors.Internal, err, nil)
		//}
		//
		//if loadAvatars {
		//	if err := notif.LoadFiles(); err != nil {
		//		return nil, errors.New(errors.Internal, err, nil)
		//	}
		//}
	}

	return containers, nil
}
