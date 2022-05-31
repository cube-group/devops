package is

import "regexp"

//是否为合法url path（目录）
func DomainPath(v string) bool {
	var pattern = "^(/[0-9a-zA-Z_.-]{1,})+(/){0,1}$" //匹配url path
	matched, _ := regexp.Match(pattern, []byte(v))
	return matched
}
