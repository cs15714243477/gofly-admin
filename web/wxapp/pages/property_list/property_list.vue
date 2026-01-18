<template>
  <view class="list-container">
    <!-- 顶部导航与搜索 -->
    <view class="header">
      <TopHeader title="房源" />

      <view class="search-row">
        <picker
          v-if="provincePickerLabels.length > 0"
          mode="selector"
          :range="provincePickerLabels"
          :value="provincePickerIndex"
          @change="onProvincePickerChange"
        >
          <view class="province-picker">
            <text class="province-text">{{ provinceName }}</text>
            <text class="material-symbols-outlined province-icon"
              >expand_more</text
            >
          </view>
        </picker>
        <view v-else class="province-picker" @tap="reloadProvinces">
          <text class="province-text">{{ provinceName }}</text>
          <text class="material-symbols-outlined province-icon"
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
            placeholder="搜索小区、商圈、地址..."
            placeholder-class="placeholder"
            @input="onKeywordInput"
            @confirm="onSearchConfirm"
          />
        </view>
        <view
          class="filter-icon-btn"
          :class="{ active: hasActiveFilters }"
          @tap="openFilter"
        >
          <text class="material-symbols-outlined">tune</text>
          <view v-if="hasActiveFilters" class="filter-dot"></view>
        </view>
      </view>

      <scroll-view scroll-x="true" class="category-scroll">
        <view class="category-row">
          <view
            class="category-item"
            v-for="(c, idx) in navTabs"
            :key="c.key || idx"
            :class="{ active: isNavActive(c) }"
            @tap="onNavTab(c)"
          >
            {{ c.label }}
          </view>
        </view>
      </scroll-view>
    </view>

    <scroll-view
      scroll-y="true"
      class="main-list"
      lower-threshold="120"
      @scrolltolower="loadMore"
    >
      <view class="list-content">
        <view
          class="property-card"
          v-for="(item, index) in displayList"
          :key="item.id || index"
          :class="{ disabled: item.sale_status === 'off_market' }"
          @click="goToDetail(item)"
        >
          <view class="card-main">
            <view
              class="image-box"
              :class="{ grayscale: isSoldOrOff(item.sale_status) }"
            >
              <image
                class="property-image"
                :src="item.image"
                mode="aspectFill"
              ></image>
              <view class="image-tag" v-if="Number(item.has_smart_lock) === 1">
                <text class="material-symbols-outlined tag-icon">lock</text>
                <text>智能门锁</text>
              </view>
              <view class="mask" v-if="isSoldOrOff(item.sale_status)">
                <text class="mask-text">{{
                  item.sale_status_label || ""
                }}</text>
              </view>
            </view>
            <view class="info-box">
              <view class="title">{{ item.title }}</view>
              <view class="meta">
                <text class="bold">{{ getLayoutText(item) }}</text>
                <text class="divider">|</text>
                <text>{{ item.area || "-" }}㎡</text>
                <text class="divider">|</text>
                <text>{{ item.orientation || "-" }}</text>
              </view>
              <view class="region-row">
                <text class="region-label">市区</text>
                <text class="region-val">{{ getDistrictText(item) }}</text>
              </view>
              <view class="tags">
                <text
                  class="tag"
                  v-for="(tag, tIdx) in ensureTags(item.tags).slice(0, 4)"
                  :key="tIdx"
                  :class="tagColorClass(tag)"
                  >{{ tag }}</text
                >
              </view>
              <view class="price-row">
                <text class="price"
                  >¥{{ item.price }}{{ item.price_unit || "" }}</text
                >
                <text class="unit-price">{{ getUnitPriceText(item) }}</text>
              </view>
              <view class="stats">
                <text>浏览: {{ item.view_count || 0 }}</text>
                <text class="stats-sep"></text>
                <text>关注: {{ item.follow_count || 0 }}</text>
              </view>
            </view>
          </view>
          <view
            class="card-footer"
            :class="
              isSoldOrOff(item.sale_status) ? 'grey-footer' : 'orange-footer'
            "
          >
            <view class="footer-left">
              <view class="footer-icon-box">
                <text class="material-symbols-outlined footer-icon"
                  >currency_yen</text
                >
              </view>
              <text class="footer-text">{{ getCommissionText(item) }}</text>
            </view>
            <button
              class="footer-btn"
              v-if="item.sale_status !== 'off_market'"
              @click.stop="handlePromote(item)"
            >
              <text class="material-symbols-outlined btn-icon">share</text>
              <text>微信推广</text>
            </button>
            <button class="footer-btn-disabled" v-else>已下架</button>
          </view>
        </view>

        <view v-if="!loading && displayList.length === 0" class="empty-wrap">
          <text class="material-symbols-outlined empty-icon">apartment</text>
          <view class="empty-title">暂无房源</view>
          <view class="empty-desc">调整筛选条件后再试试</view>
        </view>

        <view v-if="loading" class="loading-row">
          <text>加载中...</text>
        </view>
        <view
          v-else-if="finished && displayList.length > 0"
          class="loading-row muted"
        >
          <text>没有更多了</text>
        </view>
      </view>
      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部导航 -->
    <BottomTabBar active="property" />

    <!-- 筛选弹层（静态数据 + 本地过滤/排序） -->
    <view v-if="filterOpen" class="filter-mask" @tap="closeFilter">
      <view class="filter-sheet" @tap.stop>
        <view class="filter-header">
          <text class="filter-title">筛选</text>
          <text class="filter-reset" @tap="resetFilter">重置</text>
        </view>

        <scroll-view scroll-y class="filter-body">
          <view class="filter-section">
            <text class="filter-section-title">分类</text>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in categoryOptions"
                :key="'cat-' + idx"
                class="chip"
                :class="{ active: draft.category === opt.value }"
                @tap="onDraftCategoryChange(opt.value)"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">户型</text>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in roomOptions"
                :key="'room-' + idx"
                class="chip"
                :class="{ active: draft.rooms === opt.value }"
                @tap="draft.rooms = opt.value"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">特色</text>
            <view class="chip-row">
              <view
                class="chip"
                :class="{ active: !!draft.smartLock }"
                @tap="draft.smartLock = !draft.smartLock"
              >
                <text>智能门锁</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">状态</text>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in statusOptions"
                :key="'status-' + idx"
                class="chip"
                :class="{ active: draft.sale_status === opt.value }"
                @tap="draft.sale_status = opt.value"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">推荐</text>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in hotStatusOptions"
                :key="'hot-' + idx"
                class="chip"
                :class="{ active: draft.hot_status === opt.value }"
                @tap="draft.hot_status = opt.value"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">物业类型</text>
            <view v-if="filterOptionsLoading" class="muted-tip">加载中...</view>
            <view v-else class="chip-row">
              <view
                v-for="(opt, idx) in propertyTypeOptions"
                :key="'ptype-' + idx"
                class="chip"
                :class="{ active: draft.property_type === opt }"
                @tap="draft.property_type = opt"
              >
                <text>{{ opt || "不限" }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">装修类型</text>
            <view v-if="filterOptionsLoading" class="muted-tip">加载中...</view>
            <view v-else class="chip-row">
              <view
                v-for="(opt, idx) in decorationTypeOptions"
                :key="'dtype-' + idx"
                class="chip"
                :class="{ active: draft.decoration_type === opt }"
                @tap="draft.decoration_type = opt"
              >
                <text>{{ opt || "不限" }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">朝向</text>
            <view v-if="filterOptionsLoading" class="muted-tip">加载中...</view>
            <view v-else class="chip-row">
              <view
                v-for="(opt, idx) in orientationOptions"
                :key="'ori-' + idx"
                class="chip"
                :class="{ active: draft.orientation === opt }"
                @tap="draft.orientation = opt"
              >
                <text>{{ opt || "不限" }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">标签</text>
            <view v-if="filterOptionsLoading" class="muted-tip">加载中...</view>
            <view v-else class="chip-row">
              <view
                v-for="(tag, idx) in tagOptions"
                :key="'tag-' + idx"
                class="chip"
                :class="{ active: draft.more.includes(tag) }"
                @tap="toggleTag(tag)"
              >
                <text>{{ tag }}</text>
              </view>
              <view v-if="tagOptions.length === 0" class="muted-tip">
                暂无标签
              </view>
            </view>
            <view v-if="draft.more.length > 0" class="muted-tip">
              已选：{{ draft.more.join("、") }}
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">建成年份</text>
            <view v-if="buildYearStatText" class="muted-tip">
              数据范围：{{ buildYearStatText }}
            </view>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in buildYearRanges"
                :key="'year-' + idx"
                class="chip"
                :class="{ active: isRangeActive('year', opt) }"
                @tap="selectRange('year', opt)"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">价格区间（万）</text>
            <view v-if="priceStatText" class="muted-tip">
              数据范围：{{ priceStatText }}
            </view>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in priceRanges"
                :key="'price-' + idx"
                class="chip"
                :class="{ active: isRangeActive('price', opt) }"
                @tap="selectRange('price', opt)"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">面积区间（㎡）</text>
            <view v-if="areaStatText" class="muted-tip">
              数据范围：{{ areaStatText }}
            </view>
            <view class="chip-row">
              <view
                v-for="(opt, idx) in areaRanges"
                :key="'area-' + idx"
                class="chip"
                :class="{ active: isRangeActive('area', opt) }"
                @tap="selectRange('area', opt)"
              >
                <text>{{ opt.label }}</text>
              </view>
            </view>
          </view>

          <view class="filter-section">
            <text class="filter-section-title">排序</text>
            <view class="radio-list">
              <view
                v-for="(opt, idx) in sortOptions"
                :key="'sort-' + idx"
                class="radio-item"
                :class="{ active: draft.sort === opt.value }"
                @tap="draft.sort = opt.value"
              >
                <text>{{ opt.label }}</text>
                <text
                  v-if="draft.sort === opt.value"
                  class="material-symbols-outlined radio-check"
                  >check</text
                >
              </view>
            </view>
          </view>
        </scroll-view>

        <view class="filter-footer">
          <button class="filter-btn ghost" @tap="closeFilter">取消</button>
          <button class="filter-btn primary" @tap="applyFilter">确定</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import BottomTabBar from "@/components/BottomTabBar.vue";
import TopHeader from "@/components/TopHeader.vue";
import areaApi from "@/api/area";
import propertyApi from "@/api/property";

export default {
  components: { BottomTabBar, TopHeader },
  onLoad() {
    this.initPage();
  },
  onUnload() {
    if (this.searchTimer) clearTimeout(this.searchTimer);
  },
  data() {
    return {
      provinceName: "辽宁",
      provinceId: 0,
      provinces: [],
      cities: [],
      districts: [],
      currentCityId: 0,
      currentCityName: "",
      currentDistrictId: 0,
      currentDistrictName: "",
      keyword: "",
      loading: false,
      finished: false,
      pagination: { page: 1, pageSize: 10, total: 0 },
      propertyList: [],
      searchTimer: null,
      filterOptions: {
        property_types: [],
        decoration_types: [],
        orientations: [],
        tags: [],
        stats: {},
      },
      filterOptionsLoading: false,
      filterOptionsReqId: 0,
      // 筛选（静态数据）
      filters: {
        category: "all",
        rooms: "",
        smartLock: false,
        sale_status: "",
        sort: "",
        hot_status: "",
        property_type: "",
        decoration_type: "",
        orientation: "",
        build_year_min: "",
        build_year_max: "",
        price_min: "",
        price_max: "",
        area_min: "",
        area_max: "",
        more: [],
      },
      filterOpen: false,
      draft: {
        category: "all",
        rooms: "",
        smartLock: false,
        sale_status: "",
        sort: "",
        hot_status: "",
        property_type: "",
        decoration_type: "",
        orientation: "",
        build_year_min: "",
        build_year_max: "",
        price_min: "",
        price_max: "",
        area_min: "",
        area_max: "",
        more: [],
      },
      roomOptions: [
        { value: "", label: "不限" },
        { value: "1", label: "1室" },
        { value: "2", label: "2室" },
        { value: "3", label: "3室" },
        { value: "4+", label: "4室+" },
      ],
      categoryOptions: [
        { value: "all", label: "全部房源" },
        { value: "mine", label: "我的房源" },
      ],
      statusOptions: [
        { value: "", label: "不限" },
        { value: "on_sale", label: "在售" },
        { value: "sold", label: "已售" },
        { value: "off_market", label: "下架" },
      ],
      hotStatusOptions: [
        { value: "", label: "不限" },
        { value: 1, label: "推荐" },
        { value: 0, label: "不推荐" },
      ],
      sortOptions: [
        { value: "", label: "默认" },
        { value: "view_desc", label: "浏览最多" },
        { value: "follow_desc", label: "关注最多" },
        { value: "latest", label: "最新发布" },
        { value: "price_asc", label: "价格从低到高" },
        { value: "price_desc", label: "价格从高到低" },
      ],
      mockPropertyList: [],
    };
  },
  computed: {
    hasActiveFilters() {
      return !!(
        this.keyword ||
        this.filters.category !== "all" ||
        this.filters.rooms ||
        this.filters.smartLock ||
        this.filters.sale_status ||
        this.filters.sort ||
        this.filters.hot_status !== "" ||
        this.filters.property_type ||
        this.filters.decoration_type ||
        this.filters.orientation ||
        this.filters.build_year_min ||
        this.filters.build_year_max ||
        this.filters.price_min ||
        this.filters.price_max ||
        this.filters.area_min ||
        this.filters.area_max ||
        (Array.isArray(this.filters.more) && this.filters.more.length > 0)
      );
    },
    propertyTypeOptions() {
      const list = Array.isArray(this.filterOptions.property_types)
        ? this.filterOptions.property_types
        : [];
      const dedup = Array.from(
        new Set(list.map((s) => String(s || "").trim()).filter(Boolean))
      );
      return [""].concat(dedup);
    },
    decorationTypeOptions() {
      const list = Array.isArray(this.filterOptions.decoration_types)
        ? this.filterOptions.decoration_types
        : [];
      const dedup = Array.from(
        new Set(list.map((s) => String(s || "").trim()).filter(Boolean))
      );
      return [""].concat(dedup);
    },
    orientationOptions() {
      const list = Array.isArray(this.filterOptions.orientations)
        ? this.filterOptions.orientations
        : [];
      const dedup = Array.from(
        new Set(list.map((s) => String(s || "").trim()).filter(Boolean))
      );
      return [""].concat(dedup);
    },
    tagOptions() {
      const list = Array.isArray(this.filterOptions.tags)
        ? this.filterOptions.tags
        : [];
      return Array.from(
        new Set(list.map((s) => String(s || "").trim()).filter(Boolean))
      );
    },
    buildYearStatText() {
      const min = this.statNumber("build_year_min");
      const max = this.statNumber("build_year_max");
      if (!isFinite(min) || !isFinite(max) || min <= 0 || max <= 0 || max < min)
        return "";
      return `${Math.floor(min)}-${Math.floor(max)}`;
    },
    priceStatText() {
      const min = this.statNumber("price_min");
      const max = this.statNumber("price_max");
      if (!isFinite(min) || !isFinite(max) || max < min) return "";
      return `${this.formatNumber(min)}-${this.formatNumber(max)}万`;
    },
    areaStatText() {
      const min = this.statNumber("area_min");
      const max = this.statNumber("area_max");
      if (!isFinite(min) || !isFinite(max) || max < min) return "";
      return `${this.formatNumber(min)}-${this.formatNumber(max)}㎡`;
    },
    buildYearRanges() {
      const base = [{ label: "不限", min: "", max: "" }];
      const min = this.statNumber("build_year_min");
      const max = this.statNumber("build_year_max");
      if (!isFinite(min) || !isFinite(max) || min <= 0 || max <= 0 || max < min)
        return base;
      return base.concat(
        this.makeYearRanges(Math.floor(min), Math.floor(max), 5)
      );
    },
    priceRanges() {
      const base = [{ label: "不限", min: "", max: "" }];
      const min = this.statNumber("price_min");
      const max = this.statNumber("price_max");
      if (!isFinite(min) || !isFinite(max) || max < min) return base;
      return base.concat(
        this.makeValueRanges(
          Math.max(0, Math.floor(min)),
          Math.max(0, Math.ceil(max)),
          5,
          [10, 20, 30, 50, 80, 100, 150, 200, 300, 500, 800, 1000],
          "万"
        )
      );
    },
    areaRanges() {
      const base = [{ label: "不限", min: "", max: "" }];
      const min = this.statNumber("area_min");
      const max = this.statNumber("area_max");
      if (!isFinite(min) || !isFinite(max) || max < min) return base;
      return base.concat(
        this.makeValueRanges(
          Math.max(0, Math.floor(min)),
          Math.max(0, Math.ceil(max)),
          5,
          [10, 20, 30, 50, 80, 100, 120, 150, 200, 300],
          "㎡"
        )
      );
    },
    provincePickerLabels() {
      const list = Array.isArray(this.provinces) ? this.provinces : [];
      return list
        .map((p) => this.formatProvinceLabel(p && p.name))
        .filter((x) => !!x);
    },
    provincePickerIndex() {
      const list = Array.isArray(this.provinces) ? this.provinces : [];
      if (list.length === 0) return 0;
      const id = Number(this.provinceId) || 0;
      if (id) {
        const idx = list.findIndex((p) => Number(p && p.id) === id);
        if (idx >= 0) return idx;
      }
      const name = String(this.provinceName || "").trim();
      if (name) {
        const idx = list.findIndex(
          (p) => this.formatProvinceLabel(p && p.name) === name
        );
        if (idx >= 0) return idx;
      }
      return 0;
    },
    navTabs() {
      if (!this.currentCityId) {
        const list = Array.isArray(this.cities) ? this.cities : [];
        return [{ key: "", label: "全部" }].concat(
          list.map((it) => ({
            key: `city:${it.id}`,
            id: it.id,
            name: it.name,
            label: this.formatCityLabel(it.name),
          }))
        );
      }
      const list = Array.isArray(this.districts) ? this.districts : [];
      return [
        { key: "back", label: "返回" },
        { key: "all_city", label: "全市" },
      ].concat(
        list.map((it) => ({
          key: `district:${it.id}`,
          id: it.id,
          name: it.name,
          label: it.name,
        }))
      );
    },
    displayList() {
      return Array.isArray(this.propertyList) ? this.propertyList : [];
    },
  },
  methods: {
    statNumber(key) {
      const stats = (this.filterOptions && this.filterOptions.stats) || {};
      const v = stats ? stats[key] : undefined;
      const n = Number(v);
      return n;
    },
    formatNumber(v) {
      const n = Number(v);
      if (!isFinite(n)) return "-";
      if (Math.abs(n - Math.round(n)) < 1e-9) return String(Math.round(n));
      return String(Number(n.toFixed(2)));
    },
    makeNiceStep(span, candidates, targetBuckets) {
      const s = Number(span);
      if (!isFinite(s) || s <= 0) return candidates[0] || 1;
      const raw = s / Math.max(1, targetBuckets);
      const list = Array.isArray(candidates) ? candidates : [];
      for (let i = 0; i < list.length; i++) {
        const step = Number(list[i]);
        if (isFinite(step) && step > 0 && step >= raw) return step;
      }
      const last = Number(list[list.length - 1]);
      if (isFinite(last) && last > 0) return last;
      return Math.ceil(raw) || 1;
    },
    makeYearRanges(minYear, maxYear, maxItems) {
      const min = Number(minYear);
      const max = Number(maxYear);
      if (!isFinite(min) || !isFinite(max) || max < min) return [];
      const span = max - min + 1;
      const step = this.makeNiceStep(span, [5, 10, 20, 25, 50], maxItems);
      const start = Math.floor(min / step) * step;
      const out = [];
      for (let s = start; s <= max && out.length < maxItems; s += step) {
        const e = s + step - 1;
        if (e >= max) {
          out.push({ label: `${s}及以后`, min: String(s), max: "" });
          break;
        }
        out.push({ label: `${s}-${e}`, min: String(s), max: String(e) });
      }
      return out;
    },
    makeValueRanges(minValue, maxValue, maxItems, stepCandidates, unit) {
      const min = Number(minValue);
      const max = Number(maxValue);
      if (!isFinite(min) || !isFinite(max) || max < min) return [];
      const span = max - min;
      const step = this.makeNiceStep(span, stepCandidates, maxItems);
      const start = Math.max(0, Math.floor(min / step) * step);
      const out = [];
      for (let s = start; s <= max && out.length < maxItems; s += step) {
        const e = s + step;
        if (e >= max) {
          out.push({ label: `${s}${unit}以上`, min: String(s), max: "" });
          break;
        }
        out.push({
          label: `${s}-${e}${unit}`,
          min: String(s),
          max: String(e),
        });
      }
      return out;
    },
    isRangeActive(kind, opt) {
      const o = opt || {};
      if (kind === "year") {
        return (
          String(this.draft.build_year_min || "") === String(o.min || "") &&
          String(this.draft.build_year_max || "") === String(o.max || "")
        );
      }
      if (kind === "price") {
        return (
          String(this.draft.price_min || "") === String(o.min || "") &&
          String(this.draft.price_max || "") === String(o.max || "")
        );
      }
      if (kind === "area") {
        return (
          String(this.draft.area_min || "") === String(o.min || "") &&
          String(this.draft.area_max || "") === String(o.max || "")
        );
      }
      return false;
    },
    selectRange(kind, opt) {
      const o = opt || {};
      if (kind === "year") {
        this.draft.build_year_min = String(o.min || "");
        this.draft.build_year_max = String(o.max || "");
        return;
      }
      if (kind === "price") {
        this.draft.price_min = String(o.min || "");
        this.draft.price_max = String(o.max || "");
        return;
      }
      if (kind === "area") {
        this.draft.area_min = String(o.min || "");
        this.draft.area_max = String(o.max || "");
      }
    },
    async initPage() {
      const cached = uni.getStorageSync("hm_province_name");
      this.provinceName = cached || "辽宁";
      await this.loadProvinces();
      this.applyProvinceByName(this.provinceName);
      await this.loadCities();
      this.loadFilterOptions({ category: this.filters.category });
      this.fetchList(true);
    },
    async reloadProvinces() {
      uni.showLoading({ title: "加载中", mask: true });
      await this.loadProvinces();
      uni.hideLoading();
      if (!Array.isArray(this.provinces) || this.provinces.length === 0) {
        uni.showToast({ title: "省份数据获取失败", icon: "none" });
        return;
      }
      this.applyProvinceByName(this.provinceName || "辽宁");
      await this.loadCities();
      this.loadFilterOptions({ category: this.filters.category });
      this.fetchList(true);
    },
    async onProvincePickerChange(e) {
      const idx = Number(e && e.detail && e.detail.value);
      const list = Array.isArray(this.provinces) ? this.provinces : [];
      const row = list[idx];
      if (!row) return;
      const nextName = this.formatProvinceLabel(row.name);
      if (!nextName) return;
      if (nextName === this.provinceName) return;
      this.provinceName = nextName;
      this.provinceId = Number(row.id) || 0;
      uni.setStorageSync("hm_province_name", nextName);
      this.currentCityId = 0;
      this.currentCityName = "";
      this.currentDistrictId = 0;
      this.currentDistrictName = "";
      this.districts = [];
      await this.loadCities();
      this.loadFilterOptions({ category: this.filters.category });
      this.fetchList(true);
    },
    onKeywordInput() {
      if (this.searchTimer) clearTimeout(this.searchTimer);
      this.searchTimer = setTimeout(() => {
        this.fetchList(true);
      }, 350);
    },
    formatProvinceLabel(name) {
      const s = String(name || "").trim();
      if (!s) return "";
      if (s.endsWith("省")) return s.slice(0, -1);
      if (s.endsWith("市")) return s.slice(0, -1);
      if (s.endsWith("自治区")) return s.slice(0, -3);
      return s;
    },
    formatCityLabel(name) {
      const s = String(name || "").trim();
      if (!s) return "";
      return s.endsWith("市") ? s.slice(0, -1) : s;
    },
    async loadProvinces() {
      let res;
      try {
        res = await areaApi.getProvinces({});
      } catch (e) {
        return false;
      }
      if (!res || res.code !== 0) return false;
      const items = (res.data && res.data.items) || [];
      this.provinces = Array.isArray(items) ? items : [];
      return this.provinces.length > 0;
    },
    applyProvinceByName(name) {
      const n = String(name || "").trim();
      if (!n || !Array.isArray(this.provinces) || this.provinces.length === 0) {
        this.provinceId = 0;
        this.provinceName = n || "辽宁";
        return;
      }
      const hit =
        this.provinces.find((p) => String(p && p.name).includes(n)) ||
        this.provinces.find((p) => n.includes(String(p && p.name))) ||
        this.provinces[0];
      this.provinceId = Number(hit && hit.id) || 0;
      this.provinceName =
        this.formatProvinceLabel(hit && hit.name) || n || "辽宁";
      uni.setStorageSync("hm_province_name", this.provinceName);
      this.currentCityId = 0;
      this.currentCityName = "";
      this.currentDistrictId = 0;
      this.currentDistrictName = "";
      this.districts = [];
    },
    async loadCities() {
      if (!this.provinceId) {
        this.cities = [];
        return false;
      }
      let res;
      try {
        res = await areaApi.getChildren({ pid: this.provinceId });
      } catch (e) {
        return false;
      }
      if (!res || res.code !== 0) return false;
      const items = (res.data && res.data.items) || [];
      this.cities = Array.isArray(items) ? items : [];
      return true;
    },
    async loadDistricts(cityId) {
      const id = Number(cityId) || 0;
      if (!id) {
        this.districts = [];
        return false;
      }
      let res;
      try {
        res = await areaApi.getChildren({ pid: id });
      } catch (e) {
        return false;
      }
      if (!res || res.code !== 0) return false;
      const items = (res.data && res.data.items) || [];
      this.districts = Array.isArray(items) ? items : [];
      return true;
    },
    async loadFilterOptions(override) {
      const reqId = (this.filterOptionsReqId || 0) + 1;
      this.filterOptionsReqId = reqId;
      this.filterOptionsLoading = true;
      const category = String(
        (override && override.category) ||
          (this.filterOpen ? this.draft.category : this.filters.category) ||
          "all"
      ).trim();
      const params = {
        province: this.provinceName || "辽宁",
        city: this.currentCityName || "",
        district: this.currentDistrictName || "",
        category,
      };
      let res;
      try {
        res = await propertyApi.getFilterOptions(params);
      } catch (e) {
        if (reqId === this.filterOptionsReqId)
          this.filterOptionsLoading = false;
        return false;
      }
      if (reqId === this.filterOptionsReqId) this.filterOptionsLoading = false;
      if (!res || res.code !== 0) return false;
      if (reqId !== this.filterOptionsReqId) return false;
      const data = res.data || {};
      this.filterOptions = {
        property_types: Array.isArray(data.property_types)
          ? data.property_types
          : [],
        decoration_types: Array.isArray(data.decoration_types)
          ? data.decoration_types
          : [],
        orientations: Array.isArray(data.orientations) ? data.orientations : [],
        tags: Array.isArray(data.tags) ? data.tags : [],
        stats: data.stats || {},
      };
      return true;
    },
    onDraftCategoryChange(next) {
      this.draft.category = next;
      this.loadFilterOptions({ category: next });
    },
    toggleTag(tag) {
      const t = String(tag || "").trim();
      if (!t) return;
      const list = Array.isArray(this.draft.more) ? this.draft.more : [];
      const idx = list.indexOf(t);
      if (idx >= 0) {
        list.splice(idx, 1);
        this.draft.more = list;
        return;
      }
      if (list.length >= 6) {
        uni.showToast({ title: "最多选择6个标签", icon: "none" });
        return;
      }
      list.push(t);
      this.draft.more = list;
    },
    isNavActive(tab) {
      if (!tab) return false;
      if (tab.key === "" && !this.currentCityId) return true;
      if (
        tab.key === "all_city" &&
        this.currentCityId &&
        !this.currentDistrictId
      )
        return true;
      if (tab.key && String(tab.key).startsWith("city:")) {
        return (
          !!this.currentCityId && Number(tab.id) === Number(this.currentCityId)
        );
      }
      if (tab.key && String(tab.key).startsWith("district:")) {
        return (
          !!this.currentCityId &&
          Number(tab.id) === Number(this.currentDistrictId)
        );
      }
      return false;
    },
    async onNavTab(tab) {
      if (!tab) return;
      const key = String(tab.key || "");
      if (key === "") {
        this.currentCityId = 0;
        this.currentCityName = "";
        this.currentDistrictId = 0;
        this.currentDistrictName = "";
        this.districts = [];
        this.loadFilterOptions({ category: this.filters.category });
        this.fetchList(true);
        return;
      }
      if (key === "back") {
        this.currentCityId = 0;
        this.currentCityName = "";
        this.currentDistrictId = 0;
        this.currentDistrictName = "";
        this.districts = [];
        this.loadFilterOptions({ category: this.filters.category });
        this.fetchList(true);
        return;
      }
      if (key === "all_city") {
        this.currentDistrictId = 0;
        this.currentDistrictName = "";
        this.loadFilterOptions({ category: this.filters.category });
        this.fetchList(true);
        return;
      }
      if (key.startsWith("city:")) {
        this.currentCityId = Number(tab.id) || 0;
        this.currentCityName = String(tab.name || "").trim();
        this.currentDistrictId = 0;
        this.currentDistrictName = "";
        await this.loadDistricts(this.currentCityId);
        this.loadFilterOptions({ category: this.filters.category });
        this.fetchList(true);
        return;
      }
      if (key.startsWith("district:")) {
        this.currentDistrictId = Number(tab.id) || 0;
        this.currentDistrictName = String(tab.name || "").trim();
        this.loadFilterOptions({ category: this.filters.category });
        this.fetchList(true);
        return;
      }
    },
    handlePromote(item) {
      const title = (item && item.title) || "房源";
      uni.showActionSheet({
        itemList: ["生成获客海报", "复制推广文案", "模拟分享（占位）"],
        success: (res) => {
          if (res.tapIndex === 0) {
            uni.navigateTo({
              url: `/pages/marketing_posters/marketing_posters?title=${encodeURIComponent(
                title
              )}`,
            });
            return;
          }
          if (res.tapIndex === 1) {
            const id = item && (item.id || item.ID);
            const path = id
              ? `/pages/property_detail/property_detail?id=${encodeURIComponent(
                  id
                )}`
              : "/pages/property_detail/property_detail";
            const text = `【优质房源推荐】${title}\n点击查看详情：${path}`;
            uni.setClipboardData({
              data: text,
              success: () => {
                uni.showToast({ title: "已复制推广文案", icon: "none" });
              },
            });
            return;
          }
          if (res.tapIndex === 2) {
            uni.showToast({ title: "分享功能待接入", icon: "none" });
          }
        },
      });
    },
    async fetchList(reset) {
      if (this.loading) return false;
      this.loading = true;
      if (reset) {
        this.pagination.page = 1;
        this.pagination.total = 0;
        this.finished = false;
        this.propertyList = [];
      }

      const params = {
        page: this.pagination.page,
        pageSize: this.pagination.pageSize,
        province: this.provinceName || "辽宁",
        keyword: this.keyword || "",
        category: this.filters.category || "all",
        city: this.currentCityName || "",
        district: this.currentDistrictName || "",
        rooms: this.filters.rooms || "",
        sale_status: this.filters.sale_status || "",
        sort: this.filters.sort || "",
        hot_status: this.filters.hot_status,
        property_type: this.filters.property_type || "",
        decoration_type: this.filters.decoration_type || "",
        orientation: this.filters.orientation || "",
        build_year_min: this.filters.build_year_min || "",
        build_year_max: this.filters.build_year_max || "",
        price_min: this.filters.price_min || "",
        price_max: this.filters.price_max || "",
        area_min: this.filters.area_min || "",
        area_max: this.filters.area_max || "",
        filter: this.filters.smartLock ? "smart_lock" : "",
        more: Array.isArray(this.filters.more)
          ? this.filters.more.join(",")
          : "",
      };

      let res;
      try {
        res = await propertyApi.getList(params);
      } catch (e) {
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        this.loading = false;
        uni.stopPullDownRefresh();
        return false;
      }

      this.loading = false;
      uni.stopPullDownRefresh();
      if (!res || res.code !== 0) return false;

      const data = res.data || {};
      const items = Array.isArray(data.items) ? data.items : [];
      const total = Number(data.total || 0);

      this.pagination.total = total;
      if (reset) {
        this.propertyList = items;
      } else {
        this.propertyList = (
          Array.isArray(this.propertyList) ? this.propertyList : []
        ).concat(items);
      }

      if (items.length === 0) {
        this.finished = true;
        return true;
      }
      if (total > 0 && this.propertyList.length >= total) {
        this.finished = true;
      } else if (items.length < this.pagination.pageSize) {
        this.finished = true;
      }
      return true;
    },
    async loadMore() {
      if (this.loading || this.finished) return;
      const prev = this.pagination.page;
      this.pagination.page = prev + 1;
      const ok = await this.fetchList(false);
      if (!ok) this.pagination.page = prev;
    },
    onSearchConfirm() {
      this.fetchList(true);
    },
    openFilter() {
      this.filterOpen = true;
      this.draft = {
        category: this.filters.category,
        rooms: this.filters.rooms,
        smartLock: !!this.filters.smartLock,
        sale_status: this.filters.sale_status,
        sort: this.filters.sort,
        hot_status: this.filters.hot_status,
        property_type: this.filters.property_type,
        decoration_type: this.filters.decoration_type,
        orientation: this.filters.orientation,
        build_year_min: this.filters.build_year_min,
        build_year_max: this.filters.build_year_max,
        price_min: this.filters.price_min,
        price_max: this.filters.price_max,
        area_min: this.filters.area_min,
        area_max: this.filters.area_max,
        more: Array.isArray(this.filters.more) ? [...this.filters.more] : [],
      };
      this.loadFilterOptions({ category: this.draft.category });
    },
    closeFilter() {
      this.filterOpen = false;
    },
    resetFilter() {
      this.draft = {
        category: "all",
        rooms: "",
        smartLock: false,
        sale_status: "",
        sort: "",
        hot_status: "",
        property_type: "",
        decoration_type: "",
        orientation: "",
        build_year_min: "",
        build_year_max: "",
        price_min: "",
        price_max: "",
        area_min: "",
        area_max: "",
        more: [],
      };
      this.loadFilterOptions({ category: this.draft.category });
    },
    applyFilter() {
      this.filters = {
        category: this.draft.category,
        rooms: this.draft.rooms,
        smartLock: !!this.draft.smartLock,
        sale_status: this.draft.sale_status,
        sort: this.draft.sort,
        hot_status: this.draft.hot_status,
        property_type: this.draft.property_type,
        decoration_type: this.draft.decoration_type,
        orientation: this.draft.orientation,
        build_year_min: this.draft.build_year_min,
        build_year_max: this.draft.build_year_max,
        price_min: this.draft.price_min,
        price_max: this.draft.price_max,
        area_min: this.draft.area_min,
        area_max: this.draft.area_max,
        more: Array.isArray(this.draft.more) ? [...this.draft.more] : [],
      };
      this.closeFilter();
      this.loadFilterOptions({ category: this.filters.category });
      this.fetchList(true);
    },
    goToDetail(item) {
      const id = item && (item.id || item.ID);
      if (!id) {
        uni.showToast({ title: "房源ID缺失", icon: "none" });
        return;
      }
      uni.navigateTo({
        url: `/pages/property_detail/property_detail?id=${encodeURIComponent(
          id
        )}`,
      });
    },
    ensureTags(tags) {
      if (!tags) return [];
      if (Array.isArray(tags)) {
        return tags
          .map((t) => {
            if (!t) return "";
            if (typeof t === "string") return t.trim();
            if (typeof t === "object")
              return String(t.name || t.label || "").trim();
            return String(t).trim();
          })
          .filter(Boolean);
      }
      if (typeof tags === "string") {
        return tags
          .split(",")
          .map((t) => t.trim())
          .filter(Boolean);
      }
      return [];
    },
    tagColorClass(tag) {
      const t = String(tag || "").trim();
      if (!t) return "grey";
      if (t.includes("地铁") || t.includes("交通")) return "green";
      if (t.includes("学区")) return "indigo";
      if (t.includes("降价") || t.includes("急售")) return "orange";
      if (t.includes("智能") || t.includes("门锁")) return "blue";
      if (t.includes("新") || t.includes("刚需")) return "purple";
      return "grey";
    },
    isSoldOrOff(saleStatus) {
      return saleStatus === "sold" || saleStatus === "off_market";
    },
    getLayoutText(item) {
      const rooms = Number(item && item.rooms) || 0;
      const halls = Number(item && item.halls) || 0;
      const bathrooms = Number(item && item.bathrooms) || 0;
      if (!rooms && !halls && !bathrooms) return "-";
      let s = "";
      if (rooms) s += `${rooms}室`;
      if (halls) s += `${halls}厅`;
      if (bathrooms) s += `${bathrooms}卫`;
      return s || "-";
    },
    extractDistrictFromAddress(address) {
      const addr = String(address || "").trim();
      if (!addr) return "";
      const m1 = addr.match(/(?:[^省]+省)?([^市]+市)([^区县]+(?:区|县|旗))/);
      if (m1) return `${m1[1]}${m1[2]}`;
      const m2 = addr.match(/([^市]+市)(.+?(?:区|县|旗))/);
      if (m2) return `${m2[1]}${m2[2]}`;
      return "";
    },
    getDistrictText(item) {
      const addr = String((item && item.address) || "").trim();
      if (!addr) return "-";
      const district = this.extractDistrictFromAddress(addr);
      if (district) return district;
      return addr.length > 20 ? `${addr.slice(0, 20)}…` : addr;
    },
    getUnitPriceText(item) {
      const price = Number(item && item.price);
      const area = Number(item && item.area);
      if (!price || !area) return "";
      const unit = (item && item.price_unit) || "";
      let totalYuan = price;
      if (unit === "万") totalYuan = price * 10000;
      const per = totalYuan / area;
      if (!per || !isFinite(per)) return "";
      return `${Math.round(per).toLocaleString()}元/㎡`;
    },
    getCommissionText(item) {
      const rate = String((item && item.commission_rate) || "").trim();
      const reward = String((item && item.commission_reward) || "").trim();
      const hasRate =
        !!rate && rate !== "0" && rate !== "0.0" && rate !== "0.00";
      const hasReward =
        !!reward && reward !== "0" && reward !== "0.0" && reward !== "0.00";
      if (hasRate && hasReward) return `佣金${rate}% 成交奖励¥${reward}`;
      if (hasRate) return `佣金${rate}%`;
      if (hasReward) return `成交奖励¥${reward}`;
      return "佣金待定";
    },
  },
};
</script>

<style lang="scss">
.list-container {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
}

.header {
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  /* TopHeader 自带左右 padding，这里避免叠加 */
  padding: 0 0 16rpx;
  border-bottom: 1rpx solid #e2e8f0;
  z-index: 50;

  .nav-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 80rpx;
    margin-bottom: 16rpx;
    position: relative;

    .nav-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #0f172a;
      flex: 1;
      text-align: left;
      padding-right: 80rpx;
    }

    .filter-icon-btn {
      position: absolute;
      right: 0;
      width: 80rpx;
      height: 80rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 50%;
      position: relative;

      .material-symbols-outlined {
        font-size: 40rpx;
        color: #0f172a;
      }

      &:active {
        background-color: #f1f5f9;
      }

      &.active .material-symbols-outlined {
        color: #2d9cf0;
      }
    }

    .filter-dot {
      position: absolute;
      top: 14rpx;
      right: 14rpx;
      width: 14rpx;
      height: 14rpx;
      background-color: #f97316;
      border-radius: 50%;
      border: 3rpx solid rgba(255, 255, 255, 0.95);
    }
  }

  .search-row {
    margin-top: 16rpx;
    margin-bottom: 24rpx;
    padding: 0 16rpx;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    gap: 16rpx;

    .province-picker {
      display: flex;
      align-items: center;
      gap: 6rpx;
      padding: 0 16rpx;
      height: 88rpx;
      background-color: #ffffff;
      border-radius: 20rpx;
      border: 1px solid rgba(226, 232, 240, 0.95);
      box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
      flex-shrink: 0;
      position: relative;
      z-index: 2;

      &:active {
        background-color: #f1f5f9;
      }

      .province-text {
        font-size: 24rpx;
        font-weight: 700;
        color: #0f172a;
      }

      .province-icon {
        font-size: 28rpx;
        color: #64748b;
      }
    }

    .search-box {
      flex: 1;
      height: 88rpx;
      background-color: #f8fafc;
      border-radius: 20rpx;
      display: flex;
      align-items: center;
      padding: 0 24rpx;
      border: 1px solid rgba(226, 232, 240, 0.95);
      box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);

      .search-icon {
        font-size: 40rpx;
        color: #94a3b8;
      }

      .search-input {
        flex: 1;
        margin-left: 16rpx;
        font-size: 28rpx;
        color: #0f172a;
      }
    }
  }

  .category-scroll {
    white-space: nowrap;
    width: 100%;
    padding: 0 16rpx;
    box-sizing: border-box;

    .category-row {
      display: inline-flex;
      gap: 16rpx;
      padding-bottom: 4rpx;
    }

    .category-item {
      padding: 12rpx 32rpx;
      background-color: #ffffff;
      border-radius: 40rpx;
      font-size: 24rpx;
      font-weight: 500;
      color: #64748b;
      border: 1px solid #e2e8f0;
      white-space: nowrap;

      &.active {
        background-color: #2d9cf0;
        color: #ffffff;
        border-color: #2d9cf0;
        box-shadow: 0 4rpx 8rpx rgba(45, 156, 240, 0.2);
      }
    }
  }
}

