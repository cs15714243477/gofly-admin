<template>
  <div class="container">
    <page-card breadcrumb scrollPage>
      <!-- 搜索栏 -->
      <div class="search-row">
        <a-space wrap>
          <a-input style="width: 220px" v-model="formModel.name" placeholder="模板名称" allow-clear />
          <a-select style="width: 170px" v-model="formModel.template_type" placeholder="模板类型" allow-clear>
            <a-option value="property">房源</a-option>
            <a-option value="agent">经纪人</a-option>
            <a-option value="festive">节日</a-option>
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
        <div class="search-actions">
          <tabletool :showbtn="['create', 'refresh']" @create="handleCreate" @refresh="fetchData" />
        </div>
      </div>

      <!-- 列表 -->
      <a-spin :loading="loading" style="width: 100%">
        <a-table
          row-key="id"
          :data="renderData"
          :columns="columns"
          :pagination="false"
          :bordered="false"
          :stripe="true"
          :scroll="{ x: '100%' }"
        >
          <template #preview="{ record }">
            <a-image
              v-if="record.preview_url"
              :src="getImageUrl(record.preview_url)"
              width="48"
              height="48"
              fit="cover"
              :preview="true"
              class="img-thumb"
            />
            <span v-else class="muted">无</span>
          </template>

          <template #type="{ record }">
            <a-tag v-if="record.template_type === 'property'" color="arcoblue">房源</a-tag>
            <a-tag v-else-if="record.template_type === 'agent'" color="green">经纪人</a-tag>
            <a-tag v-else-if="record.template_type === 'festive'" color="orange">节日</a-tag>
            <a-tag v-else>{{ record.template_type || '-' }}</a-tag>
          </template>

          <template #status="{ record }">
            <a-switch
              type="round"
              size="small"
              v-model="record.status"
              :checked-value="0"
              :unchecked-value="1"
              @change="() => handleStatus(record)"
            />
          </template>

          <template #action="{ record }">
            <a-space>
              <a-button type="text" size="small" @click="handleEdit(record)">
                <template #icon><icon-edit /></template>
                编辑
              </a-button>
              <a-popconfirm content="确定删除该模板吗?" @ok="handleDel(record)" position="tr">
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
import { tabletool } from '/@/components/tabletool';
import { useModal } from '/@/components/Modal';
import { GetFullPath } from '@/utils/tool';

import AddForm from './AddForm.vue';
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
  template_type: '',
  status: '',
});

const { loading, setLoading } = useLoading(true);
const renderData = ref<any[]>([]);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80, align: 'center' },
  { title: '预览', dataIndex: 'preview_url', width: 90, align: 'center', slotName: 'preview' },
  { title: '模板名称', dataIndex: 'name', minWidth: 220, ellipsis: true },
  { title: '类型', dataIndex: 'template_type', width: 120, align: 'center', slotName: 'type' },
  { title: '权重', dataIndex: 'weigh', width: 100, align: 'center' },
  { title: '状态', dataIndex: 'status', width: 110, align: 'center', slotName: 'status' },
  { title: '创建时间', dataIndex: 'createtime', width: 180, ellipsis: true },
  { title: '操作', dataIndex: 'action', width: 180, fixed: 'right', align: 'center', slotName: 'action' },
] as any[];

const getImageUrl = (url: string) => {
  if (!url) return '';
  return GetFullPath(url) || url;
};

const fetchData = async () => {
  setLoading(true);
  try {
    const resp: any = await getList({ page: pagination.current, pageSize: pagination.pageSize, ...formModel.value });
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
  formModel.value = { name: '', template_type: '', status: '' };
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
    name: 'poster', // If you want the include property of keep-alive to take effect, you must name the component
  };
</script>

<style lang="less" scoped>
.search-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}
.search-actions {
  display: flex;
  justify-content: flex-end;
}
.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0 4px;
}
.img-thumb {
  border-radius: 8px;
  overflow: hidden;
}
.muted {
  color: var(--color-text-3);
}

::deep(.arco-table-th) {
  white-space: nowrap;
}
::deep(.arco-table-td) {
  white-space: nowrap;
}
</style>
