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
           <a-tooltip content="配置字段命名说明：_txt结尾是多行文本（wxapp_property_form_options 已支持“选项编辑器”），_read结尾是只读，_switch结尾开关,#注释加 &des 是字段说明。" position="top" mini>
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
             <a-col :span="getColSpan(item)">
                 <a-form-item
                 :label="FontNameAndDes(item.keyname,0)"
                 :field="item.keyfield"
                 :extra="FontNameAndDes(item.keyname,1)"
                 >
                   <div v-if="isOptionsEditorField(item)" class="options-editor">
                     <div class="options-editor-toolbar">
                       <a-space :size="10" wrap>
                         <a-radio-group
                           v-model="getOptionsEditor(item).mode"
                           type="button"
                           size="small"
                           @change="(v:any)=>handleTxtModeChange(item,v)"
                         >
                           <a-radio value="options">选项编辑</a-radio>
                           <a-radio value="raw">原始文本</a-radio>
                         </a-radio-group>
                         <a-tag color="arcoblue" size="small">共 {{ getOptionsEditor(item).options.length }} 项</a-tag>
                         <a-button
                           v-if="getOptionsEditor(item).mode==='options'"
                           type="outline"
                           size="small"
                           @click="addOption(item)"
                         >
                           新增一项
                         </a-button>
                         <a-button
                           v-if="getOptionsEditor(item).mode==='options'"
                           type="text"
                           size="small"
                           status="danger"
                           @click="clearOptions(item)"
                         >
                           清空
                         </a-button>
                       </a-space>
                       <a-tooltip content="保存时：选项编辑=自动生成 value=label 多行文本；原始文本=按文本内容保存。支持 value=label / value:label / value,label / value。" position="top" mini>
                         <icon-question-circle-fill class="configdes" />
                       </a-tooltip>
                     </div>

                     <template v-if="getOptionsEditor(item).mode==='options'">
                       <a-table
                         class="options-editor-table"
                         size="small"
                         :data="getOptionsEditor(item).options"
                         :pagination="false"
                         :bordered="{ wrapper: true, cell: true }"
                         :scroll="{ x: '100%' }"
                         row-key="id"
                       >
                         <template #columns>
                           <a-table-column title="排序" :width="120" align="center">
                             <template #cell="{ rowIndex }">
                               <a-space size="mini">
                                 <a-button
                                   size="mini"
                                   type="outline"
                                   :disabled="rowIndex===0"
                                   @click="moveOption(item,rowIndex,-1)"
                                 >
                                   上移
                                 </a-button>
                                 <a-button
                                   size="mini"
                                   type="outline"
                                   :disabled="rowIndex===getOptionsEditor(item).options.length-1"
                                   @click="moveOption(item,rowIndex,1)"
                                 >
                                   下移
                                 </a-button>
                               </a-space>
                             </template>
                           </a-table-column>
                           <a-table-column title="值(value)" :width="240">
                             <template #cell="{ record }">
                               <a-input v-model="record.value" allow-clear placeholder="例如：on_sale / 万 / 低层" />
                             </template>
                           </a-table-column>
                           <a-table-column title="显示(label)">
                             <template #cell="{ record }">
                               <a-input v-model="record.label" allow-clear placeholder="例如：在售 / 万 / 低层 (1-6)" />
                             </template>
                           </a-table-column>
                           <a-table-column title="操作" :width="90" align="center">
                             <template #cell="{ rowIndex }">
                               <a-button type="text" status="danger" @click="removeOption(item,rowIndex)">删除</a-button>
                             </template>
                           </a-table-column>
                         </template>
                       </a-table>
                       <a-alert
                         class="options-editor-tip"
                         type="info"
                         show-icon
                         title="格式说明"
                       >
                         每行一条：value=label；label 为空时默认等于 value。
                       </a-alert>
                     </template>

                     <a-textarea
                       v-else
                       v-model="item.keyvalue"
                       :placeholder="`填写${FontNameAndDes(item.keyname,0)}`"
                       :auto-size="{ minRows: 6, maxRows: 16 }"
                     />
                   </div>

                   <a-textarea
                     v-else-if="item.keyfield.indexOf('_txt')>-1"
                     v-model="item.keyvalue"
                     :placeholder="`填写${FontNameAndDes(item.keyname,0)}`"
                     :auto-size="{ minRows: 3, maxRows: 10 }"
                   />
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
 </template>
 
 <script lang="ts" setup>
   import { computed, reactive, ref,PropType, watch } from 'vue';
   //api
   import { DataObj, menuItem,saveCodeStoreConfig,upConfigStatus} from '@/api/datacenter/configuration';
   import { Message, Modal } from '@arco-design/web-vue';
   const emits = defineEmits(['ok'])
   const props = defineProps({
     ConfData: {
       type: Object as PropType<menuItem>,
     },
   })

   type TxtEditorMode = 'options' | 'raw'
   type OptionsItem = { id: string; value: string; label: string }
   type OptionsEditorState = { mode: TxtEditorMode; options: OptionsItem[] }
   const optionsEditorMap = reactive<Record<string, OptionsEditorState>>({})
   const isPropertyFormOptionsConfig = computed(() => {
     return props.ConfData?.pluginident === 'wxapp_property_form_options' || props.ConfData?.name === 'property_form_options'
   })

   const makeId = () => `${Date.now()}_${Math.random().toString(16).slice(2)}`

   const isTxtField = (item: DataObj) => {
     return !!item?.keyfield && item.keyfield.indexOf('_txt') > -1
   }

   const isOptionsEditorField = (item: DataObj) => {
     return isPropertyFormOptionsConfig.value && isTxtField(item)
   }

   const normalizeText = (raw: any) => {
     let v = String(raw ?? '').trim()
     if (!v) return ''
     // 兼容：用户可能写成字面量 \\n
     v = v.replaceAll('\\n', '\n')
     v = v.replaceAll('\r\n', '\n').replaceAll('\r', '\n')
     return v
   }

   const parseOptionsText = (raw: any): OptionsItem[] => {
     const text = normalizeText(raw)
     if (!text) return []
     const lines = text.split('\n')
     const out: OptionsItem[] = []
     for (const line of lines) {
       let l = String(line ?? '').trim()
       if (!l) continue
       // 行内注释（与后端保持一致）
       const commentIndex = l.indexOf('#')
       if (commentIndex >= 0) {
         l = l.slice(0, commentIndex).trim()
       }
       if (!l) continue
       let value = ''
       let label = ''
       if (l.includes('=')) {
         const parts = l.split('=')
         value = String(parts[0] ?? '').trim()
         label = String(parts.slice(1).join('=') ?? '').trim()
       } else if (l.includes(':')) {
         const parts = l.split(':')
         value = String(parts[0] ?? '').trim()
         label = String(parts.slice(1).join(':') ?? '').trim()
       } else if (l.includes(',')) {
         const parts = l.split(',')
         value = String(parts[0] ?? '').trim()
         label = String(parts.slice(1).join(',') ?? '').trim()
       } else {
         value = l.trim()
         label = value
       }
       if (!value) continue
       if (!label) label = value
       out.push({ id: makeId(), value, label })
     }
     return out
   }

   const formatOptionsText = (options: OptionsItem[]) => {
     const lines: string[] = []
     for (const opt of options || []) {
       const value = String(opt?.value ?? '').trim()
       if (!value) continue
       const label = String(opt?.label ?? '').trim() || value
       lines.push(`${value}=${label}`)
     }
     return lines.join('\n')
   }

   const getOptionsEditor = (item: DataObj): OptionsEditorState => {
     const key = item?.keyfield || ''
     if (!key) {
       return { mode: 'raw', options: [] }
     }
     if (!optionsEditorMap[key]) {
       optionsEditorMap[key] = {
         mode: isOptionsEditorField(item) ? 'options' : 'raw',
         options: parseOptionsText(item?.keyvalue),
       }
     }
     return optionsEditorMap[key]
   }

   const handleTxtModeChange = (item: DataObj, mode: TxtEditorMode) => {
     const editor = getOptionsEditor(item)
     editor.mode = mode
     if (mode === 'raw') {
       item.keyvalue = formatOptionsText(editor.options)
       return
     }
     // 切回选项模式时，重新从文本解析，避免“文本编辑后未同步”
     editor.options = parseOptionsText(item.keyvalue)
   }

   const addOption = (item: DataObj) => {
     const editor = getOptionsEditor(item)
     editor.options.push({ id: makeId(), value: '', label: '' })
   }

   const removeOption = (item: DataObj, index: number) => {
     const editor = getOptionsEditor(item)
     editor.options.splice(index, 1)
   }

   const moveOption = (item: DataObj, index: number, dir: -1 | 1) => {
     const editor = getOptionsEditor(item)
     const next = index + dir
     if (index < 0 || next < 0 || next >= editor.options.length) return
     const tmp = editor.options[index]
     editor.options[index] = editor.options[next]
     editor.options[next] = tmp
   }

   const clearOptions = (item: DataObj) => {
     Modal.warning({
       title: '确认清空？',
       content: '清空后将删除该字段下的全部选项，且无法恢复。',
       hideCancel: false,
       onOk: () => {
         const editor = getOptionsEditor(item)
         editor.options = []
       },
     })
   }

   const getColSpan = (item: DataObj) => {
     return isOptionsEditorField(item) ? 24 : 12
   }

   watch(
     () => props.ConfData?.name,
     () => {
       Object.keys(optionsEditorMap).forEach((key) => delete optionsEditorMap[key])
       if (props.ConfData?.data && isPropertyFormOptionsConfig.value) {
         props.ConfData.data.forEach((item) => {
           if (isOptionsEditorField(item)) {
             optionsEditorMap[item.keyfield] = { mode: 'options', options: parseOptionsText(item.keyvalue) }
           }
         })
       }
     },
     { immediate: true }
   )
   //保存配置数据
   const form = ref({});
   const submitConfig=async()=>{
     if(props.ConfData){
       try {
         // 保存前：把“选项编辑”模式的内容同步为多行文本（value=label）
         if (props.ConfData?.data && isPropertyFormOptionsConfig.value) {
           props.ConfData.data.forEach((item) => {
             if (!isOptionsEditorField(item)) return
             const editor = optionsEditorMap[item.keyfield]
             if (editor?.mode === 'options') {
               item.keyvalue = formatOptionsText(editor.options)
             }
           })
         }
         Message.loading({content:"保存中",id:"updata",duration:0})
         await saveCodeStoreConfig(props.ConfData);
         Message.success({content:"保存成功",id:"updata",duration:2000})
       } catch (error) {
         Message.error({content:"",id:"updata",duration:1})
       }
     }else{
       Message.warning({content:"配置数据不存在",id:"updata",duration:2000})
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
   .options-editor{
     width: 100%;
   }
   .options-editor-toolbar{
     display: flex;
     align-items: center;
     justify-content: space-between;
     gap: 10px;
     margin-bottom: 10px;
   }
   .options-editor-table{
     width: 100%;
   }
   .options-editor-tip{
     margin-top: 10px;
   }
 </style>
 
