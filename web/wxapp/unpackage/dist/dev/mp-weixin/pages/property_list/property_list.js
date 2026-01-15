"use strict";
const common_vendor = require("../../common/vendor.js");
const BottomTabBar = () => "../../components/BottomTabBar.js";
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { BottomTabBar, TopHeader },
  data() {
    return {
      currentCategory: "all",
      categories: [
        { key: "all", label: "全部房源" },
        { key: "mine", label: "我的房源" },
        { key: "nearby", label: "附近房源" },
        { key: "price_drop", label: "最新降价" },
        { key: "school", label: "学区房" }
      ],
      // 筛选（静态数据）
      filters: {
        rooms: "",
        smartLock: false,
        status: "",
        sort: ""
      },
      filterOpen: false,
      draft: {
        rooms: "",
        smartLock: false,
        status: "",
        sort: ""
      },
      roomOptions: [
        { value: "", label: "不限" },
        { value: "1", label: "1室" },
        { value: "2", label: "2室" },
        { value: "3", label: "3室" },
        { value: "4+", label: "4室+" }
      ],
      statusOptions: [
        { value: "", label: "不限" },
        { value: "normal", label: "正常" },
        { value: "交易中", label: "交易中" },
        { value: "已下架", label: "已下架" }
      ],
      sortOptions: [
        { value: "", label: "默认" },
        { value: "views_desc", label: "浏览最多" },
        { value: "myviews_desc", label: "我最关注" }
      ],
      propertyList: [
        {
          title: "13# 12.12开盘 核心卖点 海景高层",
          rooms: "3室2厅",
          size: "105",
          orientation: "南向",
          price: "850万",
          unitPrice: "80,952元/m²",
          views: 127,
          myViews: 1,
          tag: "智能门锁",
          tagIcon: "lock",
          image: "/static/images/img_c63bd4d9ba.png",
          tags: [
            { name: "自营", type: "blue" },
            { name: "精装", type: "green" },
            { name: "学区房", type: "orange" }
          ],
          commission: "佣金1.5%，成交奖励经纪人8000"
        },
        {
          title: "金城中心 舒适单身公寓 地铁口直达",
          rooms: "1室1厅",
          size: "45",
          orientation: "北向",
          price: "210万",
          unitPrice: "46,666元/m²",
          views: 89,
          myViews: 0,
          image: "/static/images/img_f0b09ae08c.png",
          tags: [
            { name: "地铁房", type: "indigo" },
            { name: "有电梯", type: "grey" }
          ],
          commission: "佣金1%，极速结佣",
          footerIcon: "percent"
        },
        {
          title: "滨江别墅 独栋花园 视野开阔 无遮挡",
          rooms: "4室3厅",
          size: "210",
          orientation: "东向",
          price: "1280万",
          unitPrice: "60,952元/m²",
          views: 256,
          myViews: 45,
          tag: "VR看房",
          tagIcon: "videocam",
          image: "/static/images/img_f1b597a564.png",
          tags: [
            { name: "自营", type: "blue" },
            { name: "VR看房", type: "purple" }
          ],
          commission: "高佣2%，房东急售",
          footerIcon: "local_fire_department"
        },
        {
          title: "日落大道 2室紧凑型 适合过渡",
          rooms: "2室1厅",
          size: "78",
          orientation: "西向",
          price: "420万",
          unitPrice: "53,846元/m²",
          views: 45,
          myViews: 2,
          status: "交易中",
          image: "/static/images/img_fdf5f4bc67.png",
          tags: [
            { name: "随时看房", type: "grey" }
          ],
          commission: "佣金1%",
          footerIcon: "history",
          footerClass: "grey-footer"
        },
        {
          title: "日落大道 2室紧凑型 适合过渡",
          rooms: "2室1厅",
          size: "78",
          orientation: "西向",
          price: "420万",
          unitPrice: "53,846元/m²",
          views: 45,
          myViews: 2,
          status: "已下架",
          image: "/static/images/img_fdf5f4bc67.png",
          tags: [
            { name: "随时看房", type: "grey" }
          ],
          commission: "佣金1%",
          footerIcon: "history",
          footerClass: "grey-footer"
        }
      ]
    };
  },
  computed: {
    hasActiveFilters() {
      return !!(this.filters.rooms || this.filters.smartLock || this.filters.status || this.filters.sort);
    },
    displayList() {
      let list = Array.isArray(this.propertyList) ? [...this.propertyList] : [];
      if (this.filters.rooms) {
        list = list.filter((it) => {
          const s = it && it.rooms || "";
          const m = String(s).match(/^(\d+)/);
          const n = m ? Number(m[1]) : 0;
          if (this.filters.rooms === "4+")
            return n >= 4;
          return String(n) === this.filters.rooms;
        });
      }
      if (this.filters.smartLock) {
        list = list.filter((it) => it && (it.tag === "智能门锁" || it.tagIcon === "lock"));
      }
      if (this.filters.status) {
        if (this.filters.status === "normal") {
          list = list.filter((it) => !it.status);
        } else {
          list = list.filter((it) => it && it.status === this.filters.status);
        }
      }
      if (this.filters.sort === "views_desc") {
        list.sort((a, b) => (b.views || 0) - (a.views || 0));
      } else if (this.filters.sort === "myviews_desc") {
        list.sort((a, b) => (b.myViews || 0) - (a.myViews || 0));
      }
      return list;
    }
  },
  methods: {
    setCategory(key) {
      this.currentCategory = key;
    },
    handlePromote(item) {
      const title = item && item.title || "房源";
      common_vendor.index$1.showActionSheet({
        itemList: ["生成获客海报", "复制推广文案", "模拟分享（占位）"],
        success: (res) => {
          if (res.tapIndex === 0) {
            common_vendor.index$1.navigateTo({
              url: `/pages/marketing_posters/marketing_posters?title=${encodeURIComponent(title)}`
            });
            return;
          }
          if (res.tapIndex === 1) {
            const text = `【优质房源推荐】${title}
点击查看详情（示例）：/pages/property_detail/property_detail`;
            common_vendor.index$1.setClipboardData({
              data: text,
              success: () => {
                common_vendor.index$1.showToast({ title: "已复制推广文案", icon: "none" });
              }
            });
            return;
          }
          if (res.tapIndex === 2) {
            common_vendor.index$1.showToast({ title: "分享功能待接入", icon: "none" });
          }
        }
      });
    },
    openFilter() {
      this.filterOpen = true;
      this.draft = {
        rooms: this.filters.rooms,
        smartLock: !!this.filters.smartLock,
        status: this.filters.status,
        sort: this.filters.sort
      };
    },
    closeFilter() {
      this.filterOpen = false;
    },
    resetFilter() {
      this.draft = { rooms: "", smartLock: false, status: "", sort: "" };
    },
    applyFilter() {
      this.filters = {
        rooms: this.draft.rooms,
        smartLock: !!this.draft.smartLock,
        status: this.draft.status,
        sort: this.draft.sort
      };
      this.closeFilter();
    },
    goToDetail(item) {
      if (item.status === "已下架")
        return;
      common_vendor.index$1.navigateTo({
        url: "/pages/property_detail/property_detail"
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
  return common_vendor.e({
    a: common_vendor.p({
      title: "房源"
    }),
    b: $options.hasActiveFilters
  }, $options.hasActiveFilters ? {} : {}, {
    c: $options.hasActiveFilters ? 1 : "",
    d: common_vendor.o((...args) => $options.openFilter && $options.openFilter(...args)),
    e: common_vendor.f($data.categories, (c, idx, i0) => {
      return {
        a: common_vendor.t(c.label),
        b: c.key || idx,
        c: $data.currentCategory === c.key ? 1 : "",
        d: common_vendor.o(($event) => $options.setCategory(c.key), c.key || idx)
      };
    }),
    f: common_vendor.f($options.displayList, (item, index, i0) => {
      return common_vendor.e({
        a: item.image,
        b: item.tag
      }, item.tag ? {
        c: common_vendor.t(item.tagIcon),
        d: common_vendor.t(item.tag)
      } : {}, {
        e: item.status === "交易中" || item.status === "已下架"
      }, item.status === "交易中" || item.status === "已下架" ? {
        f: common_vendor.t(item.status)
      } : {}, {
        g: common_vendor.t(item.title),
        h: common_vendor.t(item.rooms),
        i: common_vendor.t(item.size),
        j: common_vendor.t(item.orientation),
        k: common_vendor.f(item.tags, (tag, tIdx, i1) => {
          return {
            a: common_vendor.t(tag.name),
            b: tIdx,
            c: common_vendor.n(tag.type)
          };
        }),
        l: common_vendor.t(item.price),
        m: common_vendor.t(item.unitPrice),
        n: common_vendor.t(item.views),
        o: common_vendor.t(item.myViews),
        p: common_vendor.t(item.footerIcon || "currency_yen"),
        q: common_vendor.t(item.commission),
        r: item.status !== "已下架"
      }, item.status !== "已下架" ? {
        s: common_vendor.o(($event) => $options.handlePromote(item), index)
      } : {}, {
        t: common_vendor.n(item.footerClass || "orange-footer"),
        v: index,
        w: item.status === "已下架" ? 1 : "",
        x: common_vendor.o(($event) => $options.goToDetail(item), index)
      });
    }),
    g: common_vendor.p({
      active: "property"
    }),
    h: $data.filterOpen
  }, $data.filterOpen ? {
    i: common_vendor.o((...args) => $options.resetFilter && $options.resetFilter(...args)),
    j: common_vendor.f($data.roomOptions, (opt, idx, i0) => {
      return {
        a: common_vendor.t(opt.label),
        b: "room-" + idx,
        c: $data.draft.rooms === opt.value ? 1 : "",
        d: common_vendor.o(($event) => $data.draft.rooms = opt.value, "room-" + idx)
      };
    }),
    k: !!$data.draft.smartLock ? 1 : "",
    l: common_vendor.o(($event) => $data.draft.smartLock = !$data.draft.smartLock),
    m: common_vendor.f($data.statusOptions, (opt, idx, i0) => {
      return {
        a: common_vendor.t(opt.label),
        b: "status-" + idx,
        c: $data.draft.status === opt.value ? 1 : "",
        d: common_vendor.o(($event) => $data.draft.status = opt.value, "status-" + idx)
      };
    }),
    n: common_vendor.f($data.sortOptions, (opt, idx, i0) => {
      return common_vendor.e({
        a: common_vendor.t(opt.label),
        b: $data.draft.sort === opt.value
      }, $data.draft.sort === opt.value ? {} : {}, {
        c: "sort-" + idx,
        d: $data.draft.sort === opt.value ? 1 : "",
        e: common_vendor.o(($event) => $data.draft.sort = opt.value, "sort-" + idx)
      });
    }),
    o: common_vendor.o((...args) => $options.closeFilter && $options.closeFilter(...args)),
    p: common_vendor.o((...args) => $options.applyFilter && $options.applyFilter(...args)),
    q: common_vendor.o(() => {
    }),
    r: common_vendor.o((...args) => $options.closeFilter && $options.closeFilter(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/property_list/property_list.js.map
