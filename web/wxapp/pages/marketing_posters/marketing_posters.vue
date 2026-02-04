<template>
	<view class="poster-container">
		<!-- 顶部导航：使用 TopHeader 统一计算状态栏/胶囊高度，避免贴顶且不遮挡 -->
		<view class="header">
			<TopHeader title="获客海报">
				<template #right>
					<view class="search-box">
						<text class="material-symbols-outlined search-icon">search</text>
						<input class="search-input" type="text" placeholder="搜索海报主题 (例如: 端午节)" placeholder-class="placeholder" />
					</view>
				</template>
			</TopHeader>
		</view>

		<view class="main-body">
			<!-- 侧边栏导航 -->
			<scroll-view scroll-y="true" class="sidebar">
				<view class="category-list">
					<view class="cat-item" v-for="(cat, idx) in categories" :key="idx" :class="{ active: currentCat === idx }" @click="currentCat = idx">
						<view class="active-bar" v-if="currentCat === idx"></view>
						<text class="cat-name">{{ cat }}</text>
					</view>
				</view>
			</scroll-view>

			<!-- 瀑布流内容 -->
			<scroll-view scroll-y="true" class="poster-content">
				<view class="waterfall">
					<view class="poster-column">
						<view class="poster-card" v-for="(poster, pIdx) in leftPosters" :key="pIdx">
							<view class="poster-img-box" :style="{ aspectRatio: poster.aspect }">
								<image v-if="poster.image" :src="poster.image" mode="aspectFill" class="poster-image"></image>
								<view v-else class="poster-image poster-image-empty">
									<text class="material-symbols-outlined">image</text>
								</view>
								<view class="poster-tag" :class="poster.tagType" v-if="poster.tag">{{ poster.tag }}</view>
							</view>
							<view class="poster-info">
								<text class="poster-title">{{ poster.title }}</text>
								<text class="poster-desc">{{ poster.desc }}</text>
							</view>
						</view>
					</view>
					<view class="poster-column">
						<view class="poster-card" v-for="(poster, pIdx) in rightPosters" :key="pIdx">
							<view class="poster-img-box" :style="{ aspectRatio: poster.aspect }">
								<image v-if="poster.image" :src="poster.image" mode="aspectFill" class="poster-image"></image>
								<view v-else class="poster-image poster-image-empty">
									<text class="material-symbols-outlined">image</text>
								</view>
								<view class="poster-tag" :class="poster.tagType" v-if="poster.tag">{{ poster.tag }}</view>
							</view>
							<view class="poster-info">
								<text class="poster-title">{{ poster.title }}</text>
								<text class="poster-desc">{{ poster.desc }}</text>
							</view>
						</view>
					</view>
				</view>
				<view class="bottom-spacer"></view>
			</scroll-view>
		</view>

		<!-- 悬浮按钮 -->
		<view class="fab-btn">
			<text class="material-symbols-outlined fab-icon">add_photo_alternate</text>
			<text class="fab-text">我的海报</text>
		</view>
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'

	export default {
		components: { TopHeader },
		data() {
			return {
				currentCat: 0,
				categories: ['热门', '全部', '节日问候', '励志鸡汤', '早安晚安', '二手房源', '新房推荐', '豪宅鉴赏'],
				posters: [
					{ title: '豪宅鉴赏·云端之上', desc: '高端大气上档次', tag: 'HOT', tagType: 'accent', aspect: '3/4', image: '' },
					{ title: '冬日暖阳·早安', desc: '每日签到问候', tag: 'NEW', tagType: 'primary', aspect: '9/16', image: '' },
					{ title: '极简生活方式', desc: '4.2w+ 使用', aspect: '1/1', image: '' },
					{ title: '温馨小户型推荐', desc: '刚需首选', tag: 'NEW', tagType: 'primary', aspect: '3/4', image: '' },
					{ title: '奋斗正当时', desc: '励志鸡汤', aspect: '9/16', image: '' },
					{ title: '端午安康·节日签', desc: '节日问候', tag: 'HOT', tagType: 'accent', aspect: '3/4', image: '' },
					{ title: '精装样板间', desc: '实拍展示', aspect: '1/1', image: '' }
				]
			}
		},
		computed: {
			leftPosters() {
				return this.posters.filter((_, i) => i % 2 === 0)
			},
			rightPosters() {
				return this.posters.filter((_, i) => i % 2 !== 0)
			}
		}
	}
