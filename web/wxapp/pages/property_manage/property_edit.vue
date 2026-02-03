<template>
  <view class="pe">
    <TopHeader :title="id ? '编辑房源' : '新增房源'">
      <template #left>
        <view class="icon-btn" @tap="goBack">
          <text class="material-symbols-outlined">arrow_back</text>
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
                :range="priceUnitOptions.map((o) => o.label)"
                :value="priceUnitIndex"
                @change="onPriceUnitChange"
              >
                <view class="unit-picker">
                  <text class="unit-text">{{
                    priceUnitOptions[priceUnitIndex]?.label || form.price_unit || '万'
                  }}</text>
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

      <!-- 房源参数 -->
      <view class="card">
        <view class="card-title">房源参数</view>

        <view class="grid">
          <view class="form-row">
            <text class="label">建成年份</text>
            <picker
              mode="selector"
              :range="buildYearOptions.map((o) => o.label)"
              :value="buildYearIndex"
              @change="onBuildYearChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  buildYearOptions[buildYearIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">楼层位置</text>
            <picker
              mode="selector"
              :range="floorLevelOptions.map((o) => o.label)"
              :value="floorLevelIndex"
              @change="onFloorLevelChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  floorLevelOptions[floorLevelIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">总楼层</text>
            <input
              v-model="form.total_floors"
              class="input"
              type="number"
              placeholder="共多少层"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">朝向</text>
            <picker
              mode="selector"
              :range="orientationOptions.map((o) => o.label)"
              :value="orientationIndex"
              @change="onOrientationChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  orientationOptions[orientationIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">物业类型</text>
            <picker
              mode="selector"
              :range="propertyTypeOptions.map((o) => o.label)"
              :value="propertyTypeIndex"
              @change="onPropertyTypeChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  propertyTypeOptions[propertyTypeIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">装修标准</text>
            <picker
              mode="selector"
              :range="decorationTypeOptions.map((o) => o.label)"
              :value="decorationTypeIndex"
              @change="onDecorationTypeChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  decorationTypeOptions[decorationTypeIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>
        </view>
      </view>

      <!-- 房主与收房 -->
      <view class="card">
        <view class="card-title">房主与收房</view>
        <view class="grid">
          <view class="form-row">
            <text class="label">房主姓名</text>
            <input
              v-model="form.owner_name"
              class="input"
              type="text"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">房主电话</text>
            <input
              v-model="form.owner_phone"
              class="input"
              type="number"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">收房人姓名</text>
            <input
              v-model="form.receiver_name"
              class="input"
              type="text"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">收房人电话</text>
            <input
              v-model="form.receiver_phone"
              class="input"
              type="number"
              placeholder="可选"
              placeholder-class="placeholder"
            />
          </view>
        </view>
        <view class="form-row">
          <text class="label">收房价格(支付业主)</text>
          <input
            v-model="form.receiver_price"
            class="input"
            type="digit"
            placeholder="未填写"
            placeholder-class="placeholder"
          />
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
            <image class="media-img" :src="toFullMediaUrl(img)" mode="aspectFill"></image>
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
            :src="toFullMediaUrl(form.video_url)"
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
          <text class="label">省市区</text>
          <picker mode="region" :value="addressRegion" @change="onAddressRegionChange">
            <view class="picker">
              <text class="picker-text">{{ addressRegionText }}</text>
              <text class="material-symbols-outlined picker-ic">expand_more</text>
            </view>
          </picker>
        </view>

        <view class="form-row">
          <text class="label">详细地址</text>
          <input
            v-model="addressDetail"
            class="input"
            type="text"
            placeholder="街道、门牌号、楼栋单元等"
            placeholder-class="placeholder"
          />
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

      <!-- 发布与佣金 -->
      <view class="card">
        <view class="card-title">发布与佣金</view>
        <view class="grid">
          <view class="form-row">
            <text class="label">佣金比例(%)</text>
            <input
              v-model="form.commission_rate"
              class="input"
              type="digit"
              placeholder="0-100"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">成交奖励金</text>
            <input
              v-model="form.commission_reward"
              class="input"
              type="digit"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">排序权重</text>
            <input
              v-model="form.weigh"
              class="input"
              type="number"
              placeholder="0"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">智能门锁</text>
            <view class="picker disabled">
              <text class="picker-text">{{
                Number(form.has_smart_lock) === 1 ? '已绑定' : '未绑定'
              }}</text>
              <text class="material-symbols-outlined picker-ic">lock</text>
            </view>
            <view class="hint small tip">门锁请在后台/详情页绑定</view>
          </view>
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">stars</text>
            <text class="sw-text">推荐房源</text>
          </view>
          <switch
            :checked="Number(form.hot_status) === 1"
            @change="(e) => (form.hot_status = e.detail.value ? 1 : 0)"
            color="#2563eb"
          />
        </view>

        <view class="switch-row">
          <view class="switch-left">
            <text class="material-symbols-outlined sw-ic">visibility</text>
            <text class="sw-text">房源启用</text>
          </view>
          <switch
            :checked="Number(form.status) === 0"
            @change="(e) => (form.status = e.detail.value ? 0 : 1)"
            color="#2563eb"
          />
        </view>
      </view>

      <!-- 装修进度（放在最后） -->
      <view class="card">
        <view class="card-title">装修进度</view>
        <view class="hint">可选：用于展示装修进度与阶段</view>

        <view class="form-row">
          <text class="label">装修状态</text>
          <picker
            mode="selector"
            :range="renovationStatusOptions.map((o) => o.label)"
            :value="renovationStatusIndex"
            @change="onRenovationStatusChange"
          >
            <view class="picker">
              <text class="picker-text">{{
                renovationStatusOptions[renovationStatusIndex]?.label || '请选择'
              }}</text>
              <text class="material-symbols-outlined picker-ic">expand_more</text>
            </view>
          </picker>
          <view v-if="!id" class="hint small tip">请先保存房源生成ID后再维护装修进度</view>
        </view>

        <view v-if="id && renovation.renovation_status !== 'none'" class="grid">
          <view class="form-row">
            <text class="label">进度(%)</text>
            <input
              v-model="renovation.progress_percentage"
              class="input"
              type="number"
              placeholder="0-100"
              placeholder-class="placeholder"
            />
          </view>

          <view class="form-row">
            <text class="label">当前阶段</text>
            <picker
              mode="selector"
              :range="renovationStageOptions.map((o) => o.label)"
              :value="renovationStageIndex"
              @change="onRenovationStageChange"
            >
              <view class="picker">
                <text class="picker-text">{{
                  renovationStageOptions[renovationStageIndex]?.label || '请选择'
                }}</text>
                <text class="material-symbols-outlined picker-ic">expand_more</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">开始日期</text>
            <picker
              mode="date"
              :value="renovation.start_date"
              @change="(e) => (renovation.start_date = e.detail.value)"
            >
              <view class="picker">
                <text class="picker-text">{{ renovation.start_date || '请选择' }}</text>
                <text class="material-symbols-outlined picker-ic">calendar_month</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">预计完工</text>
            <picker
              mode="date"
              :value="renovation.estimated_finish_date"
              @change="(e) => (renovation.estimated_finish_date = e.detail.value)"
            >
              <view class="picker">
                <text class="picker-text">{{ renovation.estimated_finish_date || '请选择' }}</text>
                <text class="material-symbols-outlined picker-ic">event</text>
              </view>
            </picker>
          </view>

          <view class="form-row">
            <text class="label">实际完工</text>
            <picker
              mode="date"
              :value="renovation.actual_finish_date"
              @change="(e) => (renovation.actual_finish_date = e.detail.value)"
            >
              <view class="picker">
                <text class="picker-text">{{ renovation.actual_finish_date || '请选择' }}</text>
                <text class="material-symbols-outlined picker-ic">verified</text>
              </view>
            </picker>
          </view>
        </view>

        <view v-if="id && renovation.renovation_status !== 'none'">
          <view class="form-row">
            <text class="label">材料</text>
            <view class="tag-row">
              <view class="tag" v-for="(t, idx) in renovationMaterials" :key="'m-' + idx">
                <text class="tag-text">{{ t }}</text>
                <text class="material-symbols-outlined tag-x" @tap.stop="removeMaterial(idx)">close</text>
              </view>
              <view v-if="renovationMaterials.length === 0" class="hint small">暂无材料</view>
            </view>
            <view class="tag-add">
              <input
                v-model="newMaterial"
                class="input"
                type="text"
                maxlength="12"
                placeholder="输入材料（最多10个）"
                placeholder-class="placeholder"
              />
              <button class="btn add" @tap="addMaterial">添加</button>
            </view>
          </view>

          <view class="form-row">
            <text class="label">装修图片</text>
            <view class="media-grid">
              <view
                v-for="(img, idx) in renovationImages"
                :key="'rimg-' + idx"
                class="media-item"
                @tap="openRenovationImageActions(idx)"
              >
                <image class="media-img" :src="toFullMediaUrl(img)" mode="aspectFill"></image>
              </view>
              <view class="media-item add" @tap="pickRenovationImages">
                <text class="material-symbols-outlined add-ic">add_photo_alternate</text>
                <text class="add-text">添加</text>
              </view>
            </view>
          </view>

          <view class="form-row">
            <text class="label">施工说明</text>
            <textarea
              v-model="renovation.notes"
              class="textarea"
              maxlength="500"
              placeholder="可选：记录进度说明、施工要点等"
              placeholder-class="placeholder"
            ></textarea>
          </view>

          <view class="tag-add">
            <button class="btn ghost" :disabled="savingRenovation || uploading" @tap="loadRenovation">
              刷新
            </button>
            <button class="btn add" :disabled="savingRenovation || uploading" @tap="saveRenovation">
              {{ savingRenovation ? '保存中...' : '保存装修' }}
            </button>
          </view>
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
    const curYear = new Date().getFullYear()
    const buildYearOptions = [{ label: '请选择', value: '' }].concat(
      Array.from({ length: curYear - 1950 + 1 }).map((_, i) => {
        const y = String(curYear - i)
        return { label: y, value: y }
      })
    )

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
        build_year: '',
        floor_level: '',
        total_floors: '',
        orientation: '',
        property_type: '',
        decoration_type: '',
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
        owner_name: '',
        owner_phone: '',
        receiver_name: '',
        receiver_phone: '',
        receiver_price: '',
        commission_rate: '',
        commission_reward: '',
        weigh: '',
        hot_status: 0,
        has_smart_lock: 0,
        status: 0,
      },
      tags: [],
      images: [],
      newTag: '',
      // 地址分段输入：省市区 + 详细地址（最终保存时合成到 form.address）
      addressRegion: ['', '', ''],
      addressDetail: '',
      saleStatusOptions: [
        { label: '在售', value: 'on_sale' },
        { label: '已售', value: 'sold' },
        { label: '下架', value: 'off_market' },
      ],
      buildYearOptions,
      floorLevelOptions: [
        { label: '低层 (1-6)', value: '低层' },
        { label: '中层 (7-15)', value: '中层' },
        { label: '高层 (16+)', value: '高层' },
        { label: '地下', value: '地下' },
      ],
      orientationOptions: [
        '东',
        '南',
        '西',
        '北',
        '东南',
        '东北',
        '西南',
        '西北',
        '南北',
        '东西',
      ].map((it) => ({ label: it, value: it })),
      propertyTypeOptions: ['住宅', '公寓', '别墅', '商铺', '写字楼'].map((it) => ({
        label: it,
        value: it,
      })),
      decorationTypeOptions: ['毛坯', '简装', '精装', '豪装'].map((it) => ({
        label: it,
        value: it,
      })),
      priceUnitOptions: [
        { label: '万', value: '万' },
        { label: '元', value: '元' },
      ],
      // 装修进度（可维护房源）
      savingRenovation: false,
      renovation: {
        renovation_status: 'none',
        progress_percentage: 0,
        current_stage: '',
        start_date: '',
        estimated_finish_date: '',
        actual_finish_date: '',
        notes: '',
        status: 0,
      },
      renovationStatusOptions: [
        { label: '未装修', value: 'none' },
        { label: '装修中', value: 'in_progress' },
        { label: '已完成', value: 'done' },
      ],
      renovationStageOptions: [{ label: '请选择', value: '' }].concat(
        ['设计', '拆改', '水电', '泥瓦', '木工', '油漆', '安装', '软装', '验收'].map((it) => ({ label: it, value: it }))
      ),
      renovationMaterials: [],
      renovationImages: [],
      newMaterial: '',
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
      const idx = this.priceUnitOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    buildYearIndex() {
      const v = String(this.form.build_year || '')
      const idx = this.buildYearOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    floorLevelIndex() {
      const v = String(this.form.floor_level || '')
      const idx = this.floorLevelOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    orientationIndex() {
      const v = String(this.form.orientation || '')
      const idx = this.orientationOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    propertyTypeIndex() {
      const v = String(this.form.property_type || '')
      const idx = this.propertyTypeOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    decorationTypeIndex() {
      const v = String(this.form.decoration_type || '')
      const idx = this.decorationTypeOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    addressRegionText() {
      const arr = Array.isArray(this.addressRegion) ? this.addressRegion : []
      const text = arr.filter((it) => String(it || '').trim()).join(' / ')
      return text || '请选择省市区'
    },
    renovationStatusIndex() {
      const v = String(this.renovation.renovation_status || 'none')
      const idx = this.renovationStatusOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
    renovationStageIndex() {
      const v = String(this.renovation.current_stage || '')
      const idx = this.renovationStageOptions.findIndex((o) => o.value === v)
      return idx >= 0 ? idx : 0
    },
  },
  onLoad(options) {
    const id = Number((options && (options.id || options.ID)) || 0) || 0
    this.id = id
    this.bootstrap()
  },
  methods: {
    // 兼容：上传接口可能返回相对路径（/resource/uploads/...），这里统一转为可展示的完整地址
    toFullMediaUrl(url) {
      const s = String(url || '').trim()
      if (!s) return ''
      if (/^https?:\/\//i.test(s) || s.startsWith('wxfile://') || s.startsWith('file://')) return s
      const root = String(baseUrl || '').replace(/\/+$/, '')
      const path = s.startsWith('/') ? s : `/${s}`
      return root ? `${root}${path}` : path
    },
    goBack() {
      uni.navigateBack()
    },
    async ensurePermission() {
      const userStore = $store('user')
      const token = uni.getStorageSync('token')
      if (!token && !userStore.isLogin) {
        uni.reLaunch({ url: '/pages/login/login' })
        return false
      }
      try {
        const ui = userStore.userInfo || {}
        // 兼容：老版本缓存的 userInfo 可能没有 can_manage_properties，需要强制刷新一次
        if (!ui.id || typeof ui.can_manage_properties === 'undefined') {
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
        return false
      }
      return true
    },
    onSaleStatusChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.saleStatusOptions[idx] || this.saleStatusOptions[0]
      this.form.sale_status = opt.value
    },
    onPriceUnitChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.priceUnitOptions[idx] || this.priceUnitOptions[0]
      this.form.price_unit = (opt && opt.value) || '万'
    },
    // 页面初始化：权限校验 + 下拉选项缓存 + 回显
    async bootstrap() {
      const ok = await this.ensurePermission()
      if (!ok) return
      await this.loadFormOptions()
      if (this.id) {
        await this.loadContent()
        await this.loadRenovation()
      }
    },
    // 首次加载表单下拉项时写入缓存；有缓存直接使用
    async loadFormOptions() {
      const userStore = $store('user')
      const businessId = Number((userStore.userInfo || {}).business_id) || 1
      const cacheKey = `wx_property_form_options_v2_${businessId}`
      const cached = uni.getStorageSync(cacheKey)
      if (cached && typeof cached === 'object') {
        this.applyFormOptions(cached)
        return
      }

      const res = await propertyApi.getFormOptions({})
      if (!res || res.code !== 0) return
      const data = res.data || {}
      this.applyFormOptions(data)
      try {
        uni.setStorageSync(cacheKey, data)
      } catch (e) {}
    },
    applyFormOptions(data) {
      const normalize = (arr) =>
        (Array.isArray(arr) ? arr : [])
          .map((it) => ({
            label: String((it && it.label) || '').trim(),
            value: String((it && it.value) || '').trim(),
          }))
          .filter((it) => it.label && it.value)

      const sale = normalize(data.sale_status)
      if (sale.length) this.saleStatusOptions = sale

      const unit = normalize(data.price_unit)
      if (unit.length) this.priceUnitOptions = unit

      const floor = normalize(data.floor_level)
      if (floor.length) this.floorLevelOptions = floor

      const ori = normalize(data.orientation)
      if (ori.length) this.orientationOptions = ori

      const pt = normalize(data.property_type)
      if (pt.length) this.propertyTypeOptions = pt

      const deco = normalize(data.decoration_type)
      if (deco.length) this.decorationTypeOptions = deco

      const stages = normalize(data.renovation_stage)
      if (stages.length) this.renovationStageOptions = [{ label: '请选择', value: '' }].concat(stages)

      // 兜底：当前值不在选项内时，回退到第一个选项
      const fixValue = (key, options, fallback) => {
        const v = String(this.form[key] || '').trim()
        const hit = (options || []).some((o) => String(o.value) === v)
        if (!hit) this.form[key] = (options && options[0] && options[0].value) || fallback
      }
      fixValue('sale_status', this.saleStatusOptions, 'on_sale')
      fixValue('price_unit', this.priceUnitOptions, '万')
      fixValue('floor_level', this.floorLevelOptions, '')
      fixValue('orientation', this.orientationOptions, '')
      fixValue('property_type', this.propertyTypeOptions, '')
      fixValue('decoration_type', this.decorationTypeOptions, '')

      // 装修工序：允许历史数据不在选项内时追加显示，避免 picker 显示异常
      const stageVal = String(this.renovation.current_stage || '').trim()
      if (stageVal && !(this.renovationStageOptions || []).some((o) => String(o.value) === stageVal)) {
        this.renovationStageOptions = (this.renovationStageOptions || []).concat([{ label: stageVal, value: stageVal }])
      }
    },
    onBuildYearChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.buildYearOptions[idx] || this.buildYearOptions[0]
      this.form.build_year = opt.value
    },
    onFloorLevelChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.floorLevelOptions[idx] || this.floorLevelOptions[0]
      this.form.floor_level = opt.value
    },
    onOrientationChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.orientationOptions[idx] || this.orientationOptions[0]
      this.form.orientation = opt.value
    },
    onPropertyTypeChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.propertyTypeOptions[idx] || this.propertyTypeOptions[0]
      this.form.property_type = opt.value
    },
    onDecorationTypeChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.decorationTypeOptions[idx] || this.decorationTypeOptions[0]
      this.form.decoration_type = opt.value
    },
    onAddressRegionChange(e) {
      const v = (e && e.detail && e.detail.value) || []
      if (Array.isArray(v)) this.addressRegion = v
    },
    buildFullAddress() {
      const region = Array.isArray(this.addressRegion)
        ? this.addressRegion.map((it) => String(it || '').trim()).filter(Boolean)
        : []
      const detail = String(this.addressDetail || '').trim()
      return [...region, detail].filter(Boolean).join(' ')
    },
    initAddressUIFromAddress(address) {
      const raw = String(address || '').trim()
      this.addressRegion = ['', '', '']
      this.addressDetail = ''
      if (!raw) return
      const parts = raw.split(/\\s+/).filter(Boolean)
      if (parts.length >= 3) {
        this.addressRegion = [parts[0] || '', parts[1] || '', parts[2] || '']
        this.addressDetail = parts.slice(3).join(' ')
        return
      }
      this.addressDetail = raw
    },
    onRenovationStageChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.renovationStageOptions[idx] || this.renovationStageOptions[0]
      this.renovation.current_stage = (opt && opt.value) || ''
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
      if (this.form.build_year) this.form.build_year = String(this.form.build_year)
      if (this.form.total_floors) this.form.total_floors = String(this.form.total_floors)
      if (this.form.weigh || this.form.weigh === 0) this.form.weigh = String(this.form.weigh)
      if (this.form.receiver_price === null || typeof this.form.receiver_price === 'undefined') {
        this.form.receiver_price = ''
      }
      if (this.form.commission_rate === null || typeof this.form.commission_rate === 'undefined') {
        this.form.commission_rate = ''
      }
      if (
        this.form.commission_reward === null ||
        typeof this.form.commission_reward === 'undefined'
      ) {
        this.form.commission_reward = ''
      }
      // 地址拆分：用于省市区选择 + 详细地址输入
      this.initAddressUIFromAddress(this.form.address)
      // 兜底：封面未设置时，用第一张图
      if (!this.form.cover_image && this.images.length > 0) {
        this.form.cover_image = this.images[0]
      }
    },
    async loadRenovation() {
      if (!this.id) return
      const res = await propertyApi.getManageRenovation({ property_id: this.id })
      if (!res || res.code !== 0) return
      const data = res.data || {}
      this.renovation = {
        ...this.renovation,
        ...data,
      }
      this.renovationMaterials = Array.isArray(data.materials) ? data.materials : []
      this.renovationImages = Array.isArray(data.images) ? data.images : []

      // 兜底类型修正
      if (this.renovation.progress_percentage === null || typeof this.renovation.progress_percentage === 'undefined') {
        this.renovation.progress_percentage = 0
      }
      if (this.renovation.renovation_status !== 'none' && this.renovation.renovation_status !== 'in_progress' && this.renovation.renovation_status !== 'done') {
        this.renovation.renovation_status = 'none'
      }

      // 兜底：阶段值不在下拉里时，追加一项用于回显
      const stageVal = String(this.renovation.current_stage || '').trim()
      if (stageVal && !(this.renovationStageOptions || []).some((o) => String(o.value) === stageVal)) {
        this.renovationStageOptions = (this.renovationStageOptions || []).concat([{ label: stageVal, value: stageVal }])
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
    onRenovationStatusChange(e) {
      const idx = Number(e && e.detail && e.detail.value) || 0
      const opt = this.renovationStatusOptions[idx] || this.renovationStatusOptions[0]
      this.renovation.renovation_status = opt.value
      if (opt.value === 'none') {
        this.renovation.progress_percentage = 0
        this.renovation.current_stage = ''
        this.renovation.start_date = ''
        this.renovation.estimated_finish_date = ''
        this.renovation.actual_finish_date = ''
        this.renovation.notes = ''
        this.renovationMaterials = []
        this.renovationImages = []
      }
      if (opt.value === 'done') {
        this.renovation.progress_percentage = 100
      }
    },
    addMaterial() {
      const t = String(this.newMaterial || '').trim()
      if (!t) {
        uni.showToast({ title: '请输入材料', icon: 'none' })
        return
      }
      if (this.renovationMaterials.includes(t)) {
        this.newMaterial = ''
        return
      }
      if (this.renovationMaterials.length >= 10) {
        uni.showToast({ title: '最多10个材料', icon: 'none' })
        return
      }
      this.renovationMaterials = this.renovationMaterials.concat([t])
      this.newMaterial = ''
    },
    removeMaterial(idx) {
      this.renovationMaterials = this.renovationMaterials.filter((_, i) => i !== idx)
    },
    openRenovationImageActions(idx) {
      const img = this.renovationImages[idx]
      if (!img) return
      uni.showActionSheet({
        itemList: ['删除'],
        success: () => {
          this.renovationImages = this.renovationImages.filter((_, i) => i !== idx)
        },
      })
    },
    pickRenovationImages() {
      if (this.uploading) return
      if (!this.id) {
        uni.showToast({ title: '请先保存房源', icon: 'none' })
        return
      }
      uni.chooseImage({
        count: 9,
        sizeType: ['compressed'],
        sourceType: ['album', 'camera'],
        success: async (res) => {
          const files = (res && res.tempFilePaths) || []
          if (!files.length) return
          await this.uploadFiles(files, 'image', 'renovation')
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
          await this.uploadFiles(files, 'image', 'property')
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
            await this.uploadFiles([filePath], 'video', 'property')
            resolve()
          },
          fail: () => resolve(),
        })
      })
    },
    uploadFiles(filePaths = [], filetype = 'image', target = 'property') {
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
              if (target === 'renovation') {
                this.renovationImages = this.renovationImages.concat([url])
              } else {
                this.images = this.images.concat([url])
                if (!this.form.cover_image) this.form.cover_image = url
              }
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
          if (res.address && !String(this.addressDetail || '').trim()) this.addressDetail = String(res.address)
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
          address: this.buildFullAddress(),
          latitude: this.form.latitude || '',
          longitude: this.form.longitude || '',
          tags: this.tags,
          images: this.images,
          cover_image: this.form.cover_image || '',
          video_url: this.form.video_url || '',
          allow_image_download: Number(this.form.allow_image_download) === 1 ? 1 : 0,
          allow_video_download: Number(this.form.allow_video_download) === 1 ? 1 : 0,
          build_year: this.form.build_year ? Number(this.form.build_year) : 0,
          floor_level: String(this.form.floor_level || '').trim(),
          total_floors: Number(this.form.total_floors || 0),
          orientation: String(this.form.orientation || '').trim(),
          property_type: String(this.form.property_type || '').trim(),
          decoration_type: String(this.form.decoration_type || '').trim(),
          owner_name: String(this.form.owner_name || '').trim(),
          owner_phone: String(this.form.owner_phone || '').trim(),
          receiver_name: String(this.form.receiver_name || '').trim(),
          receiver_phone: String(this.form.receiver_phone || '').trim(),
          receiver_price: this.form.receiver_price,
          commission_rate: this.form.commission_rate,
          commission_reward: this.form.commission_reward,
          weigh: Number(this.form.weigh || 0),
          hot_status: Number(this.form.hot_status) === 1 ? 1 : 0,
          status: Number(this.form.status) === 0 ? 0 : 1,
        }
        const res = await propertyApi.saveManage(payload)
        if (!res || res.code !== 0) return
        const nid = Number(res && res.data && res.data.id) || 0
        if (!this.id && nid) {
          this.id = nid
          await this.loadRenovation()
        }
        uni.showToast({ title: '保存成功', icon: 'none' })
      } finally {
        this.saving = false
      }
    },
    async saveRenovation() {
      if (!this.id) {
        uni.showToast({ title: '请先保存房源', icon: 'none' })
        return
      }
      if (this.savingRenovation || this.uploading) return

      const status = String(this.renovation.renovation_status || 'none')
      const progress = Number(this.renovation.progress_percentage || 0)
      if (!isFinite(progress) || progress < 0 || progress > 100) {
        uni.showToast({ title: '进度需在0-100之间', icon: 'none' })
        return
      }

      this.savingRenovation = true
      try {
        const payload = {
          property_id: this.id,
          renovation_status: status,
          progress_percentage: Math.floor(progress),
          current_stage: String(this.renovation.current_stage || '').trim(),
          start_date: this.renovation.start_date || '',
          estimated_finish_date: this.renovation.estimated_finish_date || '',
          actual_finish_date: this.renovation.actual_finish_date || '',
          materials: this.renovationMaterials,
          images: this.renovationImages,
          notes: String(this.renovation.notes || '').trim(),
          status: Number(this.renovation.status) === 1 ? 1 : 0,
        }
        const res = await propertyApi.saveManageRenovation(payload)
        if (!res || res.code !== 0) return
        uni.showToast({ title: '保存成功', icon: 'none' })
        await this.loadRenovation()
      } finally {
        this.savingRenovation = false
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
.hint.small.tip {
  margin-top: 10rpx;
  color: #64748b;
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

.textarea {
  min-height: 160rpx;
  border-radius: 20rpx;
  border: 1px solid #e2e8f0;
  background-color: #f8fafc;
  padding: 14rpx 16rpx;
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
.picker.disabled {
  background-color: #f1f5f9;
  border-color: #e2e8f0;
  opacity: 0.9;
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

