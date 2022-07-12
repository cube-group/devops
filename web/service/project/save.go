package project

import (
	"app/library/ginutil"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := models.DB().Save(&val).Error; er != nil {
			return er
		}
		return models.TagRelProject(tx, val.ID, val.Tag)
	})
}
