import $api from '@/api';  
import $router from '@/utils/router';
import $helper from '@/utils/helper';
import wechat from '@/utils/wechat/wechat';

const utils = {
  $api,
  $router,
  $helper
};

// 加载Nx底层依赖
export async function GFInit() {
  if (process.env.NODE_ENV === 'development') {
    GFDebug();
  }
  // #ifdef MP-WEIXIN
  // 检测小程序更新
   wechat.checkMiniProgramUpdate();
  // #endif
  
}

// 开发模式
function GFDebug() {
  // 开发环境引入vconsole调试
  // #ifdef H5
  // import("vconsole").then(vconsole => {
  // 	new vconsole.default();
  // });
  // #endif
}

export default utils;
