<template>
   <a-card class="general-card contentcard" v-if="ConfData">
        <template #title>
         <div class="flex">
          <div style="flex:1;"> 
            {{ ConfData?.title }}
            <a-tag v-if="ConfData?.pluginident" >
              标识： {{ConfData?.pluginident}}
            </a-tag>
            <span class="tig" >配置文件位置：resource/config/{{ConfData?.name}}.yaml</span>
            <a-tooltip content="配置字段命名说明：_content结尾使用富文本编辑器，_txt结尾是多行文本，_read结尾是只读，_switch结尾开关,#注释加 &des 是字段说明。" position="top" mini>
              <icon-question-circle-fill  class="configdes"/>
            </a-tooltip>
          </div>
          <a-tooltip content="如果有相同配置标识文件，当点击启用时会禁用其他相同标识文件" position="tr" mini>
              <a-switch v-model="ConfData.status" @change="handleUpConfigStatus">
                <template #checked>
                  启用
                </template>
                <template #unchecked>
                  禁用
                </template>
              </a-switch>
          </a-tooltip>  
         </div>
        </template>
        <a-form :model="form"  auto-label-width>
          <a-row :gutter="20">
            <template v-for="item in ConfData.data" :key="item.keyfield">
              <a-col :span="12">
                  <a-form-item
                  :class="{ 'rich-editor-form-item': isRichTextField(item) }"
                  :label="FontNameAndDes(item.keyname,0)"
                  :field="item.keyfield"
                  :extra="FontNameAndDes(item.keyname,1)"
                  >
                    <div v-if="isRichTextField(item)" class="rich-editor-trigger">
                      <a-button type="outline" @click="openRichEditor(item)">查看/编辑文档</a-button>
                      <a-tag color="arcoblue" size="small">{{ getRichContentState(item) }}</a-tag>
                    </div>
                    <a-textarea v-model="item.keyvalue" :placeholder="`填写${FontNameAndDes(item.keyname,0)}`" v-else-if="item.keyfield.indexOf('_txt')>-1"/>
                    <a-tooltip :content="`“${FontNameAndDes(item.keyname,0)}”是只读不可编辑字段`" position="top" mini  v-else-if="item.keyfield.indexOf('_read')>-1">
                      <a-input v-model="item.keyvalue" :placeholder="`填写${FontNameAndDes(item.keyname,0)}`" readonly />
                    </a-tooltip>
                    <a-switch v-model="item.keyvalue" :checked-value="1" :unchecked-value="0" checked-text="是" unchecked-text="否" v-else-if="item.keyfield.indexOf('_switch')>-1"/>
                    <a-input v-model="item.keyvalue" :placeholder="`填写${FontNameAndDes(item.keyname,0)}`" allow-clear v-else/>
                  </a-form-item>
              </a-col>
            </template>
            <a-col :span="24">
                <a-form-item >
                  <div class="frombtn">
                      <a-button type="primary" html-type="submit" style="width: 120px;" @click="submitConfig">保存</a-button>
                  </div>
                </a-form-item>
            </a-col>
          </a-row>
        </a-form>
    </a-card>
    <a-drawer
      class="rich-editor-drawer"
      :visible="richEditorVisible"
      :title="richEditorTitle"
      :width="980"
      unmount-on-close
      @cancel="closeRichEditor"
    >
      <div v-if="activeRichItem" class="rich-editor-wrapper">
        <umo-editor
          v-bind="getRichEditorOptions(activeRichItem)"
          :theme="appStore.theme"
          @changed="(payload:any)=>handleRichEditorChanged(activeRichItem,payload)"
        />
      </div>
      <template #footer>
        <a-space>
          <a-button @click="closeRichEditor">关闭</a-button>
        </a-space>
      </template>
    </a-drawer>
  </template>
  
  <script lang="ts" setup>
    import { reactive, ref,PropType,watch } from 'vue';
    import { UmoEditor } from '@umoteam/editor';
    import '@umoteam/editor/style';
    //api
    import { DataObj, menuItem,saveCodeStoreConfig,upConfigStatus} from '@/api/datacenter/configuration';
    import { useAppStore } from '@/store';
    import { Message } from '@arco-design/web-vue';
    const appStore = useAppStore();
    const emits = defineEmits(['ok'])
    const props = defineProps({
      ConfData: {
        type: Object as PropType<menuItem>,
      },
    })
    const editorOptionsMap=reactive<Record<string,any>>({})
    const richEditorVisible = ref(false)
    const activeRichItem = ref<DataObj | null>(null)
    const richEditorTitle = ref('文档配置')
    watch(()=>props.ConfData?.name,()=>{
      Object.keys(editorOptionsMap).forEach((key)=>delete editorOptionsMap[key])
      richEditorVisible.value = false
      activeRichItem.value = null
      richEditorTitle.value = '文档配置'
    },{immediate:true})
    const isRichTextField=(item:DataObj)=>{
      return !!item?.keyfield&&item.keyfield.indexOf('_content')>-1
    }
    const getRichEditorOptions=(item:DataObj)=>{
      if(!item?.keyfield){
        return {
          height: '100%',
          document: { autofocus: false, content: '' },
          page: { layouts: ['web'] },
          toolbar: { menus: ['base', 'insert', 'tools'] },
          onSave: async()=>true
        }
      }
      const key=item.keyfield
      if(!editorOptionsMap[key]){
        editorOptionsMap[key]={
          height: '100%',
          document: {
            autofocus: false,
            content: String(item?.keyvalue ?? ''),
          },
          page: { layouts: ['web'] },
          toolbar: { menus: ['base', 'insert', 'tools'] },
          onSave: async (content:any)=>handleRichEditorSave(item,content)
        }
      } else if (editorOptionsMap[key]?.document) {
        editorOptionsMap[key].document.content = String(item?.keyvalue ?? '')
      }
      return editorOptionsMap[key]
    }
    const handleRichEditorChanged=(item:DataObj,payload:any)=>{
      const html=payload?.editor?.getHTML?.()
      if(typeof html==='string'){
        item.keyvalue=html
        const key = item?.keyfield
        if (key && editorOptionsMap[key]?.document) {
          editorOptionsMap[key].document.content = html
        }
      }
    }
    const handleRichEditorSave=async(item:DataObj,content:any)=>{
      if(content&&typeof content.html==='string'){
        item.keyvalue=content.html
      }
      const saved = await submitConfig()
      return !!saved
    }
    const getRichContentState=(item:DataObj)=>{
      return String(item?.keyvalue || '').trim() ? '已配置' : '未配置'
    }
    const openRichEditor=(item:DataObj)=>{
      activeRichItem.value = item
      richEditorTitle.value = `文档配置 - ${FontNameAndDes(item.keyname,0) || item.keyfield}`
      const options = getRichEditorOptions(item)
      if (options?.document) {
        options.document.content = String(item?.keyvalue ?? '')
      }
      richEditorVisible.value = true
    }
    const closeRichEditor=()=>{
      richEditorVisible.value = false
      activeRichItem.value = null
      richEditorTitle.value = '文档配置'
    }
    //保存配置数据
    const form = ref({});
    const submitConfig=async()=>{
      if(props.ConfData){
        try {
          Message.loading({content:"保存中",id:"updata",duration:0})
          const res:any = await saveCodeStoreConfig(props.ConfData);
          if (res === false) {
            Message.error({content:"保存失败",id:"updata",duration:2000})
            return false
          }
          Message.success({content:"保存成功",id:"updata",duration:2000})
          return true
        } catch (error) {
          Message.error({content:"保存失败",id:"updata",duration:2000})
          return false
        }
      }else{
        Message.warning({content:"配置数据不存在",id:"updata",duration:2000})
        return false
      }
    }
    //切换使用状态
    const handleUpConfigStatus=async(value:any)=>{
      if(props.ConfData){
        try {
          Message.loading({content:"更新状态中",id:"updata",duration:0})
          await upConfigStatus(props.ConfData);
          Message.success({content:"更新状态成功",id:"updata",duration:2000})
          emits('ok', true)
        } catch (error) {
          Message.loading({content:"更新状态中",id:"updata",duration:1})
        }
      }else{
        Message.warning({content:"配置数据不存在",id:"updata",duration:2000})
      }
    }
    //获取名称和描述
    const FontNameAndDes=(str:string,index:number):string=>{
      if(str){
        const str_arr=str.split("&des")
        if(index==0){
          return str_arr[index]
        }else{
          if(index==1&&str_arr.length==2)
            return str_arr[index]
          }
      }
      return ""
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
    .configdes{
      color: var(--color-neutral-4);
      margin-left: 10px;
    }
    .rich-editor-wrapper{
      width: 100%;
      height: calc(100vh - 150px);
      border-radius: 6px;
      overflow: hidden;
    }
    .rich-editor-trigger{
      width: 100%;
      display: flex;
      align-items: center;
      gap: 10px;
    }
    :deep(.rich-editor-form-item .arco-form-item-content-wrapper){
      width: 100%;
    }
    :deep(.rich-editor-drawer .arco-drawer-body){
      padding: 12px 16px;
      overflow: hidden;
    }
  </style>
  
