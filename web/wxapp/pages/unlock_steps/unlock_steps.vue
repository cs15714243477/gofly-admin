<template>
  <view class="steps-container">
    <TopHeader title="开锁步骤">
      <template #left>
        <view class="th-btn" @click="goBack">
          <text class="material-symbols-outlined">arrow_back_ios_new</text>
        </view>
      </template>
    </TopHeader>

    <view class="main-content">
      <view class="tip-card">
        <text class="material-symbols-outlined tip-icon">near_me</text>
        <text class="tip-text">请站在门锁 1-2 米内，并保持蓝牙开启。</text>
      </view>

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

        <view class="status-tag" :class="statusTagClass">
          <view class="dot-box">
            <view class="dot-ping"></view>
            <view class="dot-inner"></view>
          </view>
          <text class="status-text">{{ statusText }}</text>
        </view>

        <view class="ble-grid">
          <view class="ble-item">
            <text class="ble-lab">初始化</text>
            <text class="ble-val" :class="bleInitClass">{{ bleInitText }}</text>
          </view>
          <view class="ble-item">
            <text class="ble-lab">可用</text>
            <text class="ble-val">{{ ble.available ? "是" : "否" }}</text>
          </view>
          <view class="ble-item">
            <text class="ble-lab">搜索中</text>
            <text class="ble-val">{{ ble.discovering ? "是" : "否" }}</text>
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
            <text class="ble-log-title">最近日志</text>
            <text class="ble-log-tip"
              >{{ latestBleLogs.length }}/{{ bleLogMax }}</text
            >
          </view>
          <view class="ble-log-body">
            <view
              v-for="(it, idx) in latestBleLogs"
              :key="idx"
              class="ble-log-line"
            >
              <text class="ble-log-time">{{ it.time }}</text>
              <text class="ble-log-msg" :class="it.level">{{ it.msg }}</text>
            </view>
            <view v-if="!latestBleLogs.length" class="ble-log-empty"
              >暂无日志</view
            >
          </view>
        </view>
      </view>

      <view class="steps-section">
        <view class="steps-title">开锁流程</view>
        <view class="steps-grid">
          <view class="step-card">
            <view class="step-icon-box primary">
              <text class="material-symbols-outlined">bluetooth</text>
            </view>
            <text class="step-index">步骤 1</text>
            <text class="step-title">打开蓝牙</text>
            <text class="step-desc">确认系统蓝牙已开启</text>
          </view>
          <view class="step-card">
            <view class="step-icon-box accent">
              <text class="material-symbols-outlined">touch_app</text>
            </view>
            <text class="step-index">步骤 2</text>
            <text class="step-title">唤醒门锁</text>
            <text class="step-desc">触碰门锁点亮数字面板</text>
          </view>
          <view class="step-card">
            <view class="step-icon-box primary">
              <text class="material-symbols-outlined">phonelink_lock</text>
            </view>
            <text class="step-index">步骤 3</text>
            <text class="step-title">点击开锁</text>
            <text class="step-desc">点击下方蓝牙开锁按钮</text>
          </view>
        </view>
      </view>

      <view class="action-panel">
        <button class="bluetooth-btn" @click="bluetoothUnlock">
          <text class="material-symbols-outlined bt-icon pulse"
            >bluetooth_connected</text
          >
          <text>蓝牙开锁</text>
        </button>

        <button class="password-btn" @click="passwordUnlock">
          <text class="material-symbols-outlined btn-icon">lock_open</text>
          <text>密码开锁</text>
        </button>
      </view>
    </view>
  </view>
</template>

<script>
import TopHeader from "@/components/TopHeader.vue";
import ttlockApi from "@/api/ttlock";

