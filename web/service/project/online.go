package project

import (
	"app/library/ginutil"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

//cluster project online
func Online(c *gin.Context) (history models.History, err error) {
	var val models.History
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var project = models.GetProject(c)
	val.ID = 0
	val.ProjectId = project.ID
	val.Uid = models.UID(c)
	val.Status = models.HistoryStatusDefault
	//rollback logic
	if val.Rollback > 0 {
		rollbackHistory := models.GetHistory(val.Rollback)
		if rollbackHistory == nil {
			err = errors.New("回复命中回滚镜像历史")
			return
		}
		val.Project = rollbackHistory.Project
		val.Project.Docker.Dockerfile = ""
		val.Project.Docker.Image = rollbackHistory.RollbackImageURL()
	} else {
		val.Project = project
	}
	if err = val.Validator(); err != nil {
		return
	}
	//apply operation
	if project.IsCronjob() {
		//cronjob start
		err = models.CronjobAdd(models.ProjectCronjob{
			Uid:       val.Uid,
			Nodes:     val.Nodes,
			ProjectId: project.ID,
		})
	} else {
		//online
		if err = project.Apply(&val, true); err != nil {
			return
		}
		history = val
	}
	return
}

//project cronjob stop
func Offline(c *gin.Context) error {
	return models.CronjobStop(models.GetProject(c).ID)
}
