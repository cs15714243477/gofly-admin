<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :isPadding="false"
    :loading="loading"
    :width="860"
    :height="760"
    :title="getTitle"
    :showOkBtn="false"
    :showCancelBtn="false"
  >
    <div class="form-wrap">
      <a-form ref="formRef" :model="formData" auto-label-width layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item field="name" label="姓名" :rules="[{ required: true, message: '请填写姓名' }]">
              <a-input v-model="formData.name" placeholder="请输入姓名" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="nickname" label="昵称">
              <a-input v-model="formData.nickname" placeholder="请输入昵称" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="mobile" label="手机号" :rules="[{ required: true, message: '请填写手机号' }]">
              <a-input v-model="formData.mobile" placeholder="请输入手机号" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="email" label="邮箱">
              <a-input v-model="formData.email" placeholder="请输入邮箱" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="username" label="用户名">
              <a-input v-model="formData.username" placeholder="请输入用户名" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="store_id" label="所属门店">
              <a-select
                v-model="formData.store_id"
                allow-clear
                allow-search
                :filter-option="false"
                :loading="storeLoading"
                placeholder="输入门店名称筛选"
                @search="onStoreSearch"
                @popup-visible-change="onStorePopupVisibleChange"
              >
                <a-option v-for="it in storeOptions" :key="it.id" :value="it.id" :label="it.name">
                  <div class="store-option">
                    <div class="store-name">{{ it.name }}</div>
                    <div class="store-addr" v-if="it.address">{{ it.address }}</div>
                  </div>
                </a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item field="title" label="头衔">
              <a-input v-model="formData.title" placeholder="例如：金牌经纪人" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="24">
            <a-form-item field="introduction" label="自我介绍">
              <a-textarea v-model="formData.introduction" :auto-size="{ minRows: 3, maxRows: 6 }" placeholder="请输入自我介绍" />
            </a-form-item>
          </a-col>

          <a-col :span="24">
            <a-form-item field="remark" label="备注">
              <a-input v-model="formData.remark" placeholder="备注信息" allow-clear />
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="status" label="状态">
              <a-switch v-model="formData.status" :checked-value="0" :unchecked-value="1" type="round">
                <template #checked>启用</template>
                <template #unchecked>禁用</template>
              </a-switch>
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item field="can_manage_properties" label="可维护房源" tooltip="开启后，该经纪人可维护房源（默认不可维护）">
              <a-switch v-model="formData.can_manage_properties" :checked-value="1" :unchecked-value="0" type="round">
                <template #checked>可维护</template>
                <template #unchecked>不可维护</template>
              </a-switch>
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
import { save, getContent } from '../api';
import { Message } from '@arco-design/web-vue';
import { getList as getStoreList, getContent as getStoreContent } from '../../stores/api';

export default defineComponent({
  name: 'BrokerAddForm',
  components: { BasicModal },
  emits: ['success'],
  setup(_, { emit }) {
    const isUpdate = ref(false);
    const formRef = ref<FormInstance>();

    const baseData: any = {
      id: 0,
      username: '',
      name: '',
      nickname: '',
      remark: '',
      email: '',
      mobile: '',
      // 角色固定为 1（不在界面展示）
      role: 1,
      // 默认不可维护房源
      can_manage_properties: 0,
      store_id: undefined as any,
      title: '',
      introduction: '',
      status: 0,
    };

    const formData = ref<any>(cloneDeep(baseData));
    const { loading, setLoading } = useLoading();

    const storeOptions = ref<any[]>([]);
    const storeLoading = ref(false);
    const storeKeyword = ref('');

    const unwrapList = (resp: any) => {
      const data =
        resp?.items
          ? resp
          : resp?.data?.items
            ? resp.data
            : resp?.data?.data?.items
              ? resp.data.data
              : resp?.data
                ? resp.data
                : resp;
      return data || {};
    };

    const fetchStoreOptions = async (keyword = '') => {
      try {
        storeLoading.value = true;
        const resp: any = await getStoreList({
          page: 1,
          pageSize: 50,
          name: keyword,
          status: '',
        });
        const data = unwrapList(resp);
        storeOptions.value = (data?.items || []).map((it: any) => ({
          id: Number(it.id),
          name: it.name,
          address: it.address,
        }));
      } catch (e) {
        storeOptions.value = [];
      } finally {
        storeLoading.value = false;
      }
    };

    const ensureStoreSelectedOption = async (storeId: any) => {
      const id = Number(storeId) || 0;
      if (!id) return;
      if (storeOptions.value.some((it) => Number(it.id) === id)) return;
      try {
        const resp: any = await getStoreContent({ id });
        const detail = resp?.data?.data ? resp.data.data : resp?.data ? resp.data : resp;
        if (detail?.id) {
          storeOptions.value.unshift({
            id: Number(detail.id),
            name: detail.name,
            address: detail.address,
          });
        }
      } catch (e) {}
    };

    const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
      formRef.value?.resetFields();
      setModalProps({ confirmLoading: false });
      isUpdate.value = !!data?.isUpdate;
      formData.value = cloneDeep(baseData);
      storeKeyword.value = '';
      await fetchStoreOptions('');

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
        // 数字字段兜底
        formData.value.store_id = Number(formData.value.store_id) || undefined;
        formData.value.status = Number(formData.value.status) || 0;
        formData.value.can_manage_properties = Number((formData.value as any).can_manage_properties) || 0;
        await ensureStoreSelectedOption(formData.value.store_id);
        setLoading(false);
      }
    });

    const getTitle = computed(() => (!unref(isUpdate) ? '新增经纪人' : '编辑经纪人'));

    const onStoreSearch = async (val: string) => {
      storeKeyword.value = (val || '').trim();
      await fetchStoreOptions(storeKeyword.value);
    };

    const onStorePopupVisibleChange = async (visible: boolean) => {
      if (!visible) return;
      if (storeOptions.value.length === 0) {
        await fetchStoreOptions(storeKeyword.value);
      }
    };

    const handleSubmit = async () => {
      try {
        const res = await formRef.value?.validate();
        if (res) {
          const firstField = Object.keys(res as any)[0];
          Message.warning((res as any)[firstField]?.[0]?.message || '请先完善必填项');
          (formRef.value as any)?.scrollToField?.(firstField);
          return;
        }
        setLoading(true);
        Message.loading({ content: '提交中', id: 'broker_save', duration: 0 });
        // 角色固定为 1
        formData.value.role = 1;
        // select 可能回传 string，统一转成 number
        formData.value.store_id = Number(formData.value.store_id) || 0;
        formData.value.can_manage_properties = Number(formData.value.can_manage_properties) || 0;
        await save(cloneDeep(unref(formData)));
        Message.success({ content: '提交成功', id: 'broker_save', duration: 1500 });
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
      storeOptions,
      storeLoading,
      onStoreSearch,
      onStorePopupVisibleChange,
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

.store-option {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}
.store-name {
  font-size: 13px;
  color: var(--color-text-1);
}
.store-addr {
  margin-top: 2px;
  font-size: 12px;
  color: var(--color-text-3);
}
</style>