function pad2(n) {
  return String(n).padStart(2, "0");
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
        initStatus: "idle", // idle/initting/ready/failed/unsupported
        available: false,
        discovering: false,
        connected: false,
        deviceId: "",
        errorMsg: "",
      },
      bleLogs: [],
      bleLogMax: 50,
      _bleListening: false,
      _lastBleStateKey: "",
      _lastConnKey: "",
    };
  },
  computed: {
    bleInitText() {
      if (!this.ble.supported) return "不支持";
      switch (this.ble.initStatus) {
        case "initting":
          return "初始化中";
        case "ready":
          return "已就绪";
        case "failed":
          return "失败";
        case "unsupported":
          return "不支持";
        default:
          return "未初始化";
      }
    },
    bleInitClass() {
      if (this.ble.initStatus === "ready") return "ok";
      if (this.ble.initStatus === "failed") return "warn";
      if (this.ble.initStatus === "initting") return "loading";
      return "muted";
    },
    bleBadgeText() {
      if (!this.ble.supported || this.ble.initStatus === "unsupported")
        return "当前平台不支持";
      if (this.ble.initStatus === "failed") return "蓝牙不可用";
      if (this.ble.available) return "可用";
      return "未开启/不可用";
    },
    bleBadgeClass() {
      if (!this.ble.supported || this.ble.initStatus === "unsupported")
        return "muted";
      if (this.ble.initStatus === "failed") return "warn";
      if (this.ble.available) return "ok";
      return "warn";
    },
    bleConnectText() {
      if (this.ble.connected) return "已连接";
      if (this.unlocking) return "连接中";
      return "未连接";
    },
    statusText() {
      if (!this.ble.supported || this.ble.initStatus === "unsupported")
        return "当前平台不支持蓝牙能力";
      if (this.ble.initStatus === "initting") return "正在初始化蓝牙…";
      if (this.ble.initStatus === "failed")
        return this.ble.errorMsg
          ? `蓝牙初始化失败：${this.ble.errorMsg}`
          : "蓝牙初始化失败";
      if (!this.ble.available) return "请开启手机蓝牙后重试";
      if (this.unlocking) return "正在连接并开锁…";
      if (this.ble.connected) return "设备已连接";
      return "蓝牙就绪，等待连接";
    },
    statusTagClass() {
      if (!this.ble.supported || this.ble.initStatus === "unsupported")
        return "muted";
      if (this.ble.initStatus === "failed" || !this.ble.available)
        return "warn";
      if (this.unlocking || this.ble.initStatus === "initting")
        return "loading";
      return "ok";
    },
    latestBleLogs() {
      return this.bleLogs.slice(-3);
    },
  },
  onLoad(options) {
    const pid = Number((options && (options.property_id || options.id)) || 0);
    this.propertyId = Number.isFinite(pid) ? pid : 0;
  },
  onShow() {
    this.appendBleLog("info", "进入页面，开始检测蓝牙状态");
    this.initBluetooth();
  },
  onHide() {
    this.cleanupBluetooth();
  },
  onUnload() {
    this.cleanupBluetooth();
  },
  methods: {
    goBack() {
      uni.navigateBack();
    },
    appendBleLog(level, msg) {
      const d = new Date();
      const time = `${pad2(d.getHours())}:${pad2(d.getMinutes())}:${pad2(
        d.getSeconds()
      )}`;
      this.bleLogs.push({
        time,
        msg: String(msg || ""),
        level: level || "info",
      });
      if (this.bleLogs.length > this.bleLogMax) {
        this.bleLogs.splice(0, this.bleLogs.length - this.bleLogMax);
      }
    },
    formatBleErr(err) {
      if (!err) return "";
      const code = Number(err.errCode || err.code || 0);
      const raw = String(err.errMsg || err.message || "").trim();
      if (code === 10001) return "蓝牙适配器不可用（请开启手机蓝牙）";
      if (code === 10002) return "没有找到指定设备";
      if (code === 10003) return "连接失败";
      if (code === 10004) return "没有找到指定服务";
      if (code === 10005) return "没有找到指定特征值";
      if (code === 10006) return "当前连接已断开";
      if (code === 10007) return "当前特征值不支持该操作";
      if (code === 10008) return "其余所有系统上报的异常";
      if (code === 10009) return "系统版本低于 4.3 不支持 BLE";
      if (code === 10010) return "系统版本低于 4.4 不支持 BLE";
      if (code === 10011) return "需要授权蓝牙权限";
      if (raw) return raw;
      if (code) return `蓝牙异常（错误码：${code}）`;
      try {
        return JSON.stringify(err);
      } catch (e) {}
      return "蓝牙异常";
    },
    async reportUnlockActivity(stage, meta) {
      if (!this.propertyId) return;
      if (!ttlockApi || typeof ttlockApi.logUnlockActivity !== "function")
        return;
      try {
        await ttlockApi.logUnlockActivity({
          property_id: this.propertyId,
          stage,
          page: "unlock_steps",
          meta: meta || {},
        });
      } catch (e) {}
    },
    async initBluetooth() {
      // #ifndef MP-WEIXIN
      this.ble.supported = false;
      this.ble.initStatus = "unsupported";
      this.ble.available = false;
      this.ble.discovering = false;
      this.ble.connected = false;
      this.ble.deviceId = "";
      this.ble.errorMsg = "";
      this.appendBleLog("warn", "当前平台仅支持微信小程序蓝牙能力");
      return;
      // #endif

      // #ifdef MP-WEIXIN
      if (this.ble.initStatus === "initting" || this.ble.initStatus === "ready")
        return;
      this.ble.supported = true;
      this.ble.initStatus = "initting";
      this.ble.errorMsg = "";
      this.appendBleLog("info", "开始初始化蓝牙适配器");
      this.reportUnlockActivity("ble_init_start", {});

      try {
        await new Promise((resolve, reject) => {
          uni.openBluetoothAdapter({
            success: resolve,
            fail: reject,
          });
        });
        this.ble.initStatus = "ready";
        this.appendBleLog("success", "蓝牙适配器初始化成功");
        await this.refreshBluetoothState();
        this.bindBluetoothListeners();
        this.reportUnlockActivity("ble_init_ok", {
          available: this.ble.available,
          discovering: this.ble.discovering,
        });
      } catch (e) {
        this.ble.initStatus = "failed";
        this.ble.errorMsg = this.formatBleErr(e);
        this.appendBleLog(
          "error",
          `蓝牙初始化失败：${this.ble.errorMsg || "未知错误"}`
        );
        this.reportUnlockActivity("ble_init_fail", {
          err_code: Number((e && (e.errCode || e.code)) || 0),
          err_msg: this.formatBleErr(e),
        });
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
          });
        });
        this.ble.available = !!st.available;
        this.ble.discovering = !!st.discovering;
        const key = `${this.ble.available}_${this.ble.discovering}`;
        this._lastBleStateKey = key;
        this.appendBleLog(
          "info",
          `蓝牙状态：${this.ble.available ? "可用" : "不可用"}${
            this.ble.discovering ? "，搜索中" : ""
          }`
        );
      } catch (e) {
        this.appendBleLog("warn", `获取蓝牙状态失败：${this.formatBleErr(e)}`);
      }
      // #endif
    },
    bindBluetoothListeners() {
      // #ifdef MP-WEIXIN
      if (this._bleListening) return;
      this._bleListening = true;

      try {
        uni.onBluetoothAdapterStateChange((res) => {
          this.ble.available = !!res.available;
          this.ble.discovering = !!res.discovering;
          const key = `${this.ble.available}_${this.ble.discovering}`;
          if (key === this._lastBleStateKey) return;
          this._lastBleStateKey = key;
          this.appendBleLog(
            "info",
            `蓝牙状态变化：${this.ble.available ? "可用" : "不可用"}${
              this.ble.discovering ? "，搜索中" : ""
            }`
          );
          this.reportUnlockActivity("ble_state_change", {
            available: this.ble.available,
            discovering: this.ble.discovering,
          });
        });
      } catch (e) {}

      try {
        uni.onBLEConnectionStateChange((res) => {
          this.ble.connected = !!res.connected;
          if (res && res.deviceId) this.ble.deviceId = res.deviceId;
          const key = `${this.ble.connected}_${this.ble.deviceId || ""}`;
          if (key === this._lastConnKey) return;
          this._lastConnKey = key;
          this.appendBleLog(
            this.ble.connected ? "success" : "warn",
            `连接状态：${this.ble.connected ? "已连接" : "已断开"}${
              this.ble.deviceId ? `，设备 ${this.ble.deviceId}` : ""
            }`
          );
          this.reportUnlockActivity("ble_conn_change", {
            connected: this.ble.connected,
            deviceId: this.ble.deviceId,
          });
        });
      } catch (e) {}
      // #endif
    },
    unbindBluetoothListeners() {
      // #ifdef MP-WEIXIN
      if (!this._bleListening) return;
      this._bleListening = false;
      try {
        uni.offBluetoothAdapterStateChange();
      } catch (e) {}
      try {
        uni.offBLEConnectionStateChange();
      } catch (e) {}
      // #endif
    },
    cleanupBluetooth() {
      // #ifdef MP-WEIXIN
      this.unbindBluetoothListeners();
      try {
        uni.closeBluetoothAdapter({});
      } catch (e) {}
      // #endif
    },
    async bluetoothUnlock() {
      if (this.unlocking) return;
      if (!this.propertyId) {
        uni.showToast({ title: "缺少房源ID，无法开锁", icon: "none" });
        return;
      }

      // #ifndef MP-WEIXIN
      this.appendBleLog("warn", "仅支持微信小程序蓝牙开锁");
      uni.showToast({ title: "仅支持微信小程序蓝牙开锁", icon: "none" });
      return;
      // #endif

      // #ifdef MP-WEIXIN
      this.unlocking = true;
      let plugin = null;
      let requestId = 0;
      let finishSuccess = false;
      let finishErrCode = 0;
      let finishErrMsg = "";
      let pluginOut = null;
      try {
        this.appendBleLog("info", "用户点击蓝牙开锁");
        await this.initBluetooth();
        if (
          this.ble.initStatus === "failed" ||
          this.ble.initStatus === "unsupported"
        ) {
          finishErrMsg = this.ble.errorMsg || "蓝牙不可用，请开启后重试";
          throw new Error(finishErrMsg);
        }

        // eslint-disable-next-line no-undef
        plugin = requirePlugin("ttlock-plugin");
        if (!plugin || typeof plugin.controlLock !== "function") {
          finishErrMsg = "未检测到通通锁插件，请检查小程序插件配置";
          throw new Error(finishErrMsg);
        }

        const startRes = await ttlockApi.bleUnlockStart({
          property_id: this.propertyId,
          meta: {
            page: "unlock_steps",
            ble: {
              available: this.ble.available,
              connected: this.ble.connected,
              deviceId: this.ble.deviceId,
            },
          },
        });
        requestId =
          Number(startRes && startRes.data && startRes.data.request_id) || 0;
        if (!requestId) {
          finishErrMsg =
            startRes && startRes.message
              ? startRes.message
              : "创建开锁记录失败";
          throw new Error(finishErrMsg);
        }
        this.unlockRequestId = requestId;
        this.appendBleLog("info", `已创建开锁记录：#${requestId}`);

        this.appendBleLog("info", "正在获取开锁数据");
        const res = await ttlockApi.getPropertyBleInfo({
          property_id: this.propertyId,
        });
        if (!res || res.code !== 0) {
          finishErrMsg = res && res.message ? res.message : "获取开锁数据失败";
          throw new Error(finishErrMsg);
        }
        const data = res.data || {};
        const lockData = String(data.key_data || data.lock_data || "").trim();
        if (!lockData) {
          finishErrMsg = "缺少开锁数据，请先同步钥匙信息";
          throw new Error(finishErrMsg);
        }

        uni.showLoading({ title: "蓝牙开锁中", mask: true });
        this.appendBleLog("info", "开始调用蓝牙插件开锁（连接中）");
        pluginOut = await plugin.controlLock({
          controlAction: 3,
          lockData,
          serverTime: Date.now(),
        });
        uni.hideLoading();

        if (pluginOut && Number(pluginOut.errorCode) === 0) {
          finishSuccess = true;
          finishErrCode = 0;
          finishErrMsg = "";
          this.appendBleLog("success", "开锁成功");
          uni.showToast({ title: "开锁成功", icon: "success" });
        } else {
          finishSuccess = false;
          finishErrCode = Number(pluginOut && pluginOut.errorCode) || 0;
          const msg =
            (pluginOut &&
              (pluginOut.errorMsg ||
                pluginOut.errMsg ||
                pluginOut.description)) ||
            "开锁失败，请重试";
          finishErrMsg = msg;
          this.appendBleLog("error", `开锁失败：${msg}`);
          uni.showToast({ title: msg, icon: "none" });
        }
      } catch (e) {
        try {
          uni.hideLoading();
        } catch (err) {}
        const msg = e && e.message ? e.message : "开锁失败，请稍后重试";
        if (!finishErrMsg) finishErrMsg = msg;
        this.appendBleLog("error", msg);
        uni.showToast({ title: msg, icon: "none" });
      } finally {
        if (requestId > 0) {
          try {
            await ttlockApi.bleUnlockFinish({
              property_id: this.propertyId,
              request_id: requestId,
              success: finishSuccess ? 1 : 0,
              err_code: finishSuccess ? 0 : finishErrCode,
              err_msg: finishSuccess ? "" : finishErrMsg,
              meta: {
                page: "unlock_steps",
                ble: {
                  available: this.ble.available,
                  connected: this.ble.connected,
                  deviceId: this.ble.deviceId,
                },
                plugin: pluginOut
                  ? {
                      errorCode: pluginOut.errorCode,
                      errorMsg:
                        pluginOut.errorMsg ||
                        pluginOut.errMsg ||
                        pluginOut.description ||
                        "",
                    }
                  : null,
              },
            });
          } catch (e) {
            this.appendBleLog("warn", "写入开锁结果日志失败（不影响开锁结果）");
          }
        }
        try {
          if (plugin && typeof plugin.stopAllOperations === "function") {
            await plugin.stopAllOperations();
          }
        } catch (e) {}
        this.unlocking = false;
      }
      // #endif
    },
    passwordUnlock() {
      uni.navigateTo({
        url: `/pages/unlock_details/unlock_details?property_id=${encodeURIComponent(
          this.propertyId
        )}`,
      });
    },
  },
};
</script>

