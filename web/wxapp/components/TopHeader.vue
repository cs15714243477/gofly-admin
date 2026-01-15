<template>
	<view class="th" :style="wrapStyle">
		<view class="th__inner" :style="innerStyle">
			<view v-if="$slots.left" class="th__left">
				<slot name="left"></slot>
			</view>
			<view class="th__center">
				<slot v-if="$slots.default"></slot>
				<text v-else class="th__title">{{ title }}</text>
			</view>
			<view class="th__right">
				<slot name="right"></slot>
			</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'TopHeader',
		props: {
			title: {
				type: String,
				default: ''
			}
		},
		data() {
			return {
				padPx: 0,
				rightSafePx: 0,
				statusBarPx: 0,
				navBarPx: 0
			}
		},
		computed: {
			wrapStyle() {
				// 小程序端：把状态栏高度算进去，避免内容贴顶；其他端保持原样
				if (!this.statusBarPx || !this.navBarPx) return {}
				return {
					paddingTop: this.statusBarPx + 'px',
					height: this.statusBarPx + this.navBarPx + 'px'
				}
			},
			innerStyle() {
				// 默认走 rpx padding；仅 MP-WEIXIN 运行时用 px 覆盖，以便精确避开胶囊按钮
				const style = {}
				if (this.padPx) {
					style.paddingLeft = this.padPx + 'px'
					style.paddingRight = this.padPx + this.rightSafePx + 'px'
				}
				// 小程序端：导航栏高度跟随胶囊按钮高度（更接近原生）
				if (this.navBarPx) style.height = this.navBarPx + 'px'
				return style
			}
		},
		mounted() {
			// #ifdef MP-WEIXIN
			try {
				const info = uni.getSystemInfoSync()
				const w = Number(info.windowWidth || 375)
				const statusBar = Number(info.statusBarHeight || 0)
				// 32rpx -> px（750 设计稿）
				const pad = Math.round((w * 32) / 750)
				let rightSafe = 0
				let navBarH = 0
				if (typeof wx !== 'undefined' && wx.getMenuButtonBoundingClientRect) {
					const rect = wx.getMenuButtonBoundingClientRect()
					// 预留胶囊区域：从胶囊左边到屏幕右侧的宽度
					rightSafe = Math.max(0, w - Number(rect.left || 0))
					// 导航栏高度：更接近微信原生的计算方式
					navBarH = Math.max(44, (Number(rect.top || 0) - statusBar) * 2 + Number(rect.height || 0)) || 44
				}
				this.padPx = pad
				this.rightSafePx = rightSafe
				// 与胶囊按钮垂直居中对齐：使用原生推荐计算，不额外下移
				this.statusBarPx = statusBar
				this.navBarPx = navBarH
			} catch (e) {
				// ignore
			}
			// #endif
		}
	}
</script>

<style lang="scss" scoped>
	.th {
		width: 100%;
		background: rgba(255, 255, 255, 0.92);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1rpx solid rgba(226, 232, 240, 0.9);
		box-sizing: border-box;
		box-shadow: 0 6rpx 18rpx rgba(15, 23, 42, 0.04);
	}

	.th__inner {
		height: 96rpx;
		padding: 0 32rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
		position: relative;
		box-sizing: border-box;
	}

	.th__center {
		flex: 1;
		display: flex;
		align-items: center;
		overflow: hidden;
	}

	.th__title {
		text-align: left;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-size: 36rpx;
		font-weight: 800;
		color: #0f172a;
		letter-spacing: 0.6rpx;
		line-height: 1;
		margin: 0 24rpx 0 0;
	}

	.th__left {
		min-width: 96rpx;
		height: 96rpx;
		display: flex;
		align-items: center;
		justify-content: flex-start;
	}

	.th__right {
		min-width: 96rpx;
		height: 96rpx;
		display: flex;
		align-items: center;
		justify-content: flex-end;
	}
</style>

