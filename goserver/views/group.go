package views

import "dt/models"

type Group struct {
	ID             uint          `json:"id"`
	Creator        uint          `json:"creator"`
	Admin          uint          `json:"admin"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Parent         *uint         `json:"parent,omitempty"`
	Organization   *Org          `json:"organization"`
	Children       []int64       `json:"children"`
	CreatedAt      int64         `json:"createdAt"`
	Community      *Community    `json:"community"`
	AvatarMetaInfo *FileMetaInfo `json:"avatarMetaInfo"`
	Nickname       string        `json:"nickname"`
	ChatID         uint          `json:"chatID"`
}

func GroupFromModelShort(g *models.Group) *Group {
	if g == nil {
		return nil
	}

	return &Group{
		ID:             g.ID,
		Creator:        g.CreatorID,
		Admin:          g.AdminID,
		Title:          g.Title,
		Description:    g.Description,
		Parent:         g.ParentID,
		Community:      CommunityFromModelShort(&g.Community),
		CreatedAt:      g.CreatedAt.Unix(),
		Organization:   OrgViewFromModelShort(&g.Organization),
		Children:       g.ChildrenIDs,
		AvatarMetaInfo: FileMetaInfoViewFromModel(&g.Avatar),
		Nickname:       g.Nickname.Value,
		ChatID:         g.ChatID,
	}
}
