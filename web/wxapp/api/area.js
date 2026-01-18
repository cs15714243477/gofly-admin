import request from "@/utils/request";

export default {
  // 省份列表
  getProvinces: (data) =>
    request({
      url: "/uniapp/area/getProvinces",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),

  // 子级（城市/区县）
  getChildren: (data) =>
    request({
      url: "/uniapp/area/getChildren",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showLoading: false,
        auth: true,
      },
    }),
};
