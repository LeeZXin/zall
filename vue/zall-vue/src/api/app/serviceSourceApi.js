import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

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
    listServiceSourceRequest,
    createServiceSourceRequest,
    updateServiceSourceRequest,
    deleteServiceSourceRequest,
}