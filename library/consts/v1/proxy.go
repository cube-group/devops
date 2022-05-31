package v1

type DnsType uint

const (
	DnsNamespace = "dns"
	//svc from endpoints subsets ip
	DnsTypeIP DnsType = 0
	//svc cname from externalName
	DnsTypeCName DnsType = 1
	//svc cname from svc
	DnsTypeSvc DnsType = 2
)
