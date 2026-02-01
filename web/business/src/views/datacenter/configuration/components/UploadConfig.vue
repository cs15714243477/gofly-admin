<template>
   <a-card class="general-card contentcard" v-if="formData">
        <template #title> 
          文件上传配置
         <span class="tig" >配置文件位置：resource/config/upload.yaml</span>  
        </template>
        <a-form ref="formRef" :model="formData" auto-label-width>
          <a-row :gutter="80">
            <a-col :span="24">
              <a-form-item
                label="可上传文件类型"
                field="AllowedExt"
                tooltip="允许上传的文件类型，使用英文逗号(,)隔开，文件类型后缀前面带点，如：.jpg,.jpeg,.png,.pdf,.ico"
              >
                <a-input v-model="formData.AllowedExt" placeholder="填写允许上传文件类型" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item
                label="传输文件最大值"
                field="MaxBodySize"
                tooltip="限制上传文件的大小，超过大小的文件会被禁止提交"
              >
                <a-input v-model="formData.MaxBodySize" placeholder="填写允许上传文件的大小" allow-clear append="MB" style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24">
              <a-form-item
                label="文件存储方式"
                field="Type"
                tooltip="文件存储方式，目前支持：阿里云、腾讯云、七牛云、本地4种方式的静态云存储，如果项目需要其他云存储（华为云）可以参考当前代码扩展"
              >
                 <a-radio-group type="button" v-model="formData.Type" @change="handleChangeType">
                  <a-radio value="local">本地</a-radio>
                  <a-radio value="tencentcos">腾讯云</a-radio>
                  <a-radio value="alioss">阿里云</a-radio>
                  <a-radio value="qiniuoss">七牛云</a-radio>
                  <a-tooltip content="如果业务需要你可继续扩展" background-color="#ff9a2e">
                    <a-radio value="huaweiobs" disabled>华为云</a-radio>
                  </a-tooltip>
                </a-radio-group>
              </a-form-item>
            </a-col>
            <a-col :span="24" >
              <a-form-item
                label="文件访问路径"
                field="BaseUrl"
                tooltip="访问上传服务器的文件，访问的路径。"
              >
                <a-input v-model="formData.BaseUrl" placeholder="填写文件访问路径" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>

            <a-col :span="24" v-if="formData.Type!='local'">
              <a-form-item
                label="自定义访问域名"
                field="Endpoint"
                tooltip="自定义服务请求的访问域名"
              >
                <a-input v-model="formData.Endpoint" placeholder="填写自定义服务访问域名" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type!='local'">
              <a-form-item
                label="密钥Key或秘钥id"
                field="KeyId"
                tooltip="密钥Key或秘钥id对应AccessKey或者是appID编号"
              >
                <a-input v-model="formData.KeyId" placeholder="填写密钥Key或秘钥id" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type!='local'">
              <a-form-item
                label="密钥SecretKey"
                field="Secret"
                tooltip="对应云存储给的SecretKey,是AccessKey ID的密码"
              >
                <a-input-password v-model="formData.Secret" placeholder="填写密钥SecretKey" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type=='alioss'||formData.Type=='qiniuoss'||formData.Type=='tencentcos'">
              <a-form-item
                label="空间名称"
                field="BucketName"
                tooltip="对象存储OSS中用于存储文件（Object）的基础容器"
              >
                <a-input v-model="formData.BucketName" placeholder="填写空间名称" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type=='tencentcos'">
              <a-form-item
                label="所属地域"
                field="Region"
                tooltip="指存储桶的所属地域，选择与您最近的一个地区。例如，您在深圳，地域可以选择广州，即 ap-guangzhou"
              >
                <a-input v-model="formData.Region" placeholder="填写存储所属地域" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
              <a-col :span="24" >
              <a-form-item
                label="上传文件路径"
                field="DirPath"
                tooltip="上传到服务器目录路径，如本地：/resource/uploads/"
              >
                <a-input v-model="formData.DirPath" placeholder="填写上传附件路径" allow-clear style="max-width: 560px;"/>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type=='qiniuoss'">
              <a-form-item
                label="是否使用https"
                field="UseHTTPS"
                tooltip="是否使用https域名进行资源管理"
              >
                <a-switch type="round" v-model="formData.UseHTTPS" :checked-value="true" :unchecked-value="false">
                  <template #checked>{{ $t('cell.open') }}</template>
                  <template #unchecked>{{ $t('cell.close') }}</template>
                </a-switch>
              </a-form-item>
            </a-col>
            <a-col :span="24" v-if="formData.Type=='qiniuoss'">
              <a-form-item
                label="存储区域"
                field="Zone"
                tooltip="文件存储的区域"
              >
               <a-select v-model="formData.Zone" placeholder="选择存储区域" style="max-width: 560px;">
                  <a-option :value="1">华东机房</a-option>
                  <a-option :value="2">华东-浙江(2区)</a-option>
                  <a-option :value="3">华北-河北</a-option>
                  <a-option :value="4">华南-广东</a-option>
                  <a-option :value="5">北美机房</a-option>
                  <a-option :value="6">新加坡机房</a-option>
                </a-select>
              </a-form-item>
            </a-col>

            <a-col :span="24">
                <a-form-item style="max-width: 560px;">
                  <div class="frombtn">
                      <a-button type="primary" html-type="submit" style="width: 120px;" @click="submitAttachmentConfig">保存</a-button>
                  </div>
                </a-form-item>
            </a-col>
          </a-row>
        </a-form>
    </a-card>
  </template>
  
  <script lang="ts" setup>
    import { ref, onMounted } from 'vue';
    //api
    import { saveConfig,getConfig} from '@/api/datacenter/uploadconfig';
    import { Message } from '@arco-design/web-vue';
    import { cloneDeep } from 'lodash-es';
    //数据配置
    const formData=ref({
      Type:"local",
      MaxBodySize:600,
      AllowedExt:".jpg,.jpeg,.png,.pdf",
      BaseUrl:"",
      DirPath:"/resource/uploads/",
      Endpoint:"",
      KeyId:"",
      Secret:"",
      BucketName:"",
      DestBucketName:"",
      Region:"",
      UseHTTPS:false,
      Zone:4,

    })
    //保存邮箱
    const submitAttachmentConfig=async()=>{
      try{
        Message.loading({content:"保存中",id:"updata",duration:0})
        await saveConfig(cloneDeep(formData.value));
        InitData()
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
    var baseData:any
    const InitData=async()=>{
      const emaildata = await getConfig({});
      baseData=emaildata
      formData.value=Object.assign({},formData.value,emaildata)
      handleChangeType(formData.value.Type)
    }
    //切换上传方式
    const handleChangeType=(value:any)=>{
      if(value=="local"){
        formData.value=Object.assign({},formData.value,{BaseUrl:baseData.local.LBaseUrl,DirPath:baseData.local.LDirPath})
      }else if(value=="tencentcos"){
        formData.value=Object.assign({},formData.value,{BaseUrl:baseData.tencentcos.TBaseUrl,Endpoint:baseData.tencentcos.TEndpoint,
          KeyId:baseData.tencentcos.TKeyId,Secret:baseData.tencentcos.TSecret,BucketName:baseData.tencentcos.TBucketName,
          Region:baseData.tencentcos.TRegion,DirPath:baseData.tencentcos.TDirPath})
      }else if(value=="alioss"){
        formData.value=Object.assign({},formData.value,{BaseUrl:baseData.alioss.ABaseUrl,Endpoint:baseData.alioss.AEndpoint,
          KeyId:baseData.alioss.AKeyId,Secret:baseData.alioss.ASecret,BucketName:baseData.alioss.ABucketName,
          DirPath:baseData.alioss.ADirPath})
      }else if(value=="qiniuoss"){
        formData.value=Object.assign({},formData.value,{BaseUrl:baseData.qiniuoss.QBaseUrl,Endpoint:baseData.qiniuoss.QEndpoint,
          KeyId:baseData.qiniuoss.QKeyId,Secret:baseData.qiniuoss.QSecret,BucketName:baseData.qiniuoss.QBucketName,
        DestBucketName:baseData.qiniuoss.QDestBucketName,UseHTTPS:baseData.qiniuoss.QUseHTTPS,Zone:baseData.qiniuoss.QZone,
        DirPath:baseData.qiniuoss.QDirPath})
      }
      console.log(formData.value)
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
    .tig{
      font-size: 12px;
      color: var(--color-neutral-4);
      padding-left: 5px;
    }
  </style>
  