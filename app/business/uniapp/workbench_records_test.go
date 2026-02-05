package uniapp

import "testing"

func TestWxMetaStringAndIntSupportNestedMeta(t *testing.T) {
	meta := map[string]interface{}{
		"page": "property_detail",
		"meta": map[string]interface{}{
			"client_name":  "张三",
			"client_phone": "13800138000",
			"count":        3,
		},
	}

	if got := wxMetaString(meta, "client_name"); got != "张三" {
		t.Fatalf("wxMetaString nested client_name=%q, want=%q", got, "张三")
	}
	if got := wxMetaString(meta, "client_phone"); got != "13800138000" {
		t.Fatalf("wxMetaString nested client_phone=%q, want=%q", got, "13800138000")
	}
	if got := wxMetaInt(meta, "count"); got != 3 {
		t.Fatalf("wxMetaInt nested count=%d, want=%d", got, 3)
	}
}

