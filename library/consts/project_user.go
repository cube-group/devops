// Author: chenqionghe
// Time: 2018-11
// git访问权限等级

package consts

const (
	AccessLevelOwner     = 50 //owner，全部管理权限
	AccessLevelMaster    = 40 //master，除了能删除项目其它都能干
	AccessLevelDeveloper = 30 //developer，可以上线不可修改不可删除不可修改节点
)

func IsAccessLevelOwner(v uint) bool {
	return v == AccessLevelOwner
}

func IsAccessLevelMaster(v uint) bool {
	return v == AccessLevelMaster
}

func IsAccessLevelDeveloper(v uint) bool {
	return v == AccessLevelDeveloper
}

func IsAccessLevelProjectSuperAdmin(v uint) bool {
	return v == AccessLevelOwner
}

func IsAccessLevelProjectAdmin(v uint) bool {
	return v == AccessLevelOwner || v == AccessLevelMaster
}

func IsAccessLevelProjectDevelop(v uint) bool {
	return v >= AccessLevelDeveloper
}

func IsAccessLevelProjectSee(v uint) bool {
	return v > 0
}

func AccessLevelCn(v uint) string {
	cnMap := map[uint]string{
		AccessLevelOwner:     "owner",
		AccessLevelMaster:    "master",
		AccessLevelDeveloper: "developer",
	}
	if res, ok := cnMap[v]; ok {
		return res
	}
	return "未知"
}

func AccessLevelPermission(v uint) string {
	cnMap := map[uint]string{
		AccessLevelOwner:     "可操作该项目的所有权限",
		AccessLevelMaster:    "可操作该项目的所有权限",
		AccessLevelDeveloper: "项目内容只读权限",
	}
	if res, ok := cnMap[v]; ok {
		return res
	}
	return "未知"
}
