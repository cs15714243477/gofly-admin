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
              <a-button type="primary" @click="handleCreate">
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
      <a-card class="broker-table" :bordered="false">
        <a-spin :loading="loading" style="width: 100%">
          <a-table
            row-key="id"
            :data="renderData"
            :columns="columns"
            :pagination="false"
            :bordered="false"
            :stripe="false"
            :scroll="{ x: '100%' }"
          >
            <template #status="{ record }">
              <a-badge :status="record.status === 0 ? 'success' : 'default'" :text="record.status === 0 ? '启用' : '禁用'" />
            </template>
            <template #canManage="{ record }">
              <a-tag v-if="Number(record.can_manage_properties) === 1" color="green">可维护</a-tag>
              <a-tag v-else color="gray">不可维护</a-tag>
            </template>
            <template #canLock="{ record }">
              <a-tag v-if="Number(record.can_manage_locks) === 1" color="arcoblue">可管理</a-tag>
              <a-tag v-else color="gray">不可管理</a-tag>
            </template>
            <template #action="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  <template #icon><icon-edit /></template>
                  编辑
                </a-button>
                <a-switch
                  type="round"
                  size="small"
                  v-model="record.status"
                  :checked-value="0"
                  :unchecked-value="1"
                  @change="handleStatus(record)"
                />
                <a-popconfirm content="确定删除该经纪人吗?" @ok="handleDel(record)" position="tr">
                  <a-button type="text" size="small" status="danger">
                    <template #icon><icon-delete /></template>
                    删除
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table>

          <a-empty v-if="!loading && renderData.length === 0" description="暂无数据" />
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

import AddForm from './modal/AddForm.vue';
import { getList, upStatus, del } from './api';

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
  status: '',
  can_manage_properties: '',
  can_manage_locks: '',
});

const { loading, setLoading } = useLoading(true);
const renderData = ref<any[]>([]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80, align: 'center' },
  { title: '姓名', dataIndex: 'name', width: 120, ellipsis: true },
  { title: '昵称', dataIndex: 'nickname', width: 120, ellipsis: true },
  { title: '手机号', dataIndex: 'mobile', width: 140, ellipsis: true },
  { title: '邮箱', dataIndex: 'email', minWidth: 220, ellipsis: true },
  { title: '门店名称', dataIndex: 'store_name', minWidth: 180, ellipsis: true },
  { title: '门店地址', dataIndex: 'store_address', minWidth: 260, ellipsis: true },
  { title: '门店电话', dataIndex: 'store_contact_phone', width: 140, ellipsis: true },
  { title: '店长', dataIndex: 'store_manager_name', width: 120, ellipsis: true },
  { title: '头衔', dataIndex: 'title', minWidth: 140, ellipsis: true },
  { title: '可维护房源', dataIndex: 'can_manage_properties', width: 120, align: 'center', slotName: 'canManage' },
  { title: '可管智能锁', dataIndex: 'can_manage_locks', width: 120, align: 'center', slotName: 'canLock' },
  { title: '状态', dataIndex: 'status', width: 110, align: 'center', slotName: 'status' },
  { title: '操作', dataIndex: 'action', width: 240, fixed: 'right', align: 'center', slotName: 'action' },
] as any[];

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

    renderData.value = Array.isArray(data?.items) ? data.items : [];
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
  formModel.value = { name: '', mobile: '', username: '', status: '', can_manage_properties: '', can_manage_locks: '' };
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
  background: linear-gradient(135deg, var(--color-bg-2), rgba(var(--primary-1), 0.6));
  border: 1px solid var(--color-border-2);
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
  border: 1px solid var(--color-border-2);
}
.search-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}
.broker-table {
  border-radius: 12px;
  border: 1px solid var(--color-border-2);
}
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0 4px;
}

:deep(.arco-table-th) {
  white-space: nowrap;
}

:deep(.arco-table-td) {
  white-space: nowrap;
}
</style>
