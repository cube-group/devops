//用户
package models

import (
	v1 "app/library/consts/v1"
	"app/library/crypt/md5"
	"app/library/uuid"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

//获取所有开发者了列表
func GetAllUser() []User {
	return new(User).GetAll()
}

func GetUserByToken(token string) *User {
	var i User
	if DB().Take(&i, "token=?", token).Error != nil {
		return nil
	}
	return &i
}

type UserMarshalJSON User

//users表模型
type User struct {
	ID             uint32    `gorm:"primarykey;column:id" json:"id"`
	Username       string    `gorm:""`
	RealName       string    `gorm:""`
	Email          string    `gorm:""`
	Password       string    `gorm:"" json:"-"`
	AvatarUrl      string    `gorm:""`
	Adm            uint8     `gorm:""` //是否为超管
	Token          string    `gorm:""` //最近一次登录token
	TokenCreatedAt time.Time `gorm:""` //最近一次时间
	TokenExpiredAt time.Time `gorm:""` //过期时间

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *User) TableName() string {
	return "d_user"
}

func (t User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		UserMarshalJSON
		IsBlock bool
	}{
		UserMarshalJSON: UserMarshalJSON(t),
	})
}

func (t *User) UID() uint {
	return uint(t.ID)
}

//是否是超级管理员，拥有ID为1的角色
func (t *User) GetAll() (res []User) {
	DB().Order("id DESC").Find(&res)
	return
}

//token 是否过期
func (t *User) IsExpired() bool {
	return time.Now().After(t.TokenExpiredAt)
}

//是否是超级管理员，拥有ID为1的角色
func (t *User) IsAdm() bool {
	return t.Adm == 1
}

//登录写入cookie
func (t *User) LoginCookie(c *gin.Context) error {
	t.TokenCreatedAt = time.Now()
	t.TokenExpiredAt = time.Now().Add(24 * time.Hour)
	t.Token = uuid.GetUUID("login.token", t.ID)
	if err := DB().Save(t).Error; err != nil {
		return err
	}
	c.SetCookie(v1.SessionName, t.Token, 86400, "/", c.Request.Host, false, false)
	return nil
}

//查询用户
func (t *User) Get(option ...interface{}) *User {
	var gid uint
	var username string
	for _, v := range option {
		switch vv := v.(type) {
		case uint:
			gid = vv
		case string:
			username = vv
		case *gin.Context:
			return SessionUser(vv)
		}
	}
	var user User
	if gid > 0 {
		if err := DB().Where("id=?", gid).Take(&user).Error; err != nil {
			return nil
		}
	} else if username != "" {
		if err := DB().Where("username=?", username).Take(&user).Error; err != nil {
			return nil
		}
	} else {
		return nil
	}
	return &user
}

func (t *User) IsRoot() bool {
	return t.Username == UserRoot
}

func (t *User) PwdAddUser(password string) error {
	if t.IsRoot() {
		t.Password = password
	} else {
		t.Password = md5.MD5("corecd." + password)
	}
	return DB().Save(t).Error
}

func (t *User) PwdCheckUser(username, password string) *User {
	t.Username = username
	if t.IsRoot() {
		t.Password = password
	} else {
		t.Password = md5.MD5("corecd." + password)
	}
	if DB().Where(t).Take(&t).Error != nil {
		return nil
	}
	return t
}
