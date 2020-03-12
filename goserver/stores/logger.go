package stores

import (
	"dt/logwrap"
	"fmt"
	"strings"
)

type logger struct{}

func (*logger) Print(values ...interface{}) {
	var b strings.Builder
	for _, v := range values {
		fmt.Fprint(&b, v, ", ")
	}

	logwrap.Error("[db]: %s", b.String()[:b.Len()-2])
}
