<template>
  <view class="home-container">
    <!-- 顶部导航：使用 TopHeader 统一计算状态栏/胶囊高度，避免贴顶且不遮挡 -->
    <view class="header">
      <TopHeader title="推荐" />
    </view>

    <scroll-view scroll-y="true" class="main-content">
      <!-- 地区选择 + 搜索：放在 Banner 上方 -->
      <view class="top-tools">
        <view class="top-tools__inner">
          <!-- 左上角：省份选择（优先定位，失败默认辽宁） -->
          <view class="city-selector" @click="selectProvince">
            <text class="city-name">{{ provinceName }}</text>
            <text class="material-symbols-outlined expand-icon"
              >expand_more</text
            >
          </view>
          <view class="search-box">
            <text class="material-symbols-outlined search-icon">search</text>
            <input
              v-model="keyword"
              class="search-input"
              type="text"
              confirm-type="search"
              placeholder="搜索区、小区、地址..."
              placeholder-class="placeholder"
              @confirm="onSearchConfirm"
            />
          </view>
          <!-- 更多筛选：移到搜索栏右侧，避免下方文字拥挤 -->
          <view
            class="more-btn"
            :class="{ active: filters.more.length > 0 }"
            @click="openPanel('more')"
          >
            <text class="material-symbols-outlined">tune</text>
            <view v-if="filters.more.length > 0" class="more-dot"></view>
          </view>
        </view>
      </view>

      <!-- Banner -->
      <view class="banner-section">
        <view class="banner-card">
          <swiper class="banner-swiper" circular :indicator-dots="false">
            <swiper-item v-for="(img, idx) in bannerImages" :key="idx">
              <image class="banner-image" :src="img" mode="aspectFill"></image>
            </swiper-item>
          </swiper>
        </view>
      </view>

      <!-- 筛选条 (这里简单实现，不处理吸顶逻辑，UniApp中可以用sticky) -->
      <view class="filter-bar">
        <!-- 市 / 区：左侧二级联动选择 -->
        <view
          class="filter-item"
          :class="{ active: !!filters.location }"
          @click="openLocationPicker"
        >
          <text class="filter-text">{{ filterLabelLocation }}</text>
          <text class="material-symbols-outlined filter-icon">expand_more</text>
        </view>
        <view
          class="filter-item"
          :class="{ active: !!filters.filter }"
          @click="openPanel('filter')"
        >
          <text class="filter-text">{{ filterLabelFilter }}</text>
          <text class="material-symbols-outlined filter-icon">expand_more</text>
        </view>
        <view
          class="filter-item"
          :class="{ active: !!filters.sort }"
          @click="openPanel('sort')"
        >
          <text class="filter-text">{{ filterLabelSort }}</text>
          <text class="material-symbols-outlined filter-icon">swap_vert</text>
        </view>
      </view>

      <!-- 我关注的房源 -->
      <view class="section">
        <view class="section-header">
          <view class="section-title">
            <text class="material-symbols-outlined section-icon accent"
              >favorite</text
            >
            <text class="title-text">我关注的房源</text>
          </view>
          <view class="section-more" @click="viewAllFollowed">
            <text>全部 ({{ followedTotal }})</text>
            <text class="material-symbols-outlined more-icon"
              >chevron_right</text
            >
          </view>
        </view>
        <view
          v-if="!followedProperties || followedProperties.length === 0"
          class="followed-empty"
        >
          <view class="empty-card">
            <text class="material-symbols-outlined empty-icon">favorite</text>
            <view class="empty-title">暂无关注房源</view>
            <view class="empty-desc">在房源详情页点击关注后，会在这里展示</view>
            <view class="empty-actions">
              <button class="empty-btn ghost" @click="viewAllFollowed">
                去逛逛
              </button>
            </view>
          </view>
        </view>
        <scroll-view v-else scroll-x="true" class="horizontal-scroll">
          <view
            class="property-card-h"
            v-for="(item, index) in followedProperties"
            :key="index"
            @click="goToDetail(item)"
          >
            <view
              class="card-image"
              :class="{ 'is-disabled': isSoldOrOff(item.sale_status) }"
              :style="{ backgroundImage: 'url(' + item.image + ')' }"
            >
              <view class="image-tag">{{
                item.property_type || item.type
              }}</view>
              <view class="price-change-tag" v-if="item.priceChange">{{
                item.priceChange
              }}</view>
              <view
                class="status-tag"
                v-if="getSaleStatusLabel(item.sale_status)"
                :class="getSaleStatusClass(item.sale_status)"
              >
                {{ getSaleStatusLabel(item.sale_status) }}
              </view>
              <view class="lock-tag" v-if="Number(item.has_smart_lock) === 1"
                >智能锁</view
              >
            </view>
            <view class="card-info">
              <view class="property-name">{{ item.name }}</view>
              <view class="mini-tags" v-if="ensureTags(item.tags).length">
                <text
                  class="mini-tag"
                  v-for="(tag, tIdx) in ensureTags(item.tags).slice(0, 2)"
                  :key="tIdx"
                  >{{ tag }}</text
                >
              </view>
              <view class="property-meta">
                <text>{{ formatFollowedMeta(item) }}</text>
                <view class="price">
                  <text class="price-val">¥{{ item.price }}</text>
                  <text class="price-unit">{{ getPriceUnit(item) }}</text>
                </view>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 系统推荐 -->
      <view class="section">
        <view class="section-header">
          <view class="section-title">
            <text class="material-symbols-outlined section-icon primary"
              >recommend</text
            >
            <text class="title-text">系统推荐的房源</text>
          </view>
        </view>
        <view class="vertical-list">
          <view
            class="property-card-v"
            v-for="(item, index) in recommendedProperties"
            :key="index"
            @click="goToDetail(item)"
          >
            <view
              class="card-image-v"
              :class="{ 'is-disabled': isSoldOrOff(item.sale_status) }"
              :style="{ backgroundImage: 'url(' + item.image + ')' }"
            >
              <view
                class="status-tag-v"
                v-if="getSaleStatusLabel(item.sale_status)"
                :class="getSaleStatusClass(item.sale_status)"
              >
                {{ getSaleStatusLabel(item.sale_status) }}
              </view>
              <view class="lock-tag-v" v-if="Number(item.has_smart_lock) === 1"
                >智能锁</view
              >
            </view>
            <view class="card-info-v">
              <view class="info-top">
                <text class="property-name-v">{{ item.name }}</text>
                <view class="price-v">
                  <text class="price-val">¥{{ item.price }}</text>
                  <text class="price-unit">{{ getPriceUnit(item) }}</text>
                </view>
              </view>
              <view class="property-location">
                <text class="material-symbols-outlined loc-icon"
                  >location_on</text
                >
                <text>{{ formatRecommendedMeta(item) }}</text>
              </view>
              <view class="tags">
                <text
                  class="tag-item"
                  v-for="(tag, tIdx) in ensureTags(item.tags)"
                  :key="tIdx"
                  :class="{ 'tag-orange': tag === '新上' }"
                >
                  {{ tag }}
                </text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部导航栏 -->
    <BottomTabBar active="home" />

    <!-- 筛选弹层（静态数据 + 交互） -->
    <view v-if="panelOpen" class="mask" @click="closePanel">
      <view class="sheet" @click.stop>
        <view class="sheet-header">
          <text class="sheet-title">{{ panelTitle }}</text>
          <text class="sheet-action" @click="resetPanel">重置</text>
        </view>

        <!-- 单选：位置/筛选/排序 -->
        <view v-if="panelType !== 'more'" class="sheet-list">
          <view
            v-for="(opt, idx) in panelOptions"
            :key="idx"
            class="sheet-item"
            :class="{ active: isOptionActive(opt) }"
            @click="selectSingle(opt)"
          >
            <text class="sheet-item-text">{{ opt.label }}</text>
            <text
              v-if="isOptionActive(opt)"
              class="material-symbols-outlined sheet-check"
              >check</text
            >
          </view>
        </view>

        <!-- 多选：更多 -->
        <view v-else class="sheet-list">
          <view
            v-for="(opt, idx) in panelOptions"
            :key="idx"
            class="sheet-item"
            :class="{ active: isMoreSelected(opt) }"
            @click="toggleMore(opt)"
          >
            <text class="sheet-item-text">{{ opt.label }}</text>
            <text
              v-if="isMoreSelected(opt)"
              class="material-symbols-outlined sheet-check"
              >check</text
            >
          </view>
        </view>

        <view class="sheet-footer">
          <button class="sheet-btn ghost" @click="closePanel">取消</button>
          <button class="sheet-btn primary" @click="applyPanel">确定</button>
        </view>
      </view>
    </view>

    <!-- 市/区：左侧二级联动抽屉 -->
    <view v-if="locationOpen" class="loc-mask" @click="closeLocationPicker">
      <view class="loc-drawer" @click.stop>
        <view class="loc-head">
          <text class="loc-title">选择城市 / 区县</text>
          <text class="loc-reset" @click="resetLocationPicker">重置</text>
        </view>
        <view class="loc-body">
          <scroll-view scroll-y class="loc-col loc-left">
            <view
              v-for="c in locationCities"
              :key="c.key"
              class="loc-item"
              :class="{ active: locDraft.city === c.key }"
              @click="selectLocCity(c.key)"
            >
              {{ c.label }}
            </view>
          </scroll-view>
          <scroll-view scroll-y class="loc-col loc-right">
            <view
              v-for="d in currentDistricts"
              :key="d.key"
              class="loc-item"
              :class="{ active: locDraft.district === d.key }"
              @click="locDraft.district = d.key"
            >
              {{ d.label }}
            </view>
          </scroll-view>
        </view>
        <view class="loc-foot">
          <button class="loc-btn ghost" @click="closeLocationPicker">
            取消
          </button>
          <button class="loc-btn primary" @click="applyLocationPicker">
            确定
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import BottomTabBar from "@/components/BottomTabBar.vue";
import TopHeader from "@/components/TopHeader.vue";
import homeApi from "@/api/home";

