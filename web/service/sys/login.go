package sys

import (
	"app/library/ginutil"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

type valLogin struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//username login by password
func Login(c *gin.Context) error {
	var val valLogin
	if err := ginutil.ShouldBind(c, &val); err != nil {
		return errors.New("参数错误 " + err.Error())
	}
	var user = new(models.User)
	if user.PwdCheckUser(val.Username, val.Password) == nil {
		return errors.New("用户或密码错误")
	}
	return user.LoginCookie(c)
}
