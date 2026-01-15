"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { TopHeader },
  data() {
    return {};
  },
  methods: {
    goBack() {
      common_vendor.index$1.navigateBack();
    },
    bluetoothUnlock() {
      common_vendor.index$1.__f__("log", "at pages/unlock_steps/unlock_steps.vue:106", "蓝牙开锁");
    },
    passwordUnlock() {
      common_vendor.index$1.navigateTo({
        url: "/pages/unlock_details/unlock_details"
      });
    }
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  _component_TopHeader();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o((...args) => $options.goBack && $options.goBack(...args)),
    b: common_vendor.p({
      title: "开锁步骤"
    }),
    c: common_vendor.o((...args) => $options.bluetoothUnlock && $options.bluetoothUnlock(...args)),
    d: common_vendor.o((...args) => $options.passwordUnlock && $options.passwordUnlock(...args))
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/unlock_steps/unlock_steps.js.map
