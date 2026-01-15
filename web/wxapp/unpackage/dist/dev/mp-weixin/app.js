"use strict";
Object.defineProperty(exports, Symbol.toStringTag, { value: "Module" });
const common_vendor = require("./common/vendor.js");
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
const _sfc_main = {
  onLaunch: function() {
    common_vendor.index$1.__f__("log", "at App.vue:4", "App Launch");
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
  },
  onShow: function() {
    common_vendor.index$1.__f__("log", "at App.vue:25", "App Show");
  },
  onHide: function() {
    common_vendor.index$1.__f__("log", "at App.vue:28", "App Hide");
  }
};
const pinia = common_vendor.createPinia();
function createApp() {
  const app = common_vendor.createSSRApp(_sfc_main);
  app.use(pinia);
  return {
    app
  };
}
createApp().app.mount("#app");
exports.createApp = createApp;
//# sourceMappingURL=../.sourcemap/mp-weixin/app.js.map
