<template>
  <view class="workbench-container">
    <!-- 顶部导航 -->
    <TopHeader title="我的" />

    <!-- 名片固定在上方，下面内容允许滚动 -->
    <view class="main-content">
      <!-- 个人信息卡片 -->
      <view class="profile-section">
        <view class="profile-card">
          <view class="profile-bg-decor"></view>
          <button class="edit-btn" @click="goEditCard">
            <text class="material-symbols-outlined edit-icon">edit</text>
            <text>编辑</text>
          </button>
          <view class="profile-main">
            <view class="avatar-box">
              <image
                v-if="avatarUrl"
                class="avatar"
                :src="avatarUrl"
                mode="aspectFill"
              ></image>
              <view v-else class="avatar avatar-empty">
                <text class="material-symbols-outlined">person</text>
              </view>
              <view class="online-status"></view>
            </view>
            <view class="user-info">
              <text class="user-name">{{ displayName }}</text>
              <text class="user-role">{{ displayRoleLine }}</text>
              <text class="user-phone">{{ displayMobile }}</text>
            </view>
          </view>
        </view>
      </view>

      <scroll-view scroll-y="true" class="content-scroll">
        <view class="scroll-content">
          <view class="section">
            <view class="section-title">业务记录</view>
            <view class="records-list">
              <view
                class="property-card"
                v-for="(item, index) in displayRecords"
                :key="index"
                @click="openRecord(item)"
              >
                <view class="card-main">
                  <view
                    class="record-image-box"
                    :class="{ 'has-notice': item.hasNotice }"
                  >
                    <view class="record-icon-chip">
                      <text class="material-symbols-outlined record-icon">{{
                        item.icon
                      }}</text>
                    </view>
                    <view class="record-dot" v-if="item.hasNotice"></view>
                  </view>
                  <view class="record-info-box">
                    <view class="title">{{ item.name }}</view>
                    <view class="meta">
                      <text class="bold">{{ getRecordHint(item) }}</text>
                    </view>
                    <view class="stats">
                      <text>{{ item.countLabel || "暂无统计" }}</text>
                      <text class="stats-sep"></text>
                      <text>{{
                        item.hasNotice ? "有新动态" : "点击查看详情"
                      }}</text>
                    </view>
                  </view>
                </view>
                <view
                  class="card-footer"
                  :class="item.hasNotice ? 'orange-footer' : 'grey-footer'"
                >
                  <view class="footer-left">
                    <view class="footer-icon-box">
                      <text class="material-symbols-outlined footer-icon">{{
                        item.hasNotice ? "notifications_active" : "insights"
                      }}</text>
                    </view>
                    <text class="footer-text">{{ getRecordFooter(item) }}</text>
                  </view>
                  <view class="footer-btn">
                    <text class="material-symbols-outlined btn-icon"
                      >arrow_forward</text
                    >
                    <text>进入</text>
                  </view>
                </view>
              </view>
            </view>
          </view>

          <!-- 更多服务 -->
          <view class="section">
            <view class="section-title">更多服务</view>
            <view class="service-list">
              <view class="service-item" @click="openAbout">
                <view class="service-left">
                  <text class="material-symbols-outlined service-icon">info</text>
                  <text class="service-name">关于我们</text>
                </view>
                <text class="material-symbols-outlined arrow-icon"
                  >chevron_right</text
                >
              </view>
              <view class="service-divider"></view>
              <button
                class="service-item service-share-btn"
                open-type="share"
                @longpress="copyAgentLink"
              >
                <view class="service-left">
                  <text class="material-symbols-outlined service-icon"
                    >share</text
                  >
                  <text class="service-name">推荐给朋友</text>
                </view>
                <text class="material-symbols-outlined arrow-icon"
                  >chevron_right</text
                >
              </button>
            </view>
          </view>

          <!-- 退出按钮 -->
          <view class="logout-section">
            <button class="logout-btn" @click="handleLogout">退出登录</button>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 底部导航 -->
    <BottomTabBar active="me" />
  </view>