export default {
  components: { BottomTabBar, TopHeader },
  onLoad() {
    // 省份：优先读取缓存；定位失败默认辽宁（不做逆地理，避免引入第三方 key/域名）
    const cached = uni.getStorageSync("hm_province_name");
    this.provinceName = cached || "辽宁";
    // #ifdef MP-WEIXIN
    uni.getLocation({
      type: "wgs84",
      success: () => {},
      fail: () => {
        this.provinceName = this.provinceName || "辽宁";
      },
    });
    // #endif
    this.loadHomeData();
  },
  data() {
    return {
      // 省份（左上角）
      provinceName: "辽宁",
      provinceOptions: ["辽宁", "北京", "上海", "广东", "浙江", "江苏"],

      // 搜索关键词
      keyword: "",

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
            { key: "tiexi", label: "铁西区" },
          ],
        },
        {
          key: "dalian",
          label: "大连",
          districts: [
            { key: "zhongshan", label: "中山区" },
            { key: "xigang", label: "西岗区" },
            { key: "shahekou", label: "沙河口区" },
            { key: "ganjingzi", label: "甘井子区" },
          ],
        },
        {
          key: "anshan",
          label: "鞍山",
          districts: [
            { key: "tiedong", label: "铁东区" },
            { key: "tiexi", label: "铁西区" },
            { key: "lixian", label: "立山区" },
          ],
        },
      ],
      locDraft: { city: "", district: "" },

      // 筛选栏（静态数据 + 交互）
      filters: {
        location: "",
        filter: "",
        sort: "",
        more: [],
      },
      panelOpen: false,
      panelType: "",
      panelDraft: {
        location: "",
        filter: "",
        sort: "",
        more: [],
      },
      filterOptions: {
        filter: [
          { value: "", label: "筛选" },
          { value: "smart_lock", label: "智能锁" },
          { value: "near_subway", label: "近地铁" },
          { value: "new", label: "新上" },
          { value: "price_drop", label: "降价" },
        ],
        sort: [
          { value: "", label: "排序" },
          { value: "recommend", label: "推荐优先" },
          { value: "price_asc", label: "价格从低到高" },
          { value: "price_desc", label: "价格从高到低" },
          { value: "latest", label: "最新发布" },
        ],
        more: [
          { value: "whole_rent", label: "整租" },
          { value: "shared", label: "合租" },
          { value: "elevator", label: "有电梯" },
          { value: "south", label: "朝南" },
          { value: "pet", label: "可养宠" },
        ],
      },
      bannerImages: [
        // 示例图：可替换为你的 Banner 图地址
        "/static/images/img_28f6b412d7.png",
        "/static/images/img_2abd9934e1.png",
        "/static/images/img_48afe8538f.png",
      ],
      followedTotal: 0,
      followedProperties: [
        {
          name: "阳光花园 · 2室1厅",
          area: "朝阳区",
          size: "85",
          price: "6,500",
          type: "住宅",
          priceChange: "降价 ¥200",
          image: "/static/images/img_48afe8538f.png",
        },
        {
          name: "金域华府 A2",
          area: "东城区",
          size: "65",
          price: "5,800",
          status: "待看",
          image: "/static/images/img_cd1f0409ce.png",
        },
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
          image: "/static/images/img_28f6b412d7.png",
        },
        {
          name: "柳溪公寓 · 开间",
          price: "4,100",
          region: "丰台区",
          size: "45",
          desc: "西向",
          tags: ["随时入住", "免中介费"],
          image: "/static/images/img_f7b0d6455a.png",
        },
      ],
    };
  },
  computed: {
    filterLabelLocation() {
      if (!this.filters.location) return "市/区";
      const [cityKey, districtKey] = String(this.filters.location).split("|");
      const city = this.locationCities.find((c) => c.key === cityKey);
      if (!city) return "市/区";
      const district = (city.districts || []).find(
        (d) => d.key === districtKey
      );
      return district ? `${city.label} · ${district.label}` : city.label;
    },
    /**
     * 获取当前选中城市的所有区县列表
     *
     * @returns {Array} 当前选中城市的区县数组，如果没有城市则返回空数组
     */
    currentDistricts() {
      const city =
        this.locationCities.find((c) => c.key === this.locDraft.city) ||
        this.locationCities[0];
      return (city && city.districts) || [];
    },
    /**
     * 获取筛选条件的显示标签
     *
     * 遍历筛选选项列表，查找与当前筛选值匹配的选项，
     * 返回对应的显示标签。如果未找到匹配项，则返回默认的'筛选'文本。
     *
     * @returns {string} 筛选条件的显示标签，如果没有匹配项则返回'筛选'
     */
    filterLabelFilter() {
      const hit = this.filterOptions.filter.find(
        (x) => x.value === this.filters.filter
      );
      return hit ? hit.label : "筛选";
    },
    /**
     * 获取排序条件的显示标签
     *
     * @returns {string} 排序条件的显示标签，如果没有匹配项则返回'排序'
     */
    filterLabelSort() {
      const hit = this.filterOptions.sort.find(
        (x) => x.value === this.filters.sort
      );
      return hit ? hit.label : "排序";
    },
    filterLabelMore() {
      if (!this.filters.more || this.filters.more.length === 0) return "更多";
      return `更多(${this.filters.more.length})`;
    },
    panelTitle() {
      if (this.panelType === "location") return "选择位置";
      if (this.panelType === "filter") return "选择筛选";
      if (this.panelType === "sort") return "选择排序";
      if (this.panelType === "more") return "更多条件";
      return "";
    },
    panelOptions() {
      return (this.filterOptions && this.filterOptions[this.panelType]) || [];
    },
  },
  methods: {
    async loadHomeData() {
      let res;
      try {
        res = await homeApi.Homedata({
          keyword: this.keyword || "",
          location: this.filters.location || "",
          filter: this.filters.filter || "",
          sort: this.filters.sort || "",
          more: Array.isArray(this.filters.more) ? this.filters.more : [],
        });
      } catch (e) {
        // auth: true 且本地无 token 时会被请求拦截器 reject
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return;
      }
      if (!res || res.code !== 0) return;
      const data = res.data || {};
      if (Array.isArray(data.bannerImages))
        this.bannerImages = data.bannerImages;
      if (typeof data.followedTotal === "number")
        this.followedTotal = data.followedTotal;
      if (Array.isArray(data.followedProperties))
        this.followedProperties = data.followedProperties;
      if (Array.isArray(data.recommendedProperties))
        this.recommendedProperties = data.recommendedProperties;
    },
    onSearchConfirm() {
      this.loadHomeData();
    },
    // 省份选择
    selectProvince() {
      uni.showActionSheet({
        itemList: this.provinceOptions,
        success: (res) => {
          const name = this.provinceOptions[res.tapIndex];
          if (!name) return;
          this.provinceName = name;
          uni.setStorageSync("hm_province_name", name);
          this.loadHomeData();
        },
      });
    },

    // 市/区：左侧二级联动
    openLocationPicker() {
      this.locationOpen = true;
      const [cityKey, districtKey] = String(this.filters.location || "").split(
        "|"
      );
      this.locDraft.city =
        cityKey || (this.locationCities[0] && this.locationCities[0].key) || "";
      this.locDraft.district = districtKey || "";

      if (!this.locDraft.district) {
        const city =
          this.locationCities.find((c) => c.key === this.locDraft.city) ||
          this.locationCities[0];
        const first = city && city.districts && city.districts[0];
        this.locDraft.district = (first && first.key) || "";
      }
    },
    closeLocationPicker() {
      this.locationOpen = false;
    },
    resetLocationPicker() {
      const city = this.locationCities[0];
      this.locDraft.city = (city && city.key) || "";
      this.locDraft.district =
        (city &&
          city.districts &&
          city.districts[0] &&
          city.districts[0].key) ||
        "";
    },
    selectLocCity(cityKey) {
      this.locDraft.city = cityKey;
      const city =
        this.locationCities.find((c) => c.key === cityKey) ||
        this.locationCities[0];
      const first = city && city.districts && city.districts[0];
      this.locDraft.district = (first && first.key) || "";
    },
    applyLocationPicker() {
      if (!this.locDraft.city) {
        this.filters.location = "";
      } else {
        this.filters.location = `${this.locDraft.city}|${
          this.locDraft.district || ""
        }`;
      }
      this.closeLocationPicker();
      this.loadHomeData();
    },

    openPanel(type) {
      this.panelType = type;
      this.panelOpen = true;
      this.panelDraft = {
        location: this.filters.location,
        filter: this.filters.filter,
        sort: this.filters.sort,
        more: Array.isArray(this.filters.more) ? [...this.filters.more] : [],
      };
    },
    closePanel() {
      this.panelOpen = false;
      this.panelType = "";
    },
    resetPanel() {
      if (this.panelType === "location") this.panelDraft.location = "";
      if (this.panelType === "filter") this.panelDraft.filter = "";
      if (this.panelType === "sort") this.panelDraft.sort = "";
      if (this.panelType === "more") this.panelDraft.more = [];
    },
    applyPanel() {
      this.filters = {
        location: this.panelDraft.location,
        filter: this.panelDraft.filter,
        sort: this.panelDraft.sort,
        more: Array.isArray(this.panelDraft.more)
          ? [...this.panelDraft.more]
          : [],
      };
      this.closePanel();
      this.loadHomeData();
    },
    isOptionActive(opt) {
      if (!opt) return false;
      if (this.panelType === "location")
        return this.panelDraft.location === opt.value;
      if (this.panelType === "filter")
        return this.panelDraft.filter === opt.value;
      if (this.panelType === "sort") return this.panelDraft.sort === opt.value;
      return false;
    },
    selectSingle(opt) {
      if (!opt) return;
      if (this.panelType === "location") this.panelDraft.location = opt.value;
      if (this.panelType === "filter") this.panelDraft.filter = opt.value;
      if (this.panelType === "sort") this.panelDraft.sort = opt.value;
    },
    isMoreSelected(opt) {
      if (!opt) return false;
      return (this.panelDraft.more || []).includes(opt.value);
    },
    toggleMore(opt) {
      if (!opt) return;
      const list = Array.isArray(this.panelDraft.more)
        ? this.panelDraft.more
        : [];
      const idx = list.indexOf(opt.value);
      if (idx >= 0) list.splice(idx, 1);
      else list.push(opt.value);
      this.panelDraft.more = list;
    },
    selectCity() {
      console.log("选择城市");
    },
    viewAllFollowed() {
      uni.navigateTo({
        url: "/pages/property_list/property_list",
      });
    },
    goToDetail(item) {
      const id = item && item.id ? item.id : "";
      uni.navigateTo({
        url: id
          ? `/pages/property_detail/property_detail?id=${id}`
          : "/pages/property_detail/property_detail",
      });
    },
    ensureTags(v) {
      if (!v) return [];
      if (Array.isArray(v)) return v.filter((x) => !!x);
      if (typeof v === "string")
        return v
          .split(",")
          .map((x) => x.trim())
          .filter((x) => !!x);
      return [];
    },
    getPriceUnit(item) {
      return (item && item.price_unit ? String(item.price_unit) : "万") || "万";
    },
    getSaleStatusLabel(v) {
      const map = { on_sale: "在售", sold: "已售", off_market: "下架" };
      return map[v] || "";
    },
    getSaleStatusClass(v) {
      const map = {
        on_sale: "badge-success",
        sold: "badge-default",
        off_market: "badge-warning",
      };
      return map[v] || "badge-default";
    },
    isSoldOrOff(v) {
      return v === "sold" || v === "off_market";
    },
    getLayoutText(item) {
      if (!item) return "";
      const rooms = Number(item.rooms || 0);
      const halls = Number(item.halls || 0);
      const bathrooms = Number(item.bathrooms || 0);
      if (rooms > 0 || halls > 0 || bathrooms > 0) {
        return `${rooms}室${halls}厅${bathrooms}卫`;
      }
      return "";
    },
    formatFollowedMeta(item) {
      const parts = [];
      const loc =
        item && (item.community_name || item.area)
          ? String(item.community_name || item.area)
          : "";
      if (loc) parts.push(loc);
      const layout = this.getLayoutText(item);
      if (layout) parts.push(layout);
      const size = item && item.size ? `${item.size}㎡` : "";
      if (size) parts.push(size);
      return parts.join(" · ");
    },
    formatRecommendedMeta(item) {
      const parts = [];
      const loc =
        item && (item.region || item.community_name)
          ? String(item.region || item.community_name)
          : "";
      if (loc) parts.push(loc);
      const layout = this.getLayoutText(item);
      if (layout) parts.push(layout);
      const size = item && item.size ? `${item.size}㎡` : "";
      if (size) parts.push(size);
      const type =
        item && (item.property_type || item.desc)
          ? String(item.property_type || item.desc)
          : "";
      if (type) parts.push(type);
      return parts.join(" · ");
    },
  },
};
</script>

