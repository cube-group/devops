// Author: chenqionghe
// Time: 2018-10
// 用户相关业务逻辑
package sys

import (
	"app/library/e"
	"app/library/page"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

//退出
func Logout(c *gin.Context)  {
	models.SessionClear(c)
}

type valUser struct {
	Name string `form:"name"`
}

//分页列表
func UserList(c *gin.Context) (res gin.H) {
	var searchParams = gin.H{"name": ""}
	res = gin.H{"search": searchParams}
	var val valUser
	if err := c.ShouldBind(&val); err != nil {
		return
	}
	res = gin.H{"search": gin.H{"name": val.Name}}
	var list []models.User
	result, err := page.List(c, &list, queryUser(&val), searchParams)
	if err != nil {
		return
	}
	return result
}

func queryUser(val *valUser) *gorm.DB {
	var query = models.DB().Order("id DESC")
	if val.Name != "" {
		if id := convert.MustUint(val.Name); id > 0 {
			return query.Where("id=?", id)
		}
		like := "%" + val.Name + "%"
		query = query.Where("username LIKE ?", like)
		query = query.Or("name LIKE ?", like)
		query = query.Or("real_name LIKE ?", like)
		query = query.Or("email LIKE ?", like)
	}
	return query
}

//修改用户
func UserUpdate(c *gin.Context) (user *models.User, err error) {
	var val ValidateUpdate
	if err := c.ShouldBindWith(&val, binding.Form); err != nil {
		return nil, e.Error(err)
	}
	user = new(models.User)
	if models.DB().Where("id = ?", val.ID).Take(user).Error != nil {
		return nil, e.Error("用户不存在！")
	}
	err = models.DB().Save(user).Error
	return
}

//删除
func UserDelete(c *gin.Context) (*models.User, error) {
	var val ValidateDelete
	if err := c.ShouldBindWith(&val, binding.Form); err != nil {
		return nil, e.Error(err)
	}
	var user = new(models.User)
	if models.DB().Where("id = ?", val.ID).Take(user).Error != nil {
		return nil, e.Error("用户不存在！")
	}
	if err := models.DB().Delete(user).Error; err != nil {
		return nil, e.Error("删除失败！", err)
	}
	return user, nil
}