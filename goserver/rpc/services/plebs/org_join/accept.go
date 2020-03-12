package org_join

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

//одобрение запроса на вступление в группу. только админ группы/организации имеет право на данную операцию.
//zenrpc:requestID id запроса.
//zenrpc:3 группа не найдена
//zenrpc:66 запрос с данным id не найден.
//zenrpc:68 запрос уже закрыт.
//zenrpc:69 данный пользователь не имеет права принять данный запрос.
//zenrpc:70 некорректно указана группа (не принадлежит организации, к которой относится запрос)
//zenrpc:return при удачном выполнении запроса возвращается сообщение "ok".
func (s *Service) Accept(
	ctx context.Context,
	requestID, groupID uint,
) (*common.CodeAndMessage, *zenrpc.Error) {
	var request models.OrgJoinRequest
	if err := s.db.First(&request, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.OrgJoinRequestNotFound, err, nil) // 66
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	me := requestContext.CurrentUser(ctx)
	if !request.Organization.Admins.Contains(me.ID) {
		return nil, errors.New(errors.CantViewOrgJoinRequest, nil, nil) // 69
	}

	if request.Status != models.Pending {
		return nil, errors.New(errors.OrgJoinRequestAlreadyClosed, nil, nil) // 68
	}

	var group models.Group
	if err := s.db.First(&group, groupID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupNotFound, err, nil) // 3
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if group.OrganizationID != request.OrganizationID {
		return nil, errors.New(errors.GroupIsNotMemberOfOrg, nil, nil) // 70
	}

	if err := s.db.
		Create(&models.MembershipOfCommunity{UserID: request.InitiatorID, CommunityID: group.CommunityID}).
		Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	if err := s.db.
		Model(&models.OrgJoinRequest{}).Where("id = ?", request.ID).
		Updates(map[string]interface{}{"status": models.Confirmed, "acceptor": me.ID, "group": groupID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	s.emitter.Emit(&events.OrgJoinRequestAccepted{
		EventBase: events.EventBase{
			Context: ctx,
		},
		Request: requestID,
	})

	return common.ResultOK, nil
}
