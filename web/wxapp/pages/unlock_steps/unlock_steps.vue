<template>
	<view class="steps-container">
		<!-- 顶部导航 -->
		<TopHeader title="开锁步骤">
			<template #left>
				<view class="th-btn" @click="goBack">
					<text class="material-symbols-outlined">arrow_back_ios_new</text>
				</view>
			</template>
		</TopHeader>

		<view class="main-content">
			<!-- 提示卡片 -->
			<view class="tip-section">
				<view class="tip-card">
					<text class="material-symbols-outlined tip-icon">near_me</text>
					<text class="tip-text">请确保您已站在门锁附近 (1-2米范围内)，并保持手机蓝牙已开启。</text>
				</view>
			</view>

			<!-- 步骤流程 -->
			<view class="steps-section">
				<view class="steps-line"></view>
				
				<view class="step-item group">
					<view class="step-icon-box primary">
						<text class="material-symbols-outlined">bluetooth</text>
					</view>
					<view class="step-content">
						<text class="step-tag primary">步骤 1</text>
						<text class="step-title">第一步. 打开手机蓝牙</text>
						<text class="step-desc">打开手机系统设置，确保蓝牙功能处于开启状态。</text>
					</view>
				</view>

				<view class="step-item group">
					<view class="step-icon-box accent">
						<text class="material-symbols-outlined">touch_app</text>
					</view>
					<view class="step-content">
						<text class="step-tag accent">步骤 2</text>
						<text class="step-title">第二步. 唤醒门锁</text>
						<text class="step-desc">用手触碰门锁数字区，点亮面板屏幕以建立连接。</text>
					</view>
				</view>

				<view class="step-item group">
					<view class="step-icon-box primary">
						<text class="material-symbols-outlined">phonelink_lock</text>
					</view>
					<view class="step-content">
						<text class="step-tag primary">步骤 3</text>
						<text class="step-title">第三步. 点击下方“开锁”</text>
						<text class="step-desc">点击页面底部的蓝色开锁按钮进行开锁。</text>
					</view>
				</view>
			</view>

			<!-- 状态指示 -->
			<view class="status-indicator">
				<view class="status-tag">
					<view class="dot-box">
						<view class="dot-ping"></view>
						<view class="dot-inner"></view>
					</view>
					<text class="status-text">设备连接就绪</text>
				</view>
			</view>
			
		</view>

		<!-- 底部操作 -->
		<view class="bottom-bar">
			<view class="bar-content">
				<button class="bluetooth-btn" @click="bluetoothUnlock">
					<text class="material-symbols-outlined bt-icon pulse">bluetooth_connected</text>
					<text>蓝牙开锁</text>
				</button>
				<view class="bar-tip">若蓝牙开锁失败，点击下方按钮使用密码开锁</view>
				<button class="password-btn" @click="passwordUnlock">
					<text class="material-symbols-outlined btn-icon">lock_open</text>
					<text>密码开锁</text>
				</button>
				<view class="support-link">
					<text class="material-symbols-outlined support-icon">support_agent</text>
					<text>遇到问题？联系客服</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'

	export default {
		components: { TopHeader },
		data() {
			return {}
		},
		methods: {
			goBack() {
				uni.navigateBack()
			},
			bluetoothUnlock() {
				console.log('蓝牙开锁')
			},
			passwordUnlock() {
				uni.navigateTo({
					url: '/pages/unlock_details/unlock_details'
				})
			}
		}
	}
</script>

