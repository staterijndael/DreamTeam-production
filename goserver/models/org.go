package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

type Organization struct {
	gorm.Model
	Title       string          `gorm:"column:title"`
	Description string          `gorm:"column:description"`
	DirectorID  uint            `gorm:"column:director"`
	NicknameID  uint            `gorm:"column:nickname"`
	Nickname    Nickname        `gorm:"foreignkey:nickname"`
	FNS         *postgres.Jsonb `gorm:"column:fns"`
	Trusted     bool            `gorm:"column:trusted"`
	AvatarID    uint            `gorm:"column:avatar"`
	CommunityID uint            `gorm:"column:community"`
	Admins      Community       `gorm:"foreignkey:community"`
	Director    User            `gorm:"foreignkey:director"`
	Avatar      File            `gorm:"foreignkey:avatar"`
}

func (o *Organization) TableName() string {
	return "organizations"
}

func (o *Organization) ZDBIdxTypeDefinition() string {
	return "id " + o.ZDBIdxIDType() + ", title text, description text"
}

func (o *Organization) ZDBRowBuilder() string {
	return "(id)::" + o.ZDBIdxIDType() + ", title, description"
}

func (o *Organization) ZDBIdxIDType() string {
	return "bigint"
}

func (*Organization) FuzzyQuery(queries []string) string {
	format := "dsl.fuzzy('title', '%s', fuzziness => %d), dsl.wildcard('title', '*%s*'), " +
		"dsl.wildcard('description', '*%s*'), "
	var b strings.Builder
	fmt.Fprintf(&b, "%s ==> dsl.or(", new(Organization).TableName())
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
