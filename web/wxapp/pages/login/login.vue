<template>
	<view class="login-container">
		<!-- 顶部导航栏 -->
		<!-- <view class="nav-bar">
			<view class="back-btn" @click="goBack">
				<text class="material-symbols-outlined">arrow_back</text>
			</view>
			<view class="nav-title">经纪人登录</view>
			<view class="placeholder-box"></view>
		</view> -->

		<view class="body">
			<view class="main">
				<!-- Logo -->
				<view class="logo-box">
					<view class="logo-icon">
						<text class="material-symbols-outlined logo-symbol">real_estate_agent</text>
					</view>
				</view>

				<!-- 标题 -->
				<view class="header">
					<view class="title">极速房 欢迎你</view>
					<view class="subtitle">高效管理房源 · 智能门锁一键通</view>
				</view>

				<!-- 特别提示 -->
				<view class="alert-box">
					<view class="alert-icon">
						<text class="material-symbols-outlined warning-icon">warning</text>
					</view>
					<view class="alert-text">
						特别提示：使用本系统用户能够通过手机获得房源智能门锁开启权限，账号仅限经纪人本人使用，严禁外借他用。
					</view>
				</view>

				<!-- 登录表单 -->
				<view class="form">
					<!-- #ifdef MP-WEIXIN -->
					<!-- 微信一键登录（手机号授权）——默认只显示这个 -->
					<button class="wx-login-btn" open-type="getPhoneNumber" @getphonenumber="onGetPhoneNumber">
						<text class="material-symbols-outlined wx-icon">phone_iphone</text>
						<text>微信一键登录</text>
					</button>
					<button class="more-login-btn" @click="toggleMoreLogin">
						<text>{{ showMoreLogin ? '收起其他登录方式' : '更多登录方式' }}</text>
						<text class="material-symbols-outlined more-icon">{{ showMoreLogin ? 'expand_less' : 'expand_more' }}</text>
					</button>
					<!-- #endif -->

					<!-- #ifdef MP-WEIXIN -->
					<view v-if="showMoreLogin">
						<view class="divider-row">
							<view class="divider-line"></view>
							<text class="divider-text">手机号验证码登录</text>
							<view class="divider-line"></view>
						</view>

						<view class="form-item">
							<view class="label">手机号</view>
							<view class="input-wrapper">
								<text class="material-symbols-outlined input-icon">smartphone</text>
								<input v-model="mobile" class="input" type="tel" maxlength="11" placeholder="请输入11位手机号码" placeholder-class="placeholder" />
							</view>
						</view>

						<view class="form-item">
							<view class="label">验证码</view>
							<view class="input-wrapper">
								<text class="material-symbols-outlined input-icon">shield</text>
								<input v-model="captcha" class="input" type="number" maxlength="6" placeholder="请输入验证码" placeholder-class="placeholder" />
								<button class="code-btn">获取验证码</button>
							</view>
						</view>

						<button class="login-btn" :disabled="submitting" @click="handleLogin">{{ submitting ? '登录中...' : '登录' }}</button>

						<view class="footer-links">
							<text class="text-grey">还没有账号？</text>
							<text class="link-text" @click="goToRegister">完善信息</text>
						</view>
					</view>
					<!-- #endif -->

					<!-- #ifndef MP-WEIXIN -->
					<view class="form-item">
						<view class="label">手机号</view>
						<view class="input-wrapper">
							<text class="material-symbols-outlined input-icon">smartphone</text>
							<input v-model="mobile" class="input" type="tel" maxlength="11" placeholder="请输入11位手机号码" placeholder-class="placeholder" />
						</view>
					</view>

					<view class="form-item">
						<view class="label">验证码</view>
						<view class="input-wrapper">
							<text class="material-symbols-outlined input-icon">shield</text>
							<input v-model="captcha" class="input" type="number" maxlength="6" placeholder="请输入验证码" placeholder-class="placeholder" />
							<button class="code-btn">获取验证码</button>
						</view>
					</view>

					<button class="login-btn" :disabled="submitting" @click="handleLogin">{{ submitting ? '登录中...' : '登录' }}</button>

					<view class="footer-links">
						<text class="text-grey">还没有账号？</text>
						<text class="link-text" @click="goToRegister">完善信息</text>
					</view>
					<!-- #endif -->
				</view>
			</view>

			<!-- 底部协议（固定贴底，避免撑高出现滚动） -->
			<view class="bottom-agreements">
				<text class="agreement-link">用户协议</text>
				<text class="divider">|</text>
				<text class="agreement-link">隐私政策</text>
				<text class="divider">|</text>
				<text class="agreement-link">帮助中心</text>
			</view>
		</view>
	</view>
