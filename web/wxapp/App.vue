<script>
	// #ifdef MP-WEIXIN
	// 兼容：部分微信开发者工具/运行环境的 vue-devtools 注入会触发：
	// __VUE_DEVTOOLS_ON_SOCKET_READY__ is not a function（导致启动白屏、页面未注册）
	// 这里提前挂一个空函数兜底，避免 devtools backend 初始化直接抛错。
	;(function () {
		try {
			const g =
				(typeof globalThis !== 'undefined' && globalThis) ||
				(typeof global !== 'undefined' && global) ||
				(typeof wx !== 'undefined' && wx) ||
				{}
			if (g && typeof g.__VUE_DEVTOOLS_ON_SOCKET_READY__ !== 'function') {
				g.__VUE_DEVTOOLS_ON_SOCKET_READY__ = function () {}
			}
			// 有些实现 target 可能指向 wx 对象
			if (typeof wx !== 'undefined' && typeof wx.__VUE_DEVTOOLS_ON_SOCKET_READY__ !== 'function') {
				wx.__VUE_DEVTOOLS_ON_SOCKET_READY__ = function () {}
			}
		} catch (e) {
			// ignore
		}
	})()
	// #endif

	import $store from '@/store'

	export default {
		data() {
			return {
				__bootstrapped: false
			}
		},
		methods: {
			async bootstrapAuthRedirect() {
				// 只在启动时跑一次，避免反复 reLaunch
				if (this.__bootstrapped) return
				this.__bootstrapped = true

				const token = uni.getStorageSync('token')
				if (!token) return

				const userStore = $store('user')
				if (!userStore.isLogin) userStore.setToken(token)

				// 校验 token 是否有效：有效则从登录页跳到首页；无效则清理并停留在登录页
				let ok = false
				try {
					const info = await userStore.getInfo()
					ok = !!info
				} catch (e) {
					ok = false
				}
				if (!ok) {
					try {
						await userStore.logout(true)
					} catch (e) {}
					return
				}

				const pages = getCurrentPages ? getCurrentPages() : []
				const cur = pages && pages.length ? pages[pages.length - 1] : null
				const route = cur && cur.route ? cur.route : ''
				// 仅在“启动落到登录页”时自动跳首页，避免干扰用户在其它页面的操作
				if (route === 'pages/login/login') {
					uni.reLaunch({ url: '/pages/property_list/property_list' })
				}
			}
		},
		onLaunch: function() {
			console.log('App Launch')
			// #ifdef MP-WEIXIN
			// 微信小程序端：WXSS 不允许/不稳定地使用本地 @font-face 引入字体（会触发 do-not-use-local-path 并 500）
			// 因此这里改回使用 wx.loadFontFace 加载“稳定 CDN”字体文件
			// 注意：真机/发布需把对应域名加入小程序「downloadFile 合法域名」
			try {
				if (typeof wx !== 'undefined' && wx.loadFontFace) {
					wx.loadFontFace({
						family: 'Material Symbols Outlined',
						// 使用 jsDelivr CDN（相对 raw.githubusercontent.com 更稳定）//TODO后续放入自己的服务器中
						source:
							'url("https://cdn.jsdelivr.net/gh/google/material-design-icons@master/font/MaterialIconsOutlined-Regular.otf")',
						global: true
					})
				}
			} catch (e) {
				// ignore
			}
			// #endif

			// 启动时：若已登录则自动跳首页
			this.bootstrapAuthRedirect()
		},
		onShow: function() {
			console.log('App Show')
			// 某些情况下 onLaunch 早于页面栈就绪，这里再兜底一次
			this.bootstrapAuthRedirect()
		},
		onHide: function() {
			console.log('App Hide')
		}
	}
</script>

<style lang="scss">
	/*每个页面公共css */
	/* 字体图标：改为加载本地字体，避免外链失效 */
	/* #ifndef MP-WEIXIN */
	@font-face {
		font-family: 'Material Symbols Outlined';
		font-style: normal;
		font-weight: 400;
		src:
			url('/static/fonts/MaterialSymbolsOutlined.ttf') format('truetype'),
			url('/static/fonts/MaterialSymbolsOutlined.otf') format('opentype');
	}
	/* #endif */

	.material-symbols-outlined {
		font-family: 'Material Symbols Outlined';
		font-weight: 400;
		font-style: normal;
		line-height: 1;
		letter-spacing: normal;
		text-transform: none;
		display: inline-block;
		white-space: nowrap;
		word-wrap: normal;
		direction: ltr;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
		/* #ifdef MP-WEIXIN */
		/* 强制启用连字（ligatures），否则会直接显示英文单词 */
		font-variant-ligatures: contextual;
		-webkit-font-variant-ligatures: contextual;
		font-feature-settings: 'liga';
		-webkit-font-feature-settings: 'liga';
		/* #endif */
		/* #ifdef H5 */
		/* H5 支持变量字体参数，小程序端可能不支持，避免引发 WXSS 解析问题 */
		font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;
		/* #endif */
	}
	
	page {
		background-color: #f6f7f8;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
	}

	/* 避免使用通配选择器 `*`，部分微信编译器会直接报错 */
	view,
	text,
	image,
	button,
	input,
	textarea,
	scroll-view,
	swiper,
	swiper-item {
		box-sizing: border-box;
	}
</style>
