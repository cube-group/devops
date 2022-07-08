package user

import (
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Del(c *gin.Context) (err error) {
	var user = models.GetUser(c)
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Delete(&models.User{}, "id=?", user.ID).Error; er != nil {
			return er
		}
		if er := tx.Unscoped().Delete(&models.ProjectUser{}, "uid=?", user.ID).Error; er != nil {
			return er
		}
		return nil
	})
}
