import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 服务来源配置列表
const listServiceSourceRequest = (data, env) => request.get("/api/service/listSource", {
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
export {
    listServiceSourceRequest,
    listServiceStatusRequest,
    listStatusActionsRequest,
    doStatusActionRequest
}