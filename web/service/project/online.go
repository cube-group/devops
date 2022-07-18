package project

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
)

//cluster project online
func Online(c *gin.Context) (history models.History, err error) {
	var val models.History
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	project := models.GetProject(c)
	val.ProjectId = project.ID
	val.Project = project
	val.ID = 0
	val.Uid = models.UID(c)
	val.Status = models.HistoryStatusDefault
	if err = val.Validator(); err != nil {
		return
	}

	if project.Cronjob == "" { //normal online
		if err = project.Apply(&val, true); err != nil {
			return
		}
		history = val
	} else { //cronjob online
		err = models.CronjobAdd(models.ProjectCronjob{
			Uid:       val.Uid,
			Nodes:     val.Nodes,
			ProjectId: project.ID,
		})
	}
	return
}

//project cronjob stop
func Offline(c *gin.Context) error {
	return models.CronjobStop(models.GetProject(c).ID)
}
