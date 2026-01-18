package uniapp

import (
	"encoding/json"
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// 首页（推荐/关注）
type Home struct {
	NoNeedLogin []string
	NoNeedAuths []string
}

func init() {
	// 小程序端不走 RBAC 权限；是否需要登录由 NoNeedLogin 控制
	fpath := Home{NoNeedLogin: []string{}, NoNeedAuths: []string{"*"}}
	gf.Register(&fpath, fpath)
}

func parseStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch vv := v.(type) {
	case []string:
		return vv
	case []interface{}:
		out := make([]string, 0, len(vv))
		for _, it := range vv {
			s := strings.TrimSpace(gconv.String(it))
			if s != "" {
				out = append(out, s)
			}
		}
		return out
	default:
		s := strings.TrimSpace(gconv.String(v))
		if s == "" {
			return nil
		}
		// 兼容前端把数组以 JSON 字符串传入（例如 "[]"、"['a','b']"/'["a","b"]'）
		if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			var arr []string
			if err := json.Unmarshal([]byte(s), &arr); err == nil {
				out := make([]string, 0, len(arr))
				for _, it := range arr {
					it = strings.TrimSpace(it)
					if it != "" {
						out = append(out, it)
					}
				}
				return out
			}
			var anyArr []interface{}
			if err := json.Unmarshal([]byte(s), &anyArr); err == nil {
				out := make([]string, 0, len(anyArr))
				for _, it := range anyArr {
					ss := strings.TrimSpace(gconv.String(it))
					if ss != "" {
						out = append(out, ss)
					}
				}
				return out
			}
		}
		parts := strings.Split(s, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				out = append(out, p)
			}
		}
		return out
	}
}

func fullImgURL(path string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	return gf.GetFullUrl(p)
}

