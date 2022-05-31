// Author: chenqionghe
// Time: 2018-10
// 日志类型，所有日志枚举

package consts

const (
	//群组
	LogProjectGroupCreate = "project-group-create"
	LogProjectGroupUpdate = "project-group-update"
	LogProjectGroupDelete = "project-group-delete"

	//项目
	LogProjectUpdate        = "project-update"
	LogProjectCreate        = "project-create"
	LogProjectPass          = "project-pass"
	LogProjectDeny          = "project-deny"
	LogProjectUpdateScale   = "project-update-batch-size"
	LogProjectOnline        = "project-online"
	LogProjectRollBack      = "project-rollBack"
	LogProjectDelete        = "project-delete"
	LogProjectDeleteService = "project-delete-service"
	LogProjectChangeOwner   = "project-change-owner"
	LogProjectUpdateElastic = "project-update-elastic"
	LogProjectInfo          = "project-info"

	//用户
	LogUserCreate = "user-create"
	LogUserUpdate = "user-update"
	LogUserDelete = "user-delete"
	LogUserLogin  = "user-login"
	LogUserLogout = "user-logout"
	LogUserInfo   = "user-info"

	//角色
	LogRoleCreate = "role-create"
	LogRoleUpdate = "role-update"
	LogRoleDelete = "role-delete"

	//权限
	LogPermissionCreate = "permission-create"
	LogPermissionUpdate = "permission-update"
	LogPermissionDelete = "permission-delete"
	LogPermissionMove   = "permission-move"
	LogPermissionCopy   = "permission-copy"

	//设置
	LogSettingUpdate = "setting-update"

	//CDN
	LogCdnAddFile = "cdn-add-file"
	LogCdnDelFile = "cdn-del-file"

	//模板
	LogDockerfileUpdate = "dockerfile-update"
	LogDockerfileCreate = "dockerfile-create"
	LogDockerfileDelete = "dockerfile-delete"

	//秘钥
	LogPasswordUpdate = "password-update"
	LogPasswordCreate = "password-create"
	LogPasswordDelete = "password-delete"

	//证书
	LogCertificateUpdate = "certificate-update"
	LogCertificateCreate = "certificate-create"
	LogCertificateDelete = "certificate-delete"

	//镜像
	LogImageUpdate     = "image-update"
	LogImageCreate     = "image-create"
	LogImageDelete     = "image-delete"
	LogImageBuild      = "image-build"
	LogImageTagsDelete = "image-tags-delete"

	//配置中心
	LogMsConfigIndex          = "msConfig-index"
	LogMsConfigDoUpdate       = "msConfig-do-update"
	LogMsConfigDoCreate       = "msConfig-do-create"
	LogMsConfigDoDelete       = "msConfig-do-delete"
	LogMsConfigDetail         = "msConfig-detail"
	LogMsConfigDoGrantProject = "msConfig-do-grant-project"
	LogMsConfigDoGrantUser    = "msConfig-do-grant-user"
	LogMsConfigPublish        = "msConfig-publish"
	LogMsConfigPublishList    = "msConfig-publish-list"
	LogMsConfigUseUser        = "msConfig-use-user"
	LogMsConfigUseProject     = "msConfig-use-project"
	LogMsConfigLog            = "msConfig-log"

	//Wiki
	LogWikiUpdate = "wiki-update"
	LogWikiCreate = "wiki-create"
	LogWikiDelete = "wiki-delete"

	//服务器
	LogNodeUpdate = "node-update"
	LogNodeCreate = "node-create"
	LogNodeDelete = "node-delete"
	LogNodeCmd    = "node-cmd"

	//gitlab成员操作
	LogGitMemberUpdate      = "git-member-update"
	LogGitMemberCreate      = "git-member-create"
	LogGitMemberCreateBatch = "git-member-create-batch"
	LogGitMemberDelete      = "git-member-delete"

	//导出
	LogExportCurrentCycle = "export-current-cycle"
	LogExportCycle        = "export-cycle"
	LogExportLastCycle    = "export-last-cycle"

	//容器
	LogContainerLog         = "container-log"
	LogContainerExec        = "container-exec"
	LogProjectContainerList = "container-list"
	LogContainerExecCmd     = "container-exec-cmd"

	//代码统计
	LogCodeStatsIndex = "codeStatsIndex"

	LogCodeStatsLastDayDetail      = "codeStatsLastDayDetail"
	LogCodeStatsCurrentCycleDetail = "codeStatsCurrentCycleDetail"
	LogCodeStatsLastCycleDetail    = "codeStatsLastCycleDetail"

	LogCodeStatsLastDayProjectDetail      = "codeStatsLastDayProjectDetail"
	LogCodeStatsCurrentCycleProjectDetail = "codeStatsCurrentCycleProjectDetail"
	LogCodeStatsLastCycleProjectDetail    = "codeStatsLastCycleProjectDetail"

	LogCodeStatsLastDayProjectBranchDetail      = "codeStatsLastDayProjectBranchDetail"
	LogCodeStatsCurrentCycleProjectBranchDetail = "codeStatsCurrentCycleProjectBranchDetail"
	LogCodeStatsLastCycleProjectBranchDetail    = "codeStatsLastCycleProjectBranchDetail"

	//我的信息
	LogMyIndex = "my-index"

	//首页面板
	LogDashboardIndex = "dashboard-index"

	//历史
	LogHistoryIndex = "history-index"

	//rds
	RdsIndex = "rds-index"

	//ecs
	EcsIndex           = "ecs-index"
	EcsCreateWorkers   = "ecs-create-workers"
	EcsWorkerList      = "ecs-create-worker-list"
	EcsDoCreateWorkers = "ecs-do-create-workers"

	//日志服务
	LogLoghubSearch  = "loghub-index"
	LogLoghubProject = "loghub-project"

	//消息
	LogMessageIndex = "message-index"

	//slb黑名单
	LogSlbBlackAdd = "slb-black-add"
	LogSlbBlackDel = "slb-black-del"
	LogSlbWhiteAdd = "slb-white-add"
	LogSlbWhiteDel = "slb-white-del"

	CdnIndex   = "cdn-index"
	CdnRefresh = "cdn-refresh"

	//项目
	LogProjectMemberList          = "project-member-list"
	LogProjectClone               = "project-clone"
	LogProjectBug                 = "project-bug"
	LogProjectMonitor             = "project-monitor"
	LogProjectPreOnline           = "project-preOnline"
	LogProjectOnlineDetail        = "project-onlineDetail"
	LogProjectOnlineDetailCi      = "project-onlineDetailCi"
	LogProjectOnlineDetailCd      = "project-onlineDetailCd"
	LogProjectOnlineDetailGateway = "project-onlineDetailGateway"

	//资源监控
	LogPrometheusIndex   = "project-prometheus-index"
	LogPrometheusRefresh = "project-prometheus-refresh"
)