<style lang="scss">
.home-container {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: #ffffff;
  /* 已改为 TopHeader 承担导航栏高度与左右间距 */
  padding: 0;
  z-index: 20;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.top-tools {
  padding: 20rpx 32rpx 0;
  background-color: #ffffff;
}

.top-tools__inner {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.city-selector {
  display: flex;
  align-items: center;
  gap: 4rpx;
  flex-shrink: 0;

  .city-name {
    font-size: 32rpx;
    font-weight: 800;
    color: #0d151c;
  }

  .expand-icon {
    font-size: 40rpx;
    color: #0d151c;
  }
}

.search-box {
  flex: 1;
  height: 80rpx;
  background-color: #f6f7f8;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
  border: 1px solid transparent;

  .search-icon {
    font-size: 40rpx;
    color: #94a3b8;
  }

  .search-input {
    flex: 1;
    margin-left: 16rpx;
    font-size: 28rpx;
    color: #0d151c;
  }
}

.main-content {
  flex: 1;
  overflow: hidden;
}

.banner-section {
  padding: 16rpx 32rpx 24rpx;
  background-color: #ffffff;
}

/* 市/区二级联动：底部半屏弹层（比左侧抽屉体验更好） */
.loc-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.45);
  z-index: 999;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.loc-drawer {
  width: 100%;
  height: 62vh;
  background: rgba(255, 255, 255, 0.96);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  box-shadow: 0 -10rpx 30rpx rgba(15, 23, 42, 0.12);
  display: flex;
  flex-direction: column;
  border-radius: 32rpx 32rpx 0 0;
}

.loc-head {
  padding: 24rpx 24rpx 16rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1rpx solid rgba(226, 232, 240, 0.8);
}

.loc-title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
}

