package node

import (
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Del(c *gin.Context) (err error) {
	var obj = models.GetNode(c)
	if obj == nil {
		return errors.New("node not found")
	}
	return models.DB().Delete(&models.Node{}, "id=?", obj.ID).Error
}
