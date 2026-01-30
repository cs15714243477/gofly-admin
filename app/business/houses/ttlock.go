package houses

import (
	"encoding/json"
	"errors"
	"time"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// 通通锁（TTLock）后台管理接口
type Ttlock struct {
	NoNeedAuths  []string
	CustomRoutes gf.Map
}

func init() {
	fpath := Ttlock{
		NoNeedAuths: []string{},
		CustomRoutes: gf.Map{
			// 兼容前端更直观的短路径（避免 GET 方法名自动带 get 前缀）
			"GetStatus":       "GET:/business/houses/ttlock/status",
			"GetLockList":     "GET:/business/houses/ttlock/lockList",
			"GetPropertyLock": "GET:/business/houses/ttlock/propertyLock",
			"GetLockDetail":   "GET:/business/houses/ttlock/lockDetail",
		},
	}
	gf.Register(&fpath, fpath)
}

// 获取配置/状态
func (api *Ttlock) GetStatus(c *gf.GinCtx) {
	cfg := getTTLockCloudConfig(c)
	businessID := c.GetInt64("businessID")
	acc, _ := ttlockLoadAccount(businessID)
	ok := cfg.ClientID != "" && cfg.ClientSecret != "" && cfg.Username != "" && cfg.Password != ""
	gf.Success().SetMsg("TTLock 状态").SetData(gf.Map{
		"config_ok": ok,
		"client_id": cfg.ClientID,
		"api_base":  cfg.APIBase,
		"has_token": acc != nil && acc.AccessToken != "",
		"expire_at": func() int64 {
			if acc == nil {
				return 0
			}
			return acc.ExpireAt
		}(),
	}).Regin(c)
}

// SyncLocks 同步锁列表（从 TTLock Cloud 拉取并落库）
func (api *Ttlock) SyncLocks(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}

	pageNo := 1
	pageSize := 50
	totalSaved := 0
	for pageNo <= 10 { // 防止无限循环，最多拉 10 页
		out, err := ttlockListLocks(cfg, accessToken, pageNo, pageSize)
		if err != nil {
			gf.Failed().SetMsg("同步失败：" + err.Error()).SetData(out).Regin(c)
			return
		}
		if len(out.List) == 0 {
			break
		}

		for _, it := range out.List {
			lockID := gconv.Int64(it["lockId"])
			if lockID == 0 {
				continue
			}
			lockName := gconv.String(it["lockName"])
			lockMac := gconv.String(it["lockMac"])
			battery := gconv.Int(it["battery"])
			modelNum := gconv.String(it["modelNum"])
			raw, _ := json.Marshal(it)

			now := time.Now().Unix()
			existID, _ := gf.Model("business_smart_locks").
				Where("business_id", businessID).
				Where("ttlock_lock_id", lockID).
				Where("deletetime", 0).
				Value("id")

			data := gf.Map{
				"business_id":    businessID,
				"ttlock_lock_id": lockID,
				"lock_name":      lockName,
				"lock_mac":       lockMac,
				"battery":        battery,
				"model_num":      modelNum,
				"raw_json":       string(raw),
				"last_sync_at":   now,
				"updatetime":     now,
				"status":         0,
				"deletetime":     0,
			}

			if existID == nil {
				data["createtime"] = now
				_, _ = gf.Model("business_smart_locks").Data(data).Insert()
			} else {
				_, _ = gf.Model("business_smart_locks").Where("id", existID).Update(data)
			}
			totalSaved++
		}

		// 已经拉完
		if out.Total > 0 && pageNo*pageSize >= out.Total {
			break
		}
		pageNo++
	}

	gf.Success().SetMsg("同步完成").SetData(gf.Map{"saved": totalSaved}).Regin(c)
}

// 锁列表（本地库）
func (api *Ttlock) GetLockList(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "20"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	MDB := gf.Model("business_smart_locks").
		Where("business_id", businessID).
		Where("deletetime", 0)
	totalCount, _ := MDB.Clone().Count()
	list, _ := MDB.
		Fields("id,ttlock_lock_id,lock_name,lock_mac,battery,model_num,last_sync_at,status,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("id desc").
		Select()

	gf.Success().SetMsg("锁列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list,
	}).Regin(c)
}

