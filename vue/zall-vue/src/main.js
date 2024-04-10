import { createApp } from 'vue'
import App from './App.vue'
import 'ant-design-vue/dist/reset.css'
import './assets/css/main.css'
import Antd from 'ant-design-vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createPinia } from 'pinia'
import i18n from "./language/i8n"
const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: "",
            redirect: "/login/login"
        },
        {
            path: "/login",
            redirect: "/login/login",
            component: () => import("./layouts/LoginLayout.vue"),
            children: [
                {
                    path: "/login/login",
                    component: () => import("./pages/login/LoginPage.vue")
                },
                {
                    path: "/login/register",
                    component: () => import("./pages/login/RegisterPage.vue")
                }
            ]

        }
    ]
})

const app = createApp(App)
app.use(Antd)
app.use(router)
app.use(createPinia())
app.use(i18n);

app.mount('#app')