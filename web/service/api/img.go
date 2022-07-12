package api

import (
	"app/library/types/convert"
	"app/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ImgAvatarGet(c *gin.Context) {
	if userID := convert.MustUint32(c.Param("img_id")); userID > 0 {
		user := models.GetUser(userID)
		if user == nil {
			goto Error
		}
		c.Header("Content-Type", "image")
		c.String(http.StatusOK, user.AvatarBlob)
		return
	}

Error:
	c.String(http.StatusOK, "")
}
