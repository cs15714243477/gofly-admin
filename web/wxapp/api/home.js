import request from "@/utils/request";
export default {
  // 获取首页数据
  Homedata: (data) =>
    request({
      url: "/uniapp/home/getList",
      method: "GET",
      data,
      custom: {
        showSuccess: false,
        showError: false,
        showLoading: true,
        auth: false,
        loadingMsg: "加载中",
      },
    }),
};
