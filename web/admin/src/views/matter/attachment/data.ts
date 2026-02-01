import { computed } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  export const columns = computed<TableColumnData[]>(() => [
    {
      title: 'ID',
      dataIndex: 'id',
      align:"center",
      width:70,
    },
    {
      title:  '图片',
      dataIndex: 'mimetype',
      slotName: 'mimetype',
      width:90,
      align:'center'
    },
    {
      title: '文件名称',
      dataIndex: 'title',
      width:230,
      align:"center"
    },
    {
      title: '文件大小',
      dataIndex: 'filesize',
      slotName: 'filesize',
      width:90,
      align:"center"
    },
    {
      title: 'mime类型',
      dataIndex: 'mimetype',
      width:120,
      align:"center"
    },
    {
      title: '创建时间',
      dataIndex: 'createtime',
      slotName: 'createtime',
      width:190,
      align:"center"
    },
    {
      title: '操作',
      dataIndex: 'operations',
      slotName: 'operations',
      fixed: 'right',
      width:90,
      align:"center"
    },
  ]);