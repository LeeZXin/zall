import axios from "axios";
import { message } from "ant-design-vue";
import i18n from "../language/i8n";
import router from "../router/router";
import { getLoginUser, setLoginUser } from "./login";
const t = i18n.global.t;

let redirectTimeout = null;

// 防止refreshtoken多次触发
let refreshPromise = null;

// 刷新登录token
const refreshToken = () => {
    if (refreshPromise) {
        return refreshPromise;
    }
    refreshPromise = new Promise((resolve, reject) => {
        let fn = async() => {
            try {
                const resp = await axios.get("/api/login/refresh");
                resolve(resp);
            } catch (error) {
                reject(error);
            }
        }
        fn();
    })
    refreshPromise.finally(() => {
        refreshPromise = null;
    });
    return refreshPromise;
}

const request = axios.create({
    baseURL: "/",
    timeout: 30000
});

request.interceptors.request.use(
    async(config) => {
        // 如果还差二十分钟token过期, 则刷新token
        const user = getLoginUser();
        const now = new Date().getTime();
        if (user && user.sessionExpireAt &&
            user.sessionExpireAt > now && user.sessionExpireAt < (now + 20 * 60 * 1000)) {
            try {
                const resp = await refreshToken();
                if (resp.status === 200 && resp.data.code === 0) {
                    user.sessionExpireAt = resp.data.expireAt;
                    user.sessionId = resp.data.sessionId;
                    setLoginUser(user);
                }
            } catch (error) {
                console.log(error);
            }
        }
        return config
    },
    (error) => {
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
        return Promise.reject(new Error(t("system.internalError")))
    },
    (error) => {
        if (error.response) {
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
                    }, 1000)
                }
            } else if (status === 403) {
                message.error(t("system.request403"))
            } else if (status === 400) {
                message.error(t("system.request400"))
            } else {
                message.error(t("system.internalError"))
            }
        }
        return Promise.reject(error)
    }
)

export default request