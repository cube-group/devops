// Author: chenqionghe
// Time: 2018-10
// 自定义验证器，提供自定义验证规则，所有增加的规则均在myValidators变量里

package validator

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

//todo 注册自定义验证规则
var myValidators = map[string]validator.Func{
	"is-cn":           IsCn,
	"is-project-name": IsProjectName,
	"is-config-name":  IsConfigName,
	"is-domain":       IsDomain,
	"is-url":          IsUrl,
	"is-git-address":  IsGitAddress,
}

func init() {
	binding.Validator = new(defaultValidator)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册自定义验证规则
		for alias, function := range myValidators {
			v.RegisterValidation(alias, function)
		}
	}
}

//是否是中文
func IsCn(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile("^[\u4e00-\u9fa5]+$")
	return reg.MatchString(fl.Field().String())
}

//是否是项目名
func IsProjectName(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`)
	return reg.MatchString(fl.Field().String())
}

//是否是项目名
func IsConfigName(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile(`^([a-z]|\._-)?([-a-z0-9]*[a-z0-9])?`)
	return reg.MatchString(fl.Field().String())
}

//是否是域名
func IsDomain(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile(`^((https?|ftp|news):\/\/)?([a-z]([a-z0-9\-]*[\.。])+([a-z]{2}|aero|arpa|biz|com|coop|edu|gov|info|int|jobs|mil|museum|name|nato|net|org|pro|travel)|(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))(\/[a-z0-9_\-\.~]+)*(\/([a-z0-9_\-\.]*)(\?[a-z0-9+_\-\.%=&]*)?)?(#[a-z][a-z0-9_]*)?$`)
	return reg.MatchString(fl.Field().String())
}

//是否是url地址
func IsUrl(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile(`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	return reg.MatchString(fl.Field().String())
}

//是否是git地址
func IsGitAddress(fl validator.FieldLevel) bool {
	var reg = regexp.MustCompile(`((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?`)
	return reg.MatchString(fl.Field().String())
}
