<template>
  <div class="container">
    <page-card breadcrumb scrollPage>
      <!-- 头部（参考智能门锁页面：卡片 + 图标 + 标题 + CTA） -->
      <a-card class="broker-hero" :bordered="false">
        <div class="hero-inner">
          <div class="hero-left">
            <div class="hero-icon"><icon-user-group /></div>
            <div class="hero-text">
              <div class="hero-title">经纪人</div>
              <div class="hero-sub">账号信息、门店归属、可维护房源/智能锁权限</div>
            </div>
          </div>
          <div class="hero-actions">
            <a-space>
              <a-button @click="fetchData">
                <template #icon><icon-refresh /></template>
                刷新
              </a-button>
              <a-button type="primary" @click="handleCreate" v-permission="['add']">
                <template #icon><icon-plus /></template>
                新增经纪人
              </a-button>
            </a-space>
          </div>
        </div>
      </a-card>

      <!-- 搜索栏（卡片化） -->
      <a-card class="broker-filter" :bordered="false">
        <div class="search-row">
          <a-space wrap>
            <a-input style="width: 180px" v-model="formModel.name" placeholder="姓名" allow-clear />
            <a-input style="width: 180px" v-model="formModel.mobile" placeholder="手机号" allow-clear />
            <a-input style="width: 180px" v-model="formModel.username" placeholder="用户名" allow-clear />
            <a-select style="width: 160px" v-model="formModel.can_manage_properties" placeholder="可维护房源" allow-clear>
              <a-option :value="1">可维护</a-option>
              <a-option :value="0">不可维护</a-option>
            </a-select>
            <a-select style="width: 160px" v-model="formModel.can_manage_locks" placeholder="可管理智能锁" allow-clear>
              <a-option :value="1">可管理</a-option>
              <a-option :value="0">不可管理</a-option>
            </a-select>
            <a-select style="width: 140px" v-model="formModel.audit_status" placeholder="审核状态" allow-clear>
              <a-option value="pending">待审核</a-option>
              <a-option value="approved">已通过</a-option>
              <a-option value="rejected">已拒绝</a-option>
            </a-select>
            <a-select style="width: 140px" v-model="formModel.status" placeholder="状态" allow-clear>
              <a-option :value="0">启用</a-option>
              <a-option :value="1">禁用</a-option>
            </a-select>
            <a-button type="primary" @click="handleSearch">
              <template #icon><icon-search /></template>
              查询
            </a-button>
            <a-button @click="handleReset">重置</a-button>
          </a-space>
        </div>
      </a-card>

      <!-- 列表 -->
      <a-card class="broker-list" :bordered="false">
        <a-spin :loading="loading" style="width: 100%">
          <a-grid :cols="24" :col-gap="16" :row-gap="16" class="broker-grid" v-if="renderData.length">
            <a-grid-item
              v-for="record in renderData"
              :key="record.id"
              :span="{ xs: 24, sm: 12, md: 8, lg: 6, xl: 6, xxl: 6 }"
            >
              <a-card class="broker-card" :bordered="false" :body-style="{ padding: '14px' }">
                <div class="card-head">
                  <a-avatar :size="48" class="avatar">
                    <img v-if="record.avatar" :src="GetFullPath(record.avatar)" alt="avatar" />
                    <span v-else>{{ getAvatarText(record) }}</span>
                  </a-avatar>
                  <div class="head-main">
                    <div class="name-row">
                      <div class="name">{{ record.name || record.nickname || record.username || '未命名' }}</div>
                      <a-badge
                        :status="record.status === 0 ? 'success' : 'default'"
                        :text="record.status === 0 ? '启用' : '禁用'"
                      />
                    </div>
                    <div class="sub-row">
                      <span class="title">{{ record.title || '经纪人' }}</span>
                      <span class="sep">·</span>
                      <span class="store">{{ record.store_name || '未绑定门店' }}</span>
                    </div>
                  </div>
                </div>

                <div class="info">
                  <div class="info-line">
                    <icon-phone />
                    <span class="text">{{ record.mobile || '—' }}</span>
                  </div>
                  <div class="info-line audit-line">
                    <icon-check-circle />
                    <span class="label">审核</span>
                    <a-tag
                      size="small"
                      bordered
                      :color="auditStatusColor(record.audit_status)"
                    >
                      {{ auditStatusText(record.audit_status) }}
                    </a-tag>
                    <a-switch
                      v-permission="['audit']"
                      size="small"
                      type="round"
                      :model-value="isAuditApproved(record)"
                      :loading="isToggleLoading(record, 'audit')"
                      @change="(val) => handleAuditSwitch(record, val)"
                    >
                      <template #checked>通过</template>
                      <template #unchecked>待审</template>
                    </a-switch>
                  </div>
                  <div class="info-line">
                    <icon-location />
                    <span class="text ellipsis">{{ record.store_address || '—' }}</span>
                  </div>
                </div>

                <div class="tags">
                  <div class="tag-item">
                    <a-tag
                      size="small"
                      :color="Number(record.can_manage_properties) === 1 ? 'green' : 'gray'"
                      bordered
                    >
                      {{ Number(record.can_manage_properties) === 1 ? '可维护房源' : '不可维护房源' }}
                    </a-tag>
                    <a-switch
                      v-permission="['canManageProperties']"
                      v-model="record.can_manage_properties"
                      :checked-value="1"
                      :unchecked-value="0"
                      size="small"
                      type="round"
                      :loading="isToggleLoading(record, 'canManageProperties')"
                      @change="(val) => handleCanManageProperties(record, val)"
                    />
                  </div>
                  <div class="tag-item">
                    <a-tag size="small" :color="Number(record.can_manage_locks) === 1 ? 'arcoblue' : 'gray'" bordered>
                      {{ Number(record.can_manage_locks) === 1 ? '可管智能锁' : '不可管智能锁' }}
                    </a-tag>
                    <a-switch
                      v-permission="['canManageLocks']"
                      v-model="record.can_manage_locks"
                      :checked-value="1"
                      :unchecked-value="0"
                      size="small"
                      type="round"
                      :loading="isToggleLoading(record, 'canManageLocks')"
                      @change="(val) => handleCanManageLocks(record, val)"
                    />
                  </div>
                </div>

                <div class="actions">
                  <a-button type="text" size="small" @click="handleEdit(record)" v-permission="['add']">
                    <template #icon><icon-edit /></template>
                    编辑
                  </a-button>
                  <a-space size="mini">
                    <a-switch
                      type="round"
                      size="small"
                      v-model="record.status"
                      :checked-value="0"
                      :unchecked-value="1"
                      v-permission="['upStatus']"
                      @change="handleStatus(record)"
                    />
                    <a-popconfirm content="确定删除该经纪人吗?" @ok="handleDel(record)" position="tr" v-permission="['del']">
                      <a-button type="text" size="small" status="danger">
                        <template #icon><icon-delete /></template>
                      </a-button>
                    </a-popconfirm>
                  </a-space>
                </div>
              </a-card>
            </a-grid-item>
          </a-grid>

          <a-empty v-else-if="!loading" description="暂无数据" />
        </a-spin>
      </a-card>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="pagination.total > 0">
        <a-pagination
          :total="pagination.total"
          :current="pagination.current"
          :page-size="pagination.pageSize"
          show-total
          show-page-size
          @change="handlePageChange"
          @page-size-change="handlePageSizeChange"
        />
      </div>

      <AddForm @register="registerModal" @success="handleData" />
    </page-card>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, onMounted } from 'vue';
