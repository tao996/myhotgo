// Package common
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package common

import (
	"context"
	"hotgo/api/admin/common"
	"hotgo/internal/consts"
	"hotgo/internal/library/captcha"
	"hotgo/internal/library/token"
	"hotgo/internal/service"
	"hotgo/utility/simple"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/glog"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"
)

var Site = cSite{}

type cSite struct{}

// Ping ping
func (c *cSite) Ping(_ context.Context, _ *common.SitePingReq) (res *common.SitePingRes, err error) {
	return
}

// Config 获取配置
func (c *cSite) Config(ctx context.Context, _ *common.SiteConfigReq) (res *common.SiteConfigRes, err error) {
	request := ghttp.RequestFromCtx(ctx)
	res = &common.SiteConfigRes{
		Version: consts.VersionApp,
		WsAddr:  c.getWsAddr(ctx, request),
		Domain:  c.getDomain(ctx, request),
		Mode:    gmode.Mode(),
	}
	return
}

func (c *cSite) getWsAddr(ctx context.Context, request *ghttp.Request) string {
	// 如果是本地IP访问，则认为是调试模式，走实际请求地址，否则走配置中的地址
	// 尝试读取hostname，兼容本地运行模式
	ip := ghttp.RequestFromCtx(ctx).GetHeader("hostname")
	if len(ip) == 0 {
		ip = ghttp.RequestFromCtx(ctx).GetHost()
	}

	if validate.IsLocalIPAddr(ip) {
		return "ws://" + ip + ":" + gstr.StrEx(request.Host, ":") + g.Cfg().MustGet(ctx, "router.websocket.prefix").String()
	}

	basic, err := service.SysConfig().GetBasic(ctx)
	if err != nil || basic == nil {
		return ""
	}
	return basic.WsAddr
}

func (c *cSite) getDomain(ctx context.Context, request *ghttp.Request) string {
	// 如果是本地IP访问，则认为是调试模式，走实际请求地址，否则走配置中的地址
	// 尝试读取hostname，兼容本地运行模式
	ip := ghttp.RequestFromCtx(ctx).GetHeader("hostname")
	if len(ip) == 0 {
		ip = ghttp.RequestFromCtx(ctx).GetHost()
	}

	if validate.IsLocalIPAddr(ip) {
		return "http://" + ip + ":" + gstr.StrEx(request.Host, ":")
	}

	basic, err := service.SysConfig().GetBasic(ctx)
	if err != nil || basic == nil {
		return ""
	}
	return basic.Domain
}

// LoginConfig 登录配置
func (c *cSite) LoginConfig(ctx context.Context, _ *common.SiteLoginConfigReq) (res *common.SiteLoginConfigRes, err error) {
	res = new(common.SiteLoginConfigRes)
	login, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}

	res.LoginConfig = login
	res.I18nSwitch = g.Cfg().MustGet(ctx, "system.i18n.switch", true).Bool()
	res.DefaultLanguage = g.Cfg().MustGet(ctx, "system.i18n.defaultLanguage", consts.SysDefaultLanguage).String()
	res.ProjectName = gi18n.T(ctx, g.Cfg().MustGet(ctx, `setting.projectName`).String())
	return
}

// Captcha 登录验证码
func (c *cSite) Captcha(ctx context.Context, _ *common.LoginCaptchaReq) (res *common.LoginCaptchaRes, err error) {
	loginConf, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}
	cid, base64 := captcha.Generate(ctx, loginConf.CaptchaType)
	res = &common.LoginCaptchaRes{Cid: cid, Base64: base64}
	return
}

// AccountCode 发送登录验证码
func (c *cSite) AccountCode(ctx context.Context, req *common.AccountCodeReq) (res *common.AccountCodeRes, err error) {
	if req.Account == "" {
		if req.Mobile != "" {
			req.Account = req.Mobile
		} else if req.Email != "" {
			req.Account = req.Email
		} else {
			return nil, gerror.New("请输入账号")
		}
	}
	login, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}
	req.Mock = false
	if login.CaptchaSwitch == consts.StatusEnabled {
		// 校验 验证码
		if !captcha.Verify(req.Cid, req.Captcha, false) {
			if simple.Debug(ctx) && req.Captcha == consts.MockCaptcha {
				glog.Debug(ctx, "Debug 模式：跳过图形验证码错误")
				req.Mock = true
			} else {
				err = gerror.New("图形验证码错误")
				return
			}
		}
	}
	err = service.AdminSite().AccountCode(ctx, &req.AccountCodeInp)
	return
}

// AccountLogin 账号登录
func (c *cSite) AccountLogin(ctx context.Context, req *common.AccountLoginReq) (res *common.AccountLoginRes, err error) {

	if req.Account == "" {
		if req.Username != "" {
			req.Account = req.Username
		} else if req.Mobile != "" {
			req.Account = req.Mobile
		} else if req.Email != "" {
			req.Account = req.Email
		} else {
			return nil, gerror.New("请输入账号")
		}
	}
	login, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}
	req.Mock = false
	if login.CaptchaSwitch == consts.StatusEnabled {
		// 校验 验证码
		if !captcha.Verify(req.Cid, req.Captcha, true) {
			if simple.Debug(ctx) && req.Captcha == consts.MockCaptcha {
				req.Mock = true
				glog.Debug(ctx, "Debug 模式：跳过图形验证码错误")
			} else {
				err = gerror.New("图形验证码错误")
				return
			}
		}
	}

	model, err := service.AdminSite().AccountLogin(ctx, &req.AccountLoginInp)
	if err != nil {
		return nil, gerror.Wrap(err, "账号登录失败")
	}

	err = gconv.Scan(model, &res)
	return
}

// Logout 注销登录
func (c *cSite) Logout(ctx context.Context, _ *common.LoginLogoutReq) (res *common.LoginLogoutRes, err error) {
	err = token.Logout(ghttp.RequestFromCtx(ctx))
	return
}
