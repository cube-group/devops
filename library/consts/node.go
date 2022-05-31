// Author: chenqionghe
// Time: 2019-10
// 项目相关枚举

package consts

const (
	//待创建
	NodeStatusWait = 101

	//创建中
	NodeStatusStartCreateEcs          = 211
	NodeStatusCreateEcs               = 201
	NodeStatusStartStartEcs           = 212
	NodeStatusStartEcs                = 203
	NodeStatusStartSshd               = 204
	NodeStatusInstallLogtail          = 205
	NodeStatusAddToMachineGroup       = 206
	NodeStatusStartInitAsRancherNode  = 209
	NodeStatusInitAsRancherNode       = 207
	NodeStatusStartRenameAndAddLabels = 210
	NodeStatusRenameAndAddLabels      = 208
	//创建结束
	NodeStatusSuccess = 300
	NodeStatusFail    = 301
)

func NodeStatusCn(v uint) string {
	cnMap := map[uint]string{
		NodeStatusWait: "待创建",

		NodeStatusStartCreateEcs: "创建ECS中...",
		NodeStatusCreateEcs:      "创建ECS成功",

		NodeStatusStartStartEcs: "启动ECS中...",
		NodeStatusStartEcs:      "启动ECS成功",

		NodeStatusStartSshd: "SSH连接成功",

		NodeStatusInstallLogtail:    "安装Logtail成功",
		NodeStatusAddToMachineGroup: "加入日志机器组成功",

		NodeStatusStartInitAsRancherNode: "初始化为Rancher节点...",
		NodeStatusInitAsRancherNode:      "初始化为Rancher节点成功",

		NodeStatusStartRenameAndAddLabels: "重命名并打标签...",
		NodeStatusRenameAndAddLabels:      "重命名并打标签成功",
		NodeStatusSuccess:                 "worker初始化完成",
		NodeStatusFail:                    "worker初始化失败",
	}
	if res, ok := cnMap[v]; ok {
		return res
	}
	return "未知状态"
}

func IsNodeStatusWait(status uint) bool {
	return status == NodeStatusWait
}

func IsNodeStatusCreating(status uint) bool {
	return status >= 200 && status < 300
}

//是否是部署失败
func IsNodeStatusFail(status uint) bool {
	return status == NodeStatusFail
}

//是否已经停止
func IsNodeStatusSuccess(status uint) bool {
	return status == NodeStatusSuccess
}
