"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { TopHeader },
  onLoad(options) {
    this.mode = options && options.mode || "complete";
    const cachedPhone = common_vendor.index$1.getStorageSync("hm_phone");
    if (cachedPhone)
      this.form.phone = cachedPhone;
  },
  data() {
    return {
      agreed: true,
      mode: "complete",
      // complete | register（预留）
      region: ["北京市", "北京市", "朝阳区"],
      form: {
        name: "",
        phone: "",
        store: "链家地产 (朝阳大悦城店)"
      }
    };
  },
  computed: {
    pageTitle() {
      return this.mode === "complete" ? "完善信息" : "注册";
    },
    headerTitle() {
      return this.mode === "complete" ? "完善资料" : "欢迎加入经纪人";
    },
    headerSubtitle() {
      return this.mode === "complete" ? "首次登录请先补全信息，完成后进入系统" : "请填写以下信息以完成注册";
    },
    submitText() {
      return this.mode === "complete" ? "完成并进入" : "新用户注册";
    },
    regionText() {
      const r = Array.isArray(this.region) ? this.region : [];
      return r.filter(Boolean).join(" ");
    }
  },
  methods: {
    goBack() {
      common_vendor.index$1.navigateBack();
    },
    agreementChange(e) {
      this.agreed = e.detail.value.length > 0;
    },
    onRegionChange(e) {
      const value = e && e.detail && e.detail.value;
      if (Array.isArray(value) && value.length)
        this.region = value;
    },
    selectStore() {
      common_vendor.index$1.__f__("log", "at pages/registration/registration.vue:149", "选择门店");
    },
    handleSubmit() {
      if (!this.agreed) {
        common_vendor.index$1.showToast({
          title: "请先同意协议",
          icon: "none"
        });
        return;
      }
      if (!this.form.name) {
        common_vendor.index$1.showToast({ title: "请填写姓名", icon: "none" });
        return;
      }
      if (!this.form.phone) {
        common_vendor.index$1.showToast({ title: "请先在登录页授权手机号", icon: "none" });
        return;
      }
      common_vendor.index$1.setStorageSync("hm_profile_completed", true);
      common_vendor.index$1.reLaunch({
        url: "/pages/home/home"
      });
    }
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  _component_TopHeader();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.o((...args) => $options.goBack && $options.goBack(...args)),
    b: common_vendor.p({
      title: $options.pageTitle
    }),
    c: common_vendor.t($options.headerTitle),
    d: common_vendor.t($options.headerSubtitle),
    e: common_vendor.t($options.regionText),
    f: $data.region,
    g: common_vendor.o((...args) => $options.onRegionChange && $options.onRegionChange(...args)),
    h: $data.form.name,
    i: common_vendor.o(($event) => $data.form.name = $event.detail.value),
    j: $data.form.phone,
    k: common_vendor.o(($event) => $data.form.phone = $event.detail.value),
    l: $data.form.store,
    m: common_vendor.o((...args) => $options.selectStore && $options.selectStore(...args)),
    n: common_vendor.o((...args) => $options.agreementChange && $options.agreementChange(...args)),
    o: common_vendor.t($options.submitText),
    p: common_vendor.o((...args) => $options.handleSubmit && $options.handleSubmit(...args))
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/registration/registration.js.map
