package project

import (
	"app/models"
	"errors"
	"github.com/gin-gonic/gin"
)

func Del(c *gin.Context) (err error) {
	var obj = models.GetProject(c)
	if obj == nil {
		return errors.New("project not found")
	}
	return models.DB().Delete(&models.Project{}, "id=?", obj.ID).Error
}
