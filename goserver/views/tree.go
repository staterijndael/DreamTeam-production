package views

import "dt/models"

type GroupTreeNode struct {
	Parent      *GroupTreeNode   `json:"-"`
	Group       *models.Group    `json:"-"`
	ID          uint             `json:"id"`
	Title       string           `json:"name"`
	Admin       *User            `json:"admin"`
	GroupAvatar *File            `json:"avatar"`
	Children    []*GroupTreeNode `json:"children"`
}

func TreeNodeFromModel(gr *models.Group, admin *User, grAva *File) *GroupTreeNode {
	return &GroupTreeNode{
		Parent:      nil,
		Group:       gr,
		ID:          gr.ID,
		Title:       gr.Title,
		Admin:       admin,
		Children:    make([]*GroupTreeNode, 0),
		GroupAvatar: grAva,
	}
}

func TreeRootFromModel(
	org *models.Organization,
	director *User,
	roots []*GroupTreeNode,
	grAva *File,
) *GroupTreeNode {
	return &GroupTreeNode{
		Parent:      nil,
		Group:       nil,
		ID:          org.ID,
		Title:       org.Title,
		Admin:       director,
		Children:    roots,
		GroupAvatar: grAva,
	}
}
