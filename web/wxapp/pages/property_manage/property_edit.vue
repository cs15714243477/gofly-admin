<template>
  <view class="pe">
    <TopHeader :title="id ? '编辑房源' : '新增房源'">
      <template #left>
        <view class="icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </template>
      <template #right>
        <view class="icon-btn primary" @tap="save">
          <text class="material-symbols-outlined">save</text>
        </view>
      </template>
    </TopHeader>

    <scroll-view scroll-y="true" class="content">
      <!-- 基础信息 -->
      <view class="card">
        <view class="card-title">基础信息</view>

        <view class="form-row">
          <text class="label required">房源标题</text>
          <input
            v-model="form.title"
            class="input"
            type="text"
            maxlength="100"
            placeholder="请输入房源标题"
            placeholder-class="placeholder"
          />
        </view>

        <view class="grid">
          <view class="form-row">
            <text class="label required">销售状态</text>
            <picker
              mode="selector"
              :range="saleStatusOptions.map((o) => o.label)"
              :value="saleStatusIndex"
              @change="onSaleStatusChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  saleStatusOptions[saleStatusIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic"
                  >expand_more</text
                >
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">价格</text>
            <view class="price-row">
              <input
                v-model="form.price"
                class="input"
                type="digit"
                placeholder="0"
                placeholder-class="placeholder"
              />
              <picker
                mode="selector"
                :range="priceUnitOptions"
                :value="priceUnitIndex"
                @change="onPriceUnitChange"
              >
                <view class="unit-picker">
                  <text class="unit-text">{{ form.price_unit || '万' }}</text>
                  <text class="material-symbols-outlined unit-ic"
                    >expand_more</text
                  >
                </view>
              </picker>
            </view>
          </view>

          <view class="form-row">
            <text class="label">面积(㎡)</text>
            <input
              v-model="form.area"
              class="input"
              type="digit"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">户型</text>
            <view class="room-row">
              <input
                v-model="form.rooms"
                class="mini"
                type="number"
                placeholder="室"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">室</text>
              <input
                v-model="form.halls"
                class="mini"
                type="number"
                placeholder="厅"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">厅</text>
              <input
                v-model="form.bathrooms"
                class="mini"
                type="number"
                placeholder="卫"
                placeholder-class="placeholder"
              />
              <text class="mini-suffix">卫</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 图片 -->
      <view class="card">
        <view class="card-title">图片</view>
        <view class="hint"
          >建议上传清晰横图，点击图片可设置封面/删除</view
        >

        <view class="media-grid">
          <view
            v-for="(img, idx) in images"
            :key="'img-' + idx"
            class="media-item"
            @tap="openImageActions(idx)"
          >
            <image class="media-img" :src="img" mode="aspectFill"></image>
            <view v-if="img === form.cover_image" class="badge">封面</view>
          </view>
          <view class="media-item add" @tap="pickImages">
            <text class="material-symbols-outlined add-ic">add_photo_alternate</text>
            <text class="add-text">添加</text>
          </view>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">download</text>
            <text class="sw-text">允许下载图片</text>
          </view>
          <switch
            :checked="Number(form.allow_image_download) === 1"
            @change="(e) => (form.allow_image_download = e.detail.value ? 1 : 0)"
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 视频 -->
      <view class="card">
        <view class="card-title">视频</view>
        <view class="hint">支持单个房源视频，建议 10-30 秒</view>

        <view class="video-box" @tap="openVideoActions">
          <view v-if="!form.video_url" class="video-empty">
            <text class="material-symbols-outlined video-ic">movie</text>
            <view class="video-texts">
              <text class="video-title">未上传视频</text>
              <text class="video-sub">点击上传/更换</text>
            </view>
          </view>
          <video
            v-else
            class="video"
            :src="form.video_url"
            controls
            object-fit="cover"
          ></video>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">download</text>
            <text class="sw-text">允许下载视频</text>
          </view>
          <switch
            :checked="Number(form.allow_video_download) === 1"
            @change="(e) => (form.allow_video_download = e.detail.value ? 1 : 0)"
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 标签 -->
      <view class="card">
        <view class="card-title">标签</view>
        <view class="tag-row">
          <view class="tag" v-for="(t, idx) in tags" :key="'t-' + idx">
            <text class="tag-text">{{ t }}</text>
            <text class="material-symbols-outlined tag-x" @tap.stop="removeTag(idx)"
              >close</text
            >
          </view>
          <view v-if="tags.length === 0" class="hint small">暂无标签</view>
        </view>
        <view class="tag-add">
          <input
            v-model="newTag"
            class="input"
            type="text"
            maxlength="12"
            placeholder="输入标签（最多6个）"
            placeholder-class="placeholder"
          />
          <button class="btn add" @tap="addTag">添加</button>
        </view>
      </view>

      <!-- 位置 -->
      <view class="card">
        <view class="card-title">位置</view>
        <view class="grid">
          <view class="form-row">
            <text class="label">小区名称</text>
            <input
              v-model="form.community_name"
              class="input"
              type="text"
              placeholder="请输入小区/楼盘"
              placeholder-class="placeholder"
            />
          </view>
          <view class="form-row">
            <text class="label">详细地址</text>
            <input
              v-model="form.address"
              class="input"
              type="text"
              placeholder="请输入详细地址"
              placeholder-class="placeholder"
            />
          </view>
        </view>

        <view class="coord-row">
          <view class="coord-left">
            <text class="material-symbols-outlined coord-ic">my_location</text>
            <text class="coord-text">{{
              form.latitude && form.longitude
                ? Number(form.latitude).toFixed(6) + ', ' + Number(form.longitude).toFixed(6)
                : '未选择坐标'
            }}</text>
          </view>
          <button class="btn ghost" @tap="pickLocation">地图选点</button>
        </view>
      </view>

      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="footer">
      <button
        class="footer-btn ghost"
        :disabled="!id"
        @tap="preview"
      >
        预览
      </button>
      <button class="footer-btn primary" @tap="save">
        {{ saving ? '保存中...' : '保存' }}
      </button>
    </view>
  </view>
