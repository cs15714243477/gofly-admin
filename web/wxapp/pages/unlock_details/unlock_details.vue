<template>
	<view class="unlock-container">
		<!-- 顶部悬浮操作 -->
		<view class="float-header" :style="{ paddingTop: headerTop + 'px' }">
			<view class="header-left">
				<view class="circle-btn" @click="goBack">
					<text class="material-symbols-outlined">arrow_back</text>
				</view>
			</view>
		</view>

		<view class="main-content">
			<!-- 顶部封面 -->
			<view class="cover-section">
				<image
					v-if="coverImage"
					:src="coverImage"
					mode="aspectFill"
					class="cover-image"
				></image>
				<view v-else class="cover-image cover-image-empty">
					<text class="material-symbols-outlined">image</text>
				</view>
				<view class="cover-mask"></view>
				<view class="cover-info">
					<view class="verified-tag">
						<text class="material-symbols-outlined tag-icon">verified_user</text>
						<text>已核验</text>
					</view>
					<view class="building-name">房源ID：{{ propertyId }}</view>
					<view class="scene-link">
						<text class="material-symbols-outlined scene-icon">image</text>
						<text>房源实景图</text>
					</view>
				</view>
			</view>

			<!-- 核心操作区 -->
			<view class="card-section">
				<view class="password-card">
					<view class="card-bg-decor"></view>
					<text class="card-label">动态密码</text>
					<view class="password-row">
						<text class="password-text">{{ passcodeDisplay }}</text>
						<view class="copy-btn" @click="copyPassword">
							<text class="material-symbols-outlined">content_copy</text>
						</view>
					</view>
					
					<!-- 倒计时 -->
					<view class="timer-box">
						<view class="timer-left">
							<view class="progress-circle">
								<view class="progress-mask"></view>
								<text class="material-symbols-outlined timer-icon">timer</text>
							</view>
							<view class="timer-info">
								<text class="timer-title">本次密码有效时间</text>
								<text class="timer-val">{{ remainLabel }}</text>
							</view>
						</view>
						<button class="extend-btn" @click="extendPasscode">延时</button>
					</view>

					<!-- 时间详情 -->
					<view class="time-details">
						<view class="time-item">
							<text class="time-lab">申请时间</text>
							<text class="time-val">{{ formatTime(startDate) }}</text>
						</view>
						<view class="time-divider"></view>
						<view class="time-item">
							<text class="time-lab">截止时间</text>
							<text class="time-val accent">{{ formatTime(endDate) }}</text>
							<view class="active-dot"></view>
						</view>
					</view>
				</view>

				<!-- 开锁注意 -->
				<view class="notice-card">
					<view class="notice-icon-box">
						<text class="material-symbols-outlined">touch_app</text>
					</view>
					<view class="notice-content">
						<view class="notice-title">开锁注意</view>
						<view class="notice-text">
							按<text class="keyboard-btn">#</text>号键激活，输入密码后即可开锁。请确保操作时面板已亮起。
						</view>
						<view class="notice-tip">
							若门锁无反应，请尝试使用USB应急电源接口供电。
						</view>
						<view class="notice-tip">
							提示：部分门锁需要在门锁附近完成一次蓝牙连接/蓝牙开锁后，云端数据（如状态/记录）才会同步更新。
						</view>
					</view>
				</view>

				<!-- 操作按钮组 -->
				<view class="btn-row">
					<button class="action-btn" @click="refreshPasscode">
						<text class="material-symbols-outlined btn-icon rotate">refresh</text>
						<text>换新密码</text>
					</button>
					<button class="action-btn">
						<text class="material-symbols-outlined btn-icon">report_problem</text>
						<text>故障上报</text>
					</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import propertyApi from '@/api/property'
	import ttlockApi from '@/api/ttlock'

	export default {
		data() {
			return {
				statusBarHeight: 0,
				headerTop: 0,
				propertyId: 0,
				coverImage: '',
				passcode: '',
				startDate: 0,
				endDate: 0,
				remainSeconds: 0,
				timer: null,
			}
		},
		computed: {
			passcodeDisplay() {
				const raw = String(this.passcode || '').trim()
				if (!raw) return '——'
				const pure = raw.replace(/\s+/g, '')
				if (/^\d{6}$/.test(pure)) return `${pure.slice(0, 3)} ${pure.slice(3)}`
				return raw
			},
			remainLabel() {
				if (!this.endDate) return '—'
				const s = Number(this.remainSeconds || 0)
				if (!Number.isFinite(s) || s <= 0) return '已过期'
				const h = Math.floor(s / 3600)
				const m = Math.floor((s % 3600) / 60)
				const ss = Math.floor(s % 60)
				return `${h}小时${m}分${ss}秒`
			},
		},
		onLoad(options) {
			const pid = Number((options && (options.property_id || options.id)) || 0)
			this.propertyId = Number.isFinite(pid) ? pid : 0

			const info = uni.getSystemInfoSync()
			this.statusBarHeight = info.statusBarHeight
			// #ifdef MP-WEIXIN
			try {
				if (typeof wx !== 'undefined' && wx.getMenuButtonBoundingClientRect) {
					const rect = wx.getMenuButtonBoundingClientRect()
					this.headerTop = Number(rect.top || this.statusBarHeight) + 6
				} else {
					this.headerTop = this.statusBarHeight + 12
				}
			} catch (e) {
				this.headerTop = this.statusBarHeight + 12
			}
			// #endif
			// #ifndef MP-WEIXIN
			this.headerTop = this.statusBarHeight + 12
			// #endif

			this.loadPasscode()
			this.loadPropertyCover()
		},
		onUnload() {
			this.clearTimer()
		},
		methods: {
			normalizeImage(url) {
				const imageUrl = String(url || '').trim()
				if (!imageUrl) return ''
				if (imageUrl.indexOf('/static/images/') === 0) return ''
				return imageUrl
			},
			async loadPropertyCover() {
				if (!this.propertyId) return
				try {
					const res = await propertyApi.getDetail({ id: this.propertyId, public: 0 })
					if (!res || res.code !== 0) return
					const data = res.data || {}
					const p = data.property || {}
					const videoPoster = this.normalizeImage(p.cover_image)
					const videoUrl = String(p.video_url || '').trim()
					const imgs = Array.isArray(data.images)
						? data.images
						: Array.isArray(p.images)
							? p.images
							: []
					const normalizedImgs = imgs.map((u) => this.normalizeImage(u)).filter(Boolean)
					const coverCandidates = []
					// 与房源详情轮播图逻辑对齐：有视频时优先封面图，其次首张图片
					if (videoUrl && videoPoster) coverCandidates.push(videoPoster)
					if (normalizedImgs.length) coverCandidates.push(normalizedImgs[0])
					if (!coverCandidates.length && videoPoster) coverCandidates.push(videoPoster)
					this.coverImage = coverCandidates[0] || ''
				} catch (e) {}
			},
			goBack() {
				uni.navigateBack()
			},
			clearTimer() {
				if (this.timer) {
					clearInterval(this.timer)
					this.timer = null
				}
			},
			startTimer() {
				this.clearTimer()
				this.timer = setInterval(() => {
					if (!this.endDate) return
					const remain = Math.floor((this.endDate - Date.now()) / 1000)
					this.remainSeconds = remain > 0 ? remain : 0
				}, 1000)
			},
			formatTime(ts) {
				const t = Number(ts || 0)
				if (!t) return '--:--'
				const d = new Date(t)
				const hh = String(d.getHours()).padStart(2, '0')
				const mm = String(d.getMinutes()).padStart(2, '0')
				return `${hh}:${mm}`
			},
			async loadPasscode() {
				if (!this.propertyId) {
					uni.showToast({ title: '缺少房源ID，无法获取密码', icon: 'none' })
					return
				}
				const res = await ttlockApi.getPropertyPasscode({ property_id: this.propertyId })
				if (!res || res.code !== 0) return
				const data = res.data || {}

				this.passcode = String(data.keyboardPwd || data.keyboard_pwd || '').trim()
				this.startDate = Number(data.startDate || data.start_date || 0) || 0
				this.endDate = Number(data.endDate || data.end_date || 0) || 0
				this.remainSeconds = this.endDate ? Math.max(0, Math.floor((this.endDate - Date.now()) / 1000)) : 0
				this.startTimer()
			},
			async refreshPasscode() {
				const prev = String(this.passcode || '').trim()
				await this.loadPasscode()
				const next = String(this.passcode || '').trim()
				if (next && next !== prev) {
					uni.showToast({ title: '已更换新密码', icon: 'success' })
				} else if (next) {
					uni.showToast({ title: '密码已刷新', icon: 'none' })
				}
			},
			extendPasscode() {
				uni.showToast({ title: '暂不支持延时，可点击“换新密码”重新生成', icon: 'none' })
			},
			copyPassword() {
				const raw = String(this.passcode || '').trim()
				const pure = raw.replace(/\s+/g, '')
				if (!pure) {
					uni.showToast({ title: '暂无可复制的密码', icon: 'none' })
					return
				}
				uni.setClipboardData({
					data: pure,
					success: () => {
						uni.showToast({
							title: '复制成功',
							icon: 'success'
						})
					}
				})
			}
		}
	}
