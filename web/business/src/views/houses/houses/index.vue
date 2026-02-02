<template>
  <div class="page-container">
    <div class="content-wrapper">
      <!-- 头部状态区：采用横向卡片布局，突出关键指标 -->
      <div class="dashboard-header">
        <a-row :gutter="24">
          <a-col :span="8">
            <div class="stat-card">
              <div class="stat-icon-wrapper house-bg">
                <icon-home class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">房源总数</div>
                <div class="stat-value-row">
                  <span class="stat-value">{{ pagination.total || 0 }}</span>
                  <div class="status-dot bg-success"></div>
                </div>
                <div class="stat-sub">已录入房产信息</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card">
              <div class="stat-icon-wrapper sale-bg">
                <icon-check-circle class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">在售房源</div>
                <div class="stat-value-row">
                  <span class="stat-value">{{ getStatusCount('on_sale') }}</span>
                  <div class="status-dot bg-success"></div>
                </div>
                <div class="stat-sub">可供客户选择</div>
              </div>
            </div>
          </a-col>
          <a-col :span="8">
            <div class="stat-card info-card">
              <div class="stat-icon-wrapper sold-bg">
                <icon-trophy class="stat-icon" />
              </div>
              <div class="stat-info">
                <div class="stat-label">已售房源</div>
                <div class="stat-value-row">
                  <span class="stat-value">{{ getStatusCount('sold') }}</span>
                </div>
                <div class="stat-sub">成交完成</div>
              </div>
            </div>
          </a-col>
        </a-row>
      </div>

      <!-- 主体内容卡片 -->
      <a-card class="main-card" :bordered="false">
        <!-- 工具栏与搜索 -->
        <div class="card-toolbar">
          <div class="left-actions">
            <a-input-search
              v-model="formModel.keyword"
              placeholder="搜索房源标题、小区名称或地址..."
              style="width: 320px"
              allow-clear
              search-button
              @search="handleSearch"
              @press-enter="handleSearch"
            >
              <template #button-icon>
                <icon-search />
              </template>
            </a-input-search>
            <a-select
              v-model="formModel.sale_status"
              placeholder="全部状态"
              allow-clear
              style="width: 140px; margin-left: 12px"
              @change="handleSearch"
            >
              <a-option value="on_sale">在售</a-option>
              <a-option value="sold">已售</a-option>
              <a-option value="off_market">下架</a-option>
            </a-select>
          </div>
          <div class="right-actions">
            <a-space size="medium">
               <a-tooltip content="新增房源信息">
                <a-button type="primary" @click="handleCreate" class="action-btn">
                  <template #icon><icon-plus /></template>
                  新增房源
                </a-button>
              </a-tooltip>
              <a-divider direction="vertical" />
              <div class="icon-btn-group">
                 <a-tooltip content="刷新列表">
                  <a-button type="text" shape="circle" @click="fetchData">
                    <icon-refresh />
                  </a-button>
                </a-tooltip>
                <a-tooltip content="密度设置" v-if="!isFullscreen">
                   <a-dropdown @select="(val: any) => size = val">
                    <a-button type="text" shape="circle"><icon-menu /></a-button>
                    <template #content>
                      <a-doption value="mini">迷你</a-doption>
                      <a-doption value="medium">中等</a-doption>
                      <a-doption value="large">宽松</a-doption>
                    </template>
                  </a-dropdown>
                </a-tooltip>
              </div>
            </a-space>
          </div>
        </div>

        <!-- 数据表格 -->
        <div class="table-wrapper">
          <a-table
            row-key="id"
            :loading="loading"
            :pagination="pagination"
            :data="renderData"
            :columns="columns"
            :bordered="false"
            :size="size"
            :stripe="true"
            :hoverable="true"
            @page-change="handlePageChange"
            @page-size-change="handlePageSizeChange"
          >
            <template #empty>
              <div class="empty-state">
                <icon-home style="font-size: 48px; color: var(--color-neutral-4); margin-bottom: 16px;" />
                <div class="empty-text">暂无房源数据</div>
                <div class="empty-sub">请添加房源信息开始管理</div>
                <a-button type="primary" class="mt-4" @click="handleCreate">立即添加</a-button>
              </div>
            </template>
            <!-- 自定义列渲染 -->
            <template #cover_image="{ record }">
              <div class="cell-cover">
                <a-image
                  :src="getImageUrl(record.cover_image)"
                  width="64"
                  height="48"
                  fit="cover"
                  class="cover-img"
                />
                <div class="cover-overlay" v-if="record.sale_status === 'sold'">已售</div>
              </div>
            </template>

            <template #title="{ record }">
              <div class="cell-house-main">
                <div class="house-title">{{ record.title }}</div>
                <div class="house-location">
                  <icon-location class="location-icon" />
                  {{ record.community_name }}
                  <span v-if="record.floor_level" class="divider">|</span>
                  <span v-if="record.floor_level">{{ record.floor_level }}</span>
                </div>
              </div>
            </template>

            <template #price="{ record }">
              <div class="cell-price">
                <div class="price-main">{{ record.price }}<span class="unit">{{ record.price_unit || '万' }}</span></div>
                <div class="price-sub" v-if="record.area">≈{{ Math.round((Number(record.price) * 10000)/Number(record.area)) }}/㎡</div>
              </div>
            </template>

            <template #layout="{ record }">
              <div class="layout-tag">{{ record.rooms }}室{{ record.halls }}厅</div>
            </template>
            
            <template #tags="{ record }">
                <div class="tags-cell">
                   <a-space size="mini" wrap>
                     <a-tag
                       v-for="(tag, idx) in parseTags(record.tags).slice(0, 2)"
                       :key="idx"
                       size="small"
                       bordered
                       class="pill-tag"
                     >
                       {{ tag }}
                     </a-tag>
                     <a-tag v-if="parseTags(record.tags).length > 2" size="small" bordered class="pill-tag pill-tag-more">
                       +{{ parseTags(record.tags).length - 2 }}
                     </a-tag>
                   </a-space>
                </div>
            </template>

            <template #sale_status="{ record }">
              <div class="status-cell" :class="record.sale_status">
                <span class="status-dot"></span>
                {{ getSaleStatusLabel(record.sale_status) }}
              </div>
            </template>

            <template #status="{ record }">
               <a-switch
                  type="round"
                  size="small"
                  v-model="record.status"
                  :checked-value="0"
                  :unchecked-value="1"
                  @change="(val) => handleStatusChange(record, val)"
                />
            </template>

            <template #operations="{ record }">
              <div class="action-cell">
                 <a-button type="text" size="small" class="table-btn-primary" @click="handleDetail(record)">
                  <template #icon><icon-eye /></template>
                  详情
                </a-button>
                <a-button type="text" size="small" class="table-btn-default" @click="handleEdit(record)">
                  <template #icon><icon-edit /></template>
                </a-button>
                <a-popconfirm content="确定删除?" position="left" @ok="handleDelete(record)">
                  <a-button type="text" size="small" class="table-btn-danger">
                    <template #icon><icon-delete /></template>
                  </a-button>
                </a-popconfirm>
              </div>
            </template>
          </a-table>
        </div>
      </a-card>
    </div>

    <!-- Forms & Modals -->
    <AddForm @register="registerAddModal" @success="handleSuccess" />
    <DetailView @register="registerDetailModal" />
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue';
import useLoading from '@/hooks/loading';
import { Message } from '@arco-design/web-vue';
import { Pagination } from '@/types/global';
import { SizeProps } from '/@/components/tabletool';
import { useModal } from '/@/components/Modal';
import { GetFullPath } from '@/utils/tool';
import AddForm from './AddForm.vue';
import DetailView from './DetailView.vue';
import { getList, upStatus, del } from './api';

