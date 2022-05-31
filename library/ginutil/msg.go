// Author: chenqionghe
// Time: 2018-10
// 单例msg

package ginutil

import (
	"github.com/gin-gonic/gin"
)

const (
	MSG_KEY  = "app-msg"
	CODE_KEY = "app-code"
)

//设置消息
func Msg(c *gin.Context, msg ...string) string {
	if msg == nil {
		return c.GetString(MSG_KEY)
	}
	c.Set(MSG_KEY, msg[0])
	return msg[0]
}
