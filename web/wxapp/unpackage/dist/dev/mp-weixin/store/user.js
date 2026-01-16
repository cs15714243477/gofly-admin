"use strict";
const common_vendor = require("../common/vendor.js");
const api_user = require("../api/user.js");
const defaultUserInfo = {
  avatar: "",
  // 头像
  name: "",
  // 姓名
  nickname: "",
  // 昵称
  mobile: ""
  // 手机号
};
const user = common_vendor.defineStore({
  id: "user",
  state: () => ({
    userInfo: { ...defaultUserInfo },
    // 用户信息
    isLogin: !!common_vendor.index$1.getStorageSync("token")
    // 登录状态
  }),
  actions: {
    // 获取个人信息
    async getInfo() {
      const res = await api_user.userApi.getuserinfo();
      if (!res || res.code !== 0)
        return;
      this.userInfo = res.data;
      return Promise.resolve(res.data);
    },
    //更新用户信息
    setUserInfo(info) {
      this.userInfo = info;
    },
    // 设置token
    setToken(token = "") {
      if (token === "") {
        this.isLogin = false;
        common_vendor.index$1.removeStorageSync("token");
      } else {
        this.isLogin = true;
        common_vendor.index$1.setStorageSync("token", token);
      }
      return this.isLogin;
    },
    // 重置用户默认数据
    resetUserData() {
      this.setToken();
      this.userInfo = common_vendor.cloneDeep(defaultUserInfo);
    },
    // 登出
    async logout(force = false) {
      if (!force) {
        const res = await api_user.userApi.logout();
        if (res && res.code === 0) {
          this.resetUserData();
        }
      }
      if (force) {
        this.resetUserData();
      }
      return !this.isLogin;
    }
  }
});
exports.user = user;
//# sourceMappingURL=../../.sourcemap/mp-weixin/store/user.js.map
