package houses

import (
	"context"
	"encoding/json"
	"errors"
	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
	"gofly/utils/tools/gtime"
	"strings"
)

// 房源管理
type Houses struct{ NoNeedAuths []string }

func init() {
	// 默认都需要登录与权限校验（由 business/controller.go 的 RouterHandler 处理）
	fpath := Houses{NoNeedAuths: []string{}}
	gf.Register(&fpath, fpath)
}

func normalizeCommaText(v interface{}) string {
	s := strings.TrimSpace(gconv.String(v))
	if s == "" {
		return ""
	}
	parts := gf.SplitAndStr(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}

func normalizeTagsToString(v interface{}) string {
	if v == nil {
		return ""
	}
	// 先处理数组类型
	switch val := v.(type) {
	case []interface{}:
		parts := make([]string, 0, len(val))
		for _, it := range val {
			s := strings.TrimSpace(gconv.String(it))
			if s != "" {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, ",")
	case []string:
		parts := make([]string, 0, len(val))
		for _, it := range val {
			s := strings.TrimSpace(it)
			if s != "" {
				parts = append(parts, s)
			}
		}
		return strings.Join(parts, ",")
	}
	// 字符串类型处理
	s := gconv.String(v)
	// 如果是 JSON 数组字符串，解析它
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		var arr []string
		if err := json.Unmarshal([]byte(s), &arr); err == nil {
			parts := make([]string, 0, len(arr))
			for _, it := range arr {
				it = strings.TrimSpace(it)
				if it != "" {
					parts = append(parts, it)
				}
			}
			return strings.Join(parts, ",")
		}
	}
	// 普通逗号分隔字符串
	return normalizeCommaText(s)
}

func normalizeRenovationStageStatus(v string) string {
	s := strings.TrimSpace(strings.ToLower(v))
	switch s {
	case "done", "finished", "completed":
		return "done"
	case "doing", "in_progress", "progress", "processing":
		return "doing"
	case "todo", "pending", "not_started", "none", "":
		return "todo"
	default:
		// 兼容中文输入
		if strings.Contains(v, "完成") || strings.Contains(v, "已完") {
			return "done"
		}
		if strings.Contains(v, "进行") || strings.Contains(v, "施工") || strings.Contains(v, "处理中") {
			return "doing"
		}
		if strings.Contains(v, "未") {
			return "todo"
		}
		return "todo"
	}
}

func normalizeRenovationStageLogList(arr []map[string]interface{}) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(arr))
	for _, it := range arr {
		if it == nil {
			continue
		}

		stage := strings.TrimSpace(gconv.String(it["stage"]))
		if stage == "" {
			stage = strings.TrimSpace(gconv.String(it["stage_name"]))
		}
		if stage == "" {
			stage = strings.TrimSpace(gconv.String(it["name"]))
		}
		if stage == "" {
			continue
		}

		status := normalizeRenovationStageStatus(gconv.String(it["status"]))
		date := strings.TrimSpace(gconv.String(it["date"]))
		note := strings.TrimSpace(gconv.String(it["note"]))
		if note == "" {
			note = strings.TrimSpace(gconv.String(it["notes"]))
		}

		images := ""
		if imgV, ok := it["images"]; ok {
			images = normalizeCommaText(normalizeTagsToString(imgV))
		}

		out = append(out, map[string]interface{}{
			"stage":  stage,
			"status": status,
			"date":   date,
			"note":   note,
			"images": images,
		})
	}
	return out
}

