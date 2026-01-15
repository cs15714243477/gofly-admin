"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const BottomTabBar = () => "../../components/BottomTabBar.js";
const _sfc_main = {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      items: [
        { title: "滨江别墅 独栋花园", time: "2026-01-07 13:40", client: "张先生", type: "线下" },
        { title: "阳光花园三期 2室", time: "2026-01-06 16:20", client: "王女士", type: "自助开锁" },
        { title: "金城中心 单身公寓", time: "2026-01-05 10:30", client: "李先生", type: "线上" }
      ]
    };
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  const _component_BottomTabBar = common_vendor.resolveComponent("BottomTabBar");
  (_component_TopHeader + _component_BottomTabBar)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.p({
      title: "带看记录"
    }),
    b: common_vendor.f($data.items, (it, idx, i0) => {
      return {
        a: common_vendor.t(it.title),
        b: common_vendor.t(it.type),
        c: common_vendor.t(it.time),
        d: common_vendor.t(it.client),
        e: idx
      };
    }),
    c: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/records/record_showing.js.map