/* TopHeader 右侧按钮：统一尺寸/交互 */
.filter-icon-btn {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  position: relative;
  background-color: transparent;

  .material-symbols-outlined {
    font-size: 40rpx;
    color: #0f172a;
  }

  &:active {
    background-color: #f1f5f9;
  }

  &.active .material-symbols-outlined {
    color: #2d9cf0;
  }
}

.filter-dot {
  position: absolute;
  top: 14rpx;
  right: 14rpx;
  width: 14rpx;
  height: 14rpx;
  background-color: #f97316;
  border-radius: 50%;
  border: 3rpx solid rgba(255, 255, 255, 0.95);
}

.main-list {
  flex: 1;
  overflow: hidden;
}

.list-content {
  padding: 24rpx 16rpx;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.empty-wrap {
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 56rpx 24rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
  border: 1px dashed #e2e8f0;

  .empty-icon {
    font-size: 80rpx;
    color: #94a3b8;
  }

  .empty-title {
    font-size: 28rpx;
    font-weight: 700;
    color: #334155;
  }

  .empty-desc {
    font-size: 22rpx;
    color: #94a3b8;
  }
}

.loading-row {
  display: flex;
  justify-content: center;
  padding: 8rpx 0 16rpx;
  font-size: 22rpx;
  color: #64748b;

  &.muted {
    color: #94a3b8;
  }
}

.property-card {
  background-color: #ffffff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
  border: 1px solid #f1f5f9;

  &.disabled {
    opacity: 0.8;

    .title {
      color: #64748b;
    }

    .price {
      color: #94a3b8;
    }
  }

  .card-main {
    padding: 20rpx;
    display: flex;
    gap: 18rpx;

    .image-box {
      width: 224rpx;
      height: 224rpx;
      border-radius: 16rpx;
      overflow: hidden;
      position: relative;
      flex-shrink: 0;
      background-color: #f1f5f9;

      .property-image {
        width: 100%;
        height: 100%;
      }

      &.grayscale {
        .property-image {
          filter: grayscale(1);
        }
      }

      .image-tag {
        position: absolute;
        top: 8rpx;
        left: 8rpx;
        background-color: rgba(0, 0, 0, 0.6);
        backdrop-filter: blur(4rpx);
        color: #ffffff;
        font-size: 20rpx;
        padding: 4rpx 12rpx;
        border-radius: 8rpx;
        display: flex;
        align-items: center;
        gap: 4rpx;

        .tag-icon {
          font-size: 20rpx;
        }
      }

      .mask {
        position: absolute;
        inset: 0;
        background-color: rgba(15, 23, 42, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;

        .mask-text {
          background-color: rgba(30, 41, 59, 0.8);
          color: #ffffff;
          font-size: 24rpx;
          font-weight: bold;
          padding: 4rpx 16rpx;
          border-radius: 8rpx;
        }
      }
    }

    .info-box {
      flex: 1;
      min-width: 0;
      display: flex;
      flex-direction: column;
      gap: 6rpx;

      .title {
        font-size: 30rpx;
        font-weight: bold;
        color: #0f172a;
        line-height: 1.3;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        line-clamp: 2;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }

      .meta {
        display: flex;
        align-items: center;
        gap: 12rpx;
        font-size: 24rpx;
        color: #64748b;

        .bold {
          font-weight: 500;
          color: #334155;
        }

        .region-row {
          display: flex;
          align-items: center;
          gap: 12rpx;
          font-size: 14rpx;
          color: #64748b;

          .region-label {
            color: #94a3b8;
            flex-shrink: 0;
          }

          .region-val {
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
        }

        .divider {
          color: #e2e8f0;
        }
      }

      .tags {
        display: flex;
        flex-wrap: wrap;
        gap: 12rpx;
        margin: 2rpx 0;

        .tag {
          font-size: 20rpx;
          font-weight: 500;
          padding: 4rpx 12rpx;
          border-radius: 8rpx;
          border: 1px solid transparent;

          &.blue {
            background-color: #eff6ff;
            color: #2563eb;
            border-color: #dbeafe;
          }
          &.green {
            background-color: #f0fdf4;
            color: #16a34a;
            border-color: #dcfce7;
          }
          &.orange {
            background-color: #fff7ed;
            color: #ea580c;
            border-color: #ffedd5;
          }
          &.indigo {
            background-color: #eef2ff;
            color: #4f46e5;
            border-color: #e0e7ff;
          }
          &.purple {
            background-color: #faf5ff;
            color: #9333ea;
            border-color: #f3e8ff;
          }
          &.grey {
            background-color: #f8fafc;
            color: #64748b;
          }
        }
      }

      .price-row {
        display: flex;
        align-items: baseline;
        gap: 16rpx;
        margin-top: auto;

        .price {
          font-size: 36rpx;
          font-weight: bold;
          color: #f97316;
        }

        .unit-price {
          font-size: 20rpx;
          color: #94a3b8;
        }
      }

      .stats {
        font-size: 20rpx;
        color: #94a3b8;
        display: flex;
        align-items: center;

        .stats-sep {
          width: 1px;
          height: 16rpx;
          background-color: #e2e8f0;
          margin: 0 12rpx;
        }
      }
    }
  }

  .card-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14rpx 20rpx;
    border-top: 1rpx solid transparent;

    &.orange-footer {
      background-color: #fff7ed;
      border-color: #ffedd5;

      .footer-icon-box {
        background-color: #f97316;
      }
      .footer-text {
        color: #9a3412;
        font-weight: bold;
      }
    }

    &.grey-footer {
      background-color: #f8fafc;
      border-color: #f1f5f9;

      .footer-icon-box {
        background-color: #cbd5e1;
      }
      .footer-text {
        color: #64748b;
        font-weight: 500;
      }
    }

    .footer-left {
      display: flex;
      align-items: center;
      gap: 16rpx;
      flex: 1;
      min-width: 0;

      .footer-icon-box {
        width: 40rpx;
        height: 40rpx;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;

        .footer-icon {
          color: #ffffff;
          font-size: 28rpx;
        }
      }

      .footer-text {
        font-size: 24rpx;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .footer-btn {
      height: 52rpx;
      padding: 0 24rpx;
      background-color: #ffffff;
      border: 1px solid #fed7aa;
      border-radius: 40rpx;
      display: flex;
      align-items: center;
      gap: 8rpx;
      font-size: 22rpx;
      font-weight: bold;
      color: #ea580c;
      line-height: 1;
      margin-left: 24rpx;

      &::after {
        border: none;
      }

      .btn-icon {
        font-size: 28rpx;
      }

      &:active {
        background-color: #fff7ed;
        transform: scale(0.95);
      }
    }

    .footer-btn-disabled {
      height: 52rpx;
      padding: 0 24rpx;
      background-color: #ffffff;
      border: 1px solid #e2e8f0;
      border-radius: 40rpx;
      font-size: 22rpx;
      color: #94a3b8;
      margin-left: 24rpx;
      &::after {
        border: none;
      }
    }
  }
}

.tab-bar {
  background-color: #ffffff;
  border-top: 1rpx solid #e2e8f0;
  padding: 16rpx 48rpx calc(env(safe-area-inset-bottom) + 16rpx);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 100;

  .tab-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8rpx;
    color: #94a3b8;

    .tab-icon {
      font-size: 52rpx;

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
  }

  .add-btn-container {
    position: relative;
    top: -40rpx;

    .add-btn {
      width: 112rpx;
      height: 112rpx;
      background-color: #2d9cf0;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 8rpx 20rpx rgba(45, 156, 240, 0.4);

      .material-symbols-outlined {
        color: #ffffff;
        font-size: 48rpx;
      }

      &:active {
        transform: scale(0.9) translateY(4rpx);
      }
    }
  }
}

.bottom-spacer {
  height: 160rpx;
}

.placeholder {
  color: #94a3b8;
}

/* 筛选弹层（避免与卡片内 .mask 冲突，单独命名） */
.filter-mask {
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

.filter-sheet {
  width: 100%;
  background-color: #ffffff;
  border-radius: 32rpx 32rpx 0 0;
  padding: 24rpx 24rpx calc(env(safe-area-inset-bottom) + 24rpx);
  box-sizing: border-box;
  max-height: 82vh;
  height: 82vh;
  display: flex;
  flex-direction: column;
}

.filter-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.filter-body {
  flex: 1;
  overflow: hidden;
}

.filter-title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
}

.filter-reset {
  font-size: 26rpx;
  color: #64748b;
}

.filter-section {
  padding-top: 18rpx;
}

.filter-section-title {
  font-size: 26rpx;
  font-weight: 700;
  color: #334155;
  margin-bottom: 12rpx;
  display: block;
}

.muted-tip {
  width: 100%;
  font-size: 22rpx;
  color: #94a3b8;
  line-height: 34rpx;
}

.chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.field-row {
  display: flex;
}

.range-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.range-sep {
  color: #94a3b8;
  font-size: 22rpx;
}

.field-input {
  flex: 1;
  height: 80rpx;
  padding: 0 18rpx;
  border-radius: 16rpx;
  background-color: #f8fafc;
  border: 1px solid rgba(226, 232, 240, 0.95);
  color: #0f172a;
  font-size: 26rpx;
}

.chip {
  padding: 12rpx 22rpx;
  border-radius: 40rpx;
  background-color: #f1f5f9;
  color: #334155;
  font-size: 24rpx;
  border: 1px solid transparent;

  &.active {
    background-color: rgba(45, 156, 240, 0.12);
    color: #2d9cf0;
    border-color: rgba(45, 156, 240, 0.25);
    font-weight: 700;
  }
}

.radio-list {
  border-radius: 20rpx;
  overflow: hidden;
  border: 1px solid #f1f5f9;
}

.radio-item {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 18rpx;
  border-bottom: 1rpx solid #f1f5f9;
  color: #334155;

  &:last-child {
    border-bottom: none;
  }

  &.active {
    color: #2d9cf0;
    font-weight: 700;
  }
}

.radio-check {
  font-size: 36rpx;
}

.filter-footer {
  display: flex;
  gap: 16rpx;
  margin-top: 20rpx;
}

.filter-btn {
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
</style>