</template>

<script>
import TopHeader from '@/components/TopHeader.vue'
import propertyApi from '@/api/property'
import $store from '@/store'
import md5 from 'js-md5'
import { baseUrl } from '@/utils/config'

// base64 编码（兼容小程序端无 window.btoa）
function base64Encode(str = '') {
  try {
    // #ifdef H5
    return window.btoa(str)
    // #endif
  } catch (e) {}

  // #ifndef H5
  const utf8ToBytes = (s) => {
    const bytes = []
    for (let i = 0; i < s.length; i++) {
      let codePoint = s.charCodeAt(i)
      if (codePoint < 0x80) {
        bytes.push(codePoint)
      } else if (codePoint < 0x800) {
        bytes.push(0xc0 | (codePoint >> 6))
        bytes.push(0x80 | (codePoint & 0x3f))
      } else {
        bytes.push(0xe0 | (codePoint >> 12))
        bytes.push(0x80 | ((codePoint >> 6) & 0x3f))
        bytes.push(0x80 | (codePoint & 0x3f))
      }
    }
    return new Uint8Array(bytes)
  }
  const buf = utf8ToBytes(str).buffer
  // eslint-disable-next-line no-undef
  if (typeof wx !== 'undefined' && wx.arrayBufferToBase64) return wx.arrayBufferToBase64(buf)
  if (typeof uni !== 'undefined' && uni.arrayBufferToBase64) return uni.arrayBufferToBase64(buf)
  return ''
  // #endif
}

