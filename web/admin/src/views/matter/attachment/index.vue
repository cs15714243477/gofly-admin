<template>
  <div class="container" >
    <page-card breadcrumb fullTable :isFullscreen="isFullscreen">
      <template #searchLeft>
          <a-space>
            <a-input :style="{width:'200px'}"  v-model="formModel.title" placeholder="搜索文件名称" allow-clear  @change="fetchData"/>
            <a-range-picker v-model="formModel.createtime" :shortcuts="shortcuts" shortcuts-position="left" :style="{width:'230px'}" @change="fetchData"/>
            <a-button type="primary" @click="fetchData">
              <template #icon>
                <icon-search />
              </template>
              {{ $t("searchTable.form.search") }}
            </a-button>
            <a-button @click="reset">
              {{ $t('searchTable.form.reset') }}
            </a-button>
          </a-space>
      </template>
      <template #searchRight>
        <tabletool :showbtn="['refresh','selectdensity','fullscreen']"
          @refresh="fetchData" @selectdensity="(data)=>size=data" @fullscreen="(data)=>isFullscreen=data">
            <a-button type="outline" @click="handleDellMore" :disabled="selectIds.length==0">
              批量删除({{ selectIds.length }})
            </a-button>
        </tabletool>
      </template>
      <template #table>
        <a-table
          row-key="id"
          :loading="loading"
          :pagination="pagination"
          :columns="(cloneColumns as TableColumnData[])"
          :data="renderData"
          :bordered="{wrapper:true,cell:true}"
          :size="size"
          :default-expand-all-rows="true"
          ref="artable"
          @page-change="handlePaageChange" 
          @page-size-change="handlePaageSizeChange" 
          :row-selection="{
            type: 'checkbox',
            showCheckedAll: true
          }"
          :scroll="{x:'100%', y: '100%'}"
          v-model:selectedKeys="selectIds"
          @selection-change="handleSelect"
        >
          <template #mimetype="{ record }">
            <a-image
                v-if="record.mimetype.includes('image')||record.mimetype.includes('video')"
                height="50"
                width="50"
                :src="record.mimetype.includes('image')?record.url:record.cover_url"
              />
          </template>
          <template #createtime="{record,column}" >
            <template v-if="record[column.dataIndex]">
              {{dayjs(record[column.dataIndex]).format("YYYY-MM-DD HH:mm")}}
            </template>
          </template>
          <template #overtime="{ record }">
            {{record.overtime==0?"不限":record.overtime}}
          </template>
          <template #filesize="{ record }">
            {{filesizeFont(record.filesize)}}
          </template>
          <template #operations="{ record }">
            <ButtonGroup size="small">
              <a-popconfirm content="您确定要删除附件吗?" @ok="handleDel([record.id])">
                <a-button style="color: red;">删除</a-button>
              </a-popconfirm>
            </ButtonGroup>
          </template>
        </a-table>
      </template>
    </page-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, watch } from 'vue';
  import useLoading from '@/hooks/loading';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  //数据
  import { columns} from './data';
  //api
  import { getList,del} from '@/api/matter/attachment';
  import { Message ,Button,Modal} from '@arco-design/web-vue';
  import { Pagination } from '@/types/global';
  import { shortcuts, dayjs} from '@/utils/dayjs';
  import {tabletool,SizeProps} from '/@/components/tabletool';
  const isFullscreen = ref<boolean>(false);
  const ButtonGroup=Button.Group
  //分页
  const basePagination: Pagination = {
    current: 1,
    pageSize: 10,
  };
  const pagination = reactive({
    ...basePagination,
    showTotal:true,
    showPageSize:true,
  });
  type Column = TableColumnData & { checked?: true };
  const { loading, setLoading } = useLoading(true);
  const renderData = ref([]);
  const cloneColumns = ref<Column[]>([]);
  const showColumns = ref<Column[]>([]);
  const size = ref<SizeProps>('large');
   //查询字段
   const generateFormModel = () => {
    return {
      title: '',
      createtime: [],
      status: '',
    };
  };
  const formModel = ref(generateFormModel());
  const fetchData = async () => {
    setLoading(true);
    try {
      const data= await getList({page:pagination.current,pageSize:pagination.pageSize,...formModel.value});
      renderData.value = data.items;
      pagination.current = data.page;
      pagination.total = data.total;
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      setLoading(false);
    }
  };
  const reset = () => {
    formModel.value = generateFormModel();
    fetchData();
  };
  fetchData();

  watch(
    () => columns.value,
    (val) => {
      cloneColumns.value = cloneDeep(val);
      cloneColumns.value.forEach((item, index) => {
        item.checked = true;
      });
      showColumns.value = cloneDeep(cloneColumns.value);
    },
    { deep: true, immediate: true }
  );
  //分页
  const handlePaageChange = (page:any) => {
    pagination.current=page
    fetchData();
  }
  //分页总数
  const handlePaageSizeChange = (pageSize:any) => {
    pagination.pageSize=pageSize
    fetchData();
  }
  //删除数据
  const handleDel=async(disdata:any)=>{
    try {
      Message.loading({content:"删除中",id:"upStatus",duration:0})
      await del({ids:disdata});
      fetchData();
      Message.success({content:"删除成功",id:"upStatus",duration:2000})
    }catch (error) {
      Message.loading({content:"删除中",id:"upStatus",duration:1})
    } 
  }
  //存储单位换算
  const suffix = ['B', 'KB', 'MB', 'GB', 'TB'];
  const filesizeFont=(size:any)=>{
    const base = Math.floor(Math.log2(size) / 10);
    const index = clamp(base, 0, 4);
    return (size / 2 ** (10 * index)).toFixed(2) + suffix[index];
  }
  function clamp(v:any, min:any, max:any) {
  if (v < min) return min;
  if (v > max) return max;
  return v;
}
 //选择
 const selectIds=ref([])
  const handleSelect=(ids:any)=>{
    selectIds.value=ids
  }
  //批量删除
  const handleDellMore=()=>{
    Modal.warning({
      title: '你确定要删除吗？',
      content: '删除后将无法恢复，请检查后再操作。',
      titleAlign:"start",
      hideCancel:false,
      onOk:async()=>{
        try {
          Message.loading({content:"删除中",id:"del",duration:0})
          await del({ids:selectIds.value});
          selectIds.value=[]
          pagination.current=1
          fetchData();
          Message.success({content:"删除成功",id:"del",duration:2000})
        }catch (error) {
          Message.loading({content:"删除中",id:"v",duration:1})
        } 
       
      }
    });
  }
</script>

<script lang="ts">
  export default {
    name: 'matter',
  };
</script>

<style scoped lang="less">
  :deep(.arco-table-th) {
    &:last-child {
      .arco-table-th-item-title {
        margin-left: 16px;
      }
    }
  }
  .action-icon {
    margin-left: 12px;
    cursor: pointer;
  }
  .active {
    color: #0960bd;
    background-color: #e3f4fc;
  }
  .setting {
    display: flex;
    align-items: center;
    width: 200px;
    .title {
      margin-left: 12px;
      cursor: pointer;
    }
  }
  :deep(.general-card > .arco-card-header){
    padding: 10px 16px;
  }
  .iconbtn{
    user-select: none;
    cursor: pointer;
    opacity: .8;
    &:hover{
      opacity: 1;
    }
  }
</style>
