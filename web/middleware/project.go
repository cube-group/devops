package middleware

import (
	"app/library/consts"
	"app/library/g"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Project() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var project *models.Project
		id := convert.MustUint32(c.Param("pid"))
		if id == 0 {
			err = errors.New("pid is nil")
			goto Error
		}
		project = models.GetProject(id)
		if project == nil {
			err = errors.New("project resource not found")
			goto Error
		}
		c.Set(consts.ContextProject, project)
		c.Next()
		return

	Error:
		g.Warning(c, err)
	}
}

func ProjectPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var project = models.GetProject(c)
		if project == nil {
			err = errors.New("project resource not found")
			goto Error
		}
		if models.GetUser(c).HasPermissionProject(project.ID) != nil {
			err = errors.New("没有权限操作")
			goto Error
		}
		c.Next()
		return

	Error:
		g.Warning(c, err)
	}
}
