import { defHttp } from '@/utils/http';
import { join } from 'lodash';
import { isArray } from '/@/utils/is';
enum Api {
    // GoFly 自动路由：/business/{pkg}/{controller}/{action}
    getList = '/houses/broker/getList',
    getContent = '/houses/broker/getContent',
    save = '/houses/broker/save',
    upStatus = '/houses/broker/upStatus',
    auditSwitch = '/houses/broker/auditSwitch',
    upCanManageProperties = '/houses/broker/upCanManageProperties',
    upCanManageLocks = '/houses/broker/upCanManageLocks',
    del = '/houses/broker/del',
}

//列表数据
export function getList(params: any) {
    for(let key in params){
        if(isArray(params[key])){
            params[key]=join(params[key])
        }
    }
  return defHttp.get({ url: Api.getList, params:params }, { errorMessageMode: 'none' });
}

// 获取内容
export function getContent(params: any) {
  return defHttp.get({ url: Api.getContent, params }, { errorMessageMode: 'none' });
}

//提交数据
export function save(params: object) {
    return defHttp.post({ url: Api.save, params:params}, { errorMessageMode: 'message' });
}
//更新状态
export function upStatus(params: object) {
    return defHttp.post({ url: Api.upStatus, params:params}, { errorMessageMode: 'message' });
}
// 审核开关（pending <-> approved）
export function auditSwitch(params: object) {
    return defHttp.post({ url: Api.auditSwitch, params:params}, { errorMessageMode: 'message' });
}
// 切换“可维护房源”权限
export function upCanManageProperties(params: object) {
    return defHttp.post({ url: Api.upCanManageProperties, params:params}, { errorMessageMode: 'message' });
}
// 切换“可管理智能锁”权限
export function upCanManageLocks(params: object) {
    return defHttp.post({ url: Api.upCanManageLocks, params:params}, { errorMessageMode: 'message' });
}
//删除数据
export function del(params: object) {
    return defHttp.delete({ url: Api.del, params:params}, { errorMessageMode: 'message' });
}
/**数据类型 */
export interface DataItem {
    id:number,
    username?: string;
    name: string;
    nickname?: string;
    mobile?: string;
    email?: string;
    role?: string;
    can_manage_properties?: number;
    can_manage_locks?: number;
    store_id?: number;
    title?: string;
    status?: number;
}
