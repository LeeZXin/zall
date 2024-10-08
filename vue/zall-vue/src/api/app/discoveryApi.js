import request from '@/utils/request.js'

// 注册中心来源列表
const listDiscoverySourceRequest = (data) => request.get("/api/discoverySource/list", { params: data });
// 新增注册中心来源
const createDiscoverySourceRequest = (data) => request.post("/api/discoverySource/create", data);
// 编辑注册中心来源
const updateDiscoverySourceRequest = (data) => request.post("/api/discoverySource/update", data);
// 删除注册中心来源
const deleteDiscoverySourceRequest = (sourceId) => request.delete("/api/discoverySource/delete/" + sourceId);
// 服务列表
const listDiscoveryServiceRequest = (bindId) => request.get("/api/discoveryService/listService/" + bindId);
// 下线服务
const markAsDownServiceRequest = (data) => request.post("/api/discoveryService/markAsDown", data);
// 上线服务
const markAsUpServiceRequest = (data) => request.post("/api/discoveryService/markAsUp", data);
// 所有注册中心来源
const listAllDiscoverySourceRequest = (env) => request.get("/api/discoverySource/listAll/" + env);
// 获取注册中心来源
const listBindDiscoverySourceRequest = (data) => request.get("/api/discoverySource/listBind", { params: data });
// 绑定应用服务和注册中心来源
const bindAppAndDiscoverySourceRequest = (data) => request.post("/api/discoverySource/bindApp", data);
export {
    listDiscoverySourceRequest,
    createDiscoverySourceRequest,
    updateDiscoverySourceRequest,
    deleteDiscoverySourceRequest,
    listDiscoveryServiceRequest,
    markAsDownServiceRequest,
    markAsUpServiceRequest,
    listAllDiscoverySourceRequest,
    listBindDiscoverySourceRequest,
    bindAppAndDiscoverySourceRequest
}