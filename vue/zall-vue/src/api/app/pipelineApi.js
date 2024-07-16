import request from '@/utils/request.js'

// 流水线配置列表
const listPipelineRequest = (data) => request.get("/api/pipeline/list", { params: data });
// 创建流水线配置
const createPipelineRequest = (data) => request.post("/api/pipeline/create", data);
// 编辑流水线配置
const updatePipelineRequest = (data) => request.post("/api/pipeline/update", data);
// 删除流水线配置
const deletePipelineRequest = (pipelineId) => request.delete("/api/pipeline/delete/" + pipelineId);
// 流水线变量配置列表
const listPipelineVarsRequest = (data) => request.get("/api/pipelineVars/list", { params: data });
// 创建流水线变量
const createPipelineVarsRequest = (data) => request.post("/api/pipelineVars/create", data);
// 编辑流水线变量
const updatePipelineVarsRequest = (data) => request.post("/api/pipelineVars/update", data);
// 删除流水线变量
const deletePipelineVarsRequest = (varsId) => request.delete("/api/pipelineVars/delete/" + varsId);
// 获取流水线变量
const getPipelineVarsRequest = (varsId) => request.get("/api/pipelineVars/content/" + varsId);
export {
    listPipelineRequest,
    createPipelineRequest,
    updatePipelineRequest,
    deletePipelineRequest,
    listPipelineVarsRequest,
    createPipelineVarsRequest,
    updatePipelineVarsRequest,
    deletePipelineVarsRequest,
    getPipelineVarsRequest
}