func normalizeRenovationStageLogsToJSON(v interface{}) string {
	if v == nil {
		return ""
	}

	// 1) 已是 JSON 字符串
	if s, ok := v.(string); ok {
		raw := strings.TrimSpace(s)
		if raw == "" {
			return ""
		}
		if strings.HasPrefix(raw, "[") {
			var arr []map[string]interface{}
			if err := json.Unmarshal([]byte(raw), &arr); err == nil {
				clean := normalizeRenovationStageLogList(arr)
				if len(clean) == 0 {
					return ""
				}
				b, err := json.Marshal(clean)
				if err == nil {
					return string(b)
				}
			}
		}
		// 非合法 JSON：不保存，避免污染数据
		return ""
	}

	// 2) 数组结构
	arr := make([]map[string]interface{}, 0)
	switch val := v.(type) {
	case []interface{}:
		for _, it := range val {
			if m, ok := it.(map[string]interface{}); ok {
				arr = append(arr, m)
			}
		}
	case []map[string]interface{}:
		arr = val
	default:
		return ""
	}

	clean := normalizeRenovationStageLogList(arr)
	if len(clean) == 0 {
		return ""
	}
	b, err := json.Marshal(clean)
	if err != nil {
		return ""
	}
	return string(b)
}

func parseRenovationStageLogs(raw string) []gf.Map {
	s := strings.TrimSpace(raw)
	if s == "" {
		return make([]gf.Map, 0)
	}
	var arr []map[string]interface{}
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		return make([]gf.Map, 0)
	}
	clean := normalizeRenovationStageLogList(arr)
	out := make([]gf.Map, 0, len(clean))
	for _, it := range clean {
		out = append(out, it)
	}
	return out
}

func isUnknownColumnErr(err error, col string) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Unknown column") && strings.Contains(msg, col)
}

func pickMap(param map[string]interface{}, keys ...string) gf.Map {
	out := gf.Map{}
	for _, k := range keys {
		if v, ok := param[k]; ok {
			out[k] = v
		}
	}
	return out
}

// 获取房源列表
func (api *Houses) GetList(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "12"))
	param, _ := gf.RequestParam(c)
	whereMap := gmap.New()
	whereMap.Set("business_id", c.GetInt64("businessID"))
	if title, ok := param["title"]; ok && title != "" {
		whereMap.Set("title like ?", "%"+gconv.String(title)+"%")
	}
	if community_name, ok := param["community_name"]; ok && community_name != "" {
		whereMap.Set("community_name like ?", "%"+gconv.String(community_name)+"%")
	}
	if sale_status, ok := param["sale_status"]; ok && sale_status != "" {
		whereMap.Set("sale_status", sale_status)
	}
	if property_type, ok := param["property_type"]; ok && property_type != "" {
		whereMap.Set("property_type", property_type)
	}
	if v, ok := param["hot_status"]; ok && gconv.String(v) != "" {
		hs := gconv.Int(v)
		if hs != 0 && hs != 1 {
			gf.Failed().SetMsg("hot_status参数不合法").Regin(c)
			return
		}
		whereMap.Set("hot_status", hs)
	}
	if status, ok := param["status"]; ok && status != "" {
		whereMap.Set("status", status)
	}
	if agentID, ok := param["agent_id"]; ok && agentID != "" && gconv.Int64(agentID) != 0 {
		whereMap.Set("agent_id", agentID)
	}
	MDB := gf.Model("business_properties").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,business_id,title,price,price_unit,area,rooms,halls,bathrooms,floor_level,total_floors,orientation,build_year,property_type,decoration_type,community_name,address,latitude,longitude,tags,images,cover_image,video_url,allow_image_download,allow_video_download,has_smart_lock,commission_rate,commission_reward,owner_name,owner_phone,receiver_name,receiver_phone,receiver_price,agent_id,sale_status,hot_status,view_count,follow_count,showing_count,status,weigh,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("weigh desc, id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
	} else {
		gf.Success().SetMsg("获取房源列表").SetData(gf.Map{
			"page":     pageNo,
			"pageSize": pageSize,
			"total":    totalCount,
			"items":    list,
		}).Regin(c)
	}
}

