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

		<scroll-view class="main-content" scroll-y="true">
			<!-- 提示卡片 -->
			<view class="tip-section">
				<view class="tip-card">
					<text class="material-symbols-outlined tip-icon">near_me</text>
					<text class="tip-text">请确保您已站在门锁附近 (1-2米范围内)，并保持手机蓝牙已开启。</text>
				</view>
			</view>

			<!-- 蓝牙状态 -->
			<view class="ble-section">
				<view class="ble-card">
					<view class="ble-head">
						<view class="ble-left">
							<text class="material-symbols-outlined ble-ico">bluetooth</text>
							<text class="ble-title">蓝牙状态</text>
						</view>
						<view class="ble-badge" :class="bleBadgeClass">
							<text>{{ bleBadgeText }}</text>
						</view>
					</view>

					<view class="ble-grid">
						<view class="ble-item">
							<text class="ble-lab">初始化</text>
							<text class="ble-val" :class="bleInitClass">{{ bleInitText }}</text>
						</view>
						<view class="ble-item">
							<text class="ble-lab">可用</text>
							<text class="ble-val">{{ ble.available ? '是' : '否' }}</text>
						</view>
						<view class="ble-item">
							<text class="ble-lab">搜索中</text>
							<text class="ble-val">{{ ble.discovering ? '是' : '否' }}</text>
						</view>
						<view class="ble-item">
							<text class="ble-lab">连接</text>
							<text class="ble-val">{{ bleConnectText }}</text>
						</view>
					</view>

					<view v-if="ble.deviceId" class="ble-extra">
						<text class="ble-extra-lab">设备</text>
						<text class="ble-extra-val">{{ ble.deviceId }}</text>
					</view>
					<view v-if="ble.errorMsg" class="ble-error">
						<text class="material-symbols-outlined err-ico">error</text>
						<text class="err-text">{{ ble.errorMsg }}</text>
					</view>

					<view class="ble-log">
						<view class="ble-log-head">
							<text class="ble-log-title">实时日志</text>
							<text class="ble-log-tip">{{ bleLogs.length }}/{{ bleLogMax }}</text>
						</view>
						<scroll-view scroll-y class="ble-log-body" :scroll-into-view="bleLogAnchor" :scroll-with-animation="true">
							<view v-for="(it, idx) in bleLogs" :key="idx" class="ble-log-line">
								<text class="ble-log-time">{{ it.time }}</text>
								<text class="ble-log-msg" :class="it.level">{{ it.msg }}</text>
							</view>
							<view :id="bleLogAnchor" style="height: 1rpx;"></view>
						</scroll-view>
					</view>
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
				<view class="status-tag" :class="statusTagClass">
					<view class="dot-box">
						<view class="dot-ping"></view>
						<view class="dot-inner"></view>
					</view>
					<text class="status-text">{{ statusText }}</text>
				</view>
			</view>
			
		</scroll-view>

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
	import ttlockApi from '@/api/ttlock'

	function pad2(n) {
		return String(n).padStart(2, '0')
	}

	export default {
		components: { TopHeader },
		data() {
			return {
				propertyId: 0,
				unlocking: false,
				unlockRequestId: 0,
				ble: {
					supported: true,
					initStatus: 'idle', // idle/initting/ready/failed/unsupported
					available: false,
					discovering: false,
					connected: false,
					deviceId: '',
					errorMsg: '',
				},
				bleLogs: [],
				bleLogMax: 50,
				bleLogAnchor: 'ble-log-bottom',
				_bleListening: false,
				_lastBleStateKey: '',
				_lastConnKey: '',
			}
		},
		computed: {
			bleInitText() {
				if (!this.ble.supported) return '不支持'
				switch (this.ble.initStatus) {
					case 'initting': return '初始化中'
					case 'ready': return '已就绪'
					case 'failed': return '失败'
					case 'unsupported': return '不支持'
					default: return '未初始化'
				}
			},
			bleInitClass() {
				if (this.ble.initStatus === 'ready') return 'ok'
				if (this.ble.initStatus === 'failed') return 'warn'
				if (this.ble.initStatus === 'initting') return 'loading'
				return 'muted'
			},
			bleBadgeText() {
				if (!this.ble.supported || this.ble.initStatus === 'unsupported') return '当前平台不支持'
				if (this.ble.initStatus === 'failed') return '蓝牙不可用'
				if (this.ble.available) return '可用'
				return '未开启/不可用'
			},
			bleBadgeClass() {
				if (!this.ble.supported || this.ble.initStatus === 'unsupported') return 'muted'
				if (this.ble.initStatus === 'failed') return 'warn'
				if (this.ble.available) return 'ok'
				return 'warn'
			},
			bleConnectText() {
				if (this.ble.connected) return '已连接'
				if (this.unlocking) return '连接中'
				return '未连接'
			},
			statusText() {
				if (!this.ble.supported || this.ble.initStatus === 'unsupported') return '当前平台不支持蓝牙能力'
				if (this.ble.initStatus === 'initting') return '正在初始化蓝牙…'
				if (this.ble.initStatus === 'failed') return this.ble.errorMsg ? `蓝牙初始化失败：${this.ble.errorMsg}` : '蓝牙初始化失败'
				if (!this.ble.available) return '请开启手机蓝牙后重试'
				if (this.unlocking) return '正在连接并开锁…'
				if (this.ble.connected) return '设备已连接'
				return '蓝牙就绪，等待连接'
			},
			statusTagClass() {
				if (!this.ble.supported || this.ble.initStatus === 'unsupported') return 'muted'
				if (this.ble.initStatus === 'failed' || !this.ble.available) return 'warn'
				if (this.unlocking || this.ble.initStatus === 'initting') return 'loading'
				return 'ok'
			},
		},
		onLoad(options) {
			const pid = Number((options && (options.property_id || options.id)) || 0)
			this.propertyId = Number.isFinite(pid) ? pid : 0
		},
		onShow() {
			this.appendBleLog('info', '进入页面，开始检测蓝牙状态')
			this.initBluetooth()
		},
		onHide() {
			this.cleanupBluetooth()
		},
		onUnload() {
			this.cleanupBluetooth()
		},
		methods: {
			goBack() {
				uni.navigateBack()
			},
			appendBleLog(level, msg) {
				const d = new Date()
				const time = `${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(d.getSeconds())}`
				this.bleLogs.push({ time, msg: String(msg || ''), level: level || 'info' })
				if (this.bleLogs.length > this.bleLogMax) {
					this.bleLogs.splice(0, this.bleLogs.length - this.bleLogMax)
				}
				this.bleLogAnchor = 'ble-log-bottom'
			},
			formatBleErr(err) {
				if (!err) return ''
				const code = Number(err.errCode || err.code || 0)
				const raw = String(err.errMsg || err.message || '').trim()
				if (code === 10001) return '蓝牙适配器不可用（请开启手机蓝牙）'
				if (code === 10002) return '没有找到指定设备'
				if (code === 10003) return '连接失败'
				if (code === 10004) return '没有找到指定服务'
				if (code === 10005) return '没有找到指定特征值'
				if (code === 10006) return '当前连接已断开'
				if (code === 10007) return '当前特征值不支持该操作'
				if (code === 10008) return '其余所有系统上报的异常'
				if (code === 10009) return '系统版本低于 4.3 不支持 BLE'
				if (code === 10010) return '系统版本低于 4.4 不支持 BLE'
				if (code === 10011) return '需要授权蓝牙权限'
				if (raw) return raw
				if (code) return `蓝牙异常（错误码：${code}）`
				try { return JSON.stringify(err) } catch (e) {}
				return '蓝牙异常'
			},
			async reportUnlockActivity(stage, meta) {
				if (!this.propertyId) return
				if (!ttlockApi || typeof ttlockApi.logUnlockActivity !== 'function') return
				try {
					await ttlockApi.logUnlockActivity({
						property_id: this.propertyId,
						stage,
						page: 'unlock_steps',
						meta: meta || {},
					})
				} catch (e) {}
			},
			async initBluetooth() {
				// #ifndef MP-WEIXIN
				this.ble.supported = false
				this.ble.initStatus = 'unsupported'
				this.ble.available = false
				this.ble.discovering = false
				this.ble.connected = false
				this.ble.deviceId = ''
				this.ble.errorMsg = ''
				this.appendBleLog('warn', '当前平台仅支持微信小程序蓝牙能力')
				return
				// #endif

				// #ifdef MP-WEIXIN
				if (this.ble.initStatus === 'initting' || this.ble.initStatus === 'ready') return
				this.ble.supported = true
				this.ble.initStatus = 'initting'
				this.ble.errorMsg = ''
				this.appendBleLog('info', '开始初始化蓝牙适配器')
				this.reportUnlockActivity('ble_init_start', {})

				try {
					await new Promise((resolve, reject) => {
						uni.openBluetoothAdapter({
							success: resolve,
							fail: reject,
						})
					})
					this.ble.initStatus = 'ready'
					this.appendBleLog('success', '蓝牙适配器初始化成功')
					await this.refreshBluetoothState()
					this.bindBluetoothListeners()
					this.reportUnlockActivity('ble_init_ok', {
						available: this.ble.available,
						discovering: this.ble.discovering,
					})
				} catch (e) {
					this.ble.initStatus = 'failed'
					this.ble.errorMsg = this.formatBleErr(e)
					this.appendBleLog('error', `蓝牙初始化失败：${this.ble.errorMsg || '未知错误'}`)
					this.reportUnlockActivity('ble_init_fail', {
						err_code: Number(e && (e.errCode || e.code) || 0),
						err_msg: this.formatBleErr(e),
					})
				}
				// #endif
			},
			async refreshBluetoothState() {
				// #ifdef MP-WEIXIN
				try {
					const st = await new Promise((resolve, reject) => {
						uni.getBluetoothAdapterState({
							success: resolve,
							fail: reject,
						})
					})
					this.ble.available = !!st.available
					this.ble.discovering = !!st.discovering
					const key = `${this.ble.available}_${this.ble.discovering}`
					this._lastBleStateKey = key
					this.appendBleLog('info', `蓝牙状态：${this.ble.available ? '可用' : '不可用'}${this.ble.discovering ? '，搜索中' : ''}`)
				} catch (e) {
					this.appendBleLog('warn', `获取蓝牙状态失败：${this.formatBleErr(e)}`)
				}
				// #endif
			},
			bindBluetoothListeners() {
				// #ifdef MP-WEIXIN
				if (this._bleListening) return
				this._bleListening = true

				try {
					uni.onBluetoothAdapterStateChange((res) => {
						this.ble.available = !!res.available
						this.ble.discovering = !!res.discovering
						const key = `${this.ble.available}_${this.ble.discovering}`
						if (key === this._lastBleStateKey) return
						this._lastBleStateKey = key
						this.appendBleLog('info', `蓝牙状态变化：${this.ble.available ? '可用' : '不可用'}${this.ble.discovering ? '，搜索中' : ''}`)
						this.reportUnlockActivity('ble_state_change', {
							available: this.ble.available,
							discovering: this.ble.discovering,
						})
					})
				} catch (e) {}

				try {
					uni.onBLEConnectionStateChange((res) => {
						this.ble.connected = !!res.connected
						if (res && res.deviceId) this.ble.deviceId = res.deviceId
						const key = `${this.ble.connected}_${this.ble.deviceId || ''}`
						if (key === this._lastConnKey) return
						this._lastConnKey = key
						this.appendBleLog(this.ble.connected ? 'success' : 'warn', `连接状态：${this.ble.connected ? '已连接' : '已断开'}${this.ble.deviceId ? `，设备 ${this.ble.deviceId}` : ''}`)
						this.reportUnlockActivity('ble_conn_change', {
							connected: this.ble.connected,
							deviceId: this.ble.deviceId,
						})
					})
				} catch (e) {}
				// #endif
			},
			unbindBluetoothListeners() {
				// #ifdef MP-WEIXIN
				if (!this._bleListening) return
				this._bleListening = false
				try { uni.offBluetoothAdapterStateChange() } catch (e) {}
				try { uni.offBLEConnectionStateChange() } catch (e) {}
				// #endif
			},
			cleanupBluetooth() {
				// #ifdef MP-WEIXIN
				this.unbindBluetoothListeners()
				try { uni.closeBluetoothAdapter({}) } catch (e) {}
				// #endif
			},
			async bluetoothUnlock() {
				if (this.unlocking) return
				if (!this.propertyId) {
					uni.showToast({ title: '缺少房源ID，无法开锁', icon: 'none' })
					return
				}

				// #ifndef MP-WEIXIN
				this.appendBleLog('warn', '仅支持微信小程序蓝牙开锁')
				uni.showToast({ title: '仅支持微信小程序蓝牙开锁', icon: 'none' })
				return
				// #endif

				// #ifdef MP-WEIXIN
				this.unlocking = true
				let plugin = null
				let requestId = 0
				let finishSuccess = false
				let finishErrCode = 0
				let finishErrMsg = ''
				let pluginOut = null
				try {
					this.appendBleLog('info', '用户点击蓝牙开锁')
					await this.initBluetooth()
					if (this.ble.initStatus === 'failed' || this.ble.initStatus === 'unsupported') {
						finishErrMsg = this.ble.errorMsg || '蓝牙不可用，请开启后重试'
						throw new Error(finishErrMsg)
					}

					// eslint-disable-next-line no-undef
					plugin = requirePlugin('ttlock-plugin')
					if (!plugin || typeof plugin.controlLock !== 'function') {
						finishErrMsg = '未检测到通通锁插件，请检查小程序插件配置'
						throw new Error(finishErrMsg)
					}

					const startRes = await ttlockApi.bleUnlockStart({
						property_id: this.propertyId,
						meta: {
							page: 'unlock_steps',
							ble: {
								available: this.ble.available,
								connected: this.ble.connected,
								deviceId: this.ble.deviceId,
							},
						},
					})
					requestId = Number(startRes && startRes.data && startRes.data.request_id) || 0
					if (!requestId) {
						finishErrMsg = (startRes && startRes.message) ? startRes.message : '创建开锁记录失败'
						throw new Error(finishErrMsg)
					}
					this.unlockRequestId = requestId
					this.appendBleLog('info', `已创建开锁记录：#${requestId}`)

					this.appendBleLog('info', '正在获取开锁数据')
					const res = await ttlockApi.getPropertyBleInfo({ property_id: this.propertyId })
					if (!res || res.code !== 0) {
						finishErrMsg = (res && res.message) ? res.message : '获取开锁数据失败'
						throw new Error(finishErrMsg)
					}
					const data = res.data || {}
					const lockData = String(data.key_data || data.lock_data || '').trim()
					if (!lockData) {
						finishErrMsg = '缺少开锁数据，请先同步钥匙信息'
						throw new Error(finishErrMsg)
					}

					uni.showLoading({ title: '蓝牙开锁中', mask: true })
					this.appendBleLog('info', '开始调用蓝牙插件开锁（连接中）')
					pluginOut = await plugin.controlLock({
						controlAction: 3,
						lockData,
						serverTime: Date.now(),
					})
					uni.hideLoading()

					if (pluginOut && Number(pluginOut.errorCode) === 0) {
						finishSuccess = true
						finishErrCode = 0
						finishErrMsg = ''
						this.appendBleLog('success', '开锁成功')
						uni.showToast({ title: '开锁成功', icon: 'success' })
					} else {
						finishSuccess = false
						finishErrCode = Number(pluginOut && pluginOut.errorCode) || 0
						const msg = (pluginOut && (pluginOut.errorMsg || pluginOut.errMsg || pluginOut.description)) || '开锁失败，请重试'
						finishErrMsg = msg
						this.appendBleLog('error', `开锁失败：${msg}`)
						uni.showToast({ title: msg, icon: 'none' })
					}
				} catch (e) {
					try { uni.hideLoading() } catch (err) {}
					const msg = (e && e.message) ? e.message : '开锁失败，请稍后重试'
					if (!finishErrMsg) finishErrMsg = msg
					this.appendBleLog('error', msg)
					uni.showToast({ title: msg, icon: 'none' })
				} finally {
					if (requestId > 0) {
						try {
							await ttlockApi.bleUnlockFinish({
								property_id: this.propertyId,
								request_id: requestId,
								success: finishSuccess ? 1 : 0,
								err_code: finishSuccess ? 0 : finishErrCode,
								err_msg: finishSuccess ? '' : finishErrMsg,
								meta: {
									page: 'unlock_steps',
									ble: {
										available: this.ble.available,
										connected: this.ble.connected,
										deviceId: this.ble.deviceId,
									},
									plugin: pluginOut ? {
										errorCode: pluginOut.errorCode,
										errorMsg: pluginOut.errorMsg || pluginOut.errMsg || pluginOut.description || '',
									} : null,
								},
							})
						} catch (e) {
							this.appendBleLog('warn', '写入开锁结果日志失败（不影响开锁结果）')
						}
					}
					try {
						if (plugin && typeof plugin.stopAllOperations === 'function') {
							await plugin.stopAllOperations()
						}
					} catch (e) {}
					this.unlocking = false
				}
				// #endif
			},
			passwordUnlock() {
				uni.navigateTo({
					url: `/pages/unlock_details/unlock_details?property_id=${encodeURIComponent(this.propertyId)}`
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

	.ble-section {
		padding: 0 32rpx 24rpx;
	}

	.ble-card {
		background-color: #ffffff;
		border: 1px solid #f1f5f9;
		border-radius: 28rpx;
		padding: 24rpx;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
		display: flex;
		flex-direction: column;
		gap: 18rpx;

		.ble-head {
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 16rpx;
		}

		.ble-left {
			display: flex;
			align-items: center;
			gap: 12rpx;
		}

		.ble-ico { font-size: 34rpx; color: #2d9cf0; }
		.ble-title { font-size: 28rpx; font-weight: 800; color: #0f172a; }

		.ble-badge {
			padding: 8rpx 16rpx;
			border-radius: 999rpx;
			font-size: 22rpx;
			font-weight: 700;
			border: 1px solid transparent;

			&.ok { background-color: #ecfdf5; color: #065f46; border-color: #d1fae5; }
			&.warn { background-color: #fff7ed; color: #9a3412; border-color: #fed7aa; }
			&.muted { background-color: #f8fafc; color: #64748b; border-color: #e2e8f0; }
		}

		.ble-grid {
			display: flex;
			flex-wrap: wrap;
			row-gap: 16rpx;
		}

		.ble-item {
			width: 50%;
			display: flex;
			flex-direction: column;
			gap: 6rpx;
		}

		.ble-lab { font-size: 22rpx; color: #94a3b8; font-weight: 700; }
		.ble-val { font-size: 26rpx; color: #334155; font-weight: 800; }
		.ble-val.ok { color: #16a34a; }
		.ble-val.warn { color: #f97316; }
		.ble-val.loading { color: #2d9cf0; }
		.ble-val.muted { color: #64748b; }

		.ble-extra {
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 16rpx;
			padding-top: 6rpx;

			.ble-extra-lab { font-size: 22rpx; color: #94a3b8; font-weight: 700; }
			.ble-extra-val { font-size: 24rpx; color: #334155; font-weight: 700; }
		}

		.ble-error {
			display: flex;
			align-items: flex-start;
			gap: 12rpx;
			padding: 16rpx;
			background-color: #fef2f2;
			border: 1px solid #fee2e2;
			border-radius: 18rpx;

			.err-ico { font-size: 28rpx; color: #ef4444; }
			.err-text { font-size: 24rpx; color: #991b1b; font-weight: 700; line-height: 1.5; }
		}

		.ble-log {
			background-color: #f8fafc;
			border: 1px solid #e2e8f0;
			border-radius: 20rpx;
			overflow: hidden;

			.ble-log-head {
				display: flex;
				align-items: center;
				justify-content: space-between;
				padding: 16rpx 18rpx;
				border-bottom: 1px solid #e2e8f0;
			}

			.ble-log-title { font-size: 24rpx; font-weight: 800; color: #0f172a; }
			.ble-log-tip { font-size: 22rpx; color: #94a3b8; font-weight: 700; }

			.ble-log-body {
				height: 220rpx;
				padding: 12rpx 18rpx 16rpx;
				box-sizing: border-box;
			}

			.ble-log-line {
				display: flex;
				gap: 14rpx;
				padding: 6rpx 0;
				line-height: 1.4;
			}

			.ble-log-time { width: 110rpx; flex-shrink: 0; font-size: 20rpx; color: #94a3b8; }
			.ble-log-msg { flex: 1; font-size: 22rpx; color: #334155; font-weight: 600; }
			.ble-log-msg.success { color: #16a34a; }
			.ble-log-msg.warn { color: #f97316; }
			.ble-log-msg.error { color: #ef4444; }
		}
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
			padding: 16rpx 32rpx;
			border-radius: 40rpx;
			backdrop-filter: blur(10px);

			&.ok { background-color: #ecfdf5; border: 1px solid #d1fae5; }
			&.warn { background-color: #fff7ed; border: 1px solid #fed7aa; }
			&.loading { background-color: #eff6ff; border: 1px solid #bfdbfe; }
			&.muted { background-color: #f8fafc; border: 1px solid #e2e8f0; }
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

		.status-tag.warn .dot-ping,
		.status-tag.warn .dot-inner { background-color: #f97316; }
		.status-tag.loading .dot-ping,
		.status-tag.loading .dot-inner { background-color: #2d9cf0; }
		.status-tag.muted .dot-ping,
		.status-tag.muted .dot-inner { background-color: #94a3b8; }

		.status-text { font-size: 24rpx; font-weight: 600; }
		.status-tag.ok .status-text { color: #065f46; }
		.status-tag.warn .status-text { color: #9a3412; }
		.status-tag.loading .status-text { color: #1d4ed8; }
		.status-tag.muted .status-text { color: #475569; }
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
