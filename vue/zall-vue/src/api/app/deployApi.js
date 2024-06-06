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
// 启用配置
const enableDeployConfigRequest = (configId, env) => request.put("/api/deployConfig/enable/" + configId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 关闭配置
const disableDeployConfigRequest = (configId, env) => request.put("/api/deployConfig/disable/" + configId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listDeployConfigRequest,
    createDeployConfigRequest,
    updateDeployConfigRequest,
    enableDeployConfigRequest,
    disableDeployConfigRequest
}