// 获取房源详情
func (api *Houses) GetContent(c *gf.GinCtx) {
	id := c.DefaultQuery("id", "")
	if id == "" {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	data, err := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).Where("id", id).Find()
	if err != nil {
		gf.Failed().SetMsg("获取内容失败").SetData(err).Regin(c)
	} else {
		// 行政区域兼容：address 按 “省 市 区 详细地址...” 拼接存储，这里拆出 district + address_detail 供详情页展示
		if data != nil && data["address"] != nil {
			addr := strings.TrimSpace(data["address"].String())
			if addr != "" {
				parts := strings.Fields(addr)
				if len(parts) >= 3 {
					data["district"] = gf.VarNew(strings.Join(parts[:3], " / "))
					if len(parts) >= 4 {
						data["address_detail"] = gf.VarNew(strings.Join(parts[3:], " "))
					} else {
						data["address_detail"] = gf.VarNew("")
					}
				} else {
					data["district"] = gf.VarNew("")
					data["address_detail"] = gf.VarNew(addr)
				}
			} else {
				data["district"] = gf.VarNew("")
				data["address_detail"] = gf.VarNew("")
			}
		}
		// tags 兼容：后端落库为逗号分隔字符串，前端表单用数组
		if data != nil && data["tags"] != nil && data["tags"].String() != "" {
			data["tags"] = gf.VarNew(gf.SplitAndStr(data["tags"].String(), ","))
		} else if data != nil {
			data["tags"] = gf.VarNew(make([]string, 0))
		}
		gf.Success().SetMsg("获取内容成功").SetData(data).Regin(c)
	}
}

// 保存房源
func (api *Houses) Save(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	var f_id = gf.GetEditId(param["id"])
	// 仅允许写入的字段（避免前端传入任意字段更新）
	saveData := pickMap(param,
		"title", "price", "price_unit", "area",
		"rooms", "halls", "bathrooms",
		"floor_level", "total_floors", "orientation", "build_year",
		"property_type", "decoration_type",
		"community_name", "address", "latitude", "longitude",
		"tags", "images", "cover_image", "video_url",
		"allow_image_download", "allow_video_download",
		"has_smart_lock",
		"commission_rate", "commission_reward",
		"owner_name", "owner_phone",
		"receiver_name", "receiver_phone", "receiver_price",
		"sale_status",
		"hot_status",
		"status",
		"weigh",
	)
	if _, ok := saveData["tags"]; ok {
		saveData["tags"] = normalizeTagsToString(saveData["tags"])
	}
	if _, ok := saveData["images"]; ok {
		saveData["images"] = normalizeCommaText(saveData["images"])
	}
	if _, ok := saveData["cover_image"]; ok {
		saveData["cover_image"] = strings.TrimSpace(gconv.String(saveData["cover_image"]))
	}
	if _, ok := saveData["video_url"]; ok {
		saveData["video_url"] = strings.TrimSpace(gconv.String(saveData["video_url"]))
	}
	if _, ok := saveData["allow_image_download"]; ok {
		aid := gconv.Int(saveData["allow_image_download"])
		if aid != 0 && aid != 1 {
			gf.Failed().SetMsg("allow_image_download参数不合法").Regin(c)
			return
		}
		saveData["allow_image_download"] = aid
	}
	if _, ok := saveData["allow_video_download"]; ok {
		avd := gconv.Int(saveData["allow_video_download"])
		if avd != 0 && avd != 1 {
			gf.Failed().SetMsg("allow_video_download参数不合法").Regin(c)
			return
		}
		saveData["allow_video_download"] = avd
	}
	if _, ok := saveData["hot_status"]; ok {
		hs := gconv.Int(saveData["hot_status"])
		if hs != 0 && hs != 1 {
			gf.Failed().SetMsg("hot_status参数不合法").Regin(c)
			return
		}
		saveData["hot_status"] = hs
	}
	if f_id == 0 {
		saveData["business_id"] = c.GetInt64("businessID")
		addId, err := gf.Model("business_properties").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
		} else {
			if addId != 0 {
				gf.Model("business_properties").Where("id", addId).Update(gf.Map{"weigh": addId})
			}
			gf.Success().SetMsg("添加成功").SetData(addId).Regin(c)
		}
	} else {
		_, err := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).Where("id", f_id).Update(saveData)
		if err != nil {
			gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新成功").Regin(c)
		}
	}
}

