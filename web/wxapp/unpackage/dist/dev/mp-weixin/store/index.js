"use strict";const e=require("../common/vendor.js"),r=require("./user.js"),s=require("./wxInfo.js"),o=e.createPinia();o.use(e.index);const n={user:r.user,wxInfo:s.wxInfo};exports.$store=e=>{const r=n[e];if(!r)throw new Error(`store ${e} not found`);return r(o)},exports.setupPinia=e=>{e.use(o)};
//# sourceMappingURL=../../.sourcemap/mp-weixin/store/index.js.map
