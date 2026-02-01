import { useUserStore } from '@/store'
import { isAllUrl } from '@/utils/is';
import { useClipboard } from '@vueuse/core';
import { Message } from '@arco-design/web-vue';
//获取附件完整访问路径
export function GetFullPath(url:any):string{
    if(!url){
        return ""
    }else if(isAllUrl(url)){
        return url
    }else{
        var imgname = url.split('/');//分割url
        imgname = imgname[imgname.length-1];//通过最后一个数组下标获取图片名
        const userStore = useUserStore()
        if(imgname.startsWith("local")){
           url= userStore.rooturls?.local+url
        }else if(imgname.startsWith("alioss")){
           url= userStore.rooturls?.alioss+url
        }else if(imgname.startsWith("tencentcos")){
           url= userStore.rooturls?.tencentcos+url
        }else if(imgname.startsWith("qiniuoss")){
           url= userStore.rooturls?.qiniuoss+url
        }else{//如果地址前缀无法判断文件存储位置，则使用设置上传方式的访问地址
          url=userStore.defrooturl+url
        }
        return url
    }
}
//复制
var copyObj :any
export function CopyText(text:any,msg:string){
    if(!copyObj){
        const { copy } = useClipboard();
        copyObj=copy
    }
    copyObj(text);
    Message.success({content:msg,id:"copy"});
}