</template>

<script>
import BottomTabBar from "@/components/BottomTabBar.vue";
import TopHeader from "@/components/TopHeader.vue";
import userApi from "@/api/user";
import $store from "@/store";

export default {
  components: { BottomTabBar, TopHeader },
  onShareAppMessage() {
    const name = this.displayName || "我";
    const agentId = Number((this.userInfo && this.userInfo.id) || 0);
    return {
      title: `${name}的名片`,
      path: `/pages/agent_public_card/agent_public_card?agent_id=${encodeURIComponent(agentId)}&style=0`,
    };
  },
  data() {
    return {
      loadingUser: false,
      debugLogged: false,
      userInfo: {},
      recordSummary: {
        follow_count: 0,
        unlock_count: 0,
        showing_count: 0,
        view_count: 0,
        share_count: 0,
        call_count: 0,
        unlock_has_notice: false,
      },
      businessRecords: [
        { key: "property_manage", name: "房源管理", icon: "home_work" },
        { key: "lock_manage", name: "智能门锁", icon: "lock" },
        { key: "follow", name: "关注记录", icon: "favorite" },
        { key: "unlock", name: "开锁记录", icon: "lock_open", hasNotice: true },
        { key: "showing", name: "带看记录", icon: "location_on" },
        { key: "view", name: "浏览记录", icon: "history" },
        { key: "share", name: "分享记录", icon: "share" },
        { key: "call", name: "通话记录", icon: "call" },
      ],
    };
  },
  computed: {
    avatarUrl() {
      const avatar = String(
        (this.userInfo && this.userInfo.avatar) || "",
      ).trim();
      if (!avatar) return "";
      if (avatar.indexOf("/static/images/") === 0) return "";
      return avatar;
    },
    displayName() {
      const u = this.userInfo || {};
      // 真实姓名优先，其次昵称/用户名
      return u.name || u.nickname || u.username || "未登录";
    },
    displayRoleLine() {
      const u = this.userInfo || {};
      const title = (u.title || "").trim();
      const roleRaw = String(u.role || "").trim();
      // 门店信息由后端返回：store_name（未绑定则为“未绑定”）
      const storeName = (u.store_name || "").trim();
      const store = storeName ? storeName : "未绑定";
      // 角色/身份展示：优先头衔；否则把 role=1/user 映射为“经纪人”
      let left = title;
      if (!left) {
        if (
          roleRaw === "" ||
          roleRaw === "1" ||
          roleRaw.toLowerCase() === "user"
        )
          left = "经纪人";
        else left = roleRaw;
      }
      return `${left} | ${store}`;
    },
    displayMobile() {
      const u = this.userInfo || {};
      return u.mobile || "";
    },
    displayRecords() {
      const u = this.userInfo || {};
      const canManage = Number(u.can_manage_properties) === 1;
      const canManageLocks = Number(u.can_manage_locks) === 1;
      const summary = this.recordSummary || {};
      const countMap = {
        follow: Number(summary.follow_count || 0),
        unlock: Number(summary.unlock_count || 0),
        showing: Number(summary.showing_count || 0),
        view: Number(summary.view_count || 0),
        share: Number(summary.share_count || 0),
        call: Number(summary.call_count || 0),
      };
      // 最小权限：仅在允许维护房源时展示入口
      return (this.businessRecords || [])
        .filter((it) => {
          if (!it) return false;
          if (it.key === "property_manage") return canManage;
          if (it.key === "lock_manage") return canManageLocks;
          return true;
        })
        .map((it) => {
          const count = countMap[it.key] || 0;
          return {
            ...it,
            hasNotice:
              it.key === "unlock"
                ? !!summary.unlock_has_notice
                : !!it.hasNotice,
            countLabel: count > 0 ? `${count}条` : "",
          };
        });
    },
  },
  onShow() {
    this.ensureLoginAndLoadUser();
  },
  methods: {
    debugPrintUserInfo(tag = "") {
      if (this.debugLogged) return;
      this.debugLogged = true;
      try {
        // 关键字段 + 全量
        console.log(
          "[agent_workbench_home] " + (tag || "userInfo") + " store_name=",
          this.userInfo?.store_name,
          "store_id=",
          this.userInfo?.store_id,
          "title=",
          this.userInfo?.title,
          "role=",
          this.userInfo?.role,
        );
        console.log(
          "[agent_workbench_home] " + (tag || "userInfo") + " full=",
          JSON.parse(JSON.stringify(this.userInfo || {})),
        );
      } catch (e) {
        console.log("[agent_workbench_home] debugPrintUserInfo error", e);
      }
    },
    debugShowUserInfo() {
      try {
        const u = this.userInfo || {};
        const contentLines = [
          `name: ${u.name || ""}`,
          `nickname: ${u.nickname || ""}`,
          `username: ${u.username || ""}`,
          `title: ${u.title || ""}`,
          `role: ${u.role || ""}`,
          `store_name: ${u.store_name || ""}`,
          `store_id: ${u.store_id || ""}`,
          `store_address: ${u.store_address || ""}`,
        ];
        uni.showModal({
          title: "调试：userInfo(点击头像)",
          content: contentLines.join("\n"),
          showCancel: false,
        });
      } catch (e) {
        uni.showToast({ title: "调试弹窗失败", icon: "none" });
      }
    },
    async ensureLoginAndLoadUser() {
      const userStore = $store("user");
      // 兼容：小程序刷新后，优先用本地 token 恢复登录态
      const token = uni.getStorageSync("token");
      if (token && !userStore.isLogin) {
        userStore.setToken(token);
      }
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: "/pages/login/login" });
        return;
      }
      if (this.loadingUser) return;
      this.loadingUser = true;
      try {
        const info = await userStore.getInfo();
        this.userInfo = info || userStore.userInfo || {};
        this.debugPrintUserInfo("after getInfo");
        await this.loadWorkbenchSummary(false);
      } catch (e) {
        // 请求失败：优先保留本地态，不强制跳转（避免短暂网络抖动导致回登录）
        this.userInfo = userStore.userInfo || {};
        this.debugPrintUserInfo("fallback userStore.userInfo");
        await this.loadWorkbenchSummary(false);
      } finally {
        this.loadingUser = false;
      }
    },
    async loadWorkbenchSummary(showLoading = false) {
      try {
        const res = await userApi.getWorkbenchSummary(showLoading);
        if (!res || res.code !== 0 || !res.data) return;
        this.recordSummary = {
          follow_count: Number(res.data.follow_count || 0),
          unlock_count: Number(res.data.unlock_count || 0),
          showing_count: Number(res.data.showing_count || 0),
          view_count: Number(res.data.view_count || 0),
          share_count: Number(res.data.share_count || 0),
          call_count: Number(res.data.call_count || 0),
          unlock_has_notice: !!res.data.unlock_has_notice,
        };
      } catch (e) {}
    },
    goEditCard() {
      // 跳转到“获客-编辑资料”页（tab=1）
      uni.reLaunch({
        url: "/pages/my_business_card/my_business_card?tab=1",
      });
    },
    openAbout() {
      uni.navigateTo({ url: "/pages/doc_webview/doc_webview?key=help_center" });
    },
    async copyAgentLink() {
      const res = await userApi.getAgentUrlLink({}, true);
      const urlLink =
        res && res.code === 0 && res.data
          ? String(res.data.url_link || "").trim()
          : "";
      if (!urlLink) {
        uni.showToast({ title: "获取链接失败", icon: "none" });
        return;
      }
      uni.setClipboardData({
        data: urlLink,
        success: () => uni.showToast({ title: "已复制链接", icon: "none" }),
      });
    },
    openRecord(item) {
      const map = {
        property_manage: "/pages/property_manage/property_manage",
        lock_manage: "/pages/lock_manage/lock_manage",
        follow: "/pages/records/record_follow",
        unlock: "/pages/records/record_unlock",
        showing: "/pages/records/record_showing",
        view: "/pages/records/record_view",
        share: "/pages/records/record_share",
        call: "/pages/records/record_call",
      };
      const url = item && item.key && map[item.key];
      if (!url) return;
      uni.navigateTo({ url });
    },
    getRecordHint(item) {
      const key = String((item && item.key) || "").trim();
      const map = {
        property_manage: "维护发布、编辑、上下架",
        lock_manage: "远程开锁、门锁状态管理",
        follow: "客户关注房源动态",
        unlock: "开锁日志与异常提醒",
        showing: "带看安排与到访回溯",
        view: "浏览来源与热度追踪",
        share: "分享触达与传播记录",
        call: "客户来电与沟通统计",
      };
      return map[key] || "点击进入查看明细";
    },
    getRecordFooter(item) {
      const label = String((item && item.countLabel) || "").trim();
      if (label) return `累计 ${label}`;
      return "点击进入查看明细";
    },
    handleLogout() {
      uni.showModal({
        title: "提示",
        content: "确定要退出登录吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              await $store("user").logout(false);
            } catch (e) {}
            uni.reLaunch({
              url: "/pages/login/login",
            });
          }
        },
      });
    },
  },
};
</script>

