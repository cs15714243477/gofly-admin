"use strict";
const common_vendor = require("../../common/vendor.js");
const store_index = require("../../store/index.js");
const api_user = require("../../api/user.js");
const _sfc_main = {
  data() {
    return {
      showMoreLogin: false,
      mobile: "",
      captcha: "",
      submitting: false
    };
  },
  methods: {
    goBack() {
      common_vendor.index$1.navigateBack();
    },
    goToRegister() {
      common_vendor.index$1.navigateTo({
        url: "/pages/registration/registration?mode=complete"
      });
    },
    async handleLogin() {
      if (!this.mobile || String(this.mobile).length !== 11) {
        common_vendor.index$1.showToast({ title: "请输入11位手机号", icon: "none" });
        return;
      }
      if (!this.captcha) {
        common_vendor.index$1.showToast({ title: "请输入验证码", icon: "none" });
        return;
      }
      if (this.submitting)
        return;
      this.submitting = true;
      try {
        const res = await api_user.userApi.login({ mobile: this.mobile, captcha: this.captcha });
        if (!res || res.code !== 0)
          return;
        if (res.token)
          store_index.$store("user").setToken(res.token);
        await store_index.$store("user").getInfo().catch(() => {
        });
        this.afterLoginSuccess();
      } finally {
        this.submitting = false;
      }
    },
    /**
     * 切换更多登录方式的显示状态
     *
     * 用于控制微信小程序中更多登录方式按钮的展开/收起状态
     * 通过反转 showMoreLogin 数据属性的值来实现状态切换
     *
     * @returns {void} 无返回值
     */
    toggleMoreLogin() {
      this.showMoreLogin = !this.showMoreLogin;
    },
    onGetPhoneNumber(e) {
      const phoneCode = e && e.detail && e.detail.code;
      if (!phoneCode) {
        common_vendor.index$1.showToast({ title: "未授权手机号", icon: "none" });
        return;
      }
      if (this.submitting)
        return;
      this.submitting = true;
      common_vendor.index$1.login({
        provider: "weixin",
        success: async (loginRes) => {
          try {
            const wxCode = loginRes && loginRes.code;
            const res = await api_user.userApi.wxLogin({ wx_code: wxCode, phone_code: phoneCode });
            if (!res || res.code !== 0)
              return;
            if (res.token)
              store_index.$store("user").setToken(res.token);
            await store_index.$store("user").getInfo().catch(() => {
            });
            this.afterLoginSuccess();
          } finally {
            this.submitting = false;
          }
        },
        fail: () => {
          this.submitting = false;
          common_vendor.index$1.showToast({ title: "获取微信登录凭证失败", icon: "none" });
        }
      });
    },
    afterLoginSuccess() {
      common_vendor.index$1.reLaunch({ url: "/pages/property_list/property_list" });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.o((...args) => $options.onGetPhoneNumber && $options.onGetPhoneNumber(...args)),
    b: common_vendor.t($data.showMoreLogin ? "收起其他登录方式" : "更多登录方式"),
    c: common_vendor.t($data.showMoreLogin ? "expand_less" : "expand_more"),
    d: common_vendor.o((...args) => $options.toggleMoreLogin && $options.toggleMoreLogin(...args)),
    e: $data.showMoreLogin
  }, $data.showMoreLogin ? {
    f: $data.mobile,
    g: common_vendor.o(($event) => $data.mobile = $event.detail.value),
    h: $data.captcha,
    i: common_vendor.o(($event) => $data.captcha = $event.detail.value),
    j: common_vendor.t($data.submitting ? "登录中..." : "登录"),
    k: $data.submitting,
    l: common_vendor.o((...args) => $options.handleLogin && $options.handleLogin(...args)),
    m: common_vendor.o((...args) => $options.goToRegister && $options.goToRegister(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/login/login.js.map
