<template>
  <div class="tab-bar-container">
      <div class="tab-bar-box">
        <a-tabs  hide-content size="medium" :active-key="route.fullPath">
          <a-tab-pane  v-for="(tag, index) in tagList"  :key="tag.fullPath ">
            <template #title>    
              <tab-item  :index="index" :item-data="tag"/>
            </template>
          </a-tab-pane>
          <template #extra>
            <a-dropdown trigger="hover">
              <icon-drag-dot-vertical />
              <template #content>
                <a-doption @click="tagClose">
                  <template #icon><icon-close /></template>
                  <template #default>关闭当前</template>
                </a-doption>
                <a-doption @click="tagCloseOthers">
                  <template #icon><icon-eraser /></template>
                  <template #default>关闭其他</template>
                </a-doption>
                <a-doption @click="tagCloseAll">
                  <template #icon><icon-minus /></template>
                  <template #default>关闭全部</template>
                </a-doption>
              </template>
            </a-dropdown>
          </template>
        </a-tabs>
      </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed, watch, onUnmounted } from 'vue';
  import type { RouteLocationNormalized } from 'vue-router';
  import type { TagProps } from '@/store/modules/tab-bar/types';
  import { useRouter, useRoute } from 'vue-router';
  import { DEFAULT_ROUTE_NAME } from '@/router/constants';
  import {
    listenerRouteChange,
    removeRouteListener,
  } from '@/utils/route-listener';
  import { useAppStore, useTabBarStore } from '@/store';
  import tabItem from './tab-item.vue';
  const route = useRoute();
  const router = useRouter();
  const appStore = useAppStore();
  const tabBarStore = useTabBarStore();

  const affixRef = ref();
  const tagList = computed(() => {
    return tabBarStore.getTabList;
  });

  watch(
    () => appStore.navbar,
    () => {
      affixRef.value.updatePosition();
    }
  );
  listenerRouteChange((route: RouteLocationNormalized) => {
    if(route.name==DEFAULT_ROUTE_NAME)tabBarStore.updateTabHome(route);
    if (
      !route.meta.noAffix &&
      !tagList.value.some((tag) => tag.fullPath === route.fullPath)
    ) {
      tabBarStore.updateTabList(route);
    }
  }, true);
  //关闭当前
  const tagClose = () => {
    tagList.value.forEach((tag:TagProps,idx:number)=>{
      if(tag.fullPath==route.fullPath&&DEFAULT_ROUTE_NAME!=tag.name){
          tabBarStore.deleteTag(idx, tag);
          const latest = tagList.value[idx - 1]; // 获取队列的前一个tab
          router.push({ name: latest.name });
      }
    })
  };
  //关闭其他
  const tagCloseOthers = () => {
    tagList.value.forEach((tag:TagProps,index:number)=>{
      if(tag.fullPath==route.fullPath){
        const filterList = tagList.value.filter((el, idx) => {
          return idx === 0 || idx ===index;
        });
        tabBarStore.freshTabList(filterList);
        router.push({ name: tag.name });
      }
    })
  };
  //关闭全部
  const tagCloseAll = () => {
    tabBarStore.resetTabList();
      router.push({ name: DEFAULT_ROUTE_NAME });
  };
  onUnmounted(() => {
    removeRouteListener();
  });
</script>

<style scoped lang="less">
:deep(.arco-tabs-nav-ink) {
   display: none;
}
:deep(.arco-tabs-nav::before) {
   display: none;
}
:deep(.arco-tabs-nav-tab) {
  .arco-tabs-tab {
    margin: 0 5px;
    padding: 0;
    user-select: none;
    border-bottom-color: transparent !important;
    svg {
      display: none;
      transition: all 0.15s;
    }
    &:hover {
      svg {
        display: unset;
      }
    }
    .arco-tabs-tab-title::before {
      right:0px;
      left:0px;
    }
    &:first-child {
      .arco-tag-close-btn {
        display: none;
      }
    }
  }
  &:not(.arco-tabs-nav-tab-scroll) {
    .arco-tabs-tab:first-child {
      border-left: 0;
    }
  }
}

.tab-bar-container {
    position: relative;
    .tab-bar-box{
      padding: 4px 5px;
      background-color: var(--color-bg-2);
      border-bottom: solid 1px  var(--color-fill-2);
    }
  }
</style>
