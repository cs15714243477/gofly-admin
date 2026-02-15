<template>
  <view class="detail-container">
    <!-- 顶部悬浮操作 -->
    <view
      class="float-header"
      :style="{
        paddingTop: headerTop + 'px',
        paddingLeft: headerPadLeftPx ? headerPadLeftPx + 'px' : undefined,
        paddingRight: headerPadRightPx ? headerPadRightPx + 'px' : undefined,
      }"
    >
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
          <swiper-item v-for="(item, index) in bannerItems" :key="index">
            <video
              v-if="item.type === 'video'"
              :src="item.src"
              :poster="videoPoster"
              class="banner-video"
              controls
              show-center-play-btn
              object-fit="cover"
              playsinline
              enable-progress-gesture
            ></video>
            <image
              v-else
              :src="item.src"
              mode="aspectFill"
              class="banner-image"
            ></image>
          </swiper-item>
        </swiper>
        <view
          class="banner-actions"
          v-if="showDownloadBtn"
          @click.stop="downloadCurrent"
        >
          <view class="dl-pill">
            <text class="material-symbols-outlined dl-icon">download</text>
            <text class="dl-text">{{ downloadLabel }}</text>
          </view>
        </view>
        <view class="banner-indicator"
          >{{ bannerItems.length ? currentSwiper + 1 : 0 }}/{{
            bannerItems.length
          }}</view
        >
      </view>

      <!-- 内容卡片 -->
      <view class="content-card">
        <view class="title-section">
          <view class="title-row">
            <text class="title">{{ (property && property.title) || "-" }}</text>
            <view class="action-col">
              <view
                v-if="!isPublicView"
                class="share-btn"
                @click="handleShare"
                @longpress="debugEditState"
              >
                <text class="material-symbols-outlined share-icon">share</text>
                <text>分享</text>
              </view>
              <view class="edit-btn" v-if="canEditThisProperty" @click="goEdit">
                <text class="material-symbols-outlined edit-icon">edit</text>
                <text>编辑</text>
              </view>
            </view>
          </view>
          <view class="tags-row">
            <text
              v-for="(t, idx) in ((property && property.tags) || []).slice(
                0,
                6,
              )"
              :key="idx"
              class="tag"
              :class="tagClass(t)"
            >
              {{ t }}
            </text>
          </view>
          <view v-if="!isPublicView" class="commission-line">
            <text class="material-symbols-outlined commission-ic"
              >currency_yen</text
            >
            <text class="commission-text">{{
              getCommissionText(property)
            }}</text>
          </view>
          <view class="status-row" v-if="canEditThisProperty">
            <view
              class="status-chip"
              :class="(property && property.sale_status) || ''"
            >
              <text class="material-symbols-outlined status-ic">sell</text>
              <text>{{
                (property && property.sale_status_label) || "状态"
              }}</text>
            </view>
            <view
              v-if="Number(property && property.hot_status) === 1"
              class="status-chip hot"
            >
              <text class="material-symbols-outlined status-ic"
                >local_fire_department</text
              >
              <text>推荐</text>
            </view>
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

        <!-- 房源描述 -->
        <view class="section">
          <view class="section-title">房源描述</view>
          <view v-if="property && property.custom_desc" class="desc-card">
            <view class="desc-text">{{ property.custom_desc }}</view>
          </view>
          <view v-else class="desc-empty">暂无描述</view>
        </view>

        <!-- 房主与收房（与后台表单字段一致） -->
        <view class="section" v-if="canEditThisProperty">
          <view class="section-title">房主与收房</view>
          <view class="contact-card">
            <view class="contact-grid">
              <view class="contact-item">
                <text class="contact-label">业主姓名</text>
                <text class="contact-value">{{
                  (property && property.owner_name) || "-"
                }}</text>
              </view>
              <view class="contact-item">
                <text class="contact-label">业主电话</text>
                <text class="contact-value">{{
                  (property && property.owner_phone) || "-"
                }}</text>
              </view>
              <view class="contact-item">
                <text class="contact-label">收房人姓名</text>
                <text class="contact-value">{{
                  (property && property.receiver_name) || "-"
                }}</text>
              </view>
              <view class="contact-item">
                <text class="contact-label">收房人电话</text>
                <text class="contact-value">{{
                  (property && property.receiver_phone) || "-"
                }}</text>
              </view>
              <view class="contact-item wide">
                <text class="contact-label">收房价格（支付业主）</text>
                <text class="contact-value price">{{
                  property && property.receiver_price
                    ? "¥" + property.receiver_price
                    : "-"
                }}</text>
              </view>
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

            <!-- 工序时间线（全流程） -->
            <view
              v-if="renovation.timeline && renovation.timeline.length"
              class="reno-stage-timeline"
            >
              <view class="reno-stage-title-row">
                <text class="reno-stage-title">施工时间线</text>
                <text class="reno-stage-tip">按工序查看状态与图片</text>
              </view>

              <view class="stage-list">
                <view
                  class="stage-item"
                  v-for="(it, idx) in renovation.timeline"
                  :key="it.stage + '_' + idx"
                >
                  <view class="stage-rail">
                    <view class="stage-dot" :class="it.status"></view>
                    <view
                      class="stage-line"
                      v-if="idx !== renovation.timeline.length - 1"
                    ></view>
                  </view>
                  <view class="stage-body">
                    <view class="stage-head">
                      <text class="stage-name">{{ it.stage }}</text>
                      <text class="stage-tag" :class="it.status">{{
                        stageStatusText(it.status)
                      }}</text>
                      <text v-if="it.date" class="stage-date">{{ it.date }}</text>
                    </view>
                    <text v-if="it.note" class="stage-note">{{ it.note }}</text>
                    <view
                      class="stage-img-grid"
                      v-if="it.images && it.images.length"
                    >
                      <image
                        class="stage-thumb"
                        v-for="(img, iidx) in it.images"
                        :key="iidx"
                        :src="img"
                        mode="aspectFill"
                        @click.stop="previewStageImage(it, iidx)"
                      ></image>
                    </view>
                  </view>
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
            <map
              v-if="hasLocation"
              class="map-native"
              :latitude="mapLat"
              :longitude="mapLng"
              :scale="16"
              :markers="mapMarkers"
              :enable-zoom="false"
              :enable-scroll="false"
              :enable-rotate="false"
              :enable-overlooking="false"
              :enable-satellite="false"
              :show-location="false"
              @tap="openMap"
            ></map>
            <view v-else class="map-image map-image-empty">
              <text class="material-symbols-outlined">map</text>
            </view>
            <view class="map-mask">
              <view class="map-pin">
                <text class="material-symbols-outlined pin-icon"
                  >location_on</text
                >
                <text>{{ mapPinText() }}</text>
              </view>
              <view class="map-coords" v-if="hasLocation">
                <text class="coord-text"
                  >{{ mapLng.toFixed(6) }}, {{ mapLat.toFixed(6) }}</text
                >
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
                    v-if="normalizeImage(rec.image)"
                    :src="normalizeImage(rec.image)"
                    mode="aspectFill"
                    class="rec-image"
                  ></image>
                  <view v-else class="rec-image rec-image-empty">
                    <text class="material-symbols-outlined">image</text>
                  </view>
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
      <view v-if="!isPublicView" class="bottom-spacer"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar" v-if="!isPublicView">
      <view class="bar-content">
        <view class="action-btns">
          <view class="action-item" @click="callOwner">
            <text class="material-symbols-outlined bar-icon">call</text>
            <text class="bar-text">拨打电话</text>
          </view>
          <view class="action-item" @click="recordShowing">
            <text class="material-symbols-outlined bar-icon">location_on</text>
            <text class="bar-text">带看</text>
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
          <text>{{ hasSmartLock ? "智能锁" : "未绑定智能锁" }}</text>
        </button>
      </view>
    </view>

    <!-- 分享弹窗（转发好友/复制文案） -->
    <view v-if="shareSheetOpen" class="share-mask" @tap="closeShareSheet">
      <view class="share-sheet" @tap.stop>
        <view class="share-title">分享房源</view>
        <view class="share-actions">
          <button
            class="share-action primary"
            open-type="share"
            :data-id="propertyId"
            :data-title="String((property && property.title) || '房源')"
            :data-image="shareCover"
            @tap.stop="onTapForward"
          >
            <text class="material-symbols-outlined">ios_share</text>
            <text>转发给好友</text>
          </button>
          <view class="share-divider"></view>
          <view class="share-action" @tap="copyPromoteText">
            <text class="material-symbols-outlined">content_copy</text>
            <text>复制推广文案</text>
          </view>
        </view>
        <view class="share-cancel" @tap="closeShareSheet">取消</view>
      </view>
    </view>

    <!-- 带看登记弹窗（填写客户姓名/电话） -->
    <view v-if="showingSheetOpen" class="sheet-mask" @tap="closeShowingSheet">
      <view class="sheet-panel" @tap.stop>
        <view class="sheet-title">记录带看</view>
        <view class="sheet-tip">请填写客户姓名与电话，便于后续在带看记录中查看</view>
        <view class="sheet-form">
          <view class="sheet-field">
            <text class="sheet-label">客户姓名</text>
            <input
              class="sheet-input"
              v-model="showingForm.client_name"
              placeholder="请输入客户姓名"
              maxlength="20"
              confirm-type="next"
            />
          </view>
          <view class="sheet-divider"></view>
          <view class="sheet-field">
            <text class="sheet-label">客户电话</text>
            <input
              class="sheet-input"
              v-model="showingForm.client_phone"
              placeholder="请输入客户电话"
              type="number"
              maxlength="20"
              confirm-type="done"
              @confirm="submitShowing"
            />
          </view>
        </view>
        <view class="sheet-actions-row">
          <view class="sheet-btn cancel" @tap="closeShowingSheet">取消</view>
          <view
            class="sheet-btn primary"
            :class="{ disabled: showingSaving }"
            @tap="submitShowing"
          >
            {{ showingSaving ? "保存中..." : "确认记录" }}
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import propertyApi from "@/api/property";
import userApi from "@/api/user";
import $store from "@/store";

