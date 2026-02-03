<template>
  <view class="lm">
    <TopHeader title="智能门锁">
      <template #left>
        <view class="icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </template>
      <template #right>
        <view class="right-actions">
          <view class="icon-btn ghost" :class="{ disabled: loading }" @tap="onRefresh">
            <text class="material-symbols-outlined">refresh</text>
          </view>
          <view class="icon-btn primary" :class="{ disabled: syncing }" @tap="handleSync">
            <text class="material-symbols-outlined">sync</text>
          </view>
        </view>
      </template>
    </TopHeader>

    <view class="toolbar">
      <view class="search-box">
        <text class="material-symbols-outlined search-ic">search</text>
        <input
          v-model="keyword"
          class="search-input"
          type="text"
          confirm-type="search"
          placeholder="搜索锁名 / ID / MAC..."
          placeholder-class="placeholder"
          @confirm="handleSearch"
        />
        <view v-if="keyword" class="clear-btn" @tap="clearKeyword">
          <text class="material-symbols-outlined">close</text>
        </view>
      </view>

      <view class="status-row">
        <view class="chip" :class="{ ok: status.config_ok }">
          <text class="material-symbols-outlined chip-ic">settings</text>
          <text class="chip-txt">{{ status.config_ok ? "配置正常" : "未配置" }}</text>
        </view>
        <view class="chip" :class="{ ok: status.has_token }">
          <text class="material-symbols-outlined chip-ic">key</text>
          <text class="chip-txt">{{ status.has_token ? "Token有效" : "Token缺失" }}</text>
        </view>
        <view class="chip ghost" v-if="status.api_base">
          <text class="material-symbols-outlined chip-ic">link</text>
          <text class="chip-txt">{{ status.api_base }}</text>
        </view>
      </view>
    </view>

    <scroll-view
      scroll-y="true"
      class="list"
      :refresher-enabled="true"
      :refresher-triggered="refreshing"
      lower-threshold="160"
      @refresherrefresh="onRefresh"
      @scrolltolower="loadMore"
    >
      <view class="list-content">
        <view v-if="filteredLocks.length === 0 && !loading" class="empty">
          <text class="material-symbols-outlined empty-ic">smart_lock</text>
          <text class="empty-title">暂无设备数据</text>
          <text class="empty-sub">请先点击右上角同步或稍后再试</text>
        </view>

        <view
          class="card"
          v-for="it in filteredLocks"
          :key="it.ttlock_lock_id"
          @tap="openDetail(it)"
        >
          <view class="card-head">
            <view class="head-left">
              <view class="lock-ic">
                <text class="material-symbols-outlined">smart_lock</text>
              </view>
              <view class="head-texts">
                <text class="lock-name">{{ it.lock_name || "未命名锁" }}</text>
                <text class="lock-sub"
                  >ID: {{ it.ttlock_lock_id }} · {{ it.lock_mac || "-" }}</text
                >
              </view>
            </view>
            <view class="battery" :class="batteryClass(it.battery)">
              <text class="material-symbols-outlined bat-ic">battery_full</text>
              <text class="bat-txt">{{ normalizeBattery(it.battery) }}%</text>
            </view>
          </view>

          <view class="bind-box">
            <view v-if="it.bind_property_id" class="bind-on">
              <text class="bind-tag on">已绑定</text>
              <view class="bind-info">
                <text class="bind-title"
                  >{{ it.bind_property_title || "房源" }}（{{ it.bind_property_id }}）</text
                >
                <text class="bind-sub">{{
                  it.bind_property_community_name || "-"
                }}</text>
              </view>
              <view class="unbind-btn" @tap.stop="confirmUnbind(it)">
                <text class="material-symbols-outlined">link_off</text>
              </view>
            </view>
            <view v-else class="bind-off">
              <text class="bind-tag off">未绑定</text>
              <text class="bind-tip">点击“绑定房源”完成绑定</text>
            </view>
          </view>

          <view class="card-actions">
            <button class="btn ghost" @tap.stop="openBind(it)">
              <text class="material-symbols-outlined btn-ic">link</text>
              <text>{{ it.bind_property_id ? "更换绑定" : "绑定房源" }}</text>
            </button>
            <button class="btn ghost" @tap.stop="openRecords(it)">
              <text class="material-symbols-outlined btn-ic">history</text>
              <text>开锁记录</text>
            </button>
            <button class="btn primary" @tap.stop="remoteUnlock(it)">
              <text class="material-symbols-outlined btn-ic">lock_open</text>
              <text>远程开锁</text>
            </button>
          </view>
        </view>

        <view v-if="loading" class="loading-row">
          <text class="material-symbols-outlined spin">progress_activity</text>
          <text>加载中...</text>
        </view>
        <view v-else-if="finished && list.length > 0" class="end-row">
          <text>没有更多了</text>
        </view>
      </view>

      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 绑定房源弹窗 -->
    <view class="mask" v-if="bindVisible" @tap="closeBind">
      <view class="sheet" @tap.stop>
        <view class="sheet-head">
          <text class="sheet-title">绑定房源</text>
          <view class="sheet-close" @tap="closeBind">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>
        <view class="sheet-sub">{{ bindLockDisplay }}</view>

        <view class="p-search">
          <text class="material-symbols-outlined p-ic">search</text>
          <input
            v-model="propertyKeyword"
            class="p-input"
            type="text"
            confirm-type="search"
            placeholder="搜索房源标题/小区..."
            placeholder-class="placeholder"
            @confirm="searchProperties(true)"
          />
          <view v-if="propertyKeyword" class="p-clear" @tap="clearPropertyKeyword">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>

        <scroll-view
          scroll-y="true"
          class="p-list"
          lower-threshold="160"
          @scrolltolower="loadMoreProperties"
        >
          <view v-if="propertyLoading && propertyList.length === 0" class="p-loading">
            <text>加载中...</text>
          </view>

          <view v-else-if="!propertyLoading && propertyList.length === 0" class="p-empty">
            <text>暂无房源</text>
          </view>

          <view
            class="p-item"
            v-for="p in propertyList"
            :key="p.id"
            @tap="selectProperty(p)"
          >
            <view class="p-main">
              <text class="p-title">{{ p.title || "-" }}</text>
              <text class="p-sub"
                >{{ p.community_name || "-" }} · {{ p.sale_status_label || "-" }}</text
              >
            </view>
            <text class="material-symbols-outlined p-arrow">chevron_right</text>
          </view>

          <view v-if="propertyFinished && propertyList.length > 0" class="p-end">
            <text>没有更多了</text>
          </view>
          <view style="height: 20rpx"></view>
        </scroll-view>
      </view>
    </view>

    <!-- 开锁记录弹窗 -->
    <view class="mask" v-if="recordsVisible" @tap="closeRecords">
      <view class="sheet records" @tap.stop>
        <view class="sheet-head">
          <text class="sheet-title">开锁记录</text>
          <view class="sheet-close" @tap="closeRecords">
            <text class="material-symbols-outlined">close</text>
          </view>
        </view>
        <view class="sheet-sub">{{ recordsLockDisplay }}</view>

        <scroll-view
          scroll-y="true"
          class="r-list"
          lower-threshold="160"
          @scrolltolower="loadMoreRecords"
        >
          <view v-if="recordsLoading && records.length === 0" class="p-loading">
            <text>加载中...</text>
          </view>
          <view v-else-if="!recordsLoading && records.length === 0" class="p-empty">
            <text>暂无记录</text>
          </view>

          <view class="r-item" v-for="(r, idx) in records" :key="idx">
            <view class="r-row">
              <text class="r-title">{{ recordTitle(r) }}</text>
              <text class="r-time">{{ recordTime(r) }}</text>
            </view>
            <view class="r-sub">{{ recordSub(r) }}</view>
          </view>

          <view v-if="recordsLoading && records.length > 0" class="p-loading">
            <text>加载中...</text>
          </view>
          <view v-else-if="recordsFinished && records.length > 0" class="p-end">
            <text>没有更多了</text>
          </view>
          <view style="height: 20rpx"></view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script>
