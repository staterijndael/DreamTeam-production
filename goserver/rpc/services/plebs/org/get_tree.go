package org

import (
	"context"
	"dt/logwrap"
	"dt/models"
	"dt/requestContext"
	"dt/rpc/services/errors"
	"dt/views"
	"github.com/golang-collections/collections/stack"
	"github.com/lib/pq"
	"github.com/semrush/zenrpc"
)

//получение дерева групп организации.
//zenrpc:oid id орг-ии.
//zenrpc:36 пользователь не имеет доступа к ресурсу.
//zenrpc:return при удачном выполнении запроса возвращает дерево орг-ии.
func (s *Service) GetTree(
	ctx context.Context,
	oid uint,
) (interface{}, *zenrpc.Error) {
	me := requestContext.CurrentUser(ctx)
	var org models.Organization
	if err := s.db.First(&org, oid).Error; err != nil {
		logwrap.Debug("[gettree::findOrg]: %s", err.Error())
		return nil, errors.New(errors.Internal, err, nil)
	}

	if !org.Admins.Contains(me.ID) {
		return nil, errors.New(errors.TreeAccessDenied, nil, nil) // 36
	}

	orgAva, err := views.FileViewFromModel(&org.Avatar)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	cache := &BuildingTreeCache{
		Roots:        nil,
		AllGroups:    make([]*models.Group, 0),
		GroupsMap:    make(map[uint]*models.Group),
		Org:          &org,
		OrgAvatar:    orgAva,
		Users:        make(map[uint]*views.User),
		UserAvatars:  make(map[uint]*views.File),
		GroupAvatars: make(map[uint]*views.File),
	}

	//groups := make([]*models.GroupID, 0)
	//cache.AllGroups =

	if err := s.db.Where("organization = ?", oid).Find(&cache.AllGroups).Error; err != nil {
		logwrap.Debug("[gettree::findAllGroups]: %s", err.Error())
		return nil, errors.New(errors.Internal, err, nil)
	}

	dirAva, err := views.FileViewFromModel(&org.Director.Avatar)
	if err != nil {
		return nil, errors.New(errors.Internal, err, nil)
	}

	cache.UserAvatars[org.DirectorID] = dirAva

	var director models.User
	if err := s.db.First(&director, org.DirectorID).Error; err != nil {
		logwrap.Debug("[gettree::findDirector]: %s", err.Error())
		return nil, errors.New(errors.Internal, err, nil)
	}

	cache.Users[director.ID] = views.UserViewFromModel(&director)

	for _, g := range cache.AllGroups {
		if g.ParentID == nil {
			cache.Roots = append(cache.Roots, g)
		}

		cache.GroupsMap[g.ID] = g
		groupAva, err := views.FileViewFromModel(&g.Avatar)
		if err != nil {
			return nil, errors.New(errors.Internal, err, nil)
		}
		cache.GroupAvatars[g.ID] = groupAva

		if _, ok := cache.Users[g.AdminID]; !ok {
			var u models.User
			if err := s.db.First(&u, g.AdminID).Error; err != nil {
				logwrap.Debug(
					"[gettree::findAdmin]: group: %d, adminID: %d err: %s", g.ID, g.AdminID, err.Error())
				return nil, errors.New(errors.Internal, err, nil)
			}

			cache.Users[u.ID] = views.UserViewFromModel(&u)
			userAva, err := views.FileViewFromModel(&u.Avatar)
			if err != nil {
				return nil, errors.New(errors.Internal, err, nil)
			}
			cache.UserAvatars[g.AdminID] = userAva
		}
	}

	trees := make([]*views.GroupTreeNode, len(cache.Roots))
	for i := range cache.Roots {
		trees[i] = getTree(cache.Roots[i], cache)
	}

	return views.TreeRootFromModel(cache.Org,
		cache.Users[cache.Org.DirectorID],
		trees,
		orgAva,
	), nil
}

func getTree(root *models.Group, cache *BuildingTreeCache) *views.GroupTreeNode {
	tree := views.TreeNodeFromModel(root,
		cache.Users[root.AdminID],
		cache.GroupAvatars[root.ID],
	)

	groupsStack := stack.New()

	pushChildren := func(ch pq.Int64Array) {
		for _, c := range ch {
			groupsStack.Push(c)
		}
	}

	pushChildren(root.ChildrenIDs)

	for node := new(views.GroupTreeNode); groupsStack.Len() > 0; tree = node {
		group := cache.GroupsMap[uint(groupsStack.Pop().(int64))]

		for *group.ParentID != tree.Group.ID {
			tree = tree.Parent
		}

		node = views.TreeNodeFromModel(group,
			cache.Users[group.AdminID],
			cache.GroupAvatars[group.ID],
		)
		node.Parent = tree

		tree.Children = append(tree.Children, node)
		pushChildren(group.ChildrenIDs)
	}

	for ; tree.Parent != nil; tree = tree.Parent {
	}

	return tree
}

type BuildingTreeCache struct {
	Roots        []*models.Group
	AllGroups    []*models.Group
	GroupsMap    map[uint]*models.Group
	Org          *models.Organization
	OrgAvatar    *views.File
	Users        map[uint]*views.User
	UserAvatars  map[uint]*views.File
	GroupAvatars map[uint]*views.File
}
