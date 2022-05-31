// Author: chenqionghe
// Time: 2019-07

package consts

//项目配置修改日志类型
type MsConfigLogType uint

const (
	//项目配置修改日志类型，追加
	MsConfigLogTypeAdd MsConfigLogType = 1
	//项目配置修改日志类型，更新
	MsConfigLogTypeUpdate MsConfigLogType = 2

	MsTokenUser    = 1
	MsTokenProject = 2
)
