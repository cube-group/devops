//gin redis session middle
//author: linyang
//mail: lin2798003@sina.com
//date: 2019-08
package session

import (
	"app/library/crypt/md5"
	"app/library/log"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	SESSION_NAME             = "golib-session-id"
	SESSION_CONTEXT_INSTANCE = "golib-session-context-instance"
)

type SessionOptions struct {
	SessionName string //options

	MaxAge int    //默认为86400s
	Secret string //默认为default
	Domain string
	//Path   string
	//Secure   bool
	//HttpOnly bool

	RedisAddress  string
	RedisPassword string
	RedisDB       int
	RedisPoolSize int
}

var sessionOptions *SessionOptions
var redisConn *redis.Client

//初始化redis连接
func session_redis() {
	if sessionOptions.SessionName == "" {
		sessionOptions.SessionName = SESSION_NAME
	}
	if sessionOptions.Secret == "" {
		sessionOptions.Secret = "default"
	}
	if sessionOptions.MaxAge == 0 {
		sessionOptions.MaxAge = 86400
	}

	var initOnce sync.Once
	initOnce.Do(func() {
		//log.StdOut("MiddleWare", "Redis init")
		redisConn = redis.NewClient(&redis.Options{
			Addr:     sessionOptions.RedisAddress,
			Password: sessionOptions.RedisPassword,
			DB:       sessionOptions.RedisDB,
			PoolSize: sessionOptions.RedisPoolSize,
		})
		if err := redisConn.Ping().Err(); err != nil {
			log.StdFatal("MiddleWare", "redis connect err", err.Error())
		}
	})
}

//生成写入cookie的域
func session_hostname(c *gin.Context) string {
	if sessionOptions.Domain != "" {
		return sessionOptions.Domain
	}
	return strings.Split(c.Request.Host, ":")[0]
}

func session_ssid() (string) {
	sid := md5.MD5(fmt.Sprintf(
		"%s-%s-%d-%d",
		sessionOptions.Secret,
		sessionOptions.SessionName,
		time.Now().Nanosecond(),
		rand.Intn(100000),
	))
	return sid
}

type IRedisSession interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) (error)
	Clear() (error)
}

type RedisSession struct {
	Sid      string //session sid
	Skey     string //session cache key
	Hostname string

	context *gin.Context

	IRedisSession
}

func NewRedisSession(sid string, c *gin.Context) *RedisSession {
	s := &RedisSession{Sid: sid, context: c, Hostname: session_hostname(c)}
	s.Skey = fmt.Sprintf("%s:%s", sessionOptions.SessionName, sid)
	return s
}

//设置key
func (this *RedisSession) Set(key, value string) error {
	if err := redisConn.HSet(this.Skey, key, value).Err(); err != nil {
		return err
	}
	if err := redisConn.Expire(this.Skey, time.Duration(sessionOptions.MaxAge)*time.Second).Err(); err != nil {
		return err
	}
	this.context.SetCookie(
		sessionOptions.SessionName,
		this.Sid,
		sessionOptions.MaxAge,
		"/",
		this.Hostname,
		false,
		true,
	)
	return nil
}

//获取key
func (this *RedisSession) Get(key string) (string, error) {
	cmd := redisConn.HGet(this.Skey, key)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), nil
}

//删除key
func (this *RedisSession) Delete(key string) (error) {
	if err := redisConn.HDel(this.Skey, key).Err(); err != nil {
		return err
	}
	return nil
}

//彻底清除session
func (this *RedisSession) Clear() (error) {
	if err := redisConn.Del(this.Skey).Err(); err != nil {
		return err
	}
	this.context.SetCookie(
		sessionOptions.SessionName,
		"",
		0,
		"/",
		this.Hostname,
		false,
		true,
	)
	return nil
}

type RedisSessionNil struct {
	IRedisSession
}

//设置key
func (this *RedisSessionNil) Set(key, value string) error {
	return errors.New("没有初始化Session.GinHandler")
}

//获取key
func (this *RedisSessionNil) Get(key string) (string, error) {
	return "", errors.New("没有初始化Session.GinHandler")
}

//删除key
func (this *RedisSessionNil) Delete(key string) (error) {
	return errors.New("没有初始化Session.GinHandler")
}

//彻底清除session
func (this *RedisSessionNil) Clear() (error) {
	return errors.New("没有初始化Session.GinHandler")
}

//通过context获取session实例
func Session(c *gin.Context) (IRedisSession) {
	i, exist := c.Get(SESSION_CONTEXT_INSTANCE)
	if !exist || i == nil {
		return new(RedisSessionNil)
	}
	return i.(*RedisSession)
}

//gin redis session中间件
func GinHandler(options *SessionOptions) gin.HandlerFunc {
	sessionOptions = options
	session_redis()

	return func(c *gin.Context) {
		sid, err := c.Cookie(sessionOptions.SessionName)
		if err != nil || sid == "" {
			sid = session_ssid()
		}
		c.Set(sessionOptions.SessionName, sid)
		c.Set(SESSION_CONTEXT_INSTANCE, NewRedisSession(sid, c))
		c.Next()
	}
}
