"use strict";
const common_vendor = require("../../common/vendor.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      statusBarHeight: 0,
      headerTop: 0
    };
  },
  onLoad() {
    const info = common_vendor.index$1.getSystemInfoSync();
    this.statusBarHeight = info.statusBarHeight;
    try {
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.getMenuButtonBoundingClientRect) {
        const rect = common_vendor.wx$1.getMenuButtonBoundingClientRect();
        this.headerTop = Number(rect.top || this.statusBarHeight) + 6;
      } else {
        this.headerTop = this.statusBarHeight + 12;
      }
    } catch (e) {
      this.headerTop = this.statusBarHeight + 12;
    }
  },
  methods: {
    goBack() {
      common_vendor.index$1.navigateBack();
    },
    copyPassword() {
      common_vendor.index$1.setClipboardData({
        data: "259308",
        success: () => {
          common_vendor.index$1.showToast({
            title: "复制成功",
            icon: "success"
          });
        }
      });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o((...args) => $options.goBack && $options.goBack(...args)),
    b: $data.headerTop + "px",
    c: common_assets._imports_0$1,
    d: common_vendor.o((...args) => $options.copyPassword && $options.copyPassword(...args))
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/unlock_details/unlock_details.js.map
