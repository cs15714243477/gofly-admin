"use strict";
const common_vendor = require("../../common/vendor.js");
const common_assets = require("../../common/assets.js");
const BottomTabBar = () => "../../components/BottomTabBar.js";
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { BottomTabBar, TopHeader },
  data() {
    return {
      businessRecords: [
        { key: "follow", name: "关注记录", icon: "favorite" },
        { key: "unlock", name: "开锁记录", icon: "lock_open", hasNotice: true },
        { key: "showing", name: "带看记录", icon: "location_on" },
        { key: "view", name: "浏览记录", icon: "history" },
        { key: "share", name: "分享记录", icon: "share" },
        { key: "call", name: "通话记录", icon: "call" }
      ]
    };
  },
  methods: {
    goEditCard() {
      common_vendor.index$1.reLaunch({
        url: "/pages/my_business_card/my_business_card?tab=1"
      });
    },
    openRecord(item) {
      const map = {
        follow: "/pages/records/record_follow",
        unlock: "/pages/records/record_unlock",
        showing: "/pages/records/record_showing",
        view: "/pages/records/record_view",
        share: "/pages/records/record_share",
        call: "/pages/records/record_call"
      };
      const url = item && item.key && map[item.key];
      if (!url)
        return;
      common_vendor.index$1.navigateTo({ url });
    },
    handleLogout() {
      common_vendor.index$1.showModal({
        title: "提示",
        content: "确定要退出登录吗？",
        success: (res) => {
          if (res.confirm) {
            common_vendor.index$1.reLaunch({
              url: "/pages/login/login"
            });
          }
        }
      });
    }
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
      title: "我的"
    }),
    b: common_vendor.o((...args) => $options.goEditCard && $options.goEditCard(...args)),
    c: common_assets._imports_0$1,
    d: common_vendor.f($data.businessRecords, (item, index, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.icon),
        b: item.hasNotice
      }, item.hasNotice ? {} : {}, {
        c: common_vendor.t(item.name),
        d: index,
        e: common_vendor.o(($event) => $options.openRecord(item), index)
      });
    }),
    e: common_vendor.o((...args) => $options.handleLogout && $options.handleLogout(...args)),
    f: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/agent_workbench_home/agent_workbench_home.js.map
