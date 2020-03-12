package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
)

type Nickname struct {
	gorm.Model
	Value string `gorm:"column:value;unique"`
}

func (n *Nickname) TableName() string {
	return "nicknames"
}

func (n *Nickname) ZDBIdxTypeDefinition() string {
	return "id " + n.ZDBIdxIDType() + ", value text"
}

func (n *Nickname) ZDBRowBuilder() string {
	return "(id)::" + n.ZDBIdxIDType() + ", value"
}

func (n *Nickname) ZDBIdxIDType() string {
	return "bigint"
}

var nickRegexp = regexp.MustCompile(`[a-z_]+`)

func (n *Nickname) IsValid() bool {
	return nickRegexp.MatchString(n.Value)
}

func (n *Nickname) Validate() {
	n.Value = strings.ToLower(n.Value)
}

func (*Nickname) FuzzyQuery(queries []string) string {
	var b strings.Builder
	var fuzziness uint32
	fmt.Fprintf(&b, "%s ==> dsl.or(", new(Nickname).TableName())
	for i := range queries {
		if len(queries[i]) > 3 {
			fuzziness = 2
		} else {
			fuzziness = 1
		}

		fmt.Fprintf(
			&b,
			"dsl.fuzzy('value', '%s', fuzziness => %d), dsl.wildcard('value', '*%s*'), ",
			queries[i],
			fuzziness,
			queries[i],
		)
	}

	return b.String()[:b.Len()-2] + ")"
}
