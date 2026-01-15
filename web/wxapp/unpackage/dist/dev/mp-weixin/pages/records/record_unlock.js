"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const BottomTabBar = () => "../../components/BottomTabBar.js";
const _sfc_main = {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      items: [
        { title: "阳光花园 3栋 201室", time: "2026-01-07 15:10", method: "智能开锁", status: "成功" },
        { title: "金域华府 A2", time: "2026-01-06 11:32", method: "密码开锁", status: "成功" },
        { title: "万科四季花城", time: "2026-01-05 09:58", method: "智能开锁", status: "失败" }
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
      title: "开锁记录"
    }),
    b: common_vendor.f($data.items, (it, idx, i0) => {
      return {
        a: common_vendor.t(it.title),
        b: common_vendor.t(it.status),
        c: common_vendor.n(it.status === "成功" ? "ok" : "warn"),
        d: common_vendor.t(it.time),
        e: common_vendor.t(it.method),
        f: idx
      };
    }),
    c: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/records/record_unlock.js.map
