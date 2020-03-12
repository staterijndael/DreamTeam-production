package views

import (
	"dt/models"
)

type GroupCreationRequest struct {
	ID           uint   `json:"id"`
	Status       string `json:"status"`
	Initiator    *User  `json:"initiator"`
	Acceptor     *User  `json:"acceptor,omitempty"`
	Organization *Org   `json:"organization"`
	Parent       *Group `json:"parent,omitempty"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Nickname     string `json:"nickname"`
}

type GroupJoinRequest struct {
	ID        uint   `json:"id"`
	Status    string `json:"status"`
	Initiator *User  `json:"initiator"`
	Acceptor  *User  `json:"acceptor,omitempty"`
	Group     *Group `json:"group"`
}

type OrgJoinRequest struct {
	ID           uint   `json:"id"`
	Status       string `json:"status"`
	Initiator    *User  `json:"initiator"`
	Acceptor     *User  `json:"acceptor,omitempty"`
	Group        *Group `json:"group,omitempty"`
	Organization *Org   `json:"organization"`
}

func OrgJoinRequestFromModel(req *models.OrgJoinRequest) *OrgJoinRequest {
	return &OrgJoinRequest{
		ID:           req.ID,
		Status:       string(req.Status),
		Initiator:    UserViewFromModel(&req.Initiator),
		Acceptor:     UserViewFromModel(req.Acceptor),
		Group:        GroupFromModelShort(req.Group),
		Organization: OrgViewFromModelShort(&req.Organization),
	}
}

func GroupJoinRequestFromModel(
	req *models.GroupJoinRequest,
	initiator,
	acceptor *models.User,
	group *models.Group,
) *GroupJoinRequest {
	return &GroupJoinRequest{
		ID:        req.ID,
		Status:    string(req.Status),
		Initiator: UserViewFromModel(initiator),
		Acceptor:  UserViewFromModel(acceptor),
		Group:     GroupFromModelShort(group),
	}
}

func GroupJoinRequestFromModelShort(
	req *models.GroupJoinRequest,
) *GroupJoinRequest {
	return &GroupJoinRequest{
		ID:        req.ID,
		Status:    string(req.Status),
		Initiator: UserViewFromModel(&req.Initiator),
		Acceptor:  UserViewFromModel(req.Acceptor),
		Group:     GroupFromModelShort(&req.Group),
	}
}

func GroupCreationRequestFromModelShort(
	req *models.GroupCreationRequest,
) *GroupCreationRequest {
	return &GroupCreationRequest{
		ID:           req.ID,
		Status:       string(req.Status),
		Initiator:    UserViewFromModel(&req.Initiator),
		Acceptor:     UserViewFromModel(req.Acceptor),
		Organization: OrgViewFromModelShort(&req.Organization),
		Parent:       GroupFromModelShort(req.Parent),
		Title:        req.Title,
		Description:  req.Description,
		Nickname:     req.Nickname.Value,
	}
}
