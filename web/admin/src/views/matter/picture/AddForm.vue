<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :isPadding="false" :loading="loading" width="800px" @height-change="onHeightChange" :minHeight="modelHeight" :title="getTitle" @ok="handleSubmit">
    <div class="addFormbox" :style="{'min-height':`${windHeight}px`}">
      <div class="tabs-content addpadding"  >
        <a-form ref="formRef" :model="formData" auto-label-width>
          <div class="content_box">
              <!--基础信息-->
              <a-scrollbar v-show="activeKey==1" style="overflow: auto;" :style="{height:`${windHeight}px`}">
                <div class="besecontent" >
                  <a-row :gutter="16">
                    <a-col :span="12">
                      <a-form-item field="cid" label="图片类型" validate-trigger="input" >
                        <a-radio-group v-model="formData.type" :options="OYoptions" />
                      </a-form-item>
                    </a-col>
                    <a-col :span="16">
                      <a-form-item field="cid" label="选择分类" :rules="[{required:true,message:'请填写分类'}]">
                        <a-select v-model="formData.cid"   placeholder="请选择分类" allow-search >
                          <template v-for="item in cateList">
                            <a-option v-if="item.type==formData.type||item.type==2" :value="item.id">{{item.name}}</a-option>
                          </template>
                          </a-select>
                      </a-form-item>
                    </a-col>
                    <a-col :span="24">
											<a-form-item field="image" label="图片" style="margin-bottom:15px;">
											<a-upload
												accept="image/*"
												:show-file-list="false"
												:custom-request="customRequest"
												>
												<template #upload-button>
													<div class="upimagebox">
													<div class="imagebtn">
														<div
															class="upload-show-picture"
															v-if="formData.url"
														>
															<img :src="GetFullPath(formData.url)" />
															<div class="upload-show-picture-mask">
															<IconEdit />
															</div>
														</div>
														<div class="upload-picture-card" v-else>
															<div class="upload-picture-card-text">
															<IconPlus />
															<div style="margin-top: 10px; font-weight: 600">上传图片</div>
															</div>
														</div>
													</div>
													</div>
												</template>
												</a-upload>
											</a-form-item>
										</a-col>
										<a-col :span="16">
											<a-form-item field="title" label="图片名称" :rules="[{required:true,message:'请填写图片名称'}]" >
                        <a-input v-model="formData.title" placeholder="请填图片名称" />
											</a-form-item>
										</a-col>
                    </a-row>
                </div>
              </a-scrollbar>
          </div>
        </a-form>
      </div>
    </div>
  </BasicModal>
