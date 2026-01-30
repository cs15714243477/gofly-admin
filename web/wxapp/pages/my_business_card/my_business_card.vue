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

				<swiper class="card-swiper" circular :indicator-dots="false" :current="currentPreview" @change="onPreviewChange">
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
											<text class="company-name">{{ c.store || '未绑定门店' }}</text>
										</view>
										<view class="pro-tag" :style="{ backgroundColor: c.badgeBg }">
											<text class="material-symbols-outlined pro-icon">verified</text>
											<text>{{ c.badge }}</text>
										</view>
									</view>

									<view class="card-middle">
										<image class="card-avatar" :src="safeImage(c.avatar)" mode="aspectFill"></image>
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
											<view class="qr-code" @click="previewQrCode">
												<image v-if="qrCodeUrl" class="qr-image" :src="qrCodeUrl" mode="aspectFill"></image>
												<text v-else class="material-symbols-outlined">qr_code_2</text>
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
					<view class="form-item" @click="changeAvatar">
						<text class="label">头像</text>
						<view class="val-box">
							<image class="avatar" :src="safeImage(profile.avatar)" mode="aspectFill"></image>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item">
						<text class="label">姓名</text>
						<view class="val-box">
							<input
								class="val-input"
								v-model="profile.name"
								placeholder="请输入姓名"
								placeholder-class="val-placeholder"
							/>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item" @click="showUneditable('手机号由平台绑定，不支持修改')">
						<text class="label">电话</text>
						<view class="val-box">
							<text class="val">{{ profile.phone }}</text>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item">
						<text class="label">微信号</text>
						<view class="val-box">
							<input
								class="val-input"
								v-model="profile.wechat"
								placeholder="选填：用于名片展示"
								placeholder-class="val-placeholder"
							/>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
					<view class="form-item no-border">
						<text class="label">职位</text>
						<view class="val-box">
							<input
								class="val-input"
								v-model="profile.role"
								placeholder="如：资深置业顾问"
								placeholder-class="val-placeholder"
							/>
							<text class="material-symbols-outlined arrow-icon">chevron_right</text>
						</view>
					</view>
				</view>

				<view class="notice-section">
					<view class="form-group">
						<picker
							mode="selector"
							:range="storeOptions"
							range-key="name"
							@change="onStorePickerChange"
						>
							<view class="form-item">
								<text class="label">门店</text>
								<view class="val-box">
									<text class="val">{{ profile.store || '请选择门店' }}</text>
									<text class="material-symbols-outlined arrow-icon">chevron_right</text>
								</view>
							</view>
						</picker>
						<view class="form-item no-border" v-if="profile.store_address">
							<text class="label">门店地址</text>
							<text class="val small">{{ profile.store_address }}</text>
						</view>
					</view>
					<text class="tip-text">* 门店可修改，建议选择与你当前服务区域一致的门店。</text>
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
						<view class="add-tag-btn" @click="openAddTag">
							<text class="material-symbols-outlined add-icon">add</text>
							<text>添加标签</text>
						</view>
					</view>
					<view class="tags-row">
						<view class="tag active" v-for="(t, idx) in profile.tags" :key="'active-' + idx" @click="removeTag(idx)">
							<text>{{ t }}</text>
							<text class="material-symbols-outlined close-icon">close</text>
						</view>
						<view class="tag" v-for="(t, idx) in suggestTags" :key="'suggest-' + idx" @click="addTag(t)">
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
				<button class="share-item share-btn" open-type="share">
					<view class="share-icon green">
						<text class="material-symbols-outlined">chat_bubble</text>
					</view>
					<text class="share-label">发送给朋友</text>
				</button>
				<view class="share-item" @click="saveCardImage">
					<view class="share-icon blue">
						<text class="material-symbols-outlined">image</text>
					</view>
					<text class="share-label">保存图片</text>
				</view>
				<view class="share-item" @click="copyCardLink">
					<view class="share-icon orange">
						<text class="material-symbols-outlined">link</text>
					</view>
					<text class="share-label">复制链接</text>
				</view>
				<view class="share-item" @click="showMoreActions">
					<view class="share-icon grey">
						<text class="material-symbols-outlined">more_horiz</text>
					</view>
					<text class="share-label">更多</text>
				</view>
			</view>
		</view>

		<view class="footer" v-if="currentTab === 1">
			<button class="save-btn" :disabled="saving" @click="saveProfile">
				<text class="material-symbols-outlined">save</text>
				<text>{{ saving ? '保存中...' : '保存修改' }}</text>
			</button>
		</view>

		<!-- 添加标签弹窗（简单实现，避免额外依赖组件） -->
		<view class="mask-popup" v-if="tagPopupVisible" @click="closeAddTag">
			<view class="popup-card" @click.stop>
				<view class="popup-title">添加标签</view>
				<input
					class="popup-input"
					v-model="newTagText"
					placeholder="如：学区房专家"
					placeholder-class="val-placeholder"
					maxlength="12"
				/>
				<view class="popup-actions">
					<button class="popup-btn ghost" @click="closeAddTag">取消</button>
					<button class="popup-btn primary" @click="confirmAddTag">确定</button>
				</view>
			</view>
		</view>

		<!-- 底部导航（与首页一致） -->
		<BottomTabBar active="marketing" />

		<!-- 生成名片图片用：隐藏 canvas -->
		<canvas
			canvas-id="cardCanvas"
			id="cardCanvas"
			class="hidden-canvas"
			:style="{ width: canvasWidth + 'px', height: canvasHeight + 'px' }"
			:width="canvasWidth * canvasScale"
			:height="canvasHeight * canvasScale"
		></canvas>
	</view>
