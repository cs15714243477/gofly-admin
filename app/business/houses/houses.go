package houses

import (
	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
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
	// 兼容前端传数组（a-input-tag）或字符串
	if strings.HasPrefix(gconv.String(v), "[") && strings.HasSuffix(gconv.String(v), "]") {
		// 若是 JSON 字符串数组，交给 gf 处理可能失败，这里走最保守的字符串落库
		return normalizeCommaText(v)
	}
	switch v.(type) {
	case []interface{}:
		return gf.ArrayToStr(v, ",")
	case []string:
		arr := make([]interface{}, 0)
		for _, it := range v.([]string) {
			arr = append(arr, it)
		}
		return gf.ArrayToStr(arr, ",")
	default:
		return normalizeCommaText(v)
	}
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
	if status, ok := param["status"]; ok && status != "" {
		whereMap.Set("status", status)
	}
	if agentID, ok := param["agent_id"]; ok && agentID != "" && gconv.Int64(agentID) != 0 {
		whereMap.Set("agent_id", agentID)
	}
	MDB := gf.Model("business_properties").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	list, err := MDB.Fields("id,business_id,title,price,price_unit,area,rooms,halls,bathrooms,floor_level,total_floors,orientation,build_year,property_type,decoration_type,community_name,address,latitude,longitude,tags,images,cover_image,has_smart_lock,commission_rate,commission_reward,owner_phone,agent_id,sale_status,view_count,follow_count,showing_count,status,weigh,createtime,updatetime").
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
		"tags", "images", "cover_image",
		"has_smart_lock",
		"commission_rate", "commission_reward",
		"owner_phone", "agent_id",
		"sale_status",
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
	if param["id"] == nil || gconv.Int64(param["id"]) == 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}
	// 仅允许更新指定字段，避免“任意字段更新”
	update := gf.Map{}
	if v, ok := param["status"]; ok && v != "" {
		update["status"] = gconv.Int(v)
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
	_, err := gf.Model("business_properties").Where("business_id", c.GetInt64("businessID")).Where("id", param["id"]).Update(update)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
	} else {
		gf.Success().SetMsg("更新成功").Regin(c)
	}
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
