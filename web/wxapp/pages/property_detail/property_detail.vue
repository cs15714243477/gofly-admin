<template>
  <view class="detail-container">
    <!-- 顶部悬浮操作 -->
    <view class="float-header" :style="{ paddingTop: headerTop + 'px' }">
      <!-- 安全遮罩：顶部时用暗色渐变保证可读性；下滑后过渡到磨砂白底 -->
      <view
        class="fh-mask fh-mask--dark"
        :style="{ opacity: 1 - headerOpacity }"
      ></view>
      <view
        class="fh-mask fh-mask--light"
        :style="{ opacity: headerOpacity }"
      ></view>
      <view class="header-left">
        <view
          class="circle-btn"
          :class="{ 'is-solid': headerOpacity > 0.65 }"
          @click="goBack"
        >
          <text class="material-symbols-outlined">arrow_back</text>
        </view>
      </view>
    </view>

    <scroll-view scroll-y="true" class="detail-scroll" @scroll="onPageScroll">
      <!-- 轮播图 -->
      <view class="banner">
        <swiper
          class="swiper"
          circular
          :indicator-dots="false"
          @change="swiperChange"
        >
          <swiper-item v-for="(item, index) in images" :key="index">
            <image :src="item" mode="aspectFill" class="banner-image"></image>
          </swiper-item>
        </swiper>
        <view class="banner-indicator"
          >{{ images.length ? currentSwiper + 1 : 0 }}/{{ images.length }}</view
        >
      </view>

      <!-- 内容卡片 -->
      <view class="content-card">
        <view class="title-section">
          <view class="title-row">
            <text class="title">{{ (property && property.title) || "-" }}</text>
            <view class="share-btn" @click="handleShare">
              <text class="material-symbols-outlined share-icon">share</text>
              <text>分享</text>
            </view>
          </view>
          <view class="tags-row">
            <text
              v-for="(t, idx) in ((property && property.tags) || []).slice(
                0,
                6
              )"
              :key="idx"
              class="tag"
              :class="tagClass(t)"
            >
              {{ t }}
            </text>
          </view>
        </view>

        <!-- 核心数据展示 -->
        <view class="stats-grid">
          <view class="stats-item">
            <view class="stats-val orange">
              {{ (property && property.price) || "-"
              }}<text class="unit">{{
                (property && property.price_unit) || ""
              }}</text>
            </view>
            <view class="stats-label">售价</view>
          </view>
          <view class="stats-item border-l">
            <view class="stats-val">{{ getLayoutText(property) }}</view>
            <view class="stats-label">户型</view>
          </view>
          <view class="stats-item border-l">
            <view class="stats-val">
              {{ (property && property.area) || "-"
              }}<text class="unit">㎡</text>
            </view>
            <view class="stats-label">面积</view>
          </view>
        </view>

        <!-- 佣金横幅 -->
        <view class="commission-banner">
          <view class="banner-bg-decor"></view>
          <view class="banner-main">
            <view class="icon-box">
              <text class="material-symbols-outlined">currency_yen</text>
            </view>
            <view class="banner-info">
              <text class="banner-title">{{
                getCommissionText(property)
              }}</text>
              <text class="banner-tip">签约后7个工作日内结算</text>
            </view>
          </view>
        </view>

        <!-- 浏览数据 -->
        <view class="data-row">
          <view class="data-group">
            <view class="data-item">
              <text class="val">{{
                (property && property.view_count) || 0
              }}</text>
              <text class="lab">浏览</text>
            </view>
            <view class="data-item">
              <text class="val">{{
                (property && property.follow_count) || 0
              }}</text>
              <text class="lab">关注人数</text>
            </view>
            <view class="data-item">
              <text class="val">{{
                (property && property.showing_count) || 0
              }}</text>
              <text class="lab">带看次数</text>
            </view>
          </view>
        </view>

        <!-- 属性表 -->
        <view class="section">
          <view class="section-title">房源信息</view>
          <view class="attr-grid">
            <view
              class="attr-item"
              v-for="(attr, idx) in attributes"
              :key="idx"
            >
              <text class="attr-label">{{ attr.label }}</text>
              <text class="attr-val">{{ attr.value }}</text>
            </view>
          </view>
        </view>

        <!-- 装修状态 -->
        <view class="section">
          <view class="section-header">
            <text class="section-title">装修状态</text>
            <view class="reno-tabs">
              <view
                class="reno-tab"
                v-for="t in renovationTabs"
                :key="t.key"
                :class="{ active: renovation.status === t.key }"
                @click="setRenovationStatus(t.key)"
              >
                {{ t.label }}
              </view>
            </view>
          </view>

          <view class="reno-card">
            <view class="reno-top">
              <view class="reno-badge" :class="renovation.status">
                {{
                  renovation.status === "none"
                    ? "未装修"
                    : renovation.status === "in_progress"
                    ? "装修进行中"
                    : "装修完成"
                }}
              </view>
              <text class="reno-sub">{{ renovation.subtitle }}</text>
            </view>

            <!-- 未装修 -->
            <view v-if="renovation.status === 'none'" class="reno-empty">
              <text class="material-symbols-outlined reno-empty-icon"
                >home_repair_service</text
              >
              <view class="reno-empty-texts">
                <text class="reno-empty-title">当前房源未进行装修</text>
                <text class="reno-empty-desc"
                  >适合按个人喜好自由设计，交付为毛坯/简装（示例）。</text
                >
              </view>
            </view>

            <!-- 装修进行中 -->
            <view
              v-else-if="renovation.status === 'in_progress'"
              class="reno-progress-wrap"
            >
              <swiper
                class="reno-swiper"
                circular
                :indicator-dots="true"
                :autoplay="true"
                :interval="3500"
              >
                <swiper-item v-for="(img, idx) in renovation.images" :key="idx">
                  <image
                    :src="img"
                    mode="aspectFill"
                    class="reno-image"
                  ></image>
                </swiper-item>
              </swiper>

              <view class="reno-progress">
                <view class="reno-progress-bar">
                  <view
                    class="reno-progress-fill"
                    :style="{ width: renovation.progress + '%' }"
                  ></view>
                </view>
                <text class="reno-progress-text"
                  >{{ renovation.progress }}% · {{ renovation.stage }}</text
                >
              </view>

              <view class="reno-info">
                <view class="reno-info-row">
                  <text class="k">预计完工</text>
                  <text class="v">{{ renovation.eta }}</text>
                </view>
                <view class="reno-info-row">
                  <text class="k">材料</text>
                  <view class="chips">
                    <text
                      class="chip"
                      v-for="(m, idx) in renovation.materials"
                      :key="'m-' + idx"
                      >{{ m }}</text
                    >
                  </view>
                </view>
                <view class="reno-info-row col">
                  <text class="k">施工说明</text>
                  <text class="v desc">{{ renovation.note }}</text>
                </view>
              </view>
            </view>

            <!-- 装修完成 -->
            <view v-else class="reno-done">
              <view class="reno-info">
                <view class="reno-info-row">
                  <text class="k">完成时间</text>
                  <text class="v">{{ renovation.finishAt }}</text>
                </view>
                <view class="reno-info-row">
                  <text class="k">材料</text>
                  <view class="chips">
                    <text
                      class="chip"
                      v-for="(m, idx) in renovation.materials"
                      :key="'dm-' + idx"
                      >{{ m }}</text
                    >
                  </view>
                </view>
                <view class="reno-info-row col">
                  <text class="k">装修说明</text>
                  <text class="v desc">{{ renovation.note }}</text>
                </view>
              </view>
            </view>
          </view>
        </view>

        <!-- 交易信息 -->
        <view class="section">
          <view class="section-title">交易信息</view>
          <view class="info-list">
            <view class="info-item">
              <text class="material-symbols-outlined info-icon">payments</text>
              <view class="info-content">
                <text class="info-title">交易税费</text>
                <text class="info-desc">净得价，税费由买方承担</text>
              </view>
            </view>
            <view class="info-divider"></view>
            <view class="info-item">
              <text class="material-symbols-outlined info-icon">home_work</text>
              <view class="info-content">
                <text class="info-title">房屋现状</text>
                <text class="info-desc">目前空置，业主诚意出售</text>
              </view>
            </view>
            <view class="info-divider"></view>
            <view class="info-item">
              <text class="material-symbols-outlined info-icon primary"
                >lock_open</text
              >
              <view class="info-content">
                <view class="info-title-row">
                  <text class="info-title">看房类型</text>
                  <text class="info-tag">智能开锁</text>
                </view>
                <text class="info-desc">支持APP一键开锁，随时看房</text>
              </view>
            </view>
          </view>
        </view>

        <!-- 位置周边 -->
        <view class="section">
          <view class="section-header">
            <text class="section-title">位置周边</text>
            <view class="header-more" @click="openMap">
              <text>查看地图</text>
              <text class="material-symbols-outlined more-icon"
                >chevron_right</text
              >
            </view>
          </view>
          <view class="map-box" @click="openMap">
            <image
              src="/static/images/img_fdb9c430df.png"
              mode="aspectFill"
              class="map-image"
            ></image>
            <view class="map-mask">
              <view class="map-pin">
                <text class="material-symbols-outlined pin-icon"
                  >location_on</text
                >
                <text>{{ mapPinText() }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- 相似房源 -->
        <view class="section no-pb">
          <view class="section-title">推荐相似房源</view>
          <scroll-view scroll-x="true" class="recommend-scroll">
            <view class="recommend-row">
              <view
                class="recommend-card"
                v-for="(rec, rIdx) in recommends"
                :key="rIdx"
                @click="goToRec(rec)"
              >
                <view class="rec-img-box">
                  <image
                    :src="rec.image"
                    mode="aspectFill"
                    class="rec-image"
                  ></image>
                  <view class="rec-tag">{{ rec.size }}㎡</view>
                </view>
                <view class="rec-info">
                  <view class="rec-title">{{ rec.name }}</view>
                  <view class="rec-meta">
                    <text class="rec-desc">{{ rec.rooms }}</text>
                    <text class="rec-price">{{ rec.price }}</text>
                  </view>
                </view>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>
      <view class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar">
      <view class="bar-content">
        <view class="action-btns">
          <view class="action-item" @click="callOwner">
            <text class="material-symbols-outlined bar-icon">call</text>
            <text class="bar-text">拨打电话</text>
          </view>
          <view class="action-item" @click="toggleFollow">
            <text
              v-if="isFollowed"
              class="material-symbols-outlined bar-icon filled"
              >favorite</text
            >
            <text v-else class="material-symbols-outlined bar-icon"
              >favorite_border</text
            >
            <text class="bar-text">关注</text>
          </view>
        </view>
        <view class="bar-divider"></view>
        <button class="unlock-btn" @click="openLock">
          <text
            v-if="hasSmartLock"
            class="material-symbols-outlined unlock-icon filled"
            >lock</text
          >
          <text v-else class="material-symbols-outlined unlock-icon"
            >lock_open</text
          >
          <text>{{ hasSmartLock ? "智能锁" : "密码开锁" }}</text>
        </button>
      </view>
    </view>
  </view>
</template>

<script>
import propertyApi from "@/api/property";

export default {
  data() {
    return {
      propertyId: 0,
      loading: false,
      property: {},
      statusBarHeight: 0,
      headerTop: 0,
      // 0：顶部（暗色渐变遮罩） -> 1：下滑后（磨砂白底）
      headerOpacity: 0,
      currentSwiper: 0,
      isFollowed: false,
      hasSmartLock: false,
      ownerPhone: "",
      images: [
        "/static/images/img_cdc09ae543.png",
        "/static/images/img_cdc09ae543.png",
      ],

      attributes: [
        { label: "单价", value: "9,086元/m²" },
        { label: "楼层", value: "中层 / 共32层" },
        { label: "朝向", value: "南北通透" },
        { label: "装修", value: "精装修" },
        { label: "年代", value: "2019年建" },
        { label: "产权", value: "商品房" },
      ],

      // 未装修：展示“未装修”状态说明
      // 装修进行中：展示 轮播施工图 + 进度条 + 预计完工 + 材料标签 + 施工说明
      // 装修完成：展示完成时间 + 材料 + 装修说明
      // 实现方式：
      // 页面右上角用小的状态切换条（未装修/进行中/已完成）做演示切换（静态数据）
      // 数据在 data() 里新增 renovationTabs 和 renovation（包含 images/materials/note/progress/stage 等）
      // 已补齐对应样式（与当前“高端”风格一致）

      renovationTabs: [
        { key: "none", label: "未装修" },
        { key: "in_progress", label: "进行中" },
        { key: "done", label: "已完成" },
      ],
      renovation: {
        status: "in_progress", // none | in_progress | done
        subtitle: "高端精装，环保材料（示例）",
        progress: 68,
        stage: "水电验收完成，正在进行泥木施工",
        eta: "预计 2026-02-20",
        finishAt: "2026-01-10",
        materials: ["圣象地板", "马可波罗瓷砖", "多乐士乳胶漆"],
        note: "客厅墙面已完成找平与底漆，卫生间防水已做闭水试验；全屋线管/强弱电分离施工完成。",
        images: [
          "/static/images/img_08b712f810.png",
          "/static/images/img_8642388cea.png",
        ],
      },
      recommends: [
        {
          name: "阳光花园三期 2室",
          rooms: "2室1厅",
          price: "72万",
          size: "89",
          image: "/static/images/img_71f4809787.png",
        },
        {
          name: "金地名津 精装大三房",
          rooms: "3室2厅",
          price: "98万",
          size: "96",
          image: "/static/images/img_1112597ba0.png",
        },
        {
          name: "万科四季花城",
          rooms: "3室2厅",
          price: "105万",
          size: "110",
          image: "/static/images/img_662c3c598a.png",
        },
      ],
    };
  },
  onLoad(options) {
    const info = uni.getSystemInfoSync();
    this.statusBarHeight = info.statusBarHeight;
    // #ifdef MP-WEIXIN
    try {
      if (typeof wx !== "undefined" && wx.getMenuButtonBoundingClientRect) {
        const rect = wx.getMenuButtonBoundingClientRect();
        // 与胶囊按钮垂直居中对齐：用胶囊中心线对齐返回圆按钮中心线
        const w = Number(info.windowWidth || 375);
        const rpx2px = w / 750;
        const circleBtnSizePx = 80 * rpx2px; // .circle-btn: 80rpx
        const capsuleCenterY =
          Number(rect.top || 0) + Number(rect.height || 0) / 2;
        const top = capsuleCenterY - circleBtnSizePx / 2;
        this.headerTop = Math.max(
          Number(this.statusBarHeight || 0),
          Math.round(top)
        );
      } else {
        this.headerTop = this.statusBarHeight + 12;
      }
    } catch (e) {
      this.headerTop = this.statusBarHeight + 12;
    }
    // #endif
    // #ifndef MP-WEIXIN
    this.headerTop = this.statusBarHeight + 12;
    // #endif

    const id =
      Number(options && (options.id || options.ID || options.property_id)) || 0;
    this.propertyId = id;
    if (!this.propertyId) {
      uni.showToast({ title: "房源ID缺失", icon: "none" });
      return;
    }
    this.loadDetail();
  },
  methods: {
    tagClass(tag) {
      const t = String(tag || "").trim();
      if (!t) return "";
      if (t.includes("急售") || t.includes("降价")) return "orange";
      return "";
    },
    getLayoutText(p) {
      const rooms = Number(p && p.rooms) || 0;
      const halls = Number(p && p.halls) || 0;
      const bathrooms = Number(p && p.bathrooms) || 0;
      if (!rooms && !halls && !bathrooms) return "-";
      let s = "";
      if (rooms) s += `${rooms}室`;
      if (halls) s += `${halls}厅`;
      if (bathrooms) s += `${bathrooms}卫`;
      return s || "-";
    },
    getUnitPriceText(p) {
      const price = Number(p && p.price);
      const area = Number(p && p.area);
      if (!price || !area) return "-";
      const unit = (p && p.price_unit) || "";
      let totalYuan = price;
      if (unit === "万") totalYuan = price * 10000;
      const per = totalYuan / area;
      if (!per || !isFinite(per)) return "-";
      return `${Math.round(per).toLocaleString()}元/㎡`;
    },
    buildAttributes(p) {
      const attrs = [];
      attrs.push({ label: "单价", value: this.getUnitPriceText(p) });
      const floor = String((p && p.floor_level) || "").trim();
      const totalFloors = Number(p && p.total_floors) || 0;
      if (floor || totalFloors) {
        attrs.push({
          label: "楼层",
          value: `${floor || "-"}${totalFloors ? ` / 共${totalFloors}层` : ""}`,
        });
      }
      attrs.push({
        label: "朝向",
        value: p && p.orientation ? p.orientation : "-",
      });
      attrs.push({
        label: "装修",
        value: p && p.decoration_type ? p.decoration_type : "-",
      });
      const year = Number(p && p.build_year) || 0;
      attrs.push({ label: "年代", value: year ? `${year}年建` : "-" });
      attrs.push({
        label: "物业类型",
        value: p && p.property_type ? p.property_type : "-",
      });
      return attrs;
    },
    getCommissionText(p) {
      const rate = String((p && p.commission_rate) || "").trim();
      const reward = String((p && p.commission_reward) || "").trim();
      const hasRate =
        !!rate && rate !== "0" && rate !== "0.0" && rate !== "0.00";
      const hasReward =
        !!reward && reward !== "0" && reward !== "0.0" && reward !== "0.00";
      if (hasRate && hasReward) return `佣金${rate}% ，成交奖励￥${reward}`;
      if (hasRate) return `佣金${rate}%`;
      if (hasReward) return `成交奖励￥${reward}`;
      return "佣金待定";
    },
    mapPinText() {
      const name = String(
        (this.property && this.property.community_name) || ""
      ).trim();
      const addr = String(
        (this.property && this.property.address) || ""
      ).trim();
      if (name && addr) return `${name} · ${addr}`;
      return name || addr || "暂无位置信息";
    },
    async loadDetail() {
      if (this.loading || !this.propertyId) return false;
      this.loading = true;
      let res;
      try {
        res = await propertyApi.getDetail({ id: this.propertyId });
      } catch (e) {
        this.loading = false;
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return false;
      }
      this.loading = false;
      if (!res || res.code !== 0) {
        setTimeout(() => {
          uni.navigateBack();
        }, 600);
        return false;
      }

      const data = res.data || {};
      const p = data.property || {};
      if (!p || !p.id) {
        setTimeout(() => {
          uni.navigateBack();
        }, 600);
        return false;
      }
      this.property = p;
      this.isFollowed = !!data.is_followed;
      this.hasSmartLock = Number(p.has_smart_lock) === 1;
      this.ownerPhone = String(p.owner_phone || "").trim();

      const imgs = Array.isArray(data.images)
        ? data.images
        : Array.isArray(p.images)
        ? p.images
        : [];
      if (imgs.length > 0) this.images = imgs;
      this.currentSwiper = 0;

      this.attributes = this.buildAttributes(p);

      const r = data.renovation || {};
      const statusRaw = String(r.renovation_status || "").trim();
      let status = "none";
      if (
        statusRaw === "in_progress" ||
        statusRaw === "done" ||
        statusRaw === "none"
      ) {
        status = statusRaw;
      } else if (statusRaw === "1") {
        status = "in_progress";
      } else if (statusRaw === "2") {
        status = "done";
      }
      this.renovation = {
        status,
        subtitle: String(r.subtitle || "房源装修情况").trim() || "房源装修情况",
        progress: Number(r.progress_percentage || 0) || 0,
        stage: String(r.current_stage || "").trim() || "—",
        eta: String(r.estimated_finish_date || "").trim() || "—",
        finishAt: String(r.actual_finish_date || "").trim() || "—",
        materials: Array.isArray(r.materials) ? r.materials : [],
        note: String(r.notes || "").trim() || "—",
        images: Array.isArray(r.images) ? r.images : [],
      };

      const rec = Array.isArray(data.recommends) ? data.recommends : [];
      this.recommends = rec.map((it) => ({
        id: it.id,
        name: it.title || it.name || "房源",
        rooms: this.getLayoutText(it),
        price: `${it.price || "-"}${it.price_unit || ""}`,
        size: it.area || "-",
        image: it.image || "",
      }));
      return true;
    },
    onPageScroll(e) {
      // scroll-view 的 scrollTop（px）
      const st = Number((e && e.detail && e.detail.scrollTop) || 0);
      // 越往下越“实”（白底更明显）；回到顶部恢复暗色渐变
      const fadeDistance = 160;
      const next = Math.min(1, st / fadeDistance);
      this.headerOpacity = Math.max(0, Math.min(1, Number(next.toFixed(3))));
    },
    goBack() {
      uni.navigateBack();
    },
    goToRec(rec) {
      const id = Number(rec && (rec.id || rec.ID)) || 0;
      if (!id) return;
      if (id === this.propertyId) return;
      uni.redirectTo({
        url: `/pages/property_detail/property_detail?id=${encodeURIComponent(
          id
        )}`,
      });
    },
    swiperChange(e) {
      this.currentSwiper = e.detail.current;
    },
    handleShare() {
      const title =
        this.property && this.property.title ? this.property.title : "房源";
      const id = this.propertyId;
      uni.showActionSheet({
        itemList: ["生成获客海报", "复制房源链接", "模拟分享（占位）"],
        success: (res) => {
          if (res.tapIndex === 0) {
            uni.navigateTo({
              url: `/pages/marketing_posters/marketing_posters?title=${encodeURIComponent(
                title
              )}`,
            });
            return;
          }
          if (res.tapIndex === 1) {
            const link = id
              ? `/pages/property_detail/property_detail?id=${encodeURIComponent(
                  id
                )}`
              : `/pages/property_detail/property_detail`;
            uni.setClipboardData({
              data: `${title}\n${link}`,
              success: () =>
                uni.showToast({ title: "已复制链接", icon: "none" }),
            });
            return;
          }
          if (res.tapIndex === 2) {
            uni.showToast({ title: "分享功能待接入", icon: "none" });
          }
        },
      });
    },
    openMap() {
      const lat = Number(this.property && this.property.latitude);
      const lng = Number(this.property && this.property.longitude);
      const name = String(
        (this.property &&
          (this.property.community_name || this.property.title)) ||
          "房源"
      ).trim();
      const address = String(
        (this.property && this.property.address) || ""
      ).trim();
      if (!lat || !lng || !isFinite(lat) || !isFinite(lng)) {
        uni.showToast({ title: "暂无定位信息", icon: "none" });
        return;
      }
      uni.openLocation({
        latitude: lat,
        longitude: lng,
        name,
        address: address || name,
        fail: () => {
          uni.showToast({ title: "无法打开地图", icon: "none" });
        },
      });
    },
    async toggleFollow() {
      if (!this.propertyId) return;
      let res;
      try {
        res = await propertyApi.toggleFollow({ id: this.propertyId });
      } catch (e) {
        if (!uni.getStorageSync("token")) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return;
      }
      if (!res || res.code !== 0) return;
      const data = res.data || {};
      this.isFollowed = !!data.is_followed;
      if (this.property) {
        this.property.follow_count = data.follow_count;
      }
      uni.showToast({
        title: this.isFollowed ? "已关注" : "已取消关注",
        icon: "none",
      });
    },
    callOwner() {
      if (!this.ownerPhone) {
        uni.showToast({ title: "暂无业主电话", icon: "none" });
        return;
      }
      uni.makePhoneCall({
        phoneNumber: this.ownerPhone,
      });
    },
    openLock() {
      const ss = String(
        (this.property && this.property.sale_status) || ""
      ).trim();
      if (ss === "sold" || ss === "off_market") {
        const label =
          this.property && this.property.sale_status_label
            ? this.property.sale_status_label
            : "不可操作";
        uni.showToast({ title: `房源${label}，暂不可开锁`, icon: "none" });
        return;
      }
      // 有智能锁 -> 进入申请开锁流程；无智能锁 -> 进入密码开锁详情页
      uni.navigateTo({
        url: this.hasSmartLock
          ? "/pages/unlock_steps/unlock_steps"
          : "/pages/unlock_details/unlock_details",
      });
    },
    setRenovationStatus(key) {
      // 展示用：不允许手动切换，避免与真实数据不一致
      void key;
    },
  },
};
</script>

<style lang="scss">
.detail-container {
  height: 100vh;
  background-color: #f6f7f8;
  display: flex;
  flex-direction: column;
  position: relative;
}

.float-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  justify-content: flex-start;
  padding: 0 32rpx;
  pointer-events: none;
  padding-bottom: 8rpx;
  box-sizing: border-box;

  /* 背景遮罩：不挡点击 */
  .fh-mask {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    pointer-events: none;
    transition: opacity 0.18s ease;
  }

  /* 顶部：暗色渐变，保证返回键在封面图上更清晰 */
  .fh-mask--dark {
    background: linear-gradient(
      to bottom,
      rgba(15, 23, 42, 0.45) 0%,
      rgba(15, 23, 42, 0.18) 55%,
      rgba(15, 23, 42, 0) 100%
    );
  }

  /* 下滑后：磨砂白底（更接近普通页面导航栏） */
  .fh-mask--light {
    background: rgba(255, 255, 255, 0.92);
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
    border-bottom: 1rpx solid rgba(226, 232, 240, 0.9);
    box-shadow: 0 6rpx 18rpx rgba(15, 23, 42, 0.04);
  }

  .header-left {
    display: flex;
    gap: 24rpx;
    pointer-events: auto;
    position: relative;
    z-index: 2;
  }

  .circle-btn {
    width: 80rpx;
    height: 80rpx;
    background-color: rgba(0, 0, 0, 0.2);
    backdrop-filter: blur(10px);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ffffff;

    .material-symbols-outlined {
      font-size: 40rpx;
    }

    &:active {
      background-color: rgba(0, 0, 0, 0.4);
    }

    /* 下滑后：切到浅色按钮，避免白底上显得“脏/灰” */
    &.is-solid {
      background-color: rgba(15, 23, 42, 0.06);
      color: #0f172a;
    }
    &.is-solid:active {
      background-color: rgba(15, 23, 42, 0.12);
    }
  }
}

