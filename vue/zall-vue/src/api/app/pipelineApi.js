import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 流水线配置列表
const listPipelineRequest = (data, env) => request.get("/api/pipeline/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建流水线配置
const createPipelineRequest = (data, env) => request.post("/api/pipeline/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑流水线配置
const updatePipelineRequest = (data, env) => request.post("/api/pipeline/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除流水线配置
const deletePipelineRequest = (pipelineId, env) => request.delete("/api/pipeline/delete/" + pipelineId, {
    headers: {
        [ENV_HEADER]: env
    }
});

// 流水线变量配置列表
const listPipelineVarsRequest = (data, env) => request.get("/api/pipelineVars/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建流水线变量
const createPipelineVarsRequest = (data, env) => request.post("/api/pipelineVars/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 编辑流水线变量
const updatePipelineVarsRequest = (data, env) => request.post("/api/pipelineVars/update", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 删除流水线变量
const deletePipelineVarsRequest = (varsId, env) => request.delete("/api/pipelineVars/delete/" + varsId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 获取流水线变量
const getPipelineVarsRequest = (varsId, env) => request.get("/api/pipelineVars/content/" + varsId, {
    headers: {
        [ENV_HEADER]: env
    }
});
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