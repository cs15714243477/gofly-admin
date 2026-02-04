import request from '@/utils/request';

export default {
  // TTLock 配置状态（仅有智能锁管理权限用户可访问）
  getStatus: () =>
    request({
      url: '/uniapp/wxttlock/getStatus',
      method: 'GET',
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 同步锁列表（从云端拉取并落库）
  syncLocks: () =>
    request({
      url: '/uniapp/wxttlock/syncLocks',
      method: 'POST',
      data: {},
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '同步中',
      },
    }),

  // 锁列表（本地库）
  getLockList: (data) =>
    request({
      url: '/uniapp/wxttlock/getLockList',
      method: 'GET',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '加载中',
      },
    }),

  // 绑定房源到锁（一房一锁/一锁一房）
  bindProperty: (data) =>
    request({
      url: '/uniapp/wxttlock/bindProperty',
      method: 'POST',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '绑定中',
      },
    }),

  // 解绑房源锁
  unbindProperty: (data) =>
    request({
      url: '/uniapp/wxttlock/unbindProperty',
      method: 'POST',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '解绑中',
      },
    }),

  // 查询房源绑定的锁
  getPropertyLock: (data) =>
    request({
      url: '/uniapp/wxttlock/getPropertyLock',
      method: 'GET',
      data,
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 小程序端：获取房源蓝牙开锁所需数据（lockData/keyData 等）
  getPropertyBleInfo: (data) =>
    request({
      url: '/uniapp/wxttlock/getPropertyBleInfo',
      method: 'GET',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '加载中',
      },
    }),

  // 小程序端：生成房源开锁密码（一次性/临时）
  getPropertyPasscode: (data) =>
    request({
      url: '/uniapp/wxttlock/passcode',
      method: 'POST',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '生成中',
      },
    }),

  // 小程序端：记录开锁/蓝牙关键节点日志（写入 business_user_activity_logs）
  logUnlockActivity: (data) =>
    request({
      url: '/uniapp/wxttlock/logUnlockActivity',
      method: 'POST',
      data,
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 小程序端：开始一次蓝牙开锁（写入 business_unlock_requests）
  bleUnlockStart: (data) =>
    request({
      url: '/uniapp/wxttlock/bleUnlockStart',
      method: 'POST',
      data,
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 小程序端：结束一次蓝牙开锁（更新 business_unlock_requests + 记录活动日志）
  bleUnlockFinish: (data) =>
    request({
      url: '/uniapp/wxttlock/bleUnlockFinish',
      method: 'POST',
      data,
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 远程开锁（传 property_id 或 ttlock_lock_id）
  remoteUnlock: (data) =>
    request({
      url: '/uniapp/wxttlock/remoteUnlock',
      method: 'POST',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '开锁中',
      },
    }),

  // 锁详情（云端+本地）
  getLockDetail: (data) =>
    request({
      url: '/uniapp/wxttlock/getLockDetail',
      method: 'GET',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '加载中',
      },
    }),

  // 开锁记录（云端）
  getLockRecords: (data) =>
    request({
      url: '/uniapp/wxttlock/getLockRecords',
      method: 'GET',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '加载中',
      },
    }),
};
