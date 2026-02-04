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
