package project

import (
	"app/library/ginutil"
	"app/library/page"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type valList struct {
	Name string `form:"name"`
	Kind string `form:"kind"`
	Env  string `form:"env"`
}

type ProjectLatestHistory struct {
	ProjectId uint32 `gorm:""`
	HistoryId uint32 `gorm:""`
}

type ProjectLatestHistoryList []ProjectLatestHistory

func (t ProjectLatestHistoryList) HistoryList() []models.History {
	var list = make([]uint32, 0)
	for _, v := range t {
		list = append(list, v.HistoryId)
	}
	var res []models.History
	models.DB().Find(&res, "id IN (?)", list)
	return res
}

//k8s project list
func List(c *gin.Context) (res gin.H) {
	var obj = page.ListReturnStruct{"search": gin.H{"name": "", "env": "", "kind": ""}}
	res = gin.H(obj)
	var val valList
	if ginutil.ShouldBind(c, &val) != nil {
		return
	}
	var list models.ProjectList
	result, err := page.List(c, &list, queryList(val), obj)
	if err != nil {
		return
	}

	var historyList ProjectLatestHistoryList
	var query = models.DB().Model(&models.History{}).Select("max(id) AS history_id,project_id")
	query = query.Group("project_id").Where("project_id IN (?)", list.IDs())
	if err := query.Scan(&historyList).Error; err == nil {
		result["historyList"] = historyList.HistoryList()
	}
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
		if id := convert.MustUint32(val.Name); id > 0 {
			query = query.Where("id=?", id)
		} else {
			query = query.Where("name LIKE ?", "%"+val.Name+"%")
		}
	}
	return query
}
