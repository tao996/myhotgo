// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/simple"
	"hotgo/utility/validate"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

type sAdminSite struct{}

func NewAdminSite() service.IAdminSite {
	return &sAdminSite{}
}

func init() {
	service.RegisterAdminSite(NewAdminSite())
}

// Register 账号注册
func (s *sAdminSite) Register(ctx context.Context, in *adminin.RegisterInp) (err error) {
	config, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}

	if config.ForceInvite == 1 && in.InviteCode == "" {
		err = gerror.New("请填写邀请码")
		return
	}

	var data adminin.MemberAddInp

	// 默认上级
	data.Pid = 1

	// 存在邀请人
	if in.InviteCode != "" {
		pmb, err := service.AdminMember().GetIdByCode(ctx, &adminin.GetIdByCodeInp{Code: in.InviteCode})
		if err != nil {
			return err
		}

		if pmb == nil {
			err = gerror.New("邀请人信息不存在")
			return err
		}

		data.Pid = pmb.Id
	}

	if config.RegisterSwitch != 1 {
		err = gerror.New("管理员未开放注册")
		return
	}

	if config.RoleId < 1 {
		err = gerror.New("管理员未配置默认角色")
		return
	}

	if config.DeptId < 1 {
		err = gerror.New("管理员未配置默认部门")
		return
	}

	// 验证唯一性
	err = service.AdminMember().VerifyUnique(ctx, &adminin.VerifyUniqueInp{
		Where: g.Map{
			dao.AdminMember.Columns().Username: in.Username,
			dao.AdminMember.Columns().Mobile:   in.Mobile,
		},
	})
	if err != nil {
		return
	}

	// 验证短信验证码
	err = service.SysSmsLog().VerifyCode(ctx, &sysin.VerifyCodeInp{
		Event:  consts.SmsTemplateRegister,
		Mobile: in.Mobile,
		Code:   in.Code,
	})
	if err != nil {
		return
	}

	data.MemberEditInp = &adminin.MemberEditInp{
		Id:       0,
		RoleId:   config.RoleId,
		PostIds:  config.PostIds,
		DeptId:   config.DeptId,
		Username: in.Username,
		Password: in.Password,
		RealName: "",
		Avatar:   config.Avatar,
		Sex:      3, // 保密
		Mobile:   in.Mobile,
		Status:   consts.StatusEnabled,
	}
	data.Salt = grand.S(6)
	data.InviteCode = grand.S(12)
	data.PasswordHash = gmd5.MustEncryptString(data.Password + data.Salt)
	data.Level, data.Tree, err = service.AdminMember().GenTree(ctx, data.Pid)
	if err != nil {
		return
	}

	// 提交注册信息
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		id, err := dao.AdminMember.Ctx(ctx).Data(data).OmitEmptyData().InsertAndGetId()
		if err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return
		}

		// 更新岗位
		if err = service.AdminMemberPost().UpdatePostIds(ctx, id, config.PostIds); err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
		}
		return
	})
}

func (s *sAdminSite) getAccountType(ctx context.Context, account string) (string, error) {

	if validate.IsEmail(account) {
		return "email", nil
	} else if validate.IsMobile(account) {
		return "mobile", nil
	} else {
		if simple.Debug(ctx) {
			return "username", nil
		}
	}
	return "", gerror.New("只支持邮箱或手机号登录")
}