.loc-reset {
  font-size: 26rpx;
  color: #64748b;
  font-weight: 700;
}

.loc-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.loc-col {
  height: 100%;
}

.loc-left {
  width: 44%;
  background-color: #f8fafc;
  border-right: 1rpx solid rgba(226, 232, 240, 0.8);
}

.loc-right {
  width: 56%;
  background-color: #ffffff;
}

.loc-item {
  padding: 24rpx 20rpx;
  font-size: 28rpx;
  color: #334155;
}

.loc-item.active {
  color: #2d9cf0;
  font-weight: 800;
  background-color: rgba(45, 156, 240, 0.08);
}

.loc-foot {
  padding: 16rpx 24rpx calc(env(safe-area-inset-bottom) + 24rpx);
  display: flex;
  gap: 16rpx;
  border-top: 1rpx solid rgba(226, 232, 240, 0.8);
}

/* 顶部右侧“更多”按钮 */
.more-btn {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  background-color: #f6f7f8;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  flex-shrink: 0;
  border: 1px solid rgba(226, 232, 240, 0.95);
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);

  .material-symbols-outlined {
    font-size: 40rpx;
    color: #0f172a;
  }

  &.active .material-symbols-outlined {
    color: #2d9cf0;
  }

  &:active {
    transform: scale(0.96);
  }
}

