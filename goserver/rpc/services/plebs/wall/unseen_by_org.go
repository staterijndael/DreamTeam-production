package wall

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

//получение списка непрочитанных данным пользователем уведомлений данной орг-ии.
//zenrpc:orgID id орг-ии, уведомления которой необходимо найти.
//zenrpc:loadAvatars=false при = true отдаются в том числе и аватарки
//zenrpc:return при удачном выполнении запроса возвращает список уведомлений, в котором лежат сами уведомления и их типы.
func (s *Service) UnseenByOrg(ctx context.Context, orgID uint, loadAvatars bool) ([]*utils.Container, *zenrpc.Error) {
	initiator := requestContext.CurrentUser(ctx)

	var unseenNotifications []*models.AdminOrgWallSeen
	if err := s.db.
		Where(`"user" = ?`, initiator.ID).
		Where("seen = ?", false).
		Find(&unseenNotifications).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	containers := make([]*utils.Container, 0)
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
				containers = append(containers, notif.ContainerizedView())
				return nil
			})
	}

	return containers, nil
}
