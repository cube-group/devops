package project

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) (err error) {
	var val models.Project
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = val.Validator(c); err != nil {
		return
	}
	if val.ID == 0 {
		val.Uid = models.UID(c)
	}
	return models.DB().Save(&val).Error
}
