import request from '@/utils/request.js'

// 配置文件列表
const listPropertyFileRequest = (data) => request.get("/api/propertyFile/list", { params: data });
// 创建配置文件
const createPropertyFileRequest = (data) => request.post("/api/propertyFile/create", data);
// 版本历史
const listHistoryRequest = (data) => request.get("/api/propertyHistory/list", { params: data });
// 搜索版本号
const getHistoryByVersionRequest = (data) => request.get("/api/propertyHistory/getByVersion", { params: data });
// 新增版本
const newVersionRequest = (data) => request.post("/api/propertyHistory/newVersion", data);
// 配置来源列表
const listPropertySourceRequest = (data) => request.get("/api/propertySource/list", { params: data });
// 新增配置来源
const createPropertySourceRequest = (data) => request.post("/api/propertySource/create", data);
// 编辑配置来源
const updatePropertySourceRequest = (data) => request.post("/api/propertySource/update", data);
// 删除配置来源
const deletePropertySourceRequest = (sourceId) => request.delete("/api/propertySource/delete/" + sourceId);
// 配置来源列表
const listPropertySourceByFileIdRequest = (fileId) => request.get("/api/propertyFile/listSource/" + fileId);
// 删除配置文件
const deletePropertyFileRequest = (fileId) => request.delete("/api/propertyFile/delete/" + fileId);
// 发布配置
const deployHistoryRequest = (data) => request.post("/api/propertyHistory/deploy", data);
// 发布记录
const listDeployRequest = (historyId) => request.get("/api/propertyHistory/listDeploy/" + historyId);

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