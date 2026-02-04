<template>
  <view class="home-container">
    <view class="header">
      <TopHeader title="首页" />
    </view>

    <scroll-view scroll-y="true" class="main-content">
      <!-- Banner + 入口：筛选统一去“房源列表” -->
      <view class="hero">
        <swiper
          class="hero-swiper"
          circular
          autoplay
          interval="4000"
          :indicator-dots="true"
          indicator-color="rgba(255,255,255,0.45)"
          indicator-active-color="#ffffff"
        >
          <swiper-item v-for="(img, idx) in bannerImages" :key="idx">
            <image class="hero-img" :src="img" mode="aspectFill"></image>
          </swiper-item>
        </swiper>
        <view class="hero-mask"></view>
        <view class="hero-actions">
          <view class="hero-title">精选房源</view>
          <view class="hero-sub">搜索与筛选请在房源列表中使用</view>
          <button class="hero-btn" @click="goToList">
            <text>查看房源列表</text>
          </button>
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
            <text class="title-sub">({{ followedTotal || 0 }})</text>
          </view>
          <view class="section-more" @click="goToList">
            <text>去房源列表</text>
            <text class="material-symbols-outlined more-icon"
              >chevron_right</text
            >
          </view>
        </view>

        <view
          v-if="!followedProperties || followedProperties.length === 0"
          class="empty-wrap"
        >
          <text class="material-symbols-outlined empty-icon">favorite</text>
          <view class="empty-title">暂无关注房源</view>
          <view class="empty-desc">在房源详情页点击关注后，会在这里展示</view>
          <button class="empty-btn" @click="goToList">去房源列表看看</button>
        </view>

        <scroll-view v-else scroll-x="true" class="horizontal-scroll">
          <view
            class="property-card-h"
            v-for="(item, index) in followedProperties"
            :key="item.id || index"
            @click="goToDetail(item)"
          >
            <view
              class="card-image"
              :class="{ disabled: isSoldOrOff(item.sale_status) }"
            >
              <image
                v-if="safeImage(item.image)"
                class="card-img"
                :src="safeImage(item.image)"
                mode="aspectFill"
              />
              <view v-else class="card-img card-img-empty">
                <text class="material-symbols-outlined">image</text>
              </view>
              <view
                v-if="item.sale_status_label"
                class="status-badge"
                :class="saleStatusClass(item.sale_status)"
              >
                {{ item.sale_status_label }}
              </view>
              <view class="lock-badge" v-if="Number(item.has_smart_lock) === 1">
                智能锁
              </view>
            </view>
            <view class="card-info">
              <view class="name">{{ item.name || "-" }}</view>
              <view class="meta">{{ formatFollowedMeta(item) }}</view>
              <view class="price"
                >￥{{ item.price || "-" }}{{ item.price_unit || "" }}</view
              >
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 推荐房源 -->
      <view class="section">
        <view class="section-header">
          <view class="section-title">
            <text class="material-symbols-outlined section-icon"
              >recommend</text
            >
            <text class="title-text">推荐房源</text>
          </view>
          <view class="section-more" @click="goToList">
            <text>更多</text>
            <text class="material-symbols-outlined more-icon"
              >chevron_right</text
            >
          </view>
        </view>

        <view v-if="loading && recommendedProperties.length === 0" class="row">
          <text class="muted">加载中...</text>
        </view>

        <view
          v-else-if="!loading && recommendedProperties.length === 0"
          class="empty-wrap"
        >
          <text class="material-symbols-outlined empty-icon">apartment</text>
          <view class="empty-title">暂无推荐房源</view>
          <view class="empty-desc">请先在后台把房源设置为推荐</view>
          <button class="empty-btn" @click="goToList">去房源列表看看</button>
        </view>

        <view v-else class="vertical-list">
          <view
            class="property-card-v"
            v-for="(item, index) in recommendedProperties"
            :key="item.id || index"
            @click="goToDetail(item)"
          >
            <view
              class="card-image-v"
              :class="{ disabled: isSoldOrOff(item.sale_status) }"
            >
              <image
                v-if="safeImage(item.image)"
                class="card-img-v"
                :src="safeImage(item.image)"
                mode="aspectFill"
              />
              <view v-else class="card-img-v card-img-empty">
                <text class="material-symbols-outlined">image</text>
              </view>
              <view
                v-if="item.sale_status_label"
                class="status-tag-v"
                :class="saleStatusClass(item.sale_status)"
              >
                {{ item.sale_status_label }}
              </view>
              <view class="lock-tag-v" v-if="Number(item.has_smart_lock) === 1">
                智能锁
              </view>
            </view>

            <view class="card-info-v">
              <view class="info-top">
                <view class="property-name-v">{{ item.name || "-" }}</view>
                <view class="price-v">
                  <text class="price-val">￥{{ item.price || "-" }}</text>
                  <text class="price-unit">{{ item.price_unit || "" }}</text>
                </view>
              </view>

              <view class="property-location">
                <text class="material-symbols-outlined loc-icon"
                  >location_on</text
                >
                <text class="loc-text">{{ formatRecommendedMeta(item) }}</text>
              </view>

              <view class="tags" v-if="ensureTags(item.tags).length > 0">
                <text
                  class="tag-item"
                  v-for="(tag, tIdx) in ensureTags(item.tags).slice(0, 4)"
                  :key="tIdx"
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

    <BottomTabBar active="home" />
  </view>
