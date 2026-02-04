package houses

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// TTLockGetKeyDataForLock 获取指定锁的蓝牙钥匙数据（用于小程序蓝牙插件）
//
// 说明：
// - 内部会自动获取 access_token，并调用 key/list（sdkVersion=2）获取钥匙数据
// - 返回 keyInfo（原始结构）、keyID、keyData
func TTLockGetKeyDataForLock(c *gf.GinCtx, businessID, lockID int64, sdkVersion int) (keyInfo map[string]any, keyID int64, keyData string, err error) {
	cfg := getTTLockCloudConfig(c)
	accessToken, err := ttlockEnsureAccessToken(c, businessID)
	if err != nil {
		return nil, 0, "", err
	}
	keyInfo, err = ttlockFindKeyForLock(cfg, accessToken, lockID, sdkVersion)
	if err != nil {
		return keyInfo, 0, "", err
	}
	if keyInfo != nil {
		keyID = gf.Int64(keyInfo["keyId"])
		keyData = strings.TrimSpace(gconv.String(keyInfo["lockData"]))
	}
	return keyInfo, keyID, keyData, nil
}

// TTLockInsertEvent 写入智能锁事件日志（business_lock_events）
//
// 备注：对外暴露是为了 uniapp 包的接口复用同一套日志入库逻辑。
func TTLockInsertEvent(businessID, userID, propertyID, lockID int64, eventType string, result bool, errCode int, errMsg string, meta any) {
	ttlockInsertEvent(businessID, userID, propertyID, lockID, eventType, result, errCode, errMsg, meta)
}
