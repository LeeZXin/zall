import request from '@/utils/request.js'

const ENV_HEADER = "ZALL-ENV";

// 发布计划列表
const listDeployPlanRequest = (data, env) => request.get("/api/deployPlan/list", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 创建发布计划
const createDeployPlanRequest = (data, env) => request.post("/api/deployPlan/create", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 流水线配置列表
const listPipelineWhenCreateDeployPlanRequest = (data, env) => request.get("/api/deployPlan/listPipeline", {
    params: data,
    headers: {
        [ENV_HEADER]: env
    }
});
// 关闭发布任务
const closeDeployPlanRequest = (planId, env) => request.put("/api/deployPlan/close/" + planId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 执行发布任务
const startDeployPlanRequest = (planId, env) => request.put("/api/deployPlan/start/" + planId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 发布计划详情
const getDeployPlanDetailRequest = (planId, env) => request.get("/api/deployPlan/detail/" + planId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 发布计划流水线详情
const listDeployPlanStagesRequest = (planId, env) => request.get("/api/deployPlan/listStages/" + planId, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 重新执行stage
const redoDeployAgentStageRequest = (stageId, env) => request.put("/api/deployStage/redoAgent/" + stageId, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// kill stage
const killDeployStageRequest = (planId, index, env) => request.put(`/api/deployStage/kill/${planId}/${index}`, {}, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 确认交互阶段
const confirmInteractStageRequest = (data, env) => request.post("/api/deployStage/confirm", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
// 强制重新执行
const forceRedoStageRequest = (data, env) => request.post("/api/deployStage/forceRedoStage", data, {
    headers: {
        [ENV_HEADER]: env
    }
});
export {
    listDeployPlanRequest,
    createDeployPlanRequest,
    listPipelineWhenCreateDeployPlanRequest,
    closeDeployPlanRequest,
    getDeployPlanDetailRequest,
    startDeployPlanRequest,
    listDeployPlanStagesRequest,
    redoDeployAgentStageRequest,
    killDeployStageRequest,
    confirmInteractStageRequest,
    forceRedoStageRequest
}