package uniapp

import (
	housesPkg "gofly/app/business/houses"
	"gofly/utils/gf"
)

// wxEnsureCanManageLocks 校验小程序端是否具备“可管理智能锁”权限
//
// 说明：
// - 小程序端不走 RBAC 菜单权限（NoNeedAuths="*"），但业务必须按字段授权
// - 当前权限粒度：business_user.can_manage_locks = 1
func wxEnsureCanManageLocks(c *gf.GinCtx) (businessID int64, userID int64, ok bool) {
	businessID = wxBusinessID(c)
	// 兼容：某些链路 businessID 未写入 context
	if businessID > 0 && c.GetInt64("businessID") == 0 {
		c.Set("businessID", businessID)
	}

	userID = c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return businessID, userID, false
	}

	u, err := gf.Model("business_user").
		Fields("id,status,can_manage_locks").
		Where("id", userID).
		Where("business_id", businessID).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("校验权限失败：" + err.Error()).Regin(c)
		return businessID, userID, false
	}
	if u == nil || len(u) == 0 {
		gf.Failed().SetCode(403).SetMsg("暂无权限").Regin(c)
		return businessID, userID, false
	}
	if u["status"].Int() == 1 {
		gf.Failed().SetCode(403).SetMsg("账号已被禁用").Regin(c)
		return businessID, userID, false
	}
	if u["can_manage_locks"].Int() != 1 {
		gf.Failed().SetCode(403).SetMsg("暂无智能锁管理权限").Regin(c)
		return businessID, userID, false
	}
	return businessID, userID, true
}

// -------------------------
// 小程序端：智能锁管理接口（参考后台 /houses/ttlock）
// -------------------------

// GetStatus 获取 TTLock 配置状态（小程序端）
func (api *WxTtlock) GetStatus(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).GetStatus(c)
}

// SyncLocks 同步锁列表（小程序端）
func (api *WxTtlock) SyncLocks(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).SyncLocks(c)
}

// GetLockList 获取锁列表（本地库，小程序端）
func (api *WxTtlock) GetLockList(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).GetLockList(c)
}

// BindProperty 绑定房源到某把锁（小程序端）
func (api *WxTtlock) BindProperty(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).BindProperty(c)
}

// UnbindProperty 解绑房源锁（小程序端）
func (api *WxTtlock) UnbindProperty(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).UnbindProperty(c)
}

// GetLockDetail 获取锁详情（小程序端）
func (api *WxTtlock) GetLockDetail(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).GetLockDetail(c)
}

// GetLockRecords 获取开锁记录（小程序端）
func (api *WxTtlock) GetLockRecords(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).GetLockRecords(c)
}

// GetPropertyLock 查询房源绑定的锁（小程序端）
func (api *WxTtlock) GetPropertyLock(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).GetPropertyLock(c)
}

// RemoteUnlock 远程开锁（小程序端）
func (api *WxTtlock) RemoteUnlock(c *gf.GinCtx) {
	if _, _, ok := wxEnsureCanManageLocks(c); !ok {
		return
	}
	(&housesPkg.Ttlock{}).RemoteUnlock(c)
}