// 更新状态
func (api *Houses) UpStatus(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	propertyID := gconv.Int64(param["id"])
	if propertyID == 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	// 仅允许更新指定字段，避免“任意字段更新”
	update := gf.Map{}
	if v, ok := param["status"]; ok && v != "" {
		update["status"] = gconv.Int(v)
	}
	if v, ok := param["hot_status"]; ok && gconv.String(v) != "" {
		hs := gconv.Int(v)
		if hs != 0 && hs != 1 {
			gf.Failed().SetMsg("hot_status参数不合法").Regin(c)
			return
		}
		update["hot_status"] = hs
	}
	if v, ok := param["weigh"]; ok && v != "" {
		update["weigh"] = gconv.Int(v)
	}
	if v, ok := param["sale_status"]; ok && v != "" {
		ss := gconv.String(v)
		if ss != "on_sale" && ss != "sold" && ss != "off_market" {
			gf.Failed().SetMsg("sale_status参数不合法").Regin(c)
			return
		}
		update["sale_status"] = ss
	}
	if len(update) == 0 {
		gf.Failed().SetMsg("暂无可更新字段").Regin(c)
		return
	}
	businessID := c.GetInt64("businessID")
	userID := c.GetInt64("userID")
	remark := gconv.String(param["remark"])

	reqCtx := context.Background()
	if c != nil && c.Request != nil {
		reqCtx = c.Request.Context()
	}

	err := gf.Model("business_properties").Ctx(reqCtx).Transaction(reqCtx, func(ctx context.Context, tx gform.TX) error {
		old, err := gf.Model("business_properties").TX(tx).Ctx(ctx).
			Fields("status,hot_status,weigh,sale_status").
			Where("business_id", businessID).
			Where("id", propertyID).
			Find()
		if err != nil {
			return err
		}
		if old == nil || len(old) == 0 {
			return errors.New("房源不存在或无权限")
		}

		if _, err := gf.Model("business_properties").TX(tx).Ctx(ctx).
			Where("business_id", businessID).
			Where("id", propertyID).
			Update(update); err != nil {
			return err
		}

		now := gtime.Now().Format("Y-m-d H:i:s")
		for field, after := range update {
			beforeVar, ok := old[field]
			var before string
			if ok && beforeVar != nil {
				before = beforeVar.String()
			}
			afterStr := gconv.String(after)
			if before == afterStr {
				continue
			}
			logRow := gf.Map{
				"business_id":  businessID,
				"property_id":  propertyID,
				"user_id":      userID,
				"field":        gconv.String(field),
				"before_value": before,
				"after_value":  afterStr,
				"remark":       remark,
				"createtime":   now,
			}
			if _, err := gf.Model("business_property_status_logs").TX(tx).Ctx(ctx).Data(logRow).Insert(); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("更新成功").Regin(c)
}

// 获取房源状态变更记录
func (api *Houses) GetStatusLogs(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "12"))
	propertyID := gconv.Int64(c.DefaultQuery("property_id", "0"))
	if propertyID == 0 {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}
	businessID := c.GetInt64("businessID")

	// 校验房源归属
	existProperty, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("deletetime", nil).
		Exist()
	if !existProperty {
		gf.Failed().SetMsg("房源不存在或无权限").Regin(c)
		return
	}

	MDB := gf.Model("business_property_status_logs l").
		LeftJoin("business_user u", "u.id = l.user_id").
		Where("l.business_id", businessID).
		Where("l.property_id", propertyID)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.
		Fields("l.id,l.property_id,l.user_id,l.field,l.before_value,l.after_value,l.remark,l.createtime,u.name as user_name,u.username as user_username").
		Page(pageNo, pageSize).
		Order("l.id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg(err.Error()).Regin(c)
		return
	}
	gf.Success().SetMsg("获取状态变更记录").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    list,
	}).Regin(c)
}

