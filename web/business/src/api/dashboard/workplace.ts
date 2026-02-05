import { defHttp } from '@/utils/http';
enum Api {
    getStatistical = '/dashboard/workplace/get_statistical',
    getVisitlist = '/dashboard/workplace/get_visitlist',
    getPopular = '/dashboard/workplace/get_popular',
    saveQuick = '/dashboard/workplace/saveQuick',
    getQuick = '/dashboard/workplace/getQuick',
    delQuick = '/dashboard/workplace/delQuick',
    getSiteCount = '/dashboard/workplace/get_siteCount',
    getMessage = '/dashboard/workplace/get_message',
    getMsmContent = '/dashboard/workplace/get_msmContent',
}

// 首页概况统计
export interface WorkplaceStatistical {
  propertyTotal: number;
  propertyOnSale: number;
  lockBindTotal: number;
  unlockPendingTotal: number;
  todayPropertyAdd: number;
  todayUnlockRequests: number;
}
export function getStatistical(params: object) {
  return defHttp.get<WorkplaceStatistical>({ url: Api.getStatistical, params }, { errorMessageMode: 'message' });
}

// 趋势数据类型
export interface ContentDataRecord {
  x: string;
  y: number;
}
export function getVisitlist(params: object) {
  return defHttp.get<ContentDataRecord[]>({ url: Api.getVisitlist, params }, { errorMessageMode: 'message' });
}

// 热门房源
export interface PopularPropertyRecord {
  id: number;
  title: string;
  viewCount: number;
  followCount: number;
  showingCount: number;
  saleStatus: 'on_sale' | 'sold' | 'off_market' | string;
  price: number;
  priceUnit: string;
  area: number;
}
export function getPopular(params: object) {
  return defHttp.get<PopularPropertyRecord[]>({ url: Api.getPopular, params }, { errorMessageMode: 'message' });
}

//提交快捷操作数据
export function saveQuick(params: object) {
    return defHttp.post({ url: Api.saveQuick, params:params}, { errorMessageMode: 'message' });
}
//获取快捷操作
export function getQuick(params: object) {
    return defHttp.get<QuickItem []>({ url: Api.getQuick, params:params}, { errorMessageMode: 'message' });
}
//删除快捷操作
export function delQuick(params: object) {
    return defHttp.delete({ url: Api.delQuick, params:params}, { errorMessageMode: 'message' });
}
//获取网站常用统计数量
export function getSiteCount(params: object) {
    return defHttp.get({ url: Api.getSiteCount, params:params}, { errorMessageMode: 'message' });
}
//获取获取公告信息
export function getMessage(params: object) {
    return defHttp.get({ url: Api.getMessage, params:params}, { errorMessageMode: 'message' });
}
export function getMsmContent(params: object) {
    return defHttp.get({ url: Api.getMsmContent, params:params}, { errorMessageMode: 'message' });
}
/**数据类型 */
export interface DataItem {
    id:number,
    name: string;
}

//公告类型
export interface MessageItem {
    id:number;
    type:number;
    title:string;
    path:string;
    content:string;
    isread:number;
    createtime:number;
}
//快捷类型
export interface QuickItem {
    id:number;
    name:string;
    icon: string;
    path_url:string;
    is_common:number;
    type:number;
}
