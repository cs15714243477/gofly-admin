"use strict";
const common_vendor = require("../../common/vendor.js");
const utils_config_index = require("../config/index.js");
const store_index = require("../../store/index.js");
function base64Encode(str = "") {
  const utf8ToBytes = (s) => {
    const bytes = [];
    for (let i = 0; i < s.length; i++) {
      let codePoint = s.charCodeAt(i);
      if (codePoint < 128) {
        bytes.push(codePoint);
      } else if (codePoint < 2048) {
        bytes.push(192 | codePoint >> 6);
        bytes.push(128 | codePoint & 63);
      } else {
        bytes.push(224 | codePoint >> 12);
        bytes.push(128 | codePoint >> 6 & 63);
        bytes.push(128 | codePoint & 63);
      }
    }
    return new Uint8Array(bytes);
  };
  const buf = utf8ToBytes(str).buffer;
  if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.arrayBufferToBase64)
    return common_vendor.wx$1.arrayBufferToBase64(buf);
  if (typeof common_vendor.index$1 !== "undefined" && common_vendor.index$1.arrayBufferToBase64)
    return common_vendor.index$1.arrayBufferToBase64(buf);
  return "";
}
const options = {
  // 显示操作成功消息 默认不显示
  showSuccess: false,
  // 成功提醒 默认使用后端返回值
  successMsg: "",
  // 显示失败消息 默认显示
  showError: true,
  // 失败提醒 默认使用后端返回信息
  errorMsg: "",
  // 显示请求时loading模态框 默认显示
  showLoading: true,
  // loading提醒文字
  loadingMsg: "加载中",
  // 需要授权才能请求 默认放开
  auth: false
  // ...
};
let LoadingInstance = {
  target: null,
  count: 0
};
function closeLoading() {
  if (LoadingInstance.count > 0)
    LoadingInstance.count--;
  if (LoadingInstance.count === 0)
    common_vendor.index$1.hideLoading();
}
const http = new common_vendor.Request({
  baseURL: utils_config_index.baseUrl + utils_config_index.apiModel,
  timeout: 8e3,
  method: "GET",
  custom: options
});
http.interceptors.request.use(
  (config) => {
    const token = common_vendor.index$1.getStorageSync("token");
    if (token && !store_index.$store("user").isLogin) {
      store_index.$store("user").setToken(token);
    }
    if (config.custom.auth && !token && !store_index.$store("user").isLogin) {
      return Promise.reject();
    }
    if (config.custom.showLoading) {
      LoadingInstance.count++;
      LoadingInstance.count === 0 && common_vendor.index$1.showLoading({
        title: config.custom.loadingMsg,
        mask: true,
        fail: () => {
          common_vendor.index$1.hideLoading();
        }
      });
    }
    let timestamp = Date.parse((/* @__PURE__ */ new Date()).toString()) / 1e3;
    const passstr = common_vendor.md5("gofly@888" + timestamp);
    let baseheader = {
      Accept: "text/json",
      "Content-Type": "application/json;charset=UTF-8",
      "Businessid": "1",
      "apiverify": base64Encode(passstr + "#" + timestamp)
    };
    if (token) {
      baseheader = Object.assign({}, baseheader, {
        Authorization: `${token}`
        // 这里是token(可自行修改)
      });
    }
    config.header = baseheader;
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
http.interceptors.response.use(
  (response) => {
    if (response.header.authorization || response.header.Authorization) {
      store_index.$store("user").setToken(response.header.authorization || response.header.Authorization);
    }
    response.config.custom.showLoading && closeLoading();
    if (response.data.code !== 0) {
      if (response.config.custom.showError)
        common_vendor.index$1.showToast({
          title: response.data.message || "服务器开小差啦,请稍后再试~",
          icon: "none",
          mask: true
        });
      return Promise.resolve(response.data);
    }
    if (response.data.code !== 0 && response.data.message !== "" && response.config.custom.showSuccess) {
      common_vendor.index$1.showToast({
        title: response.config.custom.successMsg || response.data.message,
        icon: "none"
      });
    }
    if (response.data && response.data.token) {
      const userStore = store_index.$store("user");
      userStore.setToken(response.data.token);
    }
    return Promise.resolve(response.data);
  },
  (error) => {
    var _a;
    const userStore = store_index.$store("user");
    const isLogin = userStore.isLogin;
    let errorMessage = "网络请求出错";
    if (error !== void 0) {
      switch (error.statusCode) {
        case 400:
          errorMessage = "请求错误";
          break;
        case 401:
          if (isLogin) {
            errorMessage = "您的登陆已过期";
          } else {
            errorMessage = "请先登录";
          }
          userStore.logout(true);
          break;
        case 403:
          errorMessage = "拒绝访问";
          break;
        case 404:
          errorMessage = "请求出错";
          break;
        case 408:
          errorMessage = "请求超时";
          break;
        case 429:
          errorMessage = "请求频繁, 请稍后再访问";
          break;
        case 500:
          errorMessage = "服务器开小差啦,请稍后再试~";
          break;
        case 501:
          errorMessage = "服务未实现";
          break;
        case 502:
          errorMessage = "网络错误";
          break;
        case 503:
          errorMessage = "服务不可用";
          break;
        case 504:
          errorMessage = "网络超时";
          break;
        case 505:
          errorMessage = "HTTP版本不受支持";
          break;
      }
    }
    if (error && error.config) {
      if (error.config.custom.showError === false) {
        common_vendor.index$1.showToast({
          title: ((_a = error.data) == null ? void 0 : _a.message) || errorMessage,
          icon: "none",
          mask: true
        });
      }
      error.config.custom.showLoading && closeLoading();
    }
    return false;
  }
);
const request = (config) => {
  if (config && typeof config.url === "string") {
    const url = config.url;
    const model = utils_config_index.apiModel;
    const pathPrefix = utils_config_index.apiPath;
    if (model.startsWith("/") && url.startsWith(model + "/")) {
      config.url = url.slice(model.length);
    }
    if (pathPrefix.startsWith("/") && config.url.startsWith(pathPrefix + "/")) {
      config.url = config.url.slice(pathPrefix.length);
    }
    if (config.url.startsWith("/business/business/")) {
      config.url = config.url.replace("/business/business/", "/business/");
    }
    if (config.url.startsWith("/uniapp/uniapp/")) {
      config.url = config.url.replace("/uniapp/uniapp/", "/uniapp/");
    }
  }
  const needBusinessPrefix = !utils_config_index.apiModel.includes("/business") && !utils_config_index.apiPath.startsWith("/business");
  if (needBusinessPrefix && typeof config.url === "string" && config.url.startsWith("/") && !config.url.startsWith("/business/")) {
    config.url = "/business" + config.url;
  }
  if (config.url[0] !== "/") {
    config.url = utils_config_index.apiPath + config.url;
  }
  return http.middleware(config);
};
exports.request = request;
//# sourceMappingURL=../../../.sourcemap/mp-weixin/utils/request/index.js.map
