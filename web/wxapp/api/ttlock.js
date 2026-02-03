import request from '@/utils/request';

export default {
  // TTLock 配置状态（仅有智能锁管理权限用户可访问）
  getStatus: () =>
    request({
      url: '/houses/wxttlock/getStatus',
      method: 'GET',
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 同步锁列表（从云端拉取并落库）
  syncLocks: () =>
    request({
      url: '/houses/wxttlock/syncLocks',
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
      url: '/houses/wxttlock/getLockList',
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
      url: '/houses/wxttlock/bindProperty',
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
      url: '/houses/wxttlock/unbindProperty',
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
      url: '/houses/wxttlock/getPropertyLock',
      method: 'GET',
      data,
      custom: {
        showLoading: false,
        auth: true,
      },
    }),

  // 远程开锁（传 property_id 或 ttlock_lock_id）
  remoteUnlock: (data) =>
    request({
      url: '/houses/wxttlock/remoteUnlock',
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
      url: '/houses/wxttlock/getLockDetail',
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
      url: '/houses/wxttlock/getLockRecords',
      method: 'GET',
      data,
      custom: {
        showLoading: true,
        auth: true,
        loadingMsg: '加载中',
      },
    }),
};

