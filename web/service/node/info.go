package node

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"sync"
)

type valInfos struct {
	Ids []uint32 `json:"ids" form:"ids" binding:"required"`
}

func DockerVersion(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var val valInfos
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	var list []models.Node
	if err = models.DB().Find(&list, "id IN (?)", val.Ids).Error; err != nil {
		return
	}
	var maps sync.Map
	var wg sync.WaitGroup
	for _, v := range list {
		wg.Add(1)
		go func(i models.Node) {
			defer wg.Done()
			bytes, er := i.GetDockerVersion()
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

func DockerPs(c *gin.Context) (res []models.NodeContainerPsItem, err error) {
	models.GetNode(c).GetContainerRandomPort()
	return models.GetNode(c).GetDockerContainerList()
}

func DockerRestart(c *gin.Context) (err error) {
	var name = ginutil.Input(c, "name")
	if name == "" {
		return errors.New("name is nil")
	}
	var node = models.GetNode(c)
	_, err = node.Exec("docker restart " + name)
	return
}

func DockerRm(c *gin.Context) (err error) {
	var name = ginutil.Input(c, "name")
	if name == "" {
		return errors.New("name is nil")
	}
	var node = models.GetNode(c)
	_, err = node.Exec("docker rm -f " + name)
	return
}

func DockerStats(c *gin.Context) (res []interface{}, err error) {
	var node = models.GetNode(c)
	bytes, err := node.Exec("docker stats --no-stream --format='{{json .}}'")
	if err != nil {
		return
	}
	res = make([]interface{}, 0)
	for _, strItem := range strings.Split(string(bytes), "\n") {
		var resItem gin.H
		if jsoniter.Unmarshal([]byte(strItem), &resItem) == nil {
			res = append(res, resItem)
		}
	}
	return
}