export default {
  components: { TopHeader },
  data() {
    return {
      id: 0,
      saving: false,
      uploading: false,
      form: {
        title: '',
        price: '',
        price_unit: '万',
        area: '',
        rooms: 0,
        halls: 0,
        bathrooms: 0,
        community_name: '',
        address: '',
        latitude: '',
        longitude: '',
        tags: [],
        images: [],
        cover_image: '',
        video_url: '',
        allow_image_download: 1,
        allow_video_download: 1,
        sale_status: 'on_sale',
      },
      tags: [],
      images: [],
      newTag: '',
      saleStatusOptions: [
        { label: '在售', value: 'on_sale' },
        { label: '已售', value: 'sold' },
        { label: '下架', value: 'off_market' },
      ],
      priceUnitOptions: ['万', '元'],
    }
  },
  computed: {
    saleStatusIndex() {
      const v = String(this.form.sale_status || 'on_sale')
      const idx = this.saleStatusOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    priceUnitIndex() {
      const v = String(this.form.price_unit || '万')
      const idx = this.priceUnitOptions.findIndex((it) => it === v)
      return idx >= 0 ? idx : 0
    },
  },
  onLoad(options) {
    const id = Number((options && (options.id || options.ID)) || 0) || 0
    this.id = id
    this.ensurePermission().then(() => {
      if (this.id) this.loadContent()
    })
  },
  methods: {
    goBack() {
      uni.navigateBack()
    },
    async ensurePermission() {
      const userStore = $store('user')
      const token = uni.getStorageSync('token')
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: '/pages/login/login' })
        return
      }
      try {
        if (!userStore.userInfo || !userStore.userInfo.id) {
          await userStore.getInfo()
        }
      } catch (e) {}
      const canManage = Number((userStore.userInfo || {}).can_manage_properties) === 1
      if (!canManage) {
        uni.showModal({
          title: '提示',
          content: '暂无房源维护权限，请联系后台管理员开启。',
          showCancel: false,
          success: () => {
            uni.navigateBack()
          },
        })
      }
    },
    onSaleStatusChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.saleStatusOptions[idx] || this.saleStatusOptions[0]
      this.form.sale_status = opt.value
    },
    onPriceUnitChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      this.form.price_unit = this.priceUnitOptions[idx] || '万'
    },
    async loadContent() {
      const res = await propertyApi.getManageContent({ id: this.id })
      if (!res || res.code !== 0) return
      const data = res.data || {}
      this.form = {
        ...this.form,
        ...data,
      }
      this.tags = Array.isArray(data.tags) ? data.tags : []
      this.images = Array.isArray(data.images) ? data.images : []
      // 兜底：封面未设置时，用第一张图
      if (!this.form.cover_image && this.images.length > 0) {
        this.form.cover_image = this.images[0]
      }
    },
    addTag() {
      const t = String(this.newTag || '').trim()
      if (!t) {
        uni.showToast({ title: '请输入标签', icon: 'none' })
        return
      }
      if (this.tags.includes(t)) {
        this.newTag = ''
        return
      }
      if (this.tags.length >= 6) {
        uni.showToast({ title: '最多6个标签', icon: 'none' })
        return
      }
      this.tags = this.tags.concat([t])
      this.newTag = ''
    },
    removeTag(idx) {
      this.tags = this.tags.filter((_, i) => i !== idx)
    },
    openImageActions(idx) {
      const img = this.images[idx]
      if (!img) return
      const isCover = img === this.form.cover_image
      const list = isCover ? ['删除'] : ['设为封面', '删除']
      uni.showActionSheet({
        itemList: list,
        success: (res) => {
          const tap = Number(res.tapIndex)
          if (!isCover && tap === 0) {
            this.form.cover_image = img
            return
          }
          const delIndex = isCover ? 0 : 1
          if (tap === delIndex) {
            this.images = this.images.filter((_, i) => i !== idx)
            if (this.form.cover_image === img) {
              this.form.cover_image = this.images[0] || ''
            }
          }
        },
      })
    },
    openVideoActions() {
      const has = !!this.form.video_url
      const list = has ? ['更换视频', '移除视频'] : ['上传视频']
      uni.showActionSheet({
        itemList: list,
        success: async (res) => {
          const tap = Number(res.tapIndex)
          if (!has) {
            await this.pickVideo()
            return
          }
          if (tap === 0) await this.pickVideo()
          if (tap === 1) this.form.video_url = ''
        },
      })
    },
    pickImages() {
      if (this.uploading) return
      uni.chooseImage({
        count: 9,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: async (res) => {
          const files = (res && res.tempFilePaths) || []
          if (!files.length) return
          await this.uploadFiles(files, 'image')
        },
      })
    },
    pickVideo() {
      if (this.uploading) return Promise.resolve()
      return new Promise((resolve) => {
        uni.chooseVideo({
          sourceType: ['album', 'camera'],
          maxDuration: 60,
          compressed: true,
          success: async (res) => {
            const filePath = res && (res.tempFilePath || res.tempFilePaths?.[0])
            if (!filePath) return resolve()
            await this.uploadFiles([filePath], 'video')
            resolve()
          },
          fail: () => resolve(),
        })
      })
    },
    uploadFiles(filePaths = [], filetype = 'image') {
      if (!Array.isArray(filePaths) || filePaths.length === 0) return Promise.resolve()
      return new Promise(async (resolve) => {
        this.uploading = true
        uni.showLoading({ title: '上传中', mask: true })
        try {
          for (let i = 0; i < filePaths.length; i++) {
            const fp = filePaths[i]
            const out = await this.uploadSingle(fp, filetype).catch(() => null)
            const url = out && out.code === 0 && out.data ? out.data.url : ''
            if (!url) continue
            if (filetype === 'video') {
              this.form.video_url = url
            } else {
              this.images = this.images.concat([url])
              if (!this.form.cover_image) this.form.cover_image = url
            }
          }
        } finally {
          uni.hideLoading()
          this.uploading = false
          resolve()
        }
      })
    },
    uploadSingle(filePath, filetype) {
      return new Promise((resolve, reject) => {
        const token = uni.getStorageSync('token')
        const timestamp = Math.floor(Date.now() / 1000)
        const passstr = md5(import.meta.env.GF_API_SECRET + timestamp)
        const header = {
          Accept: 'text/json',
          Businessid: import.meta.env.GF_BUSINESUSSID,
          apiverify: base64Encode(passstr + '#' + timestamp),
        }
        if (token) header.Authorization = `${token}`

        uni.uploadFile({
          url: `${baseUrl}/common/upload/upFile`,
          filePath,
          name: 'file',
          formData: { filetype },
          header,
          success: (upRes) => {
            try {
              const raw = upRes && upRes.data ? upRes.data : '{}'
              const out = typeof raw === 'string' ? JSON.parse(raw) : raw
              if (!out || out.code !== 0) {
                uni.showToast({
                  title: (out && out.message) || '上传失败',
                  icon: 'none',
                })
                reject(out)
                return
              }
              resolve(out)
            } catch (e) {
              uni.showToast({ title: '解析上传结果失败', icon: 'none' })
              reject(e)
            }
          },
          fail: (err) => {
            uni.showToast({ title: '上传失败', icon: 'none' })
            reject(err)
          },
        })
      })
    },
    pickLocation() {
      uni.chooseLocation({
        success: (res) => {
          if (!res) return
          const lat = Number(res.latitude)
          const lng = Number(res.longitude)
          if (isFinite(lat) && isFinite(lng)) {
            this.form.latitude = lat
            this.form.longitude = lng
          }
          // 兼容：有 name 则给到小区名称
          if (res.name && !this.form.community_name) this.form.community_name = String(res.name)
          if (res.address) this.form.address = String(res.address)
        },
        fail: () => {
          uni.showToast({ title: '无法打开地图', icon: 'none' })
        },
      })
    },
    async save() {
      if (this.saving || this.uploading) return
      const title = String(this.form.title || '').trim()
      if (!title) {
        uni.showToast({ title: '请填写房源标题', icon: 'none' })
        return
      }
      this.saving = true
      try {
        const payload = {
          ...(this.id ? { id: this.id } : {}),
          title,
          sale_status: this.form.sale_status,
          price: this.form.price,
          price_unit: this.form.price_unit,
          area: this.form.area,
          rooms: Number(this.form.rooms || 0),
          halls: Number(this.form.halls || 0),
          bathrooms: Number(this.form.bathrooms || 0),
          community_name: String(this.form.community_name || '').trim(),
          address: String(this.form.address || '').trim(),
          latitude: this.form.latitude || '',
          longitude: this.form.longitude || '',
          tags: this.tags,
          images: this.images,
          cover_image: this.form.cover_image || '',
          video_url: this.form.video_url || '',
          allow_image_download: Number(this.form.allow_image_download) === 1 ? 1 : 0,
          allow_video_download: Number(this.form.allow_video_download) === 1 ? 1 : 0,
        }
        const res = await propertyApi.saveManage(payload)
        if (!res || res.code !== 0) return
        const nid = Number(res && res.data && res.data.id) || 0
        if (!this.id && nid) this.id = nid
        uni.showToast({ title: '保存成功', icon: 'none' })
      } finally {
        this.saving = false
      }
    },
    preview() {
      if (!this.id) return
      uni.navigateTo({ url: `/pages/property_detail/property_detail?id=${this.id}` })
    },
  },
}
</script>

