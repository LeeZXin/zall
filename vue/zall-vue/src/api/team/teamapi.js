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
const isTeamAdminRequest = (teamId) => request({
    url: "/api/team/isAdmin/" + teamId,
    method: "GET"
});
// 获取团队权限
const getTeamPermRequest = (teamId) => request({
    url: "/api/team/getTeamPerm/" + teamId,
    method: "GET"
});
// 获取单个团队信息
const getTeamRequest = (teamId) => request({
    url: "/api/team/get/" + teamId,
    method: "GET"
});
export {
    getTeamListRequest,
    createTeamRequest,
    isTeamAdminRequest,
    getTeamPermRequest,
    getTeamRequest
}