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
		//需要检测之前的history如果存在project.name不一致需要先移除container
		if val.Mode == models.ProjectModeDocker {
			if latestHistory := val.GetLatestHistory(tx); latestHistory != nil {
				if !latestHistory.CanChangeName(val.Name) {
					return errors.New("项目更名:之前服务节点尚在，请先清理节点再执行操作")
				}
			}
		}
		if er := models.DB().Save(&val).Error; er != nil {
			return er
		}
		return models.TagRelProject(tx, val.ID, val.Tag)
	})
}
