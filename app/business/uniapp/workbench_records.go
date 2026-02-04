package uniapp

import (
	"encoding/json"
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

const (
	wxRecordTypeFollow  = "follow"
	wxRecordTypeUnlock  = "unlock"
	wxRecordTypeShowing = "showing"
	wxRecordTypeView    = "view"
	wxRecordTypeShare   = "share"
	wxRecordTypeCall    = "call"
)

func wxBuildPropertyCover(cover string, images string) string {
	if u := wxFullImgURL(cover); u != "" {
		return u
	}
	parts := wxSplitComma(images)
	if len(parts) == 0 {
		return ""
	}
	return wxFullImgURL(parts[0])
}

func wxBuildLayoutText(rooms, halls, bathrooms int) string {
	if rooms <= 0 && halls <= 0 && bathrooms <= 0 {
		return ""
	}
	return gconv.String(rooms) + "室" + gconv.String(halls) + "厅" + gconv.String(bathrooms) + "卫"
}

func wxParseMeta(raw string) map[string]interface{} {
	out := map[string]interface{}{}
	s := strings.TrimSpace(raw)
	if s == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}

func wxMetaString(meta map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if meta == nil {
			return ""
		}
		if v, ok := meta[k]; ok {
			if s := strings.TrimSpace(gconv.String(v)); s != "" {
				return s
			}
		}
	}
	return ""
}

func wxMetaInt(meta map[string]interface{}, keys ...string) int {
	for _, k := range keys {
		if meta == nil {
			return 0
		}
		if v, ok := meta[k]; ok {
			n := gconv.Int(v)
			if n > 0 {
				return n
			}
		}
	}
	return 0
}

func wxParsePage(c *gf.GinCtx, param map[string]interface{}) (int, int) {
	page := gconv.Int(c.DefaultQuery("page", "1"))
	pageSize := gconv.Int(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	if param != nil {
		if v, ok := param["page"]; ok {
			if n := gconv.Int(v); n > 0 {
				page = n
			}
		}
		if v, ok := param["page_size"]; ok {
			if n := gconv.Int(v); n > 0 {
				pageSize = n
			}
		}
		if pageSize > 50 {
			pageSize = 50
		}
	}
	return page, pageSize
}

func wxUnlockStatusText(status string) string {
	switch strings.TrimSpace(status) {
	case "approved", "completed":
		return "成功"
	case "rejected", "cancelled":
		return "失败"
	default:
		return "处理中"
	}
}

func wxUnlockMethodText(method string) string {
	switch strings.TrimSpace(method) {
	case "bluetooth":
		return "智能开锁"
	case "password":
		return "密码开锁"
	default:
		return "开锁"
	}
}

// GetWorkbenchSummary 获取经纪人工作台记录汇总数据
func (api *Index) GetWorkbenchSummary(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}

	followCount, _ := gf.Model("business_favorites").Where("user_id", userID).Count()
	unlockCount, _ := gf.Model("business_unlock_requests").Where("user_id", userID).Count()

	showingCount := 0
	viewCount := 0
	shareCount := 0
	callCount := 0
	typeRows, _ := gf.Model("business_user_activity_logs").
		Fields("activity_type,count(1) as cnt").
		Where("user_id", userID).
		WhereIn("activity_type", []string{wxRecordTypeShowing, wxRecordTypeView, wxRecordTypeShare, wxRecordTypeCall}).
		Group("activity_type").
		Select()
	for _, row := range typeRows {
		if row == nil {
			continue
		}
		switch strings.TrimSpace(row["activity_type"].String()) {
		case wxRecordTypeShowing:
			showingCount = row["cnt"].Int()
		case wxRecordTypeView:
			viewCount = row["cnt"].Int()
		case wxRecordTypeShare:
			shareCount = row["cnt"].Int()
		case wxRecordTypeCall:
			callCount = row["cnt"].Int()
		}
	}

	unlockNotice := false
	lastUnlock, _ := gf.Model("business_unlock_requests").
		Fields("request_status").
		Where("user_id", userID).
		Order("id desc").
		Find()
	if lastUnlock != nil {
		unlockNotice = strings.TrimSpace(lastUnlock["request_status"].String()) == "rejected"
	}

	gf.Success().SetMsg("获取工作台汇总").SetData(gf.Map{
		"follow_count":       followCount,
		"unlock_count":       unlockCount,
		"showing_count":      showingCount,
		"view_count":         viewCount,
		"share_count":        shareCount,
		"call_count":         callCount,
		"unlock_has_notice":  unlockNotice,
		"record_types":       []string{wxRecordTypeFollow, wxRecordTypeUnlock, wxRecordTypeShowing, wxRecordTypeView, wxRecordTypeShare, wxRecordTypeCall},
		"has_any_record":     followCount+unlockCount+showingCount+viewCount+shareCount+callCount > 0,
		"has_unlock_records": unlockCount > 0,
	}).Regin(c)
}

