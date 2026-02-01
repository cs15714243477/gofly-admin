import { defHttp } from '@/utils/http';
import { UploadFileParams } from '/#/axios';
import { FilterParam } from '/@/utils';
enum Api {
    getList = '/matter/picture/getList',
    save = '/matter/picture/save',
    upStatus = '/matter/picture/upStatus',
    del = '/matter/picture/del',
}

//列表数据
export function getList(params: any) {
  params=FilterParam(params)
  return defHttp.get({ url: Api.getList, params:params }, { errorMessageMode: 'message' });
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

//上传文件
export interface UploadItem {
    id: number;
    name: string;
    response: string;
    status: string;
    time: number;
    uid: string;
    url: string;
    code: number;
    data: any;
  }
//获取配置中的上传文件路径
const DOMAIN =import.meta.env.VITE_APP_ENV=="production"? window?.globalConfig.Root_url: window?.globalConfig.Root_url_dev;
export function userUploadApi(
  params: UploadFileParams,
  onUploadProgress?: (progressEvent: any) => void
) {
  return defHttp.uploadFile<UploadItem>({ url:`${DOMAIN}/matter/picture/uploadFile`,onUploadProgress},params);
}