.more-dot {
  position: absolute;
  top: 10rpx;
  right: 10rpx;
  width: 14rpx;
  height: 14rpx;
  background-color: #f97316;
  border-radius: 50%;
  border: 3rpx solid #f6f7f8;
}

.loc-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 20rpx;
  font-size: 28rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  border: none;

  &::after {
    border: none;
  }

  &.ghost {
    background-color: #f1f5f9;
    color: #334155;
  }

  &.primary {
    background-color: #2d9cf0;
    color: #ffffff;
  }
}

.banner-card {
  width: 100%;
  height: 240rpx;
  border-radius: 24rpx;
  position: relative;
  overflow: hidden;
  box-shadow: 0 10rpx 20rpx rgba(37, 99, 235, 0.2);
}

.banner-swiper {
  width: 100%;
  height: 100%;
}

.banner-image {
  width: 100%;
  height: 100%;
  display: block;
}

.filter-bar {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  padding: 24rpx 32rpx;
  gap: 16rpx;
  border-bottom: 1rpx solid #f1f1f1;
  box-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.02);

  .filter-item {
    display: flex;
    align-items: center;
    gap: 8rpx;
    flex: 1;
    justify-content: center;
    padding: 12rpx 16rpx;
    background-color: #f6f7f8;
    border-radius: 16rpx;
    border: 1px solid transparent;

    .filter-text {
      font-size: 28rpx;
      font-weight: 500;
      color: #334155;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .filter-icon {
      font-size: 36rpx;
      color: #94a3b8;
    }

    &.active {
      background-color: rgba(45, 156, 240, 0.1);
      border-color: rgba(45, 156, 240, 0.2);

      .filter-text {
        color: #2d9cf0;
      }

      .filter-icon {
        color: #2d9cf0;
      }
    }
  }
}