<style>
.pe {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.icon-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(226, 232, 240, 0.9);
  background: rgba(255, 255, 255, 0.8);
  color: #0f172a;
}

.icon-btn.primary {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
  box-shadow: 0 12rpx 28rpx rgba(37, 99, 235, 0.22);
}

.content {
  flex: 1;
  overflow: hidden;
}

.card {
  margin: 16rpx 24rpx 0;
  padding: 18rpx 18rpx 14rpx;
  background-color: #ffffff;
  border-radius: 28rpx;
  border: 1px solid #eef2f7;
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.05);
}

.card-title {
  font-size: 30rpx;
  font-weight: 800;
  color: #0f172a;
  margin-bottom: 14rpx;
}

.hint {
  font-size: 24rpx;
  color: #94a3b8;
  margin-top: -8rpx;
  margin-bottom: 12rpx;
}
.hint.small {
  margin: 0;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
  margin-bottom: 14rpx;
}

.label {
  font-size: 24rpx;
  color: #334155;
  font-weight: 700;
}
.label.required::after {
  content: ' *';
  color: #ef4444;
}

.input {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #f8fafc;
  padding: 0 16rpx;
  font-size: 28rpx;
  color: #0f172a;
  box-sizing: border-box;
}

.placeholder {
  color: #94a3b8;
}

.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14rpx;
}

