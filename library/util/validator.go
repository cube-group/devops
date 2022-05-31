package util

import "regexp"

func IsWord(value string) bool {
	result,_:= regexp.MatchString("^\\w+$", value)
	return result
}
