import request from "@/utils/request";
export default {
  // 房源列表
  getList: (data) =>
    request({
      url: "/uniapp/properties/getList",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "加载中",
      },
    }),

  // 房源详情（小程序，走 /houses/wxproperty 控制器）
  getDetail: (data) =>
    request({
      url: "/houses/wxproperty/getDetail",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "加载中",
      },
    }),

  // 关注/取消关注
  toggleFollow: (data) =>
    request({
      url: "/houses/wxproperty/toggleFollow",
      method: "POST",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),

  // 房源筛选项（从数据库去重）
  getFilterOptions: (data) =>
    request({
      url: "/uniapp/properties/getFilterOptions",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),
};