</script>

<style lang="scss">
	.poster-container {
		height: 100vh;
		background-color: #f5f7f8;
		display: flex;
		flex-direction: column;
	}

	.header {
		background-color: #ffffff;
		/* 已改为 TopHeader 承担导航栏高度与左右间距，这里只保留容器外观 */
		padding: 0;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
		z-index: 50;

		.search-box {
			width: 50%;
			min-width: 320rpx;
			height: 80rpx;
			background-color: #f5f7f8;
			border-radius: 40rpx;
			display: flex;
			align-items: center;
			padding: 0 24rpx;
			margin-left: auto;

			.search-icon { font-size: 40rpx; color: #94a3b8; }
			.search-input { flex: 1; margin-left: 16rpx; font-size: 28rpx; color: #0f172a; }
		}
	}

	.main-body {
		flex: 1;
		display: flex;
		overflow: hidden;
	}

	.sidebar {
		width: 180rpx;
		background-color: #f5f7f8;
		
		.category-list {
			padding: 16rpx 8rpx;
			display: flex;
			flex-direction: column;
			gap: 8rpx;
		}

		.cat-item {
			height: 112rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			position: relative;
			color: #64748b;
			font-size: 28rpx;
			font-weight: 500;
			border-radius: 16rpx;
			transition: all 0.2s;

			&.active {
				background-color: #ffffff;
				color: #0d7ff2;
				font-weight: bold;
				box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
				
				.active-bar {
					position: absolute;
					left: 0;
					top: 24rpx;
					bottom: 24rpx;
					width: 8rpx;
					background-color: #0d7ff2;
					border-radius: 0 8rpx 8rpx 0;
				}
			}
			
			&:active { background-color: rgba(255, 255, 255, 0.5); }
		}
	}

	.poster-content {
		flex: 1;
		background-color: #ffffff;
		padding: 24rpx;
	}

	.waterfall {
		display: flex;
		gap: 24rpx;
		align-items: flex-start;

		.poster-column {
			flex: 1;
			display: flex;
			flex-direction: column;
			gap: 24rpx;
		}

		.poster-card {
			display: flex;
			flex-direction: column;
			gap: 16rpx;

			.poster-img-box {
				width: 100%;
				border-radius: 16rpx;
				overflow: hidden;
				position: relative;
				background-color: #f1f5f9;

				.poster-image { width: 100%; height: 100%; }
				.poster-image-empty {
					display: flex;
					align-items: center;
					justify-content: center;
					background: linear-gradient(135deg, #eef2f7, #e2e8f0);
					.material-symbols-outlined {
						font-size: 56rpx;
						color: #94a3b8;
					}
				}
				
				.poster-tag {
					position: absolute;
					top: 16rpx;
					left: 16rpx;
					padding: 4rpx 16rpx;
					border-radius: 40rpx;
					font-size: 20rpx;
					font-weight: bold;
					color: #ffffff;
					box-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.1);
					
					&.accent { background-color: #f97316; }
					&.primary { background-color: #0d7ff2; }
				}
			}

			.poster-info {
				display: flex;
				flex-direction: column;
				gap: 4rpx;
				
				.poster-title { font-size: 28rpx; font-weight: bold; color: #0f172a; line-height: 1.2; }
				.poster-desc { font-size: 24rpx; color: #94a3b8; }
			}
			
			&:active .poster-image { transform: scale(1.05); }
		}
	}

	.fab-btn {
		position: fixed;
		bottom: 48rpx;
		right: 32rpx;
		height: 96rpx;
		padding: 0 40rpx;
		background-color: #0d7ff2;
		color: #ffffff;
		border-radius: 48rpx;
		display: flex;
		align-items: center;
		gap: 16rpx;
		box-shadow: 0 12rpx 30rpx rgba(13, 127, 242, 0.3);
		z-index: 100;

		.fab-icon { font-size: 40rpx; }
		.fab-text { font-size: 28rpx; font-weight: bold; }
		
		&:active { transform: scale(0.95); }
	}

	.bottom-spacer { height: 120rpx; }
	.placeholder { color: #94a3b8; }
</style>
