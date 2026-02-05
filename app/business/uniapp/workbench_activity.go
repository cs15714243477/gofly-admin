package uniapp

import (
	"encoding/json"
	"strings"
	"time"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

func wxIsWorkbenchActivityType(raw string) bool {
	switch strings.TrimSpace(raw) {
	case wxRecordTypeShowing, wxRecordTypeView, wxRecordTypeShare, wxRecordTypeCall:
		return true
	default:
		return false
	}
}

func wxParamString(c *gf.GinCtx, param map[string]interface{}, key string) string {
	v := strings.TrimSpace(c.Query(key))
	if v != "" {
		return v
	}
	if param == nil {
		return ""
	}
	return strings.TrimSpace(gconv.String(param[key]))
}

func wxParamInt64(c *gf.GinCtx, param map[string]interface{}, key string) int64 {
	return gconv.Int64(wxParamString(c, param, key))
}

func wxMarshalMeta(meta map[string]any) string {
	if meta == nil || len(meta) == 0 {
		return ""
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return ""
	}
	return string(b)
}

func wxBuildActivityMeta(param map[string]interface{}) map[string]any {
	out := map[string]any{}
	if param == nil {
		return out
	}
	if v := strings.TrimSpace(gconv.String(param["page"])); v != "" {
		out["page"] = v
	}
	if v, ok := param["meta"]; ok && v != nil {
		out["meta"] = v
	}
	return out
}

// AddWorkbenchActivityLog 写入“带看/浏览/分享/通话”等活动日志（用于工作台记录页）
//
// 入参（query/body 均可）：
// - activity_type: showing|view|share|call（必填）
// - property_id: 房源ID（可选，不传则为0）
// - page: 页面标识（可选）
// - meta: 扩展信息（可选，JSON对象/字符串）
//
// 说明：
// - follow/unlock 走各自表/接口，不在这里写
// - 对 view 做 10 分钟去重：同一用户同一房源 10 分钟内只记 1 条
func (api *Index) AddWorkbenchActivityLog(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}
	businessID := wxBusinessID(c)

	param, _ := gf.RequestParam(c)

	activityType := wxParamString(c, param, "activity_type")
	if activityType == "" {
		activityType = wxParamString(c, param, "record_type")
	}
	if !wxIsWorkbenchActivityType(activityType) {
		gf.Failed().SetMsg("activity_type 参数无效").Regin(c)
		return
	}

	propertyID := wxParamInt64(c, param, "property_id")
	if propertyID > 0 {
		if !wxEnsurePropertyOwned(c, businessID, propertyID) {
			return
		}
	}

	// 去重：10分钟内同一用户同一房源浏览只记一次
	if activityType == wxRecordTypeView && propertyID > 0 {
		since := time.Now().Add(-10 * time.Minute)
		exists, _ := gf.Model("business_user_activity_logs").
			Where("user_id", userID).
			Where("property_id", propertyID).
			Where("activity_type", activityType).
			Where("createtime >= ?", since).
			Exist()
		if exists {
			gf.Success().SetMsg("ok").SetData(gf.Map{
				"skipped": true,
				"reason":  "10分钟内已记录",
			}).Regin(c)
			return
		}
	}

	meta := wxBuildActivityMeta(param)
	metaText := wxMarshalMeta(meta)

	now := time.Now()
	id, err := gf.Model("business_user_activity_logs").Data(gf.Map{
		"user_id":       userID,
		"property_id":   propertyID,
		"activity_type": activityType,
		"meta_data":     metaText,
		"createtime":    now,
	}).InsertAndGetId()
	if err != nil || id <= 0 {
		gf.Failed().SetMsg("记录失败").SetExdata(err).Regin(c)
		return
	}

	gf.Success().SetMsg("记录成功").SetData(gf.Map{
		"id":           id,
		"activity_type": activityType,
		"property_id":  propertyID,
	}).Regin(c)
}

