<template>
  <view class="pm">
    <TopHeader title="房源管理">
      <template #left>
        <view class="icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </template>
      <template #right>
        <view class="icon-btn primary" @tap="goAdd">
          <text class="material-symbols-outlined">add</text>
        </view>
      </template>
    </TopHeader>

    <view class="toolbar">
      <view class="search-box">
        <text class="material-symbols-outlined search-icon">search</text>
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          confirm-type="search"
          placeholder="搜索标题/小区..."
          placeholder-class="placeholder"
          @confirm="onSearchConfirm"
        />
        <view v-if="keyword" class="clear-btn" @tap="clearKeyword">
          <text class="material-symbols-outlined">close</text>
        </view>
      </view>

      <scroll-view scroll-x="true" class="tabs">
        <view class="tab-row">
          <view
            v-for="(t, idx) in saleStatusTabs"
            :key="idx"
            class="tab"
            :class="{ active: saleStatus === t.value }"
            @tap="onSaleStatus(t.value)"
          >
            {{ t.label }}
          </view>
        </view>
      </scroll-view>
    </view>

    <scroll-view
      scroll-y="true"
      class="list"
      lower-threshold="120"
      @scrolltolower="loadMore"
    >
      <view class="list-content">
        <view
          class="card"
          v-for="(item, index) in list"
          :key="item.id || index"
        >
          <view class="card-main" @tap="goDetail(item)">
            <view class="cover-box">
              <image
                class="cover"
                :src="item.image || fallbackCover"
                mode="aspectFill"
              ></image>
              <view class="status-pill" :class="statusPillClass(item)">
                {{ item.sale_status_label || '-' }}
              </view>
            </view>
            <view class="info">
              <view class="title">{{ item.title || '-' }}</view>
              <view class="sub">
                <text class="price"
                  >¥{{ item.price || '-' }}{{ item.price_unit || '' }}</text
                >
                <text class="sep">·</text>
                <text class="meta">{{
                  (item.rooms || 0) + '室' + (item.halls || 0) + '厅'
                }}</text>
                <text class="sep">·</text>
                <text class="meta">{{ item.area || '-' }}㎡</text>
              </view>
              <view class="tags">
                <text
                  class="tag"
                  v-for="(tag, tIdx) in ensureTags(item.tags).slice(0, 4)"
                  :key="tIdx"
                  >{{ tag }}</text
                >
              </view>
              <view class="addr"
                >{{ (item.community_name || '') + (item.address ? ' · ' + item.address : '') }}</view
              >
            </view>
          </view>

          <view class="actions">
            <button class="btn ghost" @tap.stop="goEdit(item)">
              <text class="material-symbols-outlined btn-ic">edit</text>
              <text>编辑</text>
            </button>

            <button class="btn danger" @tap.stop="confirmDelete(item)">
              <text class="material-symbols-outlined btn-ic">delete</text>
              <text>删除</text>
            </button>

            <button
              class="btn"
              :class="toggleBtnClass(item)"
              @tap.stop="confirmToggleSaleStatus(item)"
              :disabled="String(item.sale_status) === 'sold'"
            >
              <text class="material-symbols-outlined btn-ic">{{
                String(item.sale_status) === 'off_market' ? 'publish' : 'unpublished'
              }}</text>
              <text>{{
                String(item.sale_status) === 'sold'
                  ? '已售'
                  : String(item.sale_status) === 'off_market'
                  ? '上架'
                  : '下架'
              }}</text>
            </button>
          </view>
        </view>

        <view v-if="!loading && list.length === 0" class="empty-wrap">
          <text class="material-symbols-outlined empty-icon">home_work</text>
          <view class="empty-title">暂无可维护房源</view>
          <view class="empty-desc">你可以点击右上角新增房源</view>
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
  </view>
</template>

<script>
import TopHeader from '@/components/TopHeader.vue'
import $store from '@/store'
import propertyApi from '@/api/property'

