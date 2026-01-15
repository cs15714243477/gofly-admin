<template>
	<view class="list-container">
		<!-- 顶部导航与搜索 -->
		<view class="header">
			<TopHeader title="房源" />
			
			<view class="search-row">
				<view class="search-box">
					<text class="material-symbols-outlined search-icon">search</text>
					<input class="search-input" type="text" placeholder="搜索小区、商圈或价格..." placeholder-class="placeholder" />
				</view>
				<view class="filter-icon-btn" :class="{ active: hasActiveFilters }" @click="openFilter">
					<text class="material-symbols-outlined">tune</text>
					<view v-if="hasActiveFilters" class="filter-dot"></view>
				</view>
			</view>
			
			<scroll-view scroll-x="true" class="category-scroll">
				<view class="category-row">
					<view
						class="category-item"
						v-for="(c, idx) in categories"
						:key="c.key || idx"
						:class="{ active: currentCategory === c.key }"
						@click="setCategory(c.key)"
					>
						{{ c.label }}
					</view>
				</view>
			</scroll-view>
		</view>

		<scroll-view scroll-y="true" class="main-list">
			<view class="list-content">
				<view class="property-card" v-for="(item, index) in displayList" :key="index" :class="{ 'disabled': item.status === '已下架' }" @click="goToDetail(item)">
					<view class="card-main">
						<view class="image-box">
							<image class="property-image" :src="item.image" mode="aspectFill"></image>
							<view class="image-tag" v-if="item.tag">
								<text class="material-symbols-outlined tag-icon">{{ item.tagIcon }}</text>
								<text>{{ item.tag }}</text>
							</view>
							<view class="mask" v-if="item.status === '交易中' || item.status === '已下架'">
								<text class="mask-text">{{ item.status }}</text>
							</view>
						</view>
						<view class="info-box">
							<view class="title">{{ item.title }}</view>
							<view class="meta">
								<text class="bold">{{ item.rooms }}</text>
								<text class="divider">|</text>
								<text>{{ item.size }}m²</text>
								<text class="divider">|</text>
								<text>{{ item.orientation }}</text>
							</view>
							<view class="tags">
								<text class="tag" v-for="(tag, tIdx) in item.tags" :key="tIdx" :class="tag.type">{{ tag.name }}</text>
							</view>
							<view class="price-row">
								<text class="price">¥{{ item.price }}</text>
								<text class="unit-price">{{ item.unitPrice }}</text>
							</view>
							<view class="stats">
								<text>数据浏览: {{ item.views }}</text>
								<text class="stats-sep"></text>
								<text>我浏览过/关注: {{ item.myViews }}</text>
							</view>
						</view>
					</view>
					<view class="card-footer" :class="item.footerClass || 'orange-footer'">
						<view class="footer-left">
							<view class="footer-icon-box">
								<text class="material-symbols-outlined footer-icon">{{ item.footerIcon || 'currency_yen' }}</text>
							</view>
							<text class="footer-text">{{ item.commission }}</text>
						</view>
						<button class="footer-btn" v-if="item.status !== '已下架'" @click.stop="handlePromote(item)">
							<text class="material-symbols-outlined btn-icon">share</text>
							<text>微信推广</text>
						</button>
						<button class="footer-btn-disabled" v-else>已下架</button>
					</view>
				</view>
			</view>
			<view class="bottom-spacer"></view>
		</scroll-view>

		<!-- 底部导航 -->
		<BottomTabBar active="property" />

		<!-- 筛选弹层（静态数据 + 本地过滤/排序） -->
		<view v-if="filterOpen" class="filter-mask" @click="closeFilter">
			<view class="filter-sheet" @click.stop>
				<view class="filter-header">
					<text class="filter-title">筛选</text>
					<text class="filter-reset" @click="resetFilter">重置</text>
				</view>

				<view class="filter-section">
					<text class="filter-section-title">户型</text>
					<view class="chip-row">
						<view
							v-for="(opt, idx) in roomOptions"
							:key="'room-' + idx"
							class="chip"
							:class="{ active: draft.rooms === opt.value }"
							@click="draft.rooms = opt.value"
						>
							<text>{{ opt.label }}</text>
						</view>
					</view>
				</view>

				<view class="filter-section">
					<text class="filter-section-title">特色</text>
					<view class="chip-row">
						<view class="chip" :class="{ active: !!draft.smartLock }" @click="draft.smartLock = !draft.smartLock">
							<text>智能门锁</text>
						</view>
					</view>
				</view>

				<view class="filter-section">
					<text class="filter-section-title">状态</text>
					<view class="chip-row">
						<view
							v-for="(opt, idx) in statusOptions"
							:key="'status-' + idx"
							class="chip"
							:class="{ active: draft.status === opt.value }"
							@click="draft.status = opt.value"
						>
							<text>{{ opt.label }}</text>
						</view>
					</view>
				</view>

				<view class="filter-section">
					<text class="filter-section-title">排序</text>
					<view class="radio-list">
						<view
							v-for="(opt, idx) in sortOptions"
							:key="'sort-' + idx"
							class="radio-item"
							:class="{ active: draft.sort === opt.value }"
							@click="draft.sort = opt.value"
						>
							<text>{{ opt.label }}</text>
							<text v-if="draft.sort === opt.value" class="material-symbols-outlined radio-check">check</text>
						</view>
					</view>
				</view>

				<view class="filter-footer">
					<button class="filter-btn ghost" @click="closeFilter">取消</button>
					<button class="filter-btn primary" @click="applyFilter">确定</button>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
	import BottomTabBar from '@/components/BottomTabBar.vue'
	import TopHeader from '@/components/TopHeader.vue'

	export default {
		components: { BottomTabBar, TopHeader },
		data() {
			return {
				currentCategory: 'all',
				categories: [
					{ key: 'all', label: '全部房源' },
					{ key: 'mine', label: '我的房源' },
					{ key: 'nearby', label: '附近房源' },
					{ key: 'price_drop', label: '最新降价' },
					{ key: 'school', label: '学区房' }
				],
				// 筛选（静态数据）
				filters: {
					rooms: '',
					smartLock: false,
					status: '',
					sort: ''
				},
				filterOpen: false,
				draft: {
					rooms: '',
					smartLock: false,
					status: '',
					sort: ''
				},
				roomOptions: [
					{ value: '', label: '不限' },
					{ value: '1', label: '1室' },
					{ value: '2', label: '2室' },
					{ value: '3', label: '3室' },
					{ value: '4+', label: '4室+' }
				],
				statusOptions: [
					{ value: '', label: '不限' },
					{ value: 'normal', label: '正常' },
					{ value: '交易中', label: '交易中' },
					{ value: '已下架', label: '已下架' }
				],
				sortOptions: [
					{ value: '', label: '默认' },
					{ value: 'views_desc', label: '浏览最多' },
					{ value: 'myviews_desc', label: '我最关注' }
				],
				propertyList: [
					{
						title: '13# 12.12开盘 核心卖点 海景高层',
						rooms: '3室2厅',
						size: '105',
						orientation: '南向',
						price: '850万',
						unitPrice: '80,952元/m²',
						views: 127,
						myViews: 1,
						tag: '智能门锁',
						tagIcon: 'lock',
						image: '/static/images/img_c63bd4d9ba.png',
						tags: [
							{ name: '自营', type: 'blue' },
							{ name: '精装', type: 'green' },
							{ name: '学区房', type: 'orange' }
						],
						commission: '佣金1.5%，成交奖励经纪人8000'
					},
					{
						title: '金城中心 舒适单身公寓 地铁口直达',
						rooms: '1室1厅',
						size: '45',
						orientation: '北向',
						price: '210万',
						unitPrice: '46,666元/m²',
						views: 89,
						myViews: 0,
						image: '/static/images/img_f0b09ae08c.png',
						tags: [
							{ name: '地铁房', type: 'indigo' },
							{ name: '有电梯', type: 'grey' }
						],
						commission: '佣金1%，极速结佣',
						footerIcon: 'percent'
					},
					{
						title: '滨江别墅 独栋花园 视野开阔 无遮挡',
						rooms: '4室3厅',
						size: '210',
						orientation: '东向',
						price: '1280万',
						unitPrice: '60,952元/m²',
						views: 256,
						myViews: 45,
						tag: 'VR看房',
						tagIcon: 'videocam',
						image: '/static/images/img_f1b597a564.png',
						tags: [
							{ name: '自营', type: 'blue' },
							{ name: 'VR看房', type: 'purple' }
						],
						commission: '高佣2%，房东急售',
						footerIcon: 'local_fire_department'
					},
					{
						title: '日落大道 2室紧凑型 适合过渡',
						rooms: '2室1厅',
						size: '78',
						orientation: '西向',
						price: '420万',
						unitPrice: '53,846元/m²',
						views: 45,
						myViews: 2,
						status: '交易中',
						image: '/static/images/img_fdf5f4bc67.png',
						tags: [
							{ name: '随时看房', type: 'grey' }
						],
						commission: '佣金1%',
						footerIcon: 'history',
						footerClass: 'grey-footer'
					},
					{
						title: '日落大道 2室紧凑型 适合过渡',
						rooms: '2室1厅',
						size: '78',
						orientation: '西向',
						price: '420万',
						unitPrice: '53,846元/m²',
						views: 45,
						myViews: 2,
						status: '已下架',
						image: '/static/images/img_fdf5f4bc67.png',
						tags: [
							{ name: '随时看房', type: 'grey' }
						],
						commission: '佣金1%',
						footerIcon: 'history',
						footerClass: 'grey-footer'
					}
				]
			}
		},
		computed: {
			hasActiveFilters() {
				return !!(this.filters.rooms || this.filters.smartLock || this.filters.status || this.filters.sort)
			},
			displayList() {
				let list = Array.isArray(this.propertyList) ? [...this.propertyList] : []

				// 户型：按 rooms 前缀数字筛选
				if (this.filters.rooms) {
					list = list.filter((it) => {
						const s = (it && it.rooms) || ''
						const m = String(s).match(/^(\d+)/)
						const n = m ? Number(m[1]) : 0
						if (this.filters.rooms === '4+') return n >= 4
						return String(n) === this.filters.rooms
					})
				}

				// 智能门锁：tag 或 icon
				if (this.filters.smartLock) {
					list = list.filter((it) => it && (it.tag === '智能门锁' || it.tagIcon === 'lock'))
				}

				// 状态
				if (this.filters.status) {
					if (this.filters.status === 'normal') {
						list = list.filter((it) => !it.status)
					} else {
						list = list.filter((it) => it && it.status === this.filters.status)
					}
				}

				// 排序
				if (this.filters.sort === 'views_desc') {
					list.sort((a, b) => (b.views || 0) - (a.views || 0))
				} else if (this.filters.sort === 'myviews_desc') {
					list.sort((a, b) => (b.myViews || 0) - (a.myViews || 0))
				}

				return list
			}
		},
		methods: {
			setCategory(key) {
				this.currentCategory = key
			},
			handlePromote(item) {
				const title = (item && item.title) || '房源'
				uni.showActionSheet({
					itemList: ['生成获客海报', '复制推广文案', '模拟分享（占位）'],
					success: (res) => {
						if (res.tapIndex === 0) {
							uni.navigateTo({
								url: `/pages/marketing_posters/marketing_posters?title=${encodeURIComponent(title)}`
							})
							return
						}
						if (res.tapIndex === 1) {
							const text = `【优质房源推荐】${title}\n点击查看详情（示例）：/pages/property_detail/property_detail`
							uni.setClipboardData({
								data: text,
								success: () => {
									uni.showToast({ title: '已复制推广文案', icon: 'none' })
								}
							})
							return
						}
						if (res.tapIndex === 2) {
							uni.showToast({ title: '分享功能待接入', icon: 'none' })
						}
					}
				})
			},
			openFilter() {
				this.filterOpen = true
				this.draft = {
					rooms: this.filters.rooms,
					smartLock: !!this.filters.smartLock,
					status: this.filters.status,
					sort: this.filters.sort
				}
			},
			closeFilter() {
				this.filterOpen = false
			},
			resetFilter() {
				this.draft = { rooms: '', smartLock: false, status: '', sort: '' }
			},
			applyFilter() {
				this.filters = {
					rooms: this.draft.rooms,
					smartLock: !!this.draft.smartLock,
					status: this.draft.status,
					sort: this.draft.sort
				}
				this.closeFilter()
			},
			goToDetail(item) {
				if (item.status === '已下架') return
				uni.navigateTo({
					url: '/pages/property_detail/property_detail'
				})
			}
		}
	}
