"use strict";
const common_vendor = require("../../common/vendor.js");
const BottomTabBar = () => "../../components/BottomTabBar.js";
const TopHeader = () => "../../components/TopHeader.js";
const _sfc_main = {
  components: { BottomTabBar, TopHeader },
  onLoad() {
    const cached = common_vendor.index$1.getStorageSync("hm_province_name");
    this.provinceName = cached || "辽宁";
    common_vendor.index$1.getLocation({
      type: "wgs84",
      success: () => {
      },
      fail: () => {
        this.provinceName = this.provinceName || "辽宁";
      }
    });
  },
  data() {
    return {
      // 省份（左上角）
      provinceName: "辽宁",
      provinceOptions: ["辽宁", "北京", "上海", "广东", "浙江", "江苏"],
      // 市/区二级联动（左侧抽屉）
      locationOpen: false,
      locationCities: [
        {
          key: "shenyang",
          label: "沈阳",
          districts: [
            { key: "heping", label: "和平区" },
            { key: "shenhe", label: "沈河区" },
            { key: "huanggu", label: "皇姑区" },
            { key: "tiexi", label: "铁西区" }
          ]
        },
        {
          key: "dalian",
          label: "大连",
          districts: [
            { key: "zhongshan", label: "中山区" },
            { key: "xigang", label: "西岗区" },
            { key: "shahekou", label: "沙河口区" },
            { key: "ganjingzi", label: "甘井子区" }
          ]
        },
        {
          key: "anshan",
          label: "鞍山",
          districts: [
            { key: "tiedong", label: "铁东区" },
            { key: "tiexi", label: "铁西区" },
            { key: "lixian", label: "立山区" }
          ]
        }
      ],
      locDraft: { city: "", district: "" },
      // 筛选栏（静态数据 + 交互）
      filters: {
        location: "",
        filter: "",
        sort: "",
        more: []
      },
      panelOpen: false,
      panelType: "",
      panelDraft: {
        location: "",
        filter: "",
        sort: "",
        more: []
      },
      filterOptions: {
        filter: [
          { value: "", label: "筛选" },
          { value: "smart_lock", label: "智能锁" },
          { value: "near_subway", label: "近地铁" },
          { value: "new", label: "新上" },
          { value: "price_drop", label: "降价" }
        ],
        sort: [
          { value: "", label: "排序" },
          { value: "recommend", label: "推荐优先" },
          { value: "price_asc", label: "租金从低到高" },
          { value: "price_desc", label: "租金从高到低" },
          { value: "latest", label: "最新发布" }
        ],
        more: [
          { value: "whole_rent", label: "整租" },
          { value: "shared", label: "合租" },
          { value: "elevator", label: "有电梯" },
          { value: "south", label: "朝南" },
          { value: "pet", label: "可养宠" }
        ]
      },
      bannerImages: [
        // 示例图：可替换为你的 Banner 图地址
        "/static/images/img_28f6b412d7.png",
        "/static/images/img_2abd9934e1.png",
        "/static/images/img_48afe8538f.png"
      ],
      followedProperties: [
        {
          name: "阳光花园 · 2室1厅",
          area: "朝阳区",
          size: "85",
          price: "6,500",
          type: "住宅",
          priceChange: "降价 ¥200",
          image: "/static/images/img_48afe8538f.png"
        },
        {
          name: "金域华府 A2",
          area: "东城区",
          size: "65",
          price: "5,800",
          status: "待看",
          image: "/static/images/img_cd1f0409ce.png"
        }
      ],
      recommendedProperties: [
        {
          name: "远洋万和城 · 3室2厅",
          price: "9,200",
          region: "海淀区",
          size: "120",
          desc: "近公园",
          status: "急租",
          tags: ["近地铁", "新上"],
          image: "/static/images/img_28f6b412d7.png"
        },
        {
          name: "柳溪公寓 · 开间",
          price: "4,100",
          region: "丰台区",
          size: "45",
          desc: "西向",
          tags: ["随时入住", "免中介费"],
          image: "/static/images/img_f7b0d6455a.png"
        }
      ]
    };
  },
  computed: {
    filterLabelLocation() {
      if (!this.filters.location)
        return "市/区";
      const [cityKey, districtKey] = String(this.filters.location).split("|");
      const city = this.locationCities.find((c) => c.key === cityKey);
      if (!city)
        return "市/区";
      const district = (city.districts || []).find((d) => d.key === districtKey);
      return district ? `${city.label} · ${district.label}` : city.label;
    },
    currentDistricts() {
      const city = this.locationCities.find((c) => c.key === this.locDraft.city) || this.locationCities[0];
      return city && city.districts || [];
    },
    filterLabelFilter() {
      const hit = this.filterOptions.filter.find((x) => x.value === this.filters.filter);
      return hit ? hit.label : "筛选";
    },
    filterLabelSort() {
      const hit = this.filterOptions.sort.find((x) => x.value === this.filters.sort);
      return hit ? hit.label : "排序";
    },
    filterLabelMore() {
      if (!this.filters.more || this.filters.more.length === 0)
        return "更多";
      return `更多(${this.filters.more.length})`;
    },
    panelTitle() {
      if (this.panelType === "location")
        return "选择位置";
      if (this.panelType === "filter")
        return "选择筛选";
      if (this.panelType === "sort")
        return "选择排序";
      if (this.panelType === "more")
        return "更多条件";
      return "";
    },
    panelOptions() {
      return this.filterOptions && this.filterOptions[this.panelType] || [];
    }
  },
  methods: {
    // 省份选择
    selectProvince() {
      common_vendor.index$1.showActionSheet({
        itemList: this.provinceOptions,
        success: (res) => {
          const name = this.provinceOptions[res.tapIndex];
          if (!name)
            return;
          this.provinceName = name;
          common_vendor.index$1.setStorageSync("hm_province_name", name);
        }
      });
    },
    // 市/区：左侧二级联动
    openLocationPicker() {
      this.locationOpen = true;
      const [cityKey, districtKey] = String(this.filters.location || "").split("|");
      this.locDraft.city = cityKey || this.locationCities[0] && this.locationCities[0].key || "";
      this.locDraft.district = districtKey || "";
      if (!this.locDraft.district) {
        const city = this.locationCities.find((c) => c.key === this.locDraft.city) || this.locationCities[0];
        const first = city && city.districts && city.districts[0];
        this.locDraft.district = first && first.key || "";
      }
    },
    closeLocationPicker() {
      this.locationOpen = false;
    },
    resetLocationPicker() {
      const city = this.locationCities[0];
      this.locDraft.city = city && city.key || "";
      this.locDraft.district = city && city.districts && city.districts[0] && city.districts[0].key || "";
    },
    selectLocCity(cityKey) {
      this.locDraft.city = cityKey;
      const city = this.locationCities.find((c) => c.key === cityKey) || this.locationCities[0];
      const first = city && city.districts && city.districts[0];
      this.locDraft.district = first && first.key || "";
    },
    applyLocationPicker() {
      if (!this.locDraft.city) {
        this.filters.location = "";
      } else {
        this.filters.location = `${this.locDraft.city}|${this.locDraft.district || ""}`;
      }
      this.closeLocationPicker();
    },
    openPanel(type) {
      this.panelType = type;
      this.panelOpen = true;
      this.panelDraft = {
        location: this.filters.location,
        filter: this.filters.filter,
        sort: this.filters.sort,
        more: Array.isArray(this.filters.more) ? [...this.filters.more] : []
      };
    },
    closePanel() {
      this.panelOpen = false;
      this.panelType = "";
    },
    resetPanel() {
      if (this.panelType === "location")
        this.panelDraft.location = "";
      if (this.panelType === "filter")
        this.panelDraft.filter = "";
      if (this.panelType === "sort")
        this.panelDraft.sort = "";
      if (this.panelType === "more")
        this.panelDraft.more = [];
    },
    applyPanel() {
      this.filters = {
        location: this.panelDraft.location,
        filter: this.panelDraft.filter,
        sort: this.panelDraft.sort,
        more: Array.isArray(this.panelDraft.more) ? [...this.panelDraft.more] : []
      };
      this.closePanel();
    },
    isOptionActive(opt) {
      if (!opt)
        return false;
      if (this.panelType === "location")
        return this.panelDraft.location === opt.value;
      if (this.panelType === "filter")
        return this.panelDraft.filter === opt.value;
      if (this.panelType === "sort")
        return this.panelDraft.sort === opt.value;
      return false;
    },
    selectSingle(opt) {
      if (!opt)
        return;
      if (this.panelType === "location")
        this.panelDraft.location = opt.value;
      if (this.panelType === "filter")
        this.panelDraft.filter = opt.value;
      if (this.panelType === "sort")
        this.panelDraft.sort = opt.value;
    },
    isMoreSelected(opt) {
      if (!opt)
        return false;
      return (this.panelDraft.more || []).includes(opt.value);
    },
    toggleMore(opt) {
      if (!opt)
        return;
      const list = Array.isArray(this.panelDraft.more) ? this.panelDraft.more : [];
      const idx = list.indexOf(opt.value);
      if (idx >= 0)
        list.splice(idx, 1);
      else
        list.push(opt.value);
      this.panelDraft.more = list;
    },
    selectCity() {
      common_vendor.index$1.__f__("log", "at pages/home/home.vue:500", "选择城市");
    },
    viewAllFollowed() {
      common_vendor.index$1.navigateTo({
        url: "/pages/property_list/property_list"
      });
    },
    goToDetail(item) {
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
      title: "推荐"
    }),
    b: common_vendor.t($data.provinceName),
    c: common_vendor.o((...args) => $options.selectProvince && $options.selectProvince(...args)),
    d: $data.filters.more.length > 0
  }, $data.filters.more.length > 0 ? {} : {}, {
    e: $data.filters.more.length > 0 ? 1 : "",
    f: common_vendor.o(($event) => $options.openPanel("more")),
    g: common_vendor.f($data.bannerImages, (img, idx, i0) => {
      return {
        a: img,
        b: idx
      };
    }),
    h: common_vendor.t($options.filterLabelLocation),
    i: !!$data.filters.location ? 1 : "",
    j: common_vendor.o((...args) => $options.openLocationPicker && $options.openLocationPicker(...args)),
    k: common_vendor.t($options.filterLabelFilter),
    l: !!$data.filters.filter ? 1 : "",
    m: common_vendor.o(($event) => $options.openPanel("filter")),
    n: common_vendor.t($options.filterLabelSort),
    o: !!$data.filters.sort ? 1 : "",
    p: common_vendor.o(($event) => $options.openPanel("sort")),
    q: common_vendor.o((...args) => $options.viewAllFollowed && $options.viewAllFollowed(...args)),
    r: common_vendor.f($data.followedProperties, (item, index, i0) => {
      return common_vendor.e({
        a: common_vendor.t(item.type),
        b: item.priceChange
      }, item.priceChange ? {
        c: common_vendor.t(item.priceChange)
      } : {}, {
        d: item.status
      }, item.status ? {
        e: common_vendor.t(item.status)
      } : {}, {
        f: "url(" + item.image + ")",
        g: common_vendor.t(item.name),
        h: common_vendor.t(item.area),
        i: common_vendor.t(item.size),
        j: common_vendor.t(item.price),
        k: index,
        l: common_vendor.o(($event) => $options.goToDetail(item), index)
      });
    }),
    s: common_vendor.f($data.recommendedProperties, (item, index, i0) => {
      return common_vendor.e({
        a: item.status
      }, item.status ? {
        b: common_vendor.t(item.status)
      } : {}, {
        c: "url(" + item.image + ")",
        d: common_vendor.t(item.name),
        e: common_vendor.t(item.price),
        f: common_vendor.t(item.region),
        g: common_vendor.t(item.size),
        h: common_vendor.t(item.desc),
        i: common_vendor.f(item.tags, (tag, tIdx, i1) => {
          return {
            a: common_vendor.t(tag),
            b: tIdx,
            c: tag === "新上" ? 1 : ""
          };
        }),
        j: index,
        k: common_vendor.o(($event) => $options.goToDetail(item), index)
      });
    }),
    t: common_vendor.p({
      active: "home"
    }),
    v: $data.panelOpen
  }, $data.panelOpen ? common_vendor.e({
    w: common_vendor.t($options.panelTitle),
    x: common_vendor.o((...args) => $options.resetPanel && $options.resetPanel(...args)),
    y: $data.panelType !== "more"
  }, $data.panelType !== "more" ? {
    z: common_vendor.f($options.panelOptions, (opt, idx, i0) => {
      return common_vendor.e({
        a: common_vendor.t(opt.label),
        b: $options.isOptionActive(opt)
      }, $options.isOptionActive(opt) ? {} : {}, {
        c: idx,
        d: $options.isOptionActive(opt) ? 1 : "",
        e: common_vendor.o(($event) => $options.selectSingle(opt), idx)
      });
    })
  } : {
    A: common_vendor.f($options.panelOptions, (opt, idx, i0) => {
      return common_vendor.e({
        a: common_vendor.t(opt.label),
        b: $options.isMoreSelected(opt)
      }, $options.isMoreSelected(opt) ? {} : {}, {
        c: idx,
        d: $options.isMoreSelected(opt) ? 1 : "",
        e: common_vendor.o(($event) => $options.toggleMore(opt), idx)
      });
    })
  }, {
    B: common_vendor.o((...args) => $options.closePanel && $options.closePanel(...args)),
    C: common_vendor.o((...args) => $options.applyPanel && $options.applyPanel(...args)),
    D: common_vendor.o(() => {
    }),
    E: common_vendor.o((...args) => $options.closePanel && $options.closePanel(...args))
  }) : {}, {
    F: $data.locationOpen
  }, $data.locationOpen ? {
    G: common_vendor.o((...args) => $options.resetLocationPicker && $options.resetLocationPicker(...args)),
    H: common_vendor.f($data.locationCities, (c, k0, i0) => {
      return {
        a: common_vendor.t(c.label),
        b: c.key,
        c: $data.locDraft.city === c.key ? 1 : "",
        d: common_vendor.o(($event) => $options.selectLocCity(c.key), c.key)
      };
    }),
    I: common_vendor.f($options.currentDistricts, (d, k0, i0) => {
      return {
        a: common_vendor.t(d.label),
        b: d.key,
        c: $data.locDraft.district === d.key ? 1 : "",
        d: common_vendor.o(($event) => $data.locDraft.district = d.key, d.key)
      };
    }),
    J: common_vendor.o((...args) => $options.closeLocationPicker && $options.closeLocationPicker(...args)),
    K: common_vendor.o((...args) => $options.applyLocationPicker && $options.applyLocationPicker(...args)),
    L: common_vendor.o(() => {
    }),
    M: common_vendor.o((...args) => $options.closeLocationPicker && $options.closeLocationPicker(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/home/home.js.map
