import request from '@/utils/request.js'

// 登录
const loginRequest = (data) => request({
    url: "/api/login/login",
    method: "POST",
    data: data
});

// 注册
const registerRequest = (data) => request({
    url: "/api/login/register",
    method: "POST",
    data: data
});

// 退出登录
const logoutRequest = () => request({
    url: "/api/login/logout",
    method: "GET",
});

// 获取登录信息
const getUserInfoRequest = () => request({
    url: "/api/login/userInfo",
    method: "GET",
})

export {
    loginRequest,
    registerRequest,
    logoutRequest,
    getUserInfoRequest
}