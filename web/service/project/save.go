package project

import (
	"app/library/ginutil"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) (err error) {
	var val models.Project
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = val.Validator(); err != nil {
		return
	}
	if val.ID == 0 {
		val.Uid = models.UID(c)
	} else {
		if models.GetUser(c).HasPermissionProject(val.ID) != nil {
			return errors.New("没有权限操作")
		}
	}
	return models.DB().Save(&val).Error
}
