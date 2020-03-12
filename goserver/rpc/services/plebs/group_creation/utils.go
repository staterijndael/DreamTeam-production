package group_creation

import (
	"context"
	"dt/events"
	"dt/managers/eventEmitter"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"dt/scopes"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

func changeStatusIfLinkedOrDirector(
	ctx context.Context,
	status models.RequestStatus,
	requestID uint,
	errCode int,
	s *Service,
	emitter *eventEmitter.EventEmitter,
) (*common.CodeAndMessage, *zenrpc.Error, *uint, *models.GroupCreationRequest) {
	var req models.GroupCreationRequest
	if err := s.db.First(&req, requestID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.CreationRequestNotFound, err, nil), nil, nil
		}

		return nil, errors.New(errors.Internal, err, nil), nil, nil
	}

	if req.Status != models.Pending {
		return nil, errors.New(errCode, nil, nil), nil, nil
	}

	me := requestContext.CurrentUser(ctx)
	if req.Organization.DirectorID != me.ID &&
		!req.Organization.Admins.Contains(me.ID) &&
		(req.Parent == nil || req.Parent.AdminID != me.ID) {
		return nil, errors.New(errCode, nil, nil), nil, nil
	}

	if err := s.db.
		Model(&models.GroupCreationRequest{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{"status": status, "acceptor": me.ID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil), nil, nil
	}

	var grID uint
	if status == models.Confirmed {
		var gr models.Group
		if err := s.db.Scopes(scopes.CreateGroup(
			&gr,
			req.InitiatorID,
			req.OrganizationID,
			req.Title,
			req.Description,
			req.Nickname,
			req.ParentID)).
			Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil), nil, nil
		} else {
			grID = gr.ID
			emitter.Emit(&events.GroupCreated{
				EventBase: events.EventBase{
					Context: ctx,
				},
				Group:   gr.ID,
				Creator: req.InitiatorID,
			})
		}
	}

	return common.ResultOK, nil, &grID, &req
}
