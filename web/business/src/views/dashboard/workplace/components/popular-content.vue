<template>
  <a-spin :loading="loading" style="width: 100%">
    <a-card
      class="dash-card"
      :bordered="false"
      :header-style="{ padding: '16px 16px 0 16px', borderBottom: 'none' }"
      :body-style="{ padding: '0 16px 16px 16px' }"
      :title="$t('workplace.popularContent')"
    >
      <a-table
        :data="renderList"
        :pagination="false"
        :bordered="false"
        :scroll="{ x: '100%', y: '310px' }"
      >
        <template #columns>
          <a-table-column title="排名" :width="70">
            <template #cell="{ rowIndex }">
              {{ rowIndex + 1 }}
            </template>
          </a-table-column>
          <a-table-column title="房源标题" data-index="title" :ellipsis="true">
            <template #cell="{ record }">
              <a-typography-paragraph
                :ellipsis="{
                  rows: 1,
                }"
              >
                {{ record.title }}
              </a-typography-paragraph>
            </template>
          </a-table-column>
          <a-table-column title="浏览" data-index="viewCount" :width="90" />
          <a-table-column title="关注" data-index="followCount" :width="90" />
          <a-table-column title="带看" data-index="showingCount" :width="90" />
          <a-table-column title="状态" data-index="saleStatus" :width="90">
            <template #cell="{ record }">
              <a-tag
                size="small"
                :color="saleStatusColor(record.saleStatus)"
                bordered
              >
                {{ saleStatusText(record.saleStatus) }}
              </a-tag>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>
  </a-spin>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import useLoading from '@/hooks/loading';
  import { getPopular, PopularPropertyRecord } from '@/api/dashboard/workplace';

  const { loading, setLoading } = useLoading(true);
  const renderList = ref<PopularPropertyRecord[]>([]);

  const saleStatusText = (status: PopularPropertyRecord['saleStatus']) => {
    if (status === 'on_sale') return '在售';
    if (status === 'in_sale') return '预售';
    if (status === 'sold') return '已售';
    if (status === 'off_market') return '下架';
    return '未知';
  };
  const saleStatusColor = (status: PopularPropertyRecord['saleStatus']) => {
    if (status === 'on_sale') return 'green';
    if (status === 'in_sale') return 'gold';
    if (status === 'sold') return 'orange';
    if (status === 'off_market') return 'gray';
    return 'blue';
  };

  const fetchData = async () => {
    try {
      setLoading(true);
      renderList.value = await getPopular({});
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      setLoading(false);
    }
  };
  fetchData();
</script>

<style scoped lang="less">
  .dash-card {
    min-height: 395px;
  }
  :deep(.arco-table-tr) {
    height: 44px;
    .arco-typography {
      margin-bottom: 0;
    }
  }
</style>
