import axios from 'axios';
import { message } from 'ant-design-vue'
import { userStore } from '@/pinia/UserStore.js'
import i18n from "../language/i8n"
import { useRouter } from "vue-router";

const t = i18n.global.t;
const router = useRouter();

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
        const user = userStore()
        const now = new Date().getTime();
        if (user.sessionExpireAt > 0 && user.sessionExpireAt < now) {
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
                return Promise.reject(new error(data.message))
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