const { loading, setLoading } = useLoading(true);
const isFullscreen = ref(false);

const getStatusCount = (status: string) => {
  return renderData.value.filter(item => item.sale_status === status).length;
};
const size = ref<SizeProps>('large');

const basePagination: Pagination = { current: 1, pageSize: 12 }; // More items per page
const pagination = reactive({
  ...basePagination,
  total: 0,
  showTotal: true,
  showPageSize: true,
});

const formModel = reactive({
  keyword: '',
  sale_status: '',
});

const columns = [
  { title: '封面', slotName: 'cover_image', width: 90, align: 'center' },
  { title: '房源信息', slotName: 'title', width: 320 },
  { title: '价格行情', slotName: 'price', width: 140, sortable: { sortDirections: ['ascend', 'descend'] } },
  { title: '户型', slotName: 'layout', width: 140, align: 'center' },
  { title: '面积', dataIndex: 'area', width: 100, render: ({ record }:any) => `${record.area}㎡` },
  { title: '装修', dataIndex: 'decoration_type', width: 90, render: ({ record }:any) => record.decoration_type || '-' },
  { title: '标签', slotName: 'tags', width: 180 },
  { title: '状态', slotName: 'sale_status', width: 110 },
  { title: '上架', slotName: 'status', width: 80, align: 'center' },
  { title: '操作', slotName: 'operations', width: 180, fixed: 'right', align: 'center' },
];

