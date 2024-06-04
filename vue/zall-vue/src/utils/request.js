import axios from 'axios';
import { message } from 'ant-design-vue'
import { useUserStore } from '@/pinia/userStore.js'
import i18n from "../language/i8n"
import router from "../router/router";

const t = i18n.global.t;

let redirectTimeout = null;

const request = axios.create({
    baseURL: "/",
    timeout: 30000
});

request.interceptors.request.use(
    (config) => {
        const user = useUserStore();
        const now = new Date().getTime();
        if (user &&
            user.sessionExpireAt &&
            user.sessionExpireAt < now &&
            user.sessionExpireAt > (now - 10 * 60 * 1000)) {
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
        if (data.code === 0) {
            return data;
        }
        if (data.message) {
            message.error(data.message)
            return Promise.reject(new Error(data.message))
        }
        return Promise.reject(new Error("返回信息有误"))
    },
    (error) => {
        const status = error.response.status
        if (status === 404) {
            message.error(t("system.request404"))
        } else if (status === 401) {
            if (!redirectTimeout) {
                message.error(t("system.notLogin"))
                redirectTimeout = setTimeout(() => {
                    redirectTimeout = null;
                    router.push({
                        path: "/login/login",
                        query: {
                            redirect_uri: encodeURI(window.location.href)
                        }
                    })
                }, 3000)
            }
        } else if (status === 403) {
            message.error(t("system.request403"))
        } else if (status === 400) {
            message.error(t("system.request400"))
        } else {
            message.error(t("system.internalError"))
        }
        return Promise.reject(error)
    }
)

export default request