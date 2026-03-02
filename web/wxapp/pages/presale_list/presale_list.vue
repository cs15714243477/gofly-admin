<template>
  <view class="list-container">
    <!-- 顶部导航与搜索 -->
    <view class="header">
      <TopHeader title="预售" />

      <view class="search-row">
        <view class="search-box">
          <text class="material-symbols-outlined search-icon">search</text>
          <input
            v-model="keyword"
            class="search-input"
            type="text"
            confirm-type="search"
            placeholder="搜索小区、商圈、地址..."
            placeholder-class="placeholder"
            @confirm="onSearchConfirm"
          />
        </view>
      </view>
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
          v-for="(item, index) in list"
          :key="item.id || index"
          @click="goToDetail(item)"
        >
          <view class="card-main">
            <view class="image-box">
              <image
                v-if="safeImage(item.image)"
                class="property-image"
                :src="safeImage(item.image)"
                mode="aspectFill"
              ></image>
              <view v-else class="property-image image-empty">
                <text class="material-symbols-outlined">image</text>
              </view>
              <view class="status-badge" v-if="item.sale_status_label">
                {{ item.sale_status_label }}
              </view>
              <view class="image-tag" v-if="Number(item.has_smart_lock) === 1">
                <text class="material-symbols-outlined tag-icon">lock</text>
                <text>智能门锁</text>
              </view>
            </view>

            <view class="info-box">
              <view class="title">{{ item.title || "-" }}</view>
              <view class="meta">
                <text class="bold">{{ getLayoutText(item) }}</text>
                <text class="divider">|</text>
                <text>{{ item.area || "-" }}㎡</text>
                <text class="divider">|</text>
                <text>{{ item.orientation || "-" }}</text>
              </view>
              <view class="addr">
                <text class="material-symbols-outlined addr-icon"
                  >location_on</text
                >
                <text class="addr-text">{{
                  item.community_name || item.address || "-"
                }}</text>
              </view>
              <view class="tags" v-if="ensureTags(item.tags).length > 0">
                <text
                  class="tag"
                  v-for="(tag, tIdx) in ensureTags(item.tags).slice(0, 4)"
                  :key="tIdx"
                >
                  {{ tag }}
                </text>
              </view>
              <view class="price-row">
                <text class="price"
                  >¥{{ item.price || "-" }}{{ item.price_unit || "" }}</text
                >
              </view>
              <view class="stats">
                <text>浏览: {{ item.view_count || 0 }}</text>
                <text class="stats-sep"></text>
                <text>关注: {{ item.follow_count || 0 }}</text>
              </view>
            </view>
          </view>
        </view>

        <view v-if="!loading && list.length === 0" class="empty-wrap">
          <text class="material-symbols-outlined empty-icon">apartment</text>
          <view class="empty-title">暂无预售房源</view>
          <view class="empty-desc">请先在后台把房源设置为预售</view>
        </view>

        <view v-if="loading" class="loading-row">
          <text>加载中...</text>
        </view>
        <view v-else-if="finished && list.length > 0" class="loading-row muted">
          <text>没有更多了</text>
        </view>
      </view>
      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部导航 -->
    <BottomTabBar active="presale" />
  </view>
</template>

<script>
import BottomTabBar from "@/components/BottomTabBar.vue";
import TopHeader from "@/components/TopHeader.vue";
import propertyApi from "@/api/property";