export default {
  onShareAppMessage(options) {
    const dataset =
      options && options.target && options.target.dataset
        ? options.target.dataset
        : {};
    const id = Number(dataset.id || this.propertyId || 0) || 0;
    const title = String(dataset.title || (this.property && this.property.title) || "房源").trim() || "房源";
    const rawImg = String(dataset.image || "").trim();
    const imageUrl = /^https?:\/\//i.test(rawImg) ? rawImg : "";

    // “登录走经纪人”：携带 from_agent_id；“未登录走客户”：不带
    const token = uni.getStorageSync("token");
    const userStore = $store("user");
    const agentId = token
      ? Number(this.currentUserId || (userStore && userStore.userInfo && userStore.userInfo.id) || 0) || 0
      : 0;

    let path = `/pages/property_detail/property_detail?id=${encodeURIComponent(id)}&public=1`;
    if (agentId > 0) {
      path += `&from_agent_id=${encodeURIComponent(agentId)}&from_style=0`;
    }

    // 尝试写入分享日志（不影响转发）
    try {
      if (!this.isPublicView && id > 0) {
        this.tryLogActivity("share", id, { channel: "转发给好友" });
      }
    } catch (e) {}

    const out = { title, path };
    if (imageUrl) out.imageUrl = imageUrl;
    return out;
  },
  data() {
    return {
      propertyId: 0,
      loading: false,
      property: {},
      isPublicView: false,
      publicFromAgentId: 0,
      publicFromStyle: 0,
      currentUserId: 0,
      canManageProperties: false,
      statusBarHeight: 0,
      headerTop: 0,
      headerPadLeftPx: 0,
      headerPadRightPx: 0,
      // 0：顶部（暗色渐变遮罩） -> 1：下滑后（磨砂白底）
      headerOpacity: 0,
      currentSwiper: 0,
      isFollowed: false,
      hasSmartLock: false,
      ownerPhone: "",
      images: [],
      bannerItems: [],
      videoUrl: "",
      videoPoster: "",
      allowImageDownload: true,
      allowVideoDownload: true,
      mapLat: 0,
      mapLng: 0,
      mapMarkers: [],

      attributes: [
        { label: "单价", value: "9,086元/m²" },
        { label: "楼层", value: "中层 / 共32层" },
        { label: "朝向", value: "南北通透" },
        { label: "装修", value: "精装修" },
        { label: "年代", value: "2019年建" },
        { label: "产权", value: "商品房" },
      ],

      // 装修展示：
      // - 总状态：未装修 / 施工中 / 已完工
      // - 工序时间线：全流程按工序展示状态（未开始/进行中/已完成）
      // - 工序图片：按工序分组展示，点击可预览

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
        images: [],
        timeline: [],
      },
      recommends: [
        {
          name: "阳光花园三期 2室",
          rooms: "2室1厅",
          price: "72万",
          size: "89",
          image: "",
        },
        {
          name: "金地名津 精装大三房",
          rooms: "3室2厅",
          price: "98万",
          size: "96",
          image: "",
        },
        {
          name: "万科四季花城",
          rooms: "3室2厅",
          price: "105万",
          size: "110",
          image: "",
        },
      ],
      viewLogged: false,
      shareSheetOpen: false,
      showingSheetOpen: false,
      showingSaving: false,
      // 预售 Tab 详情模式：仅允许从预售列表进入（需 view_key）
      isPreSaleView: false,
      preSaleViewKey: "",
      showingForm: {
        client_name: "",
        client_phone: "",
      },
    };
  },
  onLoad(options) {
    const info = uni.getSystemInfoSync();
    this.statusBarHeight = info.statusBarHeight;
    // #ifdef MP-WEIXIN
    try {
      if (typeof wx !== "undefined" && wx.getMenuButtonBoundingClientRect) {
        const rect = wx.getMenuButtonBoundingClientRect();
        // 让右上角按钮避开胶囊按钮区域
        const w = Number(info.windowWidth || 375);
        const rpx2px = w / 750;
        const padLeftPx = Math.round(32 * rpx2px); // .float-header 左右 32rpx
        const rightSafePx = Math.max(0, w - Number(rect.left || 0));
        this.headerPadLeftPx = padLeftPx;
        this.headerPadRightPx = padLeftPx + rightSafePx;
        // 与胶囊按钮垂直居中对齐：用胶囊中心线对齐返回圆按钮中心线
        const circleBtnSizePx = 80 * rpx2px; // .circle-btn: 80rpx
        const capsuleCenterY =
          Number(rect.top || 0) + Number(rect.height || 0) / 2;
        const top = capsuleCenterY - circleBtnSizePx / 2;
        this.headerTop = Math.max(
          Number(this.statusBarHeight || 0),
          Math.round(top),
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

    const publicRaw = String(
      (options && (options.public || options.is_public)) || "",
    )
      .trim()
      .toLowerCase();
    this.isPublicView =
      publicRaw === "1" || publicRaw === "true" || publicRaw === "yes";

    const presaleRaw = String(
      (options && (options.presale || options.pre_sale || options.is_presale)) ||
        "",
    )
      .trim()
      .toLowerCase();
    this.isPreSaleView =
      presaleRaw === "1" || presaleRaw === "true" || presaleRaw === "yes";
    this.preSaleViewKey = String(
      (options && (options.view_key || options.viewKey)) || "",
    ).trim();
    if (this.isPreSaleView) {
      // 预售详情按“公开详情”样式展示，隐藏底部操作栏/分享等入口
      this.isPublicView = true;
    }
    this.publicFromAgentId =
      Number(options && (options.from_agent_id || options.agent_id || 0)) || 0;
    this.publicFromStyle = Number(options && options.from_style) || 0;

    const id =
      Number(options && (options.id || options.ID || options.property_id)) || 0;
    this.propertyId = id;
    if (!this.propertyId) {
      uni.showToast({ title: "房源ID缺失", icon: "none" });
      return;
    }
    if (!this.isPublicView || this.isPreSaleView) this.ensureCanManageProperties();
    this.loadDetail();
  },
  computed: {
    canEditThisProperty() {
      if (this.isPublicView && !this.isPreSaleView) return false;
      if (!this.canManageProperties) return false;
      // 已售(sold) 作为终态：小程序端只读，不提供编辑入口
      const saleStatus = String(
        (this.property && this.property.sale_status) || "",
      ).trim();
      if (saleStatus === "sold") return false;
      const pid = Number(this.propertyId || 0) || 0;
      if (!pid) return false;
      const agentId = Number(this.property && this.property.agent_id) || 0;
      const uid = Number(this.currentUserId || 0) || 0;
      // 最小权限：仅允许编辑自己维护(agent_id=自己)的房源
      if (!uid || !agentId) return false;
      return agentId === uid;
    },
    showDownloadBtn() {
      if (this.isPublicView) return false;
      const item = this.getCurrentBannerItem();
      if (!item) return false;
      if (item.type === "video") return !!this.allowVideoDownload;
      if (item.type === "image") return !!this.allowImageDownload;
      return false;
    },
    downloadLabel() {
      const item = this.getCurrentBannerItem();
      if (!item) return "下载";
      return item.type === "video" ? "下载视频" : "下载图片";
    },
    hasLocation() {
      const lat = Number(this.mapLat);
      const lng = Number(this.mapLng);
      return !!lat && !!lng && isFinite(lat) && isFinite(lng);
    },
    shareCover() {
      return this.getShareCover();
    },
  },
  methods: {
    openShareSheet() {
      this.shareSheetOpen = true;
    },
    closeShareSheet() {
      this.shareSheetOpen = false;
    },
    onTapForward() {
      this.closeShareSheet();
    },
    closeShowingSheet() {
      this.showingSheetOpen = false;
      this.showingSaving = false;
    },
    async ensureCanManageProperties() {
      // 仅用于 UI 控制；最终权限以后端校验为准
      // 公开详情（非预售）一律不展示“编辑”入口
      if (this.isPublicView && !this.isPreSaleView) {
        this.canManageProperties = false;
        this.currentUserId = 0;
        return;
      }
      const userStore = $store("user");
      const token = uni.getStorageSync("token");
      if (!token && !userStore.isLogin) {
        this.canManageProperties = false;
        this.currentUserId = 0;
        return;
      }
      try {
        const ui = userStore.userInfo || {};
        // 兼容：老版本缓存的 userInfo 可能没有 can_manage_properties，需要强制刷新一次
        if (!ui.id || typeof ui.can_manage_properties === "undefined") {
          await userStore.getInfo();
        }
      } catch (e) {}
      const u = userStore.userInfo || {};
      this.currentUserId = Number(u.id || 0) || 0;
      this.canManageProperties = Number(u.can_manage_properties) === 1;
      try {
        console.log("[property_detail] perm", {
          id: this.currentUserId,
          can_manage_properties: u.can_manage_properties,
        });
      } catch (e) {}
    },
    goEdit() {
      if (!this.propertyId) return;
      uni.navigateTo({
        url: `/pages/property_manage/property_edit?id=${this.propertyId}`,
      });
    },
    debugEditState() {
      try {
        const uid = Number(this.currentUserId || 0) || 0;
        const can = !!this.canManageProperties;
        const agentId = Number(this.property && this.property.agent_id) || 0;
        const show = !!this.canEditThisProperty;
        uni.showModal({
          title: "编辑按钮诊断",
          content: [
            `can_manage_properties: ${can ? 1 : 0}`,
            `currentUserId: ${uid}`,
            `property.agent_id: ${agentId}`,
            `canEditThisProperty: ${show ? 1 : 0}`,
            show ? "✅ 满足条件应显示" : "❌ 条件不满足所以不显示",
          ].join("\n"),
          showCancel: false,
        });
      } catch (e) {
        // ignore
      }
    },
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
        (this.property && this.property.community_name) || "",
      ).trim();
      const addr = String(
        (this.property && this.property.address) || "",
      ).trim();
      if (name && addr) return `${name} · ${addr}`;
      return name || addr || "暂无位置信息";
    },
    normalizeImage(url) {
      const imageUrl = String(url || "").trim();
      if (!imageUrl) return "";
      if (imageUrl.indexOf("/static/images/") === 0) return "";
      return imageUrl;
    },
    normalizeRenovationStageStatus(v) {
      const s = String(v || "").trim().toLowerCase();
      if (s === "done" || s === "finished" || s === "completed") return "done";
      if (s === "doing" || s === "in_progress" || s === "progress")
        return "doing";
      return "todo";
    },
    parseRenovationStageLogs(raw) {
      let arr = raw;
      if (typeof raw === "string") {
        const s = raw.trim();
        if (s && s.indexOf("[") === 0) {
          try {
            arr = JSON.parse(s);
          } catch (e) {
            arr = [];
          }
        } else {
          arr = [];
        }
      }
      if (!Array.isArray(arr)) return [];
      const out = [];
      arr.forEach((it) => {
        const stage = String(
          (it && (it.stage || it.stage_name || it.name)) || "",
        ).trim();
        if (!stage) return;
        const date = String((it && it.date) || "").trim();
        const note = String((it && (it.note || it.notes)) || "").trim();
        const st = this.normalizeRenovationStageStatus(it && it.status);

        let imgs = [];
        if (Array.isArray(it && it.images)) imgs = it.images;
        else if (typeof (it && it.images) === "string") {
          imgs = String(it.images || "")
            .split(",")
            .map((x) => String(x || "").trim())
            .filter(Boolean);
        }
        const normalizedImgs = imgs
          .map((u) => this.normalizeImage(u))
          .filter(Boolean);
        out.push({
          stage,
          status: st,
          date,
          note,
          images: Array.from(new Set(normalizedImgs)),
        });
      });
      return out;
    },
    buildRenovationStageTimeline({ overallStatus, currentStage, logs }) {
      const fallback = [
        "设计",
        "拆改",
        "水电",
        "泥瓦",
        "木工",
        "油漆",
        "安装",
        "软装",
        "验收",
      ];
      const order = fallback.slice(0);
      (Array.isArray(logs) ? logs : []).forEach((it) => {
        const s = String((it && it.stage) || "").trim();
        if (!s) return;
        if (order.indexOf(s) === -1) order.push(s);
      });

      const overall = String(overallStatus || "").trim();
      const cur = String(currentStage || "").trim();
      const curIdx = cur ? order.indexOf(cur) : -1;

      return order.map((stage, idx) => {
        const found = (Array.isArray(logs) ? logs : []).find(
          (x) => x && x.stage === stage,
        );
        let st = found ? String(found.status || "").trim() : "todo";
        const date = found ? String(found.date || "").trim() : "";
        const note = found ? String(found.note || "").trim() : "";
        const images = found && Array.isArray(found.images) ? found.images : [];

        if (!found) {
          if (overall === "done") st = "done";
          else if (overall === "none") st = "todo";
          else if (overall === "in_progress" && curIdx >= 0) {
            if (idx < curIdx) st = "done";
            else if (idx === curIdx) st = "doing";
            else st = "todo";
          }
        } else {
          // 总状态兜底覆盖（避免出现“已完工但某阶段仍未开始”）
          if (overall === "done") st = "done";
          if (overall === "none") st = "todo";
        }
        st = this.normalizeRenovationStageStatus(st);
        return { stage, status: st, date, note, images };
      });
    },
    stageStatusText(s) {
      const v = String(s || "").trim();
      if (v === "done") return "已完成";
      if (v === "doing") return "进行中";
      return "未开始";
    },
    previewStageImage(item, idx) {
      const list = item && Array.isArray(item.images) ? item.images : [];
      if (!list.length) return;
      const current = list[idx] || list[0];
      uni.previewImage({ current, urls: list });
    },
    getShareCover() {
      const cover = this.normalizeImage(this.property && this.property.cover_image);
      if (cover) return cover;
      const first = Array.isArray(this.images) ? String(this.images[0] || "").trim() : "";
      return this.normalizeImage(first);
    },
    async loadDetail() {
      if (this.loading || !this.propertyId) return false;
      this.loading = true;
      let res;
      try {
        if (this.isPreSaleView) {
          res = await propertyApi.getPreSaleDetail({
            id: this.propertyId,
            view_key: this.preSaleViewKey,
          });
        } else {
          res = await propertyApi.getDetail({
            id: this.propertyId,
            public: this.isPublicView ? 1 : 0,
          });
        }
      } catch (e) {
        this.loading = false;
        if (
          (this.isPreSaleView || !this.isPublicView) &&
          !uni.getStorageSync("token")
        ) {
          uni.reLaunch({ url: "/pages/login/login" });
        }
        return false;
      }
      this.loading = false;
      if (!res || res.code !== 0) {
        const msg = String((res && (res.message || res.msg)) || "").trim();
        if (
          (this.isPreSaleView || !this.isPublicView) &&
          /登录|token|认证/i.test(msg)
        ) {
          uni.reLaunch({ url: "/pages/login/login" });
          return false;
        }
        setTimeout(() => {
          this.goBack();
        }, 600);
        return false;
      }

      const data = res.data || {};
      if (typeof data.public_view !== "undefined") {
        this.isPublicView = !!data.public_view;
      }
      const p = data.property || {};
      // 兼容：后端也会回传 current_user（已登录用户ID）
      if (!this.currentUserId)
        this.currentUserId = Number(data.current_user || 0) || 0;
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
      if (this.isPublicView) this.ownerPhone = "";

      this.tryLogActivity("view", this.propertyId, {
        count: 1,
        from: this.isPublicView ? "public" : "agent",
      });

      // 下载权限（后端默认：1 允许）
      this.allowImageDownload = Number(p.allow_image_download) !== 0;
      this.allowVideoDownload = Number(p.allow_video_download) !== 0;

      // 视频（单个）
      this.videoUrl = String(p.video_url || "").trim();
      this.videoPoster = this.normalizeImage(p.cover_image);

      const imgs = Array.isArray(data.images)
        ? data.images
        : Array.isArray(p.images)
          ? p.images
          : [];
      const normalizedImgs = imgs
        .map((u) => this.normalizeImage(u))
        .filter(Boolean);
      this.images = normalizedImgs;
      this.currentSwiper = 0;

      // Banner：优先展示视频（若存在）
      const items = [];
      if (this.videoUrl) items.push({ type: "video", src: this.videoUrl });
      (this.images || []).forEach((u) => items.push({ type: "image", src: u }));
      this.bannerItems = items;

      // 地图坐标（小程序 map 组件）
      const lat = Number(p.latitude);
      const lng = Number(p.longitude);
      this.mapLat = isFinite(lat) ? lat : 0;
      this.mapLng = isFinite(lng) ? lng : 0;
      this.mapMarkers =
        this.mapLat && this.mapLng
          ? [
              {
                id: 1,
                latitude: this.mapLat,
                longitude: this.mapLng,
                width: 26,
                height: 26,
              },
            ]
          : [];

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
      const stageLogs = this.parseRenovationStageLogs(r.stage_logs);
      const stageTimeline = this.buildRenovationStageTimeline({
        overallStatus: status,
        currentStage: String(r.current_stage || "").trim(),
        logs: stageLogs,
      });
      this.renovation = {
        status,
        subtitle: String(r.subtitle || "房源装修情况").trim() || "房源装修情况",
        progress: Number(r.progress_percentage || 0) || 0,
        stage: String(r.current_stage || "").trim() || "—",
        eta: String(r.estimated_finish_date || "").trim() || "—",
        finishAt: String(r.actual_finish_date || "").trim() || "—",
        materials: Array.isArray(r.materials) ? r.materials : [],
        note: String(r.notes || "").trim() || "—",
        images: Array.isArray(r.images)
          ? r.images.map((u) => this.normalizeImage(u)).filter(Boolean)
          : [],
        timeline: stageTimeline,
      };

      const rec = Array.isArray(data.recommends) ? data.recommends : [];
      this.recommends = rec.map((it) => ({
        id: it.id,
        name: it.title || it.name || "房源",
        rooms: this.getLayoutText(it),
        price: `${it.price || "-"}${it.price_unit || ""}`,
        size: it.area || "-",
        image: this.normalizeImage(it.image || it.cover_image || ""),
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
      if (!this.isPublicView) {
        uni.navigateBack();
        return;
      }
      const pages = getCurrentPages ? getCurrentPages() : [];
      if (pages && pages.length > 1) {
        uni.navigateBack();
        return;
      }
      const aid = Number(this.publicFromAgentId || 0);
      const style = Number(this.publicFromStyle || 0);
      if (aid > 0) {
        uni.reLaunch({
          url: `/pages/agent_public_card/agent_public_card?agent_id=${encodeURIComponent(
            aid,
          )}&style=${encodeURIComponent(style)}`,
        });
        return;
      }
      uni.reLaunch({ url: "/pages/login/login" });
    },
    goToRec(rec) {
      const id = Number(rec && (rec.id || rec.ID)) || 0;
      if (!id) return;
      if (id === this.propertyId) return;
      const publicQuery = this.isPublicView
        ? `&public=1&from_agent_id=${encodeURIComponent(
            this.publicFromAgentId,
          )}&from_style=${encodeURIComponent(this.publicFromStyle)}`
        : "";
      uni.redirectTo({
        url: `/pages/property_detail/property_detail?id=${encodeURIComponent(
          id,
        )}${publicQuery}`,
      });
    },
    swiperChange(e) {
      this.currentSwiper = e.detail.current;
    },
    getCurrentBannerItem() {
      const list = Array.isArray(this.bannerItems) ? this.bannerItems : [];
      const idx = Number(this.currentSwiper) || 0;
      return list[idx] || null;
    },

    // 申请相册权限（用于保存图片/视频）
    ensureAlbumPermission() {
      return new Promise((resolve, reject) => {
        uni.getSetting({
          success: (st) => {
            const ok =
              st && st.authSetting && st.authSetting["scope.writePhotosAlbum"];
            if (ok) return resolve(true);
            uni.authorize({
              scope: "scope.writePhotosAlbum",
              success: () => resolve(true),
              fail: () => {
                uni.showModal({
                  title: "需要相册权限",
                  content: "保存到相册需要授权，请在设置中开启相册权限。",
                  confirmText: "去设置",
                  success: (r) => {
                    if (r.confirm) {
                      uni.openSetting({
                        success: () => resolve(true),
                        fail: () => reject(false),
                      });
                    } else {
                      reject(false);
                    }
                  },
                });
              },
            });
          },
          fail: () => reject(false),
        });
      });
    },

    async downloadCurrent() {
      if (this.isPublicView) return;
      const item = this.getCurrentBannerItem();
      if (!item) return;
      if (item.type === "video" && !this.allowVideoDownload) return;
      if (item.type === "image" && !this.allowImageDownload) return;

      try {
        await this.ensureAlbumPermission();
      } catch (e) {
        return;
      }

      const url = String(item.src || "").trim();
      if (!url) return;
      uni.showLoading({ title: "下载中..." });
      uni.downloadFile({
        url,
        success: (d) => {
          const fp = d && d.tempFilePath;
          if (!fp) {
            uni.hideLoading();
            uni.showToast({ title: "下载失败", icon: "none" });
            return;
          }
          if (item.type === "video") {
            uni.saveVideoToPhotosAlbum({
              filePath: fp,
              success: () =>
                uni.showToast({ title: "已保存视频", icon: "success" }),
              fail: () => uni.showToast({ title: "保存失败", icon: "none" }),
              complete: () => uni.hideLoading(),
            });
          } else {
            uni.saveImageToPhotosAlbum({
              filePath: fp,
              success: () =>
                uni.showToast({ title: "已保存图片", icon: "success" }),
              fail: () => uni.showToast({ title: "保存失败", icon: "none" }),
              complete: () => uni.hideLoading(),
            });
          }
        },
        fail: () => {
          uni.hideLoading();
          uni.showToast({ title: "下载失败", icon: "none" });
        },
      });
    },
    handleShare() {
      if (this.isPublicView) return;
      this.openShareSheet();
    },
    recordShowing() {
      if (this.isPublicView) return;
      const pid = Number(this.propertyId || 0) || 0;
      if (!pid) {
        uni.showToast({ title: "房源ID缺失", icon: "none" });
        return;
      }
      this.showingForm = { client_name: "", client_phone: "" };
      this.showingSheetOpen = true;
    },
    async submitShowing() {
      if (this.showingSaving) return;
      const pid = Number(this.propertyId || 0) || 0;
      if (!pid) {
        uni.showToast({ title: "房源ID缺失", icon: "none" });
        return;
      }
      const name = String(
        (this.showingForm && this.showingForm.client_name) || "",
      ).trim();
      const rawPhone = String(
        (this.showingForm && this.showingForm.client_phone) || "",
      ).trim();
      const phone = rawPhone.replace(/\D/g, "");

      if (!name) {
        uni.showToast({ title: "请输入客户姓名", icon: "none" });
        return;
      }
      if (!phone) {
        uni.showToast({ title: "请输入客户电话", icon: "none" });
        return;
      }
      if (phone.length < 6) {
        uni.showToast({ title: "客户电话格式不正确", icon: "none" });
        return;
      }

      this.showingSaving = true;
      try {
        await this.tryLogActivity("showing", pid, {
          client_name: name,
          client_phone: phone,
        });
        this.showingSheetOpen = false;
        uni.showToast({ title: "已记录带看", icon: "none" });
      } finally {
        this.showingSaving = false;
      }
    },
    async copyPromoteText() {
      const pid = Number(this.propertyId || 0) || 0;
      const title =
        this.property && this.property.title ? this.property.title : "房源";
      this.closeShareSheet();
      if (!pid) {
        uni.showToast({ title: "房源ID缺失", icon: "none" });
        return;
      }
      this.tryLogActivity("share", pid, { channel: "复制文案" });

      let urlLink = "";
      try {
        const linkRes = await propertyApi.getPropertyUrlLink(
          { property_id: pid },
          true,
        );
        if (linkRes && linkRes.code === 0 && linkRes.data) {
          urlLink = String(linkRes.data.url_link || "").trim();
        }
      } catch (e) {
        urlLink = "";
      }
      if (!urlLink) {
        urlLink = `/pages/property_detail/property_detail?id=${encodeURIComponent(
          pid,
        )}&public=1`;
      }

      const layoutText = this.getLayoutText(this.property || {});
      const areaText =
        this.property && this.property.area ? `${this.property.area}㎡` : "";
      const priceText =
        this.property && this.property.price
          ? `¥${this.property.price}${this.property.price_unit || ""}`
          : "";
      const line2 = [layoutText, areaText, priceText]
        .filter(Boolean)
        .join(" | ");

      const text = `【优质房源推荐】${title}${line2 ? `\n${line2}` : ""}\n小程序链接：${urlLink}`;
      uni.setClipboardData({
        data: text,
        success: () => uni.showToast({ title: "已复制推广文案", icon: "none" }),
        fail: () => uni.showToast({ title: "复制失败", icon: "none" }),
      });
    },
    openMap() {
      if (this.isPublicView) {
        uni.showToast({ title: "公开页不支持跳转地图", icon: "none" });
        return;
      }
      const lat = Number(this.property && this.property.latitude);
      const lng = Number(this.property && this.property.longitude);
      const name = String(
        (this.property &&
          (this.property.community_name || this.property.title)) ||
          "房源",
      ).trim();
      const address = String(
        (this.property && this.property.address) || "",
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
      if (this.isPublicView) return;
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
      if (this.isPublicView) return;
      if (!this.ownerPhone) {
        uni.showToast({ title: "暂无业主电话", icon: "none" });
        return;
      }
      this.tryLogActivity("call", this.propertyId, {
        type: "呼出",
        phone: this.ownerPhone,
        target_name: "业主",
      });
      uni.makePhoneCall({
        phoneNumber: this.ownerPhone,
      });
    },
    openLock() {
      if (this.isPublicView) return;
      const ss = String(
        (this.property && this.property.sale_status) || "",
      ).trim();
      if (ss === "sold" || ss === "off_market") {
        const label =
          this.property && this.property.sale_status_label
            ? this.property.sale_status_label
            : "不可操作";
        uni.showToast({ title: `房源${label}，暂不可开锁`, icon: "none" });
        return;
      }
      if (!this.hasSmartLock) {
        uni.showToast({ title: "该房源未绑定智能锁，无法开锁", icon: "none" });
        return;
      }
      // 进入开锁流程页（页面内可选择蓝牙开锁 / 获取密码开锁）
      uni.navigateTo({
        url: `/pages/unlock_steps/unlock_steps?property_id=${encodeURIComponent(
          this.propertyId,
        )}`,
      });
    },
    async tryLogActivity(activityType, propertyId, meta) {
      if (this.isPublicView) return;
      if (!uni.getStorageSync("token")) return;
      const type = String(activityType || "").trim();
      if (!type) return;
      if (type === "view") {
        if (this.viewLogged) return;
        this.viewLogged = true;
      }
      const pid = Number(propertyId || 0) || 0;
      try {
        await userApi.addWorkbenchActivityLog(
          {
            activity_type: type,
            property_id: pid,
            page: "property_detail",
            meta: meta || {},
          },
          false,
        );
      } catch (e) {
        // 日志写入失败不影响主流程
      }
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
  justify-content: space-between;
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

  .banner-video {
    width: 100%;
    height: 100%;
  }

  .banner-actions {
    position: absolute;
    bottom: 64rpx;
    left: 32rpx;
    z-index: 3;
  }

  .dl-pill {
    display: flex;
    align-items: center;
    gap: 10rpx;
    padding: 12rpx 18rpx;
    border-radius: 999rpx;
    background: rgba(255, 255, 255, 0.88);
    backdrop-filter: blur(10rpx);
    border: 1px solid rgba(15, 23, 42, 0.08);
    box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.18);
  }

  .dl-icon {
    font-size: 34rpx;
    color: #0f766e;
  }

  .dl-text {
    font-size: 24rpx;
    font-weight: 700;
    color: #0f172a;
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

    .action-col {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 12rpx;
      flex-shrink: 0;
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
      line-height: 1;

      &::after {
        border: none;
      }

      .share-icon {
        font-size: 28rpx;
      }

      &:active {
        transform: scale(0.98);
        background: rgba(37, 99, 235, 0.12);
      }
    }

    .edit-btn {
      display: flex;
      align-items: center;
      gap: 8rpx;
      padding: 12rpx 16rpx;
      background: rgba(15, 23, 42, 0.06);
      border: 1rpx solid rgba(15, 23, 42, 0.1);
      border-radius: 18rpx;
      font-size: 24rpx;
      color: #0f172a;
      font-weight: 700;
      flex-shrink: 0;

      .edit-icon {
        font-size: 28rpx;
      }

      &:active {
        transform: scale(0.98);
        background: rgba(15, 23, 42, 0.1);
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

  .commission-line {
    display: flex;
    align-items: center;
    gap: 8rpx;
    padding: 8rpx 12rpx;
    border-radius: 14rpx;
    background: rgba(249, 115, 22, 0.08);
    border: 1rpx solid rgba(249, 115, 22, 0.14);
    color: #9a3412;
    font-size: 22rpx;
    font-weight: 800;
    line-height: 1;

    .commission-ic {
      font-size: 26rpx;
      flex-shrink: 0;
      color: #f97316;
    }

    .commission-text {
      flex: 1;
      min-width: 0;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .status-row {
    display: flex;
    flex-wrap: wrap;
    gap: 12rpx;
    align-items: center;
  }

  .status-chip {
    display: flex;
    align-items: center;
    gap: 8rpx;
    padding: 8rpx 14rpx;
    border-radius: 999rpx;
    font-size: 22rpx;
    font-weight: 700;
    background: rgba(15, 23, 42, 0.06);
    border: 1rpx solid rgba(15, 23, 42, 0.1);
    color: #0f172a;

    .status-ic {
      font-size: 26rpx;
    }

    &.on_sale {
      background: rgba(34, 197, 94, 0.12);
      border-color: rgba(34, 197, 94, 0.18);
      color: #16a34a;
    }
    &.sold {
      background: rgba(234, 88, 12, 0.12);
      border-color: rgba(234, 88, 12, 0.18);
      color: #ea580c;
    }
    &.off_market {
      background: rgba(148, 163, 184, 0.18);
      border-color: rgba(148, 163, 184, 0.22);
      color: #64748b;
    }

    &.hot {
      background: rgba(37, 99, 235, 0.1);
      border-color: rgba(37, 99, 235, 0.16);
      color: #2563eb;
    }
  }
}

.contact-card {
  border-radius: 24rpx;
  background: linear-gradient(
    135deg,
    rgba(15, 23, 42, 0.02),
    rgba(37, 99, 235, 0.03)
  );
  border: 1rpx solid rgba(226, 232, 240, 0.9);
  padding: 18rpx 18rpx 8rpx;
}

.contact-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 14rpx;
}

.contact-item {
  background-color: #ffffff;
  border-radius: 18rpx;
  border: 1rpx solid rgba(226, 232, 240, 0.9);
  padding: 14rpx 14rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
}

.contact-item.wide {
  grid-column: span 2;
}

.contact-label {
  font-size: 22rpx;
  color: #64748b;
  font-weight: 600;
}

.contact-value {
  font-size: 26rpx;
  color: #0f172a;
  font-weight: 800;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.contact-value.price {
  color: #ea580c;
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

.desc-card {
  background-color: #f8fafc;
  padding: 32rpx;
  border-radius: 24rpx;
}

.desc-text {
  font-size: 28rpx;
  color: #0f172a;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
}

.desc-empty {
  font-size: 26rpx;
  color: #94a3b8;
  background-color: #f8fafc;
  padding: 24rpx 32rpx;
  border-radius: 24rpx;
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

.reno-stage-timeline {
  margin-top: 22rpx;
  padding-top: 20rpx;
  border-top: 1rpx dashed rgba(148, 163, 184, 0.35);
}

.reno-stage-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14rpx;
}

.reno-stage-title {
  font-size: 26rpx;
  font-weight: 900;
  color: #0f172a;
}

.reno-stage-tip {
  font-size: 22rpx;
  color: #94a3b8;
  font-weight: 700;
}

.stage-list {
  margin-top: 16rpx;
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}

.stage-item {
  display: flex;
  align-items: stretch;
  gap: 14rpx;
}

.stage-rail {
  width: 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  padding-top: 10rpx;
}

.stage-dot {
  width: 18rpx;
  height: 18rpx;
  border-radius: 999rpx;
  background: rgba(148, 163, 184, 0.35);
  border: 4rpx solid rgba(148, 163, 184, 0.18);
  box-sizing: border-box;

  &.doing {
    background: rgba(37, 99, 235, 0.9);
    border-color: rgba(37, 99, 235, 0.18);
  }

  &.done {
    background: rgba(34, 197, 94, 0.88);
    border-color: rgba(34, 197, 94, 0.18);
  }
}

.stage-line {
  width: 6rpx;
  flex: 1;
  margin-top: 10rpx;
  border-radius: 999rpx;
  background: rgba(226, 232, 240, 0.9);
}

.stage-body {
  flex: 1;
  min-width: 0;
  background-color: #ffffff;
  border-radius: 22rpx;
  padding: 18rpx;
  border: 1rpx solid rgba(226, 232, 240, 0.95);
}

.stage-head {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10rpx;
}

.stage-name {
  font-size: 26rpx;
  font-weight: 900;
  color: #0f172a;
}

.stage-tag {
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  font-weight: 900;
  letter-spacing: 0.2rpx;
  background: rgba(148, 163, 184, 0.18);
  color: #475569;

  &.doing {
    background: rgba(37, 99, 235, 0.14);
    color: #1d4ed8;
  }
  &.done {
    background: rgba(34, 197, 94, 0.14);
    color: #15803d;
  }
}

.stage-date {
  margin-left: auto;
  font-size: 22rpx;
  color: #94a3b8;
  font-weight: 700;
}

.stage-note {
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #334155;
  font-weight: 600;
  line-height: 1.6;
  word-break: break-all;
}

.stage-img-grid {
  margin-top: 12rpx;
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10rpx;
}

.stage-thumb {
  width: 100%;
  height: 160rpx;
  border-radius: 18rpx;
  overflow: hidden;
  background-color: rgba(226, 232, 240, 0.9);
  border: 1rpx solid rgba(226, 232, 240, 0.95);
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

  .map-native {
    width: 100%;
    height: 100%;
  }

  .map-image {
    width: 100%;
    height: 100%;
    opacity: 0.8;
  }
  .map-image-empty {
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #eef2f7, #e2e8f0);
    .material-symbols-outlined {
      font-size: 64rpx;
      color: #94a3b8;
    }
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

    .map-coords {
      position: absolute;
      bottom: 16rpx;
      right: 16rpx;
      padding: 8rpx 14rpx;
      border-radius: 999rpx;
      background-color: rgba(255, 255, 255, 0.92);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(15, 23, 42, 0.08);
      box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.14);
      .coord-text {
        font-size: 20rpx;
        color: #0f172a;
        font-weight: 600;
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
      .rec-image-empty {
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #eef2f7, #e2e8f0);
        .material-symbols-outlined {
          font-size: 52rpx;
          color: #94a3b8;
        }
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

.share-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.45);
  z-index: 1200;
  display: flex;
  align-items: flex-end;
}

.share-sheet {
  width: 100%;
  background-color: #ffffff;
  border-radius: 32rpx 32rpx 0 0;
  padding: 18rpx 18rpx calc(env(safe-area-inset-bottom) + 18rpx);
  box-sizing: border-box;
}

.share-title {
  font-size: 28rpx;
  font-weight: 800;
  color: #0f172a;
  padding: 12rpx 12rpx 18rpx;
}

.share-actions {
  border-radius: 24rpx;
  overflow: hidden;
  border: 1px solid #f1f5f9;
}

.share-action {
  height: 92rpx;
  padding: 0 18rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #334155;
  background: #ffffff;
  border: none;
  text-align: left;
  line-height: 1;
}

.share-action.primary {
  color: #2d9cf0;
}

.share-action::after {
  border: none;
}

.share-divider {
  height: 1px;
  background-color: #f1f5f9;
}

.share-cancel {
  margin-top: 16rpx;
  height: 88rpx;
  border-radius: 24rpx;
  background: #f1f5f9;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 800;
  color: #334155;
}

.sheet-mask {
  position: fixed;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(15, 23, 42, 0.45);
  z-index: 1200;
  display: flex;
  align-items: flex-end;
}

.sheet-panel {
  width: 100%;
  background-color: #ffffff;
  border-radius: 32rpx 32rpx 0 0;
  padding: 18rpx 18rpx calc(env(safe-area-inset-bottom) + 18rpx);
  box-sizing: border-box;
}

.sheet-title {
  font-size: 28rpx;
  font-weight: 900;
  color: #0f172a;
  padding: 12rpx 12rpx 6rpx;
}

.sheet-tip {
  padding: 0 12rpx 18rpx;
  font-size: 24rpx;
  color: #64748b;
}

.sheet-form {
  border-radius: 24rpx;
  overflow: hidden;
  border: 1px solid #f1f5f9;
  background: #ffffff;
}

.sheet-field {
  padding: 18rpx;
}

.sheet-label {
  display: block;
  font-size: 24rpx;
  font-weight: 800;
  color: #64748b;
  margin-bottom: 12rpx;
}

.sheet-input {
  height: 84rpx;
  border-radius: 18rpx;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  padding: 0 18rpx;
  box-sizing: border-box;
  font-size: 28rpx;
  color: #0f172a;
}

.sheet-divider {
  height: 1px;
  background-color: #f1f5f9;
}

.sheet-actions-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 16rpx;
}

.sheet-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 900;
  border: 1px solid transparent;
}

.sheet-btn.cancel {
  background: #f1f5f9;
  color: #334155;
}

.sheet-btn.primary {
  background: #2d9cf0;
  color: #ffffff;
  box-shadow: 0 8rpx 20rpx rgba(45, 156, 240, 0.26);
}

.sheet-btn.disabled {
  opacity: 0.65;
}
</style>
