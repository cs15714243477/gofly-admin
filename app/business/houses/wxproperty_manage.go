package houses

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
)

// wxParseOptionsText 解析配置项文本为下拉选项
//
// 支持格式（每行一条）：
// - value=label
// - value:label
// - value,label
// - value（则 value=label=value）
func wxParseOptionsText(v interface{}) []gf.Map {
	raw := strings.TrimSpace(gconv.String(v))
	if raw == "" {
		return nil
	}
	// 兼容：用户可能写成字面量 \n
	raw = strings.ReplaceAll(raw, "\\n", "\n")
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	lines := strings.Split(raw, "\n")
	out := make([]gf.Map, 0, len(lines))
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		// 行内注释
		if idx := strings.Index(l, "#"); idx >= 0 {
			l = strings.TrimSpace(l[:idx])
		}
		if l == "" {
			continue
		}

		value := ""
		label := ""
		switch {
		case strings.Contains(l, "="):
			parts := strings.SplitN(l, "=", 2)
			value = strings.TrimSpace(parts[0])
			label = strings.TrimSpace(parts[1])
		case strings.Contains(l, ":"):
			parts := strings.SplitN(l, ":", 2)
			value = strings.TrimSpace(parts[0])
			label = strings.TrimSpace(parts[1])
		case strings.Contains(l, ","):
			parts := strings.SplitN(l, ",", 2)
			value = strings.TrimSpace(parts[0])
			label = strings.TrimSpace(parts[1])
		default:
			value = strings.TrimSpace(l)
			label = value
		}

		if value == "" {
			continue
		}
		if label == "" {
			label = value
		}
		out = append(out, gf.Map{"label": label, "value": value})
	}
	return out
}

// wxPickDictOptionsByTitles 根据字典分组标题读取下拉选项
//
// 字典存储：common_dictionary_group / common_dictionary_data
// - label: keyname
// - value: keyvalue
//
// @return options, ok
func wxPickDictOptionsByTitles(businessID int64, titles []string) ([]gf.Map, bool) {
	if businessID <= 0 {
		businessID = 1
	}
	for _, title := range titles {
		t := strings.TrimSpace(title)
		if t == "" {
			continue
		}

		// 优先找 business 端配置（business_id 匹配）
		group, gerr := gf.Model("common_dictionary_group").
			Where("status", 0).
			Where("business_id", businessID).
			Where("title", t).
			Find()
		if (gerr != nil || group == nil || len(group) == 0) && businessID != 0 {
			// 兜底找公共配置
			group, _ = gf.Model("common_dictionary_group").
				Where("status", 0).
				Where("data_from", "common").
				Where("title", t).
				Find()
		}
		if group == nil || len(group) == 0 {
			continue
		}
		groupID := group["id"].Int64()
		if groupID <= 0 {
			continue
		}

		rows, err := gf.Model("common_dictionary_data").
			Where("group_id", groupID).
			Where("status", 0).
			Fields("keyname,keyvalue").
			Order("weigh asc,id asc").
			Select()
		if err != nil || len(rows) == 0 {
			continue
		}

		out := make([]gf.Map, 0, len(rows))
		for _, r := range rows {
			if r == nil || len(r) == 0 {
				continue
			}
			label := strings.TrimSpace(r["keyname"].String())
			value := strings.TrimSpace(r["keyvalue"].String())
			if label == "" || value == "" {
				continue
			}
			out = append(out, gf.Map{"label": label, "value": value})
		}
		if len(out) > 0 {
			return out, true
		}
	}
	return nil, false
}

// wxEnsureCanManageProperties 校验小程序端是否具备“可维护房源”权限
//
// 最小权限原则：仅允许维护 agent_id=当前用户 的房源
func wxEnsureCanManageProperties(c *gf.GinCtx) (businessID int64, userID int64, ok bool) {
	businessID = wxBusinessID(c)
	userID = c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return businessID, userID, false
	}

	u, err := gf.Model("business_user").
		Fields("id,status,can_manage_properties").
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
	if u["can_manage_properties"].Int() != 1 {
		gf.Failed().SetCode(403).SetMsg("暂无房源维护权限").Regin(c)
		return businessID, userID, false
	}
	return businessID, userID, true
}

