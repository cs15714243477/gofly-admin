<template>
	<view class="workbench-container">
		<!-- 顶部导航 -->
		<TopHeader title="我的" />

		<!-- 固定一屏展示：不允许上下滚动 -->
		<view class="main-content">
			<!-- 个人信息卡片 -->
			<view class="profile-section">
				<view class="profile-card">
					<view class="profile-bg-decor"></view>
					<button class="edit-btn" @click="goEditCard">
						<text class="material-symbols-outlined edit-icon">edit_square</text>
						<text>编辑</text>
					</button>
					<view class="profile-main">
						<view class="avatar-box">
							<image class="avatar" :src="avatarUrl" mode="aspectFill" ></image>
							<view class="online-status"></view>
						</view>
						<view class="user-info">
							<text class="user-name">{{ displayName }}</text>
							<text class="user-role">{{ displayRoleLine }}</text>
							<text class="user-phone">{{ displayMobile }}</text>
						</view>
					</view>
				</view>
			</view>

			<!-- 业务记录 - 横向3列 -->
			<view class="section">
				<view class="section-title">业务记录</view>
					<view class="records-grid">
					<view class="record-item" v-for="(item, index) in displayRecords" :key="index" @click="openRecord(item)">
						<view class="record-card">
							<view class="record-icon-box">
								<text
									class="material-symbols-outlined record-icon"
									:class="{ 'is-smart-lock': item.key === 'lock_manage' }"
									>{{ item.icon }}</text
								>
								<view class="record-dot" v-if="item.hasNotice"></view>
							</view>
							<text class="record-text">{{ item.name }}</text>
						</view>
					</view>
				</view>
			</view>

			<!-- 更多服务 -->
			<view class="section">
				<view class="section-title">更多服务</view>
				<view class="service-list">
					<view class="service-item">
						<view class="service-left">
							<text class="material-symbols-outlined service-icon">info</text>
							<text class="service-name">关于我们</text>
						</view>
						<text class="material-symbols-outlined arrow-icon">chevron_right</text>
					</view>
					<view class="service-divider"></view>
					<view class="service-item">
						<view class="service-left">
							<text class="material-symbols-outlined service-icon">share</text>
							<text class="service-name">推荐给朋友</text>
						</view>
						<text class="material-symbols-outlined arrow-icon">chevron_right</text>
					</view>
				</view>
			</view>

			<!-- 退出按钮 -->
			<view class="logout-section">
				<button class="logout-btn" @click="handleLogout">退出登录</button>
			</view>
		</view>

		<!-- 底部导航 -->
		<BottomTabBar active="me" />
	</view>
</template>