const renderData = ref<any[]>([]);

const fetchData = async () => {
  setLoading(true);
  try {
    const params = {
      page: pagination.current,
      pageSize: pagination.pageSize,
      title: formModel.keyword,
      sale_status: formModel.sale_status,
    };
    const resp: any = await getList(params);
    const data = resp?.items ? resp : (resp?.data?.items ? resp.data : resp?.data ?? resp);
    
    renderData.value = (data?.items || []) as any[];
    pagination.current = Number(data?.page || pagination.current);
    pagination.total = Number(data?.total || 0);
  } catch (e: any) {
    renderData.value = [];
    Message.error(e?.message || '加载列表失败');
  } finally {
    setLoading(false);
  }
};

const handleSearch = () => {
  pagination.current = 1;
  fetchData();
};

const handlePageChange = (p: number) => {
  pagination.current = p;
  fetchData();
};

const handlePageSizeChange = (ps: number) => {
  pagination.pageSize = ps;
  handleSearch();
};

const getImageUrl = (url: string) => {
   if (!url) return '';
   return GetFullPath(url);
}

const parseTags = (tags: any) => {
  if (!tags) return [];
  if (Array.isArray(tags)) return tags;
  if (typeof tags === 'string') return tags.split(',').filter(Boolean);
  return [];
};

const getSaleStatusLabel = (status: string) => {
  const map: any = { on_sale: '在售', sold: '已售', off_market: '下架' };
  return map[status] || status;
};

// Actions
const [registerAddModal, { openModal: openAddModal }] = useModal();
const [registerDetailModal, { openModal: openDetailModal }] = useModal();

const handleCreate = () => {
  openAddModal(true, { isUpdate: false });
};

const handleEdit = (record: any) => {
  openAddModal(true, { isUpdate: true, record });
};

const handleDetail = (record: any) => {
    openDetailModal(true, { record });
}

const handleSuccess = () => {
  fetchData();
};

const handleStatusChange = async (record: any, val: any) => {
    try {
        await upStatus({ id: record.id, status: val });
        Message.success('状态更新成功');
    } catch(e) {
        // revert
        record.status = val === 0 ? 1 : 0;
    }
}

const handleDelete = async (record: any) => {
    try {
        await del({ ids: [record.id] });
        Message.success('删除成功');
        fetchData();
    } catch(e) {
        // error handled by http
    }
}

onMounted(() => {
  fetchData();
});
</script>

<style scoped lang="less">
// page layout
.page-container {
  padding: 20px;
  background-color: var(--color-fill-1);
  min-height: 100vh;
}

.content-wrapper {
  max-width: 1600px;
  margin: 0 auto;
}

// Header Stats Cards
.dashboard-header {
  margin-bottom: 24px;
}

.stat-card {
  background: var(--color-bg-2);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  transition: all 0.3s ease;
  border: 1px solid transparent;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.08);
  }

  &.info-card {
     border-color: rgba(var(--arcoblue-6), 0.3);
      background: linear-gradient(135deg, var(--color-bg-2) 0%, rgba(var(--arcoblue-6), 0.05) 100%);
  }

  .stat-icon-wrapper {
    width: 56px;
    height: 56px;
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    
    &.house-bg { background-color: rgba(var(--primary-6), 0.1); color: rgb(var(--primary-6)); }
    &.sale-bg { background-color: rgba(var(--success-6), 0.1); color: rgb(var(--success-6)); }
    &.sold-bg { background-color: rgba(var(--arcoblue-6), 0.1); color: rgb(var(--arcoblue-6)); }
    
    .stat-icon {
      font-size: 28px;
    }
  }

  .stat-info {
    flex: 1;
    overflow: hidden;
    
    .stat-label {
      font-size: 13px;
      color: var(--color-text-3);
      margin-bottom: 4px;
    }
    
    .stat-value-row {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 4px;
      
      .stat-value {
        font-size: 18px;
        font-weight: 700;
        color: var(--color-text-1);
        line-height: 1.2;
      }
      
      .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        &.bg-success { background-color: rgb(var(--success-6)); box-shadow: 0 0 4px rgb(var(--success-6)); }
        &.bg-danger { background-color: rgb(var(--danger-6)); }
        &.bg-warning { background-color: rgb(var(--warning-6)); }
      }
    }
    
    .stat-sub {
      font-size: 12px;
      color: var(--color-text-4);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      opacity: 0.8;
    }
  }
}

