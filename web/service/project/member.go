package project

import (
	"app/models"
	"github.com/gin-gonic/gin"
)

//该项目成员列表
func MemberList(c *gin.Context) (res []models.ProjectUser) {
	models.DB().Order("id DESC").Find(&res, "pid=?", models.GetProject(c).ID)
	return
}

type valMemberSave struct {
	Uid         uint32                    `json:"uid" form:"uid" binding:"required"`
	AccessLevel models.ProjectAccessLevel `json:"accessLevel" form:"accessLevel"`
}

//项目成员变更
func MemberSave(c *gin.Context) error {
	var val valMemberSave
	if err := c.ShouldBind(&val); err != nil {
		return err
	}
	var project = models.GetProject(c)
	var query = models.DB().Where("uid=? AND pid=?", val.Uid, project.ID)
	if val.AccessLevel == models.ProjectAccessLevelNone {
		return query.Delete(&models.ProjectUser{}).Error
	}
	var pu models.ProjectUser
	if err := query.Last(&pu).Error; err != nil {
		return models.DB().Save(&models.ProjectUser{Uid: val.Uid, Pid: project.ID}).Error
	}
	return nil
}
