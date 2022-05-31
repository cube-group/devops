// Author: chenqionghe
// Time: 2018-10
// 验证结构体
package sys

type ValidateLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type ValidateUpdate struct {
	ID      uint   `form:"id" binding:"required"`
	RoleIds []uint `form:"roleIds[]" binding:"required"`
}

type ValidateDelete struct {
	ID uint `form:"id" binding:"required"`
}