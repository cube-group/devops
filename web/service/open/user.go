package open

import (
	"app/models"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func UserList(c *gin.Context) (res map[uint32]interface{}) {
	res = make(map[uint32]interface{})
	var ids []uint
	bytes, _ := c.GetRawData()
	if jsoniter.Unmarshal(bytes, &ids) != nil {
		ids = make([]uint, 0)
	}
	if len(ids) == 0 {
		return
	}
	var list []models.User
	if models.DB().Find(&list, "id IN (?)", ids).Error != nil {
		return
	}
	for _, v := range list {
		res[v.ID] = v
	}
	return

}
