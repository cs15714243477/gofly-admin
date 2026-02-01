<template>
  <div class="upimagebox">
      <div class="imagebtn" v-for="(img,index) in imageslist">
        <div class="upload-show-picture">
          <a-image :src="GetFullPath(img)" height="90" :preview-visible="visibleimage[index]"
            @preview-visible-change="() => { visibleimage[index] = false }" />
          <div class="upload-show-picture-mask">
            <a-space>
              <icon-eye @click="() => visibleimage[index] = true" class="opbtn" />
              <IconEdit @click="UpImage('one',index)" class="opbtn" />
              <icon-delete @click="DelImage(img)" class="opbtn"/>
            </a-space>
          </div>
        </div>
      </div>
    <div class="imagebtn">
      <div class="upload-picture-card" @click="UpImage('more',0)">
        <div class="upload-picture-card-text">
          <IconPlus />
          <div style="margin-top: 10px; font-weight: 600">上传图片</div>
        </div>
      </div>
    </div>
   <!--附件-->
    <FileManage @register="registerFileModal" @success="selectImg"/>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from 'vue'
import FileManage from '@/components/attachment/FileManage.vue';
import { useModal} from '/@/components/Modal';
import { GetFullPath } from "@/utils/tool";
const props = defineProps({
  placeholder: String,
  modelValue: {
    type: String,
    required: true
  },
})
const emits = defineEmits(['update:modelValue'])
const [registerFileModal, { openModal:openFileModal }] = useModal();
const imageslist = ref<any[]>([]);
const visibleimage = ref<boolean []>([])
//组件挂载完成后执行的函数
onMounted(() => {
})
watch(props, (nweProps) => {
  imageslist.value=[]
  if (nweProps.modelValue ) {
    imageslist.value=nweProps.modelValue.split(",")
  }
})
//上传图片
const UpImage=(getnumber:string,index:number)=>{
  openFileModal(true, {
    param:index+1,
    filetype:"image",
    getnumber: getnumber,//one 单张 more多张
    openfrom: "use",//manage管理 use选择使用
  });
}
//选择附件返回
const selectImg=(item:any)=>{
  if(item.type=="more"){
        item.list.forEach((img:any) => {
          imageslist.value.push(img)
        });
    }else if(item.type=="one"){
      imageslist.value[item.param-1]=item.url
    }
    emits('update:modelValue', imageslist.value.join(","))
}
//删除
const DelImage=(img:any)=>{
  imageslist.value= imageslist.value.filter((item)=>item!=img)
  emits('update:modelValue', imageslist.value.join(","))
}
</script>
<style lang="less" scoped>
 //上传图片
 .upimagebox{
    display: flex;
    flex-wrap: wrap;
    width: 100%;
    overflow: hidden;
    margin-bottom: -10px;
    .imagebtn{
      position: relative;
      width: 145px;
      height: 90px;
      background-color: var(--color-neutral-1);
      border-radius: 4px;
      overflow: hidden;
      flex-shrink: 0;
      //多图
      margin-right: 10px;
      margin-bottom: 10px;
      //预览
      .upload-show-picture{
        position: relative;
        box-sizing: border-box;
        width: 100%;
        height: 100%;
        overflow: hidden;
        display: flex;
        align-items: center;
        justify-content: center;
        img{
          height: 100%;
        }
        &:hover{
          .upload-show-picture-mask{
             opacity: 1;
          }
        }
        .upload-show-picture-mask{
            position: absolute;
            top: 0;
            right: 0;
            bottom: 0;
            left: 0;
            color: var(--color-white);
            font-size: 16px;
            line-height: 90px;
            text-align: center;
            background: rgba(0, 0, 0, 0.5);
            cursor: pointer;
            opacity: 0;
            transition: opacity 0.1s cubic-bezier(0, 0, 1, 1);
        }
      }
      .upload-picture-card{
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        .upload-picture-card-text{
           text-align: center;
           color:  var(--color-neutral-5);
        }
      }
    }
  }
</style>
