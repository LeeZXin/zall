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
export {
    listPropertyFileRequest,
    createPropertyFileRequest,
    listHistoryRequest,
    getHistoryByVersionRequest,
    newVersionRequest
}