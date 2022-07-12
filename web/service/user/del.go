package user

import (
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Del(c *gin.Context) (err error) {
	var uid = c.Param("user_id")
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Delete(&models.User{}, "id=?", uid).Error; er != nil {
			return er
		}
		if er := tx.Unscoped().Delete(&models.ProjectUser{}, "uid=?", uid).Error; er != nil {
			return er
		}
		return nil
	})
}