</template>

<script>
import BottomTabBar from "@/components/BottomTabBar.vue";
import TopHeader from "@/components/TopHeader.vue";
import homeApi from "@/api/home";

export default {
  components: { BottomTabBar, TopHeader },
  onLoad() {
    this.loadHomeData();
  },
  onShow() {
    this.loadHomeData();
  },
  onPullDownRefresh() {
    this.loadHomeData().finally(() => {
      uni.stopPullDownRefresh();
    });
  },
  data() {
    return {
      loading: false,
      bannerImages: [],
      followedTotal: 0,
      followedProperties: [],
      recommendedProperties: [],
    };
  },
  methods: {
    async loadHomeData() {
      if (this.loading) return false;
      this.loading = true;
      let res;
      try {
        res = await homeApi.Homedata({ page: 1, pageSize: 6 });
      } catch (e) {
        this.loading = false;
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return false;
      }
      this.loading = false;
      if (!res || res.code !== 0) return false;

      const data = res.data || {};
      this.bannerImages = Array.isArray(data.bannerImages)
        ? data.bannerImages
            .map((img) => this.safeImage(img))
            .filter(Boolean)
        : [];
      this.followedTotal = Number(data.followedTotal || 0);
      this.followedProperties = Array.isArray(data.followedProperties)
        ? data.followedProperties
        : [];
      this.recommendedProperties = Array.isArray(data.recommendedProperties)
        ? data.recommendedProperties
        : [];
      return true;
    },
    goToList() {
      uni.navigateTo({ url: "/pages/property_list/property_list" });
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
    safeImage(url) {
      const u = String(url || "").trim();
      if (!u) return "";
      if (u.indexOf("/static/images/") === 0) return "";
      return u;
    },
    ensureTags(v) {
      if (!v) return [];
      if (Array.isArray(v)) {
        return v
          .map((t) => {
            if (!t) return "";
            if (typeof t === "string") return t.trim();
            if (typeof t === "object")
              return String(t.name || t.label || "").trim();
            return String(t).trim();
          })
          .filter(Boolean);
      }
      if (typeof v === "string") {
        return v
          .split(",")
          .map((x) => x.trim())
          .filter(Boolean);
      }
      return [];
    },
    isSoldOrOff(v) {
      return v === "sold" || v === "off_market";
    },
    saleStatusClass(v) {
      const map = {
        on_sale: "badge-success",
        sold: "badge-default",
        off_market: "badge-warning",
      };
      return map[v] || "badge-default";
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
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1rpx solid #e2e8f0;
  z-index: 20;
}

.main-content {
  flex: 1;
  overflow: hidden;
}

.hero {
  margin: 24rpx 32rpx 0;
  border-radius: 28rpx;
  overflow: hidden;
  position: relative;
  height: 360rpx;
  background-color: #0f172a;
}

.hero-swiper {
  width: 100%;
  height: 100%;
}

.hero-img {
  width: 100%;
  height: 100%;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(15, 23, 42, 0.55),
    rgba(15, 23, 42, 0.15)
  );
}

.hero-actions {
  position: absolute;
  left: 24rpx;
  right: 24rpx;
  bottom: 24rpx;
  z-index: 2;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.hero-title {
  font-size: 34rpx;
  font-weight: 900;
  color: #ffffff;
}

.hero-sub {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
}

.hero-btn {
  margin-top: 8rpx;
  height: 84rpx;
  border-radius: 18rpx;
  background-color: #2d9cf0;
  color: #ffffff;
  font-size: 28rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  border: none;
}
.hero-btn::after {
  border: none;
}
.hero-btn:active {
  transform: scale(0.98);
  background-color: #1d82cc;
}

.section {
  margin-top: 22rpx;
}

.section-header {
  padding: 0 32rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14rpx;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.section-icon {
  font-size: 34rpx;
  color: #64748b;
}
.section-icon.accent {
  color: #f97316;
}

.title-text {
  font-size: 30rpx;
  font-weight: 900;
  color: #0f172a;
}

.title-sub {
  font-size: 22rpx;
  color: #94a3b8;
  font-weight: 700;
}

.section-more {
  display: flex;
  align-items: center;
  gap: 4rpx;
  font-size: 24rpx;
  color: #64748b;
}
.more-icon {
  font-size: 30rpx;
  color: #94a3b8;
}

.empty-wrap {
  margin: 0 32rpx;
  background-color: #ffffff;
  border-radius: 24rpx;
  padding: 34rpx 22rpx;
  border: 1px solid #f1f1f1;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
}
.empty-icon {
  font-size: 64rpx;
  color: rgba(249, 115, 22, 0.9);
}
.empty-title {
  font-size: 30rpx;
  font-weight: 900;
  color: #0f172a;
}
.empty-desc {
  font-size: 24rpx;
  color: #64748b;
  text-align: center;
}
.empty-btn {
  margin-top: 8rpx;
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
  background-color: #f1f5f9;
  color: #334155;
}
.empty-btn::after {
  border: none;
}

.row {
  padding: 20rpx 32rpx;
}
.muted {
  color: #94a3b8;
  font-size: 24rpx;
}

.horizontal-scroll {
  white-space: nowrap;
  padding-left: 32rpx;
  padding-right: 20rpx;
  box-sizing: border-box;
}

.property-card-h {
  width: 320rpx;
  display: inline-flex;
  flex-direction: column;
  gap: 12rpx;
  margin-right: 18rpx;
  background-color: #ffffff;
  border-radius: 20rpx;
  overflow: hidden;
  border: 1px solid #f1f1f1;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
}

.card-image {
  width: 100%;
  height: 176rpx;
  position: relative;
  background-color: #f1f5f9;
}
.card-img {
  width: 100%;
  height: 100%;
}
.card-img-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #eef2f7, #e2e8f0);
  .material-symbols-outlined {
    font-size: 52rpx;
    color: #94a3b8;
  }
}
.card-image.disabled {
  filter: grayscale(1) brightness(0.95);
}
.card-image.disabled::after {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(15, 23, 42, 0.25);
}
.status-badge {
  position: absolute;
  top: 12rpx;
  left: 12rpx;
  color: #ffffff;
  font-size: 20rpx;
  font-weight: 800;
  padding: 4rpx 14rpx;
  border-radius: 8rpx;
  background-color: rgba(100, 116, 139, 0.95);
}
.lock-badge {
  position: absolute;
  top: 12rpx;
  right: 12rpx;
  background-color: rgba(45, 156, 240, 0.92);
  color: #ffffff;
  font-size: 20rpx;
  font-weight: 700;
  padding: 4rpx 14rpx;
  border-radius: 8rpx;
}

.badge-success {
  background-color: rgba(34, 197, 94, 0.95);
}
.badge-warning {
  background-color: rgba(249, 115, 22, 0.95);
}
.badge-default {
  background-color: rgba(100, 116, 139, 0.95);
}

.card-info {
  padding: 18rpx 18rpx 20rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}
.name {
  font-size: 26rpx;
  font-weight: 900;
  color: #0f172a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.meta {
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.price {
  font-size: 24rpx;
  font-weight: 900;
  color: #2d9cf0;
}

.vertical-list {
  padding: 0 32rpx;
  display: flex;
  flex-direction: column;
  gap: 22rpx;
}

.property-card-v {
  background-color: #ffffff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
  border: 1px solid #f1f1f1;
  display: flex;
  flex-direction: column;
}

.card-image-v {
  width: 100%;
  height: 300rpx;
  background-size: cover;
  background-position: center;
  position: relative;
  background-color: #f1f5f9;
}
.card-img-v {
  width: 100%;
  height: 100%;
}
.card-image-v.disabled {
  filter: grayscale(1) brightness(0.95);
}
.card-image-v.disabled::after {
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

.card-info-v {
  padding: 20rpx;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}
.info-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 14rpx;
}
.property-name-v {
  font-size: 28rpx;
  font-weight: 900;
  color: #0f172a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}
.price-v {
  display: flex;
  align-items: baseline;
  color: #2d9cf0;
  flex-shrink: 0;
}
.price-val {
  font-size: 28rpx;
  font-weight: 900;
}
.price-unit {
  font-size: 22rpx;
  margin-left: 4rpx;
}
.property-location {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 22rpx;
  color: #64748b;
}
.loc-icon {
  font-size: 28rpx;
}
.loc-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}
.tag-item {
  padding: 4rpx 14rpx;
  background-color: #f1f5f9;
  border-radius: 8rpx;
  font-size: 20rpx;
  font-weight: 600;
  color: #475569;
}

.bottom-spacer {
  height: 120rpx;
}
</style>
