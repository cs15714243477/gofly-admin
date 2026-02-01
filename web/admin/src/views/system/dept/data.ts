import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  export const columns = computed<TableColumnData[]>(() => [
    {
      title: '部门名称',
      dataIndex: 'name',
      slotName: 'title',
      width: 230,
    },
    {
      title: '排序',
      dataIndex: 'weigh',
      width: 100,
      align:"center"
    },
    {
      title: '状态',
      dataIndex: 'status',
      slotName: 'status',
      width: 100,
      align:"center"
    },
    {
      title: '创建时间',
      dataIndex: 'createtime',
      width: 180,
      align:"center"
    },
    {
      title: '备注',
      dataIndex: 'remark',
      width: 220,
    },
    {
      title: '操作',
      dataIndex: 'operations',
      slotName: 'operations',
      fixed: 'right',
      width: 100,
      align:"center"
    },
  ]);
  export interface TreeItem  {
    id?: number;
    pid?: number;
    title?: string;
  }