import request from '@/utils/request';
export default {
	// 账号登录
	login: (data) =>
		request({
			url: '/uniapp/login',
			method: 'POST',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '登录中',
			},
	}),
	// 微信一键登录（手机号）
	wxLogin: (data) =>
		request({
			url: '/uniapp/wxLogin',
			method: 'POST',
			data,
			custom: {
				showSuccess: true,
				loadingMsg: '登录中',
			},
	}),
	// 获取用户信息
	getuserinfo: () =>request({
			url: '/uniapp/getUserinfo',
			method: 'GET',
			custom: {
				showLoading: true,
				auth: true,
			},
		}),
	// 账号登出
	logout: (data) =>request({
			url: '/uniapp/logout',
			method: 'POST',
			data,
		})
}
