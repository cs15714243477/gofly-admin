"use strict";
const common_vendor = require("../../common/vendor.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      statusBarHeight: 0,
      headerTop: 0,
      // 0：顶部（暗色渐变遮罩） -> 1：下滑后（磨砂白底）
      headerOpacity: 0,
      currentSwiper: 0,
      isFollowed: false,
      hasSmartLock: true,
      ownerPhone: "13888889999",
      images: [
        "/static/images/img_cdc09ae543.png",
        "/static/images/img_cdc09ae543.png"
      ],
      attributes: [
        { label: "单价", value: "9,086元/m²" },
        { label: "楼层", value: "中层 / 共32层" },
        { label: "朝向", value: "南北通透" },
        { label: "装修", value: "精装修" },
        { label: "年代", value: "2019年建" },
        { label: "产权", value: "商品房" }
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
        { key: "done", label: "已完成" }
      ],
      renovation: {
        status: "in_progress",
        // none | in_progress | done
        subtitle: "高端精装，环保材料（示例）",
        progress: 68,
        stage: "水电验收完成，正在进行泥木施工",
        eta: "预计 2026-02-20",
        finishAt: "2026-01-10",
        materials: ["圣象地板", "马可波罗瓷砖", "多乐士乳胶漆"],
        note: "客厅墙面已完成找平与底漆，卫生间防水已做闭水试验；全屋线管/强弱电分离施工完成。",
        images: [
          "/static/images/img_08b712f810.png",
          "/static/images/img_8642388cea.png"
        ]
      },
      recommends: [
        {
          name: "阳光花园三期 2室",
          rooms: "2室1厅",
          price: "72万",
          size: "89",
          image: "/static/images/img_71f4809787.png"
        },
        {
          name: "金地名津 精装大三房",
          rooms: "3室2厅",
          price: "98万",
          size: "96",
          image: "/static/images/img_1112597ba0.png"
        },
        {
          name: "万科四季花城",
          rooms: "3室2厅",
          price: "105万",
          size: "110",
          image: "/static/images/img_662c3c598a.png"
        }
      ]
    };
  },
  onLoad() {
    const info = common_vendor.index$1.getSystemInfoSync();
    this.statusBarHeight = info.statusBarHeight;
    try {
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.getMenuButtonBoundingClientRect) {
        const rect = common_vendor.wx$1.getMenuButtonBoundingClientRect();
        const w = Number(info.windowWidth || 375);
        const rpx2px = w / 750;
        const circleBtnSizePx = 80 * rpx2px;
        const capsuleCenterY = Number(rect.top || 0) + Number(rect.height || 0) / 2;
        const top = capsuleCenterY - circleBtnSizePx / 2;
        this.headerTop = Math.max(Number(this.statusBarHeight || 0), Math.round(top));
      } else {
        this.headerTop = this.statusBarHeight + 12;
      }
    } catch (e) {
      this.headerTop = this.statusBarHeight + 12;
    }
  },
  methods: {
    onPageScroll(e) {
      const st = Number(e && e.detail && e.detail.scrollTop || 0);
      const fadeDistance = 160;
      const next = Math.min(1, st / fadeDistance);
      this.headerOpacity = Math.max(0, Math.min(1, Number(next.toFixed(3))));
    },
    goBack() {
      common_vendor.index$1.navigateBack();
    },
    swiperChange(e) {
      this.currentSwiper = e.detail.current;
    },
    handleShare() {
      const title = "129-1# 全屋智能 拎包入住 随时看房 南北通透";
      common_vendor.index$1.showActionSheet({
        itemList: ["生成获客海报", "复制房源链接", "模拟分享（占位）"],
        success: (res) => {
          if (res.tapIndex === 0) {
            common_vendor.index$1.navigateTo({
              url: `/pages/marketing_posters/marketing_posters?title=${encodeURIComponent(title)}`
            });
            return;
          }
          if (res.tapIndex === 1) {
            const link = `/pages/property_detail/property_detail`;
            common_vendor.index$1.setClipboardData({
              data: `${title}
${link}`,
              success: () => common_vendor.index$1.showToast({ title: "已复制链接", icon: "none" })
            });
            return;
          }
          if (res.tapIndex === 2) {
            common_vendor.index$1.showToast({ title: "分享功能待接入", icon: "none" });
          }
        }
      });
    },
    openMap() {
      common_vendor.index$1.openLocation({
        latitude: 39.9042,
        longitude: 116.4074,
        name: "幸福里小区",
        address: "北京市",
        fail: () => {
          common_vendor.index$1.showToast({ title: "无法打开地图", icon: "none" });
        }
      });
    },
    toggleFollow() {
      this.isFollowed = !this.isFollowed;
      common_vendor.index$1.showToast({
        title: this.isFollowed ? "已关注" : "已取消关注",
        icon: "none"
      });
    },
    callOwner() {
      common_vendor.index$1.makePhoneCall({
        phoneNumber: this.ownerPhone
      });
    },
    openLock() {
      common_vendor.index$1.navigateTo({
        url: this.hasSmartLock ? "/pages/unlock_steps/unlock_steps" : "/pages/unlock_details/unlock_details"
      });
    },
    setRenovationStatus(key) {
      this.renovation.status = key;
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: 1 - $data.headerOpacity,
    b: $data.headerOpacity,
    c: $data.headerOpacity > 0.65 ? 1 : "",
    d: common_vendor.o((...args) => $options.goBack && $options.goBack(...args)),
    e: $data.headerTop + "px",
    f: common_vendor.f($data.images, (item, index, i0) => {
      return {
        a: item,
        b: index
      };
    }),
    g: common_vendor.o((...args) => $options.swiperChange && $options.swiperChange(...args)),
    h: common_vendor.t($data.currentSwiper + 1),
    i: common_vendor.t($data.images.length),
    j: common_vendor.o((...args) => $options.handleShare && $options.handleShare(...args)),
    k: common_vendor.f($data.attributes, (attr, idx, i0) => {
      return {
        a: common_vendor.t(attr.label),
        b: common_vendor.t(attr.value),
        c: idx
      };
    }),
    l: common_vendor.f($data.renovationTabs, (t, k0, i0) => {
      return {
        a: common_vendor.t(t.label),
        b: t.key,
        c: $data.renovation.status === t.key ? 1 : "",
        d: common_vendor.o(($event) => $options.setRenovationStatus(t.key), t.key)
      };
    }),
    m: common_vendor.t($data.renovation.status === "none" ? "未装修" : $data.renovation.status === "in_progress" ? "装修进行中" : "装修完成"),
    n: common_vendor.n($data.renovation.status),
    o: common_vendor.t($data.renovation.subtitle),
    p: $data.renovation.status === "none"
  }, $data.renovation.status === "none" ? {} : $data.renovation.status === "in_progress" ? {
    r: common_vendor.f($data.renovation.images, (img, idx, i0) => {
      return {
        a: img,
        b: idx
      };
    }),
    s: $data.renovation.progress + "%",
    t: common_vendor.t($data.renovation.progress),
    v: common_vendor.t($data.renovation.stage),
    w: common_vendor.t($data.renovation.eta),
    x: common_vendor.f($data.renovation.materials, (m, idx, i0) => {
      return {
        a: common_vendor.t(m),
        b: "m-" + idx
      };
    }),
    y: common_vendor.t($data.renovation.note)
  } : {
    z: common_vendor.t($data.renovation.finishAt),
    A: common_vendor.f($data.renovation.materials, (m, idx, i0) => {
      return {
        a: common_vendor.t(m),
        b: "dm-" + idx
      };
    }),
    B: common_vendor.t($data.renovation.note)
  }, {
    q: $data.renovation.status === "in_progress",
    C: common_vendor.o((...args) => $options.openMap && $options.openMap(...args)),
    D: common_assets._imports_0,
    E: common_vendor.o((...args) => $options.openMap && $options.openMap(...args)),
    F: common_vendor.f($data.recommends, (rec, rIdx, i0) => {
      return {
        a: rec.image,
        b: common_vendor.t(rec.size),
        c: common_vendor.t(rec.name),
        d: common_vendor.t(rec.rooms),
        e: common_vendor.t(rec.price),
        f: rIdx
      };
    }),
    G: common_vendor.o((...args) => $options.onPageScroll && $options.onPageScroll(...args)),
    H: common_vendor.o((...args) => $options.callOwner && $options.callOwner(...args)),
    I: $data.isFollowed
  }, $data.isFollowed ? {} : {}, {
    J: common_vendor.o((...args) => $options.toggleFollow && $options.toggleFollow(...args)),
    K: $data.hasSmartLock
  }, $data.hasSmartLock ? {} : {}, {
    L: common_vendor.t($data.hasSmartLock ? "智能锁" : "密码开锁"),
    M: common_vendor.o((...args) => $options.openLock && $options.openLock(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
_sfc_main.__runtimeHooks = 1;
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/property_detail/property_detail.js.map
