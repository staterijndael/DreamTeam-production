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

//Позволяет пользователю покинуть группу.
//zenrpc:gid id группы, которую пользватель должен покинуть
//zenrpc:3 группа не найдена
//zenrpc:58 админ группы не может её покинуть
//zenrpc:59 пользователь не состоит в этой группе
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) Leave(ctx context.Context, gid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var group models.Group
	if err := s.db.First(&group, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if me.ID == group.AdminID {
		return nil, errors.New(errors.LeaveGroupWhileAdmin, nil, nil) // 58
	}

	if !group.Community.Contains(me.ID) {
		return nil, errors.New(errors.NotMemberOfGroup, nil, nil) // 59
	}

	if err := s.db.Scopes(scopes.RemoveFromGroup(me.ID, &group)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.UserLeftGroup{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group: gid,
		User:  me.ID,
	})

	return common.ResultOK, nil
}