// GetManageList 可维护房源列表（小程序端）
// 入参：page,pageSize,keyword,sale_status
func (api *WxProperty) GetManageList(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	param, _ := gf.RequestParam(c)
	keyword := strings.TrimSpace(gconv.String(param["keyword"]))
	if keyword == "" {
		keyword = strings.TrimSpace(c.DefaultQuery("keyword", ""))
	}
	saleStatus := strings.TrimSpace(gconv.String(param["sale_status"]))
	if saleStatus == "" {
		saleStatus = strings.TrimSpace(c.DefaultQuery("sale_status", ""))
	}

	whereMap := gmap.New()
	whereMap.Set("business_id", businessID)
	whereMap.Set("deletetime", nil)
	// 最小权限：仅自己维护的房源
	whereMap.Set("agent_id", userID)
	if keyword != "" {
		whereMap.Set("title like ?", "%"+keyword+"%")
	}
	if saleStatus != "" {
		whereMap.Set("sale_status", saleStatus)
	}

	MDB := gf.Model("business_properties").Where(whereMap)
	totalCount, _ := MDB.Clone().Count()
	rows, err := MDB.
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,orientation,community_name,address,tags,cover_image,images,video_url,allow_image_download,allow_video_download,sale_status,status,view_count,follow_count,showing_count,createtime,updatetime").
		Page(pageNo, pageSize).
		Order("id desc").
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取房源列表失败：" + err.Error()).Regin(c)
		return
	}

	items := make([]gf.Map, 0, len(rows))
	for _, r := range rows {
		if r == nil {
			continue
		}
		cover := wxFullImgURL(r["cover_image"].String())
		if cover == "" && r["images"] != nil && strings.TrimSpace(r["images"].String()) != "" {
			first := strings.Split(r["images"].String(), ",")[0]
			cover = wxFullImgURL(first)
		}
		ss := strings.TrimSpace(r["sale_status"].String())
		items = append(items, gf.Map{
			"id":                   r["id"].Int64(),
			"title":                r["title"].String(),
			"price":                gconv.String(r["price"]),
			"price_unit":           r["price_unit"].String(),
			"area":                 gconv.String(r["area"]),
			"rooms":                r["rooms"].Int(),
			"halls":                r["halls"].Int(),
			"bathrooms":            r["bathrooms"].Int(),
			"orientation":          r["orientation"].String(),
			"community_name":       r["community_name"].String(),
			"address":              r["address"].String(),
			"tags":                 wxSplitComma(r["tags"].String()),
			"image":                cover,
			"video_url":            wxFullImgURL(r["video_url"].String()),
			"allow_image_download": r["allow_image_download"].Int(),
			"allow_video_download": r["allow_video_download"].Int(),
			"sale_status":          ss,
			"sale_status_label":    wxSaleStatusLabel(ss),
			"status":               r["status"].Int(),
			"view_count":           r["view_count"].Int(),
			"follow_count":         r["follow_count"].Int(),
			"showing_count":        r["showing_count"].Int(),
			"createtime":           r["createtime"].String(),
			"updatetime":           r["updatetime"].String(),
		})
	}

	gf.Success().SetMsg("获取房源列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    totalCount,
		"items":    items,
	}).Regin(c)
}

