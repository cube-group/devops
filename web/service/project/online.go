package project

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//cluster project online
func Online(c *gin.Context) (history *models.History, err error) {
	var val models.History
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}

	var project = models.GetProject(c)
	val.ProjectId = project.ID
	val.Project = project
	val.ID = 0
	val.Uid = models.UID(c)
	val.Status = models.HistoryStatusDefault

	if err = val.Validator(); err != nil {
		return
	}
	if err = models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Save(&val).Error; er != nil {
			return er
		}
		return val.Online()
	}); err == nil {
		history = &val
	}
	return
}
