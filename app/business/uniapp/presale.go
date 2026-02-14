package uniapp

import (
	"strings"

	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
)

// 预售房源（小程序 Tab 专用）
//
// 规则：
// - 仅展示 sale_status=in_sale
// - 详情必须携带 view_key（由列表接口下发），避免“直接带 id 访问”绕过入口
type PreSale struct {
	NoNeedLogin []string
	NoNeedAuths []string
}

func init() {
	fpath := PreSale{NoNeedLogin: []string{}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

// GetList 预售房源列表（小程序 Tab）
// 入参：page,pageSize,keyword,sort
func (api *PreSale) GetList(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	param, _ := gf.RequestParam(c)
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		keyword = strings.TrimSpace(gconv.String(param["keyword"]))
	}
	sort := strings.TrimSpace(c.Query("sort"))
	if sort == "" {
		sort = strings.TrimSpace(gconv.String(param["sort"]))
	}

	MDB := gf.Model("business_properties").
		Where("deletetime", nil).
		Where("status", 0).
		Where("business_id", businessID).
		Where("sale_status", "in_sale")

	if keyword != "" {
		kw := "%" + keyword + "%"
		MDB = MDB.Where("(title like ? OR community_name like ? OR address like ?)", kw, kw, kw)
	}

	// 排序（复用房源列表的常用排序规则）
	switch sort {
	case "price_asc":
		MDB = MDB.Order("price asc, weigh desc, id desc")
	case "price_desc":
		MDB = MDB.Order("price desc, weigh desc, id desc")
	case "latest":
		MDB = MDB.Order("createtime desc, id desc")
	case "view_desc":
		MDB = MDB.Order("view_count desc, weigh desc, id desc")
	case "follow_desc":
		MDB = MDB.Order("follow_count desc, weigh desc, id desc")
	default:
		MDB = MDB.Order("weigh desc, id desc")
	}

	total, _ := MDB.Clone().Count()
	list, err := MDB.
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,orientation,property_type,decoration_type,community_name,address,tags,cover_image,images,has_smart_lock,agent_id,sale_status,hot_status,view_count,follow_count,showing_count,weigh,createtime").
		Page(pageNo, pageSize).
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取预售房源列表失败：" + err.Error()).Regin(c)
		return
	}

	items := make([]gf.Map, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		propertyID := row["id"].Int64()
		img := fullImgURL(row["cover_image"].String())
		if img == "" && row["images"] != nil && strings.TrimSpace(row["images"].String()) != "" {
			first := strings.Split(row["images"].String(), ",")[0]
			img = fullImgURL(first)
		}
		tags := make([]string, 0)
		if row["tags"] != nil && strings.TrimSpace(row["tags"].String()) != "" {
			tags = gf.SplitAndStr(row["tags"].String(), ",")
		}
		ss := strings.TrimSpace(row["sale_status"].String())
		items = append(items, gf.Map{
			"id":                propertyID,
			"title":             row["title"].String(),
			"price":             gconv.String(row["price"]),
			"price_unit":        row["price_unit"].String(),
			"area":              gconv.String(row["area"]),
			"rooms":             row["rooms"].Int(),
			"halls":             row["halls"].Int(),
			"bathrooms":         row["bathrooms"].Int(),
			"orientation":       row["orientation"].String(),
			"property_type":     row["property_type"].String(),
			"decoration_type":   row["decoration_type"].String(),
			"community_name":    row["community_name"].String(),
			"address":           row["address"].String(),
			"tags":              tags,
			"image":             img,
			"has_smart_lock":    row["has_smart_lock"].Int(),
			"agent_id":          row["agent_id"].Int64(),
			"sale_status":       ss,
			"sale_status_label": saleStatusLabel(ss),
			"hot_status":        row["hot_status"].Int(),
			"view_count":        row["view_count"].Int(),
			"follow_count":      row["follow_count"].Int(),
			"showing_count":     row["showing_count"].Int(),
			"weigh":             row["weigh"].Int64(),
			"createtime":        row["createtime"].String(),
			// 预售详情访问校验 key（与当前登录用户绑定）
			"view_key": wxBuildPreSaleViewKey(businessID, userID, propertyID),
		})
	}

	gf.Success().SetMsg("获取预售房源列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    total,
		"items":    items,
	}).Regin(c)
}

