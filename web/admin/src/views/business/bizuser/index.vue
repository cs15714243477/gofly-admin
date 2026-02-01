<template>
  <div class="container" >
    <page-card breadcrumb fullTable :isFullscreen="isFullscreen">
      <template #searchLeft>
          <a-space>
            <a-input :style="{width:'200px'}"  v-model="formModel.name" placeholder="搜索用户名" allow-clear />
            <a-range-picker v-model="formModel.createtime" :shortcuts="shortcuts" shortcuts-position="left" :style="{width:'230px'}" @change="fetchData"/>
            <a-select v-model="formModel.status"  :options="statusOptions" placeholder="状态" :style="{width:'80px'}" @change="fetchData"/>
            <a-button type="primary" @click="search">
              <template #icon>
                <icon-search />
              </template>
              {{ $t("searchTable.form.search") }}
            </a-button>
          </a-space>
        </template>
        <template #searchRight>
        <tabletool :showbtn="['create','refresh','selectdensity','fullscreen']"
           @create="createData" @refresh="fetchData" @selectdensity="(data)=>size=data" @fullscreen="(data)=>isFullscreen=data"></tabletool>
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
          :scroll="{x:'100%', y: '100%'}"
          @page-change="handlePaageChange" 
          @page-size-change="handlePaageSizeChange" 
        >
          <template #name="{ record }">
            <CellNodeAvatar
                :avatar="record.avatar"
                :name="record.name"
              />
          </template>
          <template #createtime="{ record }">
            {{dayjs(record.createtime).format("YYYY-MM-DD")}}
          </template>
          <template #status="{ record }">
            <a-switch type="round" v-model="record.status" :checked-value="0" :unchecked-value="1" @change="handleStatus(record)">
                <template #checked>
                  开
                </template>
                <template #unchecked>
                  关
                </template>
              </a-switch>
          </template>
          <template #operations="{ record }">
            <Icon icon="svgfont-edit-o" class="iconbtn" @click="handleEdit(record)" :size="18" color="#0960bd"></Icon>
            <a-divider direction="vertical" />
            <a-popconfirm content="您确定要删除吗?" @ok="handleDel(record)">
              <Icon icon="svgfont-delete-o" class="iconbtn" :size="18" color="#ed6f6f"></Icon>
            </a-popconfirm>
          </template>
        </a-table>
      </template>
    </page-card>
    <!--表单-->
    <AddForm @register="registerModal"  @success="handleData"/>
  </div>
</template>

<script lang="ts" setup>
  import {  ref, reactive, watch, nextTick } from 'vue';
  import useLoading from '@/hooks/loading';
  import { getList,upStatus,del} from '@/api/business/bizuser';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  import { shortcuts, dayjs} from '@/utils/dayjs';
  //数据
  import { columns} from './data';
  //表单
  import { useModal } from '/@/components/Modal';
  import AddForm from './AddForm.vue';
  import {Icon} from '@/components/Icon';
  import { Message } from '@arco-design/web-vue';
  import { Pagination } from '@/types/global';
  import {tabletool,SizeProps,statusOptions} from '/@/components/tabletool';
  const [registerModal, { openModal }] = useModal();
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
  // const boxheight=document.documentElement.clientHeight;//页面高度
  type Column = TableColumnData & { checked?: true };
  const { loading, setLoading } = useLoading(true);
  const renderData = ref([]);
  const cloneColumns = ref<Column[]>([]);
  const showColumns = ref<Column[]>([]);
  const isFullscreen = ref<boolean>(false);
  const size = ref<SizeProps>('large');
  const artable=ref()
  //查询字段
  const generateFormModel = () => {
    return {
      name: '',
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
      nextTick(()=>{
        artable.value.expandAll()
      })
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      setLoading(false);
    }
  };
  const search = () => {
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
  //添加
  const createData=()=>{
    openModal(true, {
      isUpdate: false,
      record:null
    });
  }
  //编辑数据
  const handleEdit=async(record:any)=>{
    openModal(true, {
      isUpdate: true,
      record:record
    });
  }
  //更新数据
  const handleData=async()=>{
    fetchData();
  }
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
  //更新状态
  const handleStatus=async(record:any)=>{
    try {
      Message.loading({content:"更新状态中",id:"upStatus",duration:0})
      await upStatus({id:record.id,status:record.status});
      Message.success({content:"更新状态成功",id:"upStatus",duration:2000})
    }catch (error) {
      record.status=record.status==1?0:1
      Message.loading({content:"更新状态中",id:"upStatus",duration:1})
    } 
  }
  //删除数据
  const handleDel=async(record:any)=>{
    try {
      Message.loading({content:"删除中",id:"del",duration:0})
      await del({ids:[record.id]});
      fetchData();
      Message.success({content:"删除成功",id:"del",duration:2000})
    }catch (error) {
      Message.loading({content:"删除中",id:"del",duration:1})
    } 
  }
</script>

<script lang="ts">
  export default {
    name: 'bizuser',
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
</style>