.detail-scroll {
  flex: 1;
  overflow: hidden;
}

.banner {
  width: 100%;
  height: 640rpx;
  position: relative;

  .swiper {
    width: 100%;
    height: 100%;
  }

  .banner-image {
    width: 100%;
    height: 100%;
  }

  .banner-indicator {
    position: absolute;
    bottom: 64rpx;
    right: 32rpx;
    background-color: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(4rpx);
    color: #ffffff;
    font-size: 24rpx;
    font-weight: 600;
    padding: 4rpx 24rpx;
    border-radius: 40rpx;
    border: 1px solid rgba(255, 255, 255, 0.1);
  }
}

.content-card {
  background-color: #ffffff;
  border-radius: 48rpx 48rpx 0 0;
  margin-top: -48rpx;
  position: relative;
  z-index: 10;
  padding: 48rpx 40rpx;
  display: flex;
  flex-direction: column;
  gap: 48rpx;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.title-section {
  display: flex;
  flex-direction: column;
  gap: 24rpx;

  .title-row {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16rpx;

    .title {
      font-size: 40rpx;
      font-weight: bold;
      color: #0f172a;
      line-height: 1.3;
      flex: 1;
    }

    .share-btn {
      display: flex;
      align-items: center;
      gap: 8rpx;
      padding: 12rpx 16rpx;
      background: rgba(37, 99, 235, 0.08);
      border: 1rpx solid rgba(37, 99, 235, 0.14);
      border-radius: 18rpx;
      font-size: 24rpx;
      color: #2563eb;
      font-weight: 700;
      flex-shrink: 0;

      .share-icon {
        font-size: 28rpx;
      }

      &:active {
        transform: scale(0.98);
        background: rgba(37, 99, 235, 0.12);
      }
    }
  }

  .tags-row {
    display: flex;
    flex-wrap: wrap;
    gap: 16rpx;

    .tag {
      padding: 8rpx 16rpx;
      background-color: #f1f5f9;
      color: #475569;
      font-size: 22rpx;
      font-weight: 500;
      border-radius: 8rpx;

      &.orange {
        background-color: #fff7ed;
        color: #f97316;
      }
    }
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  padding: 26rpx 0;
  border-top: 1rpx solid rgba(226, 232, 240, 0.9);
  border-bottom: 1rpx solid rgba(226, 232, 240, 0.9);

  .stats-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6rpx;
    padding: 0 18rpx;

    &.border-l {
      border-left: 1rpx solid rgba(226, 232, 240, 0.9);
    }

    .stats-val {
      font-size: 40rpx;
      font-weight: 800;
      color: #0f172a;
      display: inline-flex;
      align-items: baseline;
      justify-content: center;
      line-height: 1;
      white-space: nowrap;
      max-width: 100%;

      &.orange {
        color: #f97316;
      }
      /* 非价格类字段更容易换行，略缩小字号保证一行展示 */
      &:not(.orange) {
        font-size: 34rpx;
      }

      .unit {
        font-size: 22rpx;
        font-weight: 700;
        margin-left: 4rpx;
        opacity: 0.9;
      }
    }

    .stats-label {
      font-size: 22rpx;
      color: #64748b;
      font-weight: 700;
      letter-spacing: 0.4rpx;
    }
  }
}

.commission-banner {
  background: linear-gradient(to right, #f97316, #ea580c);
  border-radius: 24rpx;
  padding: 32rpx;
  position: relative;
  overflow: hidden;
  box-shadow: 0 10rpx 20rpx rgba(249, 115, 22, 0.2);

  .banner-bg-decor {
    position: absolute;
    right: -48rpx;
    top: -48rpx;
    width: 192rpx;
    height: 192rpx;
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 50%;
    filter: blur(40rpx);
  }

  .banner-main {
    display: flex;
    align-items: center;
    gap: 24rpx;
    position: relative;
    z-index: 10;

    .icon-box {
      width: 80rpx;
      height: 80rpx;
      background-color: rgba(255, 255, 255, 0.2);
      backdrop-filter: blur(10rpx);
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #ffffff;

      .material-symbols-outlined {
        font-size: 48rpx;
      }
    }

    .banner-info {
      display: flex;
      flex-direction: column;
      gap: 4rpx;

      .banner-title {
        font-size: 32rpx;
        font-weight: bold;
        color: #ffffff;
      }

      .banner-tip {
        font-size: 22rpx;
        color: rgba(255, 255, 255, 0.8);
      }
    }
  }
}

.data-row {
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1rpx solid rgba(226, 232, 240, 0.9);
  border-radius: 24rpx;
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
  padding: 20rpx 16rpx;
  display: flex;
  align-items: stretch;

  .data-group {
    flex: 1;
    display: flex;
    align-items: stretch;
  }

  .data-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 6rpx;
    padding: 6rpx 0;
    position: relative;

    &:not(:last-child)::after {
      content: "";
      position: absolute;
      right: 0;
      top: 10rpx;
      bottom: 10rpx;
      width: 1rpx;
      background-color: rgba(226, 232, 240, 0.9);
    }

    .val {
      font-size: 40rpx;
      font-weight: 800;
      color: #0f172a;
      letter-spacing: 0.5rpx;
    }

    .lab {
      font-size: 22rpx;
      color: #64748b;
      font-weight: 600;
    }
  }
}

.section {
  display: flex;
  flex-direction: column;
  gap: 24rpx;

  &.no-pb {
    padding-bottom: 0;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .section-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #0f172a;
  }

  .header-more {
    display: flex;
    align-items: center;
    gap: 4rpx;
    font-size: 24rpx;
    color: #2d9cf0;
    font-weight: 500;

    .more-icon {
      font-size: 28rpx;
    }
  }
}

.attr-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
  background-color: #f8fafc;
  padding: 32rpx;
  border-radius: 24rpx;

  .attr-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 28rpx;

    .attr-label {
      color: #64748b;
    }
    .attr-val {
      color: #0f172a;
      font-weight: 500;
    }
  }
}

