package houses

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gvar"
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
			"GetStatus":         "GET:/business/houses/ttlock/status",
			"GetLockList":       "GET:/business/houses/ttlock/lockList",
			"GetPropertyLock":   "GET:/business/houses/ttlock/propertyLock",
			"GetLockDetail":     "GET:/business/houses/ttlock/lockDetail",
			"GetLockRecords":    "GET:/business/houses/ttlock/lockRecords",
			"GetPasscodeEvents": "GET:/business/houses/ttlock/passcodeEvents",
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
		Fields("id,ttlock_lock_id,lock_name,lock_mac,battery,model_num,raw_json,last_sync_at,status,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("id desc").
		Select()

	// Enrich list with bind info (one lock -> one property).
	if len(list) > 0 {
		lockIDs := make([]interface{}, 0, len(list))
		for _, r := range list {
			if r == nil {
				continue
			}
			if lockID := r["ttlock_lock_id"].Int64(); lockID > 0 {
				lockIDs = append(lockIDs, lockID)
			}
		}

		if len(lockIDs) > 0 {
			binds, _ := gf.Model("business_property_locks").
				Fields("id,property_id,ttlock_lock_id,bind_status,bind_time,unbind_time").
				Where("business_id", businessID).
				Where("bind_status", 1).
				Where("deletetime", 0).
				WhereIn("ttlock_lock_id", lockIDs).
				Order("id desc").
				Select()

			bindByLockID := make(map[int64]gform.Record, len(binds))
			propertyIDs := make([]interface{}, 0, len(binds))
			for _, b := range binds {
				if b == nil {
					continue
				}
				lockID := b["ttlock_lock_id"].Int64()
				if lockID <= 0 {
					continue
				}
				// Since ordered by id desc, the first one per lock is the latest binding.
				if _, ok := bindByLockID[lockID]; ok {
					continue
				}
				bindByLockID[lockID] = b
				if pid := b["property_id"].Int64(); pid > 0 {
					propertyIDs = append(propertyIDs, pid)
				}
			}

			propertyByID := make(map[int64]gform.Record, len(propertyIDs))
			if len(propertyIDs) > 0 {
				props, _ := gf.Model("business_properties").
					Fields("id,title,community_name").
					Where("business_id", businessID).
					Where("deletetime", nil).
					WhereIn("id", propertyIDs).
					Select()
				for _, p := range props {
					if p == nil {
						continue
					}
					if id := p["id"].Int64(); id > 0 {
						propertyByID[id] = p
					}
				}
			}

			for _, r := range list {
				if r == nil {
					continue
				}
				lockID := r["ttlock_lock_id"].Int64()
				if lockID <= 0 {
					continue
				}
				b, ok := bindByLockID[lockID]
				if !ok || b == nil || b["property_id"].Int64() <= 0 {
					continue
				}
				pid := b["property_id"].Int64()
				r["bind_property_id"] = gvar.New(pid)
				r["bind_time"] = gvar.New(b["bind_time"].Int64())

				if p, ok := propertyByID[pid]; ok && p != nil {
					r["bind_property_title"] = gvar.New(p["title"].String())
					r["bind_property_community_name"] = gvar.New(p["community_name"].String())
				}
			}
		}
	}

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
	reqCtx := context.Background()
	if c != nil && c.Request != nil {
		reqCtx = c.Request.Context()
	}

	var unboundOthers bool
	err := gf.Model("business_property_locks").Ctx(reqCtx).Transaction(reqCtx, func(ctx context.Context, tx gform.TX) error {
		// 1) 解绑该锁在其他房源上的绑定（确保“一锁一房”）
		otherBinds, err := gf.Model("business_property_locks").TX(tx).Ctx(ctx).
			Fields("id,property_id").
			Where("business_id", businessID).
			Where("ttlock_lock_id", lockID).
			Where("bind_status", 1).
			Where("deletetime", 0).
			WhereNot("property_id", propertyID).
			Select()
		if err != nil {
			return err
		}
		if len(otherBinds) > 0 {
			unboundOthers = true
			ids := make([]interface{}, 0, len(otherBinds))
			pids := make([]interface{}, 0, len(otherBinds))
			for _, r := range otherBinds {
				if r == nil {
					continue
				}
				if id := r["id"].Int64(); id > 0 {
					ids = append(ids, id)
				}
				if pid := r["property_id"].Int64(); pid > 0 {
					pids = append(pids, pid)
				}
			}
			if len(ids) > 0 {
				if _, e := gf.Model("business_property_locks").TX(tx).Ctx(ctx).
					WhereIn("id", ids).
					Update(gf.Map{"bind_status": 0, "unbind_time": now, "updatetime": now}); e != nil {
					return e
				}
			}
			// 同步冗余字段（列表展示用）：其他房源解除绑定
			if len(pids) > 0 {
				_, _ = gf.Model("business_properties").TX(tx).Ctx(ctx).
					Where("business_id", businessID).
					WhereIn("id", pids).
					Update(gf.Map{"has_smart_lock": 0, "updatetime": now})
			}
		}

		// 2) upsert：一房一锁（business_id + property_id 唯一）
		existRow, err := gf.Model("business_property_locks").TX(tx).Ctx(ctx).
			Fields("id").
			Where("business_id", businessID).
			Where("property_id", propertyID).
			Where("deletetime", 0).
			Find()
		if err != nil {
			return err
		}
		if existRow != nil && existRow["id"].Int64() > 0 {
			_, err = gf.Model("business_property_locks").TX(tx).Ctx(ctx).
				Where("id", existRow["id"].Int64()).
				Update(gf.Map{
					"ttlock_lock_id":  lockID,
					"bind_status":     1,
					"bind_by_user_id": userID,
					"bind_time":       now,
					"unbind_time":     nil,
					"updatetime":      now,
				})
			if err != nil {
				return err
			}
		} else {
			_, err = gf.Model("business_property_locks").TX(tx).Ctx(ctx).Data(gf.Map{
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
				return err
			}
		}

		// 3) 同步冗余字段（列表展示用）：当前房源绑定成功
		_, _ = gf.Model("business_properties").TX(tx).Ctx(ctx).
			Where("business_id", businessID).
			Where("id", propertyID).
			Update(gf.Map{"has_smart_lock": 1, "updatetime": now})
		return nil
	})
	if err != nil {
		gf.Failed().SetMsg("绑定失败").SetData(err).Regin(c)
		return
	}

	msg := "绑定成功"
	if unboundOthers {
		msg = "绑定成功（该锁已从其他房源解绑）"
	}
	gf.Success().SetMsg(msg).SetData(true).Regin(c)
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
	if bind == nil || bind["ttlock_lock_id"].Int64() == 0 || bind["bind_status"].Int() == 0 {
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

// RemoteUnlock 获取开锁密码（TTLock Cloud /v3/keyboardPwd/get）
// 入参：property_id 或 ttlock_lock_id（二选一）
// 可选：keyboard_pwd_type（默认 1 一次性）、keyboard_pwd_name、keyboard_pwd_version、start_date/end_date（秒或毫秒）
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

	keyboardPwdType := gf.Int(param["keyboard_pwd_type"])
	keyboardPwdName := strings.TrimSpace(gf.String(param["keyboard_pwd_name"]))
	keyboardPwdVersion := gf.Int(param["keyboard_pwd_version"])
	start := gf.Int64(param["start_date"])
	end := gf.Int64(param["end_date"])
	if keyboardPwdType <= 0 {
		// 1 一次性密码（6小时内可用一次）
		keyboardPwdType = 1
	}
	// 兼容秒/毫秒
	if start > 0 && start < 2000000000 {
		start = start * 1000
	}
	if end > 0 && end < 2000000000 {
		end = end * 1000
	}
	// 默认取整点（TTLock 要求按小时定义）
	if start <= 0 && keyboardPwdType == 1 {
		start = time.Now().UTC().Truncate(time.Hour).UnixMilli()
	}
	if start > 0 {
		start = time.UnixMilli(start).UTC().Truncate(time.Hour).UnixMilli()
	}
	if end > 0 {
		end = time.UnixMilli(end).UTC().Truncate(time.Hour).UnixMilli()
	}
	// 一次性密码：默认有效期 6 小时（用于前端展示）
	if end <= 0 && keyboardPwdType == 1 && start > 0 {
		end = start + 6*60*60*1000
	}

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}

	// 优先从本地锁 raw_json 解析 keyboardPwdVersion
	if keyboardPwdVersion <= 0 {
		lockRow, _ := gf.Model("business_smart_locks").
			Fields("raw_json").
			Where("business_id", businessID).
			Where("ttlock_lock_id", lockID).
			Where("deletetime", 0).
			Find()
		if lockRow != nil {
			raw := strings.TrimSpace(lockRow["raw_json"].String())
			if raw != "" {
				var m map[string]any
				if e := json.Unmarshal([]byte(raw), &m); e == nil && m != nil {
					keyboardPwdVersion = gconv.Int(m["keyboardPwdVersion"])
				}
			}
		}
	}
	// 兜底：从云端详情取一次（避免未同步导致缺少版本）
	if keyboardPwdVersion <= 0 {
		if detail, e := ttlockLockDetail(cfg, accessToken, lockID); e == nil && detail != nil {
			keyboardPwdVersion = gconv.Int(detail["keyboardPwdVersion"])
		}
	}

	out, callErr := ttlockGetKeyboardPwd(cfg, accessToken, lockID, keyboardPwdVersion, keyboardPwdType, keyboardPwdName, start, end)
	errCode, errMsg := ttlockParseErr(out)
	ok := callErr == nil
	if !ok && errMsg == "" {
		errMsg = callErr.Error()
	}
	ttlockInsertEvent(businessID, userID, propertyID, lockID, "get_passcode", ok, errCode, errMsg, gf.Map{
		"keyboard_pwd_type":    keyboardPwdType,
		"keyboard_pwd_version": keyboardPwdVersion,
		"keyboard_pwd_name":    keyboardPwdName,
		"start_date":           start,
		"end_date":             end,
		"response":             out,
	})
	if callErr != nil {
		gf.Failed().SetMsg("获取密码失败：" + errMsg).SetData(out).Regin(c)
		return
	}

	// 补充有效期信息（毫秒时间戳，前端按本地时区展示）
	if out == nil {
		out = map[string]any{}
	}
	if start > 0 {
		out["startDate"] = start
	}
	if end > 0 {
		out["endDate"] = end
	}
	out["keyboardPwdType"] = keyboardPwdType

	gf.Success().SetMsg("开锁密码已生成").SetData(out).Regin(c)
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

// GetLockRecords 获取锁开锁/操作记录（TTLock Cloud /v3/lockRecord/list）
// 入参：property_id 或 ttlock_lock_id（二选一）
// 可选：start_date/end_date（秒或毫秒），page/pageSize
// 默认：最近 7 天
func (api *Ttlock) GetLockRecords(c *gf.GinCtx) {
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

	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "20"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	start := gf.Int64(c.DefaultQuery("start_date", "0"))
	end := gf.Int64(c.DefaultQuery("end_date", "0"))
	// 兼容秒/毫秒
	if start > 0 && start < 2000000000 {
		start = start * 1000
	}
	if end > 0 && end < 2000000000 {
		end = end * 1000
	}
	// 默认最近 7 天
	if start <= 0 && end <= 0 {
		end = time.Now().UnixMilli()
		start = end - 7*24*60*60*1000
	}

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}
	out, callErr := ttlockLockRecordList(cfg, accessToken, lockID, start, end, pageNo, pageSize)
	if callErr != nil {
		gf.Failed().SetMsg("获取开锁记录失败：" + callErr.Error()).SetData(out).Regin(c)
		return
	}
	gf.Success().SetMsg("开锁记录").SetData(out).Regin(c)
}

// GetPasscodeEvents 获取开锁密码申请记录（本地业务库 business_lock_events）
// 入参：property_id 或 ttlock_lock_id（二选一）
// 可选：page/pageSize
func (api *Ttlock) GetPasscodeEvents(c *gf.GinCtx) {
	businessID := c.GetInt64("businessID")
	propertyID := gf.Int64(c.DefaultQuery("property_id", "0"))
	lockID := gf.Int64(c.DefaultQuery("ttlock_lock_id", "0"))
	if lockID == 0 && propertyID > 0 {
		if id, err := ttlockLockIDByProperty(businessID, propertyID); err == nil {
			lockID = id
		} else {
			// 房源未绑定锁：返回空列表即可
			gf.Success().SetMsg("密码申请记录").SetData(gf.Map{
				"page":     1,
				"pageSize": 20,
				"total":    0,
				"items":    []any{},
			}).Regin(c)
			return
		}
	}
	if lockID == 0 {
		gf.Failed().SetMsg("缺少参数 property_id/ttlock_lock_id").Regin(c)
		return
	}

	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "20"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	MDB := gf.Model("business_lock_events").
		Where("business_id", businessID).
		Where("ttlock_lock_id", lockID).
		Where("event_type", "get_passcode")
	if propertyID > 0 {
		MDB = MDB.Where("property_id", propertyID)
	}
	totalCount, _ := MDB.Clone().Count()
	list, _ := MDB.
		Fields("id,business_id,user_id,property_id,ttlock_lock_id,event_type,result,err_code,err_msg,meta_json,createtime").
		Page(pageNo, pageSize).
		Order("id desc").
		Select()

	// attach property title
	propertyIDs := make([]interface{}, 0, len(list))
	userIDs := make([]interface{}, 0, len(list))
	for _, r := range list {
		if r == nil {
			continue
		}
		if pid := r["property_id"].Int64(); pid > 0 {
			propertyIDs = append(propertyIDs, pid)
		}
		if uid := r["user_id"].Int64(); uid > 0 {
			userIDs = append(userIDs, uid)
		}
	}
	propertyByID := make(map[int64]gform.Record, len(propertyIDs))
	if len(propertyIDs) > 0 {
		props, _ := gf.Model("business_properties").
			Fields("id,title,community_name").
			Where("business_id", businessID).
			Where("deletetime", nil).
			WhereIn("id", propertyIDs).
			Select()
		for _, p := range props {
			if p == nil {
				continue
			}
			if id := p["id"].Int64(); id > 0 {
				propertyByID[id] = p
			}
		}
	}

	userByID := make(map[int64]gform.Record, len(userIDs))
	if len(userIDs) > 0 {
		users, _ := gf.Model("business_user").
			Fields("id,username,name,nickname").
			Where("business_id", businessID).
			Where("deletetime", nil).
			WhereIn("id", userIDs).
			Select()
		for _, u := range users {
			if u == nil {
				continue
			}
			if id := u["id"].Int64(); id > 0 {
				userByID[id] = u
			}
		}
	}

	// parse meta_json -> useful fields
	for _, r := range list {
		if r == nil {
			continue
		}
		pid := r["property_id"].Int64()
		if p, ok := propertyByID[pid]; ok && p != nil {
			r["property_title"] = gvar.New(p["title"].String())
			r["property_community_name"] = gvar.New(p["community_name"].String())
		}
		uid := r["user_id"].Int64()
		if u, ok := userByID[uid]; ok && u != nil {
			username := strings.TrimSpace(u["username"].String())
			name := strings.TrimSpace(u["name"].String())
			nickname := strings.TrimSpace(u["nickname"].String())
			display := nickname
			if display == "" {
				display = name
			}
			if display == "" {
				display = username
			}
			if display != "" {
				r["apply_user_name"] = gvar.New(display)
			}
			if username != "" {
				r["apply_user_username"] = gvar.New(username)
			}
		}

		metaStr := strings.TrimSpace(r["meta_json"].String())
		if metaStr == "" {
			continue
		}
		var meta map[string]any
		if e := json.Unmarshal([]byte(metaStr), &meta); e != nil || meta == nil {
			continue
		}

		var resp map[string]any
		if v, ok := meta["response"]; ok {
			if mm, ok2 := v.(map[string]any); ok2 {
				resp = mm
			}
		}
		keyboardPwd := ""
		keyboardPwdID := int64(0)
		startDate := int64(0)
		endDate := int64(0)
		keyboardPwdType := gconv.Int(meta["keyboard_pwd_type"])
		if keyboardPwdType <= 0 {
			keyboardPwdType = 1
		}

		if resp != nil {
			keyboardPwd = strings.TrimSpace(gconv.String(resp["keyboardPwd"]))
			if keyboardPwd == "" {
				keyboardPwd = strings.TrimSpace(gconv.String(resp["keyboard_pwd"]))
			}
			keyboardPwdID = gconv.Int64(resp["keyboardPwdId"])
			startDate = gconv.Int64(resp["startDate"])
			endDate = gconv.Int64(resp["endDate"])
		}
		if startDate <= 0 {
			startDate = gconv.Int64(meta["start_date"])
		}
		if endDate <= 0 {
			endDate = gconv.Int64(meta["end_date"])
		}

		if keyboardPwd != "" {
			r["keyboard_pwd"] = gvar.New(keyboardPwd)
		}
		if keyboardPwdID > 0 {
			r["keyboard_pwd_id"] = gvar.New(keyboardPwdID)
		}
		if startDate > 0 {
			r["start_date"] = gvar.New(startDate)
		}
		if endDate > 0 {
			r["end_date"] = gvar.New(endDate)
		}
		r["keyboard_pwd_type"] = gvar.New(keyboardPwdType)
	}

	gf.Success().SetMsg("密码申请记录").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list,
	}).Regin(c)
}