import TopHeader from "@/components/TopHeader.vue";
import ttlockApi from "@/api/ttlock";
import propertyApi from "@/api/property";
import $store from "@/store";

export default {
  components: { TopHeader },
  data() {
    return {
      permissionChecked: false,
      loading: false,
      refreshing: false,
      finished: false,
      syncing: false,
      page: 1,
      pageSize: 20,
      total: 0,
      list: [],
      keyword: "",
      status: {
        config_ok: false,
        has_token: false,
        api_base: "",
        expire_at: 0,
      },

      // bind modal
      bindVisible: false,
      bindLock: null,
      propertyKeyword: "",
      propertyLoading: false,
      propertyFinished: false,
      propertyPage: 1,
      propertyPageSize: 15,
      propertyList: [],

      // records modal
      recordsVisible: false,
      recordsLock: null,
      recordsLoading: false,
      recordsFinished: false,
      recordsPage: 1,
      recordsPageSize: 20,
      records: [],
    };
  },
  computed: {
    filteredLocks() {
      const kw = String(this.keyword || "").trim().toLowerCase();
      const items = Array.isArray(this.list) ? this.list : [];
      if (!kw) return items;
      return items.filter((it) => {
        const name = String(it && it.lock_name).toLowerCase();
        const mac = String(it && it.lock_mac).toLowerCase();
        const id = String(it && it.ttlock_lock_id);
        return name.includes(kw) || mac.includes(kw) || id.includes(kw);
      });
    },
    bindLockDisplay() {
      const l = this.bindLock || {};
      const name = l.lock_name || "未命名锁";
      const mac = l.lock_mac ? ` · ${l.lock_mac}` : "";
      return `${name}（${l.ttlock_lock_id || "-"}）${mac}`;
    },
    recordsLockDisplay() {
      const l = this.recordsLock || {};
      const name = l.lock_name || "未命名锁";
      const mac = l.lock_mac ? ` · ${l.lock_mac}` : "";
      return `${name}（${l.ttlock_lock_id || "-"}）${mac}`;
    },
  },
  onShow() {
    this.ensurePermissionAndBootstrap();
  },
  methods: {
    goBack() {
      uni.navigateBack();
    },
    async ensurePermissionAndBootstrap() {
      const userStore = $store("user");
      const token = uni.getStorageSync("token");
      if (token && !userStore.isLogin) userStore.setToken(token);
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: "/pages/login/login" });
        return;
      }

      if (!this.permissionChecked) {
        try {
          const info = await userStore.getInfo();
          const u = info || userStore.userInfo || {};
          if (Number(u.can_manage_locks) !== 1) {
            uni.showModal({
              title: "提示",
              content: "暂无智能门锁管理权限，请联系后台管理员开启。",
              showCancel: false,
              success: () => uni.navigateBack(),
            });
            return;
          }
          this.permissionChecked = true;
        } catch (e) {
          uni.showToast({ title: "获取权限失败，请重试", icon: "none" });
          return;
        }
      }

      if (!this.status || typeof this.status.config_ok === "undefined") {
        await this.fetchStatus();
      } else if (!this.status.api_base) {
        await this.fetchStatus();
      }

      if (!Array.isArray(this.list) || this.list.length === 0) {
        await this.fetchList(true);
      }
    },
    async fetchStatus() {
      try {
        const res = await ttlockApi.getStatus();
        if (!res || res.code !== 0) return;
        this.status = res.data || {};
      } catch (e) {}
    },
    handleSearch() {
      // 前端过滤即可，不额外请求
    },
    clearKeyword() {
      this.keyword = "";
    },
    normalizeBattery(v) {
      const n = Number(v);
      if (!isFinite(n) || n < 0) return 0;
      if (n > 100) return 100;
      return Math.round(n);
    },
    batteryClass(v) {
      const n = this.normalizeBattery(v);
      if (n <= 20) return "low";
      if (n <= 50) return "mid";
      return "high";
    },
    async onRefresh() {
      if (this.refreshing || this.loading) return;
      this.refreshing = true;
      await this.fetchStatus();
      await this.fetchList(true);
      this.refreshing = false;
    },
    async loadMore() {
      if (this.loading || this.finished) return;
      this.page += 1;
      await this.fetchList(false);
    },
    async fetchList(reset) {
      if (this.loading) return false;
      this.loading = true;
      try {
        if (reset) {
          this.page = 1;
          this.finished = false;
          this.list = [];
          this.total = 0;
        }
        const res = await ttlockApi.getLockList({
          page: this.page,
          pageSize: this.pageSize,
        });
        if (!res || res.code !== 0) return false;
        const data = res.data || {};
        const items = Array.isArray(data.items) ? data.items : [];
        const total = Number(data.total || 0);
        this.total = total;
        this.list = (Array.isArray(this.list) ? this.list : []).concat(items);
        if (this.list.length >= total || items.length === 0) this.finished = true;
        return true;
      } catch (e) {
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return false;
      } finally {
        this.loading = false;
      }
    },
    async handleSync() {
      if (this.syncing) return;
      this.syncing = true;
      try {
        const res = await ttlockApi.syncLocks();
        if (res && res.code === 0) {
          uni.showToast({ title: "同步完成", icon: "success" });
          await this.onRefresh();
        }
      } catch (e) {
      } finally {
        this.syncing = false;
      }
    },

    // ---- detail / actions ----
    openDetail(lock) {
      // 轻量：详情统一走弹窗/记录/绑定；这里留给后续扩展
      if (!lock) return;
      // 默认打开记录弹窗，避免“点了没反应”的感觉
      this.openRecords(lock);
    },
    async remoteUnlock(lock) {
      const lockId = Number(lock && lock.ttlock_lock_id) || 0;
      if (!lockId) return;
      uni.showModal({
        title: "远程开锁",
        content: "确认执行远程开锁？",
        success: async (r) => {
          if (!r.confirm) return;
          try {
            const resp = await ttlockApi.remoteUnlock({ ttlock_lock_id: lockId });
            if (resp && resp.code === 0) {
              uni.showToast({ title: "已发送开锁指令", icon: "success" });
            }
          } catch (e) {}
        },
      });
    },

    // ---- bind property ----
    openBind(lock) {
      if (!lock || !lock.ttlock_lock_id) return;
      this.bindLock = lock;
      this.bindVisible = true;
      this.propertyKeyword = "";
      this.propertyPage = 1;
      this.propertyFinished = false;
      this.propertyList = [];
      this.searchProperties(true);
    },
    closeBind() {
      this.bindVisible = false;
      this.bindLock = null;
    },
    clearPropertyKeyword() {
      this.propertyKeyword = "";
      this.searchProperties(true);
    },
    async searchProperties(reset) {
      if (this.propertyLoading) return;
      this.propertyLoading = true;
      try {
        if (reset) {
          this.propertyPage = 1;
          this.propertyFinished = false;
          this.propertyList = [];
        }
        const res = await propertyApi.getList({
          page: this.propertyPage,
          pageSize: this.propertyPageSize,
          keyword: String(this.propertyKeyword || "").trim(),
          category: "all",
        });
        if (!res || res.code !== 0) return;
        const data = res.data || {};
        const items = Array.isArray(data.items) ? data.items : [];
        const total = Number(data.total || 0);
        const merged = (Array.isArray(this.propertyList) ? this.propertyList : []).concat(items);
        this.propertyList = merged;
        if (merged.length >= total || items.length === 0) this.propertyFinished = true;
      } catch (e) {
      } finally {
        this.propertyLoading = false;
      }
    },
    async loadMoreProperties() {
      if (this.propertyLoading || this.propertyFinished) return;
      this.propertyPage += 1;
      await this.searchProperties(false);
    },
    async selectProperty(p) {
      const pid = Number(p && p.id) || 0;
      const lockId = Number(this.bindLock && this.bindLock.ttlock_lock_id) || 0;
      if (!pid || !lockId) return;
      try {
        const res = await ttlockApi.bindProperty({ property_id: pid, ttlock_lock_id: lockId });
        if (!res || res.code !== 0) return;
        uni.showToast({ title: "绑定成功", icon: "success" });
        this.closeBind();
        await this.fetchList(true);
      } catch (e) {}
    },
    confirmUnbind(lock) {
      const pid = Number(lock && lock.bind_property_id) || 0;
      if (!pid) return;
      uni.showModal({
        title: "解绑",
        content: "确认解绑该房源与智能锁？",
        success: async (r) => {
          if (!r.confirm) return;
          try {
            const res = await ttlockApi.unbindProperty({ property_id: pid });
            if (!res || res.code !== 0) return;
            uni.showToast({ title: "解绑成功", icon: "success" });
            await this.fetchList(true);
          } catch (e) {}
        },
      });
    },

    // ---- records ----
    openRecords(lock) {
      if (!lock || !lock.ttlock_lock_id) return;
      this.recordsLock = lock;
      this.recordsVisible = true;
      this.recordsPage = 1;
      this.recordsFinished = false;
      this.records = [];
      this.fetchRecords(true);
    },
    closeRecords() {
      this.recordsVisible = false;
      this.recordsLock = null;
      this.records = [];
      this.recordsLoading = false;
      this.recordsFinished = false;
    },
    async fetchRecords(reset) {
      if (this.recordsLoading) return;
      this.recordsLoading = true;
      try {
        if (reset) {
          this.recordsPage = 1;
          this.recordsFinished = false;
          this.records = [];
        }
        const lockId = Number(this.recordsLock && this.recordsLock.ttlock_lock_id) || 0;
        if (!lockId) return;
        const res = await ttlockApi.getLockRecords({
          ttlock_lock_id: lockId,
          page: this.recordsPage,
          pageSize: this.recordsPageSize,
        });
        if (!res || res.code !== 0) return;
        const data = res.data || {};
        const items = Array.isArray(data.list)
          ? data.list
          : Array.isArray(data.items)
            ? data.items
            : Array.isArray(data.records)
              ? data.records
              : [];
        const total = Number(data.total || data.totalCount || data.count || 0);
        const merged = (Array.isArray(this.records) ? this.records : []).concat(items);
        this.records = merged;
        if (total > 0 && merged.length >= total) this.recordsFinished = true;
        if (items.length === 0) this.recordsFinished = true;
      } catch (e) {
      } finally {
        this.recordsLoading = false;
      }
    },
    async loadMoreRecords() {
      if (this.recordsLoading || this.recordsFinished) return;
      this.recordsPage += 1;
      await this.fetchRecords(false);
    },
    recordTitle(r) {
      // 兼容不同字段：recordType/recordTypeName/description
      return (
        String((r && (r.recordTypeName || r.recordType || r.desc || r.description)) || "操作记录") ||
        "操作记录"
      );
    },
    recordTime(r) {
      const t = r && (r.lockDate || r.date || r.time || r.createDate || r.createTime);
      if (!t) return "";
      const n = Number(t);
      if (isFinite(n) && n > 0) {
        const ms = n < 2000000000 ? n * 1000 : n;
        const d = new Date(ms);
        const y = d.getFullYear();
        const m = String(d.getMonth() + 1).padStart(2, "0");
        const dd = String(d.getDate()).padStart(2, "0");
        const hh = String(d.getHours()).padStart(2, "0");
        const mm = String(d.getMinutes()).padStart(2, "0");
        return `${y}-${m}-${dd} ${hh}:${mm}`;
      }
      return String(t);
    },
    recordSub(r) {
      const op = r && (r.username || r.operator || r.openName || r.userName || "");
      const way = r && (r.openWay || r.method || r.openType || "");
      const sn = r && (r.success ? "成功" : r.success === false ? "失败" : "");
      const parts = [op, way, sn].map((s) => String(s || "").trim()).filter(Boolean);
      return parts.join(" · ");
    },
  },
};
</script>

