import { defHttp } from '@/utils/http';
import { join } from 'lodash';
import { isArray } from '/@/utils/is';
enum Api {
    getList = '/area/getList',
    getMoreList = '/area/getMoreList',
    getContent= '/area/getContent',
    save = '/area/save',
    upStatus = '/area/upStatus',
    del = '/area/del',
}

//列表数据
export function getList(params: object) {
    for(let key in params){
        if(isArray(params[key])){
            params[key]=join(params[key])
        }
    }
  return defHttp.get({ url: Api.getList, params:params }, { errorMessageMode: 'none' });
}
//获取子级数据
export function getMoreList(params: object) {
    return defHttp.get({ url: Api.getMoreList, params:params }, { errorMessageMode: 'none' });
}
//获取内容
export function getContent(params: object) {
    return defHttp.get({ url: Api.getContent, params:params }, { errorMessageMode: 'none' });
}
//提交数据
export function save(params: object) {
    return defHttp.post({ url: Api.save, params:params}, { errorMessageMode: 'message' });
}
//更新状态
export function upStatus(params: object) {
    return defHttp.post({ url: Api.upStatus, params:params}, { errorMessageMode: 'message' });
}
//删除数据
export function del(params: object) {
    return defHttp.delete({ url: Api.del, params:params}, { errorMessageMode: 'message' });
}
/**数据类型 */
export interface DataItem {
    id:number,
    name: string;
}