<style lang="scss">
	.steps-container {
		height: 100vh;
		background-color: #ffffff;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	/* 顶部栏已统一为 TopHeader，这里补充左右 slot 内部样式 */
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

	.main-content {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	.tip-section {
		padding: 28rpx 32rpx 20rpx;
	}

	.tip-card {
		background-color: #f0f9ff;
		border: 1px solid #e0f2fe;
		border-radius: 32rpx;
		padding: 32rpx;
		display: flex;
		gap: 24rpx;
		align-items: flex-start;
		box-shadow: 0 4rpx 10rpx rgba(45, 156, 240, 0.05);

		.tip-icon { color: #2d9cf0; font-size: 40rpx; }
		.tip-text { font-size: 28rpx; color: #334155; font-weight: 500; line-height: 1.6; }
	}

	.steps-section {
		padding: 0 32rpx;
		position: relative;

		.steps-line {
			position: absolute;
			left: 70rpx;
			top: 32rpx;
			bottom: 48rpx;
			width: 4rpx;
			background-color: #f1f5f9;
			z-index: 1;
		}

		.step-item {
			display: flex;
			gap: 28rpx;
			padding-bottom: 40rpx;
			position: relative;
			z-index: 10;
			
			&:last-child { padding-bottom: 0; }

			.step-icon-box {
				width: 72rpx;
				height: 72rpx;
				background-color: #ffffff;
				border-radius: 50%;
				display: flex;
				align-items: center;
				justify-content: center;
				border-width: 4rpx;
				border-style: solid;
				box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
				transition: transform 0.3s;
				
				&.primary { border-color: #2d9cf0; color: #2d9cf0; }
				&.accent { border-color: #f97316; color: #f97316; }

				.material-symbols-outlined { font-size: 40rpx; }
			}

			.step-content {
				flex: 1;
				display: flex;
				flex-direction: column;
				gap: 8rpx;
				padding-top: 8rpx;

				.step-tag {
					font-size: 20rpx;
					font-weight: bold;
					padding: 4rpx 16rpx;
					border-radius: 8rpx;
					align-self: flex-start;
					
					&.primary { background-color: #eff6ff; color: #2d9cf0; }
					&.accent { background-color: #fff7ed; color: #f97316; }
				}

				.step-title { font-size: 30rpx; font-weight: bold; color: #0f172a; }
				.step-desc { font-size: 26rpx; color: #64748b; line-height: 1.5; }
			}
			
			&:active .step-icon-box { transform: scale(1.1); }
		}
	}

	.status-indicator {
		margin-top: 28rpx;
		margin-bottom: 16rpx;
		display: flex;
		justify-content: center;
		opacity: 0;
		animation: fadeIn 0.5s ease-out 0.5s forwards;

		.status-tag {
			display: flex;
			align-items: center;
			gap: 16rpx;
			background-color: #ecfdf5;
			border: 1px solid #d1fae5;
			padding: 16rpx 32rpx;
			border-radius: 40rpx;
			backdrop-filter: blur(10px);
		}

		.dot-box {
			width: 20rpx;
			height: 20rpx;
			position: relative;
			
			.dot-ping {
				position: absolute;
				width: 100%;
				height: 100%;
				background-color: #10b981;
				border-radius: 50%;
				animation: ping 1.5s cubic-bezier(0, 0, 0.2, 1) infinite;
			}
			
			.dot-inner {
				width: 100%;
				height: 100%;
				background-color: #10b981;
				border-radius: 50%;
				position: relative;
				z-index: 10;
			}
		}

		.status-text { font-size: 24rpx; font-weight: 600; color: #065f46; }
	}

	.bottom-bar {
		background-color: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);
		border-top: 1rpx solid #f1f5f9;
		padding: 24rpx 32rpx calc(env(safe-area-inset-bottom) + 24rpx);
		z-index: 100;

		.bar-content {
			display: flex;
			flex-direction: column;
			gap: 24rpx;
			max-width: 900rpx;
			margin: 0 auto;
		}

		.bluetooth-btn {
			height: 96rpx;
			background-color: #2d9cf0;
			color: #ffffff;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 24rpx;
			font-size: 36rpx;
			font-weight: bold;
			box-shadow: 0 10rpx 25rpx rgba(45, 156, 240, 0.3);
			border: none;
			position: relative;
			overflow: hidden;
			
			&::after { border: none; }

			.bt-icon { font-size: 48rpx; }
			
			&:active { background-color: #1d82cc; transform: scale(0.98); }
		}

		.bar-tip { font-size: 24rpx; color: #94a3b8; text-align: center; }

		.password-btn {
			height: 84rpx;
			background-color: #f8fafc;
			border: 1px solid #e2e8f0;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 16rpx;
			font-size: 28rpx;
			font-weight: bold;
			color: #475569;
			
			&::after { border: none; }
			&:active { background-color: #f1f5f9; transform: scale(0.98); }

			.btn-icon { font-size: 36rpx; }
		}

		.support-link {
			margin-top: 16rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 8rpx;
			font-size: 24rpx;
			color: #94a3b8;
			
			.support-icon { font-size: 28rpx; }
			&:active { color: #2d9cf0; }
		}
	}

	@keyframes ping {
		0% { transform: scale(1); opacity: 1; }
		75%, 100% { transform: scale(2.5); opacity: 0; }
	}

	@keyframes fadeIn {
		from { opacity: 0; transform: translateY(20rpx); }
		to { opacity: 1; transform: translateY(0); }
	}
	
	.pulse {
		animation: pulse-icon 2s infinite;
	}
	
	@keyframes pulse-icon {
		0%, 100% { opacity: 1; transform: scale(1); }
		50% { opacity: 0.8; transform: scale(1.1); }
	}
</style>
