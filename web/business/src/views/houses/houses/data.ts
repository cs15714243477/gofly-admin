import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';

export const saleStatusDict = [
  { value: 'on_sale', label: '在售', color: 'success' },
  { value: 'sold', label: '已售', color: 'default' },
  { value: 'off_market', label: '下架', color: 'warning' },
];

export const columns = computed<TableColumnData[]>(() => [
  {
    title: '封面',
    dataIndex: 'cover_image',
    slotName: 'cover_image',
    align: 'left',
    width: 90,
  },
  {
    title: '标题',
    dataIndex: 'title',
    align: 'left',
    width: 260,
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '价格',
    dataIndex: 'price',
    slotName: 'price',
    align: 'left',
    width: 120,
  },
  {
    title: '面积(㎡)',
    dataIndex: 'area',
    align: 'left',
    width: 100,
  },
  {
    title: '户型',
    dataIndex: 'layout',
    slotName: 'layout',
    align: 'left',
    width: 120,
  },
  {
    title: '小区',
    dataIndex: 'community_name',
    align: 'left',
    width: 160,
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '地址',
    dataIndex: 'address',
    align: 'left',
    width: 220,
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '标签',
    dataIndex: 'tags',
    slotName: 'tags',
    align: 'left',
    width: 180,
  },
  {
    title: '销售状态',
    dataIndex: 'sale_status',
    slotName: 'sale_status',
    align: 'left',
    width: 110,
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
    align: 'center',
    width: 90,
  },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    fixed: 'right',
    width: 160,
    align: 'center',
  },
]);