<style>
.workbench-container {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* header/title 已由 TopHeader 统一 */

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding: 12rpx 24rpx 0;
  box-sizing: border-box;
  min-height: 0;
}

.profile-section {
  padding: 0;
  flex-shrink: 0;
}

.content-scroll {
  flex: 1;
  min-height: 0;
}

.scroll-content {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
  padding-bottom: calc(env(safe-area-inset-bottom) + 20rpx);
}

.profile-card {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border-radius: 32rpx;
  padding: 36rpx 32rpx;
  min-height: 240rpx;
  position: relative;
  overflow: hidden;
  box-shadow: 0 12rpx 30rpx rgba(37, 99, 235, 0.2);
  display: flex;
  align-items: center;
}

.profile-bg-decor {
  position: absolute;
  right: -90rpx;
  top: -90rpx;
  width: 360rpx;
  height: 360rpx;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  filter: blur(90rpx);
}

.edit-btn {
  position: absolute;
  top: 28rpx;
  right: 28rpx;
  height: 56rpx;
  padding: 0 20rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  color: #ffffff;
  font-size: 24rpx;
  z-index: 10;
}

.edit-icon {
  font-size: 24rpx !important;
  margin-right: 4rpx;
}

.profile-main {
  display: flex;
  align-items: center;
  gap: 24rpx;
  position: relative;
  z-index: 5;
}

