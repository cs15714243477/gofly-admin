package uniapp

import (
	"encoding/json"
	"strings"

	"gofly/utils/gf"
	"gofly/utils/tools/gconv"
)

// normalizeCommaText 把任意输入规范化为“逗号分隔字符串”
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

// normalizeTagsToString 兼容 tags/materials/images 等字段输入格式
//
// 支持：
// - []string / []interface{}
// - JSON 数组字符串（如 '["a","b"]'）
// - 逗号分隔字符串
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

// pickMap 从 param 中挑选允许写入的字段
func pickMap(param map[string]interface{}, keys ...string) gf.Map {
	out := gf.Map{}
	for _, k := range keys {
		if v, ok := param[k]; ok {
			out[k] = v
		}
	}
	return out
}