<style>
.lm {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.right-actions {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.icon-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(226, 232, 240, 0.9);
  background: rgba(255, 255, 255, 0.9);
  color: #0f172a;
}

.icon-btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
  box-shadow: 0 12rpx 26rpx rgba(37, 99, 235, 0.2);
}

.icon-btn.ghost {
  background: rgba(255, 255, 255, 0.9);
  color: #0f172a;
}

.icon-btn.disabled {
  opacity: 0.55;
}

.toolbar {
  padding: 14rpx 24rpx 10rpx;
}

.search-box {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 22rpx;
  padding: 0 16rpx;
  height: 76rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
  box-shadow: 0 10rpx 18rpx rgba(15, 23, 42, 0.03);
}

.search-ic {
  font-size: 34rpx;
  color: #94a3b8;
}

.search-input {
  flex: 1;
  height: 76rpx;
  font-size: 26rpx;
  color: #0f172a;
}

.placeholder {
  color: #94a3b8;
}

.clear-btn {
  width: 56rpx;
  height: 56rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
}

.status-row {
  margin-top: 12rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex-wrap: wrap;
}

.chip {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 14rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  font-weight: 700;
  color: #64748b;
  background: rgba(148, 163, 184, 0.14);
  border: 1px solid rgba(148, 163, 184, 0.18);
}

