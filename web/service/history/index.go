package history

import (
	"app/library/ginutil"
	"app/library/page"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type valList struct {
	ID  string `form:"id"`
	Pid string `form:"pid"`
}

//k8s node list
func List(c *gin.Context) (res gin.H) {
	var obj = page.ListReturnStruct{"search": gin.H{"id": "", "pid": ""}}
	res = gin.H(obj)
	var val valList
	if ginutil.ShouldBind(c, &val) != nil {
		return
	}
	var list []models.History
	result, _ := page.List(c, &list, queryList(val), obj)
	return result
}

func queryList(val valList) *gorm.DB {
	var query = models.DB().Order("id DESC")
	if val.ID != "" {
		if id := convert.MustUint32(val.ID); id > 0 {
			query = query.Where("id=?", id)
		}
	}
	if val.Pid != "" {
		if projectId := convert.MustUint32(val.Pid); projectId > 0 {
			if p := models.GetProject(projectId); p != nil {
				query = query.Where("project_id=?", projectId)
			}
		}
	}
	return query
}
