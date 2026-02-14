package uniapp

import "gofly/utils/gform"

// wxHiddenSaleStatuses 小程序端默认列表/入口中需要完全不可见的销售状态。
//
// 说明：
// - 业务端后台仍可查看
// - 这里的过滤仅用于 uniapp（小程序端）相关接口
// - 预售(in_sale) 仅允许在“预售 Tab”中展示，其他入口一律隐藏
func wxHiddenSaleStatuses() []string {
	return []string{"sold", "off_market", "in_sale"}
}

// wxApplyPropertyVisibility 对房源查询施加“小程序端可见性”过滤。
func wxApplyPropertyVisibility(m *gform.Model) *gform.Model {
	if m == nil {
		return m
	}
	return m.WhereNotIn("sale_status", wxHiddenSaleStatuses())
}
