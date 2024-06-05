import request from '@/utils/request.js'

// 应用服务列表
const listAppRequest = (teamId) => request.get("/api/app/list/" + teamId);
// 应用服务列表
const listAllAppByAdminRequest = (teamId) => request.get("/api/app/listAllByAdmin/" + teamId);
// 创建应用服务
const createAppRequest = (data) => request.post("/api/app/create", data);

export {
    listAppRequest,
    createAppRequest,
    listAllAppByAdminRequest
}