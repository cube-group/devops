package user

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
}

func List(c *gin.Context) (res gin.H) {
	var obj = page.ListReturnStruct{"search": gin.H{"name": ""}}
	res = gin.H(obj)
	var val valList
	if ginutil.ShouldBind(c, &val) != nil {
		return
	}
	var list []models.User
	result, err := page.List(c, &list, queryList(val), obj)
	if err != nil {
		return
	}

	//var historyList ProjectLatestHistoryList
	//var query = models.DB().Model(&models.History{}).Select("max(id) AS history_id,project_id")
	//query = query.Group("project_id").Where("project_id IN (?)", list.IDs())
	//if err := query.Scan(&historyList).Error; err == nil {
	//	result["historyList"] = historyList.HistoryList()
	//}
	//var cronjobList []models.ProjectCronjob
	//if models.DB().Find(&cronjobList, "project_id IN (?)", list.IDs()).Error == nil {
	//	result["cronjobList"] = cronjobList
	//}
	return result
}

func queryList(val valList) *gorm.DB {
	var query = models.DB().Order("id DESC")
	if val.Name != "" {
		if id := convert.MustUint32(val.Name); id > 0 {
			query = query.Where("id=?", id)
		} else {
			var like = "%" + val.Name + "%"
			query = query.Where("username LIKE ? OR real_name LIKE ?", like, like)
		}
	}
	return query
}