.chip.ok {
  color: #2563eb;
  background: rgba(37, 99, 235, 0.12);
  border: 1px solid rgba(37, 99, 235, 0.18);
}

.chip.ghost {
  color: #475569;
  background: rgba(255, 255, 255, 0.85);
  border: 1px solid rgba(226, 232, 240, 0.9);
}

.chip-ic {
  font-size: 26rpx;
}

.chip-txt {
  max-width: 420rpx;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.list {
  flex: 1;
  overflow: hidden;
}

.list-content {
  padding: 0 24rpx 18rpx;
}

.empty {
  padding: 80rpx 0 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #64748b;
}

.empty-ic {
  font-size: 72rpx;
  color: #cbd5e1;
}

.empty-title {
  margin-top: 14rpx;
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
}

.empty-sub {
  margin-top: 8rpx;
  font-size: 24rpx;
}

.card {
  background: #ffffff;
  border-radius: 28rpx;
  border: 1px solid #f1f5f9;
  box-shadow: 0 10rpx 20rpx rgba(15, 23, 42, 0.03);
  padding: 18rpx;
  margin-top: 14rpx;
}

.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.head-left {
  display: flex;
  align-items: center;
  gap: 14rpx;
  min-width: 0;
}

.lock-ic {
  width: 64rpx;
  height: 64rpx;
  border-radius: 20rpx;
  background: rgba(37, 99, 235, 0.1);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
}

.lock-ic .material-symbols-outlined {
  font-size: 38rpx;
}

.head-texts {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
}

.lock-name {
  font-size: 30rpx;
  font-weight: 900;
  color: #0f172a;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.lock-sub {
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.battery {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 14rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  font-weight: 800;
  color: #2563eb;
  background: rgba(37, 99, 235, 0.12);
  border: 1px solid rgba(37, 99, 235, 0.16);
  flex: 0 0 auto;
}

.battery.mid {
  color: #f97316;
  background: rgba(249, 115, 22, 0.12);
  border: 1px solid rgba(249, 115, 22, 0.18);
}

.battery.low {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.12);
  border: 1px solid rgba(239, 68, 68, 0.18);
}

.battery.high {
  color: #16a34a;
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.18);
}