// GetManageContent 获取可维护房源详情（编辑页回显）
// 入参：id
func (api *WxProperty) GetManageContent(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	param, _ := gf.RequestParam(c)
	idStr := strings.TrimSpace(c.Query("id"))
	if idStr == "" {
		idStr = strings.TrimSpace(gconv.String(param["id"]))
	}
	id := gconv.Int64(idStr)
	if id <= 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}

	row, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", id).
		Where("agent_id", userID).
		Where("deletetime", nil).
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取房源失败：" + err.Error()).Regin(c)
		return
	}
	if row == nil || len(row) == 0 {
		gf.Failed().SetCode(403).SetMsg("房源不存在或无权限").Regin(c)
		return
	}

	images := make([]string, 0)
	for _, p := range wxSplitComma(row["images"].String()) {
		if u := wxFullImgURL(p); u != "" {
			images = append(images, u)
		}
	}

	out := gf.Map{
		"id":                   row["id"].Int64(),
		"title":                row["title"].String(),
		"price":                gconv.String(row["price"]),
		"price_unit":           row["price_unit"].String(),
		"area":                 gconv.String(row["area"]),
		"rooms":                row["rooms"].Int(),
		"halls":                row["halls"].Int(),
		"bathrooms":            row["bathrooms"].Int(),
		"floor_level":          row["floor_level"].String(),
		"total_floors":         row["total_floors"].Int(),
		"orientation":          row["orientation"].String(),
		"build_year":           row["build_year"].Int(),
		"property_type":        row["property_type"].String(),
		"decoration_type":      row["decoration_type"].String(),
		"community_name":       row["community_name"].String(),
		"address":              row["address"].String(),
		"latitude":             gconv.Float64(row["latitude"]),
		"longitude":            gconv.Float64(row["longitude"]),
		"tags":                 wxSplitComma(row["tags"].String()),
		"images":               images,
		"cover_image":          wxFullImgURL(row["cover_image"].String()),
		"video_url":            wxFullImgURL(row["video_url"].String()),
		"allow_image_download": row["allow_image_download"].Int(),
		"allow_video_download": row["allow_video_download"].Int(),
		"owner_name":           row["owner_name"].String(),
		"owner_phone":          row["owner_phone"].String(),
		"receiver_name":        row["receiver_name"].String(),
		"receiver_phone":       row["receiver_phone"].String(),
		"receiver_price":       gconv.String(row["receiver_price"]),
		"sale_status":          row["sale_status"].String(),
		"commission_rate":      gconv.String(row["commission_rate"]),
		"commission_reward":    gconv.String(row["commission_reward"]),
		"weigh":                row["weigh"].Int64(),
		"hot_status":           row["hot_status"].Int(),
		"has_smart_lock":       row["has_smart_lock"].Int(),
		"status":               row["status"].Int(),
	}

	gf.Success().SetMsg("获取房源成功").SetData(out).Regin(c)
}

// GetManageRenovation 获取可维护房源装修信息（小程序端编辑页）
// 入参：property_id
func (api *WxProperty) GetManageRenovation(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	param, _ := gf.RequestParam(c)
	idStr := strings.TrimSpace(c.Query("property_id"))
	if idStr == "" {
		idStr = strings.TrimSpace(gconv.String(param["property_id"]))
	}
	propertyID := gconv.Int64(idStr)
	if propertyID <= 0 {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}

	// 校验房源是否属于当前用户可维护范围
	exists, _ := gf.Model("business_properties").
		Fields("id").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("agent_id", userID).
		Where("deletetime", nil).
		Find()
	if exists == nil || len(exists) == 0 {
		gf.Failed().SetCode(403).SetMsg("房源不存在或无权限").Regin(c)
		return
	}

	row, err := gf.Model("business_renovations").
		Where("property_id", propertyID).
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取装修信息失败：" + err.Error()).Regin(c)
		return
	}

	// 无记录时返回默认结构，便于前端渲染
	if row == nil || len(row) == 0 {
		gf.Success().SetMsg("获取装修信息成功").SetData(gf.Map{
			"property_id":           propertyID,
			"renovation_status":     "none",
			"progress_percentage":   0,
			"current_stage":         "",
			"start_date":            "",
			"estimated_finish_date": "",
			"actual_finish_date":    "",
			"materials":             make([]string, 0),
			"images":                make([]string, 0),
			"notes":                 "",
			"status":                0,
		}).Regin(c)
		return
	}

	images := make([]string, 0)
	for _, p := range wxSplitComma(row["images"].String()) {
		if u := wxFullImgURL(p); u != "" {
			images = append(images, u)
		}
	}

	gf.Success().SetMsg("获取装修信息成功").SetData(gf.Map{
		"property_id":           propertyID,
		"renovation_status":     row["renovation_status"].String(),
		"progress_percentage":   row["progress_percentage"].Int(),
		"current_stage":         row["current_stage"].String(),
		"start_date":            row["start_date"].String(),
		"estimated_finish_date": row["estimated_finish_date"].String(),
		"actual_finish_date":    row["actual_finish_date"].String(),
		"materials":             wxSplitComma(row["materials"].String()),
		"images":                images,
		"notes":                 row["notes"].String(),
		"status":                row["status"].Int(),
	}).Regin(c)
}

