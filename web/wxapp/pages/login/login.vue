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
				<view class="login-card">
					<!-- Logo -->
					<view class="logo-box">
						<view class="logo-icon">
							<text class="material-symbols-outlined logo-symbol">real_estate_agent</text>
						</view>
					</view>

					<!-- 标题 -->
					<view class="header">
						<view class="title">快销房
						智选经纪人助手欢迎你</view>
						<view class="subtitle">高效管理房源 · 智能门锁一键通</view>
					</view>

					<!-- 亮点（简洁三条） -->
					<view class="feature-list">
						<view class="feature-item">
							<view class="feature-dot"></view>
							<text class="feature-text">房源管理更高效</text>
						</view>
						<view class="feature-item">
							<view class="feature-dot"></view>
							<text class="feature-text">智能门锁一键开锁</text>
						</view>
						<view class="feature-item">
							<view class="feature-dot"></view>
							<text class="feature-text">记录可追溯更安心</text>
						</view>
					</view>

				<!-- 特别提示 -->
<!--				<view class="alert-box">-->
<!--					<view class="alert-icon">-->
<!--						<text class="material-symbols-outlined warning-icon">warning</text>-->
<!--					</view>-->
<!--					<view class="alert-text">-->
<!--						特别提示：使用本系统用户能够通过手机获得房源智能门锁开启权限，账号仅限经纪人本人使用，严禁外借他用。-->
<!--					</view>-->
<!--				</view>-->

				<!-- 已登录：加载并跳转首页 -->
				<view v-if="checkingLogin" class="loading-card">
					<view class="loading-row">
						<view class="loading-spinner"></view>
						<view class="loading-text">正在加载数据中...</view>
					</view>
					<view class="loading-sub">即将进入首页</view>
				</view>

				<!-- 未登录：登录表单 -->
				<view v-else class="form">
					<!-- #ifdef MP-WEIXIN -->
					<!-- 手机号快捷登录（手机号授权）——默认只显示这个 -->
					<button class="wx-login-btn" :disabled="submitting || !agreed" open-type="getPhoneNumber" @getphonenumber="onGetPhoneNumber">
						<text class="material-symbols-outlined wx-icon">phone_iphone</text>
						<text>手机号快捷登录</text>
					</button>
					<view class="agree-row">
						<view class="agree-left" @click="toggleAgree">
							<text class="material-symbols-outlined agree-icon">{{ agreed ? 'check_circle' : 'radio_button_unchecked' }}</text>
							<text class="agree-text">我已阅读并同意</text>
						</view>
						<view class="agree-links">
							<text class="agree-link" :class="{ disabled: !hasAgreementDoc('user_agreement') }" @click.stop="openAgreement('user_agreement')">《{{ agreementDocs.user_agreement.title }}》</text>
							<text class="agree-sep">和</text>
							<text class="agree-link" :class="{ disabled: !hasAgreementDoc('privacy_policy') }" @click.stop="openAgreement('privacy_policy')">《{{ agreementDocs.privacy_policy.title }}》</text>
						</view>
					</view>
					<!-- 临时注释：更多登录方式先隐藏
					<button class="more-login-btn" @click="toggleMoreLogin">
						<text>{{ showMoreLogin ? '收起其他登录方式' : '更多登录方式' }}</text>
						<text class="material-symbols-outlined more-icon">{{ showMoreLogin ? 'expand_less' : 'expand_more' }}</text>
					</button>
					-->
					<!-- #endif -->

					<!-- #ifdef MP-WEIXIN -->
					<!-- 临时注释：更多登录方式先隐藏
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
					-->
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

					<view class="agree-row">
						<view class="agree-left" @click="toggleAgree">
							<text class="material-symbols-outlined agree-icon">{{ agreed ? 'check_circle' : 'radio_button_unchecked' }}</text>
							<text class="agree-text">我已阅读并同意</text>
						</view>
						<view class="agree-links">
							<text class="agree-link" :class="{ disabled: !hasAgreementDoc('user_agreement') }" @click.stop="openAgreement('user_agreement')">《{{ agreementDocs.user_agreement.title }}》</text>
							<text class="agree-sep">和</text>
							<text class="agree-link" :class="{ disabled: !hasAgreementDoc('privacy_policy') }" @click.stop="openAgreement('privacy_policy')">《{{ agreementDocs.privacy_policy.title }}》</text>
						</view>
					</view>

					<button class="login-btn" :disabled="submitting || !agreed" @click="handleLogin">{{ submitting ? '登录中...' : '登录' }}</button>

					<view class="footer-links">
						<text class="text-grey">还没有账号？</text>
						<text class="link-text" @click="goToRegister">完善信息</text>
					</view>
					<!-- #endif -->
				</view>
				</view>
			</view>

			<!-- 底部协议（固定贴底，避免撑高出现滚动） -->
			<view class="bottom-agreements">
				<text class="agreement-link" :class="{ disabled: !hasAgreementDoc('user_agreement') }" @click="openAgreement('user_agreement')">{{ agreementDocs.user_agreement.title }}</text>
				<text class="divider">|</text>
				<text class="agreement-link" :class="{ disabled: !hasAgreementDoc('privacy_policy') }" @click="openAgreement('privacy_policy')">{{ agreementDocs.privacy_policy.title }}</text>
				<text class="divider">|</text>
				<text class="agreement-link" :class="{ disabled: !hasAgreementDoc('help_center') }" @click="openAgreement('help_center')">{{ agreementDocs.help_center.title }}</text>
			</view>
		</view>
	</view>
