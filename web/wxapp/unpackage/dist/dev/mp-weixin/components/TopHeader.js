"use strict";
const common_vendor = require("../common/vendor.js");
const _sfc_main = {
  name: "TopHeader",
  props: {
    title: {
      type: String,
      default: ""
    }
  },
  data() {
    return {
      padPx: 0,
      rightSafePx: 0,
      statusBarPx: 0,
      navBarPx: 0
    };
  },
  computed: {
    wrapStyle() {
      if (!this.statusBarPx || !this.navBarPx)
        return {};
      return {
        paddingTop: this.statusBarPx + "px",
        height: this.statusBarPx + this.navBarPx + "px"
      };
    },
    innerStyle() {
      const style = {};
      if (this.padPx) {
        style.paddingLeft = this.padPx + "px";
        style.paddingRight = this.padPx + this.rightSafePx + "px";
      }
      if (this.navBarPx)
        style.height = this.navBarPx + "px";
      return style;
    }
  },
  mounted() {
    try {
      const info = common_vendor.index$1.getSystemInfoSync();
      const w = Number(info.windowWidth || 375);
      const statusBar = Number(info.statusBarHeight || 0);
      const pad = Math.round(w * 32 / 750);
      let rightSafe = 0;
      let navBarH = 0;
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.getMenuButtonBoundingClientRect) {
        const rect = common_vendor.wx$1.getMenuButtonBoundingClientRect();
        rightSafe = Math.max(0, w - Number(rect.left || 0));
        navBarH = Math.max(44, (Number(rect.top || 0) - statusBar) * 2 + Number(rect.height || 0)) || 44;
      }
      this.padPx = pad;
      this.rightSafePx = rightSafe;
      this.statusBarPx = statusBar;
      this.navBarPx = navBarH;
    } catch (e) {
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: _ctx.$slots.left
  }, _ctx.$slots.left ? {} : {}, {
    b: _ctx.$slots.default
  }, _ctx.$slots.default ? {} : {
    c: common_vendor.t($props.title)
  }, {
    d: common_vendor.s($options.innerStyle),
    e: common_vendor.s($options.wrapStyle)
  });
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-65748c24"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/TopHeader.js.map
