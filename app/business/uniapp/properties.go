package uniapp

import (
	"sort"
	"strings"

	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
)

// 房源列表（小程序端）
type Properties struct {
	NoNeedLogin []string
	NoNeedAuths []string
}

func init() {
	fpath := Properties{NoNeedLogin: []string{}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

func (api *Properties) GetList(c *gf.GinCtx) {
	// businessID：优先从登录态，其次 header，最后默认 1
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}
	userID := c.GetInt64("userID")

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
	category := strings.TrimSpace(c.Query("category"))
	if category == "" {
		category = strings.TrimSpace(gconv.String(param["category"]))
	}
	sort := strings.TrimSpace(c.Query("sort"))
	if sort == "" {
		sort = strings.TrimSpace(gconv.String(param["sort"]))
	}
	province := strings.TrimSpace(c.Query("province"))
	if province == "" {
		province = strings.TrimSpace(gconv.String(param["province"]))
	}
	city := strings.TrimSpace(c.Query("city"))
	if city == "" {
		city = strings.TrimSpace(gconv.String(param["city"]))
	}
	district := strings.TrimSpace(c.Query("district"))
	if district == "" {
		district = strings.TrimSpace(gconv.String(param["district"]))
	}
	saleStatus := strings.TrimSpace(c.Query("sale_status"))
	if saleStatus == "" {
		saleStatus = strings.TrimSpace(gconv.String(param["sale_status"]))
	}
	hotStatusStr := strings.TrimSpace(c.Query("hot_status"))
	if hotStatusStr == "" {
		hotStatusStr = strings.TrimSpace(gconv.String(param["hot_status"]))
	}
	propertyType := strings.TrimSpace(c.Query("property_type"))
	if propertyType == "" {
		propertyType = strings.TrimSpace(gconv.String(param["property_type"]))
	}
	decorationType := strings.TrimSpace(c.Query("decoration_type"))
	if decorationType == "" {
		decorationType = strings.TrimSpace(gconv.String(param["decoration_type"]))
	}
	orientation := strings.TrimSpace(c.Query("orientation"))
	if orientation == "" {
		orientation = strings.TrimSpace(gconv.String(param["orientation"]))
	}
	buildYearMin := strings.TrimSpace(c.Query("build_year_min"))
	if buildYearMin == "" {
		buildYearMin = strings.TrimSpace(gconv.String(param["build_year_min"]))
	}
	buildYearMax := strings.TrimSpace(c.Query("build_year_max"))
	if buildYearMax == "" {
		buildYearMax = strings.TrimSpace(gconv.String(param["build_year_max"]))
	}
	priceMin := strings.TrimSpace(c.Query("price_min"))
	if priceMin == "" {
		priceMin = strings.TrimSpace(gconv.String(param["price_min"]))
	}
	priceMax := strings.TrimSpace(c.Query("price_max"))
	if priceMax == "" {
		priceMax = strings.TrimSpace(gconv.String(param["price_max"]))
	}
	areaMin := strings.TrimSpace(c.Query("area_min"))
	if areaMin == "" {
		areaMin = strings.TrimSpace(gconv.String(param["area_min"]))
	}
	areaMax := strings.TrimSpace(c.Query("area_max"))
	if areaMax == "" {
		areaMax = strings.TrimSpace(gconv.String(param["area_max"]))
	}
	roomsStr := strings.TrimSpace(c.Query("rooms"))
	if roomsStr == "" {
		roomsStr = strings.TrimSpace(gconv.String(param["rooms"]))
	}
	filter := strings.TrimSpace(c.Query("filter"))
	if filter == "" {
		filter = strings.TrimSpace(gconv.String(param["filter"]))
	}
	more := parseStringSlice(param["more"])
	if len(more) == 0 {
		more = parseStringSlice(c.Query("more"))
	}

	// 小程序端按“当前经纪人（business_user.id）”维度展示房源更符合业务：
	// business_properties.agent_id 维护经纪人ID = business_user.id
	// 若未登录（极少数场景）则回退到 business_id 维度。
	MDB := gf.Model("business_properties").
		Where("deletetime", nil).
		Where("status", 0)

	// 分类
	switch category {
	case "mine":
		if userID > 0 {
			MDB = MDB.Where("agent_id", userID)
		} else {
			MDB = MDB.Where("business_id", businessID)
		}
		// 兼容：已按 agent_id 限定，无需额外处理
	case "school":
		MDB = MDB.Where("tags like ?", "%学区房%")
	case "price_drop":
		MDB = MDB.Where("tags like ?", "%降价%")
	case "nearby", "all", "":
		// nearby 需要经纬度，前端未接入时按 all 处理
	default:
	}

	// 关键词
	if keyword != "" {
		kw := "%" + keyword + "%"
		MDB = MDB.Where("(title like ? OR community_name like ? OR address like ?)", kw, kw, kw)
	}

	// 省份（默认由前端传辽宁；为空则不过滤）
	if province != "" {
		if len([]rune(province)) > 20 {
			gf.Failed().SetMsg("province参数过长").Regin(c)
			return
		}
		MDB = MDB.Where("address like ?", "%"+province+"%")
	}

	// 城市
	if city != "" {
		if len([]rune(city)) > 50 {
			gf.Failed().SetMsg("city参数过长").Regin(c)
			return
		}
		MDB = MDB.Where("address like ?", "%"+city+"%")
	}

	// 区县
	if district != "" {
		if len([]rune(district)) > 50 {
			gf.Failed().SetMsg("district参数过长").Regin(c)
			return
		}
		MDB = MDB.Where("address like ?", "%"+district+"%")
	}

	// 推荐
	if hotStatusStr != "" {
		hs := gconv.Int(hotStatusStr)
		if hs != 0 && hs != 1 {
			gf.Failed().SetMsg("hot_status参数不合法").Regin(c)
			return
		}
		MDB = MDB.Where("hot_status", hs)
	}

	// 物业类型/装修类型/朝向
	if propertyType != "" {
		MDB = MDB.Where("property_type like ?", "%"+propertyType+"%")
	}
	if decorationType != "" {
		MDB = MDB.Where("decoration_type like ?", "%"+decorationType+"%")
	}
	if orientation != "" {
		MDB = MDB.Where("orientation like ?", "%"+orientation+"%")
	}

	// 建成年份
	if buildYearMin != "" {
		y := gconv.Int(buildYearMin)
		if y > 0 {
			MDB = MDB.Where("build_year >= ?", y)
		}
	}
	if buildYearMax != "" {
		y := gconv.Int(buildYearMax)
		if y > 0 {
			MDB = MDB.Where("build_year <= ?", y)
		}
	}

	// 价格/面积区间
	if priceMin != "" {
		v := gconv.Float64(priceMin)
		MDB = MDB.Where("price >= ?", v)
	}
	if priceMax != "" {
		v := gconv.Float64(priceMax)
		MDB = MDB.Where("price <= ?", v)
	}
	if areaMin != "" {
		v := gconv.Float64(areaMin)
		MDB = MDB.Where("area >= ?", v)
	}
	if areaMax != "" {
		v := gconv.Float64(areaMax)
		MDB = MDB.Where("area <= ?", v)
	}

	// 销售状态
	if saleStatus != "" {
		if saleStatus != "on_sale" && saleStatus != "sold" && saleStatus != "off_market" {
			gf.Failed().SetMsg("sale_status参数不合法").Regin(c)
			return
		}
		MDB = MDB.Where("sale_status", saleStatus)
	}

	// 户型 rooms（支持 4+）
	if roomsStr != "" {
		if strings.HasSuffix(roomsStr, "+") {
			minRooms := gconv.Int(strings.TrimSuffix(roomsStr, "+"))
			if minRooms > 0 {
				MDB = MDB.Where("rooms >= ?", minRooms)
			}
		} else {
			r := gconv.Int(roomsStr)
			if r > 0 {
				MDB = MDB.Where("rooms", r)
			}
		}
	}

	// filter：复用首页规则
	if filter == "smart_lock" {
		MDB = MDB.Where("has_smart_lock", 1)
	} else if filter != "" {
		MDB = MDB.Where("tags like ?", "%"+filter+"%")
	}

	// more：多标签
	for _, m := range more {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		MDB = MDB.Where("tags like ?", "%"+m+"%")
	}

	// 排序
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
	case "recommend", "":
		MDB = MDB.Order("weigh desc, id desc")
	default:
		MDB = MDB.Order("weigh desc, id desc")
	}

	total, _ := MDB.Clone().Count()
	list, err := MDB.
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,orientation,property_type,decoration_type,community_name,address,tags,cover_image,images,has_smart_lock,commission_rate,commission_reward,owner_phone,agent_id,sale_status,hot_status,view_count,follow_count,showing_count,weigh,createtime").
		Page(pageNo, pageSize).
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取房源列表失败：" + err.Error()).Regin(c)
		return
	}

	items := make([]gf.Map, 0, len(list))
	for _, row := range list {
		if row == nil {
			continue
		}
		img := fullImgURL(row["cover_image"].String())
		if img == "" && row["images"] != nil && strings.TrimSpace(row["images"].String()) != "" {
			first := strings.Split(row["images"].String(), ",")[0]
			img = fullImgURL(first)
		}
		tags := make([]string, 0)
		if row["tags"] != nil && strings.TrimSpace(row["tags"].String()) != "" {
			tags = gf.SplitAndStr(row["tags"].String(), ",")
		}
		items = append(items, gf.Map{
			"id":                row["id"].Int64(),
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
			"commission_rate":   gconv.String(row["commission_rate"]),
			"commission_reward": gconv.String(row["commission_reward"]),
			"owner_phone":       row["owner_phone"].String(),
			"agent_id":          row["agent_id"].Int64(),
			"sale_status":       row["sale_status"].String(),
			"sale_status_label": saleStatusLabel(row["sale_status"].String()),
			"hot_status":        row["hot_status"].Int(),
			"view_count":        row["view_count"].Int(),
			"follow_count":      row["follow_count"].Int(),
			"showing_count":     row["showing_count"].Int(),
			"weigh":             row["weigh"].Int64(),
			"createtime":        row["createtime"].String(),
		})
	}

	gf.Success().SetMsg("获取房源列表").SetData(gf.Map{
		"page":     pageNo,
		"pageSize": pageSize,
		"total":    total,
		"items":    items,
	}).Regin(c)
}