export default {
  components: { TopHeader },
  data() {
    return {
      fallbackCover: '/static/images/img_cdc09ae543.png',
      loading: false,
      finished: false,
      page: 1,
      pageSize: 10,
      total: 0,
      list: [],
      keyword: '',
      saleStatus: '',
      saleStatusTabs: [
        { label: '全部', value: '' },
        { label: '在售', value: 'on_sale' },
        { label: '已售', value: 'sold' },
        { label: '下架', value: 'off_market' },
      ],
      inited: false,
    }
  },
  onShow() {
    this.ensurePermissionAndLoad()
  },
  methods: {
    ensureTags(v) {
      if (Array.isArray(v)) return v
      if (!v) return []
      return String(v)
        .split(',')
        .map((s) => String(s || '').trim())
        .filter(Boolean)
    },
    statusPillClass(item) {
      const ss = String(item && item.sale_status)
      if (ss === 'sold') return 'sold'
      if (ss === 'off_market') return 'off'
      return 'on'
    },
    toggleBtnClass(item) {
      const ss = String(item && item.sale_status)
      if (ss === 'sold') return 'disabled'
      if (ss === 'off_market') return 'primary'
      return 'warn'
    },
    async ensurePermissionAndLoad() {
      const userStore = $store('user')
      const token = uni.getStorageSync('token')
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: '/pages/login/login' })
        return
      }
      try {
        if (!userStore.userInfo || !userStore.userInfo.id) {
          await userStore.getInfo()
        }
      } catch (e) {}
      const u = userStore.userInfo || {}
      const canManage = Number(u.can_manage_properties) === 1
      if (!canManage) {
        uni.showModal({
          title: '提示',
          content: '暂无房源维护权限，请联系后台管理员开启。',
          showCancel: false,
          success: () => {
            uni.navigateBack()
          },
        })
        return
      }

      if (!this.inited) this.inited = true
      this.reload()
    },
    goBack() {
      uni.navigateBack()
    },
    goAdd() {
      uni.navigateTo({ url: '/pages/property_manage/property_edit' })
    },
    goEdit(item) {
      const id = Number(item && item.id) || 0
      if (!id) return
      uni.navigateTo({ url: `/pages/property_manage/property_edit?id=${id}` })
    },
    goDetail(item) {
      const id = Number(item && item.id) || 0
      if (!id) return
      uni.navigateTo({ url: `/pages/property_detail/property_detail?id=${id}` })
    },
    clearKeyword() {
      this.keyword = ''
      this.reload()
    },
    onSearchConfirm() {
      this.reload()
    },
    onSaleStatus(v) {
      this.saleStatus = v
      this.reload()
    },
    async reload() {
      if (this.loading) return
      this.page = 1
      this.total = 0
      this.list = []
      this.finished = false
      await this.fetchList()
    },
    async loadMore() {
      await this.fetchList()
    },
    async fetchList() {
      if (this.loading || this.finished) return
      this.loading = true
      try {
        const res = await propertyApi.getManageList({
          page: this.page,
          pageSize: this.pageSize,
          keyword: String(this.keyword || '').trim(),
          sale_status: this.saleStatus,
        })
        if (!res || res.code !== 0) {
          if (res && Number(res.code) === 403) {
            uni.showToast({ title: '暂无权限', icon: 'none' })
          }
          return
        }
        const data = res.data || {}
        const items = Array.isArray(data.items) ? data.items : []
        const total = Number(data.total || 0)
        this.total = total
        if (this.page === 1) this.list = items
        else this.list = this.list.concat(items)

        if (items.length === 0 || (total > 0 && this.list.length >= total)) {
          this.finished = true
        } else {
          this.page += 1
        }
      } finally {
        this.loading = false
      }
    },
    confirmDelete(item) {
      const id = Number(item && item.id) || 0
      if (!id) return
      uni.showModal({
        title: '确认删除',
        content: '确定要删除该房源吗？（软删除，可在后台恢复）',
        confirmText: '删除',
        confirmColor: '#ef4444',
        success: async (res) => {
          if (!res.confirm) return
          const out = await propertyApi.delManage({ id })
          if (!out || out.code !== 0) return
          uni.showToast({ title: '已删除', icon: 'none' })
          this.reload()
        },
      })
    },
    confirmToggleSaleStatus(item) {
      const id = Number(item && item.id) || 0
      if (!id) return
      const cur = String(item && item.sale_status)
      if (cur === 'sold') return
      const next = cur === 'off_market' ? 'on_sale' : 'off_market'
      const nextLabel = next === 'on_sale' ? '上架' : '下架'
      uni.showModal({
        title: '确认操作',
        content: `确定要${nextLabel}该房源吗？`,
        confirmText: nextLabel,
        success: async (res) => {
          if (!res.confirm) return
          const out = await propertyApi.saveManage({ id, sale_status: next })
          if (!out || out.code !== 0) return
          uni.showToast({ title: '已更新', icon: 'none' })
          this.reload()
        },
      })
    },
  },
}
</script>

<style>
.pm {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.icon-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(226, 232, 240, 0.9);
  background: rgba(255, 255, 255, 0.8);
  color: #0f172a;
}

.icon-btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
  box-shadow: 0 12rpx 28rpx rgba(37, 99, 235, 0.22);
}