import useLoading from '@/hooks/loading';
import { Message } from '@arco-design/web-vue';
import { Pagination } from '@/types/global';
import { useModal } from '/@/components/Modal';
import { GetFullPath } from '@/utils/tool';

import AddForm from './modal/AddForm.vue';
import { getList, upStatus, del, auditSwitch, upCanManageProperties, upCanManageLocks } from './api';

const [registerModal, { openModal }] = useModal();

const basePagination: Pagination = { current: 1, pageSize: 10 };
const pagination = reactive({
  ...basePagination,
  total: 0,
  showTotal: true,
  showPageSize: true,
});

const formModel = ref({
  name: '',
  mobile: '',
  username: '',
  audit_status: '',
  status: '',
  can_manage_properties: '',
  can_manage_locks: '',
});

const { loading, setLoading } = useLoading(true);
const renderData = ref<any[]>([]);

const getAvatarText = (record: any) => {
  const name = String(record?.name || record?.nickname || '').trim();
  if (name) return name.slice(-2);
  const mobile = String(record?.mobile || '').trim();
  if (mobile) return mobile.slice(-2);
  const username = String(record?.username || '').trim();
  if (username) return username.slice(0, 2);
  return '经纪';
};

const auditStatusText = (status: any) => {
  const s = String(status || '').trim().toLowerCase();
  if (s === 'approved') return '已通过';
  if (s === 'pending') return '待审核';
  if (s === 'rejected') return '已拒绝';
  return '未设置';
};

