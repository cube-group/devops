package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func History() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var history *models.History
		var id = convert.MustUint32(c.Param("historyId"))
		if id == 0 {
			err = errors.New("id is nil")
			goto Error
		}
		history = models.GetHistory(id)
		if history == nil {
			err = errors.New("history resource not found 1")
			goto Error
		}
		c.Set(consts.ContextHistory, history)
		c.Next()
		return

	Error:
		g.Warning(c, err)
	}
}

func HistoryPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var history = models.GetHistory(c)
		if history == nil {
			g.Warning(c, errors.New("history resource not found 2"))
			return
		}
		if models.GetUser(c).HasPermissionProject(history.ProjectId) != nil {
			g.Warning(c, errors.New("没有权限操作"))
			return
		}
		c.Next()
	}
}
