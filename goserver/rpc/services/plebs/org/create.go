package org

import (
	"context"
	"dt/logwrap"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/stores"
	"dt/utils"
	"dt/views"
	"github.com/semrush/zenrpc"
)

//создание организации не привязанной к фнс.
//один из параметров названия или инн обязательны. не могут указываться одновременно.
//jsonrpc notification с данными запроса отправляется по другим соединениям данного пользователя.
//параметры уведомления - id уже созданной орг-ии, title, description, nickname.
//zenrpc:nickname nickname организации
//zenrpc:6 nickname is busy. данный nickname уже занят другим пользователем
//zenrpc:9 invalid organization title. название орг-ии не соответствует формату
//zenrpc:10 invalid organization description. описание орг-ии не соответствует формату
//zenrpc:14 invalid organization nickname. nickname орг-ии не соответствует формату
//zenrpc:return при удачном выполнении запроса возвращает полную информацию об орг-ии.
func (s *Service) Create(
	ctx context.Context,
	nickname,
	title,
	description string,
) (*views.Org, *zenrpc.Error) {
	if title != "" && !utils.IsValidOrgTitle(title) {
		return nil, errors.New(errors.InvalidOrgTitle, nil, nil)
	}

	if !utils.IsValidOrgDescription(description) {
		return nil, errors.New(errors.InvalidOrgDescription, nil, nil)
	}

	nicknameModel := models.Nickname{Value: nickname}
	if !nicknameModel.IsValid() {
		return nil, errors.New(errors.InvalidOrgNickname, nil, nil)
	}

	nicknameModel.Validate()
	if err := s.db.Create(&nicknameModel).Error; err != nil {
		if stores.IsDuplicate(err) {
			return nil, errors.New(errors.NicknameIsBusy, err, nil)
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	organization := models.Organization{
		Title:       title,
		Description: description,
		DirectorID:  me.ID,
		NicknameID:  nicknameModel.ID,
		Nickname:    nicknameModel,
		FNS:         nil,
		AvatarID:    stores.DefaultAvatars.Org.ID,
		Admins: models.Community{
			Members: []models.MembershipOfCommunity{
				{UserID: me.ID},
			},
		},
	}

	if creation := s.db.Create(&organization); creation.Error != nil {
		if stores.IsDuplicate(creation.Error) {
			return nil, errors.New(errors.NicknameIsBusy, creation.Error, nil) // 6
		}

		logwrap.Debug("err: %v", creation.Error)

		return nil, errors.New(errors.Internal, creation.Error, nil)
	}

	if err := s.db.First(&organization, organization.ID).Error; err != nil {
		logwrap.Debug("err: %v", err)
	}

	orgView := views.OrgViewFromModelShort(&organization)
	return orgView, nil
}
