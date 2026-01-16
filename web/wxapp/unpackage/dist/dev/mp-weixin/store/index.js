"use strict";
const common_vendor = require("../common/vendor.js");
const store_user = require("./user.js");
const store_wxInfo = require("./wxInfo.js");
const pinia = common_vendor.createPinia();
pinia.use(common_vendor.index);
const stores = {
  user: store_user.user,
  wxInfo: store_wxInfo.wxInfo
};
const setupPinia = (app) => {
  app.use(pinia);
};
const $store = (name) => {
  const useStore = stores[name];
  if (!useStore) {
    throw new Error(`store ${name} not found`);
  }
  return useStore(pinia);
};
exports.$store = $store;
exports.setupPinia = setupPinia;
//# sourceMappingURL=../../.sourcemap/mp-weixin/store/index.js.map
