import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
export const columns = computed<TableColumnData[]>(() => [
     {
       title:  "城市名称",
       dataIndex: "name",
       align:"left",
       width: 250,
     },
     {
      title:  "ID",
      dataIndex: "id",
      align:"left",
      width: 90,
    },
     {
       title:  "层级",
       dataIndex: "level",
       slotName:  "level",
       align:"left",
       width: 100,
     },
     {
       title:  "邮编",
       dataIndex: "zip",
       align:"left",
       width: 110,
     },
     {
       title:  "经度",
       dataIndex: "lng",
       align:"left",
       width: 110,
     },
     {
       title:  "纬度",
       dataIndex: "lat",
       align:"left",
       width: 110,
     },
    {
      title: '操作',
      dataIndex: 'operations',
      slotName: 'operations',
      fixed: 'right',
      width: 150,
      align:"center"
    },
  ]);
