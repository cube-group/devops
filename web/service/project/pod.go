package project

import (
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
	return history.Remove()
}