// BindProperty 绑定房源到某把锁（默认一房一锁）
func (api *Ttlock) BindProperty(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	businessID := c.GetInt64("businessID")
	userID := c.GetInt64("userID")
	propertyID := gf.Int64(param["property_id"])
	lockID := gf.Int64(param["ttlock_lock_id"])
	if propertyID == 0 || lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}

	// 校验房源归属
	existProperty, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("deletetime", nil).
		Exist()
	if !existProperty {
		gf.Failed().SetMsg("房源不存在").Regin(c)
		return
	}

	// 校验锁归属（必须先同步）
	existLock, _ := gf.Model("business_smart_locks").
		Where("business_id", businessID).
		Where("ttlock_lock_id", lockID).
		Where("deletetime", 0).
		Exist()
	if !existLock {
		gf.Failed().SetMsg("锁不存在，请先同步锁列表").Regin(c)
		return
	}

	now := time.Now().Unix()
	// 解绑旧绑定（确保一房一锁）
	oldBindID, _ := gf.Model("business_property_locks").
		Where("business_id", businessID).
		Where("property_id", propertyID).
		Where("deletetime", 0).
		Value("id")
	if oldBindID != nil {
		_, _ = gf.Model("business_property_locks").Where("id", oldBindID).Update(gf.Map{
			"bind_status": 0,
			"unbind_time": now,
			"updatetime":  now,
		})
	}

	// 插入新绑定
	_, err := gf.Model("business_property_locks").Data(gf.Map{
		"business_id":     businessID,
		"property_id":     propertyID,
		"ttlock_lock_id":  lockID,
		"bind_status":     1,
		"bind_by_user_id": userID,
		"bind_time":       now,
		"createtime":      now,
		"updatetime":      now,
		"deletetime":      0,
	}).Insert()
	if err != nil {
		gf.Failed().SetMsg("绑定失败").SetData(err).Regin(c)
		return
	}

	// 同步冗余字段（列表展示用）
	_, _ = gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Update(gf.Map{"has_smart_lock": 1, "updatetime": now})

	gf.Success().SetMsg("绑定成功").SetData(true).Regin(c)
}

// UnbindProperty 解绑房源锁
func (api *Ttlock) UnbindProperty(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	businessID := c.GetInt64("businessID")
	propertyID := gf.Int64(param["property_id"])
	if propertyID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id").Regin(c)
		return
	}
	now := time.Now().Unix()
	_, err := gf.Model("business_property_locks").
		Where("business_id", businessID).
		Where("property_id", propertyID).
		Where("deletetime", 0).
		Update(gf.Map{"bind_status": 0, "unbind_time": now, "updatetime": now})
	if err != nil {
		gf.Failed().SetMsg("解绑失败").SetData(err).Regin(c)
		return
	}
	_, _ = gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Update(gf.Map{"has_smart_lock": 0, "updatetime": now})
	gf.Success().SetMsg("解绑成功").SetData(true).Regin(c)
}

// 查询房源绑定的锁
func (api *Ttlock) GetPropertyLock(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	propertyID := gf.Int64(c.DefaultQuery("property_id", "0"))
	if propertyID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id").Regin(c)
		return
	}

	bind, _ := gf.Model("business_property_locks").
		Fields("id,property_id,ttlock_lock_id,bind_status,bind_time,unbind_time").
		Where("business_id", businessID).
		Where("property_id", propertyID).
		Where("deletetime", 0).
		Order("id desc").
		Find()
	if bind == nil {
		gf.Success().SetMsg("未绑定").SetData(nil).Regin(c)
		return
	}
	lock, _ := gf.Model("business_smart_locks").
		Fields("ttlock_lock_id,lock_name,lock_mac,battery,model_num,last_sync_at").
		Where("business_id", businessID).
		Where("ttlock_lock_id", bind["ttlock_lock_id"].Int64()).
		Where("deletetime", 0).
		Find()

	gf.Success().SetMsg("房源锁信息").SetData(gf.Map{
		"bind": bind,
		"lock": lock,
	}).Regin(c)
}

func ttlockLockIDByProperty(businessID, propertyID int64) (int64, error) {
	if businessID == 0 || propertyID == 0 {
		return 0, errors.New("缺少 business_id/property_id")
	}
	bind, _ := gf.Model("business_property_locks").
		Fields("ttlock_lock_id,bind_status").
		Where("business_id", businessID).
		Where("property_id", propertyID).
		Where("deletetime", 0).
		Order("id desc").
		Find()
	if bind == nil || len(bind) == 0 || bind["ttlock_lock_id"].Int64() == 0 {
		return 0, errors.New("该房源未绑定智能锁")
	}
	if bind["bind_status"].Int() == 0 {
		return 0, errors.New("该房源智能锁已解绑")
	}
	return bind["ttlock_lock_id"].Int64(), nil
}

// RemoteUnlock 远程开锁（TTLock Cloud /v3/lock/unlock）
// 入参：property_id 或 ttlock_lock_id（二选一）
func (api *Ttlock) RemoteUnlock(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	businessID := c.GetInt64("businessID")
	userID := c.GetInt64("userID")
	propertyID := gf.Int64(param["property_id"])
	lockID := gf.Int64(param["ttlock_lock_id"])
	if lockID == 0 && propertyID > 0 {
		if id, err := ttlockLockIDByProperty(businessID, propertyID); err == nil {
			lockID = id
		} else {
			gf.Failed().SetMsg(err.Error()).Regin(c)
			return
		}
	}
	if lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}
	out, callErr := ttlockUnlock(cfg, accessToken, lockID)
	errCode, errMsg := ttlockParseErr(out)
	ok := callErr == nil
	if !ok && errMsg == "" {
		errMsg = callErr.Error()
	}
	ttlockInsertEvent(businessID, userID, propertyID, lockID, "remote_unlock", ok, errCode, errMsg, gf.Map{
		"response": out,
	})
	if callErr != nil {
		gf.Failed().SetMsg("开锁失败：" + errMsg).SetData(out).Regin(c)
		return
	}
	gf.Success().SetMsg("开锁指令已下发").SetData(out).Regin(c)
}

