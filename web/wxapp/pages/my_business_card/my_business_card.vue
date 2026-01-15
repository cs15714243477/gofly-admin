<template>
	<view class="card-container">
		<!-- 顶部导航 -->
		<TopHeader title="获客" />

		<!-- 选项卡切换 -->
		<view class="tabs-box">
			<view class="tabs-inner">
				<view class="tab-item" :class="{ active: currentTab === 0 }" @click="currentTab = 0">名片预览</view>
				<view class="tab-item" :class="{ active: currentTab === 1 }" @click="currentTab = 1">编辑资料</view>
			</view>
		</view>

		<view class="main-content">
			<!-- 名片预览内容 -->
			<view v-if="currentTab === 0" class="preview-content">
				<view class="pagination-dots">
					<view
						v-for="(_, i) in previewCards"
						:key="i"
						class="dot"
						:class="{ active: currentPreview === i }"
					></view>
				</view>

				<swiper class="card-swiper" circular :indicator-dots="false" @change="onPreviewChange">
					<swiper-item v-for="(c, idx) in previewCards" :key="idx">
						<view class="card-slide">
							<view class="business-card" :style="{ background: c.bg }">
								<view class="card-decor-1"></view>
								<view class="card-decor-2" :style="{ backgroundColor: c.decor2 }"></view>

								<view class="card-inner">
									<view class="card-top">
										<view class="company-info">
											<view class="company-logo">
												<text class="material-symbols-outlined">apartment</text>
											</view>
											<text class="company-name">{{ c.company }}</text>
										</view>
										<view class="pro-tag" :style="{ backgroundColor: c.badgeBg }">
											<text class="material-symbols-outlined pro-icon">verified</text>
											<text>{{ c.badge }}</text>
										</view>
									</view>

									<view class="card-middle">
										<image class="card-avatar" :src="c.avatar" mode="aspectFill"></image>
										<view class="card-user">
											<text class="card-name">{{ c.name }}</text>
											<text class="card-role">{{ c.role }}</text>
										</view>
										<view class="card-tags">
											<text class="card-tag" v-for="(t, tIdx) in c.tags" :key="tIdx">{{ t }}</text>
										</view>
									</view>

									<view class="card-bottom">
										<view class="contact-info">
											<view class="contact-item">
												<view class="contact-icon">
													<text class="material-symbols-outlined">call</text>
												</view>
												<view class="contact-text">
													<text class="contact-label">Phone</text>
													<text class="contact-val">{{ c.phone }}</text>
												</view>
											</view>
										</view>
										<view class="extra-info">
											<view class="extra-text">
												<text class="extra-line">
													<text class="extra-label">微信号：</text>
													<text class="extra-val">{{ c.wechat }}</text>
												</text>
												<text class="extra-line">
													<text class="extra-label">门店：</text>
													<text class="extra-val">{{ c.store }}</text>
												</text>
											</view>
											<view class="qr-code">
												<text class="material-symbols-outlined">qr_code_2</text>
											</view>
										</view>
									</view>
								</view>
							</view>
						</view>
					</swiper-item>
				</swiper>

				<text class="swipe-tip">左右滑动切换名片样式</text>
			</view>

			<!-- 编辑资料内容 -->
			<scroll-view v-if="currentTab === 1" scroll-y="true" class="edit-scroll">
				<view class="edit-content">
				<view class="form-group">
					<view class="form-item">
						<text class="label">头像</text>
						<view class="val-box">
							<image class="avatar" :src="profile.avatar" mode="aspectFill"></image>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item">
						<text class="label">姓名</text>
						<view class="val-box">
							<text class="val">{{ profile.name }}</text>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item">
						<text class="label">电话</text>
						<view class="val-box">
							<text class="val">{{ profile.phone }}</text>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item">
						<text class="label">微信号</text>
						<view class="val-box">
							<text class="val">{{ profile.wechat }}</text>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item no-border">
						<text class="label">职位</text>
						<view class="val-box">
							<text class="val">{{ profile.role }}</text>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
				</view>

				<view class="notice-section">
					<view class="form-group">
						<view class="form-item">
							<text class="label">公司</text>
							<text class="val">{{ profile.company }}</text>
						</view>
						<view class="form-item no-border">
							<text class="label">门店</text>
							<text class="val">{{ profile.store }}</text>
						</view>
					</view>
					<text class="tip-text">* 公司及门店信息由平台认证，需要修改请联系平台管理员。</text>
				</view>

				<view class="form-group bio-group">
					<text class="label">自我介绍</text>
					<textarea
						v-model="profile.intro"
						class="bio-textarea"
						placeholder="介绍一下你的从业经历和服务理念..."
					></textarea>
				</view>

				<view class="form-group tags-group">
					<view class="tags-header">
						<text class="label">擅长领域</text>
						<view class="add-tag-btn">
							<text class="material-symbols-outlined add-icon">add</text>
							<text>添加标签</text>
						</view>
					</view>
					<view class="tags-row">
						<view class="tag active" v-for="(t, idx) in profile.tags" :key="'active-' + idx">
							<text>{{ t }}</text>
							<text class="material-symbols-outlined close-icon">close</text>
						</view>
						<view class="tag" v-for="(t, idx) in suggestTags" :key="'suggest-' + idx">
							<text>{{ t }}</text>
							<text class="material-symbols-outlined add-icon">add</text>
						</view>
					</view>
				</view>
				</view>
			</scroll-view>
		</view>

		<!-- 底部操作（占位在页面内，避免出现页面滚动） -->
		<view class="footer share-footer" v-if="currentTab === 0">
			<view class="share-options">
				<view class="share-item">
					<view class="share-icon green">
						<text class="material-symbols-outlined">chat_bubble</text>
					</view>
					<text class="share-label">发送给朋友</text>
				</view>
				<view class="share-item">
					<view class="share-icon blue">
						<text class="material-symbols-outlined">image</text>
					</view>
					<text class="share-label">保存图片</text>
				</view>
				<view class="share-item">
					<view class="share-icon orange">
						<text class="material-symbols-outlined">link</text>
					</view>
					<text class="share-label">复制链接</text>
				</view>
				<view class="share-item">
					<view class="share-icon grey">
						<text class="material-symbols-outlined">more_horiz</text>
					</view>
					<text class="share-label">更多</text>
				</view>
			</view>
		</view>

		<view class="footer" v-if="currentTab === 1">
			<button class="save-btn">
				<text class="material-symbols-outlined">save</text>
				<text>保存修改</text>
			</button>
		</view>

		<!-- 底部导航（与首页一致） -->
		<BottomTabBar active="marketing" />
	</view>
