import { createApp } from 'vue';
import App from './App.vue';
import 'ant-design-vue/dist/reset.css';
import './assets/css/main.css';
import Antd from 'ant-design-vue';
import router from "@/router/router";
import { createPinia } from 'pinia';
import i18n from "./language/i8n";
const app = createApp(App)
app.use(Antd)
app.use(router)
app.use(createPinia())
app.use(i18n);
app.mount('#app')