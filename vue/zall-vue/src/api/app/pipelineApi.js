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
export {
    listPipelineRequest,
    createPipelineRequest,
    updatePipelineRequest,
    deletePipelineRequest,
}