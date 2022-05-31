//// Author: chenqionghe
//// Time: 2018-10
//// 项目相关枚举
//
package consts
//
//const (
//	ProjectStatusAccess = 100 //待上线
//
//	ProjectStatusBuildOk       = 200 //构建成功
//	ProjectStatusBuilding      = 201 //构建中
//	ProjectStatusBuildFail     = 202 //构建失败
//	ProjectStatusBuildTimeout  = 203 //构建超时
//	ProjectStatusCreateJobFail = 204 //jenkins创建任务失败
//	ProjectStatusBuildCancel   = 205 //构建取消
//
//	ProjectStatusDeploying                    = 301 //部署中
//	ProjectStatusDeployFail                   = 302 //部署失败
//	ProjectStatusDeployTimeout                = 303 //部署超时
//	ProjectStatusDeployFailAndDeleteOK        = 310 //部署失败且删除pod成功
//	ProjectStatusDeployFailAndDeleteFailed    = 311 //部署失败且删除pod失败
//	ProjectStatusDeployFailAndRollBackOK      = 320 //部署失败且回滚成功
//	ProjectStatusDeployFailAndRollBackFailed  = 321 //部署失败且回滚失败
//	ProjectStatusDeployFailAndRollBackTimeout = 322 //部署失败且回滚超时
//	ProjectStatusDeployCancel                 = 323 //部署取消
//	ProjectStatusDeployCreateServiceFail      = 324 //部署失败(创建service失败)
//
//	ProjectStatusOnlineOk = 300 //上线完成
//	ProjectStatusDeployOk = 327 //部署成功
//
//	ProjectStatusDeleteService = 401 //服务已删除
//
//	ProjectUseWeb     = "web"
//	ProjectUseInner   = "inner"
//	ProjectUseJob     = "job"
//	ProjectUseCronJob = "cronjob"
//	ProjectUseBuild   = "build"
//
//	LimitFlowClientIp   = 1 //限制ip
//	LimitFlowHostHeader = 2 //限制
//
//	ScaleDefaultNum = 1   //默认节点数量
//)
//
//func IsProjectUseWebOrInner(v string) bool {
//	return v == ProjectUseWeb || v == ProjectUseInner
//}
//
//func IsProjectUseJobOrCron(v string) bool {
//	return v == ProjectUseJob || v == ProjectUseCronJob
//}
//
//func IsProjectUseJob(v string) bool {
//	return v == ProjectUseJob
//}
//
//func IsProjectUseCronjob(v string) bool {
//	return v == ProjectUseCronJob
//}
//
//func ProjectStatusCn(v uint) string {
//	cnMap := map[uint]string{
//		ProjectStatusAccess: "待上线",
//
//		ProjectStatusBuilding:      "构建中",
//		ProjectStatusBuildFail:     "构建失败",
//		ProjectStatusBuildOk:       "构建成功",
//		ProjectStatusBuildTimeout:  "构建超时",
//		ProjectStatusCreateJobFail: "创建任务失败",
//		ProjectStatusBuildCancel:   "构建取消",
//
//		ProjectStatusDeploying:                    "部署中",
//		ProjectStatusDeployOk:                     "部署成功",
//		ProjectStatusDeployFail:                   "部署失败",
//		ProjectStatusDeployTimeout:                "部署超时",
//		ProjectStatusDeployFailAndRollBackOK:      "部署失败且回滚成功",
//		ProjectStatusDeployFailAndRollBackFailed:  "部署失败且回滚失败",
//		ProjectStatusDeployFailAndRollBackTimeout: "部署失败且回滚超时",
//		ProjectStatusDeployFailAndDeleteOK:        "部署失败且删除pod成功",
//		ProjectStatusDeployFailAndDeleteFailed:    "部署失败且删除pod失败",
//		ProjectStatusDeployCancel:                 "部署取消",
//
//		ProjectStatusDeleteService: "服务删除",
//
//		ProjectStatusDeployCreateServiceFail: "部署失败(创建service失败)",
//
//		ProjectStatusOnlineOk: "上线成功",
//	}
//	if res, ok := cnMap[v]; ok {
//		return res
//	}
//	return "未知状态"
//}
//
////用途
//func LimitFlowTypes() map[int]string {
//	return map[int]string{
//		LimitFlowClientIp:   "client ip",
//		LimitFlowHostHeader: "host header",
//	}
//}
//
///*************************************************** 前端展示使用 ***************************************************/
//
//func IsStatusInBuild(status uint) bool {
//	return status >= 200 && status < 400
//}
//
//func IsStatusInDeploy(status uint) bool {
//	return status >= 300 && status < 400
//}
//
////是否已经停止
//func IsStatusDone(status uint) bool {
//	return IsStatusDeleted(status) || IsStatusBuildFail(status) || IsStatusDeployDone(status)
//}
//
///*************************************************** 构建 ***************************************************/
////是否构建完毕停止
//func IsStatusBuildDone(status uint) bool {
//	return IsStatusBuildOk(status) || IsStatusBuildFail(status)
//}
//
////是否是构建中
//func IsStatusBuilding(status uint) bool {
//	return status == ProjectStatusBuilding
//}
//
////是否是构建失败
//func IsStatusBuildFail(status uint) bool {
//	statusList := []uint{
//		ProjectStatusBuildFail,
//		ProjectStatusBuildTimeout,
//		ProjectStatusCreateJobFail,
//		ProjectStatusBuildCancel,
//	}
//	for _, v := range statusList {
//		if status == v {
//			return true
//		}
//	}
//	return false
//}
//
////是否是构建成功
//func IsStatusBuildOk(status uint) bool {
//	return status == ProjectStatusBuildOk || status >= 300
//}
//
///*************************************************** 部署 ***************************************************/
//func IsStatusDeploying(status uint) bool {
//	return status == ProjectStatusDeploying
//}
//
////是否是部署失败
//func IsStatusDeployFail(status uint) bool {
//	statusList := []uint{
//		ProjectStatusDeployFail,
//		ProjectStatusDeployTimeout,
//		ProjectStatusDeployFailAndDeleteOK,
//		ProjectStatusDeployFailAndDeleteFailed,
//		ProjectStatusDeployFailAndRollBackOK,
//		ProjectStatusDeployFailAndRollBackFailed,
//		ProjectStatusDeployFailAndRollBackTimeout,
//		ProjectStatusDeployCreateServiceFail,
//	}
//	for _, v := range statusList {
//		if status == v {
//			return true
//		}
//	}
//	return false
//}
//
////是否是部署成功
//func IsStatusDeployOk(status uint) bool {
//	return status == ProjectStatusDeployOk ||
//		status == ProjectStatusOnlineOk
//}
//
////是否是部署成功
//func IsStatusDeployDone(status uint) bool {
//	return IsStatusDeployOk(status) || IsStatusDeployFail(status)
//}
//
///*************************************************** 首页统一状态 ***************************************************/
////前端展示-loading中
//func IsStatusOnline(status uint) bool {
//	return IsOnling(status)
//}
//
////前端展示-danger
//func IsStatusWarning(v uint) bool {
//	return IsOnlineFail(v)
//}
//
////前端展示-success
//func IsStatusSuccess(status uint) bool {
//	return IsOnlineOK(status)
//}
//
////前端展示-default
//func IsStatusDefault(status uint) bool {
//	return !IsOnlineOK(status) && !IsOnlineFail(status) && !IsOnling(status)
//}
//
////前端展示-deleted
//func IsStatusDeleted(status uint) bool {
//	return status == ProjectStatusDeleteService
//}
//
////上线失败
//func IsOnlineFail(v uint) bool {
//	return IsStatusDeployFail(v) || IsStatusBuildFail(v)
//}
//
////上线中
//func IsOnling(v uint) bool {
//	return IsStatusBuilding(v) || IsStatusDeploying(v)
//}
//
////上线成功
//func IsOnlineOK(v uint) bool {
//	return v == ProjectStatusOnlineOk
//}
//
////状态，改为人可读的中文
//func ProjectOnlineCn(v uint) string {
//	if IsOnling(v) {
//		return "上线中"
//	}
//	if IsOnlineFail(v) {
//		return "上线失败"
//	}
//	if IsOnlineOK(v) {
//		return "上线成功"
//	}
//
//	return ProjectStatusCn(v)
//}