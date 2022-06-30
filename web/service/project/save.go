package project

import (
	"app/library/ginutil"
	"app/models"
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
	}
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := models.DB().Save(&val).Error; er != nil {
			return er
		}
		return models.TagRelProject(tx, val.ID, val.Tag)
	})
}