func (s *sAdminSite) AccountCode(ctx context.Context, in *adminin.AccountCodeInp) (err error) {
	accountType, err := s.getAccountType(ctx, in.Account)
	if err != nil {
		return err
	}
	if exist, err := dao.AdminMember.Ctx(ctx).Where(accountType, in.Account).Exist(); err != nil {
		return gerror.Wrap(err, "检查账号时错误")
	} else if !exist { // 账号不存在则检查是否允许注册
		// 只有手机号和邮箱能够注册
		if accountType != "email" && accountType != "mobile" {
			return gerror.New("不允许注册的账号类型")
		}
		if config, err := service.SysConfig().GetLogin(ctx); err != nil {
			return err
		} else {
			if config.RegisterSwitch != 1 {
				return gerror.New("管理员未开放注册")
			}
			if config.RoleId < 1 {
				return gerror.New("管理员未配置默认角色")
			}

			if config.DeptId < 1 {
				return gerror.New("管理员未配置默认部门")
			}
			if config.ForceInvite == 1 && in.InviteCode == "" {
				return gerror.New("请填写邀请码")
			}
			if in.InviteCode != "" {
				pmb, err := service.AdminMember().GetIdByCode(ctx, &adminin.GetIdByCodeInp{Code: in.InviteCode})
				if err != nil {
					return err
				}

				if pmb == nil {
					err = gerror.New("邀请人信息不存在")
					return err
				}
			}
		}
	}
	if accountType == "mobile" {
		if err = service.SysSmsLog().SendCode(ctx, &sysin.SendCodeInp{
			Event:  consts.SmsTemplateLogin,
			Mobile: in.Account,
			Mock:   in.Mock,
		}); err != nil {
			return err
		}

	} else if accountType == "email" {
		if err = service.SysEmsLog().Send(ctx, &sysin.SendEmsInp{
			Event: consts.EmsTemplateLogin,
			Email: in.Account,
			Mock:  in.Mock,
		}); err != nil {
			return err
		}
	}

	return nil
}

// AccountLogin 账号登录
func (s *sAdminSite) AccountLogin(ctx context.Context, in *adminin.AccountLoginInp) (res *adminin.LoginModel, err error) {
	defer func() {
		service.SysLoginLog().Push(ctx, &sysin.LoginLogPushInp{Response: res, Err: err})
	}()

	if in.Password == "" && in.Code == "" {
		err = gerror.New("请填写密码或验证码")
		return
	}
	name, err := s.getAccountType(ctx, in.Account)
	if err != nil {
		return nil, gerror.Wrap(err, "账号格式错误")
	}
	config, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(name, in.Account).Scan(&mb); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if mb == nil {
		err = gerror.New("用户名或密码错误")
		return
	}
	//	密码登录
	// 如果账号不存在 mb == nil，则直接注册
	var data adminin.MemberAddInp
	if in.Password != "" { // 密码登录则账号必须存在
		if mb == nil {
			err = gerror.New("请检查账号是否正确")
			return
		}
		if err = simple.CheckPassword(in.Password, mb.Salt, mb.PasswordHash); err != nil {
			return
		}
	} else { // 验证码登录
		if mb == nil {
			data.Pid = 1
			if in.InviteCode != "" {
				pmb, err := service.AdminMember().GetIdByCode(ctx, &adminin.GetIdByCodeInp{Code: in.InviteCode})
				if err != nil {
					return nil, err
				}

				if pmb == nil {
					err = gerror.New("邀请人信息不存在")
					return nil, err
				}

				data.Pid = pmb.Id
			}
			data.MemberEditInp = &adminin.MemberEditInp{
				Id:       0,
				RoleId:   config.RoleId,
				PostIds:  config.PostIds,
				DeptId:   config.DeptId,
				Username: "",
				Password: in.Password,
				RealName: "",
				Avatar:   config.Avatar,
				Sex:      3, // 保密
				Mobile:   "",
				Status:   consts.StatusEnabled,
			}
			if name == "mobile" {
				data.MemberEditInp.Mobile = in.Account
			} else if name == "email" {
				data.MemberEditInp.Email = in.Account
			} else {
				data.MemberEditInp.Username = in.Account
			}

			data.Salt = grand.S(6)
			data.InviteCode = grand.S(12)
			if in.Password != "" {
				data.PasswordHash = gmd5.MustEncryptString(data.Password + data.Salt)
			}
			data.Level, data.Tree, err = service.AdminMember().GenTree(ctx, data.Pid)
			if err != nil {
				return nil, gerror.Wrap(err, "查询部门信息时错误")
			}
		}

		if name == "mobile" {
			if err = service.SysSmsLog().VerifyCode(ctx, &sysin.VerifyCodeInp{
				Event:  consts.SmsTemplateLogin,
				Mobile: in.Account,
				Code:   in.Code,
			}); err != nil {
				return nil, gerror.Wrap(err, "登录验证码错误")
			}
		} else {
			if err = service.SysEmsLog().VerifyCode(ctx, &sysin.VerifyEmsCodeInp{
				Event: consts.EmsTemplateLogin,
				Email: in.Account,
				Code:  in.Code,
			}); err != nil {
				return nil, gerror.Wrap(err, "登录验证码错误")
			}
		}

	}
	if mb == nil {
		// 提交注册信息
		if err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
			id, err := dao.AdminMember.Ctx(ctx).Data(data).OmitEmptyData().InsertAndGetId()
			if err != nil {
				err = gerror.Wrap(err, consts.ErrorORM)
				return
			}

			// 更新岗位
			if err = service.AdminMemberPost().UpdatePostIds(ctx, id, config.PostIds); err != nil {
				err = gerror.Wrap(err, consts.ErrorORM)
			}
			data.Id = id
			return
		}); err != nil {
			return nil, gerror.Wrap(err, "注册失败")
		}
		if err = dao.AdminMember.Ctx(ctx).Where("id", data.Id).Scan(&mb); err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return
		}
	}

	if mb.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用")
		return
	}

	res, err = s.handleLogin(ctx, mb)
	return
}