</template>
<script lang="ts">
  import { defineComponent, ref, computed, unref} from 'vue';
  import { BasicModal, useModalInner } from '/@/components/Modal';
  import useLoading from '@/hooks/loading';
  import { cloneDeep } from 'lodash-es';
  //api
  import { save,userUploadApi} from '@/api/matter/picture';
  import { getCate,DataItem} from '@/api/matter/picturecate';
  import { FormInstance,Message } from '@arco-design/web-vue';
  import type { RequestOption} from '@arco-design/web-vue/es/upload/interfaces';
  import { GetFullPath } from "@/utils/tool";
  export default defineComponent({
    name: 'AddBook',
    components: { BasicModal },
    emits: ['success'],
    setup(_, { emit }) {
      //判断是否存在编辑器
      const isUpdate = ref(false);
      const cateList = ref<DataItem[]>([]);
      const activeKey= ref(1);
      const modelHeight= ref(450);
      const windHeight= ref(450);
      //表单
      const formRef = ref<FormInstance>();
      //表单字段
      const basedata={
            id:0,
            type:0,
            cid:0,
            url:"",
            title:"",
        }
      const formData = ref(cloneDeep(basedata))
      //编辑器
      const editorRef = ref();
      const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
          formRef.value?.resetFields()
          activeKey.value=1
          setModalProps({ confirmLoading: false });
          const mdata = await getCate({});
          const parntList_df : any=[{id: 0,name: "未选分类",type:2}];
          if(mdata){
            cateList.value=parntList_df.concat(mdata)
          }else{
            cateList.value=parntList_df
          }
          isUpdate.value = !!data?.isUpdate;
          if (unref(isUpdate)) {
            formData.value=cloneDeep(data.record)
          }else{
            formData.value=cloneDeep(basedata)
          }
      });
      const getTitle = computed(() => (!unref(isUpdate) ? '新增数据' : '编辑数据'));
     //点击确认
     const { loading, setLoading } = useLoading();
     const handleSubmit = async () => {
      try {
          const res = await formRef.value?.validate();
          if (!res) {
            setLoading(true);
            Message.loading({content:"提交中",id:"save",duration:0})
            let savedata=cloneDeep(unref(formData))
            await save(savedata);
            Message.success({content:"提交成功",id:"save",duration:2000})
            closeModal()
            emit('success');
            setLoading(false);
          }
        } catch (error) {
          setLoading(false);
          Message.loading({content:"提交中",id:"save",duration:1})
        }
      };
      //上传附件=图片
      const customRequest = (options: RequestOption) => {
          Message.loading({content:"上传中",id:"upfile",duration:0})
          const controller = new AbortController();
            (async function requestWrap() {
              const {
                onProgress,
                onError,
                onSuccess,
                fileItem,
              } = options;
              onProgress(20);
              const onUploadProgress = (event: ProgressEvent) => {
                let percent;
                if (event.total > 0) {
                  percent = (event.loaded / event.total) * 100;
                }
                onProgress(parseInt(String(percent), 10), event);
              };
              try {
                //开始手动上传
                const filename=fileItem?.name||""
                const resdata = await userUploadApi({ name: 'file', file: fileItem.file as Blob, filename,data:{cid:0,id:formData.value.id}},onUploadProgress);
                //更新图片
                if(resdata["code"]==0){
                  formData.value.url=resdata["data"].url
                  formData.value.id=resdata["data"].id
                  formData.value.title=resdata["data"].title
                }
                Message.success({content:"上传成功",id:"upfile",duration:1500})
              } catch (error) {
                onError(error);
                Message.loading({content:"上传中",id:"upfile",duration:0})
              }
            })();
            return {
              abort() {
                controller.abort();
              },
            };
      };
       //监听高度
       const onHeightChange=(val:any)=>{
        windHeight.value=val
      }
      return { 
        registerModal, 
        getTitle, 
        handleSubmit,
        formRef,
        loading,
        formData,
        cateList,
        OYoptions:[
          { label: '素材图', value: 0 },
          { label: '插图', value: 1 },
          ],
        activeKey,
        customRequest,
        modelHeight,
        editorRef,
        onHeightChange,windHeight,
        GetFullPath,
      };
    },
  });
</script>
<style lang="less" scoped>
  @import '@/assets/style/formlayer.less';
  .upfilezip{
    display: flex;
    .upbtn{
      padding-right: 10px;
    }
    .showfile{
      flex: 1;
      height: 32px;
      line-height: 32px;
      a{
        text-decoration: none;
      }
    }
  }
  //上传图片
  .upimagebox{
    display: flex;
    .imagebtn{
      position: relative;
      width: 160px;
      height: 90px;
      background-color: var(--color-neutral-1);
      border-radius: 4px;
      overflow: hidden;
      -ms-flex-negative: 0;
      flex-shrink: 0;
      //预览
      .upload-show-picture{
        position: relative;
        box-sizing: border-box;
        width: 100%;
        height: 100%;
        overflow: hidden;
        display: flex;
        align-items: center;
        justify-content: center;
        img{
          height: 100%;
        }
        &:hover{
          .upload-show-picture-mask{
             opacity: 1;
          }
        }
        .upload-show-picture-mask{
            position: absolute;
            top: 0;
            right: 0;
            bottom: 0;
            left: 0;
            color: var(--color-white);
            font-size: 16px;
            line-height: 90px;
            text-align: center;
            background: rgba(0, 0, 0, 0.5);
            cursor: pointer;
            opacity: 0;
            transition: opacity 0.1s cubic-bezier(0, 0, 1, 1);
        }
      }
      .upload-picture-card{
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        .upload-picture-card-text{
           text-align: center;
           color:  var(--color-neutral-5);
        }
      }
    }
  }
</style>

