package node

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"sync"
)

type valInfos struct {
	Ids []uint32 `json:"ids" form:"ids" binding:"required"`
}

func GetState(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valInfos
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var list []models.Node
	if err = models.DB().Find(&list, "id IN (?)", val.Ids).Error; err != nil {
		return
	}
	//并发
	var maps sync.Map
	var wg sync.WaitGroup
	for _, v := range list {
		wg.Add(1)
		go func(i models.Node) {
			defer wg.Done()
			bytes, er := i.Exec("docker version")
			var item = gin.H{"id": i.ID, "content": string(bytes), "error": ""}
			if er != nil {
				item["error"] = er.Error()
			}
			maps.Store(i.ID, item)
		}(v)
	}
	wg.Wait()
	maps.Range(func(key, value interface{}) bool {
		res[convert.MustString(key)] = value
		return true
	})
	return
}