.avatar-box {
  position: relative;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.3);
}
.avatar-empty {
  background: rgba(255, 255, 255, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  .material-symbols-outlined {
    font-size: 58rpx;
    color: rgba(255, 255, 255, 0.9);
  }
}

.online-status {
  position: absolute;
  bottom: 4rpx;
  right: 4rpx;
  width: 20rpx;
  height: 20rpx;
  background-color: #4ade80;
  border-radius: 50%;
  border: 4rpx solid #2d9cf0;
}

.user-info {
  color: #ffffff;
}

.user-name {
  font-size: 38rpx;
  font-weight: bold;
  display: block;
  letter-spacing: 0.4rpx;
}

.user-role {
  font-size: 26rpx;
  opacity: 0.9;
  display: block;
  margin-top: 6rpx;
}

.user-phone {
  font-size: 26rpx;
  opacity: 0.8;
  display: block;
  margin-top: 4rpx;
}

.section {
  padding: 0;
}

.section-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #0f172a;
  margin-bottom: 12rpx;
}

.records-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.property-card {
  background-color: #ffffff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
  border: 1px solid #f1f5f9;
}

.card-main {
  padding: 20rpx;
  display: flex;
  gap: 18rpx;
}

.record-image-box {
  width: 224rpx;
  height: 192rpx;
  border-radius: 16rpx;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
}

