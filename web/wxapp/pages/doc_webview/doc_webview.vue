<template>
	<view class="doc-webview-page">
		<!-- #ifdef MP-WEIXIN -->
		<web-view v-if="!loading && !content && isExternalUrl" :src="url"></web-view>
		<!-- #endif -->

		<view v-if="loading || content || !isExternalUrl" class="doc-content-box">
			<view v-if="loading" class="loading-box">
				<text class="loading-text">文档加载中...</text>
			</view>

			<view v-else-if="content" class="rich-box">
				<rich-text :nodes="content"></rich-text>
			</view>

			<view v-else class="fallback-box">
				<text class="fallback-title">{{ title || '文档' }}</text>
				<text class="fallback-url">{{ url || '暂无链接' }}</text>
				<button class="copy-btn" @click="copyUrl">复制链接</button>
			</view>
		</view>
	</view>
</template>

<script>
import userApi from '@/api/user'

export default {
	data() {
		return {
			key: '',
			title: '文档',
			content: '',
			url: '',
			loading: true,
		}
	},
	computed: {
		isExternalUrl() {
			return /^https?:\/\//i.test(String(this.url || '').trim())
		}
	},
	onLoad(options) {
		const rawKey = options && options.key ? decodeURIComponent(options.key) : ''
		this.key = String(rawKey || '').trim()
		this.loadDoc()
	},
	methods: {
		async loadDoc() {
			this.loading = true
			try {
				const res = await userApi.getLoginDocs(false)
				if (!res || res.code !== 0 || !res.data || !Array.isArray(res.data.docs)) {
					return
				}
				const item = res.data.docs.find((doc) => String(doc && doc.key ? doc.key : '').trim() === this.key)
				if (!item) {
					return
				}
				this.title = String(item.title || '').trim() || '文档'
				this.content = String(item.content || '').trim()
				this.url = String(item.url || '').trim()
				uni.setNavigationBarTitle({ title: this.title })
			} finally {
				this.loading = false
			}
		},
		copyUrl() {
			if (!this.url) {
				uni.showToast({ title: '暂无可复制链接', icon: 'none' })
				return
			}
			uni.setClipboardData({
				data: this.url,
				success: () => {
					uni.showToast({ title: '链接已复制', icon: 'none' })
				}
			})
		}
	}
}
</script>

<style lang="scss">
.doc-webview-page {
	height: 100vh;
	background: #ffffff;
}

.doc-content-box {
	height: 100%;
	overflow-y: auto;
}

.loading-box {
	height: 100%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.loading-text {
	color: #94a3b8;
	font-size: 28rpx;
}

.rich-box {
	padding: 32rpx;
	font-size: 28rpx;
	line-height: 1.8;
	color: #334155;
}

.fallback-box {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.fallback-title {
	font-size: 36rpx;
	font-weight: 700;
	color: #0f172a;
}

.fallback-url {
	font-size: 26rpx;
	color: #334155;
	word-break: break-all;
}

.copy-btn {
	margin-top: 16rpx;
	background-color: #2d9cf0;
	color: #ffffff;
	font-size: 30rpx;
	border-radius: 18rpx;
}
</style>