</template>

<script>
	import BottomTabBar from '@/components/BottomTabBar.vue'
	import TopHeader from '@/components/TopHeader.vue'

	export default {
		components: { BottomTabBar, TopHeader },
		onLoad(options) {
			// 支持从“我的”页跳转过来直接打开“编辑资料”
			const tab = options && (options.tab || options.currentTab)
			if (tab === '1' || tab === 1 || tab === 'edit') {
				this.currentTab = 1
			}
		},
		data() {
			return {
				// 默认先展示名片预览
				currentTab: 0,
				currentPreview: 0,
				// “我的资料”数据源：名片预览与编辑资料保持一致
				profile: {
					avatar:
						'/static/images/img_155372af18.png',
					name: '张伟',
					phone: '138 1234 5678',
					wechat: 'zhangwei_agent',
					role: '资深置业顾问',
					company: '链家地产',
					store: '朝阳区北京新天地店',
					email: 'zhangwei@lianjia.com',
					addr: '北京市朝阳区...',
					intro:
						'从事房地产行业5年，专注于朝阳区高端住宅买卖。我致力于为每一位客户提供专业、诚信的置业咨询服务，帮助您找到理想的家。',
					tags: ['二手房买卖', '学区房专家']
				},
				suggestTags: ['高端别墅'],
				// 名片样式（只改变视觉，不改变字段）
				cardStyles: [
					{ bg: 'linear-gradient(135deg, #2d9cf0, #1e40af)', decor2: 'rgba(249, 115, 22, 0.2)', badge: 'PRO', badgeBg: '#f97316' },
					{ bg: 'linear-gradient(135deg, #06b6d4, #2563eb)', decor2: 'rgba(255, 255, 255, 0.12)', badge: 'TOP', badgeBg: '#22c55e' },
					{ bg: 'linear-gradient(135deg, #111827, #334155)', decor2: 'rgba(45, 156, 240, 0.18)', badge: 'VIP', badgeBg: '#2d9cf0' }
				]
			}
		},
		computed: {
			previewCards() {
				const p = this.profile || {}
				const tags = Array.isArray(p.tags) ? p.tags.slice(0, 3) : []
				return (this.cardStyles || []).map((s) => ({
					...s,
					company: p.company,
					name: p.name,
					role: p.role,
					phone: p.phone,
					wechat: p.wechat,
					store: p.store,
					avatar: p.avatar,
					tags
				}))
			}
		},
		methods: {
			onPreviewChange(e) {
				this.currentPreview = e.detail.current
			}
		}
	}
