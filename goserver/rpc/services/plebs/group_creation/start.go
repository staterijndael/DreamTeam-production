package group_creation

import (
	"context"
	"dt/events"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/stores"
	"dt/utils"
	"dt/views"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

//открытие запроса на создание группы. без указания родителя создается запрос на создание корневой группы.
//требуется подтверждение запроса от админа род группы или слинкованного пользователя.
//если отправитель запроса яв-ся админом родительской группы запрос автоматически подтверждается, группа
//создается.
//zenrpc:oid id организации. при уведомлении заменяется сущностью.
//zenrpc:parent id родительской группы. заменяется сущностью.
//zenrpc:11 организация не найдена.
//zenrpc:22 невалидное название группы.
//zenrpc:8 невалидный никнейм
//zenrpc:3 родительская группа не найдена.
//zenrpc:4 пользователь не является членом родительской группы, не директор, не слинкован.
//zenrpc:6 данный пользователь пытается создать корневую группу в то время как ее может создать только слинкованный.
//zenrpc:return при удачном выполнении запроса от директора или слинкованного польз-я - возвращает созданную группу. при удачном выполнении запроса от участника родительской группы возвращаетя тело запроса о создании.
func (s *Service) Start(
	ctx context.Context,
	oid,
	parent uint,
	nickname,
	title,
	description string,
) (*views.GroupCreationRequest, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	if !utils.IsValidGroupTitle(title) {
		return nil, errors.New(errors.InvalidGroupTitle, nil, nil)
	}

	nicknameModel := models.Nickname{Value: nickname}
	if !nicknameModel.IsValid() {
		return nil, errors.New(errors.InvalidNickname, nil, nil)
	}

	var parentGroup models.Group
	if err := s.db.First(&parentGroup, parent).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if !parentGroup.Community.Contains(me.ID) {
		return nil, errors.New(errors.NotMemberOfParentCreateGroup, nil, nil) // 4
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

	request := models.GroupCreationRequest{
		RequestBase: models.RequestBase{
			Status:      models.Pending,
			InitiatorID: me.ID,
			AcceptorID:  nil,
		},
		Nickname:       nicknameModel,
		NicknameID:     nicknameModel.ID,
		OrganizationID: parentGroup.Organization.ID,
		ParentID:       &parentGroup.ID,
		Title:          title,
		Description:    description,
	}

	if err := s.db.Create(&request).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.First(&request, request.ID).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.GroupCreationRequestStarted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: request.ID,
	})

	return views.GroupCreationRequestFromModelShort(&request), nil
}