// 获取房源的经纪人行为记录（浏览/带看/开锁）
// 入参：
// - property_id: 房源ID（必填）
// - record_type: view|showing|unlock（必填）
// - page/pageSize
func (api *Houses) GetBehaviorLogs(c *gf.GinCtx) {
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "12"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 12
	}
	if pageSize > 50 {
		pageSize = 50
	}

	propertyID := gconv.Int64(c.DefaultQuery("property_id", "0"))
	if propertyID <= 0 {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}
	recordType := strings.TrimSpace(c.DefaultQuery("record_type", ""))
	if recordType == "" {
		recordType = strings.TrimSpace(c.DefaultQuery("type", ""))
	}
	businessID := c.GetInt64("businessID")

	// 校验房源归属
	existProperty, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("deletetime", nil).
		Exist()
	if !existProperty {
		gf.Failed().SetMsg("房源不存在或无权限").Regin(c)
		return
	}

	totalCount := 0
	list := make(gform.Result, 0)
	switch recordType {
	case "view", "showing":
		MDB := gf.Model("business_user_activity_logs a").
			LeftJoin("business_user u", "u.id = a.user_id").
			LeftJoin("business_stores s", "s.id = u.store_id AND s.deletetime IS NULL").
			Where("a.property_id", propertyID).
			Where("a.activity_type", recordType).
			Where("u.business_id", businessID)
		totalCount, _ = MDB.Clone().Count()
		rows, err := MDB.
			Fields("a.id,a.user_id,a.property_id,a.activity_type,a.meta_data,a.createtime,u.name as user_name,u.username as user_username,u.mobile as user_mobile,u.store_id,u.title as user_title,s.name as store_name").
			Page(pageNo, pageSize).
			Order("a.id desc").
			Select()
		if err != nil {
			gf.Failed().SetMsg("获取记录失败：" + err.Error()).Regin(c)
			return
		}
		list = rows
	case "unlock":
		MDB := gf.Model("business_unlock_requests r").
			LeftJoin("business_user u", "u.id = r.user_id").
			LeftJoin("business_stores s", "s.id = u.store_id AND s.deletetime IS NULL").
			Where("r.property_id", propertyID).
			Where("u.business_id", businessID)
		totalCount, _ = MDB.Clone().Count()
		rows, err := MDB.
			Fields("r.id,r.user_id,r.property_id,r.request_status,r.request_type,r.request_time,r.expires_at,r.approval_remark,r.status,r.createtime,r.updatetime,u.name as user_name,u.username as user_username,u.mobile as user_mobile,u.store_id,u.title as user_title,s.name as store_name").
			Page(pageNo, pageSize).
			Order("r.id desc").
			Select()
		if err != nil {
			gf.Failed().SetMsg("获取记录失败：" + err.Error()).Regin(c)
			return
		}
		list = rows
	default:
		gf.Failed().SetMsg("record_type 参数无效").Regin(c)
		return
	}

	gf.Success().SetMsg("获取记录成功").SetData(gf.Map{
		"record_type": recordType,
		"page":        pageNo,
		"pageSize":    pageSize,
		"total":       totalCount,
		"items":       list,
	}).Regin(c)
}

// 删除房源
func (api *Houses) Del(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	if param["ids"] == nil {
		gf.Failed().SetMsg("请传参数ids").Regin(c)
		return
	}
	_, err := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).WhereIn("id", param["ids"]).Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("删除成功").Regin(c)
	}
}

