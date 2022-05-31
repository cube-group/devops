package cfg

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Save(c *gin.Context) (err error) {
	var val map[string]string
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Scopes().Delete(&models.Cfg{},"1=1").Error; er != nil {
			return er
		}
		for k, v := range val {
			if er := tx.Save(&models.Cfg{Name: k, Value: v}).Error; er != nil {
				return er
			}
		}
		return models.ReloadCfg()
	})
}
