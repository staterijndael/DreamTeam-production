package group_join

import (
	"context"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/common"
	"dt/rpc/services/errors"
	"github.com/jinzhu/gorm"
	"github.com/semrush/zenrpc"
)

func changeStatusOfGroupJoinRequest(
	ctx context.Context,
	db *gorm.DB,
	st models.RequestStatus,
	errCode int,
	reqID uint,
	checkPermissionsFunc permissionChecker,
) (*common.CodeAndMessage, *zenrpc.Error) {
	var req models.GroupJoinRequest
	if err := db.First(&req, reqID).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.New(errors.GroupJoinRequestNotFound, err, nil) // 41
		}

		return nil, errors.New(errors.Internal, err, nil)
	}

	if req.Status != models.Pending {
		return nil, errors.New(errors.GroupJoinRequestAlreadyClosed, nil, nil) // 45
	}

	me := requestContext.CurrentUser(ctx)
	if !checkPermissionsFunc(me.ID, &req) {
		return nil, errors.New(errCode, nil, nil)
	}

	if st == models.Confirmed {
		creation := db.Create(&models.MembershipOfCommunity{UserID: req.InitiatorID, CommunityID: req.Group.CommunityID})
		if err := creation.Error; err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}
	}

	if err := db.
		Model(&models.GroupJoinRequest{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{"status": st, "acceptor": me.ID}).Error; err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	return common.ResultOK, nil
}

func isInitiator(userID uint, req *models.GroupJoinRequest) bool {
	return req.InitiatorID == userID
}

func isAdminOfGroupOrOrg(userID uint, req *models.GroupJoinRequest) bool {
	return req.Group.AdminID == userID || req.Group.Organization.Admins.Contains(userID)
}

type permissionChecker func(senderID uint, r *models.GroupJoinRequest) bool
