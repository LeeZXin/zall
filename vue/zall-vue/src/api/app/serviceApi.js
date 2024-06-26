import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 服务配置列表
const listServiceRequest = (data, env) => request.get("/api/service/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建服务配置
const createServiceRequest = (data, env) => request.post("/api/service/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑服务配置
const updateServiceRequest = (data, env) => request.post("/api/service/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除服务配置
const deleteServiceRequest = (serviceId, env) => request.delete("/api/service/delete/" + serviceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listServiceRequest,
    createServiceRequest,
    updateServiceRequest,
    deleteServiceRequest,
}