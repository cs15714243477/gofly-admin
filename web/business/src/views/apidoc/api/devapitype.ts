import { defHttp } from '@/utils/http';
//类型
export interface LoginData {
  username: string;
  password: string;
}
enum Api {
    getList = '/apidoc/devapitype/get_list',
    getTypeinfo = '/apidoc/devapitype/get_typeinfo',
    save = '/apidoc/devapitype/save',
    del = '/apidoc/devapitype/del',
}

//数据列表
export function getType(params: object) {
  return defHttp.get({ url: Api.getList, params:params }, { errorMessageMode: 'none' });
}
//获取单条数据
export function getTypeinfo(params: object) {
  return defHttp.get({ url: Api.getTypeinfo, params:params }, { errorMessageMode: 'none' });
}
//提交
export function saveType(params: any) {
  return defHttp.post({ url: Api.save, params:params}, { errorMessageMode: 'message' });
}
//删除数据
export function delType(params: object) {
    return defHttp.delete({ url: Api.del, params:params}, { errorMessageMode: 'message' });
}