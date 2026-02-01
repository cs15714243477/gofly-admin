<template>
  <BasicModal
    @register="registerModal"
    title="绑定智能锁"
    width="560px"
    :canFullscreen="false"
    :minHeight="360"
    @ok="handleOk"
  >
    <a-alert type="warning" show-icon class="mb-12">
      <template #title>提示</template>
      <template #default>
        一把锁只能绑定一个房源。
      </template>
    </a-alert>

    <a-form :model="form" layout="vertical">
      <a-form-item label="房源" required>
        <a-input :model-value="propertyDisplay" disabled />
      </a-form-item>

      <a-form-item label="选择锁" required>
        <a-select
          v-model="form.ttlock_lock_id"
          placeholder="请选择锁"
          allow-search
          allow-clear
          :loading="loadingLocks"
          :popup-container="popupContainer"
        >
          <a-option
            v-for="it in lockOptions"
            :key="it.ttlock_lock_id"
            :value="it.ttlock_lock_id"
            :label="`${it.lock_name || '未命名锁'}（${it.ttlock_lock_id}）`"
          >
            <div class="opt-title">{{ it.lock_name || '未命名锁' }}</div>
            <div class="opt-sub">{{ it.ttlock_lock_id }} · {{ it.lock_mac || '-' }} · 电量{{ it.battery ?? 0 }}%</div>
          </a-option>
        </a-select>
      </a-form-item>
    </a-form>
  </BasicModal>
</template>

<script lang="ts" setup>
import { computed, reactive, ref } from 'vue';
import { BasicModal, useModalInner } from '/@/components/Modal';
import { Message } from '@arco-design/web-vue';
import { bindProperty, getLockList } from '../api';

const emit = defineEmits<{
  (e: 'success'): void;
}>();

type LockRow = {
  ttlock_lock_id: number;
  lock_name?: string;
  lock_mac?: string;
  battery?: number;
};

const property = ref<{ id: number; title?: string } | null>(null);
const lockOptions = ref<LockRow[]>([]);
const loadingLocks = ref(false);
const popupContainer = 'body';

const form = reactive<{ property_id: number; ttlock_lock_id?: number }>({
  property_id: 0,
  ttlock_lock_id: undefined,
});

const propertyDisplay = computed(() => {
  if (!property.value?.id) return '-';
  return `${property.value.title || '房源'}（${property.value.id}）`;
});

const [registerModal, { closeModal, setModalProps }] = useModalInner(async (data: any) => {
  setModalProps({ confirmLoading: false });
  property.value = data?.property || null;
  form.property_id = Number(property.value?.id || 0);
  form.ttlock_lock_id = undefined;
  await fetchLocks();
});

const fetchLocks = async () => {
  loadingLocks.value = true;
  try {
    const resp: any = await getLockList({ page: 1, pageSize: 200 });
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
    lockOptions.value = (data?.items || []) as LockRow[];
  } catch (e: any) {
    lockOptions.value = [];
    Message.error(e?.message || '加载锁列表失败');
  } finally {
    loadingLocks.value = false;
  }
};

const handleOk = async () => {
  if (!form.property_id) {
    Message.error('缺少房源信息');
    return;
  }
  const lockId = Number(form.ttlock_lock_id);
  if (!Number.isFinite(lockId) || lockId <= 0) {
    Message.warning('请选择锁');
    return;
  }
  setModalProps({ confirmLoading: true });
  try {
    await bindProperty({ property_id: Number(form.property_id), ttlock_lock_id: lockId });
    Message.success('绑定成功');
    emit('success');
    closeModal();
  } finally {
    setModalProps({ confirmLoading: false });
  }
};
</script>

<style scoped>
.mb-12 {
  margin-bottom: 12px;
}
.opt-title {
  font-weight: 600;
  line-height: 20px;
}
.opt-sub {
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 16px;
}
</style>
