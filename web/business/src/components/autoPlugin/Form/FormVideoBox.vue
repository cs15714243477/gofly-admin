<template>
  <div class="upvideobox">
    <div class="videobtn">
      <div class="upload-show-video" v-if="modelValue">
        <div class="video-poster" @click="() => (previewVisible = true)">
          <icon-play-circle class="play-icon" :size="34" />
          <div class="video-url" :title="modelValue">{{ modelValue }}</div>
        </div>
        <div class="upload-show-video-mask">
          <a-space>
            <icon-eye @click="() => (previewVisible = true)" class="opbtn" />
            <IconEdit @click="UpVideo" class="opbtn" />
            <icon-delete @click="DelVideo" class="opbtn" />
          </a-space>
        </div>
      </div>
      <div class="upload-video-card" v-else @click="UpVideo">
        <div class="upload-video-card-text">
          <IconPlus />
          <div style="margin-top: 10px; font-weight: 600">上传视频</div>
        </div>
      </div>
    </div>
  </div>

  <a-modal v-model:visible="previewVisible" title="视频预览" width="860px" :footer="false" :maskClosable="true">
    <div class="video-preview">
      <video
        v-if="previewVisible && modelValue"
        ref="videoRef"
        :src="GetFullPath(modelValue)"
        controls
        playsinline
      />
      <a-empty v-else description="暂无视频" />
    </div>
  </a-modal>

  <!--附件-->
  <FileManage @register="registerFileModal" @success="selectVideo" />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import FileManage from '@/components/attachment/FileManage.vue';
import { useModal } from '/@/components/Modal';
import { GetFullPath } from '@/utils/tool';

defineProps({
  placeholder: String,
  modelValue: {
    type: String,
    required: true,
  },
});

const emits = defineEmits(['update:modelValue']);

const [registerFileModal, { openModal: openFileModal }] = useModal();
const previewVisible = ref(false);
const videoRef = ref<HTMLVideoElement | null>(null);

// 上传/选择视频
const UpVideo = () => {
  openFileModal(true, {
    filetype: 'video',
    getnumber: 'one', // one 单个
    openfrom: 'use', // manage管理 use选择使用
  });
};

// 选择附件返回
const selectVideo = (item: any) => {
  if (item.type === 'one') {
    emits('update:modelValue', item.url || '');
  }
};

// 删除
const DelVideo = () => {
  emits('update:modelValue', '');
};

// 关闭预览时，强制暂停并退出播放（防止音频后台继续）
watch(previewVisible, (v) => {
  if (v) return;
  const el = videoRef.value;
  if (!el) return;
  try {
    el.pause();
    el.currentTime = 0;
  } catch {}
});
</script>

<style lang="less" scoped>
.upvideobox {
  display: flex;

  .videobtn {
    position: relative;
    width: 220px;
    height: 90px;
    background-color: var(--color-neutral-1);
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;

    .upload-show-video {
      position: relative;
      box-sizing: border-box;
      width: 100%;
      height: 100%;
      overflow: hidden;
      display: flex;
      align-items: center;
      justify-content: center;

      &:hover {
        .upload-show-video-mask {
          opacity: 1;
        }
      }

      .video-poster {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-direction: column;
        gap: 6px;
        cursor: pointer;

        .play-icon {
          color: rgb(var(--primary-6));
        }

        .video-url {
          width: 92%;
          font-size: 12px;
          color: var(--color-text-3);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          text-align: center;
        }
      }

      .upload-show-video-mask {
        position: absolute;
        top: 0;
        right: 0;
        bottom: 0;
        left: 0;
        color: var(--color-white);
        font-size: 16px;
        line-height: 90px;
        text-align: center;
        background: rgba(0, 0, 0, 0.45);
        cursor: pointer;
        opacity: 0;
        transition: opacity 0.12s cubic-bezier(0, 0, 1, 1);
      }
    }

    .upload-video-card {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;

      .upload-video-card-text {
        text-align: center;
        color: var(--color-neutral-5);
      }
    }
  }
}

.video-preview {
  video {
    width: 100%;
    max-height: 62vh;
    border-radius: 10px;
    background: #000;
  }
}
</style>
