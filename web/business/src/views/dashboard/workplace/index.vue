<template>
  <div class="page-container">
    <div class="content-wrapper">
      <Banner />
      <DataPanel />
      <a-grid :cols="24" :col-gap="16" :row-gap="16" class="content-grid">
        <a-grid-item :span="{ xs: 24, sm: 24, md: 24, lg: 14, xl: 14, xxl: 14 }">
          <ContentChart />
        </a-grid-item>
        <a-grid-item :span="{ xs: 24, sm: 24, md: 24, lg: 10, xl: 10, xxl: 10 }">
          <PopularContent />
        </a-grid-item>
      </a-grid>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import Banner from './components/banner.vue';
  import DataPanel from './components/data-panel.vue';
  import ContentChart from './components/content-chart.vue';
  import PopularContent from './components/popular-content.vue';
</script>

<script lang="ts">
  export default {
    name: 'home', // If you want the include property of keep-alive to take effect, you must name the component
  };
</script>

<style lang="less" scoped>
  .page-container {
    padding: 16px;
    background-color: var(--color-fill-1);
    height: auto;
  }

  .content-wrapper {
    max-width: 1440px;
    margin: 0 auto;
  }

  .content-grid {
    margin-top: 16px;
  }

  :deep(.dash-card.arco-card) {
    border-radius: 12px;
    border: 1px solid rgb(var(--gray-2));
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
  }
</style>
