package node

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

type valInfos struct {
	Ids []uint32 `json:"ids" form:"ids" binding:"required"`
}

func getNodeRemoteShellResult(c *gin.Context, shell string) (res gin.H, err error) {
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
			bytes, er := i.Exec(shell)
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

func GetDockerVersion(c *gin.Context) (res gin.H, err error) {
	return getNodeRemoteShellResult(c, "docker version --format='{{json .}}'")
}

func GetDockerStats(c *gin.Context) (res gin.H, err error) {
	res, err = getNodeRemoteShellResult(c, "docker stats --no-stream --format='{{json .}}'")
	if err != nil {
		return
	}
	for k, v := range res {
		var item = v.(gin.H)
		var content = convert.MustString(item["content"])
		var jsonItemList = []gin.H{}
		for _, i := range strings.Split(content, "\n") {
			var jsonItem gin.H
			if er := jsoniter.Unmarshal([]byte(i), &jsonItem); er == nil {
				jsonItemList = append(jsonItemList, jsonItem)
			}
		}
		res[k] = jsonItemList
	}
	return
}
