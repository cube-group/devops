package project

import (
	"app/library/ginutil"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//cluster project online
func Online(c *gin.Context) (history *models.History, err error) {
	var val models.History
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	project := models.GetProject(val.ProjectId)
	if project == nil {
		err = errors.New("project not found")
		return
	}
	val.Uid = models.UID(c)
	val.Project = project
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
