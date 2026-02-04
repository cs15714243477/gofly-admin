<template>
  <view class="page">
    <TopHeader title="经纪人名片" />

    <view class="content">
      <view class="card" :style="{ background: currentStyle.bg }">
        <view class="decor-1"></view>
        <view
          class="decor-2"
          :style="{ backgroundColor: currentStyle.decor2 }"
        ></view>

        <view class="inner">
          <view class="top-row">
            <view class="store">
              <text class="material-symbols-outlined store-icon"
                >apartment</text
              >
              <text class="store-name">{{
                agent.store_name || "未绑定门店"
              }}</text>
            </view>
            <view
              class="badge"
              :style="{ backgroundColor: currentStyle.badgeBg }"
            >
              <text class="material-symbols-outlined badge-icon">verified</text>
              <text>{{ currentStyle.badge }}</text>
            </view>
          </view>

          <view class="mid">
            <image
              v-if="safeImage(agent.avatar)"
              class="avatar"
              :src="safeImage(agent.avatar)"
              mode="aspectFill"
            ></image>
            <view v-else class="avatar avatar-empty">
              <text class="material-symbols-outlined">person</text>
            </view>
            <view class="name">{{ agent.name || "-" }}</view>
            <view class="title">{{ agent.title || "置业顾问" }}</view>
          </view>

          <view class="bottom">
            <view class="line" v-if="agent.store_address">
              <text class="label">门店地址：</text>
              <text class="val">{{ agent.store_address }}</text>
            </view>
            <view class="line" v-if="agent.introduction">
              <text class="label">简介：</text>
              <text class="val">{{ agent.introduction }}</text>
            </view>

            <view class="actions">
              <button
                class="btn primary"
                @click="callAgent"
                :disabled="!agent.mobile"
              >
                <text class="material-symbols-outlined">call</text>
                <text>电话联系</text>
              </button>
              <button
                class="btn ghost"
                @click="copyMobile"
                :disabled="!agent.mobile"
              >
                <text class="material-symbols-outlined">content_copy</text>
                <text>复制电话</text>
              </button>
            </view>
          </view>
        </view>
      </view>

      <view v-if="loading" class="hint">加载中...</view>
      <view v-else-if="!loading && !agent.id" class="hint">暂无经纪人信息</view>

      <view v-if="!loading && agent.id" class="recommend-wrap">
        <view class="recommend-card">
          <view class="sec-title">
            <text class="material-symbols-outlined sec-icon">person</text>
            <text class="sec-text">经纪人推荐</text>
          </view>
          <view
            v-if="agentRecommendedProperties.length === 0"
            class="empty-text"
            >该经纪人暂未设置推荐房源</view
          >
          <view v-else class="list">
            <view
              class="item"
              v-for="item in agentRecommendedProperties"
              :key="'agent-' + item.id"
              @click="openProperty(item)"
            >
              <image
                v-if="safeImage(item.image || item.cover_image)"
                class="cover"
                :src="safeImage(item.image || item.cover_image)"
                mode="aspectFill"
              ></image>
              <view v-else class="cover cover-empty">
                <text class="material-symbols-outlined">image</text>
              </view>
              <view class="meta">
                <view class="title-row">
                  <text class="item-title">{{ item.title || "-" }}</text>
                  <text class="price"
                    >¥{{ item.price || "-" }}{{ item.price_unit || "" }}</text
                  >
                </view>
                <view class="sub">{{ formatPropertyMeta(item) }}</view>
              </view>
            </view>
          </view>
        </view>

        <view class="recommend-card">
          <view class="sec-title">
            <text class="material-symbols-outlined sec-icon">recommend</text>
            <text class="sec-text">系统推荐</text>
          </view>
          <view
            v-if="systemRecommendedProperties.length === 0"
            class="empty-text"
            >暂无系统推荐房源</view
          >
          <view v-else class="list">
            <view
              class="item"
              v-for="item in systemRecommendedProperties"
              :key="'sys-' + item.id"
              @click="openProperty(item)"
            >
              <image
                v-if="safeImage(item.image || item.cover_image)"
                class="cover"
                :src="safeImage(item.image || item.cover_image)"
                mode="aspectFill"
              ></image>
              <view v-else class="cover cover-empty">
                <text class="material-symbols-outlined">image</text>
              </view>
              <view class="meta">
                <view class="title-row">
                  <text class="item-title">{{ item.title || "-" }}</text>
                  <text class="price"
                    >¥{{ item.price || "-" }}{{ item.price_unit || "" }}</text
                  >
                </view>
                <view class="sub">{{ formatPropertyMeta(item) }}</view>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import TopHeader from "@/components/TopHeader.vue";