const auditStatusColor = (status: any) => {
  const s = String(status || '').trim().toLowerCase();
  if (s === 'approved') return 'green';
  if (s === 'pending') return 'orange';
  if (s === 'rejected') return 'red';
  return 'gray';
};

const normalizeAuditStatus = (status: any) => String(status || '').trim().toLowerCase();
const isAuditApproved = (record: any) => normalizeAuditStatus(record?.audit_status) === 'approved';

const toggleLoading = ref<Record<string, boolean>>({});
const toggleLoadingKey = (record: any, key: string) => `${record?.id || 0}:${key}`;
const isToggleLoading = (record: any, key: string) => !!toggleLoading.value[toggleLoadingKey(record, key)];
const setToggleLoading = (record: any, key: string, val: boolean) => {
  toggleLoading.value[toggleLoadingKey(record, key)] = val;
};

const handleAuditSwitch = async (record: any, checked: boolean) => {
  if (!record?.id) return;
  const status = checked ? 'approved' : 'pending';
  setToggleLoading(record, 'audit', true);
  try {
    await auditSwitch({ id: record.id, audit_status: status });
    record.audit_status = status;
    Message.success('更新成功');
  } catch (e: any) {
    Message.error(e?.message || '更新失败');
    fetchData();
  } finally {
    setToggleLoading(record, 'audit', false);
  }
};

const handleCanManageProperties = async (record: any, val: number) => {
  if (!record?.id) return;
  const oldVal = Number(val) === 1 ? 0 : 1;
  setToggleLoading(record, 'canManageProperties', true);
  try {
    await upCanManageProperties({ id: record.id, can_manage_properties: val });
    Message.success('更新成功');
  } catch (e: any) {
    Message.error(e?.message || '更新失败');
    record.can_manage_properties = oldVal;
  } finally {
    setToggleLoading(record, 'canManageProperties', false);
  }
};

const handleCanManageLocks = async (record: any, val: number) => {
  if (!record?.id) return;
  const oldVal = Number(val) === 1 ? 0 : 1;
  setToggleLoading(record, 'canManageLocks', true);
  try {
    await upCanManageLocks({ id: record.id, can_manage_locks: val });
    Message.success('更新成功');
  } catch (e: any) {
    Message.error(e?.message || '更新失败');
    record.can_manage_locks = oldVal;
  } finally {
    setToggleLoading(record, 'canManageLocks', false);
  }
};

const fetchData = async () => {
  setLoading(true);
  try {
    const resp: any = await getList({ page: pagination.current, pageSize: pagination.pageSize, ...formModel.value });
    // 兼容返回形态：
    // 1) {items,page,...}
    // 2) {code,data:{items,...}}
    // 3) AxiosResponse: {data:{code,data:{items,...}}}
    const data =
      resp?.items
        ? resp
        : resp?.data?.items
          ? resp.data
          : resp?.data?.data?.items
            ? resp.data.data
            : resp?.data ?? resp;

    const items = Array.isArray(data?.items) ? data.items : [];
    renderData.value = items.map((it: any) => ({
      ...it,
      status: Number(it?.status ?? 0),
      can_manage_properties: Number(it?.can_manage_properties ?? 0),
      can_manage_locks: Number(it?.can_manage_locks ?? 0),
    }));
    pagination.total = Number(data?.total) || 0;
    pagination.current = Number(data?.page) || pagination.current;
  } catch (e: any) {
    renderData.value = [];
    pagination.total = 0;
    Message.error(e?.message || '获取列表失败');
  } finally {
    setLoading(false);
  }
};

