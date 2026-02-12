import { createRouter, createWebHistory,createWebHashHistory } from 'vue-router';
import NProgress from 'nprogress'; // progress bar
import 'nprogress/nprogress.css';

import { appRoutes } from './routes';
import { REDIRECT_MAIN, NOT_FOUND_ROUTE } from './routes/base';
import createRouteGuard from './guard';
NProgress.configure({ showSpinner: false }); // NProgress Configuration
const router = createRouter({
  //go二级目录部署，位置在resource/webbusiness下
  history: createWebHashHistory(import.meta.env.BASE_URL),//hash模式带#号，自动适配 Vite base（如 /webbusiness/）
  //2独立域名部署
  // history: createWebHistory(),//history模式
  routes: [
    {
      path: '/',
      redirect: '/'+(window?.globalConfig.RouterHome||"home"),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/index.vue'),
      meta: {
        requiresAuth: false,
      },
    },
    ...appRoutes,
    REDIRECT_MAIN,
    NOT_FOUND_ROUTE,
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

createRouteGuard(router);

export default router;
