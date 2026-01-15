"use strict";
const common_vendor = require("../common/vendor.js");
const _sfc_main = {
  name: "BottomTabBar",
  props: {
    // home | property | marketing | me
    active: {
      type: String,
      default: "home"
    }
  },
  data() {
    return {
      items: [
        { key: "property", label: "房源", icon: "apartment", url: "/pages/property_list/property_list" },
        { key: "home", label: "推荐", icon: "home", url: "/pages/home/home" },
        { key: "marketing", label: "获客", icon: "campaign", url: "/pages/my_business_card/my_business_card" },
        { key: "me", label: "我的", icon: "person", url: "/pages/agent_workbench_home/agent_workbench_home" }
      ]
    };
  },
  methods: {
    onTap(item) {
      if (!item || !item.url)
        return;
      common_vendor.index$1.reLaunch({ url: item.url });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.f($data.items, (item, k0, i0) => {
      return {
        a: common_vendor.t(item.icon),
        b: $props.active === item.key ? 1 : "",
        c: $props.active === item.key ? 1 : "",
        d: common_vendor.t(item.label),
        e: item.key,
        f: $props.active === item.key ? 1 : "",
        g: common_vendor.o(($event) => $options.onTap(item), item.key)
      };
    })
  };
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render], ["__scopeId", "data-v-032475b5"]]);
wx.createComponent(Component);
//# sourceMappingURL=../../.sourcemap/mp-weixin/components/BottomTabBar.js.map