/* 筛选弹层 */
.mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.45);
  z-index: 999;
  display: flex;
  align-items: flex-end;
}

.sheet {
  width: 100%;
  background-color: #ffffff;
  border-radius: 32rpx 32rpx 0 0;
  padding: 24rpx 24rpx calc(env(safe-area-inset-bottom) + 24rpx);
  box-sizing: border-box;
}

.sheet-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.sheet-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
}

.sheet-action {
  font-size: 26rpx;
  color: #64748b;
}

.sheet-list {
  max-height: 520rpx;
  overflow: hidden;
}

.sheet-item {
  height: 92rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8rpx;
  border-bottom: 1rpx solid #f1f5f9;
  color: #334155;
  box-sizing: border-box;

  &.active {
    color: #2d9cf0;
    font-weight: 700;
  }
}

.sheet-item-text {
  font-size: 28rpx;
}

.sheet-check {
  font-size: 36rpx;
}

.sheet-footer {
  display: flex;
  gap: 16rpx;
  margin-top: 16rpx;
}

.sheet-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 20rpx;
  font-size: 28rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  border: none;

  &::after {
    border: none;
  }

  &.ghost {
    background-color: #f1f5f9;
    color: #334155;
  }

  &.primary {
    background-color: #2d9cf0;
    color: #ffffff;
  }
}