<style lang="scss">
.steps-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: linear-gradient(180deg, #f3f9ff 0%, #f8fafc 45%, #ffffff 100%);
}

.th-btn {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #0f172a;
  background: transparent;

  .material-symbols-outlined {
    font-size: 32rpx;
  }

  &:active {
    background: rgba(15, 23, 42, 0.06);
  }
}

.main-content {
  flex: 1;
  min-height: 0;
  padding: 16rpx 20rpx calc(env(safe-area-inset-bottom) + 20rpx);
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  overflow: hidden;
}

.tip-card {
  padding: 16rpx 18rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  border-radius: 18rpx;

  .tip-icon {
    font-size: 30rpx;
    color: #2563eb;
  }

  .tip-text {
    flex: 1;
    min-width: 0;
    font-size: 24rpx;
    line-height: 1.4;
    font-weight: 600;
    color: #334155;
  }
}

.ble-card {
  padding: 18rpx;
  border-radius: 20rpx;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 8rpx 24rpx rgba(15, 23, 42, 0.05);
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.ble-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.ble-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.ble-ico {
  font-size: 30rpx;
  color: #2d9cf0;
}

.ble-title {
  font-size: 26rpx;
  font-weight: 800;
  color: #0f172a;
}

.ble-badge {
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  font-weight: 700;
  border: 1px solid transparent;

  &.ok {
    background-color: #ecfdf5;
    color: #065f46;
    border-color: #d1fae5;
  }

  &.warn {
    background-color: #fff7ed;
    color: #9a3412;
    border-color: #fed7aa;
  }

  &.muted {
    background-color: #f8fafc;
    color: #64748b;
    border-color: #e2e8f0;
  }
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 12rpx 14rpx;
  border-radius: 14rpx;
  border: 1px solid transparent;
  min-height: 64rpx;

  &.ok {
    background-color: #ecfdf5;
    border-color: #d1fae5;
  }

  &.warn {
    background-color: #fff7ed;
    border-color: #fed7aa;
  }

  &.loading {
    background-color: #eff6ff;
    border-color: #bfdbfe;
  }

  &.muted {
    background-color: #f8fafc;
    border-color: #e2e8f0;
  }
}

.dot-box {
  width: 18rpx;
  height: 18rpx;
  position: relative;
  flex-shrink: 0;
}

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

.status-tag.warn .dot-ping,
.status-tag.warn .dot-inner {
  background-color: #f97316;
}

.status-tag.loading .dot-ping,
.status-tag.loading .dot-inner {
  background-color: #2d9cf0;
}

.status-tag.muted .dot-ping,
.status-tag.muted .dot-inner {
  background-color: #94a3b8;
}

.status-text {
  flex: 1;
  min-width: 0;
  font-size: 22rpx;
  line-height: 1.35;
  font-weight: 600;
  color: #334155;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.ble-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10rpx;
}

.ble-item {
  padding: 10rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  align-items: center;
  justify-content: center;
}

.ble-lab {
  font-size: 18rpx;
  color: #94a3b8;
  font-weight: 700;
}

.ble-val {
  font-size: 22rpx;
  color: #334155;
  font-weight: 800;
}

.ble-val.ok {
  color: #16a34a;
}

.ble-val.warn {
  color: #f97316;
}

.ble-val.loading {
  color: #2d9cf0;
}

.ble-val.muted {
  color: #64748b;
}

.ble-extra {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  padding: 10rpx 12rpx;
  border-radius: 12rpx;
  background: #f8fafc;
}

.ble-extra-lab {
  font-size: 20rpx;
  color: #64748b;
  font-weight: 700;
}

.ble-extra-val {
  flex: 1;
  min-width: 0;
  text-align: right;
  font-size: 20rpx;
  color: #1e293b;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ble-error {
  display: flex;
  align-items: flex-start;
  gap: 10rpx;
  padding: 12rpx;
  border-radius: 12rpx;
  border: 1px solid #fecaca;
  background: #fef2f2;
}

.err-ico {
  font-size: 24rpx;
  color: #dc2626;
}

.err-text {
  flex: 1;
  min-width: 0;
  font-size: 20rpx;
  line-height: 1.35;
  color: #b91c1c;
  font-weight: 700;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.ble-log {
  border-radius: 14rpx;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  overflow: hidden;
}

.ble-log-head {
  padding: 10rpx 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e2e8f0;
}

.ble-log-title {
  font-size: 20rpx;
  font-weight: 800;
  color: #0f172a;
}

.ble-log-tip {
  font-size: 18rpx;
  color: #94a3b8;
  font-weight: 700;
}

.ble-log-body {
  padding: 8rpx 12rpx 10rpx;
}

.ble-log-line {
  display: flex;
  gap: 8rpx;
  padding: 4rpx 0;
}

.ble-log-time {
  width: 94rpx;
  flex-shrink: 0;
  font-size: 18rpx;
  color: #94a3b8;
}

.ble-log-msg {
  flex: 1;
  min-width: 0;
  font-size: 18rpx;
  color: #334155;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ble-log-msg.success {
  color: #16a34a;
}

.ble-log-msg.warn {
  color: #f97316;
}

.ble-log-msg.error {
  color: #ef4444;
}

.ble-log-empty {
  font-size: 18rpx;
  color: #94a3b8;
  padding: 4rpx 0;
}

.steps-section {
  padding: 14rpx;
  border-radius: 18rpx;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 6rpx 20rpx rgba(15, 23, 42, 0.04);
}

.steps-title {
  font-size: 22rpx;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 12rpx;
}

.steps-grid {
  display: flex;
  gap: 10rpx;
}

.step-card {
  flex: 1;
  min-width: 0;
  padding: 12rpx 8rpx;
  border-radius: 14rpx;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 6rpx;
}

.step-icon-box {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid transparent;

  &.primary {
    border-color: #2d9cf0;
    color: #2d9cf0;
  }

  &.accent {
    border-color: #f97316;
    color: #f97316;
  }

  .material-symbols-outlined {
    font-size: 28rpx;
  }
}

.step-index {
  font-size: 18rpx;
  color: #64748b;
  font-weight: 700;
}

.step-title {
  font-size: 22rpx;
  color: #0f172a;
  font-weight: 800;
}

.step-desc {
  font-size: 18rpx;
  color: #64748b;
  line-height: 1.35;
}

.action-panel {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.bluetooth-btn {
  height: 84rpx;
  background: linear-gradient(135deg, #2d9cf0 0%, #1777d7 100%);
  color: #ffffff;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  font-size: 30rpx;
  font-weight: 800;
  border: none;
  box-shadow: 0 8rpx 24rpx rgba(37, 99, 235, 0.28);

  &::after {
    border: none;
  }

  &:active {
    transform: scale(0.98);
  }
}

.bt-icon {
  font-size: 36rpx;
}

.password-btn {
  height: 84rpx;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  color: #334155;
  font-size: 28rpx;
  font-weight: 800;

  &::after {
    border: none;
  }

  &:active {
    transform: scale(0.98);
    background: #f1f5f9;
  }
}

.btn-icon {
  font-size: 34rpx;
}

@media screen and (max-height: 700px) {
  .main-content {
    gap: 10rpx;
    padding-top: 10rpx;
  }

  .tip-card {
    padding: 12rpx 14rpx;
  }

  .ble-card {
    padding: 14rpx;
    gap: 10rpx;
  }

  .ble-log {
    display: none;
  }

  .steps-section {
    padding: 10rpx;
  }

  .step-card {
    padding: 10rpx 6rpx;
    gap: 4rpx;
  }

  .bluetooth-btn {
    height: 76rpx;
  }

  .password-btn {
    height: 64rpx;
  }
}

@keyframes ping {
  0% {
    transform: scale(1);
    opacity: 1;
  }

  75%,
  100% {
    transform: scale(2.5);
    opacity: 0;
  }
}

.pulse {
  animation: pulse-icon 2s infinite;
}

@keyframes pulse-icon {
  0%,
  100% {
    opacity: 1;
    transform: scale(1);
  }

  50% {
    opacity: 0.8;
    transform: scale(1.08);
  }
}
</style>
