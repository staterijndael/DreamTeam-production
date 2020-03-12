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

//получение кол-ва непрочитанных данным пользователем уведомлений данной орг-ии.
//zenrpc:orgID id орг-ии, уведомления которой необходимо найти.
//zenrpc:return при удачном выполнении запроса возвращает кол-во уведомлений.
func (s *Service) UnseenByOrgCount(ctx context.Context, orgID uint) (*uint64, *zenrpc.Error) {
	initiator := requestContext.CurrentUser(ctx)

	var unseenNotifications []*models.AdminOrgWallSeen
	if err := s.db.
		Where(`"user" = ?`, initiator.ID).
		Where("seen = ?", false).
		Find(&unseenNotifications).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var counter uint64 = 0
	wrapper := recoveryWrapper.Wrapper{}
	for _, n := range unseenNotifications {
		if n.Wall.OrganizationID != orgID || !n.Wall.Organization.Admins.Contains(initiator.ID) {
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
				counter++
				return nil
			})
	}

	return &counter, nil
}
