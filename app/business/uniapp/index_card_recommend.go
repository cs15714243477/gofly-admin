package uniapp

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"gofly/utils/gf"
	"gofly/utils/gform"
	"gofly/utils/tools/gconv"
	"gofly/utils/tools/gtime"
)

const (
	agentCardRecommendMaxCount = 20
	agentCardRecommendShowSize = 6
)

var (
	agentCardRecommendTableOnce sync.Once
	agentCardRecommendTableErr  error
)

func ensureAgentCardRecommendTable() error {
	agentCardRecommendTableOnce.Do(func() {
		sql := `
CREATE TABLE IF NOT EXISTS business_agent_card_recommend (
  id int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  business_id int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  agent_id int(11) NOT NULL DEFAULT 0 COMMENT '经纪人ID',
  property_id int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  sort int(11) NOT NULL DEFAULT 0 COMMENT '排序',
  status tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  createtime datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (id) USING BTREE,
  UNIQUE KEY uk_biz_agent_property (business_id, agent_id, property_id),
  KEY idx_biz_agent_sort (business_id, agent_id, sort, id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='经纪人名片推荐房源'
`
		_, agentCardRecommendTableErr = gf.DB().Exec(context.Background(), sql)
	})
	return agentCardRecommendTableErr
}

func isPublicDetailView(raw string) bool {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

func normalizePropertyIDList(raw interface{}, maxCount int) []int64 {
	if maxCount <= 0 {
		maxCount = agentCardRecommendMaxCount
	}
	values := make([]interface{}, 0)
	switch v := raw.(type) {
	case []interface{}:
		values = append(values, v...)
	case []int64:
		for _, id := range v {
			values = append(values, id)
		}
	case []int:
		for _, id := range v {
			values = append(values, id)
		}
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			values = nil
		} else if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
			var arr []interface{}
			if err := json.Unmarshal([]byte(s), &arr); err == nil {
				values = append(values, arr...)
			} else {
				values = append(values, s)
			}
		} else {
			parts := strings.Split(s, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				values = append(values, part)
			}
		}
	default:
		s := strings.TrimSpace(gconv.String(raw))
		if s != "" {
			values = append(values, s)
		}
	}

	out := make([]int64, 0, len(values))
	seen := make(map[int64]struct{}, len(values))
	for _, each := range values {
		if len(out) >= maxCount {
			break
		}
		id := gconv.Int64(each)
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func formatCardPropertyRow(row gform.Record) gf.Map {
	if row == nil {
		return nil
	}
	img := wxFullImgURL(row["cover_image"].String())
	if img == "" && row["images"] != nil && strings.TrimSpace(row["images"].String()) != "" {
		first := strings.Split(row["images"].String(), ",")[0]
		img = wxFullImgURL(first)
	}
	status := row["sale_status"].String()
	return gf.Map{
		"id":                row["id"].Int64(),
		"title":             row["title"].String(),
		"price":             gconv.String(row["price"]),
		"price_unit":        row["price_unit"].String(),
		"area":              gconv.String(row["area"]),
		"rooms":             row["rooms"].Int(),
		"halls":             row["halls"].Int(),
		"bathrooms":         row["bathrooms"].Int(),
		"community_name":    row["community_name"].String(),
		"address":           row["address"].String(),
		"image":             img,
		"cover_image":       img,
		"tags":              wxSplitComma(row["tags"].String()),
		"sale_status":       status,
		"sale_status_label": wxSaleStatusLabel(status),
		"hot_status":        row["hot_status"].Int(),
	}
}

func queryAgentRecommendedProperties(businessID, agentID int64, limit int) ([]gf.Map, []int64, error) {
	if limit <= 0 {
		limit = agentCardRecommendShowSize
	}
	if err := ensureAgentCardRecommendTable(); err != nil {
		return nil, nil, err
	}
	rows, err := gf.ModelRaw(`
SELECT
	p.id,p.title,p.price,p.price_unit,p.area,p.rooms,p.halls,p.bathrooms,
	p.community_name,p.address,p.cover_image,p.images,p.tags,p.sale_status,p.hot_status,
	r.sort
FROM business_agent_card_recommend r
INNER JOIN business_properties p ON p.id = r.property_id
WHERE
	r.business_id = ? AND r.agent_id = ? AND r.status = 0
	AND p.business_id = ? AND p.status = 0 AND p.deletetime IS NULL
ORDER BY r.sort ASC, r.id ASC
LIMIT ?
`, businessID, agentID, businessID, limit).Select()
	if err != nil {
		return nil, nil, err
	}
	items := make([]gf.Map, 0, len(rows))
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		item := formatCardPropertyRow(row)
		if item == nil {
			continue
		}
		id := gconv.Int64(item["id"])
		if id <= 0 {
			continue
		}
		items = append(items, item)
		ids = append(ids, id)
	}
	return items, ids, nil
}

