"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const BottomTabBar = () => "../../components/BottomTabBar.js";
const _sfc_main = {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      items: [
        { title: "日落大道 2室紧凑型", time: "2026-01-07 12:01", count: 3 },
        { title: "滨江别墅 独栋花园", time: "2026-01-06 20:18", count: 1 },
        { title: "金域华府 A2", time: "2026-01-05 08:33", count: 2 }
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
      title: "浏览记录"
    }),
    b: common_vendor.f($data.items, (it, idx, i0) => {
      return {
        a: common_vendor.t(it.title),
        b: common_vendor.t(it.count),
        c: common_vendor.t(it.time),
        d: idx
      };
    }),
    c: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/records/record_view.js.map
