<template>
  <div class="banner">
    <div class="left">
      <div class="title">
        <span class="welcome">{{ $t('workplace.welcome') }}</span>
        <span class="name">{{ userInfo.name }}</span>
      </div>
      <div class="sub">数据实时统计，帮助你快速掌握门店经营概况</div>
    </div>
    <div class="right">
      <div class="date">{{ todayText }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import { useUserStore } from '@/store';

  const userStore = useUserStore();
  const userInfo = computed(() => {
    return {
      name: userStore.name,
    };
  });

  const todayText = computed(() => {
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, '0');
    const d = String(now.getDate()).padStart(2, '0');
    return `${y}年${m}月${d}日`;
  });
</script>

<style scoped lang="less">
  .banner {
    width: 100%;
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    padding: 4px 4px 12px 4px;
  }

  .title {
    display: flex;
    align-items: baseline;
    color: var(--color-text-1);
    font-size: 20px;
    font-weight: 600;
    line-height: 28px;
  }

  .welcome {
    color: var(--color-text-1);
  }

  .name {
    margin-left: 8px;
    color: var(--color-text-1);
  }

  .sub {
    margin-top: 4px;
    color: var(--color-text-3);
    font-size: 12px;
    line-height: 18px;
  }

  .date {
    color: var(--color-text-3);
    font-size: 12px;
    line-height: 18px;
  }
</style>