<script>
	import BottomTabBar from '@/components/BottomTabBar.vue'
	import TopHeader from '@/components/TopHeader.vue'
	import $store from '@/store'

	export default {
		components: { BottomTabBar, TopHeader },
		data() {
			return {
				loadingUser: false,
				debugLogged: false,
				userInfo: {},
				businessRecords: [
					{ key: 'property_manage', name: '房源管理', icon: 'home_work' },
					{ key: 'lock_manage', name: '智能门锁', icon: 'smart_lock' },
					{ key: 'follow', name: '关注记录', icon: 'favorite' },
					{ key: 'unlock', name: '开锁记录', icon: 'lock_open', hasNotice: true },
					{ key: 'showing', name: '带看记录', icon: 'location_on' },
					{ key: 'view', name: '浏览记录', icon: 'history' },
					{ key: 'share', name: '分享记录', icon: 'share' },
					{ key: 'call', name: '通话记录', icon: 'call' }
				]
			}
		},
		computed: {
			avatarUrl() {
				return (this.userInfo && this.userInfo.avatar) ? this.userInfo.avatar : '/static/images/img_83575f387f.png'
			},
			displayName() {


        const u = this.userInfo || {}
				// 真实姓名优先，其次昵称/用户名
				return (u.name || u.nickname || u.username || '未登录')
			},
			displayRoleLine() {
				const u = this.userInfo || {}
				const title = (u.title || '').trim()
				const roleRaw = String(u.role || '').trim()
				// 门店信息由后端返回：store_name（未绑定则为“未绑定”）
				const storeName = (u.store_name || '').trim()
				const store = storeName ? storeName : '未绑定'
				// 角色/身份展示：优先头衔；否则把 role=1/user 映射为“经纪人”
				let left = title
				if (!left) {
					if (roleRaw === '' || roleRaw === '1' || roleRaw.toLowerCase() === 'user') left = '经纪人'
					else left = roleRaw
				}
				return `${left} | ${store}`
			},
			displayMobile() {
				const u = this.userInfo || {}
				return u.mobile || ''
			},
			displayRecords() {
				const u = this.userInfo || {}
				const canManage = Number(u.can_manage_properties) === 1
				const canManageLocks = Number(u.can_manage_locks) === 1
				// 最小权限：仅在允许维护房源时展示入口
				return (this.businessRecords || []).filter((it) => {
					if (!it) return false
					if (it.key === 'property_manage') return canManage
					if (it.key === 'lock_manage') return canManageLocks
					return true
				})
		},
		},
		onShow() {
			this.ensureLoginAndLoadUser()
		},
		methods: {
			debugPrintUserInfo(tag = '') {
				if (this.debugLogged) return
				this.debugLogged = true
				try {
					// 关键字段 + 全量
					console.log('[agent_workbench_home] ' + (tag || 'userInfo') + ' store_name=', this.userInfo?.store_name, 'store_id=', this.userInfo?.store_id, 'title=', this.userInfo?.title, 'role=', this.userInfo?.role)
					console.log('[agent_workbench_home] ' + (tag || 'userInfo') + ' full=', JSON.parse(JSON.stringify(this.userInfo || {})))
				} catch (e) {
					console.log('[agent_workbench_home] debugPrintUserInfo error', e)
				}
			},
			debugShowUserInfo() {
				try {
					const u = this.userInfo || {}
					const contentLines = [
						`name: ${u.name || ''}`,
						`nickname: ${u.nickname || ''}`,
						`username: ${u.username || ''}`,
						`title: ${u.title || ''}`,
						`role: ${u.role || ''}`,
						`store_name: ${u.store_name || ''}`,
						`store_id: ${u.store_id || ''}`,
						`store_address: ${u.store_address || ''}`,
					]
					uni.showModal({
						title: '调试：userInfo(点击头像)',
						content: contentLines.join('\n'),
						showCancel: false,
					})
				} catch (e) {
					uni.showToast({ title: '调试弹窗失败', icon: 'none' })
				}
			},
			async ensureLoginAndLoadUser() {
				const userStore = $store('user')
				// 兼容：小程序刷新后，优先用本地 token 恢复登录态
				const token = uni.getStorageSync('token')
				if (token && !userStore.isLogin) {
					userStore.setToken(token)
				}
				if (!token && !userStore.isLogin) {
					uni.reLaunch({ url: '/pages/login/login' })
					return
				}
				if (this.loadingUser) return
				this.loadingUser = true
				try {
					const info = await userStore.getInfo()
					this.userInfo = info || userStore.userInfo || {}
					this.debugPrintUserInfo('after getInfo')
				} catch (e) {
					// 请求失败：优先保留本地态，不强制跳转（避免短暂网络抖动导致回登录）
					this.userInfo = userStore.userInfo || {}
					this.debugPrintUserInfo('fallback userStore.userInfo')
				} finally {
					this.loadingUser = false
				}
			},
			goEditCard() {
				// 跳转到“获客-编辑资料”页（tab=1）
				uni.reLaunch({
					url: '/pages/my_business_card/my_business_card?tab=1'
				})
			},
			openRecord(item) {
				const map = {
					property_manage: '/pages/property_manage/property_manage',
					lock_manage: '/pages/lock_manage/lock_manage',
					follow: '/pages/records/record_follow',
					unlock: '/pages/records/record_unlock',
					showing: '/pages/records/record_showing',
					view: '/pages/records/record_view',
					share: '/pages/records/record_share',
					call: '/pages/records/record_call'
				}
				const url = item && item.key && map[item.key]
				if (!url) return
				uni.navigateTo({ url })
			},
			handleLogout() {
				uni.showModal({
					title: '提示',
					content: '确定要退出登录吗？',
					success: async (res) => {
						if (res.confirm) {
							try {
								await $store('user').logout(false)
							} catch (e) {}
							uni.reLaunch({
								url: '/pages/login/login'
							})
						}
					}
				})
			}
		}
	}
</script>

