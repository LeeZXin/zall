import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 服务来源配置列表
const listSimpleServiceSourceRequest = (data, env) => request.get("/api/service/listSource", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 服务状态配置列表
const listServiceStatusRequest = (sourceId, env) => request.get("/api/service/listStatus/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 服务操作列表
const listStatusActionsRequest = (sourceId, env) => request.get("/api/service/listActions/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 操作服务
const doStatusActionRequest = (data, env) => request.put(`/api/service/doAction`, data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 服务来源配置列表
const listServiceSourceRequest = (data, env) => request.get("/api/serviceSource/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建服务来源配置
const createServiceSourceRequest = (data, env) => request.post("/api/serviceSource/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑服务来源配置
const updateServiceSourceRequest = (data, env) => request.post("/api/serviceSource/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除服务来源配置
const deleteServiceSourceRequest = (sourceId, env) => request.delete("/api/serviceSource/delete/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
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