import request from "@/utils/request";
export default {
  // 账号登录
  login: (data) =>
    request({
      url: "/uniapp/login",
      method: "POST",
      data,
      custom: {
        showSuccess: true,
        loadingMsg: "登录中",
      },
    }),
  // 微信一键登录（手机号）
  wxLogin: (data) =>
    request({
      url: "/uniapp/wxLogin",
      method: "POST",
      data,
      custom: {
        showSuccess: true,
        loadingMsg: "登录中",
      },
    }),
  // 获取用户信息
  getuserinfo: () =>
    request({
      url: "/uniapp/getUserinfo",
      method: "GET",
      custom: {
        showLoading: true,
        auth: true,
      },
    }),
  // 获取登录页文档配置（用户协议/隐私政策/帮助中心）
  getLoginDocs: (showLoading = false) =>
    request({
      url: "/uniapp/getLoginDocs",
      method: "GET",
      custom: {
        showLoading,
        auth: false,
      },
    }),
  // 获取我的名片资料（不脱敏，给名片预览/编辑使用）
  getCardProfile: (showLoading = true) =>
    request({
      url: "/uniapp/getCardProfile",
      method: "GET",
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 更新我的名片资料
  updateCardProfile: (data, showLoading = true) =>
    request({
      url: "/uniapp/updateCardProfile",
      method: "POST",
      data,
      custom: {
        showLoading,
        auth: true,
        loadingMsg: "保存中...",
      },
    }),
  // 获取门店列表（用于名片页选择门店）
  getStores: (showLoading = true) =>
    request({
      url: "/uniapp/getStores",
      method: "GET",
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 获取经纪人小程序码（名片二维码）
  getAgentWxaCode: (data, showLoading = true) =>
    request({
      url: "/uniapp/getAgentWxaCode",
      method: "GET",
      data,
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 获取经纪人 URLLink（用于复制给客户）
  getAgentUrlLink: (data, showLoading = true) =>
    request({
      url: "/uniapp/getAgentUrlLink",
      method: "GET",
      data,
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 获取经纪人公开名片信息（客户扫码落地页）
  getAgentCard: (data, showLoading = true) =>
    request({
      url: "/uniapp/getAgentCard",
      method: "GET",
      data,
      custom: {
        showLoading,
        auth: false,
      },
    }),
  // 获取工作台记录汇总（我的页）
  getWorkbenchSummary: (showLoading = false) =>
    request({
      url: "/uniapp/getWorkbenchSummary",
      method: "GET",
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 获取工作台记录列表（关注/开锁/带看/浏览/分享/通话）
  getWorkbenchRecords: (data, showLoading = true) =>
    request({
      url: "/uniapp/getWorkbenchRecords",
      method: "GET",
      data,
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 获取名片推荐房源配置（经纪人编辑用）
  getAgentCardRecommendConfig: (showLoading = true) =>
    request({
      url: "/uniapp/getAgentCardRecommendConfig",
      method: "GET",
      custom: {
        showLoading,
        auth: true,
      },
    }),
  // 保存名片推荐房源配置（经纪人编辑用）
  saveAgentCardRecommendConfig: (data, showLoading = true) =>
    request({
      url: "/uniapp/saveAgentCardRecommendConfig",
      method: "POST",
      data,
      custom: {
        showLoading,
        auth: true,
        loadingMsg: "保存中...",
      },
    }),
  // 账号登出
  logout: (data) =>
    request({
      url: "/uniapp/logout",
      method: "POST",
      data,
    }),
};