func querySystemRecommendedProperties(businessID int64, limit int, exclude map[int64]struct{}) ([]gf.Map, error) {
	if limit <= 0 {
		limit = agentCardRecommendShowSize
	}
	items := make([]gf.Map, 0, limit)
	seen := make(map[int64]struct{}, limit*3)
	appendRows := func(rows gform.Result) {
		if len(items) >= limit {
			return
		}
		for _, row := range rows {
			if len(items) >= limit {
				return
			}
			item := formatCardPropertyRow(row)
			if item == nil {
				continue
			}
			id := gconv.Int64(item["id"])
			if id <= 0 {
				continue
			}
			if exclude != nil {
				if _, ok := exclude[id]; ok {
					continue
				}
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			items = append(items, item)
		}
	}

	hotRows, err := gf.Model("business_properties").
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,community_name,address,cover_image,images,tags,sale_status,hot_status").
		Where("business_id", businessID).
		Where("status", 0).
		Where("deletetime", nil).
		Where("hot_status", 1).
		Order("weigh desc,id desc").
		Limit(limit * 2).
		Select()
	if err != nil {
		return nil, err
	}
	appendRows(hotRows)
	if len(items) >= limit {
		return items, nil
	}

	normalRows, err := gf.Model("business_properties").
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,community_name,address,cover_image,images,tags,sale_status,hot_status").
		Where("business_id", businessID).
		Where("status", 0).
		Where("deletetime", nil).
		Order("weigh desc,id desc").
		Limit(limit * 3).
		Select()
	if err != nil {
		return nil, err
	}
	appendRows(normalRows)

	return items, nil
}

func filterAvailablePropertyIDs(businessID int64, ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return []int64{}, nil
	}
	rows, err := gf.Model("business_properties").
		Fields("id").
		Where("business_id", businessID).
		Where("status", 0).
		Where("deletetime", nil).
		WhereIn("id", ids).
		Select()
	if err != nil {
		return nil, err
	}
	okMap := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		id := row["id"].Int64()
		if id > 0 {
			okMap[id] = struct{}{}
		}
	}
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if _, ok := okMap[id]; ok {
			out = append(out, id)
		}
	}
	return out, nil
}

