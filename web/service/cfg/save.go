package cfg

import (
	"app/library/ginutil"
	"app/library/types/jsonutil"
	"app/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Save(c *gin.Context) (err error) {
	var val models.CfgStruct
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = val.Validator(); err != nil {
		return
	}
	return models.DB().Transaction(func(tx *gorm.DB) error {
		if er := tx.Unscoped().Delete(&models.Cfg{}, "1=1").Error; er != nil {
			return er
		}
		for k, v := range val.Map() {
			var value string
			switch vv := v.(type) {
			case string:
				value = vv
			default:
				value = jsonutil.ToString(vv)
			}
			if er := tx.Save(&models.Cfg{Name: k, Value: value}).Error; er != nil {
				return er
			}
		}
		return models.ReloadCfg()
	})
}