.toolbar {
  padding: 14rpx 24rpx 10rpx;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.search-box {
  background-color: #ffffff;
  border-radius: 22rpx;
  border: 1px solid #e2e8f0;
  height: 76rpx;
  display: flex;
  align-items: center;
  padding: 0 18rpx;
  box-sizing: border-box;
  gap: 10rpx;
}

.search-icon {
  font-size: 40rpx !important;
  color: #94a3b8;
}

.search-input {
  flex: 1;
  height: 76rpx;
  font-size: 28rpx;
  color: #0f172a;
}

.placeholder {
  color: #94a3b8;
}

.clear-btn {
  width: 52rpx;
  height: 52rpx;
  border-radius: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  background-color: #f1f5f9;
}

.tabs {
  width: 100%;
}

.tab-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding-bottom: 6rpx;
}

.tab {
  padding: 14rpx 22rpx;
  border-radius: 999rpx;
  background-color: #ffffff;
  border: 1px solid #e2e8f0;
  color: #334155;
  font-size: 24rpx;
  white-space: nowrap;
}

.tab.active {
  background-color: #eaf3ff;
  border-color: rgba(37, 99, 235, 0.25);
  color: #2563eb;
  font-weight: 700;
}

.list {
  flex: 1;
  overflow: hidden;
}

.list-content {
  padding: 0 24rpx 24rpx;
  box-sizing: border-box;
}

.card {
  background-color: #ffffff;
  border-radius: 28rpx;
  border: 1px solid #eef2f7;
  overflow: hidden;
  margin-bottom: 18rpx;
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.05);
}

.card-main {
  padding: 18rpx;
  display: flex;
  gap: 18rpx;
}

.cover-box {
  width: 220rpx;
  height: 168rpx;
  border-radius: 18rpx;
  overflow: hidden;
  position: relative;
  flex-shrink: 0;
  background-color: #f1f5f9;
}

.cover {
  width: 100%;
  height: 100%;
}

.status-pill {
  position: absolute;
  left: 12rpx;
  top: 12rpx;
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  color: #ffffff;
  background: rgba(15, 23, 42, 0.7);
  backdrop-filter: blur(8px);
}

.status-pill.on {
  background: rgba(34, 197, 94, 0.85);
}
.status-pill.sold {
  background: rgba(234, 88, 12, 0.9);
}
.status-pill.off {
  background: rgba(148, 163, 184, 0.9);
}

.info {
  flex: 1;
  min-width: 0;
}

.title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.25;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.sub {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #64748b;
  display: flex;
  align-items: center;
  gap: 8rpx;
  flex-wrap: wrap;
}

.price {
  color: #0f172a;
  font-weight: 800;
}

.sep {
  color: #cbd5e1;
}

.meta {
  color: #64748b;
}

.tags {
  margin-top: 10rpx;
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
}

.tag {
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  background-color: #f1f5f9;
  color: #334155;
  font-size: 22rpx;
}

.addr {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #94a3b8;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.actions {
  display: flex;
  gap: 12rpx;
  padding: 0 18rpx 18rpx;
  box-sizing: border-box;
}

.btn {
  flex: 1;
  height: 72rpx;
  border-radius: 20rpx;
  font-size: 26rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  border: 1px solid #e2e8f0;
  background-color: #ffffff;
  color: #0f172a;
}

.btn::after {
  border: none;
}

.btn-ic {
  font-size: 30rpx !important;
}

.btn.ghost {
  background-color: #f8fafc;
  color: #2563eb;
  border-color: rgba(37, 99, 235, 0.22);
}

.btn.danger {
  background-color: #fff1f2;
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.22);
}

.btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
}

.btn.warn {
  background: linear-gradient(135deg, #f97316, #ea580c);
  border: none;
  color: #ffffff;
}

.btn.disabled {
  background-color: #f1f5f9;
  color: #94a3b8;
  border-color: #e2e8f0;
}

.empty-wrap {
  padding: 120rpx 0 60rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #94a3b8;
}

.empty-icon {
  font-size: 86rpx !important;
  color: #cbd5e1;
}

.empty-title {
  margin-top: 16rpx;
  font-size: 30rpx;
  font-weight: 800;
  color: #334155;
}

.empty-desc {
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #94a3b8;
}

.loading-row {
  text-align: center;
  padding: 24rpx 0;
  color: #334155;
  font-size: 26rpx;
}

.loading-row.muted {
  color: #94a3b8;
}

.bottom-spacer {
  height: calc(env(safe-area-inset-bottom) + 24rpx);
}
</style>
