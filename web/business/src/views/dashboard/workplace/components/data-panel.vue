<template>
  <a-spin :loading="loading" style="width: 100%">
    <a-grid :cols="24" :row-gap="16" class="panel">
      <a-grid-item
        class="panel-col"
        :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }"
      >
        <a-space>
          <a-avatar :size="54" class="col-avatar avatar-primary">
            <icon-home />
          </a-avatar>
          <a-statistic
            :title="$t('workplace.onlineContent')"
            :value="stat.propertyTotal"
            :value-from="0"
            animation
            show-group-separator
          >
            <template #suffix>套</template>
          </a-statistic>
        </a-space>
      </a-grid-item>
      <a-grid-item
        class="panel-col"
        :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }"
      >
        <a-space>
          <a-avatar :size="54" class="col-avatar avatar-success">
            <icon-check-circle />
          </a-avatar>
          <a-statistic
            :title="$t('workplace.putIn')"
            :value="stat.propertyOnSale"
            :value-from="0"
            animation
            show-group-separator
          >
            <template #suffix>套</template>
          </a-statistic>
        </a-space>
      </a-grid-item>
      <a-grid-item
        class="panel-col"
        :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }"
      >
        <a-space>
          <a-avatar :size="54" class="col-avatar avatar-warning">
            <icon-lock />
          </a-avatar>
          <a-statistic
            :title="$t('workplace.newDay')"
            :value="stat.lockBindTotal"
            :value-from="0"
            animation
            show-group-separator
          >
            <template #suffix>把</template>
          </a-statistic>
        </a-space>
      </a-grid-item>
      <a-grid-item
        class="panel-col"
        :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }"
        style="border-right: none"
      >
        <a-space>
          <a-avatar :size="54" class="col-avatar avatar-danger">
            <icon-exclamation-circle />
          </a-avatar>
          <a-statistic
            :title="$t('workplace.newFromYesterday')"
            :value="stat.unlockPendingTotal"
            :value-from="0"
            animation
            show-group-separator
          >
            <template #suffix>条</template>
          </a-statistic>
        </a-space>
      </a-grid-item>
      <a-grid-item :span="24">
        <div class="panel-foot">
          <span>今日新增房源 <b>{{ stat.todayPropertyAdd }}</b> 套</span>
          <span class="dot">·</span>
          <span>今日开锁申请 <b>{{ stat.todayUnlockRequests }}</b> 条</span>
        </div>
      </a-grid-item>
      <a-grid-item :span="24">
        <a-divider class="panel-border" />
      </a-grid-item>
    </a-grid>
  </a-spin>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import useLoading from '@/hooks/loading';
  import {
    getStatistical,
    WorkplaceStatistical,
  } from '@/api/dashboard/workplace';

  const { loading, setLoading } = useLoading(true);
  const stat = ref<WorkplaceStatistical>({
    propertyTotal: 0,
    propertyOnSale: 0,
    lockBindTotal: 0,
    unlockPendingTotal: 0,
    todayPropertyAdd: 0,
    todayUnlockRequests: 0,
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      stat.value = await getStatistical({});
    } catch (err) {
      // 可以在这里统一处理错误提示
    } finally {
      setLoading(false);
    }
  };

  onMounted(fetchData);
</script>

<style lang="less" scoped>
  .arco-grid.panel {
    margin-bottom: 0;
    padding: 16px 20px 0 20px;
  }
  .panel-col {
    padding: 6px 12px;
    border-right: 1px solid rgb(var(--gray-2));
  }
  .col-avatar {
    margin-right: 12px;
    background-color: rgba(var(--arcoblue-6), 0.12);
    color: rgb(var(--arcoblue-6));
  }
  .avatar-success {
    background-color: rgba(var(--green-6), 0.12);
    color: rgb(var(--green-6));
  }
  .avatar-warning {
    background-color: rgba(var(--orange-6), 0.12);
    color: rgb(var(--orange-6));
  }
  .avatar-danger {
    background-color: rgba(var(--red-6), 0.12);
    color: rgb(var(--red-6));
  }
  .panel-foot {
    margin-top: -6px;
    padding: 0 12px;
    color: var(--color-text-2);
    font-size: 12px;
    b {
      color: var(--color-text-1);
      font-weight: 700;
      padding: 0 2px;
    }
    .dot {
      padding: 0 8px;
      color: rgb(var(--gray-5));
    }
  }
  :deep(.panel-border) {
    margin: 4px 0 0 0;
  }
</style>
