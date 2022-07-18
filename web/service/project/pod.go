package project

import (
	"app/library/ginutil"
	"app/library/types/convert"
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func PodDel(c *gin.Context) error {
	project := models.GetProject(c)
	history := project.GetLatestHistory()
	if history == nil {
		return errors.New("未找到近期部署信息")
	}
	var nodeID = convert.MustUint32(ginutil.Input(c, "nid"))
	if !history.Nodes.Has(nodeID) {
		return errors.New("node not found 1")
	}
	node := models.GetNode(nodeID)
	if node == nil {
		return errors.New("node not found 2")
	}
	return history.Remove(false, node)
}
