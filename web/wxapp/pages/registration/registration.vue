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
									<input v-model="form.name" class="input" type="text" placeholder="请输入姓名" placeholder-class="placeholder" />
								</view>
							</view>

							<!-- 手机号（微信一键登录后自动带入） -->
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
								<view class="input-wrapper" @click="selectStore">
									<text class="material-symbols-outlined input-icon">storefront</text>
									<input class="input readonly-input" type="text" placeholder="请选择或搜索门店" readonly :value="form.store" />
									<text class="material-symbols-outlined search-icon">search</text>
								</view>
							</view>

							<!-- 协议 -->
							<view class="agreement-row">
								<checkbox-group @change="agreementChange">
									<label class="checkbox-label">
										<checkbox value="agree" checked="true" color="#2d9cf0" style="transform:scale(0.7)" />
										<text class="agreement-text">我已阅读并同意 <text class="link">《秒卖房用户使用协议》</text></text>
									</label>
								</checkbox-group>
							</view>
						</view>
					</view>

					<view class="bottom-spacer"></view>
				</view>
			</scroll-view>

			<!-- 底部固定按钮（避免撑高页面出现滚动） -->
			<view class="footer">
				<button class="register-btn" @click="handleSubmit">{{ submitText }}</button>
			</view>
		</view>
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'

	export default {
		components: { TopHeader },
		onLoad(options) {
			this.mode = (options && options.mode) || 'complete'
			// 取登录页保存的手机号（占位/或后端回填的真实号）
			const cachedPhone = uni.getStorageSync('hm_phone')
			if (cachedPhone) this.form.phone = cachedPhone
		},
		data() {
			return {
				agreed: true,
				mode: 'complete', // complete | register（预留）
				region: ['北京市', '北京市', '朝阳区'],
				form: {
					name: '',
					phone: '',
					store: '链家地产 (朝阳大悦城店)'
				}
			}
		},
		computed: {
			pageTitle() {
				return this.mode === 'complete' ? '完善信息' : '注册'
			},
			headerTitle() {
				return this.mode === 'complete' ? '完善资料' : '欢迎加入经纪人'
			},
			headerSubtitle() {
				return this.mode === 'complete' ? '首次登录请先补全信息，完成后进入系统' : '请填写以下信息以完成注册'
			},
			submitText() {
				return this.mode === 'complete' ? '完成并进入' : '新用户注册'
			},
			regionText() {
				const r = Array.isArray(this.region) ? this.region : []
				return r.filter(Boolean).join(' ')
			}
		},
		methods: {
			goBack() {
				uni.navigateBack()
			},
			agreementChange(e) {
				this.agreed = e.detail.value.length > 0
			},
			onRegionChange(e) {
				const value = e && e.detail && e.detail.value
				if (Array.isArray(value) && value.length) this.region = value
			},
			selectStore() {
				console.log('选择门店')
			},
			handleSubmit() {
				if (!this.agreed) {
					uni.showToast({
						title: '请先同意协议',
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
				// 标记已完成首次资料
				uni.setStorageSync('hm_profile_completed', true)
				uni.reLaunch({
					url: '/pages/home/home'
				})
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
		}
		
		.agreement-row {
			margin-top: 20rpx;
			display: flex;
			align-items: flex-start;
			
			.checkbox-label {
				display: flex;
				align-items: center;
			}

			.agreement-text {
				font-size: 24rpx;
				color: #4b789b;
				
				.link {
					color: #2d9cf0;
				}
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
