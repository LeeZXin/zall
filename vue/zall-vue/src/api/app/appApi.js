import request from '@/utils/request.js'

// 应用服务列表
const listAppRequest = (teamId) => request.get("/api/app/list/" + teamId);
// 应用服务列表
const listAllAppByAdminRequest = (teamId) => request.get("/api/app/listAllByAdmin/" + teamId);
// 应用服务列表
const listAllAppBySaRequest = () => request.get("/api/app/listAllBySa");
// 创建应用服务
const createAppRequest = (data) => request.post("/api/app/create", data);
// 获取应用服务
const getAppRequest = (appId) => request.get("/api/app/get/" + appId);
// 编辑应用服务
const updateAppRequest = (data) => request.post("/api/app/update", data);
// 删除应用服务
const deleteAppRequest = (appId) => request.delete("/api/app/delete/" + appId);
// 迁移团队
const transferAppRequest = (data) => request.put("/api/app/transferTeam", data);
export {
    listAppRequest,
    createAppRequest,
    listAllAppByAdminRequest,
    listAllAppBySaRequest,
    getAppRequest,
    updateAppRequest,
    deleteAppRequest,
    transferAppRequest
}