.record-image-box.has-notice {
  background: linear-gradient(135deg, #fb923c, #f97316);
}

.record-icon-chip {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  border: 2rpx solid rgba(255, 255, 255, 0.28);
  display: flex;
  align-items: center;
  justify-content: center;
}

.record-info-box {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.record-info-box .title {
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

.record-info-box .meta {
  display: flex;
  align-items: center;
  gap: 12rpx;
  font-size: 24rpx;
  color: #64748b;
}

.record-info-box .meta .bold {
  font-weight: 500;
  color: #334155;
}

.record-info-box .stats {
  font-size: 20rpx;
  color: #94a3b8;
  display: flex;
  align-items: center;
  margin-top: auto;
}

.record-info-box .stats .stats-sep {
  width: 1px;
  height: 16rpx;
  background-color: #e2e8f0;
  margin: 0 12rpx;
}

.record-icon {
  font-size: 56rpx !important;
  color: #ffffff;
  line-height: 1;
  display: block;
}

.record-dot {
  position: absolute;
  top: 14rpx;
  right: 14rpx;
  width: 16rpx;
  height: 16rpx;
  background-color: #ef4444;
  border-radius: 50%;
  border: 3rpx solid #ffffff;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14rpx 20rpx;
  border-top: 1rpx solid transparent;
}

.card-footer.orange-footer {
  background-color: #fff7ed;
  border-color: #ffedd5;
}

.card-footer.orange-footer .footer-icon-box {
  background-color: #f97316;
}

.card-footer.orange-footer .footer-text {
  color: #9a3412;
  font-weight: bold;
}

.card-footer.grey-footer {
  background-color: #f8fafc;
  border-color: #f1f5f9;
}

.card-footer.grey-footer .footer-icon-box {
  background-color: #cbd5e1;
}

.card-footer.grey-footer .footer-text {
  color: #64748b;
  font-weight: 500;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
  flex: 1;
  min-width: 0;
}

.footer-icon-box {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.footer-icon {
  color: #ffffff;
  font-size: 28rpx;
}

.footer-text {
  font-size: 24rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
}

.footer-btn .btn-icon {
  font-size: 28rpx;
}

.service-list {
  background-color: #ffffff;
  border-radius: 24rpx;
  border: 1px solid #f1f5f9;
}

.service-item {
  padding: 18rpx 24rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.service-share-btn {
  width: 100%;
  border: none;
  background: #ffffff;
  text-align: left;
}
.service-share-btn::after {
  border: none;
}

.service-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.service-icon {
  font-size: 40rpx !important;
  color: #94a3b8;
}

.service-name {
  font-size: 28rpx;
  color: #334155;
}

.arrow-icon {
  font-size: 32rpx !important;
  color: #cbd5e1;
}

.service-divider {
  height: 1px;
  background-color: #f1f5f9;
  margin: 0 32rpx;
}

.logout-section {
  padding: 0;
}

.logout-btn {
  width: 100%;
  height: 80rpx;
  background-color: #fff1f2;
  color: #ef4444;
  font-size: 28rpx;
  font-weight: bold;
  border-radius: 20rpx;
  border: 1px solid #fee2e2;
}

.tab-bar {
  background-color: #ffffff;
  border-top: 1rpx solid #e2e8f0;
  padding: 12rpx 0 calc(env(safe-area-inset-bottom) + 12rpx);
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.tab-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  color: #94a3b8;
  flex: 1;
}

.tab-icon {
  font-size: 48rpx !important;
}

.tab-text {
  font-size: 20rpx;
  margin-top: 4rpx;
}

.tab-item.active {
  color: #2d9cf0;
}

.active-bg {
  background-color: #eef7ff;
  padding: 8rpx 32rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fill {
  font-variation-settings: "FILL" 1;
}

/* 页面改为“名片固定 + 内容滚动” */
</style>
