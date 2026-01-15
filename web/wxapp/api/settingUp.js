import request from '@/utils/request';
export default {
	// 修改头像
	changeAvatar: (data) =>
		request({
			url: '/api/custom/user/up_info',
			method: 'POST',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '加载中',
			},
		}),
	// 修改手机号
	changeTel: (data) =>
		request({
			url: '/api/custom/user/up_info',
			method: 'POST',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '加载中',
			},
		}),
	// 修改昵称
	changeName: (data) =>
		request({
			url: '/api/custom/user/up_info',
			method: 'POST',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '加载中',
			},
		}),

}