// SaveManageRenovation 保存可维护房源装修信息（小程序端编辑页）
// 入参：property_id + 装修信息字段
func (api *WxProperty) SaveManageRenovation(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	param, _ := gf.RequestParam(c)
	propertyID := gconv.Int64(param["property_id"])
	if propertyID <= 0 {
		gf.Failed().SetMsg("请传参数property_id").Regin(c)
		return
	}

	// 校验房源是否属于当前用户可维护范围
	exists, _ := gf.Model("business_properties").
		Fields("id").
		Where("business_id", businessID).
		Where("id", propertyID).
		Where("agent_id", userID).
		Where("deletetime", nil).
		Find()
	if exists == nil || len(exists) == 0 {
		gf.Failed().SetCode(403).SetMsg("房源不存在或无权限").Regin(c)
		return
	}

	saveData := pickMap(param,
		"renovation_status", "progress_percentage", "current_stage",
		"start_date", "estimated_finish_date", "actual_finish_date",
		"materials", "images", "notes", "status",
	)
	if len(saveData) == 0 {
		gf.Failed().SetMsg("暂无可保存字段").Regin(c)
		return
	}

	if _, ok := saveData["renovation_status"]; ok {
		rs := strings.TrimSpace(gconv.String(saveData["renovation_status"]))
		if rs == "" {
			rs = "none"
		}
		if rs != "none" && rs != "in_progress" && rs != "done" {
			gf.Failed().SetMsg("renovation_status参数不合法").Regin(c)
			return
		}
		saveData["renovation_status"] = rs
	}
	if _, ok := saveData["progress_percentage"]; ok {
		p := gconv.Int(saveData["progress_percentage"])
		if p < 0 || p > 100 {
			gf.Failed().SetMsg("progress_percentage参数不合法").Regin(c)
			return
		}
		saveData["progress_percentage"] = p
	}
	if _, ok := saveData["materials"]; ok {
		saveData["materials"] = normalizeTagsToString(saveData["materials"])
	}
	if _, ok := saveData["images"]; ok {
		saveData["images"] = normalizeTagsToString(saveData["images"])
	}
	if _, ok := saveData["status"]; ok {
		v := gconv.Int(saveData["status"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("status参数不合法").Regin(c)
			return
		}
		saveData["status"] = v
	}

	existing, _ := gf.Model("business_renovations").
		Where("property_id", propertyID).
		Find()
	if existing == nil || len(existing) == 0 {
		saveData["property_id"] = propertyID
		_, err := gf.Model("business_renovations").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加装修信息失败").SetData(err).Regin(c)
			return
		}
		gf.Success().SetMsg("添加装修信息成功").Regin(c)
		return
	}

	_, err := gf.Model("business_renovations").
		Where("property_id", propertyID).
		Update(saveData)
	if err != nil {
		gf.Failed().SetMsg("更新装修信息失败").SetData(err).Regin(c)
		return
	}
	gf.Success().SetMsg("更新装修信息成功").Regin(c)
}

