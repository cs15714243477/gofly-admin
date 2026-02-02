package houses

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
)

// 小程序房源详情（避免走 RBAC：NoNeedAuths 传 "*"）
type WxProperty struct{ NoNeedAuths []string }

func init() {
	fpath := WxProperty{NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

func wxBusinessID(c *gf.GinCtx) int64 {
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}
	return businessID
}

// wxSaleStatusLabel 根据销售状态码返回对应的中文状态描述
//
// @param status 销售状态码，支持 "on_sale"、"sold"、"off_market"
// @return 对应的中文状态描述，如果状态码不在支持范围内则返回空字符串
func wxSaleStatusLabel(status string) string {
	switch status {
	case "on_sale":
		return "在售"
	case "sold":
		return "已售"
	case "off_market":
		return "下架"
	default:
		return ""
	}
}

func wxFullImgURL(path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	return gf.GetFullUrl(p)
}

func wxSplitComma(s string) []string {
	raw := strings.TrimSpace(s)
	if raw == "" {
		return nil
	}
	parts := gf.SplitAndStr(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func wxBuildImages(cover string, images string) []string {
	out := make([]string, 0)
	if u := wxFullImgURL(cover); u != "" {
		out = append(out, u)
	}
	for _, p := range wxSplitComma(images) {
		if u := wxFullImgURL(p); u != "" {
			out = append(out, u)
		}
	}
	// 去重
	seen := make(map[string]struct{}, len(out))
	uniq := make([]string, 0, len(out))
	for _, u := range out {
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		uniq = append(uniq, u)
	}
	return uniq
}

// GetDetail 房源详情（小程序）
// 入参：id
// 返回：property、images、is_followed、renovation、recommends
func (api *WxProperty) GetDetail(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")

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
		Where("deletetime", nil).
		Where("status", 0).
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,floor_level,total_floors,orientation,build_year,property_type,decoration_type,community_name,address,latitude,longitude,tags,images,cover_image,video_url,allow_image_download,allow_video_download,has_smart_lock,commission_rate,commission_reward,owner_phone,agent_id,sale_status,hot_status,view_count,follow_count,showing_count,createtime").
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取房源详情失败：" + err.Error()).Regin(c)
		return
	}
	if row == nil || len(row) == 0 {
		gf.Failed().SetMsg("房源不存在或已下架").Regin(c)
		return
	}

	// 浏览+1（不阻塞响应）
	go func() {
		_, _ = gf.Model("business_properties").
			Where("business_id", businessID).
			Where("id", id).
			Update(gf.Map{"view_count": gform.Raw("view_count + 1")})
	}()

	tags := wxSplitComma(row["tags"].String())
	images := wxBuildImages(row["cover_image"].String(), row["images"].String())
	saleStatus := row["sale_status"].String()

	isFollowed := false
	if userID > 0 {
		fav, _ := gf.Model("business_favorites").
			Where("user_id", userID).
			Where("property_id", id).
			Find()
		isFollowed = fav != nil && len(fav) > 0
	}

	// 装修信息（可为空）
	renovation, _ := gf.Model("business_renovations").Where("property_id", id).Find()
	if renovation != nil {
		if renovation["materials"] != nil && strings.TrimSpace(renovation["materials"].String()) != "" {
			renovation["materials"] = gf.VarNew(wxSplitComma(renovation["materials"].String()))
		} else {
			renovation["materials"] = gf.VarNew(make([]string, 0))
		}
		if renovation["images"] != nil && strings.TrimSpace(renovation["images"].String()) != "" {
			imgs := make([]string, 0)
			for _, p := range wxSplitComma(renovation["images"].String()) {
				if u := wxFullImgURL(p); u != "" {
					imgs = append(imgs, u)
				}
			}
			renovation["images"] = gf.VarNew(imgs)
		} else {
			renovation["images"] = gf.VarNew(make([]string, 0))
		}
	}

	// 推荐（同商户 hot_status=1 优先；不足则按 weigh）
	recRows, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("deletetime", nil).
		Where("status", 0).
		Where("id <>", id).
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,cover_image,images,sale_status,weigh,createtime").
		Order("hot_status desc, weigh desc, id desc").
		Limit(6).
		Select()
	recommends := make([]gf.Map, 0, len(recRows))
	for _, r := range recRows {
		if r == nil {
			continue
		}
		img := wxFullImgURL(r["cover_image"].String())
		if img == "" && r["images"] != nil && strings.TrimSpace(r["images"].String()) != "" {
			first := strings.Split(r["images"].String(), ",")[0]
			img = wxFullImgURL(first)
		}
		ss := r["sale_status"].String()
		recommends = append(recommends, gf.Map{
			"id":                r["id"].Int64(),
			"title":             r["title"].String(),
			"price":             gconv.String(r["price"]),
			"price_unit":        r["price_unit"].String(),
			"area":              gconv.String(r["area"]),
			"rooms":             r["rooms"].Int(),
			"halls":             r["halls"].Int(),
			"bathrooms":         r["bathrooms"].Int(),
			"image":             img,
			"sale_status":       ss,
			"sale_status_label": wxSaleStatusLabel(ss),
		})
	}

	property := gf.Map{
		"id":                row["id"].Int64(),
		"title":             row["title"].String(),
		"price":             gconv.String(row["price"]),
		"price_unit":        row["price_unit"].String(),
		"area":              gconv.String(row["area"]),
		"rooms":             row["rooms"].Int(),
		"halls":             row["halls"].Int(),
		"bathrooms":         row["bathrooms"].Int(),
		"floor_level":       row["floor_level"].String(),
		"total_floors":      row["total_floors"].Int(),
		"orientation":       row["orientation"].String(),
		"build_year":        row["build_year"].Int(),
		"property_type":     row["property_type"].String(),
		"decoration_type":   row["decoration_type"].String(),
		"community_name":    row["community_name"].String(),
		"address":           row["address"].String(),
		"latitude":          row["latitude"],
		"longitude":         row["longitude"],
		"tags":              tags,
		"images":            images,
		"cover_image":       wxFullImgURL(row["cover_image"].String()),
		"video_url":         wxFullImgURL(row["video_url"].String()),
		"allow_image_download": row["allow_image_download"].Int(),
		"allow_video_download": row["allow_video_download"].Int(),
		"has_smart_lock":    row["has_smart_lock"].Int(),
		"commission_rate":   gconv.String(row["commission_rate"]),
		"commission_reward": gconv.String(row["commission_reward"]),
		"owner_phone":       row["owner_phone"].String(),
		"agent_id":          row["agent_id"].Int64(),
		"sale_status":       saleStatus,
		"sale_status_label": wxSaleStatusLabel(saleStatus),
		"hot_status":        row["hot_status"].Int(),
		"view_count":        row["view_count"].Int(),
		"follow_count":      row["follow_count"].Int(),
		"showing_count":     row["showing_count"].Int(),
		"createtime":        row["createtime"].String(),
	}

	gf.Success().SetMsg("获取房源详情").SetData(gf.Map{
		"property":     property,
		"images":       images,
		"is_followed":  isFollowed,
		"renovation":   renovation,
		"recommends":   recommends,
		"view_added":   true,
		"business_id":  businessID,
		"current_user": userID,
	}).Regin(c)
}

// ToggleFollow 关注/取消关注
// 入参：id
// 返回：is_followed、follow_count
func (api *WxProperty) ToggleFollow(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetMsg("请先登录").Regin(c)
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

	// 校验房源存在且属于当前商户
	exists, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", id).
		Where("deletetime", nil).
		Where("status", 0).
		Fields("id").
		Find()
	if exists == nil || len(exists) == 0 {
		gf.Failed().SetMsg("房源不存在或已下架").Regin(c)
		return
	}

	favMDB := gf.Model("business_favorites").Where("user_id", userID).Where("property_id", id)
	fav, _ := favMDB.Find()
	isFollowed := false
	if fav != nil && len(fav) > 0 {
		_, _ = favMDB.Delete()
		_, _ = gf.Model("business_properties").
			Where("business_id", businessID).
			Where("id", id).
			Update(gf.Map{"follow_count": gform.Raw("if(follow_count>0, follow_count-1, 0)")})
		isFollowed = false
	} else {
		_, _ = gf.Model("business_favorites").Data(gf.Map{
			"user_id":     userID,
			"property_id": id,
		}).Insert()
		_, _ = gf.Model("business_properties").
			Where("business_id", businessID).
			Where("id", id).
			Update(gf.Map{"follow_count": gform.Raw("follow_count + 1")})
		isFollowed = true
	}

	// 返回最新关注数
	latest, _ := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", id).
		Fields("follow_count").
		Find()
	followCount := 0
	if latest != nil {
		followCount = latest["follow_count"].Int()
	}

	gf.Success().SetMsg("操作成功").SetData(gf.Map{
		"is_followed":  isFollowed,
		"follow_count": followCount,
	}).Regin(c)
}
