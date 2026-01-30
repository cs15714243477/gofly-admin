package houses

import (
	"encoding/json"
	"strings"

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
	idStr := strings.TrimSpace(c.Query("property_id"))
	if idStr == "" {
		idStr = strings.TrimSpace(gconv.String(param["property_id"]))
	}
	propertyID := gconv.Int64(idStr)
	if propertyID <= 0 {
		gf.Failed().SetMsg("缺少参数 property_id").Regin(c)
		return
	}

	// 校验房源归属
	existProperty, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("deletetime", nil).
		Where("status", 0).
		Exist()
	if !existProperty {
		gf.Failed().SetMsg("房源不存在或已下架").Regin(c)
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

	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		gf.Failed().SetMsg("获取 access_token 失败：" + err.Error()).Regin(c)
		return
	}

	// key/list + sdkVersion=2：用于小程序蓝牙插件（WxTTLockPlugin）
	keyInfo, keyErr := ttlockFindKeyForLock(cfg, accessToken, lockID, 2)
	var keyData string
	var keyID int64
	if keyInfo != nil {
		keyID = gf.Int64(keyInfo["keyId"])
		keyData = strings.TrimSpace(gconv.String(keyInfo["lockData"]))
	}

	ok := keyErr == nil
	errCode, errMsg := ttlockParseErr(nil)
	if keyErr != nil {
		// 这里没有云端返回体，先把错误信息写入日志，避免前端卡死
		errMsg = keyErr.Error()
	}
	ttlockInsertEvent(businessID, userID, propertyID, lockID, "ble_info", ok, errCode, errMsg, gf.Map{
		"has_lock_data": lockData != "",
		"has_key_data":  keyData != "",
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
