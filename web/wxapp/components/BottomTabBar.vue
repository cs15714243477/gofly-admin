<template>
	<view class="btb">
		<view class="btb__inner">
		<view
			v-for="item in items"
			:key="item.key"
			class="btb__item"
			:class="{ 'is-active': active === item.key }"
			@click="onTap(item)"
		>
			<view class="btb__icon-wrap" :class="{ 'is-active': active === item.key }">
				<text class="material-symbols-outlined btb__icon" :class="{ fill: active === item.key }">{{ item.icon }}</text>
			</view>
			<text class="btb__text">{{ item.label }}</text>
		</view>
		</view>
	</view>
</template>

<script>
	export default {
		name: 'BottomTabBar',
		props: {
			// home | property | marketing | me
			active: {
				type: String,
				default: 'home'
			}
		},
		data() {
			return {
				items: [
					{ key: 'property', label: '房源', icon: 'apartment', url: '/pages/property_list/property_list' },
					{ key: 'home', label: '推荐', icon: 'home', url: '/pages/home/home' },
					{ key: 'marketing', label: '获客', icon: 'campaign', url: '/pages/my_business_card/my_business_card' },
					{ key: 'me', label: '我的', icon: 'person', url: '/pages/agent_workbench_home/agent_workbench_home' }
				]
			}
		},
		methods: {
			onTap(item) {
				if (!item || !item.url) return
				uni.reLaunch({ url: item.url })
			}
		}
	}
</script>

<style lang="scss" scoped>
	.btb {
		position: relative;
		z-index: 30;
		padding: 10rpx 16rpx calc(env(safe-area-inset-bottom) + 12rpx);
		background: transparent;
	}

	.btb__inner {
		background: rgba(255, 255, 255, 0.92);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border: 1rpx solid rgba(226, 232, 240, 0.9);
		border-radius: 28rpx;
		box-shadow: 0 -10rpx 24rpx rgba(15, 23, 42, 0.06);
		padding: 18rpx 18rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.btb__item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8rpx;
		color: #94a3b8;
		flex: 1;
		min-width: 0;

		&.is-active {
			color: #2563eb;
		}

		&:active .btb__icon-wrap {
			transform: scale(1.1);
		}
	}

	.btb__icon-wrap {
		width: 88rpx;
		height: 56rpx;
		border-radius: 28rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: transform 0.18s ease, background-color 0.18s ease;

		&.is-active {
			background: rgba(37, 99, 235, 0.12);
			border: 1rpx solid rgba(37, 99, 235, 0.18);
		}
	}

	.btb__icon {
		font-size: 48rpx;
		transition: transform 0.2s;

		&.fill {
			font-variation-settings: 'FILL' 1;
		}
	}

	.btb__text {
		font-size: 20rpx;
		font-weight: 700;
		letter-spacing: 0.4rpx;
		line-height: 1;
	}
</style>

