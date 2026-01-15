import request from '@/utils/request';
export default {
	// 获取首页数据
	Homedata: (data) =>
		request({
			url: '/home/getList',
			method: 'get',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '加载中',
			},
		}),
}
