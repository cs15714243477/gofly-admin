import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
export const columns = computed<TableColumnData[]>(() => [
     {
       title:  "ID",
       dataIndex: "id",
       align:"left",
       width: 80,
     },
     {
       title:  "门店名称",
       dataIndex: "name",
       align:"left",
       width: 220,
     },
     {
       title:  "门店地址",
       dataIndex: "address",
       slotName:  "cellcopy",
       align:"left",
       width: 320,
     },
     {
       title:  "联系电话",
       dataIndex: "contact_phone",
       slotName:  "cellcopy",
       align:"left",
       width: 160,
     },
     {
       title:  "店长姓名",
       dataIndex: "manager_name",
       slotName:  "cellcopy",
       align:"left",
       width: 140,
     },
     {
       title:  "图片",
       dataIndex: "images",
       slotName:  "images",
       align:"center",
       width: 120,
     },
     {
       title:  "状态",
       dataIndex: "status",
       slotName:  "switchstatus",
       align:"center",
       width: 100,
     },
     {
       title:  "排序权重",
       dataIndex: "weigh",
       align:"center",
       width: 100,
     },
     {
       title:  "创建时间",
       dataIndex: "createtime",
       align:"left",
       width: 170,
     },
    {
      title: '操作',
      dataIndex: 'operations',
      slotName: 'operations',
      fixed: 'right',
      width: 120,
      align:"center"
    },
  ]);
