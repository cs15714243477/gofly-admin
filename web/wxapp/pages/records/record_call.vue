<template>
  <view class="record-page">
    <TopHeader title="通话记录" />
    <scroll-view scroll-y class="list">
      <view v-if="loading && items.length === 0" class="hint">加载中...</view>
      <view v-else-if="!loading && items.length === 0" class="hint"
        >暂无通话记录</view
      >
      <view
        class="card"
        v-for="(it, idx) in items"
        :key="idx"
        @click="callItem(it)"
      >
        <view class="row">
          <text class="title">{{ it.name || it.title || "通话记录" }}</text>
          <text class="tag" :class="it.type === '呼入' ? 'in' : 'out'">{{
            it.type || "呼出"
          }}</text>
        </view>
        <view class="sub"
          >{{ it.time || "-" }} · {{ it.phone || "暂无号码" }}</view
        >
      </view>
      <view class="bottom-spacer"></view>
    </scroll-view>
    <BottomTabBar active="me" />
  </view>
</template>

<script>
import TopHeader from "@/components/TopHeader.vue";
import BottomTabBar from "@/components/BottomTabBar.vue";
import userApi from "@/api/user";

export default {
  components: { TopHeader, BottomTabBar },
  data() {
    return {
      loading: false,
      items: [],
    };
  },
  onShow() {
    this.loadData();
  },
  methods: {
    async loadData() {
      this.loading = true;
      try {
        const res = await userApi.getWorkbenchRecords(
          { record_type: "call", page: 1, page_size: 50 },
          true,
        );
        const data = (res && res.code === 0 && res.data) || {};
        this.items = Array.isArray(data.items) ? data.items : [];
      } finally {
        this.loading = false;
      }
    },
    callItem(it) {
      const phone = String((it && it.phone) || "").trim();
      if (!phone) {
        uni.showToast({ title: "暂无可拨打号码", icon: "none" });
        return;
      }
      uni.makePhoneCall({ phoneNumber: phone });
    },
  },
};
</script>

<style lang="scss">
.record-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f6f7f8;
}
.list {
  flex: 1;
  overflow: hidden;
  padding: 24rpx 16rpx;
  box-sizing: border-box;
}
.card {
  background: #fff;
  border-radius: 24rpx;
  padding: 20rpx;
  border: 1px solid #f1f5f9;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
  margin-bottom: 16rpx;
}
.row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
}
.title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tag {
  font-size: 22rpx;
  font-weight: 900;
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
}
.tag.out {
  color: #2563eb;
  background: rgba(37, 99, 235, 0.12);
  border: 1rpx solid rgba(37, 99, 235, 0.18);
}
.tag.in {
  color: #16a34a;
  background: rgba(34, 197, 94, 0.12);
  border: 1rpx solid rgba(34, 197, 94, 0.18);
}
.sub {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #64748b;
}
.hint {
  padding: 50rpx 0;
  text-align: center;
  font-size: 24rpx;
  color: #94a3b8;
}
.bottom-spacer {
  height: 40rpx;
}
</style>