onMounted(() => {
  fetchData();
});

const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

const handleReset = () => {
  pagination.current = 1;
  formModel.value = { name: '', mobile: '', username: '', audit_status: '', status: '', can_manage_properties: '', can_manage_locks: '' };
  fetchData();
};

const handleCreate = () => {
  openModal(true, { isUpdate: false, record: null });
};

const handleEdit = (record: any) => {
  openModal(true, { isUpdate: true, record });
};

const handleData = () => {
  pagination.current = 1;
  fetchData();
};

const handleStatus = async (record: any) => {
  try {
    await upStatus({ id: record.id, status: record.status });
    Message.success('更新成功');
  } catch (e: any) {
    Message.error(e?.message || '更新失败');
    fetchData();
  }
};

const handleDel = async (record: any) => {
  try {
    await del({ ids: [record.id] });
    Message.success('删除成功');
    fetchData();
  } catch (e: any) {
    Message.error(e?.message || '删除失败');
  }
};

const handlePageChange = (page: number) => {
  pagination.current = page;
  fetchData();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  fetchData();
};
</script>
<script lang="ts">
  export default {
    name: 'broker', // If you want the include property of keep-alive to take effect, you must name the component
  };
</script>

<style lang="less" scoped>
.broker-hero {
  margin-bottom: 12px;
  border-radius: 12px;
  background: var(--color-bg-2);
  border: 1px solid rgb(var(--gray-2));
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}
.hero-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.hero-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.hero-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(var(--primary-6), 0.12);
  color: rgb(var(--primary-6));
}
.hero-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--color-text-1);
}
.hero-sub {
  margin-top: 2px;
  font-size: 12px;
  color: var(--color-text-3);
}
.broker-filter {
  margin-bottom: 12px;
  border-radius: 12px;
  border: 1px solid rgb(var(--gray-2));
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}
.search-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}
.broker-list {
  border-radius: 12px;
  border: 1px solid rgb(var(--gray-2));
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0 4px;
}

.broker-grid {
  width: 100%;
}

.broker-card {
  border-radius: 12px;
  border: 1px solid rgb(var(--gray-2));
  background: var(--color-fill-1);
  transition: box-shadow 0.2s ease, transform 0.2s ease, border-color 0.2s ease;
}
.broker-card:hover {
  border-color: rgb(var(--arcoblue-3));
  box-shadow: 0 10px 26px rgba(0, 0, 0, 0.08);
  transform: translateY(-1px);
}

.card-head {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}
.avatar {
  flex: 0 0 auto;
  background: rgba(var(--arcoblue-6), 0.12);
  color: rgb(var(--arcoblue-6));
  font-weight: 700;
}
.head-main {
  flex: 1;
  min-width: 0;
}
.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.name {
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-text-1);
  font-size: 15px;
  font-weight: 700;
  line-height: 22px;
}
.sub-row {
  margin-top: 4px;
  color: var(--color-text-3);
  font-size: 12px;
  line-height: 18px;
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}
.store {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.sep {
  color: rgb(var(--gray-5));
}

.info {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 66px;
}
.info-line {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--color-text-2);
  font-size: 12px;
  line-height: 18px;
}
.audit-line .label {
  color: var(--color-text-3);
}
.info-line :deep(.arco-icon) {
  font-size: 14px;
  color: var(--color-text-3);
}
.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tags {
  margin-top: 10px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.tag-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.actions {
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px dashed rgb(var(--gray-2));
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
