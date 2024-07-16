import request from '@/utils/request.js'

// 注册中心来源列表
const listDiscoverySourceRequest = (data) => request.get("/api/discoverySource/list", { params: data });
// 新增注册中心来源
const createDiscoverySourceRequest = (data) => request.post("/api/discoverySource/create", data);
// 编辑注册中心来源
const updateDiscoverySourceRequest = (data) => request.post("/api/discoverySource/update", data);
// 删除注册中心来源
const deleteDiscoverySourceRequest = (sourceId) => request.delete("/api/discoverySource/delete/" + sourceId);
// 注册中心来源列表
const listSimpleDiscoverySourceRequest = (data) => request.get("/api/discoveryService/listSource", { params: data });
// 服务列表
const listDiscoveryServiceRequest = (sourceId) => request.get("/api/discoveryService/listService/" + sourceId);
// 下线服务
const deregisterServiceRequest = (data) => request.post("/api/discoveryService/deregister", data);
// 上线服务
const reRegisterServiceRequest = (data) => request.post("/api/discoveryService/reRegister", data);
// 删除下线服务
const deleteDownServiceRequest = (data) => request.delete("/api/discoveryService/deleteDownService", { params: data });
export {
    listDiscoverySourceRequest,
    createDiscoverySourceRequest,
    updateDiscoverySourceRequest,
    deleteDiscoverySourceRequest,
    listSimpleDiscoverySourceRequest,
    listDiscoveryServiceRequest,
    deregisterServiceRequest,
    reRegisterServiceRequest,
    deleteDownServiceRequest
}