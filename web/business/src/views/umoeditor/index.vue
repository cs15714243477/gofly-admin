<template>
  <div class="container">
        <umo-editor v-bind="options" :theme="appStore.theme"/>
  </div>
</template>

<script lang="ts" setup>
import { ref,onMounted,onUnmounted, nextTick } from 'vue';
import { UmoEditor } from '@umoteam/editor';
import { useAppStore, } from '@/store';
const appStore = useAppStore();
 const options = ref({
   // 配置项
   height: '100%',
   document: {
    autofocus: false,
   }
 });
 //加载组件
 onMounted(()=>{
       appStore.footer=false
  })
  //离开
  onUnmounted(() => {
    appStore.footer=true
  })
</script>
<script lang="ts">
  export default {
    name: 'wpeditor', // If you want the include property of keep-alive to take effect, you must name the component
  };
</script>

<style lang="less" scoped>
.container{
  position: relative;
  display: block;
  padding: 0;
  margin: 0;
  width: 100%;
  height: calc(100vh - 50px);
  background-color: var(--color-bg-1);
}
</style>
