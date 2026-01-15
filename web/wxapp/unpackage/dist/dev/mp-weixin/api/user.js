"use strict";
const utils_request_index = require("../utils/request/index.js");
const userApi = {
  // 账号登录
  login: (data) => utils_request_index.request({
    url: "/user/login",
    method: "POST",
    data,
    custom: {
      showSuccess: true,
      loadingMsg: "登录中"
    }
  }),
  // 微信一键登录（手机号）
  wxLogin: (data) => utils_request_index.request({
    url: "/user/wxLogin",
    method: "POST",
    data,
    custom: {
      showSuccess: true,
      loadingMsg: "登录中"
    }
  }),
  // 获取用户信息
  getuserinfo: () => utils_request_index.request({
    url: "/user/getUserinfo",
    method: "GET",
    custom: {
      showLoading: true,
      auth: true
    }
  }),
  // 账号登出
  logout: (data) => utils_request_index.request({
    url: "/user/logout",
    method: "POST",
    data
  })
};
exports.userApi = userApi;
//# sourceMappingURL=../../.sourcemap/mp-weixin/api/user.js.map