.section {
  padding: 32rpx 0;
  display: flex;
  flex-direction: column;
  gap: 24rpx;

  .section-header {
    padding: 0 32rpx;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .section-title {
      display: flex;
      align-items: center;
      gap: 16rpx;

      .section-icon {
        font-size: 40rpx;

        &.accent {
          color: #f97316;
        }

        &.primary {
          color: #2d9cf0;
        }
      }

      .title-text {
        font-size: 36rpx;
        font-weight: bold;
        color: #0f172a;
      }
    }

    .section-more {
      display: flex;
      align-items: center;
      font-size: 24rpx;
      color: #64748b;

      .more-icon {
        font-size: 28rpx;
      }
    }
  }
}

.horizontal-scroll {
  white-space: nowrap;
  padding: 0 32rpx 8rpx;
  width: 100%;
  box-sizing: border-box;

  .property-card-h {
    display: inline-flex;
    flex-direction: column;
    width: 560rpx;
    background-color: #ffffff;
    border-radius: 24rpx;
    overflow: hidden;
    margin-right: 24rpx;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
    border: 1px solid #f1f1f1;

    .card-image {
      width: 100%;
      height: 240rpx;
      background-size: cover;
      background-position: center;
      position: relative;

      &.is-disabled {
        filter: grayscale(1) brightness(0.95);
      }

      &.is-disabled::after {
        content: "";
        position: absolute;
        inset: 0;
        background: rgba(15, 23, 42, 0.25);
      }

      .image-tag {
        position: absolute;
        top: 16rpx;
        left: 16rpx;
        background-color: rgba(15, 23, 42, 0.6);
        backdrop-filter: blur(4rpx);
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 600;
        padding: 4rpx 16rpx;
        border-radius: 8rpx;
        letter-spacing: 2rpx;
      }

      .price-change-tag {
        position: absolute;
        bottom: 16rpx;
        right: 16rpx;
        background-color: rgba(45, 156, 240, 0.9);
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 500;
        padding: 4rpx 16rpx;
        border-radius: 8rpx;
      }

      .status-tag {
        position: absolute;
        top: 16rpx;
        right: 16rpx;
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 800;
        padding: 4rpx 16rpx;
        border-radius: 8rpx;
        background-color: rgba(100, 116, 139, 0.95);

        &.badge-success {
          background-color: rgba(34, 197, 94, 0.95);
        }
        &.badge-warning {
          background-color: rgba(249, 115, 22, 0.95);
        }
        &.badge-default {
          background-color: rgba(100, 116, 139, 0.95);
        }
      }

      .lock-tag {
        position: absolute;
        bottom: 16rpx;
        left: 16rpx;
        background-color: rgba(45, 156, 240, 0.92);
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 700;
        padding: 4rpx 14rpx;
        border-radius: 8rpx;
      }
    }

    .card-info {
      padding: 24rpx;
      display: flex;
      flex-direction: column;
      gap: 8rpx;

      .property-name {
        font-size: 28rpx;
        font-weight: bold;
        color: #0f172a;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .mini-tags {
        display: flex;
        gap: 10rpx;
        flex-wrap: wrap;
        margin-top: 2rpx;

        .mini-tag {
          padding: 2rpx 10rpx;
          background-color: rgba(241, 245, 249, 0.95);
          border-radius: 8rpx;
          font-size: 22rpx;
          color: #475569;
          line-height: 1.4;
        }
      }

      .property-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 24rpx;
        color: #64748b;

        .price {
          display: flex;
          align-items: baseline;

          .price-val {
            font-size: 28rpx;
            font-weight: bold;
            color: #2d9cf0;
          }

          .price-unit {
            font-size: 20rpx;
            font-weight: normal;
            color: #64748b;
            margin-left: 2rpx;
          }
        }
      }
    }
  }
}

