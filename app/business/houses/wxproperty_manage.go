package houses

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gmap"
)

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
		"sale_status":          row["sale_status"].String(),
		"status":               row["status"].Int(),
	}

	gf.Success().SetMsg("获取房源成功").SetData(out).Regin(c)
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
