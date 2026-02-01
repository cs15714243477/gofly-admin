import { defHttp } from '@/utils/http';
enum Api {
    getConfig = '/datacenter/uploadconfig/getConfig',
    saveConfig = '/datacenter/uploadconfig/saveConfig',
}

//列表数据
export function getConfig(params: object) {
  return defHttp.get({ url: Api.getConfig, params:params }, { errorMessageMode: 'message' });
}
//提交数据
export function saveConfig(params: any) {
    return defHttp.post({ url: Api.saveConfig, params:params}, { errorMessageMode: 'message' });
}
