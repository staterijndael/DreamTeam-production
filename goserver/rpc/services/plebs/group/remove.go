package group

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//Удаление пользователя из группы. Право на данную операцию имеет только админ группы и админы организации.
//zenrpc:gid id группы, из которой необходиму удалить пользователя
//zenrpc:uid id удаляемого пользователя
//zenrpc:3 group not found
//zenrpc:37 not enough rights
//zenrpc:58 админ группы не может её покинуть
//zenrpc:59 пользователь не состоит в этой группе
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) RemoveUser(ctx context.Context, gid, uid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if me.ID != group.AdminID && !group.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.NotAssociated, nil, nil) // 37
	}

	if uid == group.AdminID {
		return nil, errors.New(errors.LeaveGroupWhileAdmin, nil, nil) // 58
	}

	if !group.Community.Contains(uid) {
		return nil, errors.New(errors.NotMemberOfGroup, nil, nil) // 59
	}

	if err := s.db.Scopes(scopes.RemoveFromGroup(uid, &group)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserRemovedFromGroup{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group:     group.ID,
		Removed:   uid,
		RemovedBy: me.ID,
	})

	return common.ResultOK, nil
}