</template>

<script>
	import $store from '@/store'
	import userApi from '@/api/user'
	export default {
		data() {
			return {
				// #ifdef MP-WEIXIN
				showMoreLogin: false,
				// #endif
				mobile: '',
				captcha: '',
				submitting: false,
			}
		},
		methods: {
			goBack() {
				uni.navigateBack()
			},
			goToRegister() {
				uni.navigateTo({
					url: '/pages/registration/registration?mode=complete'
				})
			},
			async handleLogin() {
				if (!this.mobile || String(this.mobile).length !== 11) {
					uni.showToast({ title: '请输入11位手机号', icon: 'none' })
					return
				}
				if (!this.captcha) {
					uni.showToast({ title: '请输入验证码', icon: 'none' })
					return
				}
				if (this.submitting) return
				this.submitting = true
				try {
					const res = await userApi.login({ mobile: this.mobile, captcha: this.captcha })
					if (!res || res.code !== 0) return
					// token 由后端 token 字段返回，拦截器也会自动 setToken；这里兜底一次
					if (res.token) $store('user').setToken(res.token)
					await $store('user').getInfo().catch(() => {})
					this.afterLoginSuccess()
				} finally {
					this.submitting = false
				}
			},
			// #ifdef MP-WEIXIN
			/**
			 * 切换更多登录方式的显示状态
			 *
			 * 用于控制微信小程序中更多登录方式按钮的展开/收起状态
			 * 通过反转 showMoreLogin 数据属性的值来实现状态切换
			 *
			 * @returns {void} 无返回值
			 */
			toggleMoreLogin() {
				this.showMoreLogin = !this.showMoreLogin
			},
			// #endif
			// #ifdef MP-WEIXIN
			onGetPhoneNumber(e) {
				// 需要用户点击按钮触发；e.detail.code 需传后端换取手机号（解密）
				const phoneCode = e && e.detail && e.detail.code
				if (!phoneCode) {
					uni.showToast({ title: '未授权手机号', icon: 'none' })
					return
				}
				if (this.submitting) return
				this.submitting = true
				uni.login({
					provider: 'weixin',
					success: async (loginRes) => {
						try {
							const wxCode = loginRes && loginRes.code
							const res = await userApi.wxLogin({ wx_code: wxCode, phone_code: phoneCode })
							if (!res || res.code !== 0) return
							if (res.token) $store('user').setToken(res.token)
							await $store('user').getInfo().catch(() => {})
							this.afterLoginSuccess()
						} finally {
							this.submitting = false
						}
					},
					fail: () => {
						this.submitting = false
						uni.showToast({ title: '获取微信登录凭证失败', icon: 'none' })
					}
				})
			},
			// #endif
			afterLoginSuccess() {
				// 登录成功后默认进入房源列表（tab 最左）
				uni.reLaunch({ url: '/pages/property_list/property_list' })
			}
		}
	}
</script>

