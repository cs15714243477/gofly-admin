<template>
  <BasicModal
    @register="registerModal"
    :title="modalTitle"
    width="520px"
    :canFullscreen="false"
    :minHeight="420"
    @ok="handleOk"
  >
    <a-alert type="warning" show-icon class="mb-12">
      <template #title>{{ lock?.bind_property_id ? '更换绑定' : '提示' }}</template>
      <template #default>
        <div class="alert-line">一把锁只能绑定一个房源。</div>
        <div class="alert-line">
          <a-space size="mini" align="center">
            <a-tag :color="lock?.bind_property_id ? 'green' : 'gray'" size="small">
              {{ lock?.bind_property_id ? '已绑定' : '未绑定' }}
            </a-tag>
            <template v-if="lock?.bind_property_id">
              <span>
                当前绑定：{{ lock.bind_property_title || '房源' }}（{{ lock.bind_property_id }}）
              </span>
            </template>
            <template v-else>
              <span>当前未绑定房源</span>
            </template>
          </a-space>
        </div>
      </template>
    </a-alert>

    <a-form :model="form" layout="vertical">
      <a-form-item label="锁" required>
        <a-input :model-value="lockDisplay" disabled />
      </a-form-item>

      <a-form-item label="房源" required>
        <a-select
          v-model="form.property_id"
          placeholder="搜索并选择房源"
          allow-search
          allow-clear
          :filter-option="false"
          :loading="searching"
          @search="onSearch"
          @popup-visible-change="onPopupVisibleChange"
          :popup-container="popupContainer"
        >
          <a-option
            v-for="it in propertyOptions"
            :key="it.id"
            :value="it.id"
            :label="`${it.title}（${it.id}）`"
          >
            <div class="opt-title">{{ it.title }}</div>
            <div class="opt-sub">{{ it.community_name || '-' }} · {{ it.sale_status_label || '-' }}</div>
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
import { bindProperty } from '../api';
import { getList as getPropertyList } from '../../houses/api';

const emit = defineEmits<{
  (e: 'success'): void;
}>();

type LockRow = {
  ttlock_lock_id: number;
  lock_name?: string;
  lock_mac?: string;
  bind_property_id?: number;
  bind_property_title?: string;
  bind_property_community_name?: string;
};

type PropertyRow = {
  id: number;
  title: string;
  community_name?: string;
  sale_status?: string;
  sale_status_label?: string;
};

const lock = ref<LockRow | null>(null);
const searching = ref(false);
const propertyOptions = ref<PropertyRow[]>([]);
const popupContainer = 'body';

const baseOptions = ref<PropertyRow[]>([]);
const form = reactive<{ property_id?: number }>({
  property_id: undefined,
});

const modalTitle = computed(() => (lock.value?.bind_property_id ? '更换绑定房源' : '绑定房源'));

const lockDisplay = computed(() => {
  if (!lock.value) return '-';
  const name = lock.value.lock_name || '未命名锁';
  const mac = lock.value.lock_mac ? ` · ${lock.value.lock_mac}` : '';
  return `${name}（${lock.value.ttlock_lock_id}）${mac}`;
});

const [registerModal, { closeModal, setModalProps }] = useModalInner(async (data: any) => {
  setModalProps({ confirmLoading: false });
  lock.value = data?.lock || null;
  const pid = Number(lock.value?.bind_property_id || 0);
  form.property_id = pid > 0 ? pid : undefined;

  baseOptions.value = [];
  if (pid > 0) {
    baseOptions.value.push({
      id: pid,
      title: lock.value?.bind_property_title || `房源（${pid}）`,
      community_name: lock.value?.bind_property_community_name,
    });
  }
  propertyOptions.value = [...baseOptions.value];
  await doSearch('');
});

let searchTimer: any = null;
const onSearch = (kw: string) => {
  if (searchTimer) clearTimeout(searchTimer);
  searchTimer = setTimeout(() => doSearch(kw), 250);
};

const onPopupVisibleChange = (visible: boolean) => {
  if (!visible) return;
  if (propertyOptions.value.length <= baseOptions.value.length) {
    doSearch('');
  }
};

const doSearch = async (kw: string) => {
  const keyword = String(kw || '').trim();
  searching.value = true;
  try {
    const resp: any = await getPropertyList({
      page: 1,
      pageSize: 20,
      ...(keyword ? { title: keyword } : {}),
    });
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
    const items = (data?.items || []) as PropertyRow[];
    const merged = [...baseOptions.value];
    const exist = new Set(merged.map((i) => i.id));
    for (const it of items) {
      if (!it?.id || exist.has(it.id)) continue;
      exist.add(it.id);
      merged.push(it);
    }
    propertyOptions.value = merged;
  } catch (e: any) {
    propertyOptions.value = [...baseOptions.value];
    Message.error(e?.message || '加载房源列表失败');
  } finally {
    searching.value = false;
  }
};

const handleOk = async () => {
  if (!lock.value?.ttlock_lock_id) {
    Message.error('缺少锁信息');
    return;
  }
  const pid = Number(form.property_id);
  if (!Number.isFinite(pid) || pid <= 0) {
    Message.warning('请选择要绑定的房源');
    return;
  }
  setModalProps({ confirmLoading: true });
  try {
    await bindProperty({ property_id: pid, ttlock_lock_id: Number(lock.value.ttlock_lock_id) });
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
.alert-line {
  line-height: 18px;
}
.alert-line + .alert-line {
  margin-top: 6px;
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
