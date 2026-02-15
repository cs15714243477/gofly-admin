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

// wxManageHiddenSaleStatuses 小程序端“可维护房源”需要隐藏的销售状态。
//
// 说明：
// - 预售(in_sale) 允许维护者编辑，因此不在此隐藏范围内。
// - 下架(off_market) 允许维护者重新上架/编辑，因此不在此隐藏范围内。
// - 已售(sold) 作为终态，在小程序维护入口中保持只读（不允许编辑/删除）。
func wxManageHiddenSaleStatuses() []string {
	return []string{"sold"}
}
