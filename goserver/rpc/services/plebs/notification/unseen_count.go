package notification

import (
	"context"
	"dt/models"
	"dt/notification"
	"dt/recovery_wrapper"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"github.com/semrush/zenrpc"
)

//Получение количества непрочитанных уведомлений пользователя, отправившего запрос
//zenrpc:return при удачном выполнении запроса возвращает количество непрочитанных уведомлений
func (s *Service) UnseenCount(ctx context.Context) (*int, *zenrpc.Error) {
	var count int
	me := requestContext.CurrentUser(ctx)
	var notifications []*models.UserNotificationSeen
	if err := s.db.
		Where(`"user" = ?`, me.ID).
		Where(`seen = ?`, false).
		Find(&notifications).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	wrapper := recoveryWrapper.Wrapper{}
	for _, n := range notifications {
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
				count++
				return nil
			})
	}

	return &count, nil
}
