import request from '@/utils/request.js'

// 服务状态配置列表
const listServiceStatusRequest = (sourceId) => request.get("/api/service/listStatus/" + sourceId);
// 杀死服务
const killServiceRequest = (data) => request.put(`/api/service/kill`, data);
// 重启服务
const restartServiceRequest = (data) => request.put(`/api/service/restart`, data);
// 服务来源配置列表
const listServiceSourceRequest = (data) => request.get("/api/serviceSource/list", { params: data });
// 创建服务来源配置
const createServiceSourceRequest = (data) => request.post("/api/serviceSource/create", data);
// 编辑服务来源配置
const updateServiceSourceRequest = (data) => request.post("/api/serviceSource/update", data);
// 删除服务来源配置
const deleteServiceSourceRequest = (sourceId) => request.delete("/api/serviceSource/delete/" + sourceId);
// 所有服务来源
const listAllServiceSourceRequest = (env) => request.get("/api/serviceSource/listAll/" + env);
// 获取绑定服务来源
const listBindServiceSourceRequest = (data) => request.get("/api/serviceSource/listBind", { params: data });
// 绑定应用服务和服务来源
const bindAppAndServiceSourceRequest = (data) => request.post("/api/serviceSource/bindApp", data);
export {
    listServiceStatusRequest,
    killServiceRequest,
    restartServiceRequest,
    listServiceSourceRequest,
    createServiceSourceRequest,
    updateServiceSourceRequest,
    deleteServiceSourceRequest,
    listAllServiceSourceRequest,
    listBindServiceSourceRequest,
    bindAppAndServiceSourceRequest
}