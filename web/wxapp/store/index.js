import { createPinia } from 'pinia';
import piniaPersist from 'pinia-plugin-persist-uni';

// 显式导入（避免 import.meta.globEager 在小程序运行时报错）
import user from './user';
import wxInfo from './wxInfo';

// 单例 pinia，保证在非组件上下文也可安全使用
export const pinia = createPinia();
pinia.use(piniaPersist);

const stores = {
  user,
  wxInfo,
};

export const setupPinia = (app) => {
  app.use(pinia);
};

export default (name) => {
  const useStore = stores[name];
  if (!useStore) {
    throw new Error(`store ${name} not found`);
  }
  // 传入 pinia 单例，避免“no active pinia”
  return useStore(pinia);
};
