<template>
  <a-card class="general-card contentcard" v-if="formData">
       <template #title> 应用(系统)主要配置(对应resource/config.yaml) 是整个项目共用配置</template>
       <a-row :gutter="80">
         <a-col :span="24">
             <a-form-item
             label="前端代码开发根目录位置"
             field="vueobjroot"
             extra="配置代码生成时-前端代码根目录位置(开发必改)"
             >
             <a-input v-model="formData.vueobjroot" placeholder="填写附件请求路径前缀" allow-clear/>
             </a-form-item>
         </a-col>
          <a-col :span="12">
            <a-form-item
              label="登录启用人机验证"
              field="loginCaptcha"
              tooltip="登录是否启用人机验证，启用后登录时需要输入正确的数字验证码才可继续登录（作用域两个端登录后台）"
              >
                 <a-switch type="round" v-model="formData.loginCaptcha" :checked-value="true" :unchecked-value="false">
                    <template #checked>{{ $t('cell.open') }}</template>
                    <template #unchecked>{{ $t('cell.close') }}</template>
                  </a-switch>
              </a-form-item>
          </a-col>
          <a-col :span="12">
              <a-form-item
              label="启用请求对称加密"
              field="validityApi"
              tooltip="是否启用接口请求对称加密，启用后请求后端接口需要检验请求头的对称密文(作用域为全局)"
              >
                 <a-switch type="round" v-model="formData.validityApi" :checked-value="true" :unchecked-value="false">
                    <template #checked>{{ $t('cell.open') }}</template>
                    <template #unchecked>{{ $t('cell.close') }}</template>
                  </a-switch>
              </a-form-item>
         </a-col>
         <a-col :span="24">
             <a-form-item >
             <div class="frombtn">
                 <a-button type="primary" html-type="submit" style="width: 120px;" @click="submitAttachmentConfig">保存</a-button>
             </div>
             </a-form-item>
         </a-col>
       </a-row>
   </a-card>
 </template>
 
 <script lang="ts" setup>
   import {  ref, onMounted } from 'vue';
   //api
   import { saveConfig,getConfig} from '@/api/datacenter/common_config';
   import { Message } from '@arco-design/web-vue';
   //数据配置
   const formData=ref({
    vueobjroot:"",
    loginCaptcha:false,
    validityApi:false,
   })
   //保存邮箱
   const submitAttachmentConfig=async()=>{
     try{
       Message.loading({content:"保存中",id:"updata",duration:0})
       await saveConfig(formData.value);
       Message.success({content:"保存成功",id:"updata",duration:2000})
     } catch (error) {
       Message.loading({content:"保存中",id:"updata",duration:1})
     }
   }
   //组件挂载完成后执行的函数
   onMounted(()=>{
     InitData()
   })
   //加载数据
   const InitData=async()=>{
     const appdata = await getConfig({});
     formData.value=Object.assign({},formData.value,appdata)
   }
 </script>
 
 <style scoped lang="less">
 .contentcard{
       overflow: hidden;
   }
   :deep(.general-card > .arco-card-header){
     padding: 10px 16px;
   }
   .iconbtn{
     user-select: none;
     cursor: pointer;
     opacity: .8;
     &:hover{
       opacity: 1;
     }
   }
   .frombtn{
     width: 100%;
     text-align: center;
   }
 </style>
 