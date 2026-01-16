package uniapp

import "gofly/utils/gf"

// 兼容路由：/business/uniapp/user/{action}
// 说明：部分前端可能按“pkg/ctrl/action”三段式调用（uniapp/user/wxLogin）。
// 这里统一转发到 Index，实现一处逻辑多处入口。
type User struct {
	NoNeedLogin []string
	NoNeedAuths []string
}

func init() {
	fpath := User{NoNeedLogin: []string{"login", "wxLogin", "logout"}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

func (api *User) Login(c *gf.GinCtx)       { (&Index{}).Login(c) }
func (api *User) WxLogin(c *gf.GinCtx)     { (&Index{}).WxLogin(c) }
func (api *User) GetUserinfo(c *gf.GinCtx) { (&Index{}).GetUserinfo(c) }
func (api *User) Logout(c *gf.GinCtx)      { (&Index{}).Logout(c) }