func (api *Properties) GetFilterOptions(c *gf.GinCtx) {
	// businessID：优先从登录态，其次 header，最后默认 1
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}
	userID := c.GetInt64("userID")

	param, _ := gf.RequestParam(c)

	province := strings.TrimSpace(c.Query("province"))
	if province == "" {
		province = strings.TrimSpace(gconv.String(param["province"]))
	}
	city := strings.TrimSpace(c.Query("city"))
	if city == "" {
		city = strings.TrimSpace(gconv.String(param["city"]))
	}
	district := strings.TrimSpace(c.Query("district"))
	if district == "" {
		district = strings.TrimSpace(gconv.String(param["district"]))
	}
	category := strings.TrimSpace(c.Query("category"))
	if category == "" {
		category = strings.TrimSpace(gconv.String(param["category"]))
	}

	// 小程序端筛选项也按“当前经纪人（business_user.id）”维度统计
	MDB := gf.Model("business_properties").
		Where("deletetime", nil).
		Where("status", 0)
	if userID > 0 {
		MDB = MDB.Where("agent_id", userID)
	} else {
		MDB = MDB.Where("business_id", businessID)
	}

	// 分类（与列表一致，便于筛选项联动）
	switch category {
	case "mine":
		if userID > 0 {
			// 兼容：已按 agent_id 限定，无需额外处理
		}
	case "school":
		MDB = MDB.Where("tags like ?", "%学区房%")
	case "price_drop":
		MDB = MDB.Where("tags like ?", "%降价%")
	case "nearby", "all", "":
	default:
	}

	// 省/市/区（简单用 address like 匹配）
	if province != "" {
		MDB = MDB.Where("address like ?", "%"+province+"%")
	}
	if city != "" {
		MDB = MDB.Where("address like ?", "%"+city+"%")
	}
	if district != "" {
		MDB = MDB.Where("address like ?", "%"+district+"%")
	}

	extractDistinct := func(field string) []string {
		out := make([]string, 0)
		rows, err := MDB.Clone().
			Where(field + " <> ''").
			Fields(field).
			Group(field).
			Order(field + " asc").
			Select()
		if err != nil {
			return out
		}
		for _, r := range rows {
			if r == nil {
				continue
			}
			s := strings.TrimSpace(r[field].String())
			if s == "" {
				continue
			}
			out = append(out, s)
		}
		return out
	}

	propertyTypes := extractDistinct("property_type")
	decorationTypes := extractDistinct("decoration_type")
	orientations := extractDistinct("orientation")

	// tags：取样后拆分去重，并按频次排序取前 N
	tagCounts := make(map[string]int)
	tagRows, _ := MDB.Clone().
		Where("tags <> ''").
		Fields("tags").
		Limit(2000).
		Select()
	for _, r := range tagRows {
		if r == nil {
			continue
		}
		raw := strings.TrimSpace(r["tags"].String())
		if raw == "" {
			continue
		}
		for _, t := range gf.SplitAndStr(raw, ",") {
			t = strings.TrimSpace(t)
			if t == "" {
				continue
			}
			tagCounts[t] = tagCounts[t] + 1
		}
	}
	type tagPair struct {
		Name  string
		Count int
	}
	pairs := make([]tagPair, 0, len(tagCounts))
	for k, v := range tagCounts {
		pairs = append(pairs, tagPair{Name: k, Count: v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Count != pairs[j].Count {
			return pairs[i].Count > pairs[j].Count
		}
		return pairs[i].Name < pairs[j].Name
	})
	tags := make([]string, 0, 30)
	for i := 0; i < len(pairs) && i < 30; i++ {
		tags = append(tags, pairs[i].Name)
	}

	stats, _ := MDB.Clone().
		Fields("min(price) as price_min,max(price) as price_max,min(area) as area_min,max(area) as area_max").
		Find()
	yearStats, _ := MDB.Clone().
		Where("build_year >", 0).
		Fields("min(build_year) as build_year_min,max(build_year) as build_year_max").
		Find()
	if stats == nil {
		stats = gform.Record{}
	}
	if yearStats != nil {
		for k, v := range yearStats {
			stats[k] = v
		}
	}

	gf.Success().SetMsg("获取筛选项").SetData(gf.Map{
		"property_types":   propertyTypes,
		"decoration_types": decorationTypes,
		"orientations":     orientations,
		"tags":             tags,
		"stats":            stats,
	}).Regin(c)
}
