package consts

const (
	IngressControllerTypePublic = "public" //公网类入口
	IngressControllerTypeApi    = "api"    //微服务类入口

	IngressRulePath            = "Path"            //绝对命中，且传递全路径
	IngressRulePathPrefix      = "PathPrefix"      //前缀命中，且传递全路径
	IngressRulePathStrip       = "PathStrip"       //绝对命中，传递时丢弃配置路径
	IngressRulePathPrefixStrip = "PathPrefixStrip" //前缀命中，传递时丢弃配置路径
)

//ingress rule list
var IngressRules = []string{
	IngressRulePath,
	IngressRulePathPrefix,
	IngressRulePathStrip,
	IngressRulePathPrefixStrip,
}