import userApi from "@/api/user";

function parseScene(scene = "") {
  const s = String(scene || "").trim();
  // 兼容 a1_s0 / a=1&s=0
  if (s.includes("=") && s.includes("&")) {
    const out = {};
    s.split("&").forEach((pair) => {
      const [k, v] = pair.split("=");
      if (k) out[k] = v;
    });
    return {
      agentId: Number(out.a || out.agent_id || 0),
      style: Number(out.s || 0),
    };
  }
  const m = s.match(/a(\d+)_s(\d+)/);
  if (m) return { agentId: Number(m[1] || 0), style: Number(m[2] || 0) };
  const m2 = s.match(/a(\d+)/);
  if (m2) return { agentId: Number(m2[1] || 0), style: 0 };
  return { agentId: 0, style: 0 };
}

export default {
  components: { TopHeader },
  data() {
    return {
      loading: true,
      agent: {},
      agentId: 0,
      styleIndex: 0,
      systemRecommendedProperties: [],
      agentRecommendedProperties: [],
      styles: [
        {
          bg: "linear-gradient(135deg, #2d9cf0, #1e40af)",
          decor2: "rgba(249, 115, 22, 0.2)",
          badge: "PRO",
          badgeBg: "#f97316",
        },
        {
          bg: "linear-gradient(135deg, #06b6d4, #2563eb)",
          decor2: "rgba(255, 255, 255, 0.12)",
          badge: "TOP",
          badgeBg: "#22c55e",
        },
        {
          bg: "linear-gradient(135deg, #111827, #334155)",
          decor2: "rgba(45, 156, 240, 0.18)",
          badge: "VIP",
          badgeBg: "#2d9cf0",
        },
      ],
    };
  },
  computed: {
    currentStyle() {
      return this.styles[this.styleIndex] || this.styles[0];
    },
  },
  onLoad(options) {
    const scene =
      options && options.scene ? decodeURIComponent(options.scene) : "";
    const agentId = Number(options.agent_id || 0);
    const parsed = parseScene(scene);
    this.styleIndex = Math.max(0, Number(options.style || parsed.style || 0));
    this.agentId = agentId || parsed.agentId;
    this.fetchAgent(this.agentId);
  },
  methods: {
    safeImage(url) {
      const u = String(url || "").trim();
      if (!u) return "";
      if (u.indexOf("/static/images/") === 0) return "";
      return u;
    },
    formatPropertyMeta(item) {
      const area = item && item.area ? `${item.area}㎡` : "";
      const rooms = Number(item && item.rooms) || 0;
      const halls = Number(item && item.halls) || 0;
      const bathrooms = Number(item && item.bathrooms) || 0;
      let layout = "";
      if (rooms || halls || bathrooms)
        layout = `${rooms}室${halls}厅${bathrooms}卫`;
      const community = String((item && item.community_name) || "").trim();
      return (
        [community, layout, area].filter(Boolean).join(" · ") || "基础信息"
      );
    },
    async fetchAgent(agentId) {
      this.loading = true;
      try {
        if (!agentId) {
          this.agent = {};
          this.systemRecommendedProperties = [];
          this.agentRecommendedProperties = [];
          return;
        }
        const res = await userApi.getAgentCard({ agent_id: agentId }, true);
        if (!res || res.code !== 0) {
          this.agent = {};
          this.systemRecommendedProperties = [];
          this.agentRecommendedProperties = [];
          return;
        }
        const data = res.data || {};
        this.agent = data;
        this.systemRecommendedProperties = Array.isArray(
          data.system_recommended_properties,
        )
          ? data.system_recommended_properties
          : [];
        this.agentRecommendedProperties = Array.isArray(
          data.agent_recommended_properties,
        )
          ? data.agent_recommended_properties
          : [];
      } finally {
        this.loading = false;
      }
    },
    openProperty(item) {
      const id = Number(item && (item.id || item.ID)) || 0;
      if (!id) return;
      const agentID = Number(
        this.agentId || (this.agent && this.agent.id) || 0,
      );
      const style = Number(this.styleIndex || 0);
      uni.navigateTo({
        url: `/pages/property_detail/property_detail?id=${encodeURIComponent(id)}&public=1&from_agent_id=${encodeURIComponent(agentID)}&from_style=${encodeURIComponent(style)}`,
      });
    },
    callAgent() {
      const mobile = String((this.agent && this.agent.mobile) || "").trim();
      if (!mobile) return;
      uni.makePhoneCall({ phoneNumber: mobile });
    },
    copyMobile() {
      const mobile = String((this.agent && this.agent.mobile) || "").trim();
      if (!mobile) return;
      uni.setClipboardData({
        data: mobile,
        success: () => uni.showToast({ title: "已复制", icon: "none" }),
      });
    },
  },
};
</script>

