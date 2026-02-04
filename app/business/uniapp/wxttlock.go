package uniapp

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	housesPkg "gofly/app/business/houses"
	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// WxTtlock 小程序端 TTLock 数据接口（离线蓝牙插件使用）
// - 避免 RBAC：NoNeedAuths = ["*"]
// - 仍建议要求登录：接口内会校验 userID
type WxTtlock struct{ NoNeedAuths []string }

func init() {
	fpath := WxTtlock{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

func wxParsePropertyID(c *gf.GinCtx, param map[string]any, key string) int64 {
	idStr := strings.TrimSpace(c.Query(key))
	if idStr == "" {
		idStr = strings.TrimSpace(gconv.String(param[key]))
	}
	return gconv.Int64(idStr)
}

func wxEnsurePropertyOwned(c *gf.GinCtx, businessID, propertyID int64) bool {
	if propertyID <= 0 {
		gf.Failed().SetMsg("缺少参数 property_id").Regin(c)
		return false
	}
	existProperty, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("deletetime", nil).
		Where("status", 0).
		Exist()
	if !existProperty {
		gf.Failed().SetMsg("房源不存在或已下架").Regin(c)
		return false
	}
	return true
}

func wxInsertUserUnlockActivity(userID, propertyID int64, meta any) {
	_, _ = gf.Model("business_user_activity_logs").Data(gf.Map{
		"user_id":       userID,
		"property_id":   propertyID,
		"activity_type": "unlock",
		"meta_data":     meta,
		"createtime":    time.Now(),
	}).Insert()
}

// GetPropertyBleInfo 获取房源绑定锁的离线蓝牙插件所需数据
// 入参：property_id
// 返回：lock_id/lock_mac/lock_data + key_data（sdkVersion=2）
func (api *WxTtlock) GetPropertyBleInfo(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetMsg("请先登录").Regin(c)
		return
	}

	param, _ := gf.RequestParam(c)
	propertyID := wxParsePropertyID(c, param, "property_id")
	if !wxEnsurePropertyOwned(c, businessID, propertyID) {
		return
	}

	// 获取绑定锁
	bind, _ := gf.Model("business_property_locks").
		Fields("ttlock_lock_id,bind_status").
		Where("business_id", businessID).
		Where("property_id", propertyID).
		Where("deletetime", 0).
		Order("id desc").
		Find()
	if bind == nil || bind["ttlock_lock_id"].Int64() == 0 {
		gf.Failed().SetMsg("该房源未绑定智能锁").Regin(c)
		return
	}
	if bind["bind_status"].Int() == 0 {
		gf.Failed().SetMsg("该房源智能锁已解绑").Regin(c)
		return
	}
	lockID := bind["ttlock_lock_id"].Int64()

	lockRow, _ := gf.Model("business_smart_locks").
		Fields("ttlock_lock_id,lock_name,lock_mac,battery,model_num,raw_json,last_sync_at").
		Where("business_id", businessID).
		Where("ttlock_lock_id", lockID).
		Where("deletetime", 0).
		Find()
	if lockRow == nil || lockRow["ttlock_lock_id"].Int64() == 0 {
		gf.Failed().SetMsg("锁不存在，请先在后台同步锁列表").Regin(c)
		return
	}

	lockData := ""
	raw := strings.TrimSpace(lockRow["raw_json"].String())
	if raw != "" {
		var m map[string]any
		if e := json.Unmarshal([]byte(raw), &m); e == nil && m != nil {
			if v, ok := m["lockData"]; ok {
				lockData = strings.TrimSpace(gconv.String(v))
			}
		}
	}

	keyInfo, keyID, keyData, keyErr := housesPkg.TTLockGetKeyDataForLock(c, businessID, lockID, 2)
	if keyInfo == nil {
		keyInfo = map[string]any{}
	}

	housesPkg.TTLockInsertEvent(businessID, userID, propertyID, lockID, "ble_info", keyErr == nil, 0, func() string {
		if keyErr != nil {
			return keyErr.Error()
		}
		return ""
	}(), gf.Map{
		"has_lock_data": lockData != "",
		"has_key_data":  strings.TrimSpace(keyData) != "",
	})

	if keyErr != nil {
		gf.Failed().SetMsg("获取钥匙数据失败：" + keyErr.Error()).Regin(c)
		return
	}

	gf.Success().SetMsg("ok").SetData(gf.Map{
		"property_id":     propertyID,
		"ttlock_lock_id":  lockID,
		"lock_name":       lockRow["lock_name"].String(),
		"lock_mac":        lockRow["lock_mac"].String(),
		"battery":         lockRow["battery"].Int(),
		"model_num":       lockRow["model_num"].String(),
		"last_sync_at":    lockRow["last_sync_at"].Int64(),
		"lock_data":       lockData,
		"key_id":          keyID,
		"key_data":        keyData,
		"key_status":      gf.Int(keyInfo["keyStatus"]),
		"key_sender_name": gconv.String(keyInfo["senderUsername"]),
	}).Regin(c)
}

// LogUnlockActivity 小程序端：记录开锁/蓝牙过程的关键节点日志
//
// 入参：
// - property_id: 房源ID（必填）
// - stage: 阶段标识（必填，如 ble_init_start/ble_init_ok/ble_init_fail/ble_state_change/ble_conn_change）
// - meta: 扩展信息（可选，JSON对象/字符串）
// - page: 页面标识（可选）
func (api *WxTtlock) LogUnlockActivity(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	param, _ := gf.RequestParam(c)
	propertyID := wxParsePropertyID(c, param, "property_id")
	if !wxEnsurePropertyOwned(c, businessID, propertyID) {
		return
	}

	stage := strings.TrimSpace(gconv.String(param["stage"]))
	if stage == "" {
		gf.Failed().SetMsg("缺少参数 stage").Regin(c)
		return
	}

	meta := gf.Map{
		"stage": stage,
		"page":  strings.TrimSpace(gconv.String(param["page"])),
	}
	if v, ok := param["meta"]; ok && v != nil {
		meta["meta"] = v
	}

	wxInsertUserUnlockActivity(userID, propertyID, meta)
	gf.Success().SetMsg("ok").Regin(c)
}

// BleUnlockStart 小程序端：创建一次“蓝牙开锁”请求记录
//
// 入参：
// - property_id: 房源ID（必填）
// - meta: 扩展信息（可选）
// 返回：request_id
func (api *WxTtlock) BleUnlockStart(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	param, _ := gf.RequestParam(c)
	propertyID := wxParsePropertyID(c, param, "property_id")
	if !wxEnsurePropertyOwned(c, businessID, propertyID) {
		return
	}

	now := time.Now()
	requestID, err := gf.Model("business_unlock_requests").Data(gf.Map{
		"user_id":        userID,
		"property_id":    propertyID,
		"request_status": "pending",
		"request_type":   "bluetooth",
		"request_time":   now,
		"status":         0,
		"createtime":     now,
		"updatetime":     now,
	}).InsertAndGetId()
	if err != nil || requestID <= 0 {
		gf.Failed().SetMsg("创建开锁请求失败").SetExdata(err).Regin(c)
		return
	}

	meta := gf.Map{
		"stage":      "ble_unlock_start",
		"request_id": requestID,
	}
	if v, ok := param["meta"]; ok && v != nil {
		meta["meta"] = v
	}
	wxInsertUserUnlockActivity(userID, propertyID, meta)

	gf.Success().SetMsg("ok").SetData(gf.Map{"request_id": requestID}).Regin(c)
}

func wxLimitLen(s string, max int) string {
	raw := strings.TrimSpace(s)
	if raw == "" {
		return ""
	}
	rs := []rune(raw)
	if len(rs) <= max {
		return raw
	}
	return string(rs[:max])
}

// BleUnlockFinish 小程序端：完成一次“蓝牙开锁”请求（成功/失败）
//
// 入参：
// - request_id: BleUnlockStart 返回的ID（必填）
// - property_id: 房源ID（必填）
// - success: 1/0 或 true/false（必填）
// - err_code: 错误码（可选）
// - err_msg: 错误信息（可选）
// - meta: 扩展信息（可选）
func (api *WxTtlock) BleUnlockFinish(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	param, _ := gf.RequestParam(c)
	propertyID := wxParsePropertyID(c, param, "property_id")
	if !wxEnsurePropertyOwned(c, businessID, propertyID) {
		return
	}

	requestID := gconv.Int64(strings.TrimSpace(gconv.String(param["request_id"])))
	if requestID <= 0 {
		gf.Failed().SetMsg("缺少参数 request_id").Regin(c)
		return
	}

	success := false
	switch v := param["success"].(type) {
	case bool:
		success = v
	default:
		success = gf.Int(v) == 1
	}
	errCode := gf.Int(param["err_code"])
	errMsg := wxLimitLen(gconv.String(param["err_msg"]), 255)

	now := time.Now()
	requestStatus := "completed"
	remark := "开锁成功"
	if !success {
		requestStatus = "rejected"
		if strings.TrimSpace(errMsg) == "" {
			errMsg = "开锁失败"
		}
		remark = errMsg
	}

	_, _ = gf.Model("business_unlock_requests").
		Where("id", requestID).
		Where("user_id", userID).
		Where("property_id", propertyID).
		Data(gf.Map{
			"request_status":  requestStatus,
			"approval_remark": remark,
			"updatetime":      now,
		}).
		Update()

	meta := gf.Map{
		"stage":      "ble_unlock_finish",
		"request_id": requestID,
		"success":    success,
		"err_code":   errCode,
		"err_msg":    errMsg,
	}
	if v, ok := param["meta"]; ok && v != nil {
		meta["meta"] = v
	}
	wxInsertUserUnlockActivity(userID, propertyID, meta)

	gf.Success().SetMsg("ok").Regin(c)
}

// Passcode 生成房源开锁密码（一次性密码/临时密码）
//
// 入参：property_id
// 备注：此接口面向小程序端“看房经纪人”，按你确认的权限方案A：只要求已登录 + 房源归属校验，不额外要求 can_manage_locks。
func (api *WxTtlock) Passcode(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	// 兼容：后续内部会复用 houses.Ttlock.RemoteUnlock，该方法会再次读取 request body。
	// 因此这里先把 body 读出来并回填，避免出现 “缺少参数 property_id/ttlock_lock_id”。
	body, _ := io.ReadAll(c.Request.Body)
	if len(bytes.TrimSpace(body)) == 0 {
		body = []byte("{}")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	param := make(map[string]any)
	_ = json.Unmarshal(body, &param)

	propertyID := wxParsePropertyID(c, param, "property_id")
	if !wxEnsurePropertyOwned(c, businessID, propertyID) {
		return
	}

	// 兼容：某些链路 businessID 未写入 context（Ttlock.RemoteUnlock 内部取 c.GetInt64("businessID")）
	if businessID > 0 && c.GetInt64("businessID") == 0 {
		c.Set("businessID", businessID)
	}

	// 复用后台现有逻辑：RemoteUnlock 实际为“获取键盘开锁密码”
	(&housesPkg.Ttlock{}).RemoteUnlock(c)
}
