<template>
  <view class="workbench-container">
    <!-- 顶部导航 -->
    <TopHeader title="我的" />

    <!-- 名片固定在上方，下面内容允许滚动 -->
    <view class="main-content">
      <!-- 个人信息卡片 -->
      <view class="profile-section">
        <view v-if="!isGuest" class="profile-card">
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
        <view v-else class="guest-card">
          <view class="guest-title">当然未登录</view>
          <view class="guest-desc">你可以先浏览推荐、房源等服务，需要时再自主登录。</view>
          <view class="guest-actions">
            <button class="guest-login-btn" @click="goLogin">去登录</button>
            <button class="guest-back-btn" @click="goHome">返回首页</button>
          </view>
        </view>
      </view>

      <scroll-view scroll-y="true" class="content-scroll">
        <view class="scroll-content">
          <view class="section">
            <view class="section-title">快捷工作台</view>
            <view class="records-grid" :class="{ 'is-guest': isGuest }">
              <view
                class="record-grid-item"
                v-for="(item, index) in displayRecords"
                :key="index"
                @click="openRecord(item)"
              >
                <view class="record-grid-icon-wrap">
                  <text class="material-symbols-outlined record-grid-icon">{{
                    item.icon
                  }}</text>
                  <view class="record-grid-dot" v-if="item.hasNotice"></view>
                </view>
                <text class="record-grid-name">{{ item.name }}</text>
              </view>
            </view>
            <view v-if="isGuest" class="guest-hint">
              工作台功能需登录后使用，你也可以先继续浏览其他页面。
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
                v-if="!isGuest"
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
              <view v-else class="service-item" @click="promptLogin('登录后可推荐给朋友')">
                <view class="service-left">
                  <text class="material-symbols-outlined service-icon"
                    >share</text
                  >
                  <text class="service-name">推荐给朋友（登录后）</text>
                </view>
                <text class="material-symbols-outlined arrow-icon"
                  >chevron_right</text
                >
              </view>
            </view>
          </view>

          <!-- 退出按钮 -->
          <view class="logout-section" v-if="!isGuest">
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
    if (this.isGuest) {
      return {
        title: "快销房智选",
        path: "/pages/home/home",
      };
    }
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
      isGuest: true,
      lastLoginPromptAt: 0,
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
    setGuestState() {
      this.isGuest = true;
      this.userInfo = {
        name: "游客",
        title: "",
        role: "",
        mobile: "",
        store_name: "",
      };
      this.recordSummary = {
        follow_count: 0,
        unlock_count: 0,
        showing_count: 0,
        view_count: 0,
        share_count: 0,
        call_count: 0,
        unlock_has_notice: false,
      };
    },
    goLogin() {
      uni.navigateTo({ url: "/pages/login/login" });
    },
    goHome() {
      uni.reLaunch({ url: "/pages/home/home" });
    },
    promptLogin(content = "登录后可使用该功能") {
      const now = Date.now();
      if (now - Number(this.lastLoginPromptAt || 0) < 1200) return;
      this.lastLoginPromptAt = now;
      uni.showModal({
        title: "需要登录",
        content,
        cancelText: "继续浏览",
        confirmText: "去登录",
        success: (res) => {
          if (res && res.confirm) this.goLogin();
        },
      });
    },
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
        this.setGuestState();
        return;
      }
      if (this.loadingUser) return;
      this.loadingUser = true;
      try {
        const info = await userStore.getInfo();
        this.userInfo = info || userStore.userInfo || {};
        this.isGuest = !(
          this.userInfo && Number(this.userInfo.id || 0) > 0
        );
        if (this.isGuest) {
          this.setGuestState();
          return;
        }
        this.debugPrintUserInfo("after getInfo");
        await this.loadWorkbenchSummary(false);
      } catch (e) {
        // 请求失败：优先保留本地态，不强制跳转（避免短暂网络抖动导致回登录）
        this.userInfo = userStore.userInfo || {};
        this.isGuest = !(
          this.userInfo && Number(this.userInfo.id || 0) > 0
        );
        if (this.isGuest) {
          this.setGuestState();
          return;
        }
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
      if (this.isGuest) {
        this.promptLogin("登录后可编辑你的名片资料");
        return;
      }
      // 跳转到“获客-编辑资料”页（tab=1）
      uni.reLaunch({
        url: "/pages/my_business_card/my_business_card?tab=1",
      });
    },
    openAbout() {
      uni.navigateTo({ url: "/pages/doc_webview/doc_webview?key=about_us" });
    },
    async copyAgentLink() {
      if (this.isGuest) {
        this.promptLogin("登录后可复制你的专属推广链接");
        return;
      }
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
      if (this.isGuest) {
        this.promptLogin("登录后可查看工作台记录");
        return;
      }
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
      if (this.isGuest) {
        this.goLogin();
        return;
      }
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

.guest-card {
  background: linear-gradient(145deg, #f8fbff, #eef6ff);
  border: 1px solid #dbeafe;
  border-radius: 28rpx;
  padding: 28rpx 24rpx;
  box-shadow: 0 8rpx 20rpx rgba(37, 99, 235, 0.08);
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.guest-title {
  font-size: 34rpx;
  font-weight: 800;
  color: #0f172a;
}

.guest-desc {
  font-size: 24rpx;
  color: #475569;
  line-height: 1.6;
}

.guest-actions {
  display: flex;
  gap: 14rpx;
}

.guest-login-btn,
.guest-back-btn {
  flex: 1;
  height: 72rpx;
  border-radius: 14rpx;
  font-size: 26rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

.guest-login-btn {
  background: #2563eb;
  color: #ffffff;
  border: none;
}

.guest-login-btn::after {
  border: none;
}

.guest-back-btn {
  background: #ffffff;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.guest-back-btn::after {
  border: none;
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

.records-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0;
  background: #ffffff;
  border: 1px solid #f1f5f9;
  border-radius: 24rpx;
  overflow: hidden;
}

.records-grid.is-guest {
  opacity: 0.92;
}

.record-grid-item {
  background: transparent;
  border: none;
  border-radius: 0;
  padding: 18rpx 8rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  min-height: 176rpx;
  box-shadow: none;
  position: relative;
  transition: background-color 200ms ease, opacity 200ms ease;
}
.record-grid-item:active {
  background-color: #f8fafc;
  opacity: 0.92;
}

/* 图标之间用“小竖线”分隔（每行 3 个） */
.record-grid-item::after {
  content: "";
  position: absolute;
  right: 0;
  top: 28rpx;
  bottom: 28rpx;
  width: 1px;
  background: #e2e8f0;
}
.record-grid-item:nth-child(3n)::after,
.record-grid-item:last-child::after {
  display: none;
}

.record-grid-icon-wrap {
  width: 72rpx;
  height: 72rpx;
  border-radius: 20rpx;
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.record-grid-icon {
  color: #ffffff;
  font-size: 42rpx !important;
  line-height: 1;
}

.record-grid-dot {
  width: 12rpx;
  height: 12rpx;
  border-radius: 50%;
  background: #ef4444;
  border: 2rpx solid #ffffff;
  position: absolute;
  right: -2rpx;
  top: -2rpx;
}

.record-grid-name {
  font-size: 24rpx;
  font-weight: 600;
  color: #0f172a;
  text-align: center;
  line-height: 1.25;
  min-height: 50rpx;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.guest-hint {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #64748b;
  line-height: 1.6;
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
