package models

import (
	"app/library/consts"
	v1 "app/library/consts/v1"
	"app/library/crypt/md5"
	"app/library/ginutil"
	"app/library/types/convert"
	"app/library/types/jsonutil"
	"app/library/uuid"
	"app/setting"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"sync"
	"time"
)

var userTokenMap sync.Map
var userTokenCheckMap sync.Map

//获取用户ID
func SessionUid(c *gin.Context) (res uint) {
	if user := SessionUser(c); user != nil {
		return uint(user.ID)
	}
	return 0
}

func UID(values ...interface{}) uint32 {
	var user *User
	for _, v := range values {
		switch vv := v.(type) {
		case *gin.Context:
			user = SessionUser(vv)
		}
	}
	if user == nil {
		return 0
	}
	return user.ID
}

//获取用户登录名
func SessionUsername(c *gin.Context) string {
	if user := SessionUser(c); user != nil {
		return user.Username
	}
	return ""
}

//获取用户姓名
func SessionUserRealName(c *gin.Context) string {
	if user := SessionUser(c); user != nil {
		return user.RealName
	}
	return ""
}

//清空session
func SessionClear(c *gin.Context) {
	c.SetCookie(v1.SessionName, "", -1, "/", setting.SysWebDomain, false, false)
}

// 获取用户session
func SessionUser(c *gin.Context) *User {
	//load from context
	if i, ok := c.Get(consts.ContextUser); ok {
		return i.(*User)
	}
	sessionToken, _ := c.Cookie(v1.SessionName)
	if sessionToken == "" {
		return nil
	}
	var user *User
	//load from map
	if i, ok := userTokenMap.Load(sessionToken); ok {
		if jsoniter.Unmarshal(i.([]byte), &user) == nil && !user.IsExpired() {
			go checkUserStatus(c, user)
			goto Success
		}
	}
	//load from db
	if user = GetUserByToken(sessionToken); user == nil || user.IsExpired() {
		return nil
	}
	userTokenMap.Store(sessionToken, jsonutil.ToBytes(user))
	c.Set(consts.ContextUser, user)
Success:
	return user
}

// 是否为超级管理员
func IsAdm(c *gin.Context) bool {
	return SessionUser(c).IsAdm()
}

//延迟检测用户实际登录状态
func checkUserStatus(c *gin.Context, user *User) {
	var md5 = uuid.GetUUID("checkUserStatus", convert.MustString(user.ID))
	userTokenCheckMap.Store(user.ID, md5)
	time.Sleep(5 * time.Second)
	if i, ok := userTokenCheckMap.Load(user.ID); ok {
		if i.(string) != md5 {
			return
		}
	}
	latestUser := GetUser(user.ID)
	if latestUser == nil {
		return
	}
	if !latestUser.IsExpired() {
		return
	}
	SessionClear(c)
}

func SessionURLSignGet(c *gin.Context) string {
	timestamp := time.Now().Unix()
	sign := md5.MD5(fmt.Sprintf("t=%d&uid=%d", timestamp, SessionUid(c)))
	return fmt.Sprintf("t=%d&sign=%s", timestamp, sign)
}

func SessionURLSignCheck(c *gin.Context) error {
	timestamp := convert.MustInt64(ginutil.Input(c, "t"))
	if time.Now().After(time.Unix(timestamp, 0).Add(time.Hour * 6)) {
		return errors.New("url安全校验时间戳过期")
	}
	sign := md5.MD5(fmt.Sprintf("t=%d&uid=%d", timestamp, SessionUid(c)))
	if sign != ginutil.Input(c, "sign") {
		return errors.New("url安全校验签名错误")
	}
	return nil
}
