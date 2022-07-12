package open

import (
	"app/library/api"
	"app/library/crypt/base64"
	"app/models"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

// gitlab callback 携带参数结构体
type GitlabCallbackVal struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// 处理gitlab callback oauth 请求
func OauthCallback(c *gin.Context) (ref string, err error) {
	var val GitlabCallbackVal
	if err = c.ShouldBind(&val); err != nil {
		return
	}
	var state api.GitlabState
	if err = json.Unmarshal([]byte(base64.Base64Decode(val.State)), &state); err != nil {
		return
	}
	if state.Value != "test" {
		err = errors.New("gitlab oauth state invalid")
		return
	}

	gitlab, err := models.Gitlab()
	if err != nil {
		return
	}

	gitlabUser, err := gitlab.GetUserInfo(val.Code)
	if err != nil {
		return
	}

	//处理webUrl
	gitlabUser.WebUrl = gitlabUser.GetWebURL(gitlab.Option.GitlabAddress)

	user, err := models.SyncGitlabUser(gitlabUser)
	if err != nil {
		return
	}

	if err = user.LoginCookie(c); err != nil {
		return
	}
	ref = state.Ref
	return
}
