package history

import (
	"app/library/ginutil"
	"app/models"
	"github.com/gin-gonic/gin"
)

type valState struct {
	List []uint32 `form:"list" binding:"required"`
}

func State(c *gin.Context) (res []models.History, err error) {
	var val valState
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = models.DB().Find(&res, "id IN (?)", val.List).Error; err != nil {
		return
	}
	return
}
