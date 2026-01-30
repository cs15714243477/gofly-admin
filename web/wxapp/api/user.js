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
	// 获取我的名片资料（不脱敏，给名片预览/编辑使用）
	getCardProfile: (showLoading = true) =>
		request({
			url: '/uniapp/getCardProfile',
			method: 'GET',
			custom: {
				showLoading,
				auth: true,
			},
		}),
	// 更新我的名片资料
	updateCardProfile: (data, showLoading = true) =>
		request({
			url: '/uniapp/updateCardProfile',
			method: 'POST',
			data,
			custom: {
				showLoading,
				auth: true,
				loadingMsg: '保存中...',
			},
		}),
	// 获取门店列表（用于名片页选择门店）
	getStores: (showLoading = true) =>
		request({
			url: '/uniapp/getStores',
			method: 'GET',
			custom: {
				showLoading,
				auth: true,
			},
		}),
	// 获取经纪人小程序码（名片二维码）
	getAgentWxaCode: (data, showLoading = true) =>
		request({
			url: '/uniapp/getAgentWxaCode',
			method: 'GET',
			data,
			custom: {
				showLoading,
				auth: true,
			},
		}),
	// 获取经纪人 URLLink（用于复制给客户）
	getAgentUrlLink: (data, showLoading = true) =>
		request({
			url: '/uniapp/getAgentUrlLink',
			method: 'GET',
			data,
			custom: {
				showLoading,
				auth: true,
			},
		}),
	// 获取经纪人公开名片信息（客户扫码落地页）
	getAgentCard: (data, showLoading = true) =>
		request({
			url: '/uniapp/getAgentCard',
			method: 'GET',
			data,
			custom: {
				showLoading,
				auth: false,
			},
		}),
	// 账号登出
	logout: (data) =>request({
			url: '/uniapp/logout',
			method: 'POST',
			data,
		})
}