// AddKeyboardPwd 通过网关下发临时密码（TTLock Cloud /v3/keyboardPwd/add, addType=2）
// 入参：property_id 或 ttlock_lock_id；keyboard_pwd（必填）；keyboard_pwd_name/start_date/end_date（可选，毫秒或秒）
func (api *Ttlock) AddKeyboardPwd(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	businessID := c.GetInt64("businessID")
	userID := c.GetInt64("userID")
	propertyID := gf.Int64(param["property_id"])
	lockID := gf.Int64(param["ttlock_lock_id"])
	if lockID == 0 && propertyID > 0 {
		if id, err := ttlockLockIDByProperty(businessID, propertyID); err == nil {
			lockID = id
		} else {
			gf.Failed().SetMsg(err.Error()).Regin(c)
			return
		}
	}
	if lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}
	keyboardPwd := gf.String(param["keyboard_pwd"])
	keyboardPwdName := gf.String(param["keyboard_pwd_name"])
	start := gf.Int64(param["start_date"])
	end := gf.Int64(param["end_date"])
	// 兼容秒/毫秒
	if start > 0 && start < 2000000000 {
		start = start * 1000
	}
	if end > 0 && end < 2000000000 {
		end = end * 1000
	}

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}
	out, callErr := ttlockAddKeyboardPwdGateway(cfg, accessToken, lockID, keyboardPwd, keyboardPwdName, start, end)
	errCode, errMsg := ttlockParseErr(out)
	ok := callErr == nil
	if !ok && errMsg == "" {
		errMsg = callErr.Error()
	}
	ttlockInsertEvent(businessID, userID, propertyID, lockID, "add_passcode", ok, errCode, errMsg, gf.Map{
		"keyboard_pwd_name": keyboardPwdName,
		"start_date":        start,
		"end_date":          end,
		"response":          out,
	})
	if callErr != nil {
		gf.Failed().SetMsg("下发密码失败：" + errMsg).SetData(out).Regin(c)
		return
	}
	gf.Success().SetMsg("密码已下发").SetData(out).Regin(c)
}

// SendEkey 下发电子钥匙（TTLock Cloud /v3/key/send）
// 入参：property_id 或 ttlock_lock_id；receiver_username（必填）；key_name/remarks/start_date/end_date（可选，毫秒或秒）
func (api *Ttlock) SendEkey(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	businessID := c.GetInt64("businessID")
	userID := c.GetInt64("userID")
	propertyID := gf.Int64(param["property_id"])
	lockID := gf.Int64(param["ttlock_lock_id"])
	if lockID == 0 && propertyID > 0 {
		if id, err := ttlockLockIDByProperty(businessID, propertyID); err == nil {
			lockID = id
		} else {
			gf.Failed().SetMsg(err.Error()).Regin(c)
			return
		}
	}
	if lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}
	receiver := gf.String(param["receiver_username"])
	keyName := gf.String(param["key_name"])
	remarks := gf.String(param["remarks"])
	start := gf.Int64(param["start_date"])
	end := gf.Int64(param["end_date"])
	if start > 0 && start < 2000000000 {
		start = start * 1000
	}
	if end > 0 && end < 2000000000 {
		end = end * 1000
	}

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}
	out, callErr := ttlockSendEkey(cfg, accessToken, lockID, receiver, keyName, remarks, start, end)
	errCode, errMsg := ttlockParseErr(out)
	ok := callErr == nil
	if !ok && errMsg == "" {
		errMsg = callErr.Error()
	}
	ttlockInsertEvent(businessID, userID, propertyID, lockID, "send_ekey", ok, errCode, errMsg, gf.Map{
		"receiver_username": receiver,
		"key_name":          keyName,
		"start_date":        start,
		"end_date":          end,
		"response":          out,
	})
	if callErr != nil {
		gf.Failed().SetMsg("下发钥匙失败：" + errMsg).SetData(out).Regin(c)
		return
	}
	gf.Success().SetMsg("钥匙已下发").SetData(out).Regin(c)
}

// GetLockDetail 获取云端锁详情（TTLock Cloud /v3/lock/detail）
// 入参：property_id 或 ttlock_lock_id（二选一）
func (api *Ttlock) GetLockDetail(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	propertyID := gf.Int64(c.DefaultQuery("property_id", "0"))
	lockID := gf.Int64(c.DefaultQuery("ttlock_lock_id", "0"))
	if lockID == 0 && propertyID > 0 {
		if id, err := ttlockLockIDByProperty(businessID, propertyID); err == nil {
			lockID = id
		} else {
			gf.Failed().SetMsg(err.Error()).Regin(c)
			return
		}
	}
	if lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}
	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}
	out, callErr := ttlockLockDetail(cfg, accessToken, lockID)
	if callErr != nil {
		gf.Failed().SetMsg("获取锁详情失败：" + callErr.Error()).SetData(out).Regin(c)
		return
	}
	gf.Success().SetMsg("锁详情").SetData(out).Regin(c)
}