.bat-ic {
  font-size: 24rpx;
}

.bind-box {
  margin-top: 14rpx;
  padding: 14rpx 14rpx;
  border-radius: 22rpx;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.bind-on {
  display: flex;
  gap: 12rpx;
  align-items: flex-start;
}

.bind-off {
  display: flex;
  gap: 12rpx;
  align-items: center;
}

.bind-tag {
  font-size: 20rpx;
  font-weight: 900;
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
}

.bind-tag.on {
  color: #16a34a;
  background: rgba(34, 197, 94, 0.12);
  border: 1px solid rgba(34, 197, 94, 0.18);
}

.bind-tag.off {
  color: #64748b;
  background: rgba(148, 163, 184, 0.12);
  border: 1px solid rgba(148, 163, 184, 0.18);
}

.bind-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
}

.unbind-btn {
  width: 56rpx;
  height: 56rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.18);
  flex: 0 0 auto;
}

.unbind-btn .material-symbols-outlined {
  font-size: 30rpx;
}

.bind-title {
  font-size: 24rpx;
  font-weight: 800;
  color: #0f172a;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.bind-sub {
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.bind-tip {
  font-size: 22rpx;
  color: #64748b;
}

.card-actions {
  margin-top: 14rpx;
  display: flex;
  gap: 12rpx;
  align-items: center;
}

.btn {
  flex: 1;
  height: 66rpx;
  border-radius: 18rpx;
  font-size: 24rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
}

.btn::after {
  border: none;
}

.btn-ic {
  font-size: 28rpx;
}

.btn.ghost {
  background: #ffffff;
  border: 1px solid rgba(226, 232, 240, 0.9);
  color: #0f172a;
}

.btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  color: #ffffff;
  box-shadow: 0 12rpx 26rpx rgba(37, 99, 235, 0.2);
}

