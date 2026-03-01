<template>
  <a-spin :loading="loading" style="width: 100%">
    <a-card
      class="dash-card metrics-card"
      :bordered="false"
      title="核心指标"
      :header-style="{ padding: '16px 16px 0 16px', borderBottom: 'none' }"
      :body-style="{ padding: '16px 16px 12px 16px' }"
    >
      <template #extra>
        <div class="extra">
          <span class="dot" />
          <span class="text">实时数据</span>
        </div>
      </template>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.onlineContent') }}</div>
              <div class="metric-icon icon-primary"><icon-home /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.propertyTotal) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">全部已录入房源</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.putIn') }}</div>
              <div class="metric-icon icon-success"><icon-check-circle /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.propertyOnSale) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">可对客户展示</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.inSale') }}</div>
              <div class="metric-icon icon-warning"><icon-clock-circle /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.propertyInSale) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">预售中</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.sold') }}</div>
              <div class="metric-icon icon-purple"><icon-trophy /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.propertySold) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">已成交</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.offMarket') }}</div>
              <div class="metric-icon icon-gray"><icon-close /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.propertyOffMarket) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">已下架</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.lockBindProperty') }}</div>
              <div class="metric-icon icon-info"><icon-lock /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.lockBindPropertyTotal) }}<span class="metric-unit">套</span>
            </div>
            <div class="metric-sub">房源已绑定智能锁</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.newDay') }}</div>
              <div class="metric-icon icon-warning"><icon-lock /></div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.lockBindTotal) }}<span class="metric-unit">把</span>
            </div>
            <div class="metric-sub">绑定锁数量</div>
          </div>
        </a-col>

        <a-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6" :xxl="6">
          <div class="metric-tile">
            <div class="metric-top">
              <div class="metric-label">{{ $t('workplace.newFromYesterday') }}</div>
              <div class="metric-icon icon-danger">
                <icon-exclamation-circle />
              </div>
            </div>
            <div class="metric-value">
              {{ formatNumber(stat.unlockPendingTotal) }}<span class="metric-unit">条</span>
            </div>
            <div class="metric-sub">待审核开锁申请</div>
          </div>
        </a-col>
      </a-row>

      <div class="today-strip">
        <div class="today-item">
          今日新增房源 <b>{{ formatNumber(stat.todayPropertyAdd) }}</b> 套
        </div>
        <div class="today-item">
          今日开锁申请 <b>{{ formatNumber(stat.todayUnlockRequests) }}</b> 条
        </div>
        <div class="today-item">
          今日浏览 <b>{{ formatNumber(stat.todayViewCount) }}</b> 次
        </div>
        <div class="today-item">
          今日带看 <b>{{ formatNumber(stat.todayShowingCount) }}</b> 次
        </div>
      </div>
    </a-card>
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
    propertyInSale: 0,
    propertySold: 0,
    propertyOffMarket: 0,
    lockBindTotal: 0,
    lockBindPropertyTotal: 0,
    unlockPendingTotal: 0,
    todayPropertyAdd: 0,
    todayUnlockRequests: 0,
    todayViewCount: 0,
    todayShowingCount: 0,
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

  const formatNumber = (v: number) => {
    return Number(v || 0).toLocaleString();
  };
</script>

<style lang="less" scoped>
  .metrics-card {
    margin-top: 8px;
  }

  .extra {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    color: var(--color-text-3);
    font-size: 12px;
    line-height: 18px;
  }
  .extra .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: rgb(var(--green-6));
  }

  .metric-tile {
    height: 112px;
    padding: 14px 14px 12px;
    border-radius: 12px;
    border: 1px solid rgb(var(--gray-2));
    background: var(--color-fill-1);
    transition: box-shadow 0.2s ease, border-color 0.2s ease,
      background-color 0.2s ease, transform 0.2s ease;
  }
  .metric-tile:hover {
    background: var(--color-fill-2);
    border-color: rgb(var(--arcoblue-3));
    box-shadow: 0 10px 26px rgba(0, 0, 0, 0.08);
    transform: translateY(-1px);
  }

  .metric-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 10px;
  }

  .metric-label {
    color: var(--color-text-3);
    font-size: 12px;
    line-height: 18px;
  }

  .metric-icon {
    width: 34px;
    height: 34px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .metric-icon :deep(.arco-icon) {
    font-size: 16px;
  }

  .icon-primary {
    background: rgba(var(--arcoblue-6), 0.12);
    color: rgb(var(--arcoblue-6));
  }
  .icon-success {
    background: rgba(var(--green-6), 0.12);
    color: rgb(var(--green-6));
  }
  .icon-warning {
    background: rgba(var(--orange-6), 0.12);
    color: rgb(var(--orange-6));
  }
  .icon-danger {
    background: rgba(var(--red-6), 0.12);
    color: rgb(var(--red-6));
  }
  .icon-purple {
    background: rgba(var(--purple-6), 0.12);
    color: rgb(var(--purple-6));
  }
  .icon-gray {
    background: rgba(var(--gray-6), 0.12);
    color: rgb(var(--gray-8));
  }
  .icon-info {
    background: rgba(var(--cyan-6), 0.12);
    color: rgb(var(--cyan-6));
  }

  .metric-value {
    color: var(--color-text-1);
    font-size: 24px;
    font-weight: 700;
    line-height: 30px;
    letter-spacing: 0.2px;
  }

  .metric-unit {
    margin-left: 4px;
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-3);
  }

  .metric-sub {
    margin-top: 6px;
    color: var(--color-text-3);
    font-size: 12px;
    line-height: 18px;
  }

  .today-strip {
    margin-top: 16px;
    padding-top: 12px;
    border-top: 1px dashed rgb(var(--gray-2));
    display: flex;
    gap: 10px;
    flex-wrap: wrap;
  }

  .today-item {
    padding: 6px 10px;
    border-radius: 999px;
    border: 1px solid rgb(var(--gray-2));
    background: var(--color-fill-1);
    color: var(--color-text-2);
    font-size: 12px;
    line-height: 18px;
  }

  .today-item b {
    padding: 0 2px;
    color: var(--color-text-1);
    font-weight: 700;
  }
</style>
