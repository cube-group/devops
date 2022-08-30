package user

import (
	"app/library/ginutil"
	"app/models"
	"errors"
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
	val.From = "sys." + models.SessionUsername(c)
	if val.ID > 0 {
		if user := models.GetUser(val.ID); user != nil {
			val.AvatarBlob = user.AvatarBlob
		} else {
			return errors.New("用户不存在¬")
		}
	}
	return models.DB().Save(&val).Error
}
