"use strict";
const common_vendor = require("../../common/vendor.js");
const store_index = require("../../store/index.js");
const BottomTabBar = () => "../../components/BottomTabBar.js";
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { BottomTabBar, TopHeader },
  data() {
    return {
      loadingUser: false,
      debugLogged: false,
      userInfo: {},
      businessRecords: [
        { key: "follow", name: "关注记录", icon: "favorite" },
        { key: "unlock", name: "开锁记录", icon: "lock_open", hasNotice: true },
        { key: "showing", name: "带看记录", icon: "location_on" },
        { key: "view", name: "浏览记录", icon: "history" },
        { key: "share", name: "分享记录", icon: "share" },
        { key: "call", name: "通话记录", icon: "call" }
      ]
    };
  },
  computed: {
    avatarUrl() {
      return this.userInfo && this.userInfo.avatar ? this.userInfo.avatar : "/static/images/img_83575f387f.png";
    },
    displayName() {
      const u = this.userInfo || {};
      return u.name || u.nickname || u.username || "未登录";
    },
    displayRoleLine() {
      const u = this.userInfo || {};
      const title = (u.title || "").trim();
      const roleRaw = String(u.role || "").trim();
      const storeName = (u.store_name || "").trim();
      const store = storeName ? storeName : "未绑定";
      let left = title;
      if (!left) {
        if (roleRaw === "" || roleRaw === "1" || roleRaw.toLowerCase() === "user")
          left = "经纪人";
        else
          left = roleRaw;
      }
      return `${left} | ${store}`;
    },
    displayMobile() {
      const u = this.userInfo || {};
      return u.mobile || "";
    }
  },
  onShow() {
    common_vendor.index$1.__f__("log", "at pages/agent_workbench_home/agent_workbench_home.vue:133", "App 1");
    this.ensureLoginAndLoadUser();
  },
  methods: {
    debugPrintUserInfo(tag = "") {
      var _a, _b, _c, _d;
      if (this.debugLogged)
        return;
      this.debugLogged = true;
      try {
        common_vendor.index$1.__f__("log", "at pages/agent_workbench_home/agent_workbench_home.vue:142", "[agent_workbench_home] " + (tag || "userInfo") + " store_name=", (_a = this.userInfo) == null ? void 0 : _a.store_name, "store_id=", (_b = this.userInfo) == null ? void 0 : _b.store_id, "title=", (_c = this.userInfo) == null ? void 0 : _c.title, "role=", (_d = this.userInfo) == null ? void 0 : _d.role);
        common_vendor.index$1.__f__("log", "at pages/agent_workbench_home/agent_workbench_home.vue:143", "[agent_workbench_home] " + (tag || "userInfo") + " full=", JSON.parse(JSON.stringify(this.userInfo || {})));
      } catch (e) {
        common_vendor.index$1.__f__("log", "at pages/agent_workbench_home/agent_workbench_home.vue:145", "[agent_workbench_home] debugPrintUserInfo error", e);
      }
    },
    debugShowUserInfo() {
      try {
        const u = this.userInfo || {};
        const contentLines = [
          `name: ${u.name || ""}`,
          `nickname: ${u.nickname || ""}`,
          `username: ${u.username || ""}`,
          `title: ${u.title || ""}`,
          `role: ${u.role || ""}`,
          `store_name: ${u.store_name || ""}`,
          `store_id: ${u.store_id || ""}`,
          `store_address: ${u.store_address || ""}`
        ];
        common_vendor.index$1.showModal({
          title: "调试：userInfo(点击头像)",
          content: contentLines.join("\n"),
          showCancel: false
        });
      } catch (e) {
        common_vendor.index$1.showToast({ title: "调试弹窗失败", icon: "none" });
      }
    },
    async ensureLoginAndLoadUser() {
      common_vendor.index$1.__f__("log", "at pages/agent_workbench_home/agent_workbench_home.vue:172", "App 2");
      const userStore = store_index.$store("user");
      const token = common_vendor.index$1.getStorageSync("token");
      if (token && !userStore.isLogin) {
        userStore.setToken(token);
      }
      if (!token && !userStore.isLogin) {
        common_vendor.index$1.reLaunch({ url: "/pages/login/login" });
        return;
      }
      if (this.loadingUser)
        return;
      this.loadingUser = true;
      try {
        const info = await userStore.getInfo();
        this.userInfo = info || userStore.userInfo || {};
        this.debugPrintUserInfo("after getInfo");
      } catch (e) {
        this.userInfo = userStore.userInfo || {};
        this.debugPrintUserInfo("fallback userStore.userInfo");
      } finally {
        this.loadingUser = false;
      }
    },
    goEditCard() {
      common_vendor.index$1.reLaunch({
        url: "/pages/my_business_card/my_business_card?tab=1"
      });
    },
    openRecord(item) {
      const map = {
        follow: "/pages/records/record_follow",
        unlock: "/pages/records/record_unlock",
        showing: "/pages/records/record_showing",
        view: "/pages/records/record_view",
        share: "/pages/records/record_share",
        call: "/pages/records/record_call"
      };
      const url = item && item.key && map[item.key];
      if (!url)
        return;
      common_vendor.index$1.navigateTo({ url });
    },
    handleLogout() {
      common_vendor.index$1.showModal({
        title: "提示",
        content: "确定要退出登录吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              await store_index.$store("user").logout(false);
            } catch (e) {
            }
            common_vendor.index$1.reLaunch({
              url: "/pages/login/login"
            });
          }
        }
      });
    }
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  const _component_BottomTabBar = common_vendor.resolveComponent("BottomTabBar");
  (_component_TopHeader + _component_BottomTabBar)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.p({
      title: "我的"
    }),
    b: common_vendor.o((...args) => $options.goEditCard && $options.goEditCard(...args)),
    c: $options.avatarUrl,
    d: common_vendor.t($options.displayName),
    e: common_vendor.t($options.displayRoleLine),
    f: common_vendor.t($options.displayMobile),
    g: common_vendor.f($data.businessRecords, (item, index, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.icon),
        b: item.hasNotice
      }, item.hasNotice ? {} : {}, {
        c: common_vendor.t(item.name),
        d: index,
        e: common_vendor.o(($event) => $options.openRecord(item), index)
      });
    }),
    h: common_vendor.o((...args) => $options.handleLogout && $options.handleLogout(...args)),
    i: common_vendor.p({
      active: "me"
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/agent_workbench_home/agent_workbench_home.js.map
