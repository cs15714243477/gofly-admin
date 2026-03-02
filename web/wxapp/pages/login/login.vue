<template>
	<view class="login-container">
		<view class="body">
			<view class="main">
				<!-- Logo区域 - Exaggerated Minimalism -->
				<view class="logo-section">
					<view class="logo-icon">
						<text class="material-symbols-outlined logo-symbol">real_estate_agent</text>
					</view>
					<view class="brand-name">快销房智选</view>
					<view class="brand-slogan">经纪人助手</view>
				</view>

				<!-- 已登录：加载状态 -->
				<view v-if="checkingLogin" class="loading-section">
					<view class="loading-spinner"></view>
					<view class="loading-text">正在加载中...</view>
				</view>

				<!-- 未登录：登录表单 -->
				<view v-else class="form-section">
					<!-- #ifdef MP-WEIXIN -->
					<view class="login-form">
						<button 
							class="wx-login-btn" 
							:disabled="submitting || !agreed" 
							open-type="getPhoneNumber" 
							@getphonenumber="onGetPhoneNumber"
						>
							<text class="material-symbols-outlined wx-icon">phone_iphone</text>
							<text>{{ submitting ? '登录中...' : '手机号一键登录' }}</text>
						</button>
					</view>

					<view class="agree-section">
						<view class="agree-line">
							<view class="agree-left" @click="toggleAgree">
								<text class="material-symbols-outlined agree-icon">{{ agreed ? 'check_circle' : 'radio_button_unchecked' }}</text>
								<text class="agree-text">我已阅读并同意</text>
							</view>
							<text
								class="agree-link"
								:class="{ disabled: !hasAgreementDoc('user_agreement') }"
								@click.stop="openAgreement('user_agreement')"
							>《{{ agreementDocs.user_agreement.title }}》</text>
							<text class="agree-sep">和</text>
							<text
								class="agree-link"
								:class="{ disabled: !hasAgreementDoc('privacy_policy') }"
								@click.stop="openAgreement('privacy_policy')"
							>《{{ agreementDocs.privacy_policy.title }}》</text>
						</view>
					</view>
					<!-- #endif -->

					<!-- #ifndef MP-WEIXIN -->
					<view class="login-form">
						<view class="form-item">
							<text class="form-label">手机号</text>
							<view class="input-wrapper">
								<text class="material-symbols-outlined input-icon">smartphone</text>
								<input 
									v-model="mobile" 
									class="input" 
									type="tel" 
									inputmode="numeric"
									maxlength="11" 
									placeholder="请输入手机号" 
									placeholder-class="placeholder" 
								/>
							</view>
						</view>

						<view class="form-item">
							<text class="form-label">验证码</text>
							<view class="input-wrapper">
								<text class="material-symbols-outlined input-icon">shield</text>
								<input 
									v-model="captcha" 
									class="input" 
									type="number" 
									inputmode="numeric"
									maxlength="6" 
									placeholder="请输入验证码" 
									placeholder-class="placeholder" 
								/>
								<button class="code-btn" :disabled="true">获取验证码</button>
							</view>
						</view>

						<view class="agree-section">
							<view class="agree-line">
								<view class="agree-left" @click="toggleAgree">
									<text class="material-symbols-outlined agree-icon">{{ agreed ? 'check_circle' : 'radio_button_unchecked' }}</text>
									<text class="agree-text">我已阅读并同意</text>
								</view>
								<text
									class="agree-link"
									:class="{ disabled: !hasAgreementDoc('user_agreement') }"
									@click.stop="openAgreement('user_agreement')"
								>《{{ agreementDocs.user_agreement.title }}》</text>
								<text class="agree-sep">和</text>
								<text
									class="agree-link"
									:class="{ disabled: !hasAgreementDoc('privacy_policy') }"
									@click.stop="openAgreement('privacy_policy')"
								>《{{ agreementDocs.privacy_policy.title }}》</text>
							</view>
						</view>

						<button class="login-btn" :disabled="submitting || !agreed" @click="handleLogin">
							{{ submitting ? '登录中...' : '登录' }}
						</button>

						<view class="footer-links">
							<text class="text-grey">还没有账号？</text>
							<text class="link-text" @click="goToRegister">去完善信息</text>
						</view>
					</view>
					<!-- #endif -->
				</view>
			</view>

			<!-- 底部协议 -->
			<view class="bottom-agreements">
				<text 
					class="agreement-link" 
					:class="{ disabled: !hasAgreementDoc('user_agreement') }" 
					@click="openAgreement('user_agreement')"
				>{{ agreementDocs.user_agreement.title }}</text>
				<text class="divider">·</text>
				<text 
					class="agreement-link" 
					:class="{ disabled: !hasAgreementDoc('privacy_policy') }" 
					@click="openAgreement('privacy_policy')"
				>{{ agreementDocs.privacy_policy.title }}</text>
				<text class="divider">·</text>
				<text 
					class="agreement-link" 
					:class="{ disabled: !hasAgreementDoc('help_center') }" 
					@click="openAgreement('help_center')"
				>{{ agreementDocs.help_center.title }}</text>
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
				const registerTicket = String(data.register_ticket || exdata.register_ticket || extra.register_ticket || '').trim()

				if (mobile) {
					try { uni.setStorageSync('hm_phone', mobile) } catch (e) {}
				}
				if (reason) {
					try { uni.setStorageSync('wxapp_register_reject_reason', reason) } catch (e) {}
				} else {
					try { uni.removeStorageSync('wxapp_register_reject_reason') } catch (e) {}
				}

				if (registerTicket) {
					try { uni.setStorageSync('wxapp_register_ticket', registerTicket) } catch (e) {}
				} else {
					try { uni.removeStorageSync('wxapp_register_ticket') } catch (e) {}
				}
				try { uni.removeStorageSync('wxapp_register_phone_code') } catch (e) {}

				const q = [`mobile=${encodeURIComponent(mobile || '')}`]
				if (registerTicket) q.push(`ticket=${encodeURIComponent(registerTicket)}`)
				uni.navigateTo({
					url: `/pages/registration/registration?${q.join('&')}`
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
						if (this.handleAuditGate(res, { mobile: this.mobile })) return
						return
					}
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
				uni.reLaunch({ url: '/pages/property_list/property_list' })
			}
		}
	}