/* 装修状态 */
.reno-tabs {
  display: flex;
  gap: 12rpx;
  background-color: #f1f5f9;
  padding: 6rpx;
  border-radius: 18rpx;
}

.reno-tab {
  padding: 10rpx 16rpx;
  border-radius: 14rpx;
  font-size: 22rpx;
  color: #64748b;
  font-weight: 700;
  background-color: transparent;

  &.active {
    background-color: #ffffff;
    color: #0f172a;
    box-shadow: 0 6rpx 14rpx rgba(15, 23, 42, 0.06);
  }
}

.reno-card {
  background-color: #f8fafc;
  border-radius: 24rpx;
  padding: 24rpx;
  border: 1px solid #f1f5f9;
}

.reno-top {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.reno-badge {
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
  font-weight: 800;
  letter-spacing: 0.4rpx;

  &.none {
    background: rgba(148, 163, 184, 0.18);
    color: #475569;
  }
  &.in_progress {
    background: rgba(245, 158, 11, 0.18);
    color: #b45309;
  }
  &.done {
    background: rgba(34, 197, 94, 0.16);
    color: #15803d;
  }
}

.reno-sub {
  font-size: 24rpx;
  color: #64748b;
  font-weight: 600;
}

.reno-empty {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  padding: 8rpx 4rpx 4rpx;
}

.reno-empty-icon {
  font-size: 44rpx;
  color: #94a3b8;
}

.reno-empty-texts {
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.reno-empty-title {
  font-size: 28rpx;
  font-weight: 800;
  color: #0f172a;
}

.reno-empty-desc {
  font-size: 24rpx;
  color: #64748b;
  line-height: 1.5;
}

.reno-swiper {
  width: 100%;
  height: 260rpx;
  border-radius: 20rpx;
  overflow: hidden;
}

.reno-image {
  width: 100%;
  height: 100%;
}

.reno-progress {
  margin-top: 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
}

.reno-progress-bar {
  flex: 1;
  height: 16rpx;
  border-radius: 999rpx;
  background-color: rgba(226, 232, 240, 0.9);
  overflow: hidden;
}

.reno-progress-fill {
  height: 100%;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #60a5fa, #2563eb);
}

.reno-progress-text {
  font-size: 22rpx;
  color: #334155;
  font-weight: 700;
  white-space: nowrap;
}

.reno-info {
  margin-top: 16rpx;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.reno-info-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16rpx;

  &.col {
    flex-direction: column;
  }

  .k {
    font-size: 24rpx;
    color: #64748b;
    font-weight: 700;
    flex-shrink: 0;
  }

  .v {
    font-size: 24rpx;
    color: #0f172a;
    font-weight: 600;
    text-align: right;
    flex: 1;

    &.desc {
      text-align: left;
      color: #334155;
      font-weight: 500;
      line-height: 1.55;
    }
  }
}

.chips {
  flex: 1;
  display: flex;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 10rpx;
}

.chip {
  padding: 6rpx 12rpx;
  border-radius: 999rpx;
  background-color: rgba(37, 99, 235, 0.08);
  border: 1rpx solid rgba(37, 99, 235, 0.12);
  color: #2563eb;
  font-size: 22rpx;
  font-weight: 700;
}

.info-list {
  display: flex;
  flex-direction: column;
  padding: 32rpx;
  border-radius: 24rpx;
  border: 1px solid #f1f5f9;

  .info-item {
    display: flex;
    gap: 24rpx;
    align-items: flex-start;

    .info-icon {
      font-size: 36rpx;
      color: #94a3b8;
      margin-top: 4rpx;

      &.primary {
        color: #2d9cf0;
      }
    }

    .info-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      gap: 8rpx;
    }

    .info-title-row {
      display: flex;
      align-items: center;
      gap: 16rpx;
    }

    .info-title {
      font-size: 28rpx;
      font-weight: 500;
      color: #0f172a;
    }

    .info-tag {
      padding: 4rpx 12rpx;
      background-color: #eff6ff;
      color: #2d9cf0;
      font-size: 20rpx;
      font-weight: bold;
      border-radius: 4rpx;
    }

    .info-desc {
      font-size: 24rpx;
      color: #64748b;
    }
  }

  .info-divider {
    height: 1px;
    background-color: #f1f5f9;
    margin: 32rpx 0;
  }
}

.map-box {
  width: 100%;
  height: 288rpx;
  border-radius: 24rpx;
  overflow: hidden;
  position: relative;

  .map-image {
    width: 100%;
    height: 100%;
    opacity: 0.8;
  }

  .map-mask {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: rgba(0, 0, 0, 0.05);

    .map-pin {
      display: flex;
      align-items: center;
      gap: 12rpx;
      background-color: rgba(255, 255, 255, 0.95);
      backdrop-filter: blur(10px);
      padding: 12rpx 24rpx;
      border-radius: 16rpx;
      box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
      font-size: 24rpx;
      font-weight: bold;
      color: #0f172a;

      .pin-icon {
        color: #2d9cf0;
        font-size: 28rpx;
      }
    }
  }
}

.recommend-scroll {
  width: 100%;
  white-space: nowrap;
  // 左侧保持与内容区对齐（40rpx），右侧拉满到屏幕边缘
  margin-left: -40rpx;
  margin-right: -40rpx;
  padding-left: 40rpx;
  padding-right: 0;
  box-sizing: border-box;

  .recommend-row {
    display: inline-flex;
    gap: 24rpx;
    padding-bottom: 8rpx;
  }

  .recommend-card {
    width: 320rpx;
    display: flex;
    flex-direction: column;
    gap: 16rpx;

    .rec-img-box {
      width: 100%;
      height: 192rpx;
      border-radius: 16rpx;
      overflow: hidden;
      position: relative;
      background-color: #f1f5f9;

      .rec-image {
        width: 100%;
        height: 100%;
      }
      .rec-tag {
        position: absolute;
        bottom: 8rpx;
        left: 8rpx;
        background-color: rgba(0, 0, 0, 0.5);
        backdrop-filter: blur(4rpx);
        color: #ffffff;
        font-size: 20rpx;
        padding: 4rpx 12rpx;
        border-radius: 8rpx;
      }
    }

    .rec-info {
      display: flex;
      flex-direction: column;
      gap: 4rpx;

      .rec-title {
        font-size: 28rpx;
        font-weight: 500;
        color: #0f172a;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .rec-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;

        .rec-desc {
          font-size: 24rpx;
          color: #94a3b8;
        }
        .rec-price {
          font-size: 24rpx;
          font-weight: bold;
          color: #f97316;
        }
      }
    }
  }
}

.bottom-bar {
  background-color: #ffffff;
  border-top: 1rpx solid #f1f5f9;
  padding: 16rpx 32rpx calc(env(safe-area-inset-bottom) + 16rpx);
  box-shadow: 0 -10rpx 30rpx rgba(0, 0, 0, 0.05);

  .bar-content {
    display: flex;
    align-items: center;
    gap: 20rpx;
    max-width: 960rpx;
    margin: 0 auto;
  }

  .action-btns {
    display: flex;
    gap: 40rpx;
    padding: 0 8rpx;
  }

  .action-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4rpx;
    color: #64748b;
    min-width: 80rpx;

    .bar-icon {
      font-size: 44rpx;
    }
    .bar-text {
      font-size: 20rpx;
      font-weight: 500;
      white-space: nowrap;
    }

    .filled {
      font-variation-settings: "FILL" 1;
    }

    &:active {
      color: #0f172a;
    }
  }

  .bar-divider {
    width: 1px;
    height: 64rpx;
    background-color: #f1f5f9;
    margin: 0 8rpx;
  }

  .unlock-btn {
    flex: 1;
    height: 88rpx;
    background-color: #2d9cf0;
    color: #ffffff;
    font-size: 30rpx;
    font-weight: bold;
    border-radius: 16rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12rpx;
    box-shadow: 0 8rpx 20rpx rgba(45, 156, 240, 0.3);
    border: none;

    &::after {
      border: none;
    }

    .unlock-icon {
      font-size: 36rpx;
      font-variation-settings: "FILL" 1;
    }

    &:active {
      background-color: #1d82cc;
      transform: scale(0.98);
    }
  }
}

.bottom-spacer {
  height: 64rpx;
}
</style>
