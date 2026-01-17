import { loadEnv } from 'vite';
import uni from '@dcloudio/vite-plugin-uni';
import path from 'path';
// import viteCompression from 'vite-plugin-compression';
import uniReadPagesV3Plugin from './utils/router/utils/uni-read-pages-v3';

// https://vitejs.dev/config/
export default (command, mode) => {
	const env = loadEnv(mode, __dirname, 'GF_');

	return {
		envPrefix: "GF_",
		plugins: [
			uni(),
			// viteCompression({
			// 	verbose: false
			// }),
			uniReadPagesV3Plugin({
				pagesJsonDir: path.resolve(__dirname, './pages.json'),
				includes: ['path', 'aliasPath', 'name', 'meta'],
			}),
		],
		server: {
			host: true,
			open: false,
			port: env.GF_DEV_PORT || 9000, // 默认端口 9000
			hmr: {
				overlay: true,
			},
		},
		build: {
			// 增加构建超时时间（单位：毫秒）
			chunkSizeWarningLimit: 2000,
			// 优化构建性能
			minify: 'terser',
			terserOptions: {
				compress: {
					drop_console: false, // 开发环境保留 console
				},
			},
		},
		// 优化依赖预构建
		optimizeDeps: {
			include: ['vue', 'pinia', 'dayjs', 'lodash'],
		},
	};
};