<style lang="scss">
	.login-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.nav-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 40rpx 30rpx 20rpx;
		position: absolute;
		top: 0;
		width: 100%;
		z-index: 10;
		box-sizing: border-box;

		.back-btn {
			width: 96rpx;
			height: 96rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			border-radius: 50%;

			.material-symbols-outlined {
				color: #64748b;
				font-size: 48rpx;
			}

			&:active {
				background-color: rgba(226, 232, 240, 0.5);
			}
		}

		.nav-title {
			font-size: 28rpx;
			font-weight: 500;
			color: #64748b;
		}

		.placeholder-box {
			width: 96rpx;
		}
	}

	.body {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: stretch;
		padding: 140rpx 48rpx 0;
		max-width: 960rpx;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
	}

	.main {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 0;
	}

	.logo-box {
		margin-bottom: 48rpx;

		.logo-icon {
			width: 160rpx;
			height: 160rpx;
			background: linear-gradient(135deg, #2d9cf0 0%, #1a7ab5 100%);
			border-radius: 32rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			box-shadow: 0 10rpx 30rpx rgba(45, 156, 240, 0.2);
			transition: transform 0.3s;

			&:active {
				transform: scale(1.05);
			}

			.logo-symbol {
				color: #ffffff;
				font-size: 80rpx;
			}
		}
	}

	.header {
		text-align: center;
		margin-bottom: 48rpx;

		.title {
			font-size: 64rpx;
			font-weight: bold;
			color: #0d151c;
			line-height: 1.2;
			margin-bottom: 24rpx;
		}

		.subtitle {
			font-size: 32rpx;
			color: #64748b;
			font-weight: 400;
		}
	}

	.alert-box {
		width: 100%;
		background-color: #fef2f2;
		border: 1px solid #fee2e2;
		border-radius: 24rpx;
		padding: 28rpx;
		display: flex;
		gap: 24rpx;
		margin-bottom: 40rpx;
		box-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.05);
		box-sizing: border-box;

		.warning-icon {
			color: #dc2626;
			font-size: 40rpx;
		}

		.alert-text {
			flex: 1;
			font-size: 28rpx;
			color: #dc2626;
			font-weight: 500;
			line-height: 1.6;
			text-align: justify;
		}
	}

	.form {
		width: 100%;

		.wx-login-btn {
			width: 100%;
			height: 104rpx;
			background-color: #10b981;
			color: #ffffff;
			font-size: 36rpx;
			font-weight: bold;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 16rpx;
			box-shadow: 0 8rpx 20rpx rgba(16, 185, 129, 0.22);
			border: none;
			margin-bottom: 26rpx;

			&::after { border: none; }
			&:active { background-color: #059669; transform: scale(0.98); }

			.wx-icon { font-size: 44rpx; }
		}

		.more-login-btn {
			width: 100%;
			height: 84rpx;
			background-color: #ffffff;
			border: 1px solid #e2e8f0;
			color: #334155;
			font-size: 28rpx;
			font-weight: 800;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 10rpx;
			box-shadow: 0 6rpx 16rpx rgba(15, 23, 42, 0.05);
			margin-bottom: 22rpx;

			&::after { border: none; }
			&:active { background-color: #f8fafc; transform: scale(0.99); }

			.more-icon { font-size: 36rpx; color: #64748b; }
		}

		.divider-row {
			display: flex;
			align-items: center;
			gap: 16rpx;
			margin-bottom: 22rpx;
			color: #94a3b8;
			font-size: 24rpx;
			font-weight: 600;

			.divider-line {
				flex: 1;
				height: 1px;
				background-color: #e2e8f0;
			}

			.divider-text { white-space: nowrap; }
		}

		.form-item {
			margin-bottom: 28rpx;

			.label {
				font-size: 28rpx;
				font-weight: 500;
				color: #0d151c;
				margin-bottom: 16rpx;
			}

			.input-wrapper {
				position: relative;
				display: flex;
				align-items: center;

				.input-icon {
					position: absolute;
					left: 28rpx;
					color: #94a3b8;
					font-size: 40rpx;
					z-index: 1;
				}

				.input {
					width: 100%;
					height: 104rpx;
					background-color: #ffffff;
					border: 1px solid #cfdde8;
					border-radius: 24rpx;
					padding-left: 88rpx;
					padding-right: 32rpx;
					font-size: 32rpx;
					color: #0d151c;
					transition: all 0.3s;
					box-sizing: border-box;

					&:focus {
						border-color: #2d9cf0;
						box-shadow: 0 0 0 4rpx rgba(45, 156, 240, 0.1);
					}
				}

				.code-btn {
					position: absolute;
					right: 16rpx;
					height: 72rpx;
					background-color: rgba(45, 156, 240, 0.05);
					color: #2d9cf0;
					font-size: 28rpx;
					font-weight: 600;
					padding: 0 24rpx;
					border-radius: 16rpx;
					display: flex;
					align-items: center;
					justify-content: center;
					border: none;
					line-height: 1;
					z-index: 2;

					&::after {
						border: none;
					}

					&:active {
						background-color: rgba(45, 156, 240, 0.1);
					}
				}
			}
		}

		.login-btn {
			width: 100%;
			height: 104rpx;
			background-color: #2d9cf0;
			color: #ffffff;
			font-size: 36rpx;
			font-weight: bold;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-top: 32rpx;
			box-shadow: 0 8rpx 20rpx rgba(45, 156, 240, 0.2);
			border: none;

			&::after {
				border: none;
			}

			&:active {
				background-color: #1a7ab5;
				transform: scale(0.98);
			}
		}

		.footer-links {
			margin-top: 36rpx;
			display: flex;
			justify-content: center;
			align-items: center;
			font-size: 28rpx;

			.text-grey {
				color: #64748b;
			}

			.link-text {
				color: #2d9cf0;
				font-weight: bold;
				margin-left: 8rpx;
			}
		}
	}

	.bottom-agreements {
		padding: 24rpx 0 calc(env(safe-area-inset-bottom) + 24rpx);
		display: flex;
		justify-content: center;
		gap: 24rpx;

		.agreement-link {
			font-size: 24rpx;
			color: #94a3b8;
			font-weight: 500;
		}

		.divider {
			font-size: 24rpx;
			color: #e2e8f0;
		}
	}

	.placeholder {
		color: #94a3b8;
	}
</style>
