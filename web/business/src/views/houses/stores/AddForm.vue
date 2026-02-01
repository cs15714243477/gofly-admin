<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :isPadding="false"
    :loading="loading"
    width="900px"
    @height-change="onHeightChange"
    :minHeight="modelHeight"
    :title="getTitle"
    @ok="handleSubmit"
  >
    <div class="addFormbox" :style="{'min-height':`${windHeight}px`}">
      <div class="tabs-content" >
        <a-form ref="formRef" :model="formData" auto-label-width>
          <div class="content_box">
              <!--基础信息-->
              <a-scrollbar   style="overflow: auto;" :style="{height:`${windHeight}px`}">
                <div class="besecontent" >
                  <a-row :gutter="16">
										<a-col :span="12">
											<a-form-item field="name" label="门店名称" :rules="[{ required: true, message: '请填写门店名称' }]" >
													<a-input v-model="formData.name" placeholder="输入门店名称" />
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="area_ids" label="省市区" :rules="[cascaderRule('请选择省市区')]"  >
                        <a-cascader
                          v-model="formData.area_ids"
                          :options="areaOptions"
                          :load-more="loadAreaMore"
                          :field-names="{ value: 'value', label: 'label', children: 'children', isLeaf: 'isLeaf' }"
                          path-mode
                          placeholder="选择省 / 市 / 区"
                          allow-clear
                        />
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="address_detail" label="详细地址" :rules="[{ required: true, message: '请填写详细地址' }]"  >
													<a-input v-model="formData.address_detail" placeholder="街道/门牌号/楼栋等" allow-clear />
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="contact_phone" label="联系电话" :rules="[{ required: true, message: '请填写联系电话' }]" >
													<a-input v-model="formData.contact_phone" placeholder="输入联系电话" />
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="manager_name" label="店长姓名" :rules="[{ required: true, message: '请填写店长姓名' }]"  >
													<a-input v-model="formData.manager_name" placeholder="输入店长姓名" />
											</a-form-item>
										</a-col>
										<a-col :span="24">
											<a-form-item field="images" label="门店图片(多图,逗号分隔)"  >
													<FormImagesBox v-model="formData.images" placeholder="请选择门店图片"/>
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="status" label="状态"  >
													<a-radio-group v-model="formData.status">
                            <a-radio :value="0">启用</a-radio>
                            <a-radio :value="1">禁用</a-radio>
                          </a-radio-group>
											</a-form-item>
										</a-col>
										<a-col :span="12">
											<a-form-item field="weigh" label="排序权重"  >
													<a-input-number v-model="formData.weigh" placeholder="输入排序权重" />
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
  import FormImagesBox from '@/components/autoPlugin/Form/FormImagesBox.vue';
  //api
  import { save,getContent } from './api';
  import { FormInstance,Message} from '@arco-design/web-vue';
  import { defHttp } from '@/utils/http';
  export default defineComponent({
    name: 'AddBook',
    components: { BasicModal, FormImagesBox},
    emits: ['success'],
    setup(_, { emit }) {
      const isUpdate = ref(false);
      const modelHeight= ref(520);
      const windHeight= ref(520);
      //表单
      const formRef = ref<FormInstance>();
      //表单字段
      const basedata={
            id:0,
            name:"",
            address:"",
            contact_phone:"",
            manager_name:"",
            images:"",
            status:0,
            weigh:0,
            area_ids: [] as (number | string)[],
            address_detail: "",

        }
      const formData = ref(basedata)
      const areaOptions = ref<any[]>([]);
      const cascaderRule = (message: string) => ({
        validator: (value: any, callback: (error?: string) => void) => {
          if (Array.isArray(value) && value.length >= 3) {
            callback();
            return;
          }
          callback(message);
        },
      });

      const ensureAreaRootLoaded = async () => {
        if (areaOptions.value.length > 0) return;
        const list: any = await defHttp.get({ url: '/area/getMoreList', params: { pid: 0 } }, { errorMessageMode: 'none' });
        areaOptions.value = (list || []).map((it: any) => ({
          ...it,
          value: it.id,
          label: it.name,
          isLeaf: it.isLeaf === true,
          children: undefined,
        }));
      };

      const loadAreaMore = async (option: any, done: (children?: any[]) => void) => {
        try {
          const pid = option?.id ?? option?.value ?? option;
          const list: any = await defHttp.get({ url: '/area/getMoreList', params: { pid } }, { errorMessageMode: 'none' });
          const children = (list || []).map((it: any) => ({
            ...it,
            value: it.id,
            label: it.name,
            isLeaf: it.isLeaf === true,
            children: undefined,
          }));
          done(children);
        } catch (e) {
          done([]);
        }
      };

      const ensureAreaPathLoaded = async (ids: any[]) => {
        await ensureAreaRootLoaded();
        let list: any[] = areaOptions.value;
        for (let i = 0; i < (ids || []).length; i++) {
          const id = ids[i];
          const node = Array.isArray(list) ? list.find((it) => String(it.value ?? it.id) === String(id)) : null;
          if (!node) return;
          const hasNext = i < ids.length - 1;
          if (hasNext) {
            if (!Array.isArray(node.children) || node.children.length === 0) {
              const children: any = await defHttp.get({ url: '/area/getMoreList', params: { pid: node.id } }, { errorMessageMode: 'none' });
              node.children = (children || []).map((it: any) => ({
                ...it,
                value: it.id,
                label: it.name,
                isLeaf: it.isLeaf === true,
                children: undefined,
              }));
            }
            list = node.children || [];
          }
        }
      };

      const getAreaLabelsByIds = (ids: any[]): string[] => {
        const out: string[] = [];
        let list: any[] = areaOptions.value;
        for (const id of ids || []) {
          const node = Array.isArray(list) ? list.find((it) => String(it.value ?? it.id) === String(id)) : null;
          if (!node) break;
          out.push(node.label ?? node.name);
          list = node.children || [];
        }
        return out;
      };

      const buildFullAddress = () => {
        const ids = formData.value.area_ids || [];
        const names = getAreaLabelsByIds(ids);
        const detail = (formData.value.address_detail || '').toString().trim();
        return [names.join(' '), detail].filter(Boolean).join(' ').trim();
      };

      const initAreaFromAddress = async (address: string) => {
        const raw = (address || '').toString().trim();
        if (!raw) return;
        const parts = raw.split(/\s+/).filter(Boolean);
        if (parts.length < 4) {
          formData.value.address_detail = raw;
          return;
        }
        await ensureAreaRootLoaded();
        formData.value.address_detail = parts.slice(3).join(' ');
      };
      const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
          formRef.value?.resetFields()
          setModalProps({ confirmLoading: false });
          isUpdate.value = !!data?.isUpdate;
          await ensureAreaRootLoaded();
          if (unref(isUpdate)) {
            formData.value=cloneDeep(data.record)
            setLoading(true);
            const mewdata = await getContent({id:data.record.id});
            formData.value=Object.assign({},formData.value,mewdata)
            // 地址回填到“省市区 + 详细地址”
            await initAreaFromAddress(formData.value.address || "");
            setLoading(false);
          }else{
            formData.value=cloneDeep(basedata)
          }
      });
      const getTitle = computed(() => (!unref(isUpdate) ? '新增数据' : '编辑数据'));
     //提交数据
     const { loading, setLoading } = useLoading();
     const handleSubmit = async () => {
      try {
          const res = await formRef.value?.validate();
          if (!res) {
            setLoading(true);
            Message.loading({content:"提交中",id:"submit",duration:0})
            let savedata=cloneDeep(unref(formData))
            // 保存前把结构化地址拼回 address 字段
            await ensureAreaPathLoaded(savedata.area_ids || []);
            savedata.address = buildFullAddress();
            await save(savedata);
            Message.success({content:"提交成功",id:"submit",duration:1500})
            closeModal()
            emit('success');
            setLoading(false);
          }
        } catch (error) {
          setLoading(false);
          Message.loading({content:"提交中",id:"submit",duration:1})
        }
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
        modelHeight,
        onHeightChange,windHeight,
        areaOptions,
        loadAreaMore,
        cascaderRule,
      };
    },
  });
</script>
<style lang="less" scoped>
  @import '@/assets/style/formlayer.less';
  .besecontent{
    padding: 15px 30px 5px 30px!important;
  }
</style>