//获取日志类型中文
func LogTypeCn(v string) string {
	cnMap := LogTypeCnMap()
	if res, ok := cnMap[v]; ok {
		return res
	}
	return "未知"
}

//获取所有的日志映射类型
func LogTypeCnMap() map[string]string {
	return map[string]string{
		LogProjectGroupCreate: "(项目群组）创建",
		LogProjectGroupUpdate: "(项目群组）更新",
		LogProjectGroupDelete: "(项目群组）移除",

		LogProjectUpdate:        "（项目）修改",
		LogProjectCreate:        "（项目）添加",
		LogProjectPass:          "（项目）审核通过",
		LogProjectDeny:          "（项目）审核拒绝",
		LogProjectOnline:        "（项目）上线",
		LogProjectUpdateScale:   "（项目）修改节点",
		LogProjectDelete:        "（项目）删除",
		LogProjectDeleteService: "（项目）删除服务",
		LogProjectChangeOwner:   "（项目）转移管理员",
		LogProjectUpdateElastic: "（项目）修改弹性扩容",
		LogProjectInfo:          "（项目）查看",

		LogUserCreate: "（用户）添加",
		LogUserUpdate: "（用户）修改",
		LogUserDelete: "（用户）删除",
		LogUserLogin:  "（用户）登录",
		LogUserLogout: "（用户）退出",
		LogUserInfo:   "（用户）查看信息",

		LogRoleCreate: "（角色）添加",
		LogRoleUpdate: "（角色）修改",
		LogRoleDelete: "（角色）删除",

		LogPermissionCreate: "（权限）添加",
		LogPermissionUpdate: "（权限）修改",
		LogPermissionDelete: "（权限）删除",
		LogPermissionMove:   "（权限）移动",
		LogPermissionCopy:   "（权限）复制",

		LogCdnAddFile: "（CDN）上传文件",
		LogCdnDelFile: "（CDN）删除文件",

		LogDockerfileCreate: "（模板）添加",
		LogDockerfileUpdate: "（模板）修改",
		LogDockerfileDelete: "（模板）删除",

		LogMsConfigDoCreate:       "（配置中心）执行新增",
		LogMsConfigDoUpdate:       "（配置中心）执行更新",
		LogMsConfigDoDelete:       "（配置中心）删除",
		LogMsConfigLog:            "（配置中心）修改历史",
		LogMsConfigDoGrantProject: "（配置中心）项目授权",
		LogMsConfigDoGrantUser:    "（配置中心）用户授权",
		LogMsConfigPublishList:    "（配置中心）发布历史",
		LogMsConfigUseUser:        "（配置中心）用户调用历史",
		LogMsConfigUseProject:     "（配置中心）项目调用历史",
		LogMsConfigPublish:        "（配置中心）发布",
		LogMsConfigDetail:         "（配置中心）详情",
		LogMsConfigIndex:          "（配置中心）列表",

		LogImageCreate:     "（镜像）添加",
		LogImageUpdate:     "（镜像）修改",
		LogImageDelete:     "（镜像）删除",
		LogImageBuild:      "（镜像）构建",
		LogImageTagsDelete: "（镜像）删除tag",

		LogGitMemberCreate:      "（git成员）添加",
		LogGitMemberCreateBatch: "（git成员）批量添加",
		LogGitMemberUpdate:      "（git成员）修改",
		LogGitMemberDelete:      "（git成员）删除",
		LogExportCurrentCycle:   "（导出）当前考勤周期",
		LogExportLastCycle:      "（导出）上个考勤周期",
		LogExportCycle:          "（导出）指定考勤周期",

		LogContainerLog:     "（容器）查看日志",
		LogContainerExec:    "（容器）执行命令",
		LogContainerExecCmd: "（容器）执行命令详细",

		LogSettingUpdate: "（系统配置）更新",

		LogNodeCreate: "（服务器）添加",
		LogNodeUpdate: "（服务器）修改",
		LogNodeDelete: "（服务器）删除",
		LogNodeCmd:    "（服务器）执行命令",

		LogCodeStatsIndex: "（代码统计）查看",

		LogCodeStatsLastDayDetail:      "（代码统计）今日统计详情",
		LogCodeStatsCurrentCycleDetail: "（代码统计）当前周期统计详情",
		LogCodeStatsLastCycleDetail:    "（代码统计）上个周期统计详情",

		LogCodeStatsLastDayProjectDetail:      "（代码统计）今日项目统计详情",
		LogCodeStatsCurrentCycleProjectDetail: "（代码统计）当前周期项目统计详情",
		LogCodeStatsLastCycleProjectDetail:    "（代码统计）上个周期项目统计详情",

		LogCodeStatsLastDayProjectBranchDetail:      "（代码统计）今日项目分支统计详情",
		LogCodeStatsCurrentCycleProjectBranchDetail: "（代码统计）当前周期项目分支统计详情",
		LogCodeStatsLastCycleProjectBranchDetail:    "（代码统计）上个周期项目分支统计详情",

		LogMyIndex: "（个人中心）我的信息",

		LogDashboardIndex: "（Dashboard）查看",

		LogHistoryIndex: "（上线历史）查看",

		LogMessageIndex: "（我的消息）查看",

		LogLoghubSearch:  "（日志服务）搜索",
		LogLoghubProject: "（项目）流量监控",

		LogSlbBlackAdd: "（slb）黑名单添加",
		LogSlbBlackDel: "（slb）黑名单删除",
		LogSlbWhiteAdd: "（slb）白名单添加",
		LogSlbWhiteDel: "（slb）白名单删除",

		LogProjectContainerList:       "（项目）节点管理",
		LogProjectMemberList:          "（项目）人员管理",
		LogProjectClone:               "（项目）克隆页面",
		LogProjectBug:                 "（项目）漏洞查看",
		LogProjectMonitor:             "（项目）WEB监控",
		LogProjectPreOnline:           "（项目）准备上线",
		LogProjectOnlineDetail:        "（项目）上线详情",
		LogProjectOnlineDetailCi:      "（项目）上线详情-构建",
		LogProjectOnlineDetailCd:      "（项目）上线详情-部署",
		LogProjectOnlineDetailGateway: "（项目）上线详情-网关服务发现",

		CdnIndex:   "（CDN管理）查看",
		CdnRefresh: "（CDN管理）刷新缓存",

		LogPrometheusIndex: "（项目）资源监控",

		EcsIndex:           "（ECS）列表",
		EcsCreateWorkers:   "（ECS）添加Worker节点",
		EcsDoCreateWorkers: "（ECS）执行添加Worker节点",
		EcsWorkerList:      "（ECS）worker列表",

		RdsIndex: "（RDS）列表",
	}
}
