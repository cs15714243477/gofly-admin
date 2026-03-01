<template>
	<view class="register-container">
		<!-- 顶部导航：使用 TopHeader，自动与胶囊对齐（MP-WEIXIN） -->
		<TopHeader :title="pageTitle">
			<template #left>
				<view class="th-btn" @click="goBack">
					<text class="material-symbols-outlined">arrow_back_ios_new</text>
				</view>
			</template>
		</TopHeader>

		<view class="body">
			<scroll-view scroll-y="true" class="content">
				<view class="content-inner">
					<!-- 顶部欢迎卡片 -->
					<view class="header-card">
						<view class="header-decor"></view>
						<view class="header-main">
							<view class="brand-icon">
								<text class="material-symbols-outlined">real_estate_agent</text>
							</view>
							<view class="header-texts">
								<view class="title">{{ headerTitle }}</view>
								<view class="subtitle">{{ headerSubtitle }}</view>
							</view>
						</view>
					</view>

					<!-- 审核状态提示 -->
					<view v-if="audit.status" class="audit-card" :class="audit.status">
						<view class="audit-top">
							<text class="material-symbols-outlined audit-icon">{{ auditIcon }}</text>
							<view class="audit-texts">
								<view class="audit-title">{{ auditTitle }}</view>
								<view class="audit-desc">{{ auditDesc }}</view>
							</view>
						</view>
						<view v-if="audit.status === 'rejected' && audit.reason" class="audit-reason">原因：{{ audit.reason }}</view>
						<view class="audit-actions">
							<button v-if="audit.status === 'pending'" class="audit-btn" @click="fetchAuditStatus(true)">刷新状态</button>
							<button v-if="audit.status === 'approved'" class="audit-btn primary" @click="goLogin">去登录</button>
						</view>
					</view>

					<!-- 表单卡片（加留白，避免内容过满） -->
					<view class="form-card">
						<view class="form">
							<!-- 所在地区（省/市/区） -->
							<view class="form-item">
								<view class="label">所在地区</view>
								<picker mode="region" :value="region" @change="onRegionChange">
									<view class="picker-wrapper">
										<text class="material-symbols-outlined icon-fill">location_on</text>
										<text class="picker-text">{{ regionText }}</text>
										<text class="material-symbols-outlined arrow-icon">expand_more</text>
									</view>
								</picker>
							</view>

							<!-- 真实姓名 -->
							<view class="form-item">
								<view class="label">真实姓名</view>
								<view class="input-wrapper">
									<text class="material-symbols-outlined input-icon">person</text>
									<input v-model="form.name" class="input" type="text" placeholder="请输入姓名" placeholder-class="placeholder" :disabled="formLocked" />
								</view>
							</view>

							<!-- 手机号（授权后自动带入） -->
							<view class="form-item">
								<view class="label">手机号</view>
								<view class="input-wrapper">
									<text class="material-symbols-outlined input-icon">smartphone</text>
									<input v-model="form.phone" class="input readonly-input" type="text" placeholder="请先在登录页授权手机号" placeholder-class="placeholder" readonly />
								</view>
							</view>

							<!-- 所属门店/公司 -->
							<view class="form-item">
								<view class="label">所属门店/公司</view>
								<view class="store-mode">
									<view class="mode-item" :class="{ active: storeMode === 'select' }" @click="setStoreMode('select')">选择门店</view>
									<view class="mode-item" :class="{ active: storeMode === 'manual' }" @click="setStoreMode('manual')">手动填写</view>
								</view>

								<view v-if="storeMode === 'select'">
									<picker :range="storeOptions" range-key="name" :value="storeIndex" @change="onStorePick" :disabled="formLocked">
										<view class="picker-wrapper">
											<text class="material-symbols-outlined icon-fill">storefront</text>
											<text class="picker-text">{{ storeDisplay }}</text>
											<text class="material-symbols-outlined arrow-icon">expand_more</text>
										</view>
									</picker>
									<view class="store-addr-tip" v-if="storeSelected && storeSelected.address">{{ storeSelected.address }}</view>
								</view>

								<view v-else>
									<view class="input-wrapper">
										<text class="material-symbols-outlined input-icon">storefront</text>
										<input v-model="form.store_name_text" class="input" type="text" placeholder="请输入门店名称" placeholder-class="placeholder" :disabled="formLocked" />
									</view>
									<view class="input-wrapper mt">
										<text class="material-symbols-outlined input-icon">pin_drop</text>
										<input v-model="form.store_address_text" class="input" type="text" placeholder="门店地址（可选）" placeholder-class="placeholder" :disabled="formLocked" />
									</view>
								</view>
							</view>

							<!-- 协议 -->
							<view class="agreement-row">
								<checkbox-group @change="agreementChange">
									<label class="checkbox-label">
										<checkbox value="agree" :checked="agreed" color="#2d9cf0" style="transform:scale(0.7)" />
										<text class="agreement-text">我已阅读并同意</text>
									</label>
								</checkbox-group>
								<view class="agreement-links">
									<text class="link" :class="{ disabled: !hasAgreementDoc('user_agreement') }" @click.stop="openAgreement('user_agreement')">《{{ agreementDocs.user_agreement.title }}》</text>
									<text class="sep">和</text>
									<text class="link" :class="{ disabled: !hasAgreementDoc('privacy_policy') }" @click.stop="openAgreement('privacy_policy')">《{{ agreementDocs.privacy_policy.title }}》</text>
								</view>
							</view>
						</view>
					</view>

					<view class="bottom-spacer"></view>
				</view>
			</scroll-view>

			<!-- 底部固定按钮（避免撑高页面出现滚动） -->
			<view class="footer">
				<button class="register-btn" :disabled="submitting" @click="handleSubmit">{{ submitText }}</button>
			</view>
		</view>
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'
	import userApi from '@/api/user'

	export default {
		components: { TopHeader },
		onLoad(options) {
			this.loadAgreementDocs()
			this.mode = (options && options.mode) || 'complete'
			const mobile = options && options.mobile ? String(options.mobile).trim() : ''
			if (mobile) this.form.phone = mobile
			// 取登录页保存的手机号（兜底）
			const cachedPhone = uni.getStorageSync('hm_phone')
			if (!this.form.phone && cachedPhone) this.form.phone = cachedPhone

			this.fetchStores()
			this.fetchAuditStatus(false)
		},
		data() {
			return {
				agreed: false,
				agreementDocs: {
					user_agreement: { title: '用户协议', content: '', url: '' },
					privacy_policy: { title: '隐私政策', content: '', url: '' },
				},
				mode: 'complete', // complete | register（预留）
				submitting: false,
				region: ['北京市', '北京市', '朝阳区'],
				storeMode: 'select', // select | manual
				storeOptions: [],
				storeIndex: -1,
				audit: {
					status: '', // pending | approved | rejected | ''
					reason: '',
				},
				form: {
					name: '',
					phone: '',
					store_id: 0,
					store_name_text: '',
					store_address_text: ''
				}
			}
		},
		computed: {
			pageTitle() {
				return '注册申请'
			},
			headerTitle() {
				if (this.audit.status === 'pending') return '资料审核中'
				if (this.audit.status === 'approved') return '审核已通过'
				if (this.audit.status === 'rejected') return '审核未通过'
				return '提交资料'
			},
			headerSubtitle() {
				if (this.audit.status === 'pending') return '我们将在 1 个工作日内完成审核，请耐心等待'
				if (this.audit.status === 'approved') return '请返回登录页进行登录'
				if (this.audit.status === 'rejected') return '请修改资料后重新提交审核'
				return '首次登录请先填写资料并提交审核，通过后才可登录'
			},
			submitText() {
				if (this.audit.status === 'pending') return '审核中'
				if (this.audit.status === 'approved') return '去登录'
				return '提交审核'
			},
			regionText() {
				const r = Array.isArray(this.region) ? this.region : []
				return r.filter(Boolean).join(' ')
			},
			formLocked() {
				return this.audit.status === 'pending' || this.audit.status === 'approved'
			},
			storeSelected() {
				if (this.storeIndex < 0) return null
				return this.storeOptions && this.storeOptions[this.storeIndex] ? this.storeOptions[this.storeIndex] : null
			},
			storeDisplay() {
				if (this.form.store_id && this.storeSelected) return this.storeSelected.name || '请选择门店'
				if (this.form.store_id) return '已选择门店'
				return '请选择门店'
			},
			auditIcon() {
				if (this.audit.status === 'pending') return 'hourglass_top'
				if (this.audit.status === 'approved') return 'check_circle'
				if (this.audit.status === 'rejected') return 'cancel'
				return 'info'
			},
			auditTitle() {
				if (this.audit.status === 'pending') return '资料审核中'
				if (this.audit.status === 'approved') return '审核已通过'
				if (this.audit.status === 'rejected') return '审核未通过'
				return ''
			},
			auditDesc() {
				if (this.audit.status === 'pending') return '提交成功，请等待管理员审核'
				if (this.audit.status === 'approved') return '你已通过审核，现在可以登录'
				if (this.audit.status === 'rejected') return '请根据原因修改后重新提交'
				return ''
			}
		},
		methods: {
			goBack() {
				uni.navigateBack()
			},
			goLogin() {
				uni.reLaunch({ url: '/pages/login/login' })
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
				if (!doc) return
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
			agreementChange(e) {
				this.agreed = e.detail.value.length > 0
			},
			onRegionChange(e) {
				if (this.formLocked) return
				const value = e && e.detail && e.detail.value
				if (Array.isArray(value) && value.length) this.region = value
			},
			setStoreMode(mode) {
				if (this.formLocked) return
				const m = String(mode || '').trim()
				if (m !== 'select' && m !== 'manual') return
				this.storeMode = m
				if (m === 'select') {
					this.form.store_name_text = ''
					this.form.store_address_text = ''
				} else {
					this.form.store_id = 0
					this.storeIndex = -1
				}
			},
			onStorePick(e) {
				if (this.formLocked) return
				const idx = e && e.detail ? Number(e.detail.value) : -1
				if (!Number.isFinite(idx) || idx < 0 || idx >= (this.storeOptions || []).length) {
					this.storeIndex = -1
					this.form.store_id = 0
					return
				}
				this.storeIndex = idx
				const it = this.storeOptions[idx]
				this.form.store_id = it && it.id ? Number(it.id) : 0
			},
			async fetchStores() {
				try {
					const res = await userApi.getRegisterStores(false)
					if (!res || res.code !== 0) return
					const list = Array.isArray(res.data) ? res.data : []
					this.storeOptions = list.map((it) => ({
						id: Number(it.id || 0) || 0,
						name: it.name || '',
						address: it.address || '',
					})).filter((it) => it.id > 0 && it.name)
					// 若已有 store_id，尝试回显 index
					if (this.form.store_id) {
						const idx = this.storeOptions.findIndex((x) => Number(x.id) === Number(this.form.store_id))
						if (idx >= 0) this.storeIndex = idx
					}
				} catch (e) {}
			},
			async fetchAuditStatus(showLoading = false) {
				const mobile = String(this.form.phone || '').trim()
				if (!mobile) {
					this.audit = { status: '', reason: '' }
					return
				}
				try {
					const res = await userApi.getRegisterStatus({ mobile }, showLoading)
					if (!res) return
					if (res.code !== 0) {
						// 未注册：允许填写资料
						const rejectReason = uni.getStorageSync('wxapp_register_reject_reason')
						this.audit = { status: '', reason: rejectReason ? String(rejectReason) : '' }
						return
					}
					const data = res.data || {}
					const status = String(data.audit_status || '').trim()
					const reason = String(data.audit_reason || '').trim()
					this.audit = { status, reason }

					// 回显：姓名/门店信息（仅在未锁定时回显，避免覆盖用户输入）
					if (status !== 'pending' && status !== 'approved') {
						if (!this.form.name && data.name) this.form.name = String(data.name)

						const storeId = Number(data.store_id || 0) || 0
						if (storeId > 0) {
							this.storeMode = 'select'
							this.form.store_id = storeId
							const idx = this.storeOptions.findIndex((x) => Number(x.id) === storeId)
							if (idx >= 0) this.storeIndex = idx
						} else if (data.store_name) {
							this.storeMode = 'manual'
							this.form.store_id = 0
							this.storeIndex = -1
							this.form.store_name_text = String(data.store_name || '').trim()
							this.form.store_address_text = String(data.store_address || '').trim()
						}
					}

					// 记住拒绝原因（方便从登录页跳转时展示）
					if (reason) {
						try { uni.setStorageSync('wxapp_register_reject_reason', reason) } catch (e) {}
					} else {
						try { uni.removeStorageSync('wxapp_register_reject_reason') } catch (e) {}
					}
				} catch (e) {}
			},
			async handleSubmit() {
				if (this.submitting) return

				if (this.audit.status === 'pending') {
					uni.showToast({ title: '资料审核中，请耐心等待', icon: 'none' })
					return
				}
				if (this.audit.status === 'approved') {
					this.goLogin()
					return
				}

				if (!this.agreed) {
					uni.showToast({
						title: '请先阅读并同意协议',
						icon: 'none'
					})
					return
				}
				if (!this.form.name) {
					uni.showToast({ title: '请填写姓名', icon: 'none' })
					return
				}
				if (!this.form.phone) {
					uni.showToast({ title: '请先在登录页授权手机号', icon: 'none' })
					return
				}

				const storeId = Number(this.form.store_id || 0) || 0
				const storeNameText = String(this.form.store_name_text || '').trim()
				if (storeId <= 0 && !storeNameText) {
					uni.showToast({ title: '请选择门店或填写门店名称', icon: 'none' })
					return
				}

				const phoneCode = uni.getStorageSync('wxapp_register_phone_code')
				if (!phoneCode) {
					uni.showToast({ title: '请返回登录页授权手机号后再提交', icon: 'none' })
					return
				}

				this.submitting = true
				try {
					const payload = {
						phone_code: String(phoneCode),
						name: String(this.form.name).trim(),
						region_province: this.region && this.region[0] ? String(this.region[0]) : '',
						region_city: this.region && this.region[1] ? String(this.region[1]) : '',
						region_district: this.region && this.region[2] ? String(this.region[2]) : '',
						store_id: storeId,
						store_name_text: this.storeMode === 'manual' ? storeNameText : '',
						store_address_text: this.storeMode === 'manual' ? String(this.form.store_address_text || '').trim() : '',
					}
					const res = await userApi.registerApply(payload)
					if (!res) return
					if (res.code !== 0) return
					uni.showToast({ title: '提交成功，请等待审核', icon: 'none' })
					try { uni.setStorageSync('hm_phone', String(this.form.phone)) } catch (e) {}
					try { uni.removeStorageSync('wxapp_register_phone_code') } catch (e) {}
					await this.fetchAuditStatus(false)
				} finally {
					this.submitting = false
				}
			}
		}
	}
</script>

<style lang="scss">
	.register-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	/* TopHeader 左侧按钮 */
	.th-btn {
		width: 80rpx;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		color: #0f172a;
		background: transparent;

		.material-symbols-outlined { font-size: 32rpx; }
		&:active { background: rgba(15, 23, 42, 0.06); }
	}

	.body {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
	}

	.content {
		flex: 1;
		overflow: hidden;
	}

	.content-inner {
		padding: 20rpx 24rpx 0;
		max-width: 920rpx;
		margin: 0 auto;
		box-sizing: border-box;
	}

	/* 顶部欢迎卡片：增加留白与层次 */
	.header-card {
		position: relative;
		overflow: hidden;
		border-radius: 28rpx;
		padding: 28rpx 26rpx;
		background: linear-gradient(135deg, rgba(45, 156, 240, 0.14) 0%, rgba(37, 99, 235, 0.08) 60%, rgba(249, 115, 22, 0.06) 100%);
		border: 1rpx solid rgba(226, 232, 240, 0.9);
		box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
		margin-bottom: 18rpx;

		.header-decor {
			position: absolute;
			right: -60rpx;
			top: -60rpx;
			width: 240rpx;
			height: 240rpx;
			border-radius: 50%;
			background: rgba(45, 156, 240, 0.16);
			filter: blur(30rpx);
		}

		.header-main {
			position: relative;
			z-index: 2;
			display: flex;
			align-items: center;
			gap: 18rpx;
		}

		.brand-icon {
			width: 92rpx;
			height: 92rpx;
			border-radius: 22rpx;
			background: rgba(45, 156, 240, 0.18);
			border: 1rpx solid rgba(45, 156, 240, 0.22);
			display: flex;
			align-items: center;
			justify-content: center;
			color: #2d9cf0;

			.material-symbols-outlined {
				font-size: 48rpx;
				font-variation-settings: 'FILL' 1;
			}
		}

		.title {
			font-size: 40rpx;
			font-weight: 900;
			color: #0f172a;
			line-height: 1.2;
		}

		.subtitle {
			font-size: 24rpx;
			color: #64748b;
			margin-top: 6rpx;
			font-weight: 600;
		}
	}

	/* 审核状态卡片 */
	.audit-card {
		border-radius: 24rpx;
		padding: 22rpx 22rpx;
		border: 1rpx solid rgba(226, 232, 240, 0.9);
		background: #fff;
		box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
		margin-bottom: 18rpx;
	}
	.audit-card.pending { border-color: rgba(245, 158, 11, 0.26); background: rgba(245, 158, 11, 0.06); }
	.audit-card.approved { border-color: rgba(34, 197, 94, 0.26); background: rgba(34, 197, 94, 0.06); }
	.audit-card.rejected { border-color: rgba(239, 68, 68, 0.26); background: rgba(239, 68, 68, 0.06); }

	.audit-top {
		display: flex;
		align-items: flex-start;
		gap: 16rpx;
	}
	.audit-icon {
		font-size: 44rpx;
		color: #0f172a;
		opacity: 0.9;
	}
	.audit-texts { flex: 1; }
	.audit-title { font-size: 30rpx; font-weight: 700; color: #0f172a; }
	.audit-desc { margin-top: 8rpx; font-size: 24rpx; color: #475569; line-height: 1.6; }
	.audit-reason { margin-top: 12rpx; font-size: 24rpx; color: #b91c1c; line-height: 1.6; }
	.audit-actions { margin-top: 14rpx; display: flex; gap: 12rpx; }
	.audit-btn {
		height: 64rpx;
		line-height: 64rpx;
		padding: 0 20rpx;
		border-radius: 16rpx;
		border: 1rpx solid rgba(15, 23, 42, 0.12);
		background: rgba(255, 255, 255, 0.78);
		font-size: 26rpx;
		color: #0f172a;
	}
	.audit-btn.primary {
		border-color: rgba(45, 156, 240, 0.28);
		background: rgba(45, 156, 240, 0.12);
		color: #1a7ab5;
	}

	/* 表单卡片：让内容“不贴边、不拥挤” */
	.form-card {
		background-color: rgba(255, 255, 255, 0.96);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-radius: 28rpx;
		border: 1rpx solid rgba(226, 232, 240, 0.9);
		box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
		padding: 22rpx 20rpx;
	}

	.form {
		.form-item {
			margin-bottom: 18rpx;

			.label {
				font-size: 26rpx;
				font-weight: 500;
				color: #334155;
				margin-bottom: 10rpx;
				display: block;
			}
			
			.label-row {
				display: flex;
				justify-content: space-between;
				align-items: center;
				margin-bottom: 12rpx;
				
				.label-tip {
					font-size: 24rpx;
					color: #4b789b;
				}
			}

			.picker-wrapper {
				display: flex;
				align-items: center;
				gap: 24rpx;
				background-color: #f8fafc;
				padding: 18rpx 24rpx;
				border-radius: 24rpx;
				border: 1px solid rgba(226, 232, 240, 0.95);
				transition: border-color 0.3s;

				.icon-fill {
					color: #2d9cf0;
					font-size: 40rpx;
				}

				.picker-text {
					flex: 1;
					font-size: 32rpx;
					color: #0f172a;
					font-weight: 700;
				}

				.arrow-icon {
					color: #94a3b8;
					font-size: 40rpx;
				}

				&:active {
					border-color: #2d9cf0;
				}
			}

			.input-wrapper {
				position: relative;
				display: flex;
				align-items: center;

				.input-icon {
					position: absolute;
					left: 32rpx;
					color: #4b789b;
					font-size: 40rpx;
					z-index: 1;
				}

				.input {
					width: 100%;
					height: 84rpx;
					background-color: #f8fafc;
					border: 1px solid rgba(226, 232, 240, 0.95);
					border-radius: 24rpx;
					padding-left: 96rpx;
					padding-right: 32rpx;
					font-size: 32rpx;
					color: #0d151c;
					box-sizing: border-box;
				}
				
				.readonly-input {
					pointer-events: none;
				}

				.search-icon {
					position: absolute;
					right: 24rpx;
					color: #4b789b;
					font-size: 36rpx;
				}
			}

			/* 门店选择方式 */
			.store-mode {
				display: flex;
				gap: 12rpx;
				margin-bottom: 14rpx;
			}
			.mode-item {
				flex: 1;
				height: 64rpx;
				border-radius: 18rpx;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 26rpx;
				color: #64748b;
				border: 1rpx solid rgba(226, 232, 240, 0.9);
				background: rgba(248, 250, 252, 0.9);
			}
			.mode-item.active {
				color: #1a7ab5;
				border-color: rgba(45, 156, 240, 0.28);
				background: rgba(45, 156, 240, 0.10);
			}
			.store-addr-tip {
				margin-top: 10rpx;
				font-size: 24rpx;
				color: #64748b;
				line-height: 1.5;
			}
			.mt { margin-top: 14rpx; }
		}
		
		.agreement-row {
			margin-top: 20rpx;
			display: flex;
			align-items: flex-start;
			flex-wrap: wrap;
			gap: 8rpx;
			
			.checkbox-label {
				display: flex;
				align-items: center;
			}

			.agreement-text {
				font-size: 24rpx;
				color: #4b789b;
			}

			.agreement-links {
				display: flex;
				align-items: center;
				flex-wrap: wrap;
				gap: 8rpx;
			}

			.link {
				color: #2d9cf0;
			}

			.sep {
				color: #64748b;
			}

			.disabled {
				opacity: 0.5;
			}
		}
	}

	.footer {
		padding: 20rpx 32rpx calc(env(safe-area-inset-bottom) + 20rpx);
		border-top: 1rpx solid #f1f1f1;
		background-color: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);

		.register-btn {
			width: 100%;
			height: 88rpx;
			background-color: #2d9cf0;
			color: #ffffff;
			font-size: 32rpx;
			font-weight: bold;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			box-shadow: 0 10rpx 30rpx rgba(45, 156, 240, 0.3);
			border: none;

			&::after { border: none; }

			&:active {
				background-color: #1a85d6;
				transform: scale(0.98);
			}
		}
	}

	.placeholder {
		color: rgba(75, 120, 155, 0.6);
	}

	.bottom-spacer {
		height: 28rpx;
	}
</style>
