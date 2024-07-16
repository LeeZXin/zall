import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 注册中心来源列表
const listDiscoverySourceRequest = (data, env) => request.get("/api/discoverySource/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 新增注册中心来源
const createDiscoverySourceRequest = (data, env) => request.post("/api/discoverySource/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑注册中心来源
const updateDiscoverySourceRequest = (data, env) => request.post("/api/discoverySource/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除注册中心来源
const deleteDiscoverySourceRequest = (sourceId, env) => request.delete("/api/discoverySource/delete/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 注册中心来源列表
const listSimpleDiscoverySourceRequest = (data, env) => request.get("/api/discoveryService/listSource", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 服务列表
const listDiscoveryServiceRequest = (sourceId, env) => request.get("/api/discoveryService/listService/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 下线服务
const deregisterServiceRequest = (data, env) => request.post("/api/discoveryService/deregister", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 上线服务
const reRegisterServiceRequest = (data, env) => request.post("/api/discoveryService/reRegister", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除下线服务
const deleteDownServiceRequest = (data, env) => request.delete("/api/discoveryService/deleteDownService", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
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