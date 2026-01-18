"use strict";const e=require("../utils/request/index.js"),s={getProvinces:s=>e.request({url:"/uniapp/area/getProvinces",method:"GET",data:s,custom:{showSuccess:!1,showLoading:!1,auth:!0}}),getChildren:s=>e.request({url:"/uniapp/area/getChildren",method:"GET",data:s,custom:{showSuccess:!1,showLoading:!1,auth:!0}})};exports.areaApi=s;
//# sourceMappingURL=../../.sourcemap/mp-weixin/api/area.js.map