<style lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7f8;
}

.content {
  padding: 24rpx;
}

.card {
  position: relative;
  border-radius: 40rpx;
  overflow: hidden;
  box-shadow: 0 40rpx 80rpx rgba(0, 0, 0, 0.18);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.decor-1 {
  position: absolute;
  top: -80rpx;
  right: -80rpx;
  width: 256rpx;
  height: 256rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

.decor-2 {
  position: absolute;
  bottom: -90rpx;
  left: -90rpx;
  width: 320rpx;
  height: 320rpx;
  border-radius: 50%;
  background: rgba(249, 115, 22, 0.2);
}

.inner {
  position: relative;
  z-index: 2;
  padding: 40rpx;
  color: #fff;
}

.top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.store {
  display: flex;
  align-items: center;
  gap: 12rpx;
  .store-icon {
    font-size: 34rpx;
    opacity: 0.9;
  }
  .store-name {
    font-size: 28rpx;
    font-weight: 800;
    opacity: 0.95;
  }
}

.badge {
  padding: 6rpx 16rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 22rpx;
  font-weight: 800;
  .badge-icon {
    font-size: 26rpx;
  }
}

.mid {
  padding: 50rpx 0 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  .avatar {
    width: 180rpx;
    height: 180rpx;
    border-radius: 50%;
    border: 6rpx solid rgba(255, 255, 255, 0.25);
    background: rgba(255, 255, 255, 0.06);
  }
  .avatar-empty {
    display: flex;
    align-items: center;
    justify-content: center;
    .material-symbols-outlined {
      font-size: 76rpx;
      color: rgba(255, 255, 255, 0.85);
    }
  }
  .name {
    font-size: 44rpx;
    font-weight: 900;
  }
  .title {
    font-size: 28rpx;
    font-weight: 600;
    opacity: 0.86;
  }
}

.bottom {
  background: rgba(255, 255, 255, 0.1);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  border-radius: 28rpx;
  padding: 22rpx 22rpx;
  .line {
    display: flex;
    gap: 8rpx;
    line-height: 1.5;
    margin-bottom: 12rpx;
  }
  .label {
    font-size: 24rpx;
    opacity: 0.78;
  }
  .val {
    font-size: 26rpx;
    font-weight: 600;
    opacity: 0.92;
  }
}

.actions {
  margin-top: 8rpx;
  display: flex;
  gap: 16rpx;
  .btn {
    flex: 1;
    height: 84rpx;
    border-radius: 22rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10rpx;
    font-size: 30rpx;
    font-weight: 800;
    border: none;
    &::after {
      border: none;
    }
  }
  .primary {
    background: rgba(255, 255, 255, 0.95);
    color: #0f172a;
  }
  .ghost {
    background: rgba(255, 255, 255, 0.14);
    color: #ffffff;
  }
}

.hint {
  padding: 28rpx 0;
  text-align: center;
  color: #64748b;
  font-size: 26rpx;
}

.recommend-wrap {
  margin-top: 24rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.recommend-card {
  background: #ffffff;
  border-radius: 24rpx;
  padding: 22rpx;
  border: 1rpx solid #e2e8f0;
  box-shadow: 0 8rpx 24rpx rgba(15, 23, 42, 0.05);
}

.sec-title {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
  line-height: 1;
  margin-bottom: 14rpx;
  .sec-icon {
    font-size: 34rpx;
    color: #2d9cf0;
    line-height: 1;
    flex-shrink: 0;
  }
  .sec-text {
    line-height: 1.1;
  }
}

.empty-text {
  font-size: 24rpx;
  color: #94a3b8;
  padding: 16rpx 0;
}

.list {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.item {
  display: flex;
  align-items: center;
  gap: 14rpx;
  padding: 8rpx 0;
}

.cover {
  width: 172rpx;
  height: 118rpx;
  border-radius: 14rpx;
  background: #f1f5f9;
  flex-shrink: 0;
}
.cover-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  .material-symbols-outlined {
    font-size: 40rpx;
    color: #94a3b8;
  }
}

.meta {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.title-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12rpx;
}

.item-title {
  flex: 1;
  font-size: 28rpx;
  font-weight: 800;
  color: #0f172a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  flex-shrink: 0;
  font-size: 24rpx;
  font-weight: 800;
  color: #f97316;
}

.sub {
  font-size: 22rpx;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
