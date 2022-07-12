// 基础与大数据部
// API Control Instance: gitlab 操作实例
// Author: linyang
// Date: 2018-10
package api

import (
	"app/library/crypt/md5"
	"app/library/types/jsonutil"
	"context"
	"errors"
	"fmt"
	"github.com/imroc/req"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// gitlab user info
// access_level master所对应的access_level为40，developer的权限为30，将access_level从40改成30就实现了降权
type GitlabUser struct {
	Id               uint      `json:"id"`                 //用户id
	Name             string    `json:"name"`               //用户昵称
	Username         string    `json:"username"`           //用户登录名称
	State            string    `json:"state"`              //用户状态,active、blocked
	AvatarUrl        string    `json:"avatar_url"`         //用户头像地址
	WebUrl           string    `json:"web_url"`            //用户主页地址
	CreateAt         time.Time `json:"create_at"`          //用户创建时间
	LastSignInAt     time.Time `json:"last_sign_in_at"`    //用户最近一次登录时间
	Email            string    `json:"email"`              //用户邮箱
	CanCreateGroup   bool      `json:"can_create_group"`   //用户是否具备权限创建群组
	CanCreateProject bool      `json:"can_create_project"` //用户是否具备权限创建项目
	IsAdmin          bool      `json:"is_admin"`           //用户是否为管理员
	AccessLevel      uint      `json:"access_level"`       //用户所处权限
	AccessToken      string    `json:"-"`
}

// gitlab operate struct
type Gitlab struct {
	Option GitlabOption
}

type GitlabOption struct {
	GitlabAddress     string
	GitlabAppId       string
	GitlabAppSecret   string
	GitlabRedirectUri string
	State             GitlabState
}

type GitlabState struct {
	Value string `json:"v"`
	Ref   string `json:"ref"`
}

// create Gitlab sync Instance
func NewGitlab(option GitlabOption) *Gitlab {
	return &Gitlab{Option: option}
}

// 获取gitlab授权跳转地址
func (t *Gitlab) GetAuthURL(ref string) string {
	redirectUrl := fmt.Sprintf(
		"%s/oauth/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s",
		t.Option.GitlabAddress,
		t.Option.GitlabAppId,
		md5.Base64Encode(jsonutil.ToString(GitlabState{Ref: ref, Value: "test"})),
		url.QueryEscape(t.Option.GitlabRedirectUri),
	)
	return redirectUrl
}

// 获取用户信息
// param code: gitlab redirect code
func (t *Gitlab) GetUserInfo(code string) (user *GitlabUser, err error) {
	c := &oauth2.Config{
		ClientID:     t.Option.GitlabAppId,
		ClientSecret: t.Option.GitlabAppSecret,
		RedirectURL:  t.Option.GitlabRedirectUri,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth/authorize", t.Option.GitlabAddress),
			TokenURL: fmt.Sprintf("%s/oauth/token", t.Option.GitlabAddress),
		},
	}
	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		return
	}

	requestUrl := fmt.Sprintf("%s/api/v4/user?access_token=%s", t.Option.GitlabAddress, token.AccessToken)
	resp, err := req.Get(requestUrl)
	if err != nil {
		return
	}
	if resp.Response().StatusCode != http.StatusOK {
		err = errors.New(resp.String())
		return
	}
	var u GitlabUser
	if err = jsoniter.Unmarshal(resp.Bytes(), &u); err != nil {
		return
	}
	u.AccessToken = token.AccessToken
	return &u, nil
}

func (t *GitlabUser) getGitlabAddressURL(gitlabAddress, requestURL string) string {
	parseURL, err := url.Parse(requestURL)
	if err != nil {
		return t.WebUrl
	}
	gitlabAddressURL, err := url.Parse(gitlabAddress)
	if err != nil {
		return t.WebUrl
	}
	parseURL.Scheme = gitlabAddressURL.Scheme
	parseURL.Host = gitlabAddressURL.Host
	return parseURL.String()
}

func (t *GitlabUser) GetAvatar() (imgContent string, imgURL string) {
	var requestURL = t.AvatarUrl
	if strings.Contains(requestURL, "?") {
		requestURL += "&access_token=" + t.AccessToken
	} else {
		requestURL += "?access_token=" + t.AccessToken
	}
	resp, err := req.Get(requestURL)
	if err == nil && resp.Response().StatusCode == http.StatusOK {
		imgContent = resp.String()
	}
	if imgContent == "" {
		imgURL = "/public/assets/img/avatar.png"
	}
	return
}

func (t *GitlabUser) GetWebURL(gitlabAddress string) string {
	return t.getGitlabAddressURL(gitlabAddress, t.WebUrl)
}
