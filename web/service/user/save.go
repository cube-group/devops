package user

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) (err error) {
	var val models.User
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = val.Validator(); err != nil {
		return err
	}
	return models.DB().Save(&val).Error
}
