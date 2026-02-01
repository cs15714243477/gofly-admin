import { defHttp } from '@/utils/http';

enum Api {
  status = '/houses/ttlock/status',
  syncLocks = '/houses/ttlock/syncLocks',
  lockList = '/houses/ttlock/lockList',
  lockRecords = '/houses/ttlock/lockRecords',
  passcodeEvents = '/houses/ttlock/passcodeEvents',
  bindProperty = '/houses/ttlock/bindProperty',
  unbindProperty = '/houses/ttlock/unbindProperty',
  propertyLock = '/houses/ttlock/propertyLock',
  remoteUnlock = '/houses/ttlock/remoteUnlock',
  addKeyboardPwd = '/houses/ttlock/addKeyboardPwd',
  sendEkey = '/houses/ttlock/sendEkey',
  lockDetail = '/houses/ttlock/lockDetail',
}

export function getStatus() {
  return defHttp.get({ url: Api.status }, { errorMessageMode: 'none' });
}

export function syncLocks() {
  return defHttp.post({ url: Api.syncLocks, params: {} }, { errorMessageMode: 'message' });
}

export function getLockList(params: Record<string, any>) {
  return defHttp.get({ url: Api.lockList, params }, { errorMessageMode: 'none' });
}

export function getLockRecords(params: { property_id?: number; ttlock_lock_id?: number; start_date?: number; end_date?: number; page?: number; pageSize?: number }) {
  return defHttp.get({ url: Api.lockRecords, params }, { errorMessageMode: 'none' });
}

export function getPasscodeEvents(params: { property_id?: number; ttlock_lock_id?: number; page?: number; pageSize?: number }) {
  return defHttp.get({ url: Api.passcodeEvents, params }, { errorMessageMode: 'none' });
}

export function bindProperty(params: { property_id: number; ttlock_lock_id: number }) {
  return defHttp.post({ url: Api.bindProperty, params }, { errorMessageMode: 'message' });
}

export function unbindProperty(params: { property_id: number }) {
  return defHttp.post({ url: Api.unbindProperty, params }, { errorMessageMode: 'message' });
}

export function getPropertyLock(params: { property_id: number }) {
  return defHttp.get({ url: Api.propertyLock, params }, { errorMessageMode: 'none' });
}

export function remoteUnlock(params: { property_id?: number; ttlock_lock_id?: number }) {
  return defHttp.post({ url: Api.remoteUnlock, params }, { errorMessageMode: 'message' });
}

export function addKeyboardPwd(params: Record<string, any>) {
  return defHttp.post({ url: Api.addKeyboardPwd, params }, { errorMessageMode: 'message' });
}

export function sendEkey(params: Record<string, any>) {
  return defHttp.post({ url: Api.sendEkey, params }, { errorMessageMode: 'message' });
}

export function getLockDetail(params: { property_id?: number; ttlock_lock_id?: number }) {
  return defHttp.get({ url: Api.lockDetail, params }, { errorMessageMode: 'none' });
}
