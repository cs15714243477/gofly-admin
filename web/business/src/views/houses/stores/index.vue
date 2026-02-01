<template>
  <div class="container" >
    <page-card breadcrumb fullTable :isFullscreen="isFullscreen">
      <template #searchBar>
      <div class="search-row" :class="{unhasesearchbtn:!showsearchdownbtn,hasesearchbtn:showsearchdownbtn}">
        <div class="search-field">
          <div class="search-line search-drown" :class="{searchdowncss:searchdown}">
            <div class="search-drown-box" ref="searchRef">
              <a-space wrap>
                <a-input style="width: 180px" v-model="formModel.name" placeholder="门店名称" allow-clear />
                <a-input style="width: 220px" v-model="formModel.address" placeholder="门店地址" allow-clear />
                <a-input style="width: 160px" v-model="formModel.contact_phone" placeholder="联系电话" allow-clear />
                <a-input style="width: 160px" v-model="formModel.manager_name" placeholder="店长姓名" allow-clear />
                <a-select style="width: 120px" v-model="formModel.status" placeholder="状态" allow-clear>
                  <a-option :value="0">启用</a-option>
                  <a-option :value="1">禁用</a-option>
                </a-select>
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
                {{ $t('searchTable.form.search') }}
              </a-button>
              <a-button @click="handleReset">
                {{ $t('searchTable.form.reset') }}
              </a-button>
            </a-space>
        </div>
        <div class="search-option">
          <tabletool :showbtn="['create','refresh','selectdensity','fullscreen']"
                    @create="createData" @refresh="fetchData" @selectdensity="(data)=>size=data" @fullscreen="(data)=>isFullscreen=data"></tabletool>
        </div>
      </div>
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
          @page-change="handlePageChange" 
          @page-size-change="handlePageSizeChange" 
        >
          <template #cellcopy="{ record, column }">
            <a-typography-paragraph
              :ellipsis="{ rows: 1, showTooltip: true }"
              copyable
              style="margin: 0"
            >
              {{ String(record[column.dataIndex] ?? '') }}
            </a-typography-paragraph>
          </template>

          <template #images="{ record }">
            <div class="img-cell" v-if="parseImages(record.images).length > 0">
              <a-image-preview-group>
                <a-image
                  :src="getImageUrl(parseImages(record.images)[0])"
                  width="44"
                  height="44"
                  fit="cover"
                  class="img-thumb"
                />
                <a-image
                  v-for="(img, idx) in parseImages(record.images).slice(1)"
                  :key="idx"
                  :src="getImageUrl(img)"
                  style="display: none"
                />
              </a-image-preview-group>
              <a-tag size="small" color="arcoblue">{{ parseImages(record.images).length }}张</a-tag>
            </div>
            <span v-else class="muted">无</span>
          </template>

          <template #switchstatus="{ record }">
            <a-switch
              type="round"
              size="small"
              v-model="record.status"
              :checked-value="0"
              :unchecked-value="1"
              @change="() => handleStatus(record, 'status')"
            />
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
    <AddForm @register="registerModal" @success="handleData"/>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, watch, onMounted } from 'vue';
  import useLoading from '@/hooks/loading';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  import { GetFullPath } from '@/utils/tool';
  //api
  import { getList,upStatus,del} from './api';
  //数据
  import { columns} from './data';
  import { shortcuts, dayjs} from '@/utils/dayjs';
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
  const isFullscreen = ref(false);
  const size = ref<SizeProps>('large');
   //查询字段
   const generateFormModel = () => {
    return {
      name: "",
      address: "",
      contact_phone: "",
      manager_name: "",
      status: "",
    };
  };
  const formModel = ref(generateFormModel());

  const parseImages = (v:any):string[] => {
    if(!v) return [];
    if(typeof v === "string") return v.split(",").map((i)=>i.trim()).filter(Boolean);
    return [];
  }
  const getImageUrl = (url:string) => {
    if(!url) return "";
    return GetFullPath(url) || url;
  }
  const fetchData = async () => {
    setLoading(true);
    searchdown.value=false
    try {
      const resp:any = await getList({page:pagination.current,pageSize:pagination.pageSize,...formModel.value});
      const data =
        resp?.items
          ? resp
          : resp?.data?.items
            ? resp.data
            : resp?.data?.data?.items
              ? resp.data.data
              : resp?.data
                ? resp.data
                : resp;
      renderData.value = (data?.items || []) as any;
      pagination.current = Number(data?.page || pagination.current);
      pagination.total = Number(data?.total || 0);
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
    searchbar()
    window.onresize = () => {
      searchbar()
    }
  })
 //更新搜索框按钮状态
 const searchbar=()=>{
  if(searchRef.value){
    const searchboxheight=searchRef.value["offsetHeight"]
    if(searchboxheight>40){
      showsearchdownbtn.value=true
    }else{
      showsearchdownbtn.value=false
      searchdown.value=false
    }
  }
 }
 //查找
 const handleSearch = () => {
    pagination.current=1
    fetchData();
  };
  //重置
  const handleReset = () => {
    pagination.current=1
    formModel.value = generateFormModel();
    fetchData();
  };

  watch(
    () => columns.value,
    (val) => {
      cloneColumns.value = cloneDeep(val);
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
    pagination.current=1
    fetchData();
  }
  //分页
  const handlePageChange = (page:any) => {
    pagination.current=page
    fetchData();
  }
  //分页总数
  const handlePageSizeChange = (pageSize:any) => {
    pagination.pageSize=pageSize
    fetchData();
  }
  //更新状态
  const handleStatus=async(record:any,field:string)=>{
    try {
      Message.loading({content:"更新状态中",id:"up",duration:0})
      await upStatus({id:record.id,[field]:record[field]});
      Message.success({content:"更新状态成功",id:"up",duration:2000})
    }catch (error) {
      record.status=record.status==1?0:1
      Message.loading({content:"更新状态中",id:"up",duration:1})
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
    name: '{vuetplname}', // If you want the include property of keep-alive to take effect, you must name the component
  };
</script>
<style scoped lang="less">
  .img-cell{
    display:flex;
    align-items:center;
    gap:8px;
  }
  .img-thumb{
    border-radius:6px;
    overflow:hidden;
  }
  .muted{
    color: var(--color-text-3);
  }
  :deep(.arco-table-th) {
    &:last-child {
      .arco-table-th-item-title {
        margin-left: 16px;
      }
    }
  }
</style>
