import axios from 'axios';
import { message } from 'ant-design-vue'
import { useUserStore } from '@/pinia/userStore.js'
import i18n from "../language/i8n"
import router from "../router/router";

const t = i18n.global.t;

const request = axios.create({
    baseURL: "/",
    timeout: 30000
})

request.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8'
request.defaults.transformRequest = [
    (data) => {
        return JSON.stringify(data);
    }
]

request.interceptors.request.use(
    (config) => {
        const user = useUserStore();
        const now = new Date().getTime();
        if (user && user.sessionExpireAt && user.sessionExpireAt < now && user.sessionExpireAt > (now - 10 * 60 * 1000)) {
            console.log("refresh token")
        }
        return config
    },
    (error) => {
        console.log(error)
        message.error(t("system.requestFailed"));
        return Promise.reject(error)
    }
)

request.interceptors.response.use(
    (response) => {
        let data = response.data
        if (data.code !== 0) {
            if (data.message) {
                message.error(data.message)
                return Promise.reject(new Error(data.message))
            }
        }
        return data
    },
    (error) => {
        console.log(error)
        const status = error.response.status
        if (status === 404) {
            message.error(t("system.request404"))
        } else if (status === 401) {
            message.error(t("system.notLogin"))
            setTimeout(() => {
                router.push({
                    path: "/login/login",
                    query: {
                        redirect_uri: encodeURI(window.location.href)
                    }
                })
            }, 3000)
        } else if (status === 403) {
            message.error(t("system.request403"))
        } else {
            message.error(t("system.internalError"))
        }
        return Promise.reject(error)
    }
)

export default request