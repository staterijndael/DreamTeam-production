package views

import "dt/models"

type Community struct {
	ID      uint    `json:"id"`
	Members []*User `json:"members"`
}

func CommunityFromModelShort(community *models.Community) *Community {
	membersView := make([]*User, len(community.Members))
	for i := range community.Members {
		membersView[i] = UserViewFromModel(&community.Members[i].User)
	}

	return &Community{
		ID:      community.ID,
		Members: membersView,
	}
}

func CommunityFromModel(community *models.Community, members []*models.User) *Community {
	membersView := make([]*User, len(members))
	for i := range members {
		membersView[i] = UserViewFromModel(members[i])
	}

	return &Community{
		ID:      community.ID,
		Members: membersView,
	}
}