</template>

<script>
	import $store from '@/store'
	import userApi from '@/api/user'
	export default {
		onLoad() {
			this.loadAgreementDocs()
			this.bootstrapIfLoggedIn()
		},
		data() {
			return {
				mobile: '',
				captcha: '',
				submitting: false,
				checkingLogin: false,
				agreed: false,
				agreementDocs: {
					user_agreement: { title: '用户协议', content: '', url: '' },
					privacy_policy: { title: '隐私政策', content: '', url: '' },
					help_center: { title: '帮助中心', content: '', url: '' }
				}
			}
		},
		methods: {
			toggleAgree() {
				this.agreed = !this.agreed
			},
			ensureAgreed() {
				if (this.agreed) return true
				uni.showToast({ title: '请先阅读并同意协议', icon: 'none' })
				return false
			},
			goBack() {
				uni.navigateBack()
			},
			goToRegister() {
				uni.navigateTo({
					url: '/pages/registration/registration?mode=complete'
				})
			},
			async loadAgreementDocs() {
				try {
					const res = await userApi.getLoginDocs(false)
					if (!res || res.code !== 0 || !res.data || !Array.isArray(res.data.docs)) return
					const nextDocs = { ...this.agreementDocs }
					res.data.docs.forEach((item) => {
						const key = String(item && item.key ? item.key : '').trim()
						if (!key || !nextDocs[key]) return
						const title = String(item && item.title ? item.title : '').trim()
						const content = String(item && item.content ? item.content : '').trim()
						const url = String(item && item.url ? item.url : '').trim()
						nextDocs[key] = {
							title: title || nextDocs[key].title,
							content,
							url,
						}
					})
					this.agreementDocs = nextDocs
				} catch (e) {}
			},
			hasAgreementDoc(key) {
				const doc = (this.agreementDocs && this.agreementDocs[key]) ? this.agreementDocs[key] : null
				if (!doc) return false
				return !!(String(doc.content || '').trim() || String(doc.url || '').trim())
			},
			openAgreement(key) {
				const doc = (this.agreementDocs && this.agreementDocs[key]) ? this.agreementDocs[key] : null
				const title = doc && doc.title ? String(doc.title).trim() : '文档'
				const content = doc && doc.content ? String(doc.content).trim() : ''
				const url = doc && doc.url ? String(doc.url).trim() : ''
				if (!content && !url) {
					uni.showToast({ title: `${title}暂未配置`, icon: 'none' })
					return
				}
				if (content) {
					uni.navigateTo({
						url: `/pages/doc_webview/doc_webview?key=${encodeURIComponent(key)}`
					})
					return
				}
				if (/^https?:\/\//i.test(url)) {
					uni.navigateTo({
						url: `/pages/doc_webview/doc_webview?key=${encodeURIComponent(key)}`
					})
					return
				}
				uni.navigateTo({ url })
			},
			async bootstrapIfLoggedIn() {
				const token = uni.getStorageSync('token')
				if (!token) return
				if (this.checkingLogin) return

				this.checkingLogin = true
				const userStore = $store('user')
				try {
					if (!userStore.isLogin) userStore.setToken(token)
					if (userStore.userInfo && userStore.userInfo.id) {
						this.afterLoginSuccess()
						return
					}
					await userStore.getInfo()
					this.afterLoginSuccess()
				} catch (e) {
					try {
						await userStore.logout(true)
					} catch (e2) {}
				} finally {
					this.checkingLogin = false
				}
			},
			handleAuditGate(res, extra = {}) {
				const code = Number(res && res.code)
				if (![10001, 10002, 10003].includes(code)) return false

				const data = (res && res.data) ? res.data : {}
				const exdata = (res && res.exdata) ? res.exdata : {}
				const mobile = String(data.mobile || exdata.mobile || extra.mobile || '').trim()
				const reason = String(data.audit_reason || '').trim()

				if (mobile) {
					try { uni.setStorageSync('hm_phone', mobile) } catch (e) {}
				}
				if (reason) {
					try { uni.setStorageSync('wxapp_register_reject_reason', reason) } catch (e) {}
				} else {
					try { uni.removeStorageSync('wxapp_register_reject_reason') } catch (e) {}
				}

				// 登录页拿到的手机号授权 code：用于注册页提交审核（一次性）
				if (extra && extra.phone_code) {
					try { uni.setStorageSync('wxapp_register_phone_code', String(extra.phone_code)) } catch (e) {}
				}

				uni.navigateTo({
					url: `/pages/registration/registration?mobile=${encodeURIComponent(mobile || '')}`
				})
				return true
			},
			async handleLogin() {
				if (!this.ensureAgreed()) return
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
					if (!res) return
					if (res.code !== 0) {
						// 未注册/审核中/已拒绝：跳转注册页
						if (this.handleAuditGate(res, { mobile: this.mobile })) return
						return
					}
					// token 由后端 token 字段返回，拦截器也会自动 setToken；这里兜底一次
					if (res.token) $store('user').setToken(res.token)
					await $store('user').getInfo().catch(() => {})
					this.afterLoginSuccess()
				} finally {
					this.submitting = false
				}
			},
			// #ifdef MP-WEIXIN
			onGetPhoneNumber(e) {
				if (!this.ensureAgreed()) return
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
							if (!res) return
							if (res.code !== 0) {
								// 未注册/审核中/已拒绝：跳转注册页
								if (this.handleAuditGate(res, { phone_code: phoneCode })) return
								return
							}
							if (res.token) $store('user').setToken(res.token)
							await $store('user').getInfo().catch(() => {})
							this.afterLoginSuccess()
						} finally {
							this.submitting = false
						}
					},
					fail: () => {
						this.submitting = false
						uni.showToast({ title: '获取登录凭证失败', icon: 'none' })
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
		position: relative;
		background: linear-gradient(180deg, #f0fdfa 0%, #ecfeff 46%, #f8fafc 100%);
		background-image:
			radial-gradient(circle at 18% 12%, rgba(20, 184, 166, 0.18), transparent 56%),
			radial-gradient(circle at 84% 0%, rgba(3, 105, 161, 0.14), transparent 52%);
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
		padding: calc(env(safe-area-inset-top) + 96rpx) 48rpx 0;
		max-width: 920rpx;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
		position: relative;
		z-index: 1;
	}

	.main {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: stretch;
		justify-content: center;
		min-height: 0;
	}

	.login-card {
		width: 100%;
		padding: 56rpx 44rpx 44rpx;
		border-radius: 32rpx;
		background: rgba(255, 255, 255, 0.72);
		border: 1rpx solid rgba(255, 255, 255, 0.9);
		box-shadow: 0 26rpx 70rpx rgba(15, 118, 110, 0.14);
		backdrop-filter: blur(18px);
		-webkit-backdrop-filter: blur(18px);
		box-sizing: border-box;
	}

	.logo-box {
		display: flex;
		justify-content: center;
		margin-bottom: 36rpx;

		.logo-icon {
			width: 148rpx;
			height: 148rpx;
			background: linear-gradient(135deg, #14b8a6 0%, #0369a1 100%);
			border-radius: 30rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			box-shadow: 0 14rpx 36rpx rgba(3, 105, 161, 0.18);
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
		margin-bottom: 28rpx;

		.title {
			font-size: 58rpx;
			font-weight: 900;
			color: #0f172a;
			line-height: 1.2;
			margin-bottom: 18rpx;
		}

		.subtitle {
			font-size: 28rpx;
			color: #475569;
			font-weight: 600;
		}
	}

	.feature-list {
		display: flex;
		flex-direction: column;
		gap: 14rpx;
		margin: 0 0 38rpx;
		padding: 0 10rpx;
	}

	.feature-item {
		display: flex;
		align-items: center;
		gap: 12rpx;
	}

	.feature-dot {
		width: 10rpx;
		height: 10rpx;
		border-radius: 999px;
		background: linear-gradient(135deg, #14b8a6 0%, #0369a1 100%);
		box-shadow: 0 0 0 10rpx rgba(20, 184, 166, 0.08);
		flex-shrink: 0;
	}

	.feature-text {
		font-size: 26rpx;
		color: #475569;
		font-weight: 600;
	}

	.loading-card {
		width: 100%;
		background: rgba(255, 255, 255, 0.76);
		border: 1rpx solid rgba(255, 255, 255, 0.9);
		border-radius: 24rpx;
		padding: 32rpx 28rpx;
		box-shadow: 0 18rpx 44rpx rgba(15, 118, 110, 0.10);
	}

	.loading-row {
		display: flex;
		align-items: center;
		gap: 18rpx;
	}

	.loading-spinner {
		width: 36rpx;
		height: 36rpx;
		border-radius: 50%;
		border: 4rpx solid rgba(3, 105, 161, 0.18);
		border-top-color: #0369a1;
		animation: spin 0.9s linear infinite;
	}

	.loading-text {
		font-size: 30rpx;
		font-weight: 600;
		color: #0f172a;
	}

	.loading-sub {
		margin-top: 14rpx;
		font-size: 24rpx;
		color: #64748b;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
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
		margin-top: 6rpx;

		.wx-login-btn {
			width: 100%;
			height: 104rpx;
			background: linear-gradient(135deg, #07c160 0%, #059669 100%);
			color: #ffffff;
			font-size: 34rpx;
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
				&:active { background-color: #059669; transform: scale(0.985); }
				&[disabled] {
					background: #9ae6b4;
					box-shadow: none;
					opacity: 0.95;
				}

				.wx-icon { font-size: 44rpx; }
			}

		.agree-row {
			margin-top: 10rpx;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 16rpx;
			padding: 0 4rpx;
			flex-wrap: wrap;

			.agree-left {
				display: flex;
				align-items: center;
				gap: 10rpx;
				min-width: 0;
			}

			.agree-icon {
				font-size: 34rpx;
				color: #0f766e;
				line-height: 1;
			}

			.agree-text {
				font-size: 24rpx;
				color: #475569;
				font-weight: 600;
			}

			.agree-links {
				display: flex;
				align-items: center;
				gap: 8rpx;
				flex-wrap: wrap;
			}

			.agree-link {
				font-size: 24rpx;
				color: #0369a1;
				font-weight: 700;

				&:active {
					color: #075985;
				}

				&.disabled {
					color: #cbd5e1;
				}
			}

			.agree-sep {
				font-size: 24rpx;
				color: #94a3b8;
				font-weight: 600;
			}
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
						border-color: #0369a1;
						box-shadow: 0 0 0 4rpx rgba(3, 105, 161, 0.12);
					}
				}

				.code-btn {
					position: absolute;
					right: 16rpx;
					height: 72rpx;
					background-color: rgba(3, 105, 161, 0.06);
					color: #0369a1;
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
						background-color: rgba(3, 105, 161, 0.10);
					}
				}
			}
		}

		.login-btn {
			width: 100%;
			height: 104rpx;
			background-color: #0369a1;
			color: #ffffff;
			font-size: 36rpx;
			font-weight: bold;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-top: 32rpx;
			box-shadow: 0 10rpx 26rpx rgba(3, 105, 161, 0.18);
			border: none;

			&[disabled] {
				background-color: #93c5fd;
				box-shadow: none;
				opacity: 0.9;
			}

			&::after {
				border: none;
			}

			&:active {
				background-color: #075985;
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
				color: #0369a1;
				font-weight: bold;
				margin-left: 8rpx;
			}
		}
	}

	.bottom-agreements {
		padding: 22rpx 0 calc(env(safe-area-inset-bottom) + 22rpx);
		display: flex;
		justify-content: center;
		gap: 24rpx;
		opacity: 0.92;

		.agreement-link {
			font-size: 22rpx;
			color: #94a3b8;
			font-weight: 500;

			&:active {
				color: #0369a1;
			}

			&.disabled {
				color: #cbd5e1;
			}
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
