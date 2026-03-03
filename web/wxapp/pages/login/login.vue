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

					<view class="experience-tip">登录后可使用完整功能，也可先浏览</view>
					<view class="visitor-entry">
						<button class="visitor-btn visitor-main" @click="goHome">暂不登录</button>
						<button class="visitor-btn" @click="goBack">返回上一页</button>
					</view>
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
				const pages = typeof getCurrentPages === 'function' ? getCurrentPages() : []
				if (pages.length > 1) {
					uni.navigateBack({ delta: 1 })
					return
				}
				this.goHome()
			},
			goHome() {
				uni.reLaunch({ url: '/pages/home/home' })
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
	$primary: #0a65d1;
	$primary-strong: #084ea7;
	$text: #0f172a;
	$muted: #64748b;

	.login-container {
		min-height: 100vh;
		display: flex;
		position: relative;
		overflow: hidden;
		background: linear-gradient(180deg, #f7fbff 0%, #eef5ff 52%, #eaf2ff 100%);

		&::before,
		&::after {
			content: '';
			position: absolute;
			border-radius: 50%;
			filter: blur(10rpx);
			pointer-events: none;
		}

		&::before {
			width: 560rpx;
			height: 560rpx;
			top: -180rpx;
			left: -150rpx;
			background: radial-gradient(circle, rgba(10, 101, 209, 0.12) 0%, rgba(10, 101, 209, 0) 72%);
		}

		&::after {
			width: 620rpx;
			height: 620rpx;
			right: -230rpx;
			bottom: -240rpx;
			background: radial-gradient(circle, rgba(21, 128, 61, 0.1) 0%, rgba(21, 128, 61, 0) 74%);
		}
	}

	.body {
		flex: 1;
		display: flex;
		flex-direction: column;
		width: 100%;
		box-sizing: border-box;
		padding: calc(env(safe-area-inset-top) + 24rpx) 44rpx calc(env(safe-area-inset-bottom) + 18rpx);
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
		padding-bottom: 44rpx;
		animation: fadeUp 0.45s ease both;

		.logo-icon {
			width: 110rpx;
			height: 110rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-bottom: 16rpx;
			background: transparent;

			.logo-symbol {
				font-size: 84rpx;
				color: $primary;
			}
		}

		.brand-name {
			font-size: 56rpx;
			font-weight: 800;
			letter-spacing: 2rpx;
			color: $text;
			line-height: 1.2;
		}

		.brand-slogan {
			margin-top: 8rpx;
			font-size: 24rpx;
			color: $muted;
		}
	}

	.loading-section,
	.form-section {
		width: 100%;
	}

	.loading-section {
		min-height: 240rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;

		.loading-spinner {
			width: 54rpx;
			height: 54rpx;
			border-radius: 50%;
			border: 5rpx solid rgba(10, 101, 209, 0.16);
			border-top-color: $primary;
			animation: spin 0.8s linear infinite;
			margin-bottom: 18rpx;
		}

		.loading-text {
			font-size: 26rpx;
			color: $muted;
		}
	}

	.form-section {
		display: flex;
		flex-direction: column;
		gap: 14rpx;
		animation: fadeUp 0.5s ease both;
	}

	.login-form {
		width: 100%;
	}

	.wx-login-btn,
	.login-btn,
	.visitor-btn {
		width: 100%;
		min-height: 92rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 10rpx;
		border: none;
		padding: 0;
		margin: 0;
		background: transparent;
		border-radius: 0;
		line-height: 1.25;

		&::after {
			border: none;
		}

		&:active {
			opacity: 0.65;
		}

		&[disabled] {
			opacity: 0.38;
		}
	}

	.wx-login-btn {
		color: #0d9c56;
		font-size: 38rpx;
		font-weight: 700;
		letter-spacing: 1rpx;

		.wx-icon {
			font-size: 42rpx;
		}
	}

	.form-item {
		margin-bottom: 18rpx;

		.form-label {
			display: block;
			font-size: 24rpx;
			color: $text;
			font-weight: 600;
			margin-bottom: 10rpx;
		}

		.input-wrapper {
			min-height: 90rpx;
			padding: 0 16rpx;
			display: flex;
			align-items: center;
			background: rgba(255, 255, 255, 0.76);
			border-radius: 14rpx;
			box-sizing: border-box;

			&:focus-within {
				background: rgba(255, 255, 255, 0.94);
			}

			.input-icon {
				font-size: 32rpx;
				margin-right: 10rpx;
				color: rgba(10, 101, 209, 0.88);
			}

			.input {
				flex: 1;
				height: 90rpx;
				line-height: 90rpx;
				font-size: 30rpx;
				color: $text;
				font-weight: 500;
			}

			.placeholder {
				color: rgba(100, 116, 139, 0.85);
			}

			.code-btn {
				min-width: 144rpx;
				text-align: right;
				padding: 0 6rpx;
				font-size: 24rpx;
				font-weight: 600;
				color: $primary;
				border: none;
				background: transparent;
				border-radius: 0;

				&::after {
					border: none;
				}

				&[disabled] {
					color: rgba(100, 116, 139, 0.6);
				}
			}
		}
	}

	.login-btn {
		margin-top: 14rpx;
		margin-bottom: 10rpx;
		color: $primary;
		font-size: 40rpx;
		font-weight: 800;
		letter-spacing: 2rpx;
	}

	.agree-section {
		margin-top: 150rpx;
	}

	.agree-line {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-wrap: nowrap;
		white-space: nowrap;
		overflow: hidden;
		gap: 6rpx;

		.agree-left {
			display: flex;
			align-items: center;
			gap: 6rpx;
			min-width: 0;
		}

		.agree-icon {
			font-size: 28rpx;
			color: $primary;
		}

		.agree-text {
			font-size: 22rpx;
			color: $muted;
		}

		.agree-link {
			font-size: 22rpx;
			font-weight: 600;
			color: $primary;

			&:active {
				opacity: 0.7;
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
		margin-top: 4rpx;

		.text-grey {
			font-size: 24rpx;
			color: $muted;
		}

		.link-text {
			font-size: 24rpx;
			font-weight: 600;
			color: $primary;

			&:active {
				opacity: 0.7;
			}
		}
	}

	.experience-tip {
		text-align: center;
		font-size: 22rpx;
		color: $muted;
	}

	.visitor-entry {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 28rpx;
		margin-top: 4rpx;
	}

	.visitor-btn {
		width: auto;
		min-height: 72rpx;
		padding: 0 4rpx;
		font-size: 28rpx;
		color: $muted;
		font-weight: 600;
	}

	.visitor-main {
		color: $primary-strong;
		font-weight: 700;
	}

	.bottom-agreements {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-wrap: wrap;
		gap: 10rpx;
		padding: 22rpx 0 4rpx;

		.agreement-link {
			font-size: 22rpx;
			color: rgba(71, 85, 105, 0.92);

			&:active {
				opacity: 0.7;
			}

			&.disabled {
				opacity: 0.52;
			}
		}

		.divider {
			font-size: 20rpx;
			color: rgba(100, 116, 139, 0.5);
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
		from {
			opacity: 0;
			transform: translateY(14rpx);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>
