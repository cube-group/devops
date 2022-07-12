//用户
package models

import (
	"app/library/api"
	v1 "app/library/consts/v1"
	"app/library/crypt/md5"
	"app/library/types/slice"
	"app/library/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

const (
	UserRoot = "root"
	UserTest = "test"
)

var UserBossList = []string{
	UserRoot, UserTest,
}

func GetUser(values ...interface{}) (res *User) {
	for _, v := range values {
		switch vv := v.(type) {
		case uint32:
			var i User
			if err := DB().Take(&i, "id=?", vv).Error; err != nil {
				return nil
			}
			res = &i
		case string:
			var i User
			if err := DB().Take(&i, "username=?", vv).Error; err != nil {
				return nil
			}
			res = &i
		case *gin.Context:
			res = SessionUser(vv)
		}
	}
	return
}

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
	Username       string    `gorm:"" json:"username" form:"username"`
	RealName       string    `gorm:"" json:"realName" form:"realName"`
	From           string    `gorm:"" json:"from" form:"from"`
	Email          string    `gorm:"" json:"email" form:"email"`
	Password       string    `gorm:"" json:"password" form:"password"`
	AvatarUrl      string    `gorm:"" json:"avatarUrl" form:"avatarUrl"`
	AvatarBlob     string    `gorm:"" json:"-" form:"-"`
	WebUrl         string    `gorm:"" json:"webUrl" form:"webUrl"`
	Adm            uint8     `gorm:"" json:"adm" form:"adm"`          //是否为超管
	Token          string    `gorm:"" json:"-" form:"-"`              //最近一次登录token
	TokenCreatedAt time.Time `gorm:"" json:"tokenCreatedAt" form:"-"` //最近一次时间
	TokenExpiredAt time.Time `gorm:"" json:"-" form:"-"`              //过期时间

	CreatedAt time.Time      `json:"createdAt" form:"-"`
	UpdatedAt time.Time      `json:"updatedAt" form:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" form:"-"`
}

func (t *User) TableName() string {
	return "d_user"
}

func (t *User) Validator() error {
	if t.Username == "" || t.Password == "" {
		return errors.New("用户名或密码不能为空")
	}
	if !slice.InArrayString(t.Username, UserBossList) {
		t.Password = md5.MD5("visible." + t.Password)
	}
	if t.AvatarUrl == "" {
		t.AvatarUrl = "/public/assets/img/avatar.png"
	}
	return nil
}

func (t User) MarshalJSON() ([]byte, error) {
	//TODO 自定义AvatarUrl
	if t.AvatarBlob != "" {
		t.AvatarUrl = fmt.Sprintf("/api/img/avatar/%d", t.ID)
	}

	return json.Marshal(struct {
		UserMarshalJSON
		Password string `json:"password"`
	}{
		UserMarshalJSON: UserMarshalJSON(t),
		Password:        "",
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

func (t *User) IsRoot() bool {
	return t.Username == UserRoot
}

func (t *User) PwdAddUser(password string) error {
	if t.IsRoot() {
		t.Password = password
	} else {
		t.Password = md5.MD5("visible." + password)
	}
	return DB().Save(t).Error
}

func (t *User) PwdCheckUser(username, password string) *User {
	t.Username = username
	if slice.InArrayString(username, UserBossList) {
		t.Password = password
	} else {
		t.Password = md5.MD5("visible." + password)
	}
	if DB().Where(t).Take(&t).Error != nil {
		return nil
	}
	return t
}

func (t *User) HasPermissionProject(pid uint32) error {
	if t.IsAdm() {
		return nil
	}
	return DB().Take(&ProjectUser{}, "uid=? AND pid=?", t.ID, pid).Error
}

func (t *User) HasOwnPermissionProject(pid uint32) error {
	if t.IsAdm() {
		return nil
	}
	if i := GetProject(pid); i != nil {
		if i.Uid == t.ID {
			return nil
		}
	}
	return errors.New("no own permission")
}

func CreateUser(username string) (user *User, err error) {
	var rootRandPwd string
	var find = GetUser(username)
	if find != nil {
		user = find
	} else {
		rootRandPwd = md5.MD5(uuid.GetUUID("devops", username))
		var newUser = &User{
			Username: username,
			RealName: username,
			From:     "init",
			Adm:      1,
		}
		if err = newUser.PwdAddUser(rootRandPwd); err != nil {
			return
		}
		user = newUser
	}
	return
}

//同步gitlab用户
func SyncGitlabUser(gitlabUser *api.GitlabUser) (u *User, err error) {
	imgContent, imgURL := gitlabUser.GetAvatar()
	var find User
	err = DB().Transaction(func(tx *gorm.DB) error {
		if tx.Take(&find, "username=?", gitlabUser.Username).Error != nil {
			find = User{
				Username: gitlabUser.Username,
			}
		}
		if imgContent != "" {
			find.AvatarBlob = imgContent
		} else {
			find.AvatarUrl = imgURL
		}
		find.From = "gitlab"
		find.RealName = gitlabUser.Name
		find.Email = gitlabUser.Email
		find.WebUrl = gitlabUser.WebUrl
		return tx.Save(&find).Error
	})
	if err != nil {
		return
	}
	return &find, nil
}