.followed-empty {
  padding: 0 32rpx 8rpx;

  .empty-card {
    width: 100%;
    background-color: #ffffff;
    border-radius: 24rpx;
    padding: 36rpx 28rpx;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12rpx;
    border: 1px solid #f1f1f1;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
  }

  .empty-icon {
    font-size: 64rpx;
    color: rgba(249, 115, 22, 0.9);
  }

  .empty-title {
    font-size: 30rpx;
    font-weight: 800;
    color: #0f172a;
  }

  .empty-desc {
    font-size: 24rpx;
    color: #64748b;
    text-align: center;
  }

  .empty-actions {
    margin-top: 10rpx;
    width: 100%;
    display: flex;
    justify-content: center;
  }

  .empty-btn {
    height: 72rpx;
    padding: 0 28rpx;
    border-radius: 18rpx;
    font-size: 26rpx;
    font-weight: 800;
    display: flex;
    align-items: center;
    justify-content: center;
    line-height: 1;
    border: none;

    &::after {
      border: none;
    }

    &.ghost {
      background-color: #f1f5f9;
      color: #334155;
    }
  }
}

.vertical-list {
  padding: 0 32rpx;
  display: flex;
  flex-direction: column;
  gap: 32rpx;

  .property-card-v {
    background-color: #ffffff;
    border-radius: 24rpx;
    overflow: hidden;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
    border: 1px solid #f1f1f1;
    display: flex;
    flex-direction: column;

    .card-image-v {
      width: 100%;
      aspect-ratio: 16 / 9;
      background-size: cover;
      background-position: center;
      position: relative;

      &.is-disabled {
        filter: grayscale(1) brightness(0.95);
      }

      &.is-disabled::after {
        content: "";
        position: absolute;
        inset: 0;
        background: rgba(15, 23, 42, 0.25);
      }

      .status-tag-v {
        position: absolute;
        top: 16rpx;
        left: 16rpx;
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 800;
        padding: 4rpx 16rpx;
        border-radius: 8rpx;
        background-color: rgba(100, 116, 139, 0.95);

        &.badge-success {
          background-color: rgba(34, 197, 94, 0.95);
        }
        &.badge-warning {
          background-color: rgba(249, 115, 22, 0.95);
        }
        &.badge-default {
          background-color: rgba(100, 116, 139, 0.95);
        }
      }

      .lock-tag-v {
        position: absolute;
        top: 16rpx;
        right: 16rpx;
        background-color: rgba(45, 156, 240, 0.92);
        color: #ffffff;
        font-size: 20rpx;
        font-weight: 700;
        padding: 4rpx 14rpx;
        border-radius: 8rpx;
      }
    }

    .card-info-v {
      padding: 32rpx;
      display: flex;
      flex-direction: column;
      gap: 16rpx;

      .info-top {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;

        .property-name-v {
          font-size: 32rpx;
          font-weight: bold;
          color: #0f172a;
        }

        .price-v {
          display: flex;
          align-items: baseline;
          color: #2d9cf0;

          .price-val {
            font-size: 32rpx;
            font-weight: bold;
          }

          .price-unit {
            font-size: 24rpx;
            margin-left: 4rpx;
          }
        }
      }

      .property-location {
        display: flex;
        align-items: center;
        gap: 8rpx;
        font-size: 24rpx;
        color: #64748b;

        .loc-icon {
          font-size: 28rpx;
        }
      }

      .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 16rpx;
        margin-top: 8rpx;

        .tag-item {
          padding: 4rpx 16rpx;
          background-color: #f1f5f9;
          border-radius: 8rpx;
          font-size: 24rpx;
          font-weight: 500;
          color: #475569;

          &.tag-orange {
            background-color: #fff7ed;
            color: #f97316;
          }
        }
      }
    }
  }
}

.tab-bar {
  background-color: #ffffff;
  border-top: 1rpx solid #e2e8f0;
  padding: 24rpx 48rpx 64rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 30;

  .tab-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8rpx;
    color: #94a3b8;

    .tab-icon {
      font-size: 52rpx;
      transition: transform 0.2s;

      &.fill {
        font-variation-settings: "FILL" 1;
      }
    }

    .tab-text {
      font-size: 20rpx;
      font-weight: 500;
    }

    &.active {
      color: #2d9cf0;
    }

    &:active .tab-icon {
      transform: scale(1.1);
    }
  }

  .icon-with-dot {
    position: relative;

    .dot {
      position: absolute;
      top: -2rpx;
      right: -2rpx;
      width: 16rpx;
      height: 16rpx;
      background-color: #f97316;
      border-radius: 50%;
      border: 2rpx solid #ffffff;
    }
  }
}

.bottom-spacer {
  height: 120rpx;
}

.placeholder {
  color: #94a3b8;
}
</style>
