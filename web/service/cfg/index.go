package cfg

import (
	"app/models"
	"github.com/gin-gonic/gin"
)

//k8s node list
func List(c *gin.Context) (res gin.H, err error) {
	res = gin.H{}
	var list []models.Cfg
	if err = models.DB().Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		res[v.Name] = v.Value
	}
	return
}
