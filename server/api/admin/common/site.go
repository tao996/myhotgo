// Package common
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package common

import (
	"hotgo/internal/model"
	"hotgo/internal/model/input/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

// LoginLogoutReq 注销登录
type LoginLogoutReq struct {
	g.Meta `path:"/site/logout" method:"post" tags:"后台基础" summary:"注销登录"`
}

type LoginLogoutRes struct{}

// RegisterReq 提交账号注册
type RegisterReq struct {
	g.Meta `path:"/site/register" method:"post" tags:"后台基础" summary:"账号注册"`
	adminin.RegisterInp
}

type RegisterRes struct {
	*adminin.LoginModel
}

// LoginCaptchaReq 获取登录验证码
type LoginCaptchaReq struct {
	g.Meta `path:"/site/captcha" method:"get" tags:"后台基础" summary:"获取登录验证码"`
}

type LoginCaptchaRes struct {
	Cid    string `json:"cid" dc:"验证码ID"`
	Base64 string `json:"base64" dc:"验证码"`
}

// AccountLoginReq 提交账号登录
type AccountLoginReq struct {
	g.Meta `path:"/site/accountLogin" method:"post" tags:"后台基础" summary:"账号登录"`
	adminin.AccountLoginInp
	Username string `json:"username" dc:"用户名"`
	Email    string `json:"email" dc:"邮箱地址"`
	Mobile   string `json:"mobile" dc:"手机号"`
}

type AccountLoginRes struct {
	*adminin.LoginModel
	Email  string `json:"email" dc:"邮箱地址"`
	Mobile string `json:"mobile" dc:"手机号"`
}

// SiteConfigReq 获取配置
type SiteConfigReq struct {
	g.Meta `path:"/site/config" method:"get" tags:"后台基础" summary:"获取配置"`
}

type SiteConfigRes struct {
	Version string `json:"version"        dc:"系统版本"`
	WsAddr  string `json:"wsAddr"         dc:"客户端websocket地址"`
	Domain  string `json:"domain"         dc:"对外域名"`
	Mode    string `json:"mode"           dc:"运行模式"`
}

// SiteLoginConfigReq 获取登录配置
type SiteLoginConfigReq struct {
	g.Meta `path:"/site/loginConfig" method:"get" tags:"后台基础" summary:"获取登录配置"`
}

type SiteLoginConfigRes struct {
	*model.LoginConfig
	I18nSwitch      bool   `json:"i18nSwitch" dc:"国际化开关"`
	DefaultLanguage string `json:"defaultLanguage" dc:"默认语言设置"`
	ProjectName     string `json:"projectName" dc:"项目名称"`
}

// SitePingReq ping
type SitePingReq struct {
	g.Meta `path:"/site/ping" method:"get" tags:"后台基础" summary:"ping"`
}

type SitePingRes struct{}

type AccountCodeReq struct {
	g.Meta `path:"/site/accountCode" method:"post" tags:"后台基础" summary:"发送登录验证码"`
	adminin.AccountCodeInp
	Email  string `json:"email" dc:"邮箱地址"`
	Mobile string `json:"mobile" dc:"手机号"`
}

type AccountCodeRes struct {
}
