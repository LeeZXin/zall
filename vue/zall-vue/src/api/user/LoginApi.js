import request from '@/utils/request.js'

// 登录
const loginRequest = (data) => request.post("/api/login/login", data);
// 注册
const registerRequest = (data) => request.post("/api/login/register", data);
// 退出登录
const logoutRequest = () => request.get("/api/login/logout");
// 获取登录信息
const getUserInfoRequest = () => request.get("/api/login/userInfo");

export {
    loginRequest,
    registerRequest,
    logoutRequest,
    getUserInfoRequest
}