// handleLogin .
func (s *sAdminSite) handleLogin(ctx context.Context, mb *entity.AdminMember) (res *adminin.LoginModel, err error) {
	role, dept, err := s.getLoginRoleAndDept(ctx, mb.RoleId, mb.DeptId)
	if err != nil {
		return nil, err
	}

	user := &model.Identity{
		Id:       mb.Id,
		Pid:      mb.Pid,
		DeptId:   dept.Id,
		DeptType: dept.Type,
		RoleId:   role.Id,
		RoleKey:  role.Key,
		Username: mb.Username,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      consts.AppAdmin,
		LoginAt:  gtime.Now(),
	}

	lt, expires, err := token.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	res = &adminin.LoginModel{
		Username: user.Username,
		Id:       user.Id,
		Token:    lt,
		Expires:  expires,
	}
	if res.Username == "" {
		if mb.Email != "" {
			res.Username = mb.Email
		} else if mb.Mobile != "" {
			res.Username = mb.Mobile
		}
	}
	return
}

// getLoginRoleAndDept 获取登录的角色和部门信息
func (s *sAdminSite) getLoginRoleAndDept(ctx context.Context, roleId, deptId int64) (role *entity.AdminRole, dept *entity.AdminDept, err error) {
	if err = dao.AdminRole.Ctx(ctx).Fields(dao.AdminRole.Columns().Id, dao.AdminRole.Columns().Key, dao.AdminRole.Columns().Status).WherePri(roleId).Scan(&role); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if role == nil {
		err = gerror.New("角色不存在或已被删除")
		return
	}

	if role.Status != consts.StatusEnabled {
		err = gerror.New("角色已被禁用，如有疑问请联系管理员")
		return
	}

	if err = dao.AdminDept.Ctx(ctx).Fields(dao.AdminDept.Columns().Id, dao.AdminDept.Columns().Type, dao.AdminDept.Columns().Status).WherePri(deptId).Scan(&dept); err != nil {
		err = gerror.Wrap(err, "获取部门信息失败，请稍后重试！")
		return
	}

	if dept == nil {
		err = gerror.New("部门不存在或已被删除")
		return
	}

	if dept.Status != consts.StatusEnabled {
		err = gerror.New("部门已被禁用，如有疑问请联系管理员")
		return
	}
	return
}

// BindUserContext 绑定用户上下文
func (s *sAdminSite) BindUserContext(ctx context.Context, claims *model.Identity) (err error) {
	//// 如果不想每次访问都重新加载用户信息，可以放开注释。但在本次登录未失效前，用户信息不会刷新
	// contexts.SetUser(ctx, claims)
	// return

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(claims.Id).Scan(&mb); err != nil {
		err = gerror.Wrap(err, "获取用户信息失败，请稍后重试！")
		return
	}

	if mb == nil {
		err = gerror.Wrap(err, "账号不存在或已被删除！")
		return
	}

	if mb.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用，如有疑问请联系管理员")
		return
	}

	role, dept, err := s.getLoginRoleAndDept(ctx, mb.RoleId, mb.DeptId)
	if err != nil {
		return err
	}

	user := &model.Identity{
		Id:       mb.Id,
		Pid:      mb.Pid,
		DeptId:   dept.Id,
		DeptType: dept.Type,
		RoleId:   mb.RoleId,
		RoleKey:  role.Key,
		Username: mb.Username,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      claims.App,
		LoginAt:  claims.LoginAt,
	}

	contexts.SetUser(ctx, user)
	return
}
