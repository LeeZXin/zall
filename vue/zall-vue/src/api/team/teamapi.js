import request from '@/utils/request.js'
// 获取团队列表
const getTeamListRequest = () => request.get("/api/team/list");
// 创建团队
const createTeamRequest = (data) => request.post("/api/team/create", data);
// 是否是团队管理员
const isTeamAdminRequest = (teamId) => request.get("/api/team/isAdmin/" + teamId);
// 获取团队权限
const getTeamPermRequest = (teamId) => request.get("/api/team/getTeamPerm/" + teamId);
// 获取单个团队信息
const getTeamRequest = (teamId) => request.get("/api/team/get/" + teamId);
// 获取团队成员账号
const listAccountRequest = (teamId) => request.get("/api/team/listAccount/" + teamId);
export {
    getTeamListRequest,
    createTeamRequest,
    isTeamAdminRequest,
    getTeamPermRequest,
    getTeamRequest,
    listAccountRequest
}