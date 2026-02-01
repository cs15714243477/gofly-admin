import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
export const columns = computed<TableColumnData[]>(() => [
     {
       title:  'ID',
       dataIndex: 'id',
       width: 80,
       align:'center'
     },
     {
      title:  '图片',
      dataIndex: 'url',
      slotName: 'image',
      width: 90,
      align:'center'
    },
     {
       title:  '文件名称',
       dataIndex: 'title',
       slotName: 'des',
       width: 280,
       align:'center'
     },
     {
      title:  '类型',
      dataIndex: 'type',
      slotName: 'type',
      width: 100,
      align:'center'
     },
     {
       title:  '分类',
       dataIndex: 'catename',
       width: 100,
       align:'center'
     },
     {
       title:  '大小',
       dataIndex: 'filesize',
       slotName: 'filesize',
       width: 100,
       align:'center'
     },
     {
       title:  '状态',
       dataIndex: 'status',
       slotName: 'status',
       width: 90,
       align:'center'
     },
     {
      title: '创建时间',
      dataIndex: 'createtime',
      slotName: 'createtime',
      width: 150,
      align:"center"
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
