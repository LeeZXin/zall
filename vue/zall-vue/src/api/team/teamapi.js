import request from '@/utils/request.js'
// 获取团队列表
const getTeamListRequest = () => request({
    url: "/api/team/list",
    method: "GET"
});
// 创建团队
const createTeamRequest = (data) => request({
    url: "/api/team/create",
    method: "POST",
    data: data
});
// 是否是团队管理员
const isTeamAdminRequest = (data) => request({
    url: "/api/team/isAdmin",
    method: "POST",
    data: data
});
// 获取团队权限
const getTeamPermRequest = (data) => request({
    url: "/api/team/getTeamPerm",
    method: "POST",
    data: data
});
// 获取单个团队信息
const getTeamRequest = (data) => request({
    url: "/api/team/get",
    method: "POST",
    data: data
});
export {
    getTeamListRequest,
    createTeamRequest,
    isTeamAdminRequest,
    getTeamPermRequest,
    getTeamRequest
}