"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const BottomTabBar = () => "../../components/BottomTabBar.js";
const _sfc_main = {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      items: [
        { title: "阳光花园三期 2室1厅", time: "2026-01-07 14:20" },
        { title: "金地名津 精装大三房", time: "2026-01-06 19:02" },
        { title: "万科四季花城", time: "2026-01-05 10:11" }
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
      title: "关注记录"
    }),
    b: common_vendor.f($data.items, (it, idx, i0) => {
      return {
        a: common_vendor.t(it.title),
        b: common_vendor.t(it.time),
        c: idx
      };
    }),
    c: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/records/record_follow.js.map
