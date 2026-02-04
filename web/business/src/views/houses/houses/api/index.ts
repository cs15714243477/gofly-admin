import { defHttp } from '@/utils/http';
import { join } from 'lodash';
import { isArray } from '/@/utils/is';

enum Api {
  // GoFly 自动路由：/business/{pkg}/{controller}/{action}
  getList = '/houses/houses/getList',
  getContent = '/houses/houses/getContent',
  save = '/houses/houses/save',
  upStatus = '/houses/houses/upStatus',
  del = '/houses/houses/del',
  getStatusLogs = '/houses/houses/getStatusLogs',
  // 表单下拉选项（统一配置/字典兜底）
  getFormOptions = '/uniapp/wxproperty/getFormOptions',
  // 装修信息
  getRenovation = '/houses/houses/getRenovation',
  saveRenovation = '/houses/houses/saveRenovation',
  // 地区
  getAreaMoreList = '/area/getMoreList',
}

// 列表数据
export function getList(params: Record<string, any>) {
  for (const key in params) {
    if (isArray(params[key])) {
      params[key] = join(params[key]);
    }
  }
  return defHttp.get({ url: Api.getList, params }, { errorMessageMode: 'none' });
}

// 获取内容
export function getContent(params: Record<string, any>) {
  return defHttp.get({ url: Api.getContent, params }, { errorMessageMode: 'none' });
}

// 提交数据
export function save(params: Record<string, any>) {
  return defHttp.post({ url: Api.save, params }, { errorMessageMode: 'message' });
}

// 更新状态/销售状态/排序
export function upStatus(params: Record<string, any>) {
  return defHttp.post({ url: Api.upStatus, params }, { errorMessageMode: 'message' });
}

// 获取状态变更记录
export function getStatusLogs(params: Record<string, any>) {
  return defHttp.get({ url: Api.getStatusLogs, params }, { errorMessageMode: 'none' });
}

// 删除数据
export function del(params: Record<string, any>) {
  return defHttp.delete({ url: Api.del, params }, { errorMessageMode: 'message' });
}

// 获取装修信息
export function getRenovation(params: Record<string, any>) {
  return defHttp.get({ url: Api.getRenovation, params }, { errorMessageMode: 'none' });
}

// 保存装修信息
export function saveRenovation(params: Record<string, any>) {
  return defHttp.post({ url: Api.saveRenovation, params }, { errorMessageMode: 'message' });
}

// 获取房源表单下拉选项（后台配置/字典）
export function getFormOptions(params: Record<string, any>) {
  return defHttp.get({ url: Api.getFormOptions, params }, { errorMessageMode: 'none' });
}

// 获取地区子级（pid=0 为省）
export function getAreaMoreList(params: Record<string, any>) {
  return defHttp.get({ url: Api.getAreaMoreList, params }, { errorMessageMode: 'none' });
}