</script>

<style lang="scss">
	.unlock-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.float-header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 100;
		display: flex;
		justify-content: flex-start;
		padding: 0 32rpx;
		pointer-events: none;
		padding-bottom: 8rpx;

		.header-left {
			display: flex;
			gap: 24rpx;
			pointer-events: auto;
		}

		.circle-btn {
			width: 80rpx;
			height: 80rpx;
			background-color: rgba(0, 0, 0, 0.2);
			backdrop-filter: blur(10px);
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			color: #ffffff;
			
			.material-symbols-outlined {
				font-size: 40rpx;
			}
		}
	}

	.main-content {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	.cover-section {
		width: 100%;
		height: 480rpx;
		position: relative;
		
		.cover-image {
			width: 100%;
			height: 100%;
		}
		.cover-image-empty {
			background: linear-gradient(135deg, #0ea5e9, #1e3a8a);
			display: flex;
			align-items: center;
			justify-content: center;
			.material-symbols-outlined {
				font-size: 72rpx;
				color: rgba(255, 255, 255, 0.45);
			}
		}
		
		.cover-mask {
			position: absolute;
			inset: 0;
			background: linear-gradient(to bottom, rgba(0,0,0,0.4) 0%, transparent 40%, rgba(0,0,0,0.8) 100%);
		}
		
		.cover-info {
			position: absolute;
			bottom: 72rpx;
			left: 32rpx;
			right: 32rpx;
			display: flex;
			flex-direction: column;
			gap: 12rpx;
			color: #ffffff;
			
			.verified-tag {
				display: flex;
				align-items: center;
				gap: 8rpx;
				background-color: rgba(45, 156, 240, 0.9);
				padding: 4rpx 16rpx;
				border-radius: 40rpx;
				font-size: 20rpx;
				font-weight: bold;
				align-self: flex-start;
				
				.tag-icon { font-size: 24rpx; font-variation-settings: 'FILL' 1; }
			}
			
			.building-name {
				font-size: 40rpx;
				font-weight: bold;
				text-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.5);
			}
			
			.scene-link {
				display: flex;
				align-items: center;
				gap: 8rpx;
				font-size: 24rpx;
				color: rgba(255, 255, 255, 0.9);
				
				.scene-icon { font-size: 28rpx; }
			}
		}
	}

	.card-section {
		padding: 0 24rpx 24rpx;
		margin-top: -48rpx;
		position: relative;
		z-index: 20;
		display: flex;
		flex-direction: column;
		gap: 20rpx;
		flex: 1;
		min-height: 0;
	}

	.password-card {
		background-color: #ffffff;
		border-radius: 40rpx;
		padding: 36rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		position: relative;
		overflow: hidden;
		box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.5);

		.card-bg-decor {
			position: absolute;
			top: 0;
			left: 50%;
			transform: translateX(-50%);
			width: 75%;
			height: 160rpx;
			background-color: rgba(245, 158, 11, 0.05);
			border-radius: 50%;
			filter: blur(60rpx);
		}

		.card-label {
			font-size: 24rpx;
			font-weight: bold;
			color: #94a3b8;
			letter-spacing: 4rpx;
			margin-bottom: 20rpx;
		}

		.password-row {
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 16rpx;
			margin-bottom: 24rpx;
			position: relative;
			width: 100%;

			.password-text {
				font-size: 84rpx;
				font-weight: 900;
				color: #f59e0b;
				letter-spacing: 12rpx;
				font-family: 'Inter', sans-serif;
			}

			.copy-btn {
				position: absolute;
				right: -24rpx;
				padding: 16rpx;
				color: #cbd5e1;
				
				&:active { color: #f59e0b; }
			}
		}

		.timer-box {
			width: 100%;
			background-color: #f8fafc;
			border-radius: 24rpx;
			padding: 24rpx;
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 20rpx;
			border: 1px solid #f1f5f9;

			.timer-left {
				display: flex;
				align-items: center;
				gap: 24rpx;
			}

			.progress-circle {
				width: 72rpx;
				height: 72rpx;
				border-radius: 50%;
				border: 6rpx solid #e2e8f0;
				position: relative;
				display: flex;
				align-items: center;
				justify-content: center;

				.progress-mask {
					position: absolute;
					width: 100%;
					height: 100%;
					border-radius: 50%;
					border: 6rpx solid #f59e0b;
					border-bottom-color: transparent;
					border-left-color: transparent;
					transform: rotate(45deg);
				}

				.timer-icon { color: #f59e0b; font-size: 36rpx; }
			}

			.timer-info {
				display: flex;
				flex-direction: column;
				gap: 4rpx;

				.timer-title { font-size: 28rpx; font-weight: bold; color: #0f172a; }
				.timer-val { font-size: 24rpx; color: #64748b; font-weight: 500; }
			}

			.extend-btn {
				height: 52rpx;
				padding: 0 24rpx;
				background-color: rgba(45, 156, 240, 0.1);
				color: #2d9cf0;
				font-size: 24rpx;
				font-weight: bold;
				border-radius: 12rpx;
				border: none;
				
				&::after { border: none; }
			}
		}

		.time-details {
			width: 100%;
			display: flex;
			background-color: #f8fafc;
			border-radius: 20rpx;
			overflow: hidden;
			border: 1px solid #f1f5f9;

			.time-item {
				flex: 1;
				padding: 20rpx;
				display: flex;
				flex-direction: column;
				align-items: center;
				gap: 4rpx;
				position: relative;
			}

			.time-divider {
				width: 1px;
				height: 56rpx;
				background-color: #e2e8f0;
				align-self: center;
			}

			.time-lab { font-size: 20rpx; font-weight: bold; color: #94a3b8; }
			.time-val { font-size: 28rpx; font-weight: bold; color: #334155; }
			.accent { color: #f59e0b; }
			
			.active-dot {
				position: absolute;
				top: 16rpx;
				right: 16rpx;
				width: 12rpx;
				height: 12rpx;
				background-color: #f59e0b;
				border-radius: 50%;
			}
		}
	}

	.notice-card {
		background-color: #ffffff;
		border-radius: 24rpx;
		padding: 28rpx;
		display: flex;
		gap: 20rpx;
		border: 1px solid #f1f5f9;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.02);

		.notice-icon-box {
			width: 76rpx;
			height: 76rpx;
			background-color: rgba(45, 156, 240, 0.1);
			color: #2d9cf0;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			flex-shrink: 0;
			box-shadow: 0 0 20rpx rgba(45, 156, 240, 0.1);
			
			.material-symbols-outlined { font-size: 40rpx; }
		}

		.notice-content {
			flex: 1;
			display: flex;
			flex-direction: column;
			gap: 12rpx;

			.notice-title { font-size: 28rpx; font-weight: bold; color: #1e293b; }
			
			.notice-text {
				font-size: 22rpx;
				color: #64748b;
				line-height: 1.5;
				
				.keyboard-btn {
					background-color: #f1f5f9;
					padding: 2rpx 12rpx;
					border-radius: 8rpx;
					font-weight: bold;
					color: #1e293b;
					margin: 0 8rpx;
				}
			}

			.notice-tip {
				font-size: 22rpx;
				color: #94a3b8;
				font-style: italic;
				padding-top: 12rpx;
				border-top: 1rpx dashed #e2e8f0;
			}
		}
	}

	.btn-row {
		display: flex;
		gap: 20rpx;
		padding-top: 8rpx;

		.action-btn {
			flex: 1;
			height: 88rpx;
			background-color: #ffffff;
			border: 1px solid #e2e8f0;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 16rpx;
			font-size: 26rpx;
			font-weight: bold;
			color: #475569;
			box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.02);
			
			&::after { border: none; }

			.btn-icon { font-size: 36rpx; color: #94a3b8; }
			
			&:active {
				border-color: #2d9cf0;
				color: #2d9cf0;
				.btn-icon { color: #2d9cf0; }
			}
			
			.rotate {
				transition: transform 0.5s;
			}
			
			&:active .rotate { transform: rotate(180deg); }
		}
	}
</style>