// 获取装修信息
func (api *Houses) GetRenovation(c *gf.GinCtx) {
	propertyId := c.DefaultQuery("property_id", "")
	if propertyId == "" {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}
	// 先校验房源是否属于当前商户
	property, _ := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).Where("id", propertyId).Find()
	if property == nil || len(property) == 0 {
		gf.Failed().SetMsg("房源不存在或无权限").Regin(c)
		return
	}
	data, err := gf.Model("business_renovations").Where("property_id", propertyId).Find()
	if err != nil {
		gf.Failed().SetMsg("获取装修信息失败").SetData(err).Regin(c)
		return
	}
	// 处理 materials 字段：后端逗号分隔，前端数组
	if data != nil && data["materials"] != nil && data["materials"].String() != "" {
		data["materials"] = gf.VarNew(gf.SplitAndStr(data["materials"].String(), ","))
	} else if data != nil {
		data["materials"] = gf.VarNew(make([]string, 0))
	}
	// 处理 images 字段：保持字符串（兼容业务端表单组件），但做一次去重/去空格清洗
	if data != nil && data["images"] != nil && strings.TrimSpace(data["images"].String()) != "" {
		data["images"] = gf.VarNew(normalizeCommaText(data["images"].String()))
	} else if data != nil {
		data["images"] = gf.VarNew("")
	}
	// 处理 stage_logs 字段：JSON 数组（用于“工序时间线”）
	if data != nil {
		raw := ""
		if data["stage_logs"] != nil {
			raw = strings.TrimSpace(data["stage_logs"].String())
		}
		if raw != "" {
			data["stage_logs"] = gf.VarNew(parseRenovationStageLogs(raw))
		} else {
			data["stage_logs"] = gf.VarNew(make([]gf.Map, 0))
		}
	}
	gf.Success().SetMsg("获取装修信息成功").SetData(data).Regin(c)
}

// 保存装修信息
func (api *Houses) SaveRenovation(c *gf.GinCtx) {
	param, _ := gf.RequestParam(c)
	propertyId := gconv.Int64(param["property_id"])
	if propertyId == 0 {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}
	// 校验房源是否属于当前商户
	property, _ := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).Where("id", propertyId).Find()
	if property == nil || len(property) == 0 {
		gf.Failed().SetMsg("房源不存在或无权限").Regin(c)
		return
	}
	// 仅允许写入的字段
	saveData := pickMap(param,
		"renovation_status", "progress_percentage", "current_stage",
		"start_date", "estimated_finish_date", "actual_finish_date",
		"materials", "images", "stage_logs", "notes", "status",
	)
	// 处理 materials 字段（数组转逗号分隔字符串）
	if _, ok := saveData["materials"]; ok {
		saveData["materials"] = normalizeTagsToString(saveData["materials"])
	}
	// 处理 images 字段
	if _, ok := saveData["images"]; ok {
		saveData["images"] = normalizeCommaText(saveData["images"])
	}
	// 处理 stage_logs 字段（JSON 数组字符串）
	if _, ok := saveData["stage_logs"]; ok {
		saveData["stage_logs"] = normalizeRenovationStageLogsToJSON(saveData["stage_logs"])
	}
	// 查询是否已存在装修记录
	existing, _ := gf.Model("business_renovations").Where("property_id", propertyId).Find()
	if existing == nil || len(existing) == 0 {
		// 新增
		saveData["property_id"] = propertyId
		_, err := gf.Model("business_renovations").Data(saveData).InsertAndGetId()
		if err != nil && isUnknownColumnErr(err, "stage_logs") {
			delete(saveData, "stage_logs")
			_, err = gf.Model("business_renovations").Data(saveData).InsertAndGetId()
		}
		if err != nil {
			gf.Failed().SetMsg("添加装修信息失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("添加装修信息成功").Regin(c)
		}
	} else {
		// 更新
		_, err := gf.Model("business_renovations").Where("property_id", propertyId).Update(saveData)
		if err != nil && isUnknownColumnErr(err, "stage_logs") {
			delete(saveData, "stage_logs")
			_, err = gf.Model("business_renovations").Where("property_id", propertyId).Update(saveData)
		}
		if err != nil {
			gf.Failed().SetMsg("更新装修信息失败").SetData(err).Regin(c)
		} else {
			gf.Success().SetMsg("更新装修信息成功").Regin(c)
		}
	}
}