</template>

<script>
	import BottomTabBar from '@/components/BottomTabBar.vue'
	import TopHeader from '@/components/TopHeader.vue'
	import userApi from '@/api/user'
	import $store from '@/store'
	import md5 from 'js-md5'
	import { baseUrl } from '@/utils/config'

	function base64Encode(str = '') {
		try {
			// #ifdef H5
			return window.btoa(str)
			// #endif
		} catch (e) {}
		// #ifndef H5
		const utf8ToBytes = (s) => {
			const bytes = []
			for (let i = 0; i < s.length; i++) {
				let codePoint = s.charCodeAt(i)
				if (codePoint < 0x80) {
					bytes.push(codePoint)
				} else if (codePoint < 0x800) {
					bytes.push(0xc0 | (codePoint >> 6))
					bytes.push(0x80 | (codePoint & 0x3f))
				} else {
					bytes.push(0xe0 | (codePoint >> 12))
					bytes.push(0x80 | ((codePoint >> 6) & 0x3f))
					bytes.push(0x80 | (codePoint & 0x3f))
				}
			}
			return new Uint8Array(bytes)
		}
		const buf = utf8ToBytes(str).buffer
		// eslint-disable-next-line no-undef
		if (typeof wx !== 'undefined' && wx.arrayBufferToBase64) return wx.arrayBufferToBase64(buf)
		if (typeof uni !== 'undefined' && uni.arrayBufferToBase64) return uni.arrayBufferToBase64(buf)
		return ''
		// #endif
	}

	export default {
		components: { BottomTabBar, TopHeader },
		onLoad(options) {
			// 支持从“我的”页跳转过来直接打开“编辑资料”
			const tab = options && (options.tab || options.currentTab)
			if (tab === '1' || tab === 1 || tab === 'edit') {
				this.currentTab = 1
			}
			this.initCanvasSize()
			this.loadProfile().catch(() => {})
		},
		onShow() {
			// 返回页面后，防止缓存导致信息不更新
			this.loadProfile(false).catch(() => {})
		},
		onShareAppMessage() {
			const name = (this.profile && this.profile.name) || '我'
			const agentId = Number(this.loadedUserId || 0)
			const style = Number(this.currentPreview || 0)
			return {
				title: `${name}的名片`,
				path: `/pages/agent_public_card/agent_public_card?agent_id=${encodeURIComponent(agentId)}&style=${encodeURIComponent(style)}`
			}
		},
		data() {
			return {
				// 默认先展示名片预览
				currentTab: 0,
				currentPreview: 0,
				qrCodeUrl: '',
				styleChangeTimer: null,
				saving: false,
				savingImage: false,
				tagPopupVisible: false,
				newTagText: '',
				loadedUserId: 0,
				storeOptions: [],
				canvasWidth: 900,
				canvasHeight: 1320,
				canvasScale: 2,
				// “我的资料”数据源：名片预览与编辑资料保持一致
				profile: {
					avatar: '',
					name: '',
					phone: '',
					wechat: '',
					role: '',
					store_id: 0,
					store: '',
					store_address: '',
					email: '',
					intro: '',
					tags: []
				},
				suggestTags: ['学区房专家', '二手房买卖', '新房咨询', '高端别墅', '婚房首选', '投资置业'],
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
		watch: {
			currentTab(v) {
				if (Number(v) === 1) this.loadStores().catch(() => {})
			}
		},
		methods: {
			initCanvasSize() {
				try {
					const sys = uni.getSystemInfoSync ? uni.getSystemInfoSync() : {}
					const pr = Number(sys.pixelRatio || 2)
					// 固定导出海报尺寸（与绘制坐标系一致：900x1320）
					this.canvasWidth = 900
					this.canvasHeight = 1320

					// 部分机型/环境下 canvas 的像素边长过大，会导致导出只截取左上角。
					// 这里按最大边长做一次兜底，确保“能完整导出”优先。
					const maxSidePx = 2000
					let scale = Math.min(3, pr || 2)
					const byW = Math.floor(maxSidePx / this.canvasWidth)
					const byH = Math.floor(maxSidePx / this.canvasHeight)
					scale = Math.min(scale, Math.max(1, byW), Math.max(1, byH))
					this.canvasScale = Math.max(1, scale)
				} catch (e) {}
			},
			onPreviewChange(e) {
				this.currentPreview = e.detail.current
				this.persistDefaultStyle()
				this.scheduleWxaCodeRefresh()
			},
			safeImage(url) {
				const u = String(url || '').trim()
				return u || '/static/images/img_155372af18.png'
			},
			extrasStorageKey(uid) {
				return `wxapp_my_business_card_extras_${uid || 0}`
			},
			loadExtras(uid) {
				const key = this.extrasStorageKey(uid)
				const raw = uni.getStorageSync(key)
				if (!raw) return {}
				try {
					if (typeof raw === 'string') return JSON.parse(raw)
					return raw || {}
				} catch (e) {
					return {}
				}
			},
			saveExtras(uid) {
				const key = this.extrasStorageKey(uid)
				const payload = {
					wechat: String(this.profile.wechat || '').trim(),
					tags: Array.isArray(this.profile.tags) ? this.profile.tags : []
				}
				try {
					uni.setStorageSync(key, JSON.stringify(payload))
				} catch (e) {
					uni.setStorageSync(key, payload)
				}
			},
			async loadProfile(showLoading = true) {
				const token = uni.getStorageSync('token')
				if (!token && !$store('user').isLogin) return
				const res = await userApi.getCardProfile(!!showLoading)
				if (!res || res.code !== 0) return
				const d = res.data || {}
				const uid = Number(d.id || d.ID || 0)
				this.loadedUserId = uid
				this.profile.avatar = d.avatar || ''
				this.profile.name = d.name || d.nickname || ''
				this.profile.phone = d.mobile || ''
				this.profile.role = d.title || d.role || ''
				this.profile.store_id = Number(d.store_id || 0)
				this.profile.store = d.store_name || ''
				this.profile.store_address = d.store_address || ''
				this.profile.email = d.email || ''
				this.profile.intro = d.introduction || ''

				// 读取上次选择的样式
				const idx = Number(uni.getStorageSync(this.styleStorageKey(uid)) || 0)
				if (!Number.isNaN(idx) && idx >= 0) this.currentPreview = idx

				const extras = this.loadExtras(uid)
				if (extras) {
					if (typeof extras.wechat === 'string') this.profile.wechat = extras.wechat
					if (Array.isArray(extras.tags)) this.profile.tags = extras.tags
				}

				// 刷新小程序码（与当前样式绑定，scene 内携带样式下标）
				this.loadWxaCode(false).catch(() => {})
			},
			styleStorageKey(uid) {
				return `wxapp_my_business_card_style_${uid || 0}`
			},
			persistDefaultStyle() {
				if (!this.loadedUserId) return
				uni.setStorageSync(this.styleStorageKey(this.loadedUserId), Number(this.currentPreview || 0))
			},
			scheduleWxaCodeRefresh() {
				if (this.styleChangeTimer) {
					clearTimeout(this.styleChangeTimer)
					this.styleChangeTimer = null
				}
				this.styleChangeTimer = setTimeout(() => {
					this.loadWxaCode(false).catch(() => {})
				}, 220)
			},
			/**
			 * 加载代理商小程序码
			 *
			 * 根据当前选择的样式生成代理商专属的小程序码URL，支持开发环境跳过路径校验
			 *	
			 * @param {boolean} [showLoading=true] - 是否显示加载提示
			 * @returns {Promise<void>} 无返回值，通过组件属性更新小程序码URL
			 */
			async loadWxaCode(showLoading = true) {
				// 检查用户登录状态
				const token = uni.getStorageSync('token')
				if (!token && !$store('user').isLogin) return
				// 构建请求参数，包含当前选择的样式
				const data = { style: Number(this.currentPreview || 0) }
				// 开发/测试阶段：允许跳过页面路径校验，避免新页面未发布导致生成失败
				// 注意：小程序环境可能没有全局 process，不能直接写 `process && ...`，否则会抛 ReferenceError
				const isDev =
					(typeof process !== 'undefined' && process.env && process.env.NODE_ENV === 'development') ||
					(typeof wx !== 'undefined' &&
						wx.getAccountInfoSync &&
						['develop', 'trial'].includes((wx.getAccountInfoSync().miniProgram || {}).envVersion))
				if (isDev) data.check_path = false
				// 调用API获取小程序码
				const res = await userApi.getAgentWxaCode(data, !!showLoading)
				// 检查响应结果
				if (!res || res.code !== 0) return
				// 更新小程序码URL
				this.qrCodeUrl = (res.data && res.data.url) || ''
			},
			previewQrCode() {
				if (!this.qrCodeUrl) {
					uni.showToast({ title: '小程序码生成中...', icon: 'none' })
					this.loadWxaCode(true).catch(() => {})
					return
				}
				uni.previewImage({ urls: [this.qrCodeUrl] })
			},
			async loadStores() {
				const token = uni.getStorageSync('token')
				if (!token && !$store('user').isLogin) return
				const res = await userApi.getStores(false)
				if (!res || res.code !== 0) return
				const items = Array.isArray(res.data) ? res.data : []
				this.storeOptions = items
			},
			onStorePickerChange(e) {
				const idx = Number(e && e.detail ? e.detail.value : -1)
				if (!Array.isArray(this.storeOptions) || this.storeOptions.length === 0) {
					uni.showToast({ title: '暂无可选门店', icon: 'none' })
					return
				}
				if (idx < 0 || idx >= this.storeOptions.length) return
				const s = this.storeOptions[idx] || {}
				this.profile.store_id = Number(s.id || 0)
				this.profile.store = s.name || ''
				this.profile.store_address = s.address || ''
			},
			showUneditable(msg) {
				uni.showToast({ title: msg || '不支持修改', icon: 'none' })
			},
			openAddTag() {
				this.newTagText = ''
				this.tagPopupVisible = true
			},
			closeAddTag() {
				this.tagPopupVisible = false
				this.newTagText = ''
			},
			confirmAddTag() {
				const t = String(this.newTagText || '').trim()
				if (!t) {
					uni.showToast({ title: '请输入标签内容', icon: 'none' })
					return
				}
				this.addTag(t)
				this.closeAddTag()
			},
			addTag(t) {
				const tag = String(t || '').trim()
				if (!tag) return
				if (!Array.isArray(this.profile.tags)) this.profile.tags = []
				if (this.profile.tags.includes(tag)) return
				if (this.profile.tags.length >= 6) {
					uni.showToast({ title: '最多添加6个标签', icon: 'none' })
					return
				}
				this.profile.tags = this.profile.tags.concat([tag])
			},
			removeTag(idx) {
				if (!Array.isArray(this.profile.tags)) return
				this.profile.tags = this.profile.tags.filter((_, i) => i !== idx)
			},
			async saveProfile() {
				if (this.saving) return
				const name = String(this.profile.name || '').trim()
				if (!name) {
					uni.showToast({ title: '请填写姓名', icon: 'none' })
					return
				}
				this.saving = true
				try {
					const payload = {
						avatar: this.profile.avatar,
						name,
						title: String(this.profile.role || '').trim(),
						introduction: String(this.profile.intro || '').trim(),
						store_id: Number(this.profile.store_id || 0)
					}
					const res = await userApi.updateCardProfile(payload, true)
					if (!res || res.code !== 0) return
					if (this.loadedUserId) this.saveExtras(this.loadedUserId)
					const userStore = $store('user')
					userStore.setUserInfo({
						...(userStore.userInfo || {}),
						avatar: this.profile.avatar,
						name: this.profile.name
					})
					uni.showToast({ title: '保存成功', icon: 'none' })
				} finally {
					this.saving = false
				}
			},
			async changeAvatar() {
				if (this.saving) return
				uni.chooseImage({
					count: 1,
					sizeType: ['compressed'],
					sourceType: ['album', 'camera'],
					success: async (res) => {
						const filePath = res && res.tempFilePaths && res.tempFilePaths[0]
						if (!filePath) return
						await this.uploadAvatar(filePath)
					}
				})
			},
			uploadAvatar(filePath) {
				return new Promise((resolve, reject) => {
					const token = uni.getStorageSync('token')
					const timestamp = Math.floor(Date.now() / 1000)
					const passstr = md5(import.meta.env.GF_API_SECRET + timestamp)
					const header = {
						Accept: 'text/json',
						Businessid: import.meta.env.GF_BUSINESUSSID,
						apiverify: base64Encode(passstr + '#' + timestamp)
					}
					if (token) header.Authorization = `${token}`
					uni.uploadFile({
						url: `${baseUrl}/common/upload/upFile`,
						filePath,
						name: 'file',
						formData: { filetype: 'image' },
						header,
						success: (upRes) => {
							try {
								const raw = upRes && upRes.data ? upRes.data : '{}'
								const out = typeof raw === 'string' ? JSON.parse(raw) : raw
								if (!out || out.code !== 0) {
									uni.showToast({
										title: (out && out.message) || '上传失败',
										icon: 'none'
									})
									reject(out)
									return
								}
								const url = out && out.data ? out.data.url : ''
								this.profile.avatar = url || ''
								uni.showToast({ title: '头像已更新', icon: 'none' })
								resolve(out)
							} catch (e) {
								uni.showToast({ title: '解析上传结果失败', icon: 'none' })
								reject(e)
							}
						},
						fail: (err) => {
							uni.showToast({ title: '上传失败', icon: 'none' })
							reject(err)
						}
					})
				})
			},
			saveCardImage() {
				this.saveCardToAlbum().catch(() => {})
			},
			copyCardLink() {
				this.copyUrlLink().catch(() => {})
			},
			async copyUrlLink() {
				const style = Number(this.currentPreview || 0)
				let envVersion = ''
				// 小程序端：按当前运行环境生成对应的 URLLink（develop|trial|release）
				if (typeof wx !== 'undefined' && wx.getAccountInfoSync) {
					const ev = (wx.getAccountInfoSync().miniProgram || {}).envVersion
					if (ev === 'develop' || ev === 'trial' || ev === 'release') envVersion = ev
				}
				const res = await userApi.getAgentUrlLink({ style, env_version: envVersion }, true)
				const urlLink = res && res.code === 0 && res.data ? res.data.url_link : ''
				if (!urlLink) {
					uni.showToast({ title: '获取链接失败', icon: 'none' })
					return
				}
				uni.setClipboardData({
					data: urlLink,
					success: () => uni.showToast({ title: '已复制链接', icon: 'none' })
				})
			},
			showMoreActions() {
				uni.showActionSheet({
					itemList: ['刷新资料', '清空本地名片扩展信息'],
					success: async (res) => {
						if (res.tapIndex === 0) {
							await this.loadProfile(true).catch(() => {})
						}
						if (res.tapIndex === 1) {
							if (!this.loadedUserId) return
							uni.removeStorageSync(this.extrasStorageKey(this.loadedUserId))
							this.profile.wechat = ''
							this.profile.tags = []
							uni.showToast({ title: '已清空', icon: 'none' })
						}
					}
				})
			},
			getGradientColors(bg) {
				const s = String(bg || '')
				const matches = s.match(/#(?:[0-9a-fA-F]{3}){1,2}/g) || []
				if (matches.length >= 2) return [matches[0], matches[1]]
				return ['#2d9cf0', '#1e40af']
			},
			drawRoundRectPath(ctx, x, y, w, h, r) {
				const rr = Math.max(0, Math.min(r, w / 2, h / 2))
				ctx.beginPath()
				ctx.moveTo(x + rr, y)
				ctx.arcTo(x + w, y, x + w, y + h, rr)
				ctx.arcTo(x + w, y + h, x, y + h, rr)
				ctx.arcTo(x, y + h, x, y, rr)
				ctx.arcTo(x, y, x + w, y, rr)
				ctx.closePath()
			},
			getImageInfo(src) {
				return new Promise((resolve, reject) => {
					if (!src) return reject(new Error('empty src'))
					uni.getImageInfo({
						src,
						success: (res) => resolve(res),
						fail: (err) => reject(err)
					})
				})
			},
			async saveCardToAlbum() {
				if (this.savingImage) return
				this.savingImage = true
				try {
					if (!this.qrCodeUrl) {
						await this.loadWxaCode(true).catch(() => {})
					}
					if (!this.qrCodeUrl) {
						uni.showToast({ title: '小程序码生成失败', icon: 'none' })
						return
					}

					const style = this.previewCards[this.currentPreview] || this.previewCards[0]
					const w = this.canvasWidth * this.canvasScale
					const h = this.canvasHeight * this.canvasScale
					const scale = this.canvasScale
					// 以 900x1320 的设计稿坐标系绘制，再按实际导出尺寸等比缩放，避免出现裁切/偏移
					const baseW = 900
					const s = w / baseW
					const px = (v) => Number(v) * (w / 900)
          const fs = (v) => ath.round(Number(v) * (w / 900))

					const ctx = uni.createCanvasContext('cardCanvas', this)
					ctx.save()
					// ctx.scale(scale, scale)

					// 背景 + 圆角裁剪
					this.drawRoundRectPath(ctx, 0, 0, w, h, px(40))
					ctx.clip()
					const [c1, c2] = this.getGradientColors(style.bg)
					const grd = ctx.createLinearGradient(0, 0, w, h)
					grd.addColorStop(0, c1)
					grd.addColorStop(1, c2)
					ctx.setFillStyle(grd)
					ctx.fillRect(0, 0, w, h)

					// 装饰圆
					ctx.setFillStyle('rgba(255,255,255,0.10)')
					ctx.beginPath()
					ctx.arc(w - px(80), px(80), px(140), 0, Math.PI * 2)
					ctx.fill()
					ctx.setFillStyle(String(style.decor2 || 'rgba(249,115,22,0.18)'))
					ctx.beginPath()
					ctx.arc(px(80), h - px(120), px(180), 0, Math.PI * 2)
					ctx.fill()

					// 顶部门店 + badge
					ctx.setFillStyle('rgba(255,255,255,0.22)')
					this.drawRoundRectPath(ctx, px(40), px(40), px(520), px(70), px(16))
					ctx.fill()
					ctx.setFillStyle('#ffffff')
					ctx.setFontSize(fs(28))
					ctx.fillText(String(style.store || '未绑定门店'), px(60), px(86))

					ctx.setFillStyle(String(style.badgeBg || '#f97316'))
					this.drawRoundRectPath(ctx, w - px(180), px(44), px(140), px(60), px(30))
					ctx.fill()
					ctx.setFillStyle('#ffffff')
					ctx.setFontSize(fs(22))
					ctx.fillText(String(style.badge || 'PRO'), w - px(140), px(82))

					// 头像
					const avatarInfo = await this.getImageInfo(this.safeImage(style.avatar))
					const avatarSize = px(200)
					const avatarX = (w - avatarSize) / 2
					const avatarY = px(170)
					ctx.save()
					ctx.beginPath()
					ctx.arc(w / 2, avatarY + avatarSize / 2, avatarSize / 2, 0, Math.PI * 2)
					ctx.clip()
					ctx.drawImage(avatarInfo.path, avatarX, avatarY, avatarSize, avatarSize)
					ctx.restore()
					ctx.setStrokeStyle('rgba(255,255,255,0.28)')
					ctx.setLineWidth(Math.max(2, px(6)))
					ctx.beginPath()
					ctx.arc(w / 2, avatarY + avatarSize / 2, avatarSize / 2, 0, Math.PI * 2)
					ctx.stroke()

					// 姓名/职位
					ctx.setFillStyle('#ffffff')
					ctx.setFontSize(fs(44))
					const name = String(style.name || '')
					ctx.fillText(name, (w - ctx.measureText(name).width) / 2, px(430))
					ctx.setFontSize(fs(28))
					ctx.setFillStyle('rgba(255,255,255,0.85)')
					const role = String(style.role || '置业顾问')
					ctx.fillText(role, (w - ctx.measureText(role).width) / 2, px(475))

					// 电话块
					ctx.setFillStyle('rgba(255,255,255,0.12)')
					this.drawRoundRectPath(ctx, px(60), px(520), w - px(120), px(120), px(24))
					ctx.fill()
					ctx.setFillStyle('rgba(255,255,255,0.72)')
					ctx.setFontSize(fs(22))
					ctx.fillText('电话', px(90), px(570))
					ctx.setFillStyle('#ffffff')
					ctx.setFontSize(fs(30))
					const phone = String(style.phone || '')
					ctx.fillText(phone, px(90), px(612))

					// 文字信息（微信/门店）
					ctx.setFillStyle('rgba(255,255,255,0.92)')
					ctx.setFontSize(fs(26))
					const wechatLine = '微信号：' + String(style.wechat || '')
					ctx.fillText(wechatLine, px(70), h - px(230))
					const storeLine = '门店：' + String(style.store || '')
					ctx.fillText(storeLine, px(70), h - px(190))

					// 小程序码
					const qrInfo = await this.getImageInfo(this.qrCodeUrl)
					ctx.setFillStyle('#ffffff')
					this.drawRoundRectPath(ctx, w - px(260), h - px(300), px(200), px(200), px(16))
					ctx.fill()
					ctx.drawImage(qrInfo.path, w - px(250), h - px(290), px(180), px(180))
					ctx.setFillStyle('rgba(255,255,255,0.78)')
					ctx.setFontSize(fs(20))
					ctx.fillText('微信识别进入', w - px(250), h - px(80))

					ctx.restore()

					await new Promise((resolve) => ctx.draw(false, resolve))

					const temp = await new Promise((resolve, reject) => {
						// 注意：canvas 实际像素尺寸是 w/h 乘以 canvasScale（模板里绑定了 :width/:height）
						// 导出时用实际像素尺寸，避免只导出左上角
						const pxW = Math.round(w * scale)
						const pxH = Math.round(h * scale)
						uni.canvasToTempFilePath(
							{
								canvasId: 'cardCanvas',
								x: 0,
								y: 0,
								width: pxW,
								height: pxH,
								destWidth: pxW,
								destHeight: pxH,
								fileType: 'png',
								quality: 1,
								success: (res) => resolve(res),
								fail: (err) => reject(err)
							},
							this
						)
					})

					await new Promise((resolve, reject) => {
						uni.saveImageToPhotosAlbum({
							filePath: temp.tempFilePath,
							success: () => resolve(true),
							fail: (err) => reject(err)
						})
					})
					uni.showToast({ title: '已保存到相册', icon: 'none' })
				} catch (e) {
					// 常见：未授权保存到相册
					uni.showModal({
						title: '保存失败',
						content: '请允许保存到相册权限后重试',
						confirmText: '去设置',
						cancelText: '取消',
						success: (res) => {
							if (res.confirm) uni.openSetting({})
						}
					})
				} finally {
					this.savingImage = false
				}
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
				.val-input {
					width: 420rpx;
					font-size: 30rpx;
					color: #0f172a;
					text-align: right;
				}
				.val-placeholder { color: #94a3b8; }
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
		.val.small { font-size: 26rpx; }
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
						overflow: hidden;
						.material-symbols-outlined { font-size: 56rpx; }

						.qr-image {
							width: 84rpx;
							height: 84rpx;
						}
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
			background: transparent;
			border: none;
			padding: 0;

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

		.share-btn::after {
			border: none;
		}
	}

	.mask-popup {
		position: fixed;
		left: 0;
		top: 0;
		right: 0;
		bottom: 0;
		background: rgba(15, 23, 42, 0.45);
		z-index: 999;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 24rpx;

		.popup-card {
			width: 100%;
			max-width: 620rpx;
			background: #ffffff;
			border-radius: 28rpx;
			padding: 28rpx 24rpx;
			box-shadow: 0 30rpx 80rpx rgba(0, 0, 0, 0.18);
		}

		.popup-title {
			font-size: 32rpx;
			font-weight: 800;
			color: #0f172a;
			margin-bottom: 16rpx;
		}

		.popup-input {
			height: 88rpx;
			border-radius: 18rpx;
			background: #f8fafc;
			border: 1rpx solid #e2e8f0;
			padding: 0 20rpx;
			font-size: 30rpx;
			color: #0f172a;
		}

		.popup-actions {
			margin-top: 18rpx;
			display: flex;
			gap: 16rpx;
		}

		.popup-btn {
			flex: 1;
			height: 82rpx;
			border-radius: 18rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 30rpx;
			font-weight: 700;
			border: none;

			&::after { border: none; }

			&.ghost {
				background: #f1f5f9;
				color: #334155;
			}

			&.primary {
				background: #2d9cf0;
				color: #ffffff;
			}
		}
	}

	.hidden-canvas {
		position: fixed;
		left: -9999px;
		top: -9999px;
		opacity: 0;
		pointer-events: none;
	}

	/* 不再需要 bottom spacer（页面不滚动） */
</style>