// GetAgentCardRecommendConfig 获取名片推荐房源配置（经纪人可编辑）
func (api *Index) GetAgentCardRecommendConfig(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	businessID := c.GetInt64("businessID")
	if userID <= 0 || businessID <= 0 {
		gf.Failed().SetMsg("请先登录").Regin(c)
		return
	}
	if err := ensureAgentCardRecommendTable(); err != nil {
		gf.Failed().SetMsg("初始化推荐配置失败：" + err.Error()).Regin(c)
		return
	}

	selectedProperties, selectedIDs, err := queryAgentRecommendedProperties(businessID, userID, agentCardRecommendMaxCount)
	if err != nil {
		gf.Failed().SetMsg("获取经纪人推荐失败：" + err.Error()).Regin(c)
		return
	}
	selectedIDMap := make(map[int64]struct{}, len(selectedIDs))
	for _, id := range selectedIDs {
		selectedIDMap[id] = struct{}{}
	}

	candidateRows, err := gf.Model("business_properties").
		Fields("id,title,price,price_unit,area,rooms,halls,bathrooms,community_name,address,cover_image,images,tags,sale_status,hot_status,weigh,createtime").
		Where("business_id", businessID).
		Where("status", 0).
		Where("deletetime", nil).
		Order("hot_status desc,weigh desc,id desc").
		Limit(200).
		Select()
	if err != nil {
		gf.Failed().SetMsg("获取候选房源失败：" + err.Error()).Regin(c)
		return
	}
	candidates := make([]gf.Map, 0, len(candidateRows))
	for _, row := range candidateRows {
		item := formatCardPropertyRow(row)
		if item == nil {
			continue
		}
		id := gconv.Int64(item["id"])
		_, checked := selectedIDMap[id]
		item["selected"] = checked
		candidates = append(candidates, item)
	}

	gf.Success().SetMsg("获取推荐配置").SetData(gf.Map{
		"max_count":           agentCardRecommendMaxCount,
		"selected_ids":        selectedIDs,
		"selected_properties": selectedProperties,
		"candidates":          candidates,
	}).Regin(c)
}

// SaveAgentCardRecommendConfig 保存名片推荐房源配置（经纪人可编辑）
// 入参：property_ids（数组或逗号分隔字符串）
func (api *Index) SaveAgentCardRecommendConfig(c *gf.GinCtx) {
	userID := c.GetInt64("userID")
	businessID := c.GetInt64("businessID")
	if userID <= 0 || businessID <= 0 {
		gf.Failed().SetMsg("请先登录").Regin(c)
		return
	}
	if err := ensureAgentCardRecommendTable(); err != nil {
		gf.Failed().SetMsg("初始化推荐配置失败：" + err.Error()).Regin(c)
		return
	}

	param, _ := gf.RequestParam(c)
	rawIDs, ok := param["property_ids"]
	if !ok {
		rawIDs, ok = param["propertyIds"]
	}
	if !ok {
		rawIDs = ""
	}

	propertyIDs := normalizePropertyIDList(rawIDs, agentCardRecommendMaxCount)
	propertyIDs, err := filterAvailablePropertyIDs(businessID, propertyIDs)
	if err != nil {
		gf.Failed().SetMsg("校验房源失败：" + err.Error()).Regin(c)
		return
	}

	reqCtx := context.Background()
	if c != nil && c.Request != nil {
		reqCtx = c.Request.Context()
	}
	err = gf.Model("business_agent_card_recommend").Ctx(reqCtx).Transaction(reqCtx, func(ctx context.Context, tx gform.TX) error {
		if _, err := gf.Model("business_agent_card_recommend").TX(tx).Ctx(ctx).
			Where("business_id", businessID).
			Where("agent_id", userID).
			Delete(); err != nil {
			return err
		}
		if len(propertyIDs) == 0 {
			return nil
		}
		for i, propertyID := range propertyIDs {
			if _, err := gf.Model("business_agent_card_recommend").TX(tx).Ctx(ctx).Data(gf.Map{
				"business_id": businessID,
				"agent_id":    userID,
				"property_id": propertyID,
				"sort":        i + 1,
				"status":      0,
				"createtime":  gtime.Datetime(),
			}).Insert(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		gf.Failed().SetMsg("保存推荐配置失败：" + err.Error()).Regin(c)
		return
	}

	gf.Success().SetMsg("保存成功").SetData(gf.Map{
		"property_ids": propertyIDs,
		"count":        len(propertyIDs),
		"max_count":    agentCardRecommendMaxCount,
	}).Regin(c)
}
