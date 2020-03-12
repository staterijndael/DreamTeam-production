package group

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//добавление человека в группу. операция доступа только для админов группы.
//zenrpc:gid id группы.
//zenrpc:targetPerson поль-ль, которого необходимо добавить.
//zenrpc:3 группа не найдена.
//zenrpc:38 пользователь уже в группе.
//zenrpc:40 операция доступа только для админов группы/организации.
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) AddUser(
	ctx context.Context,
	gid,
	targetPerson uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	// check not already member of group
	if group.Community.Contains(targetPerson) {
		return nil, errors.New(errors.AlreadyInGroup, nil, nil) // 38
	}

	me := requestContext.CurrentUser(ctx)
	if group.AdminID != me.ID && !group.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.InvitePersonWhileNotAdminAndNotAssociated, nil, nil) // 40
	}

	creation := s.db.Create(&models.MembershipOfCommunity{
		CommunityID: group.CommunityID,
		UserID:      targetPerson,
	})

	if creation.Error != nil {
		return nil, errors.New(errors.Internal, creation.Error, nil)
	}

	s.emitter.Emit(&events.UserAddedToGroup{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group:   gid,
		Added:   targetPerson,
		AddedBy: me.ID,
	})

	return common.ResultOK, nil
}
