<template>
  <a-modal
    v-model:visible="visibleProxy"
    title="地图选点"
    width="980px"
    :footer="false"
    :maskClosable="true"
    :unmountOnClose="true"
    @close="handleClose"
  >
    <div class="picker-wrap">
      <div class="picker-toolbar">
        <a-space wrap>
          <a-input
            v-model="keyword"
            style="width: 360px"
            placeholder="输入地址/关键字（可选，可用腾讯地图检索定位）"
            allow-clear
          >
            <template #prefix><icon-location /></template>
          </a-input>
          <a-button :loading="searching" @click="searchByKeyword">
            <template #icon><icon-search /></template>
            搜索
          </a-button>
          <a-button @click="locateByBrowser">
            <template #icon><icon-compass /></template>
            使用当前位置
          </a-button>
          <a-button @click="clearPick">清空坐标</a-button>
        </a-space>

        <div class="picker-toolbar-right">
          <a-space>
            <div class="coord-pill" v-if="pickedText">
              <span class="label">已选</span>
              <span class="val">{{ pickedText }}</span>
            </div>
            <a-button @click="visibleProxy = false">取消</a-button>
            <a-button type="primary" :disabled="!picked" @click="applyPick">使用该坐标</a-button>
          </a-space>
        </div>
      </div>

      <div class="picker-body">
        <div class="map-shell">
          <div class="map" ref="mapEl"></div>
          <div class="map-hint">
            <div class="hint-title">提示</div>
            <div class="hint-text">在地图上单击即可选点并回填经纬度</div>
          </div>
          <div class="map-error" v-if="loadError">
            <div class="err-title">地图加载失败</div>
            <div class="err-text">{{ loadError }}</div>
            <div class="err-text">请在 `web/business/public/config.js` 配置 `TencentMapKey` 后重试</div>
          </div>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script lang="ts">
import { defineComponent, ref, computed, watch, nextTick, onBeforeUnmount } from 'vue';
import { Message } from '@arco-design/web-vue';

declare global {
  interface Window {
    qq?: any;
    globalConfig?: any;
  }
}

type Picked = { lng: number; lat: number } | null;

