import { mergeConfig } from 'vite';
import baseConfig from './vite.config.base';

export default mergeConfig(
  {
    mode: 'development',
    server: {
      open: false,
      fs: {
        strict: true,
      },
      host: true,
      port:9200,
      proxy: {
        '/business': {
          target: 'http://youdomain.cn/business',//代理的地址
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/business/, '')//这里的/需要转义
        }
      }
    },
    plugins: [
    ],
  },
  baseConfig
);
