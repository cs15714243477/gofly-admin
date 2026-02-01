<template>
    <div class="hcontent" v-show="activeKey==subnum" :style="{height:`${winHeights}px`}">
      温馨提示：如果您的admin端需要富文本编辑请自己引入，参考business端，admin端以简单轻量为主，如果您的业务需要自己引入。
    </div>
</template>

<script lang="ts" setup>
import { ref,onUnmounted,watch } from 'vue'
const props = defineProps({
  activeKey: Number,
  subnum: Number,
  winHeights: Number,
  placeholder: String,
  modelValue: {
    type: String,
    required: true
  },
})
const emits = defineEmits(['update:modelValue'])
const editorRef = ref();
const lockdit= ref(true);
const handleFocus=()=>{
  lockdit.value=false
}
onUnmounted(() => {
  lockdit.value=true
});
//编辑器
watch(props, (nweProps) => {
  if(editorRef.value&&lockdit.value){
    editorRef.value.setVal(nweProps.modelValue)
  }
})
//编辑器返回数据
const handleEditUpdta=(val:string)=>{
  emits('update:modelValue', val)
}
</script>
<style lang="less" scoped>
</style>
