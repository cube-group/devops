package project

import (
	"app/library/ginutil"
	"app/library/page"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type valList struct {
	Name string `form:"name"`
	Kind string `form:"kind"`
	Env  string `form:"env"`
}

//k8s project list
func List(c *gin.Context) (res gin.H) {
	var obj = page.ListReturnStruct{"search": gin.H{"name": "", "env": "", "kind": ""}}
	res = gin.H(obj)
	var val valList
	if ginutil.ShouldBind(c, &val) != nil {
		return
	}
	var list []models.Project
	result, _ := page.List(c, &list, queryList(val), obj)
	return result
}

func queryList(val valList) *gorm.DB {
	var query = models.DB().Order("id DESC")
	if val.Env != "" {
		query = query.Where("env=?", val.Env)
	}
	if val.Kind != "" {
		query = query.Where("kind=?", val.Kind)
	}
	if val.Name != "" {
		query = query.Where("name LIKE ?", "%"+val.Name+"%")
	}
	return query
}
