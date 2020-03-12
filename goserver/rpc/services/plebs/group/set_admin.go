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

//zenrpc:23 group not found
//zenrpc:74 только админ группы/орги может назначать нового админа
//zenrpc:75 данный пользователь не состоит в группе
//zenrpc:return возвращает "ok" в случае успеха
func (s *Service) SetAdmin(ctx context.Context, gid, uid uint) (*common.CodeAndMessage, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var gr models.Group
	if err := s.db.First(&gr, gid).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if gr.AdminID != me.ID && !gr.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.OnlyAdminOfGroupCanSetNewAdmin, nil, nil)
	}

	var user models.User
	for _, m := range gr.Community.Members {
		if m.UserID == uid {
			user = m.User
			break
		}
	}

	if user.ID == 0 {
		return nil, errors.New(errors.UserNotInGroup, nil, nil)
	}

	oldAdmin := gr.AdminID
	if err := s.db.
		Model(&models.Group{}).Where("id = ?", gr.ID).
		Update("admin", uid).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.NewGroupAdmin{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group:      gid,
		NewAdmin:   uid,
		OldAdmin:   oldAdmin,
		AssignedBy: me.ID,
	})

	return common.ResultOK, nil
}