<style>
	.workbench-container {
		height: 100vh;
		background-color: #f6f7f8;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	/* header/title 已由 TopHeader 统一 */

	.main-content {
		flex: 1;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		gap: 12rpx;
		padding: 12rpx 24rpx 10rpx;
		box-sizing: border-box;
	}

	.profile-section {
		padding: 0;
	}

	.profile-card {
		background: linear-gradient(135deg, #2d9cf0, #2563eb);
		border-radius: 32rpx;
		padding: 36rpx 32rpx;
		min-height: 240rpx;
		position: relative;
		overflow: hidden;
		box-shadow: 0 12rpx 30rpx rgba(37, 99, 235, 0.2);
		display: flex;
		align-items: center;
	}

	.profile-bg-decor {
		position: absolute;
		right: -90rpx;
		top: -90rpx;
		width: 360rpx;
		height: 360rpx;
		background-color: rgba(255, 255, 255, 0.1);
		border-radius: 50%;
		filter: blur(90rpx);
	}

	.edit-btn {
		position: absolute;
		top: 28rpx;
		right: 28rpx;
		height: 56rpx;
		padding: 0 20rpx;
		background-color: rgba(255, 255, 255, 0.2);
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		color: #ffffff;
		font-size: 24rpx;
		z-index: 10;
	}

	.edit-icon {
		font-size: 24rpx !important;
		margin-right: 4rpx;
	}

	.profile-main {
		display: flex;
		align-items: center;
		gap: 24rpx;
		position: relative;
		z-index: 5;
	}

	.avatar-box {
		position: relative;
	}

	.avatar {
		width: 120rpx;
		height: 120rpx;
		border-radius: 50%;
		border: 4rpx solid rgba(255, 255, 255, 0.3);
	}

	.online-status {
		position: absolute;
		bottom: 4rpx;
		right: 4rpx;
		width: 20rpx;
		height: 20rpx;
		background-color: #4ade80;
		border-radius: 50%;
		border: 4rpx solid #2d9cf0;
	}

	.user-info {
		color: #ffffff;
	}

	.user-name {
		font-size: 38rpx;
		font-weight: bold;
		display: block;
		letter-spacing: 0.4rpx;
	}

	.user-role {
		font-size: 26rpx;
		opacity: 0.9;
		display: block;
		margin-top: 6rpx;
	}

	.user-phone {
		font-size: 26rpx;
		opacity: 0.8;
		display: block;
		margin-top: 4rpx;
	}

	.section {
		padding: 0;
	}

	.section-title {
		font-size: 30rpx;
		font-weight: bold;
		color: #0f172a;
		margin-bottom: 12rpx;
	}

	.records-grid {
		display: flex;
		flex-wrap: wrap;
		margin: 0 -12rpx;
	}

	.record-item {
		width: 33.33%;
		padding: 12rpx;
		box-sizing: border-box;
	}

	.record-card {
		background-color: #ffffff;
		border-radius: 24rpx;
		padding: 18rpx 10rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		border: 1px solid #f1f5f9;
	}

	.record-icon-box {
		width: 72rpx;
		height: 72rpx;
		background-color: #eff6ff;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #2d9cf0;
		margin-bottom: 12rpx;
		position: relative;
	}

	.record-icon {
		font-size: 40rpx !important;
		line-height: 1;
		display: block;
	}

	/* Material Symbols 某些图标字形左右/上下偏移，做一点定向微调 */
	.record-icon.is-smart-lock {
		transform: translateY(2rpx);
	}

	.record-dot {
		position: absolute;
		top: 0;
		right: 0;
		width: 14rpx;
		height: 14rpx;
		background-color: #ef4444;
		border-radius: 50%;
		border: 3rpx solid #ffffff;
	}

	.record-text {
		font-size: 24rpx;
		color: #334155;
		font-weight: 500;
	}

	.service-list {
		background-color: #ffffff;
		border-radius: 24rpx;
		border: 1px solid #f1f5f9;
	}

	.service-item {
		padding: 18rpx 24rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.service-left {
		display: flex;
		align-items: center;
		gap: 20rpx;
	}

	.service-icon {
		font-size: 40rpx !important;
		color: #94a3b8;
	}

	.service-name {
		font-size: 28rpx;
		color: #334155;
	}

	.arrow-icon {
		font-size: 32rpx !important;
		color: #cbd5e1;
	}

	.service-divider {
		height: 1px;
		background-color: #f1f5f9;
		margin: 0 32rpx;
	}

	.logout-section {
		padding: 0;
	}

	.logout-btn {
		width: 100%;
		height: 80rpx;
		background-color: #fff1f2;
		color: #ef4444;
		font-size: 28rpx;
		font-weight: bold;
		border-radius: 20rpx;
		border: 1px solid #fee2e2;
	}

	.tab-bar {
		background-color: #ffffff;
		border-top: 1rpx solid #e2e8f0;
		padding: 12rpx 0 calc(env(safe-area-inset-bottom) + 12rpx);
		display: flex;
		justify-content: space-around;
		align-items: center;
	}

	.tab-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		color: #94a3b8;
		flex: 1;
	}

	.tab-icon {
		font-size: 48rpx !important;
	}

	.tab-text {
		font-size: 20rpx;
		margin-top: 4rpx;
	}

	.tab-item.active {
		color: #2d9cf0;
	}

	.active-bg {
		background-color: #eef7ff;
		padding: 8rpx 32rpx;
		border-radius: 40rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.fill {
		font-variation-settings: 'FILL' 1;
	}

	/* 已取消 bottom spacer：页面不滚动 */
</style>