// Main Card Styling
.main-card {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.02);
  transition: all 0.3s;
  
  :deep(.arco-card-body) {
    padding: 24px;
  }
}

.card-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  
  .left-actions {
    flex: 1;
    display: flex;
    align-items: center;
  }
  .right-actions {
    margin-left: 16px;
  }
  
  .action-btn {
    border-radius: 6px;
    font-weight: 500;
  }
  
  .icon-btn-group {
    display: flex;
    gap: 8px;
    color: var(--color-text-2);
  }
}

// Table Styling Overrides
.table-wrapper {
  :deep(.arco-table-container) {
    border-radius: 8px;
  }
  :deep(.arco-table-th) {
    background-color: var(--color-fill-1);
    font-weight: 600;
  }
}

// Cell Renderers
.cell-cover {
  position: relative;
  width: 64px;
  height: 48px;
  border-radius: 6px;
  overflow: hidden;
  border: 1px solid var(--color-border-2);
  
  .cover-img {
    display: block;
  }
  
  .cover-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0,0,0,0.5);
    color: #fff;
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
  }
}

.cell-house-main {
  .house-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--color-text-1);
    margin-bottom: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .house-location {
    font-size: 12px;
    color: var(--color-text-3);
    display: flex;
    align-items: center;
    
    .location-icon { 
      margin-right: 4px; 
      font-size: 12px; 
    }
    .divider { 
      margin: 0 8px; 
      color: var(--color-border-3); 
    }
  }
}

.cell-price {
  .price-main {
    color: rgb(var(--danger-6));
    font-weight: 700;
    font-size: 16px;
    line-height: 1;
    margin-bottom: 2px;
    
    .unit { 
      font-size: 12px; 
      font-weight: normal; 
      margin-left: 2px; 
    }
  }
  .price-sub {
     font-size: 12px;
     color: var(--color-text-4);
  }
}

.layout-tag {
  background: var(--color-fill-2);
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 13px;
  color: var(--color-text-2);
  font-weight: 500;
  display: inline-block;
}

.tags-cell {
  .pill-tag {
    border-radius: 999px;
    padding: 0 10px;
    border-color: var(--color-border-2);
    background: rgba(var(--primary-6), 0.08);
    color: rgb(var(--primary-6));

    :deep(.arco-tag-content) {
      max-width: 120px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    &.pill-tag-more {
      background: var(--color-fill-2);
      color: var(--color-text-2);
    }
  }
}

.status-cell {
   display: inline-flex;
   align-items: center;
   font-size: 13px;
   font-weight: 500;
   
   .status-dot { 
     width: 6px; 
     height: 6px; 
     border-radius: 50%; 
     margin-right: 6px; 
   }
   
   &.on_sale { 
     color: rgb(var(--success-6)); 
     .status-dot { 
       background: rgb(var(--success-6)); 
       box-shadow: 0 0 0 2px rgba(var(--success-6), 0.2); 
     } 
   }
   &.sold { 
     color: var(--color-text-3); 
     .status-dot { 
       background: var(--color-text-4); 
     } 
   }
   &.off_market { 
     color: rgb(var(--warning-6)); 
     .status-dot { 
       background: rgb(var(--warning-6)); 
     } 
   }
}

.action-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  
  .table-btn-primary {
    color: rgb(var(--primary-6));
    &:hover { background-color: var(--color-fill-2); }
  }
  .table-btn-default {
    color: var(--color-text-2);
    &:hover { color: var(--color-text-1); background-color: var(--color-fill-2); }
  }
  .table-btn-danger {
    color: rgb(var(--danger-6));
    &:hover { background-color: rgba(var(--danger-6), 0.1); }
  }
}

// Empty State
.empty-state {
  padding: 40px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  
  .empty-text {
    font-size: 16px;
    color: var(--color-text-1);
    font-weight: 500;
    margin-bottom: 4px;
  }
  .empty-sub {
    font-size: 13px;
    color: var(--color-text-3);
  }
  .mt-4 {
    margin-top: 16px;
  }
}
</style>
