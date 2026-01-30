<template>
	<view class="page">
		<TopHeader title="经纪人名片" />

		<view class="content">
			<view class="card" :style="{ background: currentStyle.bg }">
				<view class="decor-1"></view>
				<view class="decor-2" :style="{ backgroundColor: currentStyle.decor2 }"></view>

				<view class="inner">
					<view class="top-row">
						<view class="store">
							<text class="material-symbols-outlined store-icon">apartment</text>
							<text class="store-name">{{ agent.store_name || '未绑定门店' }}</text>
						</view>
						<view class="badge" :style="{ backgroundColor: currentStyle.badgeBg }">
							<text class="material-symbols-outlined badge-icon">verified</text>
							<text>{{ currentStyle.badge }}</text>
						</view>
					</view>

					<view class="mid">
						<image class="avatar" :src="safeImage(agent.avatar)" mode="aspectFill"></image>
						<view class="name">{{ agent.name || '-' }}</view>
						<view class="title">{{ agent.title || '置业顾问' }}</view>
					</view>

					<view class="bottom">
						<view class="line" v-if="agent.store_address">
							<text class="label">门店地址：</text>
							<text class="val">{{ agent.store_address }}</text>
						</view>
						<view class="line" v-if="agent.introduction">
							<text class="label">简介：</text>
							<text class="val">{{ agent.introduction }}</text>
						</view>

						<view class="actions">
							<button class="btn primary" @click="callAgent" :disabled="!agent.mobile">
								<text class="material-symbols-outlined">call</text>
								<text>电话联系</text>
							</button>
							<button class="btn ghost" @click="copyMobile" :disabled="!agent.mobile">
								<text class="material-symbols-outlined">content_copy</text>
								<text>复制电话</text>
							</button>
						</view>
					</view>
				</view>
			</view>

			<view v-if="loading" class="hint">加载中...</view>
			<view v-else-if="!loading && !agent.id" class="hint">暂无经纪人信息</view>
		</view>
	</view>
</template>

<script>
	import TopHeader from '@/components/TopHeader.vue'
	import userApi from '@/api/user'

	function parseScene(scene = '') {
		const s = String(scene || '').trim()
		// 兼容 a1_s0 / a=1&s=0
		if (s.includes('=') && s.includes('&')) {
			const out = {}
			s.split('&').forEach((pair) => {
				const [k, v] = pair.split('=')
				if (k) out[k] = v
			})
			return {
				agentId: Number(out.a || out.agent_id || 0),
				style: Number(out.s || 0)
			}
		}
		const m = s.match(/a(\d+)_s(\d+)/)
		if (m) return { agentId: Number(m[1] || 0), style: Number(m[2] || 0) }
		const m2 = s.match(/a(\d+)/)
		if (m2) return { agentId: Number(m2[1] || 0), style: 0 }
		return { agentId: 0, style: 0 }
	}

	export default {
		components: { TopHeader },
		data() {
			return {
				loading: true,
				agent: {},
				styleIndex: 0,
				styles: [
					{ bg: 'linear-gradient(135deg, #2d9cf0, #1e40af)', decor2: 'rgba(249, 115, 22, 0.2)', badge: 'PRO', badgeBg: '#f97316' },
					{ bg: 'linear-gradient(135deg, #06b6d4, #2563eb)', decor2: 'rgba(255, 255, 255, 0.12)', badge: 'TOP', badgeBg: '#22c55e' },
					{ bg: 'linear-gradient(135deg, #111827, #334155)', decor2: 'rgba(45, 156, 240, 0.18)', badge: 'VIP', badgeBg: '#2d9cf0' }
				]
			}
		},
		computed: {
			currentStyle() {
				return this.styles[this.styleIndex] || this.styles[0]
			}
		},
		onLoad(options) {
			const scene = options && options.scene ? decodeURIComponent(options.scene) : ''
			const agentId = Number(options.agent_id || 0)
			const parsed = parseScene(scene)
			this.styleIndex = Math.max(0, Number(options.style || parsed.style || 0))
			const id = agentId || parsed.agentId
			this.fetchAgent(id)
		},
		methods: {
			safeImage(url) {
				const u = String(url || '').trim()
				return u || '/static/images/img_155372af18.png'
			},
			async fetchAgent(agentId) {
				this.loading = true
				try {
					if (!agentId) {
						this.agent = {}
						return
					}
					const res = await userApi.getAgentCard({ agent_id: agentId }, true)
					if (!res || res.code !== 0) {
						this.agent = {}
						return
					}
					this.agent = res.data || {}
				} finally {
					this.loading = false
				}
			},
			callAgent() {
				const mobile = String((this.agent && this.agent.mobile) || '').trim()
				if (!mobile) return
				uni.makePhoneCall({ phoneNumber: mobile })
			},
			copyMobile() {
				const mobile = String((this.agent && this.agent.mobile) || '').trim()
				if (!mobile) return
				uni.setClipboardData({
					data: mobile,
					success: () => uni.showToast({ title: '已复制', icon: 'none' })
				})
			}
		}
	}
