package consts

//超级管理员
const RoleSuperAdmin = 1

//默认
const RoleDefault = 2

func IsRoleSuperAdmin(id uint) bool {
	return id == RoleSuperAdmin
}

func IsRoleDefault(id uint) bool {
	return id == RoleDefault
}
