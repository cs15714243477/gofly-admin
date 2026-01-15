"use strict";
const common_vendor = require("../../common/vendor.js");
const BottomTabBar = () => "../../components/BottomTabBar.js";
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { BottomTabBar, TopHeader },
  onLoad(options) {
    const tab = options && (options.tab || options.currentTab);
    if (tab === "1" || tab === 1 || tab === "edit") {
      this.currentTab = 1;
    }
  },
  data() {
    return {
      // 默认先展示名片预览
      currentTab: 0,
      currentPreview: 0,
      // “我的资料”数据源：名片预览与编辑资料保持一致
      profile: {
        avatar: "/static/images/img_155372af18.png",
        name: "张伟",
        phone: "138 1234 5678",
        wechat: "zhangwei_agent",
        role: "资深置业顾问",
        company: "链家地产",
        store: "朝阳区北京新天地店",
        email: "zhangwei@lianjia.com",
        addr: "北京市朝阳区...",
        intro: "从事房地产行业5年，专注于朝阳区高端住宅买卖。我致力于为每一位客户提供专业、诚信的置业咨询服务，帮助您找到理想的家。",
        tags: ["二手房买卖", "学区房专家"]
      },
      suggestTags: ["高端别墅"],
      // 名片样式（只改变视觉，不改变字段）
      cardStyles: [
        { bg: "linear-gradient(135deg, #2d9cf0, #1e40af)", decor2: "rgba(249, 115, 22, 0.2)", badge: "PRO", badgeBg: "#f97316" },
        { bg: "linear-gradient(135deg, #06b6d4, #2563eb)", decor2: "rgba(255, 255, 255, 0.12)", badge: "TOP", badgeBg: "#22c55e" },
        { bg: "linear-gradient(135deg, #111827, #334155)", decor2: "rgba(45, 156, 240, 0.18)", badge: "VIP", badgeBg: "#2d9cf0" }
      ]
    };
  },
  computed: {
    previewCards() {
      const p = this.profile || {};
      const tags = Array.isArray(p.tags) ? p.tags.slice(0, 3) : [];
      return (this.cardStyles || []).map((s) => ({
        ...s,
        company: p.company,
        name: p.name,
        role: p.role,
        phone: p.phone,
        wechat: p.wechat,
        store: p.store,
        avatar: p.avatar,
        tags
      }));
    }
  },
  methods: {
    onPreviewChange(e) {
      this.currentPreview = e.detail.current;
    }
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  const _component_BottomTabBar = common_vendor.resolveComponent("BottomTabBar");
  (_component_TopHeader + _component_BottomTabBar)();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      title: "获客"
    }),
    b: $data.currentTab === 0 ? 1 : "",
    c: common_vendor.o(($event) => $data.currentTab = 0),
    d: $data.currentTab === 1 ? 1 : "",
    e: common_vendor.o(($event) => $data.currentTab = 1),
    f: $data.currentTab === 0
  }, $data.currentTab === 0 ? {
    g: common_vendor.f($options.previewCards, (_, i, i0) => {
      return {
        a: i,
        b: $data.currentPreview === i ? 1 : ""
      };
    }),
    h: common_vendor.f($options.previewCards, (c, idx, i0) => {
      return {
        a: c.decor2,
        b: common_vendor.t(c.company),
        c: common_vendor.t(c.badge),
        d: c.badgeBg,
        e: c.avatar,
        f: common_vendor.t(c.name),
        g: common_vendor.t(c.role),
        h: common_vendor.f(c.tags, (t, tIdx, i1) => {
          return {
            a: common_vendor.t(t),
            b: tIdx
          };
        }),
        i: common_vendor.t(c.phone),
        j: common_vendor.t(c.wechat),
        k: common_vendor.t(c.store),
        l: c.bg,
        m: idx
      };
    }),
    i: common_vendor.o((...args) => $options.onPreviewChange && $options.onPreviewChange(...args))
  } : {}, {
    j: $data.currentTab === 1
  }, $data.currentTab === 1 ? {
    k: $data.profile.avatar,
    l: common_vendor.t($data.profile.name),
    m: common_vendor.t($data.profile.phone),
    n: common_vendor.t($data.profile.wechat),
    o: common_vendor.t($data.profile.role),
    p: common_vendor.t($data.profile.company),
    q: common_vendor.t($data.profile.store),
    r: $data.profile.intro,
    s: common_vendor.o(($event) => $data.profile.intro = $event.detail.value),
    t: common_vendor.f($data.profile.tags, (t, idx, i0) => {
      return {
        a: common_vendor.t(t),
        b: "active-" + idx
      };
    }),
    v: common_vendor.f($data.suggestTags, (t, idx, i0) => {
      return {
        a: common_vendor.t(t),
        b: "suggest-" + idx
      };
    })
  } : {}, {
    w: $data.currentTab === 0
  }, $data.currentTab === 0 ? {} : {}, {
    x: $data.currentTab === 1
  }, $data.currentTab === 1 ? {} : {}, {
    y: common_vendor.p({
      active: "marketing"
    })
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/my_business_card/my_business_card.js.map