// SaveManage 保存可维护房源（新增/编辑）
//
// 说明：
// - agent_id 由后端维护：新增时自动写入当前 userID；更新不允许修改 agent_id
// - 仅允许写入部分字段，避免前端传任意字段更新
func (api *WxProperty) SaveManage(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	param, _ := gf.RequestParam(c)
	editID := gf.GetEditId(param["id"])

	saveData := pickMap(param,
		"title", "price", "price_unit", "area",
		"rooms", "halls", "bathrooms",
		"floor_level", "total_floors", "orientation", "build_year",
		"property_type", "decoration_type",
		"community_name", "address", "latitude", "longitude",
		"tags", "images", "cover_image", "video_url",
		"allow_image_download", "allow_video_download",
		"owner_name", "owner_phone",
		"receiver_name", "receiver_phone", "receiver_price",
		"sale_status",
		"commission_rate", "commission_reward",
		"weigh",
		"hot_status",
		"status",
	)
	if len(saveData) == 0 {
		gf.Failed().SetMsg("暂无可保存字段").Regin(c)
		return
	}

	if _, ok := saveData["tags"]; ok {
		saveData["tags"] = normalizeTagsToString(saveData["tags"])
	}
	if _, ok := saveData["images"]; ok {
		saveData["images"] = normalizeTagsToString(saveData["images"])
	}
	if _, ok := saveData["cover_image"]; ok {
		saveData["cover_image"] = strings.TrimSpace(gconv.String(saveData["cover_image"]))
	}
	if _, ok := saveData["video_url"]; ok {
		saveData["video_url"] = strings.TrimSpace(gconv.String(saveData["video_url"]))
	}
	if _, ok := saveData["allow_image_download"]; ok {
		v := gconv.Int(saveData["allow_image_download"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("allow_image_download参数不合法").Regin(c)
			return
		}
		saveData["allow_image_download"] = v
	}
	if _, ok := saveData["allow_video_download"]; ok {
		v := gconv.Int(saveData["allow_video_download"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("allow_video_download参数不合法").Regin(c)
			return
		}
		saveData["allow_video_download"] = v
	}
	if _, ok := saveData["receiver_price"]; ok {
		raw := strings.TrimSpace(gconv.String(saveData["receiver_price"]))
		if raw == "" {
			saveData["receiver_price"] = nil
		} else {
			saveData["receiver_price"] = gconv.Float64(raw)
		}
	}
	if _, ok := saveData["sale_status"]; ok {
		ss := strings.TrimSpace(gconv.String(saveData["sale_status"]))
		if ss != "" && ss != "on_sale" && ss != "sold" && ss != "off_market" {
			gf.Failed().SetMsg("sale_status参数不合法").Regin(c)
			return
		}
		if ss != "" {
			saveData["sale_status"] = ss
		} else {
			delete(saveData, "sale_status")
		}
	}
	if _, ok := saveData["hot_status"]; ok {
		v := gconv.Int(saveData["hot_status"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("hot_status参数不合法").Regin(c)
			return
		}
		saveData["hot_status"] = v
	}
	if _, ok := saveData["status"]; ok {
		v := gconv.Int(saveData["status"])
		if v != 0 && v != 1 {
			gf.Failed().SetMsg("status参数不合法").Regin(c)
			return
		}
		saveData["status"] = v
	}
	if _, ok := saveData["commission_rate"]; ok {
		raw := strings.TrimSpace(gconv.String(saveData["commission_rate"]))
		if raw == "" {
			saveData["commission_rate"] = nil
		} else {
			v := gconv.Float64(raw)
			if v < 0 || v > 100 {
				gf.Failed().SetMsg("commission_rate参数不合法").Regin(c)
				return
			}
			saveData["commission_rate"] = v
		}
	}
	if _, ok := saveData["commission_reward"]; ok {
		raw := strings.TrimSpace(gconv.String(saveData["commission_reward"]))
		if raw == "" {
			saveData["commission_reward"] = nil
		} else {
			v := gconv.Float64(raw)
			if v < 0 {
				gf.Failed().SetMsg("commission_reward参数不合法").Regin(c)
				return
			}
			saveData["commission_reward"] = v
		}
	}
	if _, ok := saveData["weigh"]; ok {
		v := gconv.Int64(saveData["weigh"])
		if v < 0 {
			gf.Failed().SetMsg("weigh参数不合法").Regin(c)
			return
		}
		saveData["weigh"] = v
	}

	if editID == 0 {
		// 新增：由后端维护 agent_id
		saveData["business_id"] = businessID
		saveData["agent_id"] = userID
		addID, err := gf.Model("business_properties").Data(saveData).InsertAndGetId()
		if err != nil {
			gf.Failed().SetMsg("添加失败").SetData(err).Regin(c)
			return
		}
		if addID != 0 {
			_, _ = gf.Model("business_properties").Where("id", addID).Update(gf.Map{"weigh": addID})
		}
		gf.Success().SetMsg("添加成功").SetData(gf.Map{"id": addID}).Regin(c)
		return
	}

	// 更新：仅允许更新自己维护的房源
	res, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", editID).
		Where("agent_id", userID).
		Update(saveData)
	if err != nil {
		gf.Failed().SetMsg("更新失败").SetData(err).Regin(c)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		gf.Failed().SetCode(403).SetMsg("房源不存在或无权限").Regin(c)
		return
	}
	gf.Success().SetMsg("更新成功").SetData(true).Regin(c)
}

// DelManage 删除可维护房源（软删除）
// 入参：id
func (api *WxProperty) DelManage(c *gf.GinCtx) {
	businessID, userID, ok := wxEnsureCanManageProperties(c)
	if !ok {
		return
	}

	param, _ := gf.RequestParam(c)
	idStr := strings.TrimSpace(c.Query("id"))
	if idStr == "" {
		idStr = strings.TrimSpace(gconv.String(param["id"]))
	}
	id := gconv.Int64(idStr)
	if id <= 0 {
		gf.Failed().SetMsg("请传参数id").Regin(c)
		return
	}

	res, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", id).
		Where("agent_id", userID).
		Delete()
	if err != nil {
		gf.Failed().SetMsg("删除失败").SetData(err).Regin(c)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		gf.Failed().SetCode(403).SetMsg("房源不存在或无权限").Regin(c)
		return
	}
	gf.Success().SetMsg("删除成功").SetData(true).Regin(c)
}

// GetFormOptions 房源表单下拉选项（小程序端缓存用）
//
// 标准做法（企业版统一配置）：
// - 在 resource/config 下新增配置文件（参考 confdemo.yaml）
// - 本项目已新增：resource/config/property_form_options.yaml
// - data 中使用 *_txt 字段配置选项文本（每行一条 value=label）
//
// 兼容做法（历史方案）：
// - 仍支持通过“数据中心 → 字典”维护：common_dictionary_group / common_dictionary_data
func (api *WxProperty) GetFormOptions(c *gf.GinCtx) {
	businessID := wxBusinessID(c)

	// 默认值（字典未配置时兜底）
	out := gf.Map{
		"sale_status": []gf.Map{
			{"label": "在售", "value": "on_sale"},
			{"label": "已售", "value": "sold"},
			{"label": "下架", "value": "off_market"},
		},
		"price_unit": []gf.Map{
			{"label": "万", "value": "万"},
			{"label": "元", "value": "元"},
		},
		"floor_level": []gf.Map{
			{"label": "低层 (1-6)", "value": "低层"},
			{"label": "中层 (7-15)", "value": "中层"},
			{"label": "高层 (16+)", "value": "高层"},
			{"label": "地下", "value": "地下"},
		},
		"orientation": []gf.Map{
			{"label": "东", "value": "东"},
			{"label": "南", "value": "南"},
			{"label": "西", "value": "西"},
			{"label": "北", "value": "北"},
			{"label": "东南", "value": "东南"},
			{"label": "东北", "value": "东北"},
			{"label": "西南", "value": "西南"},
			{"label": "西北", "value": "西北"},
			{"label": "南北", "value": "南北"},
			{"label": "东西", "value": "东西"},
		},
		"property_type": []gf.Map{
			{"label": "住宅", "value": "住宅"},
			{"label": "公寓", "value": "公寓"},
			{"label": "别墅", "value": "别墅"},
			{"label": "商铺", "value": "商铺"},
			{"label": "写字楼", "value": "写字楼"},
		},
		"decoration_type": []gf.Map{
			{"label": "毛坯", "value": "毛坯"},
			{"label": "简装", "value": "简装"},
			{"label": "精装", "value": "精装"},
			{"label": "豪装", "value": "豪装"},
		},
		// 装修工序（当前阶段）
		"renovation_stage": []gf.Map{
			{"label": "设计", "value": "设计"},
			{"label": "拆改", "value": "拆改"},
			{"label": "水电", "value": "水电"},
			{"label": "泥瓦", "value": "泥瓦"},
			{"label": "木工", "value": "木工"},
			{"label": "油漆", "value": "油漆"},
			{"label": "安装", "value": "安装"},
			{"label": "软装", "value": "软装"},
			{"label": "验收", "value": "验收"},
		},
	}

	hasSaleStatus := false
	hasPriceUnit := false
	hasFloorLevel := false
	hasOrientation := false
	hasPropertyType := false
	hasDecorationType := false
	hasRenovationStage := false

	// 1) 标准配置：resource/config/*.yaml（后台“配置管理”可维护）
	if conf, err := gf.GetConfByFile("property_form_options"); err == nil && conf != nil {
		if m, ok := conf.(map[string]interface{}); ok {
			if gconv.String(m["conftype"]) == "configuration" && gf.Bool(m["status"]) {
				if data, ok := m["data"].(map[string]interface{}); ok {
					if list := wxParseOptionsText(data["sale_status_txt"]); len(list) > 0 {
						out["sale_status"] = list
						hasSaleStatus = true
					}
					if list := wxParseOptionsText(data["price_unit_txt"]); len(list) > 0 {
						out["price_unit"] = list
						hasPriceUnit = true
					}
					if list := wxParseOptionsText(data["floor_level_txt"]); len(list) > 0 {
						out["floor_level"] = list
						hasFloorLevel = true
					}
					if list := wxParseOptionsText(data["orientation_txt"]); len(list) > 0 {
						out["orientation"] = list
						hasOrientation = true
					}
					if list := wxParseOptionsText(data["property_type_txt"]); len(list) > 0 {
						out["property_type"] = list
						hasPropertyType = true
					}
					if list := wxParseOptionsText(data["decoration_type_txt"]); len(list) > 0 {
						out["decoration_type"] = list
						hasDecorationType = true
					}
					if list := wxParseOptionsText(data["renovation_stage_txt"]); len(list) > 0 {
						out["renovation_stage"] = list
						hasRenovationStage = true
					}
				}
			}
		}
	}

	// 2) 兼容：字典覆盖（历史方案，有配置则覆盖默认）
	if !hasSaleStatus {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源销售状态", "销售状态", "房源-销售状态"}); ok {
			out["sale_status"] = list
		}
	}
	if !hasPriceUnit {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源价格单位", "价格单位", "房源-价格单位"}); ok {
			out["price_unit"] = list
		}
	}
	if !hasFloorLevel {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源楼层位置", "楼层位置", "房源-楼层位置"}); ok {
			out["floor_level"] = list
		}
	}
	if !hasOrientation {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源朝向", "朝向", "房源-朝向"}); ok {
			out["orientation"] = list
		}
	}
	if !hasPropertyType {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源物业类型", "物业类型", "房源-物业类型"}); ok {
			out["property_type"] = list
		}
	}
	if !hasDecorationType {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源装修标准", "装修标准", "房源-装修标准"}); ok {
			out["decoration_type"] = list
		}
	}
	if !hasRenovationStage {
		if list, ok := wxPickDictOptionsByTitles(businessID, []string{"房源装修工序", "装修工序", "房源-装修工序", "装修阶段", "房源-装修阶段"}); ok {
			out["renovation_stage"] = list
		}
	}

	gf.Success().SetMsg("获取表单选项").SetData(out).Regin(c)
}