// GetWorkbenchRecords 获取经纪人工作台记录列表
// 入参：
// - record_type: follow|unlock|showing|view|share|call
// - page/page_size
func (api *Index) GetWorkbenchRecords(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	if userID <= 0 {
		gf.Failed().SetCode(401).SetMsg("请先登录").Regin(c)
		return
	}
	businessID := wxBusinessID(c)

	param, _ := gf.RequestParam(c)
	recordType := strings.TrimSpace(c.DefaultQuery("record_type", ""))
	if recordType == "" {
		recordType = strings.TrimSpace(gconv.String(param["record_type"]))
	}
	pageNo, pageSize := wxParsePage(c, param)

	items := make([]gf.Map, 0)
	total := 0

	switch recordType {
	case wxRecordTypeFollow:
		mdb := gf.Model("business_favorites f").
			LeftJoin("business_properties p", "p.id=f.property_id AND p.deletetime IS NULL").
			Fields("f.id,f.property_id,f.createtime,p.business_id,p.title,p.community_name,p.cover_image,p.images,p.area,p.rooms,p.halls,p.bathrooms").
			Where("f.user_id", userID)
		total, _ = mdb.Clone().Count()
		rows, err := mdb.Page(pageNo, pageSize).Order("f.id desc").Select()
		if err != nil {
			gf.Failed().SetMsg("获取关注记录失败：" + err.Error()).Regin(c)
			return
		}
		for _, row := range rows {
			if row == nil {
				continue
			}
			if bid := row["business_id"].Int64(); bid > 0 && bid != businessID {
				continue
			}
			title := strings.TrimSpace(row["title"].String())
			if title == "" {
				title = strings.TrimSpace(row["community_name"].String())
			}
			if title == "" {
				title = "房源"
			}
			layout := wxBuildLayoutText(row["rooms"].Int(), row["halls"].Int(), row["bathrooms"].Int())
			area := strings.TrimSpace(gconv.String(row["area"]))
			subParts := make([]string, 0, 3)
			if v := strings.TrimSpace(row["community_name"].String()); v != "" {
				subParts = append(subParts, v)
			}
			if layout != "" {
				subParts = append(subParts, layout)
			}
			if area != "" {
				subParts = append(subParts, area+"㎡")
			}
			items = append(items, gf.Map{
				"id":          row["id"].Int64(),
				"property_id": row["property_id"].Int64(),
				"title":       title,
				"time":        row["createtime"].String(),
				"sub":         strings.Join(subParts, " · "),
				"image":       wxBuildPropertyCover(row["cover_image"].String(), row["images"].String()),
			})
		}
	case wxRecordTypeUnlock:
		mdb := gf.Model("business_unlock_requests u").
			LeftJoin("business_properties p", "p.id=u.property_id AND p.deletetime IS NULL").
			Fields("u.id,u.property_id,u.request_status,u.request_type,u.request_time,u.updatetime,p.business_id,p.title,p.community_name,p.cover_image,p.images").
			Where("u.user_id", userID)
		total, _ = mdb.Clone().Count()
		rows, err := mdb.Page(pageNo, pageSize).Order("u.id desc").Select()
		if err != nil {
			gf.Failed().SetMsg("获取开锁记录失败：" + err.Error()).Regin(c)
			return
		}
		for _, row := range rows {
			if row == nil {
				continue
			}
			if bid := row["business_id"].Int64(); bid > 0 && bid != businessID {
				continue
			}
			title := strings.TrimSpace(row["title"].String())
			if title == "" {
				title = strings.TrimSpace(row["community_name"].String())
			}
			if title == "" {
				title = "房源"
			}
			timeText := strings.TrimSpace(row["updatetime"].String())
			if timeText == "" {
				timeText = strings.TrimSpace(row["request_time"].String())
			}
			statusText := wxUnlockStatusText(row["request_status"].String())
			items = append(items, gf.Map{
				"id":          row["id"].Int64(),
				"property_id": row["property_id"].Int64(),
				"title":       title,
				"time":        timeText,
				"method":      wxUnlockMethodText(row["request_type"].String()),
				"status":      statusText,
				"image":       wxBuildPropertyCover(row["cover_image"].String(), row["images"].String()),
			})
		}
	case wxRecordTypeShowing, wxRecordTypeView, wxRecordTypeShare, wxRecordTypeCall:
		mdb := gf.Model("business_user_activity_logs a").
			LeftJoin("business_properties p", "p.id=a.property_id AND p.deletetime IS NULL").
			Fields("a.id,a.property_id,a.meta_data,a.createtime,p.business_id,p.title,p.community_name,p.cover_image,p.images").
			Where("a.user_id", userID).
			Where("a.activity_type", recordType)
		total, _ = mdb.Clone().Count()
		rows, err := mdb.Page(pageNo, pageSize).Order("a.id desc").Select()
		if err != nil {
			gf.Failed().SetMsg("获取记录失败：" + err.Error()).Regin(c)
			return
		}
		for _, row := range rows {
			if row == nil {
				continue
			}
			if bid := row["business_id"].Int64(); bid > 0 && bid != businessID {
				continue
			}
			meta := wxParseMeta(row["meta_data"].String())
			title := strings.TrimSpace(row["title"].String())
			if title == "" {
				title = strings.TrimSpace(row["community_name"].String())
			}
			if title == "" {
				title = wxMetaString(meta, "title", "property_title", "name", "contact_name")
			}
			if title == "" {
				title = "记录"
			}
			item := gf.Map{
				"id":          row["id"].Int64(),
				"property_id": row["property_id"].Int64(),
				"title":       title,
				"time":        row["createtime"].String(),
				"image":       wxBuildPropertyCover(row["cover_image"].String(), row["images"].String()),
			}
			switch recordType {
			case wxRecordTypeShowing:
				client := wxMetaString(meta, "client", "client_name", "customer", "name")
				if client == "" {
					client = "客户"
				}
				showType := wxMetaString(meta, "showing_type", "type", "channel")
				if showType == "" {
					showType = "线下"
				}
				item["client"] = client
				item["type"] = showType
			case wxRecordTypeView:
				viewCount := wxMetaInt(meta, "count", "view_count", "times")
				if viewCount <= 0 {
					viewCount = 1
				}
				item["count"] = viewCount
			case wxRecordTypeShare:
				channel := wxMetaString(meta, "channel", "type", "from")
				if channel == "" {
					channel = "微信"
				}
				clicks := wxMetaInt(meta, "clicks", "click_count", "count")
				item["channel"] = channel
				item["clicks"] = clicks
			case wxRecordTypeCall:
				callType := wxMetaString(meta, "type", "direction")
				if callType == "" {
					callType = "呼出"
				}
				phone := wxMetaString(meta, "phone", "mobile", "target_phone")
				name := wxMetaString(meta, "contact_name", "name", "target_name")
				if name == "" {
					name = title
				}
				item["type"] = callType
				item["phone"] = phone
				item["name"] = name
			}
			items = append(items, item)
		}
		// 兜底：若“带看/浏览”无行为日志，则回退展示当前经纪人房源累计统计
		if len(items) == 0 && (recordType == wxRecordTypeShowing || recordType == wxRecordTypeView) {
			countField := "view_count"
			if recordType == wxRecordTypeShowing {
				countField = "showing_count"
			}
			fallbackMDB := gf.Model("business_properties").
				Fields("id,title,community_name,cover_image,images,view_count,showing_count,updatetime").
				Where("business_id", businessID).
				Where("agent_id", userID).
				Where("deletetime", nil).
				Where("status", 0).
				Where(countField+" > ?", 0)
			total, _ = fallbackMDB.Clone().Count()
			fallbackRows, ferr := fallbackMDB.Page(pageNo, pageSize).Order(countField + " desc,id desc").Select()
			if ferr == nil {
				for _, row := range fallbackRows {
					if row == nil {
						continue
					}
					title := strings.TrimSpace(row["title"].String())
					if title == "" {
						title = strings.TrimSpace(row["community_name"].String())
					}
					if title == "" {
						title = "房源"
					}
					item := gf.Map{
						"id":          row["id"].Int64(),
						"property_id": row["id"].Int64(),
						"title":       title,
						"time":        strings.TrimSpace(row["updatetime"].String()),
						"image":       wxBuildPropertyCover(row["cover_image"].String(), row["images"].String()),
					}
					if recordType == wxRecordTypeShowing {
						item["type"] = "累计带看"
						item["client"] = "累计" + gconv.String(row["showing_count"].Int()) + "次"
					} else {
						item["count"] = row["view_count"].Int()
					}
					items = append(items, item)
				}
			}
		}
	default:
		gf.Failed().SetMsg("record_type 参数无效").Regin(c)
		return
	}

	gf.Success().SetMsg("获取记录成功").SetData(gf.Map{
		"record_type": recordType,
		"page":        pageNo,
		"page_size":   pageSize,
		"total":       total,
		"items":       items,
	}).Regin(c)
}
