export enum ApiEnum {
  // api前缀
  Prefix = '/api',

  // 基础
  SiteRegister = '/site/register', // 账号注册
  SiteAccountLogin = '/site/accountLogin', // 账号(只支持邮箱+手机号)登录
  SiteAccountCode = '/site/accountCode', // 账号登录验证码
  SiteLoginConfig = '/site/loginConfig', // 登录配置
  SiteLogout = '/site/logout', // 注销
  SiteConfig = '/site/config', // 配置信息

  // 用户
  MemberInfo = '/member/info', // 登录用户信息

  // 角色
  RoleDynamic = '/role/dynamic', // 动态路由
}
