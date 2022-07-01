package cfg

import (
	"app/library/ginutil"
	"app/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

func TagList(c *gin.Context) []models.Tag {
	return models.TagList()
}

func TagSave(c *gin.Context) (err error) {
	var val models.Tag
	if err = ginutil.ShouldBind(c, &val); err != nil {
		return
	}
	if err = val.Validator(); err != nil {
		return
	}
	val.Uid = models.UID(c)
	return models.DB().Save(&val).Error
}

func TagDel(c *gin.Context) (err error) {
	var tag = models.GetTag(c)
	var count int64
	if models.DB().Model(&models.TagRel{}).Where("tid=?", tag.ID).Count(&count).Error != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该标签有%d个项目在使用，请先解除关联", count)
	}
	return models.DB().Unscoped().Delete(&models.Tag{}, "id=?", tag.ID).Error
}
