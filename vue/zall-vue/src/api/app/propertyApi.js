import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 配置文件列表
const listPropertyFileRequest = (data, env) => request.get("/api/propertyFile/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建配置文件
const createPropertyFileRequest = (data, env) => request.post("/api/propertyFile/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 版本历史
const listHistoryRequest = (data, env) => request.get("/api/propertyHistory/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 搜索版本号
const getHistoryByVersionRequest = (data, env) => request.get("/api/propertyHistory/getByVersion", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 新增版本
const newVersionRequest = (data, env) => request.post("/api/propertyHistory/newVersion", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 配置来源列表
const listPropertySourceRequest = (data, env) => request.get("/api/propertySource/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 新增配置来源
const createPropertySourceRequest = (data, env) => request.post("/api/propertySource/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑配置来源
const updatePropertySourceRequest = (data, env) => request.post("/api/propertySource/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除配置来源
const deletePropertySourceRequest = (sourceId, env) => request.delete("/api/propertySource/delete/" + sourceId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 配置来源列表
const listPropertySourceByFileIdRequest = (fileId, env) => request.get("/api/propertyFile/listSource/" + fileId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除配置文件
const deletePropertyFileRequest = (fileId, env) => request.delete("/api/propertyFile/delete/" + fileId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 发布配置
const deployHistoryRequest = (data, env) => request.post("/api/propertyHistory/deploy", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 发布记录
const listDeployRequest = (historyId, env) => request.get("/api/propertyHistory/listDeploy/" + historyId, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listPropertyFileRequest,
    createPropertyFileRequest,
    listHistoryRequest,
    getHistoryByVersionRequest,
    newVersionRequest,
    listPropertySourceRequest,
    createPropertySourceRequest,
    updatePropertySourceRequest,
    deletePropertySourceRequest,
    listPropertySourceByFileIdRequest,
    deployHistoryRequest,
    listDeployRequest,
    deletePropertyFileRequest
}