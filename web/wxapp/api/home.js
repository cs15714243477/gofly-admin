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
        showLoading: true,
        auth: true,
        loadingMsg: "加载中",
      },
    }),
};
