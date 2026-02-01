<template>
  <div class="container" >
    <page-card breadcrumb scrollPage :isFullscreen="isFullscreen">
      <div class="search-row" :class="{unhasesearchbtn:!showsearchdownbtn,hasesearchbtn:showsearchdownbtn}">
        <div class="search-field">
          <div class="search-line search-drown" :class="{searchdowncss:searchdown}">
            <div class="search-drown-box" ref="searchRef">
              <a-space wrap>
								<a-input style="width: 120px;" v-model="formModel.name" placeholder="输入名称" allow-clear />
								<a-select style="width: 120px;" v-model="formModel.level" placeholder="选择层级" allow-clear>
									<a-option :value="1" >省</a-option>
									<a-option :value="2" >市</a-option>
									<a-option :value="3" >区/县</a-option>
								</a-select>
								<a-input style="width: 120px;" v-model="formModel.zip" placeholder="输入邮编" allow-clear />
              </a-space>
            </div>
          </div>
        </div>
        <div class="search-btn">
          <a-space>
            <a-tooltip mini :content="searchdown?'收起查询条件':'展开查询条件'">
              <a-button v-if="showsearchdownbtn" @click="searchdown=!searchdown">
                <template #icon>
                  <icon-up v-if="searchdown"/>
                  <icon-down v-else/>
                </template>
              </a-button>
            </a-tooltip>
            <a-button type="primary" @click="handleSearch">
                <template #icon>
                  <icon-search />
                </template>
                查询
              </a-button>
              <a-button @click="handleReset">
                {{ $t('searchTable.form.reset') }}
              </a-button>
            </a-space>
        </div>
        <div class="search-option">
          <tabletool :showbtn="['refresh','selectdensity','fullscreen']"
               @refresh="fetchData" @selectdensity="(data)=>size=data" @fullscreen="(data)=>isFullscreen=data"></tabletool>
        </div>
      </div>
      <a-table
         row-key="id"
        :loading="loading"
        :pagination="false"
        :columns="(cloneColumns as TableColumnData[])"
        :data="renderData"
        :bordered="{wrapper:true,cell:true}"
        :size="size"
        :default-expand-all-rows="true"
        ref="artable"
        :scroll="{x:'100%'}"
        @page-change="handlePaageChange" 
        @page-size-change="handlePaageSizeChange" 
        :load-more="loadMore"
      >
        <template #level="{record,column}">
          <cell-tag :value="record[column.dataIndex]" :dict="[{value:1,label:'省',color:'primary'},{value:2,label:'市',color:'success'},{value:3,label:'区/县',color:'warning'}]"/>
        </template>
        <template #operations="{ record }">
          <a-link :hoverable="false" @click="handleAdd(record)">添加</a-link>
          <a-divider direction="vertical" />
          <Icon icon="svgfont-edit-o" class="iconbtn" @click="handleEdit(record)" :size="18" color="#0960bd"></Icon>
          <a-divider direction="vertical" />
          <a-popconfirm content="您确定要删除吗?" @ok="handleDel(record)">
            <Icon icon="svgfont-delete-o" class="iconbtn" :size="18" color="#ed6f6f"></Icon>
          </a-popconfirm>
        </template>
      </a-table>
    </page-card>
    <!--表单-->
    <AddForm @register="registerModal" @success="handleData"/>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, watch, onMounted } from 'vue';
  import useLoading from '@/hooks/loading';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  //api
  import { getList,getMoreList,upStatus,del} from './api';
  //数据
  import { columns} from './data';
  import dayjs from 'dayjs';
  //表单
  import { useModal } from '/@/components/Modal';
  import AddForm from './AddForm.vue';
  import {Icon} from '@/components/Icon';
  import { Message } from '@arco-design/web-vue';
  import { Pagination } from '@/types/global';
  import {tabletool,SizeProps} from '/@/components/tabletool';
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
  type Column = TableColumnData & { checked?: true };
  const { loading, setLoading } = useLoading(true);
  const renderData = ref([]);
  const cloneColumns = ref<Column[]>([]);
  const showColumns = ref<Column[]>([]);
  const isFullscreen = ref(false);
  const size = ref<SizeProps>('large');
   //查询字段
   const generateFormModel = () => {
    return {
			name:"",
			level:"",
			zip:"",

    };
  };
  const formModel = ref(generateFormModel());
  const fetchData = async () => {
    setLoading(true);
    searchdown.value=false
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
  //组件挂载完成后执行的函数
  const searchRef=ref(null)
  const searchdown = ref(false);
  const showsearchdownbtn=ref(false)
  onMounted(()=>{
    fetchData();
    if(searchRef.value){
      const searchboxheight=searchRef.value["offsetHeight"]
      if(searchboxheight>40){
        showsearchdownbtn.value=true
      }else{
        showsearchdownbtn.value=false
      }
    }
  })
 //查找
 const handleSearch = () => {
    pagination.current =1
    fetchData();
  };
  //重置
  const handleReset = () => {
    pagination.current =1
    formModel.value = generateFormModel();
    fetchData();
  };

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
  const handleAdd=(record:any)=>{
    openModal(true, {
      isUpdate: false,
      pid: record.id
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
    pagination.current =1
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
        Message.loading({content:"更新状态中",id:"upStatus"})
       const res= await upStatus({id:record.id,status:record.status});
       if(res){
         Message.success({content:"更新状态成功",id:"upStatus"})
       }
    }catch (error) {
      Message.clear("top")
    } 
  }
  //删除数据
  const handleDel=async(record:any)=>{
    try {
        Message.loading({content:"删除中",id:"upStatus"})
       const res= await del({ids:[record.id]});
       if(res){
        fetchData();
         Message.success({content:"删除成功",id:"upStatus"})
       }
    }catch (error) {
      Message.clear("top")
    } 
}
  //加载更多
  const loadMore=async(record :any, done:any)=>{
      const res= await getMoreList({pid:record.id});
      done(res)
  }
</script>

<script lang="ts">
  export default {
    name: 'index',
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