export default defineComponent({
  name: 'LatLngPickerModal',
  props: {
    visible: { type: Boolean, default: false },
    latitude: { type: [String, Number], default: '' },
    longitude: { type: [String, Number], default: '' },
  },
  emits: ['update:visible', 'confirm', 'clear'],
  setup(props, { emit }) {
    const mapEl = ref<HTMLDivElement | null>(null);
    const keyword = ref('');
    const picked = ref<Picked>(null);
    const searching = ref(false);
    const loadError = ref('');

    const visibleProxy = computed({
      get: () => props.visible,
      set: (v: boolean) => emit('update:visible', v),
    });

    const pickedText = computed(() => {
      if (!picked.value) return '';
      return `${picked.value.lng.toFixed(6)}, ${picked.value.lat.toFixed(6)}`;
    });

    const parseNum = (v: any) => {
      const n = Number(v);
      return Number.isFinite(n) ? n : NaN;
    };

    const guessCenter = () => {
      const lng = parseNum(props.longitude);
      const lat = parseNum(props.latitude);
      if (Number.isFinite(lng) && Number.isFinite(lat)) return { lng, lat, zoom: 16 };
      // 默认：北京
      return { lng: 116.4074, lat: 39.9042, zoom: 11 };
    };

    const getTencentKey = () => {
      // 兼容配置入口：public/config.js
      const fromGlobal = window?.globalConfig?.TencentMapKey;
      return (fromGlobal || '').toString().trim();
    };

    const ensureTencentLoaded = async () => {
      if (window.qq?.maps) return window.qq;
      const key = getTencentKey();
      console.log(key);
      if (!key) throw new Error('未配置 TencentMapKey');
      const src = `https://map.qq.com/api/js?v=2.exp&key=${encodeURIComponent(key)}&libraries=service`;
      await new Promise<void>((resolve, reject) => {
        const exists = Array.from(document.querySelectorAll('script')).some((s) => (s as any).src?.includes('map.qq.com/api/js'));
        if (exists) return resolve();
        const script = document.createElement('script');
        script.src = src;
        script.async = true;
        script.onload = () => resolve();
        script.onerror = () => reject(new Error('腾讯地图 SDK 加载失败'));
        document.body.appendChild(script);
      });
      if (!window.qq?.maps) throw new Error('腾讯地图 SDK 初始化失败');
      return window.qq;
    };

    let mapIns: any = null;
    let markerIns: any = null;
    let searchService: any = null;

    const destroyMap = () => {
      try {
        if (window.qq?.maps?.event && mapIns) {
          try {
            window.qq.maps.event.clearListeners(mapIns, 'click');
          } catch {}
        }
      } catch {}
      mapIns = null;
      markerIns = null;
      searchService = null;
    };

    const setMarker = (lng: number, lat: number) => {
      const qq = window.qq;
      if (!qq?.maps || !mapIns) return;
      const ll = new qq.maps.LatLng(lat, lng);
      if (!markerIns) {
        markerIns = new qq.maps.Marker({ map: mapIns, position: ll });
      } else {
        markerIns.setPosition(ll);
      }
      mapIns.setCenter(ll);
    };

    const initMap = async () => {
      const el = mapEl.value;
      if (!el) return;
      destroyMap();
      loadError.value = '';
      const qq = await ensureTencentLoaded();
      if (!qq?.maps) return;

      const center = guessCenter();
      mapIns = new qq.maps.Map(el, {
        center: new qq.maps.LatLng(center.lat, center.lng),
        zoom: center.zoom,
      });

      // 关键字检索服务（可选）
      try {
        searchService = new qq.maps.SearchService({
          map: mapIns,
          pageCapacity: 1,
          complete: (results: any) => {
            const list = results?.detail?.pois || results?.detail?.results || [];
            const first = list?.[0];
            const ll = first?.latLng || first?.latlng || first?.location;
            if (!ll) {
              Message.warning('未找到匹配地点');
              return;
            }
            const lat = Number(ll.getLat ? ll.getLat() : ll.lat);
            const lng = Number(ll.getLng ? ll.getLng() : ll.lng);
            if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;
            picked.value = { lng, lat };
            setMarker(lng, lat);
            mapIns.setZoom(Math.max(mapIns.getZoom?.() || 16, 14));
          },
          error: () => Message.warning('搜索失败，请稍后重试'),
        });
      } catch {}

      const lng = parseNum(props.longitude);
      const lat = parseNum(props.latitude);
      if (Number.isFinite(lng) && Number.isFinite(lat)) {
        picked.value = { lng, lat };
        setMarker(lng, lat);
      }

      qq.maps.event.addListener(mapIns, 'click', (e: any) => {
        const ll = e?.latLng;
        if (!ll) return;
        const lat = Number(ll.getLat());
        const lng = Number(ll.getLng());
        if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;
        picked.value = { lng, lat };
        setMarker(lng, lat);
      });
    };

    const locateByBrowser = () => {
      if (!navigator.geolocation) {
        Message.warning('当前浏览器不支持定位');
        return;
      }
      navigator.geolocation.getCurrentPosition(
        (pos) => {
          const lat = Number(pos.coords.latitude);
          const lng = Number(pos.coords.longitude);
          if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;
          picked.value = { lng, lat };
          setMarker(lng, lat);
        },
        () => Message.warning('定位失败，请检查浏览器权限'),
        { enableHighAccuracy: true, timeout: 8000 }
      );
    };

    const clearPick = () => {
      picked.value = null;
      emit('clear');
      if (markerIns && mapIns) {
        try {
          mapIns.removeLayer(markerIns);
        } catch {}
        markerIns = null;
      }
    };

    const searchByKeyword = async () => {
      const q = keyword.value.trim();
      if (!q) {
        Message.warning('请输入关键字');
        return;
      }
      if (!searchService) {
        Message.warning('地图未就绪或未配置腾讯地图 Key');
        return;
      }
      try {
        searching.value = true;
        searchService.search(q);
      } finally {
        setTimeout(() => (searching.value = false), 300);
      }
    };

    const applyPick = () => {
      if (!picked.value) return;
      emit('confirm', {
        longitude: picked.value.lng.toFixed(6),
        latitude: picked.value.lat.toFixed(6),
        keyword: keyword.value.trim(),
      });
      visibleProxy.value = false;
    };

    const handleClose = () => {
      // unmount-on-close=true 时会自动销毁 DOM，这里也做一次主动释放
      destroyMap();
    };

    watch(
      () => props.visible,
      async (v) => {
        if (!v) return;
        await nextTick();
        try {
          await initMap();
        } catch (e: any) {
          loadError.value = e?.message || '未知错误';
        }
      }
    );

    onBeforeUnmount(() => destroyMap());

    return {
      mapEl,
      keyword,
      picked,
      pickedText,
      searching,
      loadError,
      visibleProxy,
      searchByKeyword,
      locateByBrowser,
      clearPick,
      applyPick,
      handleClose,
    };
  },
});
</script>

<style lang="less" scoped>
.picker-wrap {
  padding: 12px 16px 16px;
}
.picker-toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}
.picker-toolbar-right {
  display: flex;
  justify-content: flex-end;
}
.coord-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 999px;
  background: var(--color-fill-2);
  border: 1px solid var(--color-border-2);
  .label {
    font-size: 12px;
    color: var(--color-text-3);
  }
  .val {
    font-size: 12px;
    font-weight: 600;
    color: var(--color-text-1);
    font-variant-numeric: tabular-nums;
  }
}
.map-shell {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--color-border-2);
  background: var(--color-bg-2);
}
.map {
  width: 100%;
  height: 520px;
}
.map-hint {
  position: absolute;
  left: 12px;
  bottom: 12px;
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(0, 0, 0, 0.06);
  backdrop-filter: blur(6px);
  .hint-title {
    font-size: 12px;
    font-weight: 700;
    color: var(--color-text-1);
  }
  .hint-text {
    margin-top: 2px;
    font-size: 12px;
    color: var(--color-text-3);
  }
}

.map-error {
  position: absolute;
  inset: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px dashed var(--color-border-2);
  text-align: center;
  padding: 16px;
  .err-title {
    font-size: 14px;
    font-weight: 700;
    color: var(--color-text-1);
  }
  .err-text {
    font-size: 12px;
    color: var(--color-text-3);
    word-break: break-all;
  }
}
</style>
