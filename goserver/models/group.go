package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"strings"
)

type Group struct {
	gorm.Model
	CreatorID      uint          `gorm:"column:creator"`
	AdminID        uint          `gorm:"column:admin"`
	Title          string        `gorm:"column:title"`
	Description    string        `gorm:"column:description"`
	ParentID       *uint         `gorm:"column:parent"`
	OrganizationID uint          `gorm:"column:organization"`
	ChildrenIDs    pq.Int64Array `gorm:"column:children;type:integer[]"`
	CommunityID    uint          `gorm:"column:community"`
	AvatarID       uint          `gorm:"column:avatar"`
	NicknameID     uint          `gorm:"column:nickname"`
	ChatID         uint          `gorm:"column:chat"`
	Nickname       Nickname      `gorm:"foreignkey:nickname"`
	Creator        User          `gorm:"foreignkey:creator"`
	Admin          User          `gorm:"foreignkey:admin"`
	Parent         *Group        `gorm:"foreignkey:id;PRELOAD:false"`
	Organization   Organization  `gorm:"foreignkey:organization"`
	Community      Community     `gorm:"foreignkey:community"`
	Avatar         File          `gorm:"foreignkey:avatar"`
}

func (g *Group) TableName() string {
	return "groups"
}
func (g *Group) ZDBIdxTypeDefinition() string {
	return "id " + g.ZDBIdxIDType() + ", title text, description text"
}

func (g *Group) ZDBRowBuilder() string {
	return "(id)::" + g.ZDBIdxIDType() + ", title, description"
}

func (g *Group) ZDBIdxIDType() string {
	return "bigint"
}

func (*Group) FuzzyQuery(queries []string) string {
	format := "dsl.fuzzy('title', '%s', fuzziness => %d), dsl.wildcard('title', '*%s*'), " +
		"dsl.wildcard('description', '*%s*'), "
	var b strings.Builder
	fmt.Fprintf(&b, "%s ==> dsl.or(", new(Group).TableName())
	var fuzziness uint32
	for i := range queries {
		if len(queries[i]) > 2 {
			fuzziness = 2
		} else {
			fuzziness = 1
		}

		fmt.Fprintf(
			&b,
			format,
			queries[i],
			fuzziness,
			queries[i],
			queries[i],
		)
	}

	return b.String()[:b.Len()-2] + ")"
}
