package project

import (
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Del(c *gin.Context) (err error) {
	var obj = models.GetProject(c)
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Delete(&models.Project{}, "id=?", obj.ID).Error; er != nil {
			return er
		}
		if h := obj.GetLatestHistory(); h != nil {
			if er := h.Remove(true, tx); er != nil {
				return er
			}
		}
		return nil
	})
}
