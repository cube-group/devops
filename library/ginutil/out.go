// Author: chenqionghe
// Time: 2018-10
// 标准输出

package ginutil

import (
	"app/library/types/convert"
	"app/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JsonSuccess(c *gin.Context, msg interface{}, data interface{}) {
	c.Header("x-server", setting.Local.GetString("name"))
	c.Header("x-server-version", setting.Local.GetString("version"))
	if gMsg := Msg(c); gMsg != "" {
		msg = gMsg
	}
	c.JSON(http.StatusOK, JsonData(CODE_OK, data, msg))
	c.Abort()
}

/**
 * 输出错误json
 * @param c gin请求实例
 * @param msg 错误信息
 * @param data 返回数据
 * @code 返回code码，如果为空，返回3000
 */
func JsonError(c *gin.Context, msg interface{}, data interface{}, code ...int) {
	c.Header("x-server", setting.Local.GetString("name"))
	c.Header("x-server-version", setting.Local.GetString("version"))
	resCode := CODE_ERR
	codeContext, exist := c.Get(CODE_KEY)
	if exist {
		resCode = convert.MustInt(codeContext)
	}
	if code != nil {
		resCode = code[0]
	}
	if gMsg := Msg(c); gMsg != "" {
		msg = gMsg
	}
	c.JSON(http.StatusOK, JsonData(resCode, data, msg))
	c.Abort()
}

//如果err不为nil,输入JsonError，否则输出JsonSuccess,msg为successMsg
func JsonAuto(c *gin.Context, successMsg string, err error, data interface{}, code ...int) {
	if err != nil {
		JsonError(c, err, data, code...)
	} else {
		JsonSuccess(c, successMsg, data)
	}
}

//获取json数据结构
func JsonData(code int, data interface{}, msg interface{}) map[string]interface{} {
	return gin.H{
		"data": data,
		"code": code,
		"msg":  fmt.Sprintf("%v", msg),
	}
}

//显示HTML，加上头部公共信息(如登录用户)
func HTML(c *gin.Context, template string, data map[string]interface{}, code ...int) {
	c.Header("x-server", setting.Local.GetString("name"))
	c.Header("x-server-version", setting.Local.GetString("version"))
	if data == nil {
		data = gin.H{}
	}
	//data["_u"] = models.SessionUser(c)
	//data["appConf"], _ = config.GetAll(c)
	data["randomNum"] = time.Now().Unix()
	//设置错误消息，如果appErr存在
	setAppErr(data)
	//设置成功消息，如果appOk
	setAppOk(data)
	//设置http状态码，如果没传状态码，默认用200
	httpCode := getHttpCode(code...)

	//如果appErr存在，直接显示err模板，屏蔽不必要的错误
	if data["appErr"] != "" {
		data["title"] = "Render Global Err"
		c.HTML(httpCode, "errors/exception.html", data)
		return
	}

	c.HTML(httpCode, template, data)
	c.Abort()
}

func getHttpCode(code ...int) int {
	httpCode := http.StatusOK
	if code != nil {
		httpCode = code[0]
	}
	return httpCode
}

func setAppErr(data gin.H) {
	if err, ok := data["appErr"]; ok {
		if err != nil {
			data["appErr"] = fmt.Sprint(err)
			return
		}
	}
	data["appErr"] = ""
	return

}

func setAppOk(data gin.H) {
	if err, ok := data["appOk"]; ok {
		if err != nil {
			data["appOk"] = fmt.Sprint(err)
		} else {
			data["appOk"] = ""
		}
	}
}
