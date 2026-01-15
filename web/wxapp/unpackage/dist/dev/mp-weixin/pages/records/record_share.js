"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const BottomTabBar = () => "../../components/BottomTabBar.js";
const _sfc_main = {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      items: [
        { title: "阳光花园三期 2室1厅", time: "2026-01-07 10:10", channel: "微信", clicks: 12 },
        { title: "滨江别墅 独栋花园", time: "2026-01-06 15:30", channel: "朋友圈", clicks: 8 },
        { title: "金城中心 单身公寓", time: "2026-01-05 09:20", channel: "复制链接", clicks: 3 }
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
      title: "分享记录"
    }),
    b: common_vendor.f($data.items, (it, idx, i0) => {
      return {
        a: common_vendor.t(it.title),
        b: common_vendor.t(it.channel),
        c: common_vendor.t(it.time),
        d: common_vendor.t(it.clicks),
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
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/records/record_share.js.map