// GetDetail 预售房源详情（小程序 Tab）
// 入参：id,view_key
func (api *PreSale) GetDetail(c *gf.GinCtx) {
	businessID := wxBusinessID(c)
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
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

	viewKey := strings.TrimSpace(c.Query("view_key"))
	if viewKey == "" {
		viewKey = strings.TrimSpace(gconv.String(param["view_key"]))
	}
	if !wxCheckPreSaleViewKey(businessID, userID, id, viewKey) {
		gf.Failed().SetMsg("房源不存在或不可见").Regin(c)
		return
	}

	row, err := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("id", id).
		Where("deletetime", nil).
		Where("status", 0).
		Where("sale_status", "in_sale").
		Fields("id,title,custom_desc,price,price_unit,area,rooms,halls,bathrooms,floor_level,total_floors,orientation,build_year,property_type,decoration_type,community_name,address,latitude,longitude,tags,images,cover_image,video_url,allow_image_download,allow_video_download,has_smart_lock,commission_rate,commission_reward,owner_name,owner_phone,receiver_name,receiver_phone,receiver_price,agent_id,sale_status,hot_status,view_count,follow_count,showing_count,createtime").
		Find()
	if err != nil {
		gf.Failed().SetMsg("获取房源详情失败：" + err.Error()).Regin(c)
		return
	}
	if row == nil || len(row) == 0 {
		gf.Failed().SetMsg("房源不存在或不可见").Regin(c)
		return
	}

	// 浏览+1（不阻塞响应）
	go func() {
		_, _ = gf.Model("business_properties").
			Where("business_id", businessID).
			Where("id", id).
			Update(gf.Map{"view_count": gform.Raw("view_count + 1")})
	}()

	// 预售 Tab 仅做展示：按公开详情处理（隐藏隐私字段 & 底部操作栏）
	tags := wxSplitComma(row["tags"].String())
	images := wxBuildImages(row["cover_image"].String(), row["images"].String())
	saleStatus := strings.TrimSpace(row["sale_status"].String())

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
		// 工序时间线（stage_logs）：JSON 数组
		if renovation["stage_logs"] != nil && strings.TrimSpace(renovation["stage_logs"].String()) != "" {
			list := parseRenovationStageLogs(renovation["stage_logs"].String())
			out := make([]gf.Map, 0, len(list))
			for _, it := range list {
				imgs := make([]string, 0)
				for _, p := range wxSplitComma(gconv.String(it["images"])) {
					if u := wxFullImgURL(p); u != "" {
						imgs = append(imgs, u)
					}
				}
				it["images"] = imgs
				out = append(out, it)
			}
			renovation["stage_logs"] = gf.VarNew(out)
		} else {
			renovation["stage_logs"] = gf.VarNew(make([]gf.Map, 0))
		}
	}

	// 经纪人公开信息（供详情页展示）
	agentInfo := gf.Map{}
	agentID := row["agent_id"].Int64()
	if agentID > 0 {
		agent, _ := gf.Model("business_user").
			Fields("id,name,nickname,mobile,avatar,title,status,store_id,business_id").
			Where("id", agentID).
			Where("deletetime", nil).
			Find()
		if agent != nil && agent["status"].Int() == 0 {
			agentName := agent["name"].String()
			if strings.TrimSpace(agentName) == "" {
				agentName = agent["nickname"].String()
			}
			agentStoreName := "未绑定"
			if agent["store_id"].Int64() > 0 {
				store, _ := gf.Model("business_stores").
					Fields("name,status").
					Where("id", agent["store_id"].Int64()).
					Where("business_id", agent["business_id"].Int64()).
					Where("deletetime", nil).
					Find()
				if store != nil && store["status"].Int() == 0 && strings.TrimSpace(store["name"].String()) != "" {
					agentStoreName = store["name"].String()
				}
			}
			agentAvatar := agent["avatar"].String()
			if strings.TrimSpace(agentAvatar) == "" {
				agentAvatar = gf.GetLocalUrl() + "resource/uploads/static/avatar.png"
			} else if !strings.HasPrefix(agentAvatar, "http://") && !strings.HasPrefix(agentAvatar, "https://") {
				agentAvatar = gf.GetFullUrl(agentAvatar)
			}
			agentInfo = gf.Map{
				"id":         agent["id"].Int64(),
				"name":       agentName,
				"title":      agent["title"].String(),
				"mobile":     agent["mobile"].String(),
				"avatar":     agentAvatar,
				"store_name": agentStoreName,
			}
		}
	}

	property := gf.Map{
		"id":                   row["id"].Int64(),
		"title":                row["title"].String(),
		"custom_desc":          gconv.String(row["custom_desc"]),
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
		"latitude":             row["latitude"],
		"longitude":            row["longitude"],
		"tags":                 tags,
		"images":               images,
		"cover_image":          wxFullImgURL(row["cover_image"].String()),
		"video_url":            wxFullImgURL(row["video_url"].String()),
		"allow_image_download": row["allow_image_download"].Int(),
		"allow_video_download": row["allow_video_download"].Int(),
		"has_smart_lock":       row["has_smart_lock"].Int(),
		"commission_rate":      "",
		"commission_reward":    "",
		"owner_name":           "",
		"owner_phone":          "",
		"receiver_name":        "",
		"receiver_phone":       "",
		"receiver_price":       "",
		"agent_id":             row["agent_id"].Int64(),
		"sale_status":          saleStatus,
		"sale_status_label":    wxSaleStatusLabel(saleStatus),
		"hot_status":           row["hot_status"].Int(),
		"view_count":           row["view_count"].Int(),
		"follow_count":         row["follow_count"].Int(),
		"showing_count":        row["showing_count"].Int(),
		"createtime":           row["createtime"].String(),
		"agent_name":           gconv.String(agentInfo["name"]),
		"agent_mobile":         gconv.String(agentInfo["mobile"]),
		"agent_title":          gconv.String(agentInfo["title"]),
		"agent_store_name":     gconv.String(agentInfo["store_name"]),
	}

	gf.Success().SetMsg("获取房源详情").SetData(gf.Map{
		"property":     property,
		"images":       images,
		"is_followed":  false,
		"renovation":   renovation,
		"recommends":   make([]gf.Map, 0),
		"agent":        agentInfo,
		"public_view":  true,
		"view_added":   true,
		"business_id":  businessID,
		"current_user": userID,
	}).Regin(c)
}