</script>

<style lang="scss">
	.page {
		min-height: 100vh;
		background: #f6f7f8;
	}

	.content {
		padding: 24rpx;
	}

	.card {
		position: relative;
		border-radius: 40rpx;
		overflow: hidden;
		box-shadow: 0 40rpx 80rpx rgba(0, 0, 0, 0.18);
		border: 1rpx solid rgba(255, 255, 255, 0.08);
	}

	.decor-1 {
		position: absolute;
		top: -80rpx;
		right: -80rpx;
		width: 256rpx;
		height: 256rpx;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.08);
	}

	.decor-2 {
		position: absolute;
		bottom: -90rpx;
		left: -90rpx;
		width: 320rpx;
		height: 320rpx;
		border-radius: 50%;
		background: rgba(249, 115, 22, 0.2);
	}

	.inner {
		position: relative;
		z-index: 2;
		padding: 40rpx;
		color: #fff;
	}

	.top-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.store {
		display: flex;
		align-items: center;
		gap: 12rpx;
		.store-icon { font-size: 34rpx; opacity: 0.9; }
		.store-name { font-size: 28rpx; font-weight: 800; opacity: 0.95; }
	}

	.badge {
		padding: 6rpx 16rpx;
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		gap: 8rpx;
		font-size: 22rpx;
		font-weight: 800;
		.badge-icon { font-size: 26rpx; }
	}

	.mid {
		padding: 50rpx 0 30rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16rpx;
		.avatar {
			width: 180rpx;
			height: 180rpx;
			border-radius: 50%;
			border: 6rpx solid rgba(255, 255, 255, 0.25);
			background: rgba(255, 255, 255, 0.06);
		}
		.name { font-size: 44rpx; font-weight: 900; }
		.title { font-size: 28rpx; font-weight: 600; opacity: 0.86; }
	}

	.bottom {
		background: rgba(255, 255, 255, 0.10);
		border: 1rpx solid rgba(255, 255, 255, 0.08);
		border-radius: 28rpx;
		padding: 22rpx 22rpx;
		.line {
			display: flex;
			gap: 8rpx;
			line-height: 1.5;
			margin-bottom: 12rpx;
		}
		.label { font-size: 24rpx; opacity: 0.78; }
		.val { font-size: 26rpx; font-weight: 600; opacity: 0.92; }
	}

	.actions {
		margin-top: 8rpx;
		display: flex;
		gap: 16rpx;
		.btn {
			flex: 1;
			height: 84rpx;
			border-radius: 22rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			gap: 10rpx;
			font-size: 30rpx;
			font-weight: 800;
			border: none;
			&::after { border: none; }
		}
		.primary {
			background: rgba(255, 255, 255, 0.95);
			color: #0f172a;
		}
		.ghost {
			background: rgba(255, 255, 255, 0.14);
			color: #ffffff;
		}
	}

	.hint {
		padding: 28rpx 0;
		text-align: center;
		color: #64748b;
		font-size: 26rpx;
	}
</style>

