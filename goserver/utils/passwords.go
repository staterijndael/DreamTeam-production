package utils

import (
	"regexp"
)

var (
	pwdRegex = regexp.MustCompile(`\w{6,}`)
)

func IsValidPassword(pwd string) bool {
	return pwdRegex.MatchString(pwd)
}
