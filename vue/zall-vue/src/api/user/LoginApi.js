import request from '@/utils/request.js'
export default {
    // 登录
    Login: (data) => request({
        url: "/api/login/login",
        method: "POST",
        data: data
    }),
    // 注册
    Register: (data) => request({
        url: "/api/login/register",
        method: "POST",
        data: data
    }),
    // 退出登录
    LoginOut: () => request({
        url: "/api/login/loginOut",
        method: "GET",
    })
}