func saleStatusLabel(status string) string {
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

// GetList 获取首页数据
// 入参（query）：page,pageSize,keyword,filter,sort,location,more
// 返回：bannerImages、followedTotal、followedProperties、recommendedProperties
// GetList 获取首页数据，包括banner图片、用户关注的房源和推荐房源列表
// 支持分页查询、关键词搜索、筛选和排序功能
func (api *Home) GetList(c *gf.GinCtx) {
	// 获取业务ID，兼容未登录情况
	businessID := c.GetInt64("businessID")
	if businessID == 0 {
		businessID = gf.Int64(c.GetHeader("Businessid"))
		if businessID == 0 {
			businessID = 1
		}
	}
	userID := c.GetInt64("userID")

	// 分页参数处理
	pageNo := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("pageSize", "10"))
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}

	// 查询参数处理，优先从query获取，其次兼容body
	param, _ := gf.RequestParam(c)
	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		keyword = strings.TrimSpace(gconv.String(param["keyword"]))
	}
	filter := strings.TrimSpace(c.Query("filter"))
	if filter == "" {
		filter = strings.TrimSpace(gconv.String(param["filter"]))
	}
	sort := strings.TrimSpace(c.Query("sort"))
	if sort == "" {
		sort = strings.TrimSpace(gconv.String(param["sort"]))
	}
	location := strings.TrimSpace(c.Query("location"))
	if location == "" {
		location = strings.TrimSpace(gconv.String(param["location"]))
	}
	more := parseStringSlice(param["more"])
	if len(more) == 0 {
		more = parseStringSlice(c.Query("more"))
	}

	// Banner图片配置
	bannerImages := []string{
		"/static/images/img_28f6b412d7.png",
		"/static/images/img_2abd9934e1.png",
		"/static/images/img_48afe8538f.png",
	}

	// 获取用户关注的房源（仅登录用户）
	followedTotal := 0
	followedProperties := make([]gf.Map, 0)
	if userID > 0 {
		favMDB := gf.Model("business_favorites").Where("user_id", userID)
		followedTotal, _ = favMDB.Clone().Count()
		favs, _ := favMDB.Fields("property_id").Order("id desc").Page(1, 5).Select()
		propIDs := make([]interface{}, 0, len(favs))
		for _, row := range favs {
			if row != nil && row["property_id"].Int64() > 0 {
				propIDs = append(propIDs, row["property_id"].Int64())
			}
		}
		if len(propIDs) > 0 {
			props, _ := gf.Model("business_properties").
				Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,property_type,community_name,address,tags,has_smart_lock,sale_status,view_count,follow_count,cover_image,images").
				Where("business_id", businessID).
				Where("status", 0).
				Where("deletetime", nil).
				WhereIn("id", propIDs).
				Order("weigh desc, id desc").
				Select()
			for _, p := range props {
				if p == nil {
					continue
				}
				img := fullImgURL(p["cover_image"].String())
				if img == "" && p["images"] != nil && strings.TrimSpace(p["images"].String()) != "" {
					first := strings.Split(p["images"].String(), ",")[0]
					img = fullImgURL(first)
				}
				tags := make([]string, 0)
				if p["tags"] != nil && strings.TrimSpace(p["tags"].String()) != "" {
					tags = gf.SplitAndStr(p["tags"].String(), ",")
				}
				ss := p["sale_status"].String()
				followedProperties = append(followedProperties, gf.Map{
					"id":   p["id"].Int64(),
					"name": p["title"].String(),
					// 兼容旧字段：area=小区/区域名，size=面积（平米）
					"area":        p["community_name"].String(),
					"size":        gconv.String(p["area"]),
					"price":       gconv.String(p["price"]),
					"price_unit":  p["price_unit"].String(),
					"type":        p["property_type"].String(),
					"priceChange": "",
					// 兼容旧字段：status 用于卡片角标展示
					"status":            saleStatusLabel(ss),
					"sale_status":       ss,
					"sale_status_label": saleStatusLabel(ss),
					"has_smart_lock":    p["has_smart_lock"].Int(),
					"rooms":             p["rooms"].Int(),
					"halls":             p["halls"].Int(),
					"bathrooms":         p["bathrooms"].Int(),
					"property_type":     p["property_type"].String(),
					"community_name":    p["community_name"].String(),
					"address":           p["address"].String(),
					"tags":              tags,
					"view_count":        p["view_count"].Int(),
					"follow_count":      p["follow_count"].Int(),
					"image":             img,
				})
			}
		}
	}

	// 构建推荐房源查询条件
	recMDB := gf.Model("business_properties").
		Where("business_id", businessID).
		Where("hot_status", 1).
		Where("status", 0).
		Where("deletetime", nil)

	// 关键词搜索：标题/小区/地址
	if keyword != "" {
		kw := "%" + keyword + "%"
		recMDB = recMDB.Where("(title like ? OR community_name like ? OR address like ?)", kw, kw, kw)
	}

	// location参数暂不使用，保留兼容性
	_ = location

	// filter筛选条件处理
	if filter == "smart_lock" {
		recMDB = recMDB.Where("has_smart_lock", 1)
	} else if filter == "new" {
		// 新房源默认按创建时间排序
		if sort == "" {
			sort = "latest"
		}
	} else if filter != "" {
		// 其他filter通过tags匹配
		recMDB = recMDB.Where("tags like ?", "%"+filter+"%")
	}

	// more标签条件处理
	for _, m := range more {
		if m == "" {
			continue
		}
		recMDB = recMDB.Where("tags like ?", "%"+m+"%")
	}

	// 排序处理
	switch sort {
	case "price_asc":
		recMDB = recMDB.Order("price asc, weigh desc, id desc")
	case "price_desc":
		recMDB = recMDB.Order("price desc, weigh desc, id desc")
	case "latest":
		recMDB = recMDB.Order("createtime desc, id desc")
	case "recommend", "":
		recMDB = recMDB.Order("weigh desc, id desc")
	default:
		recMDB = recMDB.Order("weigh desc, id desc")
	}

	// 获取推荐房源列表
	recs, err := recMDB.
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,property_type,community_name,address,tags,has_smart_lock,sale_status,hot_status,view_count,follow_count,cover_image,images,createtime").
		Page(pageNo, pageSize).
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取推荐房源失败：" + err.Error()).Regin(c)
		return
	}

	// 格式化推荐房源数据
	recommendedProperties := make([]gf.Map, 0, len(recs))
	for _, p := range recs {
		if p == nil {
			continue
		}
		img := fullImgURL(p["cover_image"].String())
		if img == "" && p["images"] != nil && strings.TrimSpace(p["images"].String()) != "" {
			first := strings.Split(p["images"].String(), ",")[0]
			img = fullImgURL(first)
		}
		tags := make([]string, 0)
		if p["tags"] != nil && strings.TrimSpace(p["tags"].String()) != "" {
			tags = gf.SplitAndStr(p["tags"].String(), ",")
		}
		region := strings.TrimSpace(p["community_name"].String())
		if region == "" {
			region = strings.TrimSpace(p["address"].String())
		}
		ss := p["sale_status"].String()
		recommendedProperties = append(recommendedProperties, gf.Map{
			"id":         p["id"].Int64(),
			"name":       p["title"].String(),
			"price":      gconv.String(p["price"]),
			"price_unit": p["price_unit"].String(),
			"region":     region,
			"size":       gconv.String(p["area"]),
			"desc":       strings.TrimSpace(p["property_type"].String()),
			// 兼容旧字段：status 用于卡片角标展示
			"status":            saleStatusLabel(ss),
			"sale_status":       ss,
			"sale_status_label": saleStatusLabel(ss),
			"hot_status":        p["hot_status"].Int(),
			"has_smart_lock":    p["has_smart_lock"].Int(),
			"rooms":             p["rooms"].Int(),
			"halls":             p["halls"].Int(),
			"bathrooms":         p["bathrooms"].Int(),
			"property_type":     p["property_type"].String(),
			"community_name":    p["community_name"].String(),
			"address":           p["address"].String(),
			"tags":              tags,
			"view_count":        p["view_count"].Int(),
			"follow_count":      p["follow_count"].Int(),
			"image":             img,
		})
	}

	// 返回首页数据
	gf.Success().SetMsg("获取首页数据").SetData(gf.Map{
		"bannerImages":          bannerImages,
		"followedTotal":         followedTotal,
		"followedProperties":    followedProperties,
		"recommendedProperties": recommendedProperties,
	}).Regin(c)
}
