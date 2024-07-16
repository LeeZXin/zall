import request from '@/utils/request.js'

// 服务来源配置列表
const listSimpleServiceSourceRequest = (data) => request.get("/api/service/listSource", { params: data });
// 服务状态配置列表
const listServiceStatusRequest = (sourceId) => request.get("/api/service/listStatus/" + sourceId);
// 服务操作列表
const listStatusActionsRequest = (sourceId) => request.get("/api/service/listActions/" + sourceId);
// 操作服务
const doStatusActionRequest = (data) => request.put(`/api/service/doAction`, data);
// 服务来源配置列表
const listServiceSourceRequest = (data) => request.get("/api/serviceSource/list", { params: data });
// 创建服务来源配置
const createServiceSourceRequest = (data) => request.post("/api/serviceSource/create", data);
// 编辑服务来源配置
const updateServiceSourceRequest = (data) => request.post("/api/serviceSource/update", data);
// 删除服务来源配置
const deleteServiceSourceRequest = (sourceId) => request.delete("/api/serviceSource/delete/" + sourceId);

export {
    listSimpleServiceSourceRequest,
    listServiceStatusRequest,
    listStatusActionsRequest,
    doStatusActionRequest,
    listServiceSourceRequest,
    createServiceSourceRequest,
    updateServiceSourceRequest,
    deleteServiceSourceRequest,
}