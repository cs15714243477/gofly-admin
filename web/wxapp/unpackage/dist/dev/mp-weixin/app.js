"use strict";
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
const common_vendor = require("./common/vendor.js");
const store_index = require("./store/index.js");
if (!Math) {
  "./pages/login/login.js";
  "./pages/registration/registration.js";
  "./pages/home/home.js";
  "./pages/property_list/property_list.js";
  "./pages/property_detail/property_detail.js";
  "./pages/agent_workbench_home/agent_workbench_home.js";
  "./pages/unlock_details/unlock_details.js";
  "./pages/unlock_steps/unlock_steps.js";
  "./pages/my_business_card/my_business_card.js";
  "./pages/marketing_posters/marketing_posters.js";
  "./pages/records/record_follow.js";
  "./pages/records/record_unlock.js";
  "./pages/records/record_showing.js";
  "./pages/records/record_view.js";
  "./pages/records/record_share.js";
  "./pages/records/record_call.js";
}
(function() {
  try {
    const g = typeof globalThis !== "undefined" && globalThis || typeof global !== "undefined" && global || typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1 || {};
    if (g && typeof g.__VUE_DEVTOOLS_ON_SOCKET_READY__ !== "function") {
      g.__VUE_DEVTOOLS_ON_SOCKET_READY__ = function() {
      };
    }
    if (typeof common_vendor.wx$1 !== "undefined" && typeof common_vendor.wx$1.__VUE_DEVTOOLS_ON_SOCKET_READY__ !== "function") {
      common_vendor.wx$1.__VUE_DEVTOOLS_ON_SOCKET_READY__ = function() {
      };
    }
  } catch (e) {
  }
})();
const _sfc_main = {
  data() {
    return {
      __bootstrapped: false
    };
  },
  methods: {
    async bootstrapAuthRedirect() {
      if (this.__bootstrapped)
        return;
      this.__bootstrapped = true;
      const token = common_vendor.index$1.getStorageSync("token");
      if (!token)
        return;
      const userStore = store_index.$store("user");
      if (!userStore.isLogin)
        userStore.setToken(token);
      let ok = false;
      try {
        const info = await userStore.getInfo();
        ok = !!info;
      } catch (e) {
        ok = false;
      }
      if (!ok) {
        try {
          await userStore.logout(true);
        } catch (e) {
        }
        return;
      }
      const pages = getCurrentPages ? getCurrentPages() : [];
      const cur = pages && pages.length ? pages[pages.length - 1] : null;
      const route = cur && cur.route ? cur.route : "";
      if (route === "pages/login/login") {
        common_vendor.index$1.reLaunch({ url: "/pages/property_list/property_list" });
      }
    }
  },
  onLaunch: function() {
    common_vendor.index$1.__f__("log", "at App.vue:71", "App Launch");
    try {
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.loadFontFace) {
        common_vendor.wx$1.loadFontFace({
          family: "Material Symbols Outlined",
          // 使用 jsDelivr CDN（相对 raw.githubusercontent.com 更稳定）//TODO后续放入自己的服务器中
          source: 'url("https://cdn.jsdelivr.net/gh/google/material-design-icons@master/font/MaterialIconsOutlined-Regular.otf")',
          global: true
        });
      }
    } catch (e) {
    }
    this.bootstrapAuthRedirect();
  },
  onShow: function() {
    common_vendor.index$1.__f__("log", "at App.vue:95", "App Show");
    this.bootstrapAuthRedirect();
  },
  onHide: function() {
    common_vendor.index$1.__f__("log", "at App.vue:100", "App Hide");
  }
};
function createApp() {
  const app = common_vendor.createSSRApp(_sfc_main);
  store_index.setupPinia(app);
  return {
    app
  };
}
createApp().app.mount("#app");
exports.createApp = createApp;
//# sourceMappingURL=../.sourcemap/mp-weixin/app.js.map