.picker {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #f8fafc;
  padding: 0 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.picker-text {
  font-size: 28rpx;
  color: #0f172a;
}

.picker-ic {
  font-size: 36rpx !important;
  color: #94a3b8;
}

.price-row {
  display: flex;
  gap: 10rpx;
  align-items: center;
}

.unit-picker {
  width: 140rpx;
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #f8fafc;
  padding: 0 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
}

.unit-text {
  font-size: 28rpx;
  color: #0f172a;
}
.unit-ic {
  font-size: 34rpx !important;
  color: #94a3b8;
}

.room-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.mini {
  width: 120rpx;
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #f8fafc;
  padding: 0 12rpx;
  font-size: 28rpx;
  color: #0f172a;
  box-sizing: border-box;
}

.mini-suffix {
  font-size: 26rpx;
  color: #64748b;
  margin-right: 8rpx;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12rpx;
}

.media-item {
  height: 200rpx;
  border-radius: 20rpx;
  overflow: hidden;
  position: relative;
  background-color: #f1f5f9;
  border: 1px solid #e2e8f0;
}

.media-item.add {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8rpx;
  background-color: #f8fafc;
}

.media-img {
  width: 100%;
  height: 100%;
}

.badge {
  position: absolute;
  left: 10rpx;
  top: 10rpx;
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  background: rgba(37, 99, 235, 0.85);
  color: #ffffff;
  backdrop-filter: blur(8px);
}

.add-ic {
  font-size: 54rpx !important;
  color: #2563eb;
}

.add-text {
  font-size: 24rpx;
  color: #2563eb;
  font-weight: 700;
}

.switch-row {
  margin-top: 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 10rpx;
  border-top: 1px dashed #eef2f7;
}

.switch-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  color: #334155;
}

.sw-ic {
  font-size: 34rpx !important;
  color: #64748b;
}

.sw-text {
  font-size: 26rpx;
  font-weight: 700;
}

.video-box {
  height: 280rpx;
  border-radius: 24rpx;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  background-color: #0b1220;
}

.video-empty {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
  color: rgba(255, 255, 255, 0.9);
}

.video-ic {
  font-size: 64rpx !important;
  opacity: 0.85;
}

.video-texts {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.video-title {
  font-size: 30rpx;
  font-weight: 800;
}

.video-sub {
  font-size: 24rpx;
  opacity: 0.8;
}

.video {
  width: 100%;
  height: 100%;
}

.tag-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10rpx;
}

.tag {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 14rpx;
  border-radius: 999rpx;
  background-color: #f1f5f9;
  border: 1px solid #e2e8f0;
}

.tag-text {
  font-size: 24rpx;
  color: #0f172a;
  font-weight: 700;
}

.tag-x {
  font-size: 26rpx !important;
  color: #64748b;
}

.tag-add {
  margin-top: 12rpx;
  display: flex;
  gap: 12rpx;
  align-items: center;
}

.btn {
  height: 78rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #ffffff;
  font-size: 26rpx;
  color: #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 18rpx;
}
.btn::after {
  border: none;
}
.btn.add {
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  border: none;
  color: #ffffff;
  font-weight: 800;
  min-width: 140rpx;
}

.btn.ghost {
  background-color: #f8fafc;
  color: #2563eb;
  border-color: rgba(37, 99, 235, 0.22);
  min-width: 160rpx;
}

.coord-row {
  margin-top: 12rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
  padding-top: 10rpx;
  border-top: 1px dashed #eef2f7;
}

.coord-left {
  display: flex;
  align-items: center;
  gap: 10rpx;
  min-width: 0;
}

.coord-ic {
  font-size: 34rpx !important;
  color: #64748b;
}

.coord-text {
  font-size: 24rpx;
  color: #334155;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.bottom-spacer {
  height: calc(env(safe-area-inset-bottom) + 120rpx);
}

.footer {
  padding: 14rpx 24rpx calc(env(safe-area-inset-bottom) + 14rpx);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(12px);
  border-top: 1rpx solid rgba(226, 232, 240, 0.9);
  display: flex;
  gap: 14rpx;
}

.footer-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 24rpx;
  font-size: 28rpx;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
  background-color: #ffffff;
  color: #0f172a;
}
.footer-btn::after {
  border: none;
}

.footer-btn.primary {
  border: none;
  color: #ffffff;
  background: linear-gradient(135deg, #2d9cf0, #2563eb);
  box-shadow: 0 14rpx 28rpx rgba(37, 99, 235, 0.22);
}

.footer-btn.ghost {
  background-color: #f8fafc;
  color: #2563eb;
  border-color: rgba(37, 99, 235, 0.22);
}
</style>