export default {
  components: { BottomTabBar, TopHeader },
  data() {
    return {
      keyword: "",
      page: 1,
      pageSize: 10,
      total: 0,
      list: [],
      loading: false,
      finished: false,
    };
  },
  onLoad() {
    this.loadList(true);
  },
  onShow() {
    // 这里不强制刷新，避免频繁请求
  },
  onPullDownRefresh() {
    this.loadList(true).finally(() => {
      uni.stopPullDownRefresh();
    });
  },
  methods: {
    safeImage(url) {
      const u = String(url || "").trim();
      if (!u) return "";
      if (/^https?:\/\//i.test(u)) return u;
      return u;
    },
    ensureTags(tags) {
      if (!tags) return [];
      if (Array.isArray(tags)) return tags.filter((t) => String(t || "").trim());
      const s = String(tags || "").trim();
      if (!s) return [];
      return s
        .split(",")
        .map((t) => String(t || "").trim())
        .filter(Boolean);
    },
    getLayoutText(item) {
      const rooms = Number(item && item.rooms) || 0;
      const halls = Number(item && item.halls) || 0;
      const bathrooms = Number(item && item.bathrooms) || 0;
      if (!rooms && !halls && !bathrooms) return "-";
      return `${rooms || 0}室${halls || 0}厅${bathrooms || 0}卫`;
    },
    onSearchConfirm() {
      this.loadList(true);
    },
    async loadList(reset = false) {
      if (this.loading) return false;
      if (reset) {
        this.page = 1;
        this.finished = false;
      }
      if (this.finished) return true;

      this.loading = true;
      let res;
      try {
        res = await propertyApi.getPreSaleList({
          page: this.page,
          pageSize: this.pageSize,
          keyword: String(this.keyword || "").trim(),
        });
      } catch (e) {
        this.loading = false;
        return false;
      }
      this.loading = false;

      if (!res || res.code !== 0) {
        return false;
      }
      const data = res.data || {};
      const items = Array.isArray(data.items) ? data.items : [];
      const total = Number(data.total || 0) || 0;
      this.total = total;
      if (reset) this.list = items;
      else this.list = (this.list || []).concat(items);

      // 分页结束判断：优先用 total，其次用本次返回数量
      if (total > 0) {
        this.finished = this.list.length >= total;
      } else {
        this.finished = items.length < this.pageSize;
      }
      if (!this.finished) this.page += 1;
      return true;
    },
    loadMore() {
      if (this.loading || this.finished) return;
      this.loadList(false);
    },
    goToDetail(item) {
      const id = Number(item && item.id) || 0;
      if (!id) return;
      const viewKey = String((item && item.view_key) || "").trim();
      // A：只能从预售 Tab 进入；这里使用 view_key 做后端校验
      uni.navigateTo({
        url: `/pages/property_detail/property_detail?id=${encodeURIComponent(
          id,
        )}&presale=1&public=1&view_key=${encodeURIComponent(viewKey)}`,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.list-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
  overflow: hidden;
}

.header {
  background: #f8fafc;
}

.search-row {
  padding: 16rpx 24rpx 20rpx;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 18rpx 20rpx;
  border-radius: 18rpx;
  background: rgba(255, 255, 255, 0.96);
  border: 1rpx solid rgba(226, 232, 240, 0.9);
  box-shadow: 0 8rpx 18rpx rgba(15, 23, 42, 0.04);
}

.search-icon {
  font-size: 36rpx;
  color: #94a3b8;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  color: #0f172a;
}

.placeholder {
  color: #94a3b8;
}

.main-list {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.list-content {
  padding: 0 24rpx;
}

.property-card {
  margin-top: 18rpx;
  border-radius: 22rpx;
  background: rgba(255, 255, 255, 0.96);
  border: 1rpx solid rgba(226, 232, 240, 0.9);
  overflow: hidden;
  box-shadow: 0 10rpx 22rpx rgba(15, 23, 42, 0.05);
}

.card-main {
  display: flex;
  padding: 18rpx;
  gap: 18rpx;
}

.image-box {
  width: 210rpx;
  height: 210rpx;
  border-radius: 18rpx;
  overflow: hidden;
  position: relative;
  background: #e2e8f0;
}

.property-image {
  width: 100%;
  height: 100%;
}

.image-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15, 23, 42, 0.35);
}

.status-badge {
  position: absolute;
  top: 12rpx;
  left: 12rpx;
  padding: 8rpx 12rpx;
  border-radius: 999rpx;
  background: rgba(37, 99, 235, 0.92);
  color: #fff;
  font-size: 22rpx;
  font-weight: 700;
}

.image-tag {
  position: absolute;
  bottom: 12rpx;
  left: 12rpx;
  padding: 6rpx 10rpx;
  border-radius: 999rpx;
  background: rgba(15, 23, 42, 0.66);
  color: #fff;
  font-size: 20rpx;
  display: flex;
  align-items: center;
  gap: 6rpx;
}

.tag-icon {
  font-size: 22rpx;
}

.info-box {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.2;
}

.meta {
  font-size: 24rpx;
  color: #334155;
  display: flex;
  align-items: center;
  gap: 10rpx;
  flex-wrap: wrap;
}

.bold {
  font-weight: 800;
}

.divider {
  color: rgba(148, 163, 184, 0.9);
}

.addr {
  display: flex;
  align-items: center;
  gap: 6rpx;
  color: #64748b;
  font-size: 24rpx;
}

.addr-icon {
  font-size: 26rpx;
}

.addr-text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
}

.tag {
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  background: rgba(226, 232, 240, 0.7);
  color: #334155;
  font-size: 22rpx;
  font-weight: 600;
}

.price-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
}

.price {
  color: #f97316;
  font-weight: 900;
  font-size: 30rpx;
}

.stats {
  color: #94a3b8;
  font-size: 22rpx;
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.stats-sep {
  width: 1rpx;
  height: 18rpx;
  background: rgba(148, 163, 184, 0.6);
}

.empty-wrap {
  padding: 120rpx 0 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #94a3b8;
}

.empty-icon {
  font-size: 80rpx;
  margin-bottom: 12rpx;
}

.empty-title {
  color: #0f172a;
  font-size: 30rpx;
  font-weight: 800;
  margin-bottom: 8rpx;
}

.empty-desc {
  font-size: 24rpx;
}

.loading-row {
  text-align: center;
  padding: 24rpx 0 18rpx;
  color: #475569;
  font-size: 24rpx;
}

.loading-row.muted {
  color: #94a3b8;
}

.bottom-spacer {
  height: 140rpx;
}
</style>