</script>

<style lang="scss">
	.list-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
	}

	.header {
		background-color: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);
		/* TopHeader 自带左右 padding，这里避免叠加 */
		padding: 0 0 16rpx;
		border-bottom: 1rpx solid #e2e8f0;
		z-index: 50;

		.nav-row {
			display: flex;
			align-items: center;
			justify-content: space-between;
			height: 80rpx;
			margin-bottom: 16rpx;
			position: relative;

			.nav-title {
				font-size: 32rpx;
				font-weight: bold;
				color: #0f172a;
				flex: 1;
				text-align: left;
				padding-right: 80rpx;
			}

			.filter-icon-btn {
				position: absolute;
				right: 0;
				width: 80rpx;
				height: 80rpx;
				display: flex;
				align-items: center;
				justify-content: center;
				border-radius: 50%;
				position: relative;

				.material-symbols-outlined {
					font-size: 40rpx;
					color: #0f172a;
				}
				
				&:active {
					background-color: #f1f5f9;
				}

				&.active .material-symbols-outlined {
					color: #2d9cf0;
				}
			}

			.filter-dot {
				position: absolute;
				top: 14rpx;
				right: 14rpx;
				width: 14rpx;
				height: 14rpx;
				background-color: #f97316;
				border-radius: 50%;
				border: 3rpx solid rgba(255, 255, 255, 0.95);
			}
		}

		.search-row {
			margin-top: 16rpx;
			margin-bottom: 24rpx;
			padding: 0 16rpx;
			box-sizing: border-box;
			display: flex;
			align-items: center;
			gap: 16rpx;

			.search-box {
				flex: 1;
				height: 88rpx;
				background-color: #f8fafc;
				border-radius: 20rpx;
				display: flex;
				align-items: center;
				padding: 0 24rpx;
				border: 1px solid rgba(226, 232, 240, 0.95);
				box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);

				.search-icon {
					font-size: 40rpx;
					color: #94a3b8;
				}

				.search-input {
					flex: 1;
					margin-left: 16rpx;
					font-size: 28rpx;
					color: #0f172a;
				}
			}
		}

		.category-scroll {
			white-space: nowrap;
			width: 100%;
			padding: 0 16rpx;
			box-sizing: border-box;

			.category-row {
				display: inline-flex;
				gap: 16rpx;
				padding-bottom: 4rpx;
			}

			.category-item {
				padding: 12rpx 32rpx;
				background-color: #ffffff;
				border-radius: 40rpx;
				font-size: 24rpx;
				font-weight: 500;
				color: #64748b;
				border: 1px solid #e2e8f0;
				white-space: nowrap;

				&.active {
					background-color: #2d9cf0;
					color: #ffffff;
					border-color: #2d9cf0;
					box-shadow: 0 4rpx 8rpx rgba(45, 156, 240, 0.2);
				}
			}
		}
	}

	/* TopHeader 右侧按钮：统一尺寸/交互 */
	.filter-icon-btn {
		width: 80rpx;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		position: relative;
		background-color: transparent;

		.material-symbols-outlined {
			font-size: 40rpx;
			color: #0f172a;
		}

		&:active {
			background-color: #f1f5f9;
		}

		&.active .material-symbols-outlined {
			color: #2d9cf0;
		}
	}

	.filter-dot {
		position: absolute;
		top: 14rpx;
		right: 14rpx;
		width: 14rpx;
		height: 14rpx;
		background-color: #f97316;
		border-radius: 50%;
		border: 3rpx solid rgba(255, 255, 255, 0.95);
	}

	.main-list {
		flex: 1;
		overflow: hidden;
	}

	.list-content {
		padding: 24rpx 16rpx;
		display: flex;
		flex-direction: column;
		gap: 24rpx;
	}

	.property-card {
		background-color: #ffffff;
		border-radius: 24rpx;
		overflow: hidden;
		box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.02);
		border: 1px solid #f1f5f9;

		&.disabled {
			opacity: 0.8;
			
			.title {
				color: #64748b;
			}
			
			.price {
				color: #94a3b8;
			}
		}

		.card-main {
			padding: 20rpx;
			display: flex;
			gap: 18rpx;

			.image-box {
				width: 224rpx;
				height: 224rpx;
				border-radius: 16rpx;
				overflow: hidden;
				position: relative;
				flex-shrink: 0;
				background-color: #f1f5f9;

				.property-image {
					width: 100%;
					height: 100%;
				}
				
				.image-tag {
					position: absolute;
					top: 8rpx;
					left: 8rpx;
					background-color: rgba(0, 0, 0, 0.6);
					backdrop-filter: blur(4rpx);
					color: #ffffff;
					font-size: 20rpx;
					padding: 4rpx 12rpx;
					border-radius: 8rpx;
					display: flex;
					align-items: center;
					gap: 4rpx;

					.tag-icon {
						font-size: 20rpx;
					}
				}

				.mask {
					position: absolute;
					inset: 0;
					background-color: rgba(15, 23, 42, 0.4);
					display: flex;
					align-items: center;
					justify-content: center;

					.mask-text {
						background-color: rgba(30, 41, 59, 0.8);
						color: #ffffff;
						font-size: 24rpx;
						font-weight: bold;
						padding: 4rpx 16rpx;
						border-radius: 8rpx;
					}
				}
				
				.grayscale {
					filter: grayscale(1);
				}
			}

			.info-box {
				flex: 1;
				min-width: 0;
				display: flex;
				flex-direction: column;
				gap: 6rpx;

				.title {
					font-size: 30rpx;
					font-weight: bold;
					color: #0f172a;
					line-height: 1.3;
					overflow: hidden;
					text-overflow: ellipsis;
					display: -webkit-box;
					line-clamp: 2;
					-webkit-line-clamp: 2;
					-webkit-box-orient: vertical;
				}

				.meta {
					display: flex;
					align-items: center;
					gap: 12rpx;
					font-size: 24rpx;
					color: #64748b;

					.bold {
						font-weight: 500;
						color: #334155;
					}

					.divider {
						color: #e2e8f0;
					}
				}

				.tags {
					display: flex;
					flex-wrap: wrap;
					gap: 12rpx;
					margin: 2rpx 0;

					.tag {
						font-size: 20rpx;
						font-weight: 500;
						padding: 4rpx 12rpx;
						border-radius: 8rpx;
						border: 1px solid transparent;

						&.blue { background-color: #eff6ff; color: #2563eb; border-color: #dbeafe; }
						&.green { background-color: #f0fdf4; color: #16a34a; border-color: #dcfce7; }
						&.orange { background-color: #fff7ed; color: #ea580c; border-color: #ffedd5; }
						&.indigo { background-color: #eef2ff; color: #4f46e5; border-color: #e0e7ff; }
						&.purple { background-color: #faf5ff; color: #9333ea; border-color: #f3e8ff; }
						&.grey { background-color: #f8fafc; color: #64748b; }
					}
				}

				.price-row {
					display: flex;
					align-items: baseline;
					gap: 16rpx;
					margin-top: auto;

					.price {
						font-size: 36rpx;
						font-weight: bold;
						color: #f97316;
					}

					.unit-price {
						font-size: 20rpx;
						color: #94a3b8;
					}
				}

				.stats {
					font-size: 20rpx;
					color: #94a3b8;
					display: flex;
					align-items: center;

					.stats-sep {
						width: 1px;
						height: 16rpx;
						background-color: #e2e8f0;
						margin: 0 12rpx;
					}
				}
			}
		}

		.card-footer {
			display: flex;
			align-items: center;
			justify-content: space-between;
			padding: 14rpx 20rpx;
			border-top: 1rpx solid transparent;

			&.orange-footer {
				background-color: #fff7ed;
				border-color: #ffedd5;
				
				.footer-icon-box { background-color: #f97316; }
				.footer-text { color: #9a3412; font-weight: bold; }
			}

			&.grey-footer {
				background-color: #f8fafc;
				border-color: #f1f5f9;
				
				.footer-icon-box { background-color: #cbd5e1; }
				.footer-text { color: #64748b; font-weight: 500; }
			}

			.footer-left {
				display: flex;
				align-items: center;
				gap: 16rpx;
				flex: 1;
				min-width: 0;

				.footer-icon-box {
					width: 40rpx;
					height: 40rpx;
					border-radius: 50%;
					display: flex;
					align-items: center;
					justify-content: center;
					flex-shrink: 0;

					.footer-icon {
						color: #ffffff;
						font-size: 28rpx;
					}
				}

				.footer-text {
					font-size: 24rpx;
					overflow: hidden;
					text-overflow: ellipsis;
					white-space: nowrap;
				}
			}

			.footer-btn {
				height: 52rpx;
				padding: 0 24rpx;
				background-color: #ffffff;
				border: 1px solid #fed7aa;
				border-radius: 40rpx;
				display: flex;
				align-items: center;
				gap: 8rpx;
				font-size: 22rpx;
				font-weight: bold;
				color: #ea580c;
				line-height: 1;
				margin-left: 24rpx;
				
				&::after { border: none; }

				.btn-icon {
					font-size: 28rpx;
				}

				&:active {
					background-color: #fff7ed;
					transform: scale(0.95);
				}
			}
			
			.footer-btn-disabled {
				height: 52rpx;
				padding: 0 24rpx;
				background-color: #ffffff;
				border: 1px solid #e2e8f0;
				border-radius: 40rpx;
				font-size: 22rpx;
				color: #94a3b8;
				margin-left: 24rpx;
				&::after { border: none; }
			}
		}
	}

	.tab-bar {
		background-color: #ffffff;
		border-top: 1rpx solid #e2e8f0;
		padding: 16rpx 48rpx calc(env(safe-area-inset-bottom) + 16rpx);
		display: flex;
		justify-content: space-between;
		align-items: center;
		position: relative;
		z-index: 100;

		.tab-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 8rpx;
			color: #94a3b8;

			.tab-icon {
				font-size: 52rpx;
				
				&.fill {
					font-variation-settings: 'FILL' 1;
				}
			}

			.tab-text {
				font-size: 20rpx;
				font-weight: 500;
			}

			&.active {
				color: #2d9cf0;
			}
		}

		.add-btn-container {
			position: relative;
			top: -40rpx;
			
			.add-btn {
				width: 112rpx;
				height: 112rpx;
				background-color: #2d9cf0;
				border-radius: 50%;
				display: flex;
				align-items: center;
				justify-content: center;
				box-shadow: 0 8rpx 20rpx rgba(45, 156, 240, 0.4);

				.material-symbols-outlined {
					color: #ffffff;
					font-size: 48rpx;
				}

				&:active {
					transform: scale(0.9) translateY(4rpx);
				}
			}
		}
	}

	.bottom-spacer {
		height: 160rpx;
	}

	.placeholder {
		color: #94a3b8;
	}

	/* 筛选弹层（避免与卡片内 .mask 冲突，单独命名） */
	.filter-mask {
		position: fixed;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		background-color: rgba(15, 23, 42, 0.45);
		z-index: 999;
		display: flex;
		align-items: flex-end;
	}

	.filter-sheet {
		width: 100%;
		background-color: #ffffff;
		border-radius: 32rpx 32rpx 0 0;
		padding: 24rpx 24rpx calc(env(safe-area-inset-bottom) + 24rpx);
		box-sizing: border-box;
	}

	.filter-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 8rpx;
	}

	.filter-title {
		font-size: 30rpx;
		font-weight: 800;
		color: #0f172a;
	}

	.filter-reset {
		font-size: 26rpx;
		color: #64748b;
	}

	.filter-section {
		padding-top: 18rpx;
	}

	.filter-section-title {
		font-size: 26rpx;
		font-weight: 700;
		color: #334155;
		margin-bottom: 12rpx;
		display: block;
	}

	.chip-row {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;
	}

	.chip {
		padding: 12rpx 22rpx;
		border-radius: 40rpx;
		background-color: #f1f5f9;
		color: #334155;
		font-size: 24rpx;
		border: 1px solid transparent;

		&.active {
			background-color: rgba(45, 156, 240, 0.12);
			color: #2d9cf0;
			border-color: rgba(45, 156, 240, 0.25);
			font-weight: 700;
		}
	}

	.radio-list {
		border-radius: 20rpx;
		overflow: hidden;
		border: 1px solid #f1f5f9;
	}

	.radio-item {
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 0 18rpx;
		border-bottom: 1rpx solid #f1f5f9;
		color: #334155;

		&:last-child {
			border-bottom: none;
		}

		&.active {
			color: #2d9cf0;
			font-weight: 700;
		}
	}

	.radio-check {
		font-size: 36rpx;
	}

	.filter-footer {
		display: flex;
		gap: 16rpx;
		margin-top: 20rpx;
	}

	.filter-btn {
		flex: 1;
		height: 80rpx;
		border-radius: 20rpx;
		font-size: 28rpx;
		font-weight: 800;
		display: flex;
		align-items: center;
		justify-content: center;
		line-height: 1;
		border: none;

		&::after { border: none; }

		&.ghost {
			background-color: #f1f5f9;
			color: #334155;
		}

		&.primary {
			background-color: #2d9cf0;
			color: #ffffff;
		}
	}
</style>
