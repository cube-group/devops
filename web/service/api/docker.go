package api

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/library/types/jsonutil"
	"app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
)

type valInspect struct {
	List []valInspectItem `json:"list"`
}

type valInspectItem struct {
	NodeID        uint32 `json:"nodeId"`
	ContainerName string `json:"containerName"`
}

func DockerInspect(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valInspect
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var wg sync.WaitGroup
	var resMaps sync.Map
	for _, item := range val.List {
		wg.Add(1)
		go func(i valInspectItem) {
			defer wg.Done()
			var node = models.GetNode(i.NodeID)
			if node == nil {
				return
			}
			if bytes, er := node.Exec(fmt.Sprintf("docker inspect %s --format='{{json .}}'", i.ContainerName)); er == nil {
				var instance = jsonutil.ToJson(bytes)
				if instance == nil {
					return
				}
				instance["nodeId"] = i.NodeID
				resMaps.Store(fmt.Sprintf("%d.%s", i.NodeID, i.ContainerName), instance)
			}
		}(item)
	}
	wg.Wait()

	resMaps.Range(func(key, value interface{}) bool {
		res[convert.MustString(key)] = value
		return true
	})
	return
}