.loading-row,
.end-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  padding: 18rpx 0 10rpx;
  color: #64748b;
  font-size: 24rpx;
}

.spin {
  font-size: 30rpx;
  color: #2563eb;
}

.bottom-spacer {
  height: 30rpx;
}

/* ---------- modal ---------- */
.mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.42);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 999;
}

.sheet {
  width: 100%;
  max-height: 78vh;
  background: #ffffff;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
  overflow: hidden;
  box-shadow: 0 -18rpx 40rpx rgba(15, 23, 42, 0.18);
  display: flex;
  flex-direction: column;
}

.sheet.records {
  max-height: 82vh;
}

.sheet-head {
  padding: 18rpx 18rpx 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(226, 232, 240, 0.9);
}

.sheet-title {
  font-size: 30rpx;
  font-weight: 900;
  color: #0f172a;
}

.sheet-close {
  width: 64rpx;
  height: 64rpx;
  border-radius: 18rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  background: rgba(148, 163, 184, 0.12);
}

.sheet-sub {
  padding: 10rpx 18rpx 0;
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.p-search {
  margin: 14rpx 18rpx 10rpx;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 20rpx;
  height: 72rpx;
  padding: 0 14rpx;
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.p-ic {
  font-size: 30rpx;
  color: #94a3b8;
}

.p-input {
  flex: 1;
  height: 72rpx;
  font-size: 26rpx;
  color: #0f172a;
}

.p-clear {
  width: 52rpx;
  height: 52rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
}

.p-list {
  flex: 1;
  overflow: hidden;
  padding: 0 18rpx 0;
}

.p-item {
  padding: 16rpx 12rpx;
  border-radius: 20rpx;
  border: 1px solid rgba(226, 232, 240, 0.9);
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10rpx;
  margin-bottom: 12rpx;
}

.p-main {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  flex: 1;
  min-width: 0;
}

.p-title {
  font-size: 26rpx;
  font-weight: 900;
  color: #0f172a;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.p-sub {
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.p-arrow {
  font-size: 34rpx;
  color: #cbd5e1;
}

.p-loading,
.p-empty,
.p-end {
  padding: 22rpx 0 12rpx;
  text-align: center;
  color: #64748b;
  font-size: 24rpx;
}

.r-list {
  flex: 1;
  overflow: hidden;
  padding: 10rpx 18rpx 0;
}

.r-item {
  padding: 16rpx 14rpx;
  border-radius: 20rpx;
  border: 1px solid rgba(226, 232, 240, 0.9);
  background: #ffffff;
  margin-bottom: 12rpx;
}

.r-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.r-title {
  font-size: 26rpx;
  font-weight: 900;
  color: #0f172a;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.r-time {
  font-size: 22rpx;
  color: #64748b;
  flex: 0 0 auto;
}

.r-sub {
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #64748b;
}
</style>