</script>

<style lang="scss">
	.card-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.header {
		/* header 已由 TopHeader 统一 */
		display: none;
	}

	.tabs-box {
		padding: 24rpx 32rpx;
		background-color: #ffffff;
		
		.tabs-inner {
			display: flex;
			background-color: #f1f5f9;
			padding: 8rpx;
			border-radius: 20rpx;
		}

		.tab-item {
			flex: 1;
			height: 72rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 28rpx;
			font-weight: 500;
			color: #64748b;
			border-radius: 16rpx;
			transition: all 0.3s;

			&.active {
				background-color: #ffffff;
				color: #0f172a;
				box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.05);
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

	.edit-scroll {
		flex: 1;
		min-height: 0;
	}

	.edit-content {
		padding: 20rpx 24rpx 16rpx;
		display: flex;
		flex-direction: column;
		gap: 16rpx;
	}

	.form-group {
		background-color: #ffffff;
		border-radius: 24rpx;
		overflow: hidden;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);
		
		.form-item {
			padding: 24rpx;
			display: flex;
			justify-content: space-between;
			align-items: center;
			border-bottom: 1rpx solid #f1f5f9;

			&.no-border { border-bottom: none; }
			
			.label { font-size: 30rpx; font-weight: 500; color: #0f172a; }
			
			.val-box {
				display: flex;
				align-items: center;
				gap: 16rpx;
				color: #64748b;

				.avatar { width: 96rpx; height: 96rpx; border-radius: 50%; border: 1rpx solid #e2e8f0; }
				.val { font-size: 30rpx; }
				.arrow-icon { font-size: 40rpx; }
			}
			
			&:active { background-color: #f8fafc; }
		}
	}

	.notice-section {
		display: flex;
		flex-direction: column;
		gap: 12rpx;

		.val { color: #64748b; font-size: 30rpx; }
		.tip-text { padding: 0 16rpx; font-size: 24rpx; color: #94a3b8; line-height: 1.6; }
	}

	.bio-group {
		padding: 24rpx;
		display: flex;
		flex-direction: column;
		gap: 16rpx;

		.bio-textarea {
			width: 100%;
			height: 160rpx;
			background-color: #f8fafc;
			border: 1px solid #e2e8f0;
			border-radius: 16rpx;
			padding: 24rpx;
			font-size: 30rpx;
			color: #334155;
			line-height: 1.5;
		}
	}

	.tags-group {
		padding: 24rpx;
		display: flex;
		flex-direction: column;
		gap: 20rpx;

		.tags-header {
			display: flex;
			justify-content: space-between;
			align-items: center;

			.add-tag-btn {
				display: flex;
				align-items: center;
				gap: 4rpx;
				font-size: 24rpx;
				color: #2d9cf0;
				font-weight: 600;

				.add-icon { font-size: 32rpx; }
			}
		}

		.tags-row {
			display: flex;
			flex-wrap: wrap;
			gap: 20rpx;

			.tag {
				padding: 12rpx 24rpx;
				background-color: #f1f5f9;
				border-radius: 40rpx;
				font-size: 24rpx;
				color: #64748b;
				display: flex;
				align-items: center;
				gap: 8rpx;
				border: 1px solid transparent;

				&.active {
					background-color: rgba(45, 156, 240, 0.1);
					color: #2d9cf0;
					border-color: rgba(45, 156, 240, 0.2);
				}

				.close-icon, .add-icon { font-size: 28rpx; }
			}
		}
	}

	.preview-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 24rpx 24rpx 16rpx;
		flex: 1;
		min-height: 0;

		.pagination-dots {
			display: flex;
			gap: 16rpx;
			margin-bottom: 16rpx;

			.dot {
				width: 12rpx;
				height: 12rpx;
				border-radius: 50%;
				background-color: #cbd5e1;

				&.active { background-color: #2d9cf0; }
			}
		}

		.card-swiper {
			width: 100%;
			flex: 1;
			min-height: 0;
			/* swiper 本身不参与布局，居中由 swiper-item 内部容器完成 */
		}

		.card-slide {
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
		}

		.business-card {
			/* 小程序端可视区较小：整体缩小一点，避免被底部操作/导航遮挡 */
			width: 600rpx;
			height: 880rpx;
			border-radius: 40rpx;
			position: relative;
			overflow: hidden;
			box-shadow: 0 40rpx 80rpx rgba(0, 0, 0, 0.2);
			border: 1rpx solid rgba(255, 255, 255, 0.08);

			.card-decor-1 {
				position: absolute;
				top: -80rpx;
				right: -80rpx;
				width: 256rpx;
				height: 256rpx;
				background-color: rgba(255, 255, 255, 0.05);
				border-radius: 50%;
				filter: blur(40rpx);
			}

			.card-decor-2 {
				position: absolute;
				bottom: -80rpx;
				left: -80rpx;
				width: 320rpx;
				height: 320rpx;
				background-color: rgba(249, 115, 22, 0.2);
				border-radius: 50%;
				filter: blur(80rpx);
			}

			.card-inner {
				position: relative;
				z-index: 10;
				height: 100%;
				display: flex;
				flex-direction: column;
				padding: 40rpx;
				color: #ffffff;
			}

			.card-top {
				display: flex;
				justify-content: space-between;
				align-items: center;

				.company-info {
					display: flex;
					align-items: center;
					gap: 16rpx;

					.company-logo {
						width: 64rpx;
						height: 64rpx;
						background-color: rgba(255, 255, 255, 0.2);
						backdrop-filter: blur(10px);
						border-radius: 12rpx;
						display: flex;
						align-items: center;
						justify-content: center;
						.material-symbols-outlined { font-size: 32rpx; }
					}
					
					.company-name { font-size: 28rpx; font-weight: bold; opacity: 0.9; }
				}

				.pro-tag {
					background-color: #f97316;
					padding: 4rpx 16rpx;
					border-radius: 40rpx;
					display: flex;
					align-items: center;
					gap: 8rpx;
					font-size: 20rpx;
					font-weight: bold;
					
					.pro-icon { font-size: 24rpx; }
				}
			}

			.card-middle {
				flex: 1;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				gap: 24rpx;

				.card-avatar {
					width: 168rpx;
					height: 168rpx;
					border-radius: 50%;
					border: 6rpx solid rgba(255, 255, 255, 0.3);
					box-shadow: 0 16rpx 32rpx rgba(0, 0, 0, 0.2);
				}

				.card-user {
					text-align: center;
					.card-name { font-size: 44rpx; font-weight: bold; display: block; }
					.card-role { font-size: 28rpx; font-weight: 500; opacity: 0.8; margin-top: 8rpx; display: block; }
				}

				.card-tags {
					display: flex;
					gap: 12rpx;
					
					.card-tag {
						font-size: 20rpx;
						padding: 4rpx 16rpx;
						background-color: rgba(255, 255, 255, 0.1);
						border: 1px solid rgba(255, 255, 255, 0.1);
						backdrop-filter: blur(10px);
						border-radius: 8rpx;
					}
				}
			}

			.card-bottom {
				display: flex;
				flex-direction: column;
				gap: 36rpx;

				.contact-info {
					background-color: rgba(255, 255, 255, 0.1);
					backdrop-filter: blur(10px);
					border: 1px solid rgba(255, 255, 255, 0.05);
					border-radius: 24rpx;
					padding: 24rpx;
					
					.contact-item {
						display: flex;
						align-items: center;
						gap: 24rpx;
						
						.contact-icon {
							width: 64rpx;
							height: 64rpx;
							background-color: #ffffff;
							color: #2d9cf0;
							border-radius: 50%;
							display: flex;
							align-items: center;
							justify-content: center;
							.material-symbols-outlined { font-size: 32rpx; }
						}
						
						.contact-text {
							display: flex;
							flex-direction: column;
							.contact-label { font-size: 20rpx; text-transform: uppercase; opacity: 0.6; }
							.contact-val { font-size: 28rpx; font-weight: bold; }
						}
					}
				}

				.extra-info {
					display: flex;
					justify-content: space-between;
					align-items: flex-end;

					.extra-text {
						display: flex;
						flex-direction: column;
						gap: 4rpx;
						font-size: 26rpx;
						opacity: 0.92;
						letter-spacing: 0.5rpx;
					}

					.extra-line {
						display: block;
						line-height: 1.5;
					}

					.extra-label {
						opacity: 0.8;
					}

					.extra-val {
						font-weight: 700;
					}

					.qr-code {
						width: 84rpx;
						height: 84rpx;
						background-color: #ffffff;
						border-radius: 12rpx;
						display: flex;
						align-items: center;
						justify-content: center;
						color: #0f172a;
						.material-symbols-outlined { font-size: 56rpx; }
					}
				}
			}
		}

		.swipe-tip { font-size: 24rpx; color: #94a3b8; margin-top: 12rpx; }
	}

	.footer {
		background-color: rgba(255, 255, 255, 0.9);
		backdrop-filter: blur(10px);
		/* safe-area 交给 BottomTabBar 处理，避免底部过高 */
		padding: 20rpx 24rpx 20rpx;
		border-top: 1rpx solid #e2e8f0;
		z-index: 100;

		.save-btn {
			width: 100%;
			height: 88rpx;
			background-color: #2d9cf0;
			color: #ffffff;
			border-radius: 24rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 16rpx;
			font-size: 32rpx;
			font-weight: bold;
			box-shadow: 0 10rpx 25rpx rgba(45, 156, 240, 0.3);
			border: none;
			
			&::after { border: none; }
			&:active { background-color: #1d82cc; transform: scale(0.98); }
		}
	}

	.share-footer {
		/* 与 BottomTabBar 的“胶囊卡片”风格统一：本层透明，只保留外边距 */
		background: transparent;
		border-top: none;
		backdrop-filter: none;
		padding: 8rpx 16rpx 0;
		box-shadow: none;

		.share-options {
			display: flex;
			justify-content: space-between;
			align-items: center;
			gap: 0;
			/* 复用 BottomTabBar 的磨砂/边框/圆角/阴影 */
			background: rgba(255, 255, 255, 0.92);
			backdrop-filter: blur(12px);
			-webkit-backdrop-filter: blur(12px);
			border: 1rpx solid rgba(226, 232, 240, 0.9);
			border-radius: 28rpx;
			box-shadow: 0 -10rpx 24rpx rgba(15, 23, 42, 0.06);
			padding: 18rpx 18rpx;
		}

		.share-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 10rpx;
			flex: 1;
			min-width: 0;

			.share-icon {
				/* 改为扁平 pill，和底部导航 icon-wrap 对齐 */
				width: 88rpx;
				height: 56rpx;
				border-radius: 28rpx;
				display: flex;
				align-items: center;
				justify-content: center;
				color: #0f172a;
				transition: transform 0.18s ease, background-color 0.18s ease, border-color 0.18s ease;
				border: 1rpx solid transparent;
				
				&.green { background-color: rgba(34, 197, 94, 0.10); border-color: rgba(34, 197, 94, 0.18); color: #16a34a; }
				&.blue { background-color: rgba(37, 99, 235, 0.10); border-color: rgba(37, 99, 235, 0.18); color: #2563eb; }
				&.orange { background-color: rgba(249, 115, 22, 0.10); border-color: rgba(249, 115, 22, 0.18); color: #f97316; }
				&.grey { background-color: rgba(148, 163, 184, 0.14); border-color: rgba(148, 163, 184, 0.18); color: #64748b; }

				.material-symbols-outlined { font-size: 40rpx; }
			}

			.share-label {
				font-size: 22rpx;
				font-weight: 700;
				color: #475569;
				white-space: nowrap;
				overflow: hidden;
				text-overflow: ellipsis;
			}
			
			&:active .share-icon { transform: scale(1.06); }
		}
	}

	/* 不再需要 bottom spacer（页面不滚动） */
</style>
