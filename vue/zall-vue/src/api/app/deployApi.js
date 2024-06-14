import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 部署配置列表
const listDeployConfigRequest = (data, env) => request.get("/api/deployConfig/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建配置
const createDeployConfigRequest = (data, env) => request.post("/api/deployConfig/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑配置
const updateDeployConfigRequest = (data, env) => request.post("/api/deployConfig/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除配置
const deleteDeployConfigRequest = (configId, env) => request.delete("/api/deployConfig/delete/" + configId, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listDeployConfigRequest,
    createDeployConfigRequest,
    updateDeployConfigRequest,
    deleteDeployConfigRequest
}