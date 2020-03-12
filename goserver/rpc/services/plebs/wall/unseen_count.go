package wall

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
	initiator := requestContext.CurrentUser(ctx)

	var unseenNotifications []*models.AdminOrgWallSeen
	if err := s.db.
		Where(`"user" = ?`, initiator.ID).
		Where("seen = ?", false).
		Find(&unseenNotifications).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	count := 0
	wrapper := recoveryWrapper.Wrapper{}
	for _, n := range unseenNotifications {
		if !n.Wall.Organization.Admins.Contains(initiator.ID) {
			continue
		}

		notif, err := notification.GetNotification(n.Wall.Notification)
		if err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}

		wrapper.
			Clear().
			Do(func() error {
				return notif.Load(s.db, &n.Wall.Notification)
			}).
			Do(func() error {
				count++
				return nil
			})
	}

	return &count, nil
}
