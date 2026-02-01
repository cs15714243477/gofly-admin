<template>
    <Icon :icon="icon+(modelValue==val?'-select':'')" size="35" @click="selectIcon"></Icon>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import {  Icon} from '@/components/Icon';
const props = defineProps({
  icon: String,
  val: Number,
  modelValue: {
    type: Number,
    required: true
  },
})
const emits = defineEmits(['update:modelValue'])
const iconval = ref("");
//组件挂载完成后执行的函数
onMounted(() => {
})
//选择图标
const selectIcon=()=>{
      emits('update:modelValue', props.val)
}
</script>
