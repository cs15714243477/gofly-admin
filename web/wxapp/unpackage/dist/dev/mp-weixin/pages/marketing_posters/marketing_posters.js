"use strict";
const common_vendor = require("../../common/vendor.js");
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { TopHeader },
  data() {
    return {
      currentCat: 0,
      categories: ["热门", "全部", "节日问候", "励志鸡汤", "早安晚安", "二手房源", "新房推荐", "豪宅鉴赏"],
      posters: [
        { title: "豪宅鉴赏·云端之上", desc: "高端大气上档次", tag: "HOT", tagType: "accent", aspect: "3/4", image: "/static/images/img_19380f982c.png" },
        { title: "冬日暖阳·早安", desc: "每日签到问候", tag: "NEW", tagType: "primary", aspect: "9/16", image: "/static/images/img_0685ece1a5.png" },
        { title: "极简生活方式", desc: "4.2w+ 使用", aspect: "1/1", image: "/static/images/img_08b712f810.png" },
        { title: "温馨小户型推荐", desc: "刚需首选", tag: "NEW", tagType: "primary", aspect: "3/4", image: "/static/images/img_b62c39564d.png" },
        { title: "奋斗正当时", desc: "励志鸡汤", aspect: "9/16", image: "/static/images/img_610af521bc.png" },
        { title: "端午安康·节日签", desc: "节日问候", tag: "HOT", tagType: "accent", aspect: "3/4", image: "/static/images/img_3d4235e8d0.png" },
        { title: "精装样板间", desc: "实拍展示", aspect: "1/1", image: "/static/images/img_8642388cea.png" }
      ]
    };
  },
  computed: {
    leftPosters() {
      return this.posters.filter((_, i) => i % 2 === 0);
    },
    rightPosters() {
      return this.posters.filter((_, i) => i % 2 !== 0);
    }
  }
};
if (!Array) {
  const _component_TopHeader = common_vendor.resolveComponent("TopHeader");
  _component_TopHeader();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return {
    a: common_vendor.p({
      title: "获客海报"
    }),
    b: common_vendor.f($data.categories, (cat, idx, i0) => {
      return common_vendor.e({
        a: $data.currentCat === idx
      }, $data.currentCat === idx ? {} : {}, {
        b: common_vendor.t(cat),
        c: idx,
        d: $data.currentCat === idx ? 1 : "",
        e: common_vendor.o(($event) => $data.currentCat = idx, idx)
      });
    }),
    c: common_vendor.f($options.leftPosters, (poster, pIdx, i0) => {
      return common_vendor.e({
        a: poster.image,
        b: poster.tag
      }, poster.tag ? {
        c: common_vendor.t(poster.tag),
        d: common_vendor.n(poster.tagType)
      } : {}, {
        e: poster.aspect,
        f: common_vendor.t(poster.title),
        g: common_vendor.t(poster.desc),
        h: pIdx
      });
    }),
    d: common_vendor.f($options.rightPosters, (poster, pIdx, i0) => {
      return common_vendor.e({
        a: poster.image,
        b: poster.tag
      }, poster.tag ? {
        c: common_vendor.t(poster.tag),
        d: common_vendor.n(poster.tagType)
      } : {}, {
        e: poster.aspect,
        f: common_vendor.t(poster.title),
        g: common_vendor.t(poster.desc),
        h: pIdx
      });
    })
  };
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/marketing_posters/marketing_posters.js.map
