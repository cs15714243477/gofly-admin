import { createApp } from 'vue';
import ArcoVue from '@arco-design/web-vue';
import ArcoVueIcon from '@arco-design/web-vue/es/icon';
import globalComponents from '@/components';
import router from './router';
import store from './store';
import i18n from './locale';
import directive from './directive';
import App from './App.vue';
//检测元素大小变化
import directives from '@/utils/directive/index'
import '@arco-design/web-vue/dist/arco.css';
import '@/assets/style/global.less';
import '@/assets/style/gfui.less';
import '@/assets/style/cover-arco.less';//覆盖arco样式
const app = createApp(App);
// 屏蔽警告信息
app.config.warnHandler = () => null;

app.use(ArcoVue, {});
app.use(ArcoVueIcon);

app.use(router);
app.use(store);
app.use(i18n);
app.use(globalComponents);
app.use(directive);
app.use(directives);

app.mount('#app');
