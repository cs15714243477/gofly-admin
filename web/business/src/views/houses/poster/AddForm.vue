<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :isPadding="false"
    :loading="loading"
    :width="920"
    :height="800"
    :title="getTitle"
    :showOkBtn="false"
    :showCancelBtn="false"
  >
    <div class="form-wrap">
      <a-form ref="formRef" :model="formData" auto-label-width layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="模板名称" :rules="[{ required: true, message: '请填写模板名称' }]">
              <a-input v-model="formData.name" placeholder="请输入模板名称" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="template_type" label="模板类型" :rules="[{ required: true, message: '请选择模板类型' }]">
              <a-select v-model="formData.template_type" placeholder="选择类型" allow-clear>
                <a-option value="property">房源</a-option>
                <a-option value="agent">经纪人</a-option>
                <a-option value="festive">节日</a-option>
              </a-select>
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="preview_url" label="预览图">
              <FormImageBox v-model="formData.preview_url" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="weigh" label="排序权重">
              <a-input-number v-model="formData.weigh" :min="0" :precision="0" style="width: 100%" />
            </a-form-item>
          </a-col>

          <a-col :span="24">
            <a-form-item field="template_config" label="配置(JSON)">
              <a-textarea
                v-model="formData.template_config"
                placeholder='例如：{"bg":"#fff","blocks":[]}'
                :auto-size="{ minRows: 6, maxRows: 12 }"
                show-word-limit
              />
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-radio-group v-model="formData.status">
                <a-radio :value="0">启用</a-radio>
                <a-radio :value="1">禁用</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>

      <div class="footer">
        <a-space>
          <a-button @click="closeModal">取消</a-button>
          <a-button type="primary" :loading="loading" @click="handleSubmit">保存</a-button>
        </a-space>
      </div>
    </div>
  </BasicModal>
</template>

<script lang="ts">
import { defineComponent, ref, computed, unref } from 'vue';
import { BasicModal, useModalInner } from '/@/components/Modal';
import type { FormInstance } from '@arco-design/web-vue';
import useLoading from '@/hooks/loading';
import { cloneDeep } from 'lodash-es';
import { Message } from '@arco-design/web-vue';

import FormImageBox from '@/components/autoPlugin/Form/FormImageBox.vue';
import { save, getContent } from './api';

export default defineComponent({
  name: 'PosterAddForm',
  components: { BasicModal, FormImageBox },
  emits: ['success'],
  setup(_, { emit }) {
    const isUpdate = ref(false);
    const formRef = ref<FormInstance>();

    const baseData: any = {
      id: 0,
      name: '',
      preview_url: '',
      template_config: '',
      template_type: 'property',
      status: 0,
      weigh: 0,
    };

    const formData = ref<any>(cloneDeep(baseData));
    const { loading, setLoading } = useLoading();

    const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
      formRef.value?.resetFields();
      setModalProps({ confirmLoading: false });
      isUpdate.value = !!data?.isUpdate;
      formData.value = cloneDeep(baseData);

      if (unref(isUpdate) && data?.record?.id) {
        setLoading(true);
        const resp: any = await getContent({ id: data.record.id });
        const detail =
          resp?.data?.data
            ? resp.data.data
            : resp?.data
              ? resp.data
              : resp;
        formData.value = Object.assign({}, formData.value, detail || {});
        formData.value.status = Number(formData.value.status) || 0;
        formData.value.weigh = Number(formData.value.weigh) || 0;
        setLoading(false);
      }
    });

    const getTitle = computed(() => (!unref(isUpdate) ? '新增海报模板' : '编辑海报模板'));

    const handleSubmit = async () => {
      try {
        const res = await formRef.value?.validate();
        if (res) {
          const firstField = Object.keys(res as any)[0];
          Message.warning((res as any)[firstField]?.[0]?.message || '请先完善必填项');
          (formRef.value as any)?.scrollToField?.(firstField);
          return;
        }
        // 简单校验 JSON
        const cfg = (formData.value.template_config || '').toString().trim();
        if (cfg) {
          try {
            JSON.parse(cfg);
          } catch (e) {
            Message.warning('配置(JSON)格式不正确');
            return;
          }
        }

        setLoading(true);
        Message.loading({ content: '提交中', id: 'poster_save', duration: 0 });
        await save(cloneDeep(unref(formData)));
        Message.success({ content: '提交成功', id: 'poster_save', duration: 1500 });
        closeModal();
        emit('success');
        setLoading(false);
      } catch (e) {
        setLoading(false);
        Message.clear('top');
      }
    };

    return {
      registerModal,
      closeModal,
      getTitle,
      formRef,
      formData,
      loading,
      handleSubmit,
    };
  },
});
</script>

<style lang="less" scoped>
.form-wrap {
  padding: 16px 20px 8px;
}
.footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  border-top: 1px solid var(--color-border);
}
</style>

