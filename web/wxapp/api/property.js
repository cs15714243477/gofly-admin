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

  // 房源详情（小程序，走 /uniapp/wxproperty 控制器）
  getDetail: (data) =>
    request({
      url: "/uniapp/wxproperty/getDetail",
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
      url: "/uniapp/wxproperty/toggleFollow",
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

  // ------------------------
  // 可维护房源（经纪人端 CRUD）
  // ------------------------

  // 可维护房源列表
  getManageList: (data) =>
    request({
      url: "/uniapp/wxproperty/getManageList",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "加载中",
      },
    }),

  // 获取可维护房源详情（编辑页回显）
  getManageContent: (data) =>
    request({
      url: "/uniapp/wxproperty/getManageContent",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "加载中",
      },
    }),

  // 获取可维护房源装修信息（编辑页）
  getManageRenovation: (data) =>
    request({
      url: "/uniapp/wxproperty/getManageRenovation",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),

  // 保存可维护房源装修信息（编辑页）
  saveManageRenovation: (data) =>
    request({
      url: "/uniapp/wxproperty/saveManageRenovation",
      method: "POST",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "保存中",
      },
    }),

  // 保存可维护房源（新增/编辑）
  saveManage: (data) =>
    request({
      url: "/uniapp/wxproperty/saveManage",
      method: "POST",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "保存中",
      },
    }),

  // 房源表单下拉选项（后台字典配置，前端缓存）
  getFormOptions: (data) =>
    request({
      url: "/uniapp/wxproperty/getFormOptions",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),

  // 删除可维护房源（软删除）
  delManage: (data) =>
    request({
      url: "/uniapp/wxproperty/delManage",
      method: "DELETE",
      data,
      custom: {
        showSuccess: false,
        showLoading: true,
        auth: true,
        loadingMsg: "删除中",
      },
    }),
};
