package views

import (
	"dt/models"
)

type Org struct {
	ID              uint          `json:"id"`
	Date            int64         `json:"date"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Director        uint          `json:"director"`
	Nickname        string        `json:"nickname"`
	AssociatedUsers *Community    `json:"associatedUsers"`
	FNS             interface{}   `json:"fns,omitempty"`
	AvatarMetaInfo  *FileMetaInfo `json:"avatarMetaInfo"`
}

func OrgViewFromModelShort(org *models.Organization) *Org {
	return &Org{
		ID:              org.ID,
		Date:            org.CreatedAt.Unix(),
		Title:           org.Title,
		Description:     org.Description,
		Director:        org.DirectorID,
		Nickname:        org.Nickname.Value,
		AssociatedUsers: CommunityFromModelShort(&org.Admins),
		FNS:             org.FNS,
		AvatarMetaInfo:  FileMetaInfoViewFromModel(&org.Avatar),
	}
}

func OrgViewFromModel(organization *models.Organization, associated []*models.User, com *models.Community) *Org {
	return &Org{
		ID:              organization.ID,
		Date:            organization.CreatedAt.Unix(),
		Title:           organization.Title,
		Description:     organization.Description,
		Director:        organization.DirectorID,
		Nickname:        organization.Nickname.Value,
		AssociatedUsers: CommunityFromModel(com, associated),
		FNS:             organization.FNS,
		AvatarMetaInfo:  FileMetaInfoViewFromModel(&organization.Avatar),
	}
}