</script>

<style lang="scss">
	$primary: #0f766e;
	$primary-2: #0ea5a4;
	$accent: #0b79d0;
	$text: #0f172a;
	$muted: #5f6f81;
	$surface: rgba(255, 255, 255, 0.88);

	.login-container {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		position: relative;
		overflow: hidden;
		background: linear-gradient(165deg, #eef8f5 0%, #edf7ff 52%, #f7fbfe 100%);

		&::before,
		&::after {
			content: '';
			position: absolute;
			border-radius: 50%;
			filter: blur(8rpx);
			opacity: 0.75;
			animation: auraMove 11s ease-in-out infinite;
		}

		&::before {
			width: 520rpx;
			height: 520rpx;
			top: -170rpx;
			left: -120rpx;
			background: radial-gradient(circle, rgba(14, 165, 164, 0.28) 0%, rgba(14, 165, 164, 0) 72%);
		}

		&::after {
			width: 620rpx;
			height: 620rpx;
			right: -230rpx;
			bottom: -220rpx;
			background: radial-gradient(circle, rgba(11, 121, 208, 0.22) 0%, rgba(11, 121, 208, 0) 70%);
			animation-delay: 1.7s;
		}
	}

	.body {
		flex: 1;
		display: flex;
		flex-direction: column;
		box-sizing: border-box;
		width: 100%;
		padding: 0 40rpx;
		padding-top: calc(env(safe-area-inset-top) + 28rpx);
		padding-bottom: calc(env(safe-area-inset-bottom) + 14rpx);
		position: relative;
		z-index: 1;
	}

	.main {
		flex: 1;
		width: 100%;
		max-width: 690rpx;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.logo-section {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 0 0 44rpx;
		animation: fadeUp 0.62s ease both;

		.logo-icon {
			width: 132rpx;
			height: 132rpx;
			border-radius: 34rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-bottom: 24rpx;
			background: linear-gradient(135deg, $primary 0%, $accent 100%);
			box-shadow: 0 18rpx 44rpx rgba(15, 118, 110, 0.26);

			.logo-symbol {
				color: #ffffff;
				font-size: 74rpx;
				font-weight: 700;
			}
		}

		.brand-name {
			font-size: 54rpx;
			font-weight: 900;
			color: $text;
			letter-spacing: 2rpx;
			margin-bottom: 6rpx;
			text-shadow: 0 4rpx 14rpx rgba(15, 23, 42, 0.1);
		}

		.brand-slogan {
			font-size: 24rpx;
			color: $muted;
			letter-spacing: 2rpx;
		}
	}

	.loading-section,
	.form-section {
		width: 100%;
		background: $surface;
		border-radius: 30rpx;
		backdrop-filter: blur(16rpx);
		-webkit-backdrop-filter: blur(16rpx);
	}

	.loading-section {
		flex: 0 0 auto;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		min-height: 320rpx;

		.loading-spinner {
			width: 56rpx;
			height: 56rpx;
			border-radius: 50%;
			border: 6rpx solid rgba(15, 118, 110, 0.18);
			border-top-color: $primary;
			animation: spin 0.76s linear infinite;
			margin-bottom: 20rpx;
		}

		.loading-text {
			font-size: 27rpx;
			color: $muted;
			font-weight: 500;
		}
	}

	.form-section {
		flex: 0 0 auto;
		display: flex;
		flex-direction: column;
		padding: 28rpx;
		animation: fadeUp 0.7s ease both;
	}

	.login-form {
		width: 100%;
	}

	.wx-login-btn,
	.login-btn {
		width: 100%;
		height: 96rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12rpx;
		font-size: 33rpx;
		font-weight: 700;
		border: none;
		padding: 0;
		margin: 0;
		letter-spacing: 1rpx;
		border-radius: 18rpx;
		color: #ffffff;

		&::after {
			border: none;
		}

		&:active {
			opacity: 0.86;
			transform: scale(0.992);
		}

		&[disabled] {
			opacity: 0.5;
		}
	}

	.wx-login-btn {
		background: linear-gradient(135deg, #07c160 0%, #059669 100%);
		margin-bottom: 16rpx;

		.wx-icon {
			font-size: 42rpx;
		}
	}

	.form-item {
		margin-bottom: 20rpx;

		.form-label {
			display: block;
			font-size: 25rpx;
			font-weight: 600;
			color: $text;
			margin-bottom: 10rpx;
		}

		.input-wrapper {
			display: flex;
			align-items: center;
			min-height: 92rpx;
			box-sizing: border-box;
			padding: 0 22rpx;
			background: #ffffff;
			border-radius: 16rpx;

			&:focus-within {
				background: #ffffff;
				box-shadow: 0 0 0 4rpx rgba(15, 118, 110, 0.12);
			}

			.input-icon {
				font-size: 34rpx;
				color: rgba(15, 118, 110, 0.9);
				margin-right: 12rpx;
			}

			.input {
				flex: 1;
				height: 92rpx;
				line-height: 92rpx;
				font-size: 30rpx;
				color: $text;
				font-weight: 500;
			}

			.placeholder {
				color: rgba(100, 116, 139, 0.86);
			}

			.code-btn {
				min-width: 150rpx;
				background: rgba(15, 118, 110, 0.1);
				color: $primary;
				font-size: 24rpx;
				font-weight: 600;
				padding: 0 18rpx;
				border-radius: 999rpx;
				text-align: center;

				&::after {
					border: none;
				}

				&[disabled] {
					color: rgba(100, 116, 139, 0.58);
					background: rgba(100, 116, 139, 0.11);
				}
			}
		}
	}

	.login-btn {
		margin-top: 24rpx;
		margin-bottom: 14rpx;
		background: linear-gradient(135deg, $primary 0%, $primary-2 100%);
	}

	.agree-section {
		margin-top: 8rpx;
	}

	.agree-line {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6rpx;
		flex-wrap: nowrap;
		white-space: nowrap;
		overflow: hidden;

		.agree-left {
			display: flex;
			align-items: center;
			gap: 6rpx;
			min-width: 0;
		}

		.agree-icon {
			font-size: 30rpx;
			color: $primary;
		}

		.agree-text {
			font-size: 22rpx;
			color: $muted;
		}

		.agree-link {
			font-size: 22rpx;
			color: $primary;
			font-weight: 600;

			&:active {
				opacity: 0.72;
			}

			&.disabled {
				color: rgba(100, 116, 139, 0.56);
			}
		}

		.agree-sep {
			font-size: 22rpx;
			color: $muted;
		}
	}

	.footer-links {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 10rpx;
		margin-top: 8rpx;

		.text-grey {
			font-size: 24rpx;
			color: $muted;
		}

		.link-text {
			font-size: 24rpx;
			color: $primary;
			font-weight: 600;

			&:active {
				opacity: 0.72;
			}
		}
	}

	.bottom-agreements {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 10rpx;
		padding: 22rpx 0 4rpx;
		flex-wrap: wrap;

		.agreement-link {
			font-size: 22rpx;
			color: rgba(80, 97, 118, 0.9);

			&:active {
				opacity: 0.72;
			}

			&.disabled {
				opacity: 0.52;
			}
		}

		.divider {
			font-size: 20rpx;
			color: rgba(80, 97, 118, 0.42);
		}
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}

	@keyframes fadeUp {
		0% {
			opacity: 0;
			transform: translateY(20rpx);
		}
		100% {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes auraMove {
		0%,
		100% {
			transform: translate3d(0, 0, 0) scale(1);
		}
		50% {
			transform: translate3d(10rpx, -10rpx, 0) scale(1.05);
		}
	}
</style>
