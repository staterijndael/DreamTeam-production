package group

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/scopes"
	"dt/stores"
	"dt/utils"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//создание дочерней группы. только админ родительской группы имеет право на эту операцию.
//zenrpc:oid id организации
//zenrpc:parent=-1 id родительской группы
//zenrpc:8 invalid nickname
//zenrpc: невалидное название группы.
//zenrpc:3 родительская группа не найдена.
//zenrpc:5 нет прав на выполнение операции. (не админ этой организации и не админ родительской)
//zenrpc:return созданная группа.
func (s *Service) Create(
	ctx context.Context,
	oid uint,
	parent int,
	nickname,
	title,
	description string,
) (*views.Group, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	if !utils.IsValidGroupTitle(title) {
		return nil, errors.New(errors.InvalidGroupTitle, nil, nil) //
	}

	nicknameModel := models.Nickname{Value: nickname}
	if !nicknameModel.IsValid() {
		return nil, errors.New(errors.InvalidNickname, nil, nil)
	}

	var org models.Organization
	err := s.db.First(&org, oid).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgNotFound, err, nil) // 11
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	var parentID *uint = nil
	if parent != -1 {
		var parentGroup models.Group
		if err := s.db.First(&parentGroup, parent).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, errors.New(errors.GroupNotFound, err, nil) // 3
			}

			return nil, errors.New(errors.Internal, err, nil)
		}

		if me.ID != parentGroup.AdminID && org.DirectorID != me.ID && !org.Admins.Contains(me.ID) {
			return nil, errors.New(errors.CantCreateSubgroup, nil, nil) // 5
		}

		parentID = &parentGroup.ID
	} else {
		if org.DirectorID != me.ID && !org.Admins.Contains(me.ID) {
			return nil, errors.New(errors.CantCreateSubgroup, nil, nil) // 5
		}
	}

	if err := s.db.Create(&nicknameModel).Error; err != nil {
		if stores.IsDuplicate(err) {
			return nil, errors.New(errors.NicknameIsBusy, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.First(&nicknameModel, nicknameModel.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	var gr models.Group
	if err := s.db.Scopes(scopes.CreateGroup(&gr, me.ID, oid, title, description, nicknameModel, parentID)).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.GroupCreated{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Group:   gr.ID,
		Creator: me.ID,
	})

	return views.GroupFromModelShort(&gr), nil
}
