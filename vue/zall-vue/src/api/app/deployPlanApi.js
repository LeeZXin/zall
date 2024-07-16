import request from '@/utils/request.js'

// 发布计划列表
const listDeployPlanRequest = (data) => request.get("/api/deployPlan/list", { params: data });
// 创建发布计划
const createDeployPlanRequest = (data) => request.post("/api/deployPlan/create", data);
// 流水线配置列表
const listPipelineWhenCreateDeployPlanRequest = (data) => request.get("/api/deployPlan/listPipeline", { params: data });
// 关闭发布任务
const closeDeployPlanRequest = (planId) => request.put("/api/deployPlan/close/" + planId);
// 执行发布任务
const startDeployPlanRequest = (planId) => request.put("/api/deployPlan/start/" + planId);
// 发布计划详情
const getDeployPlanDetailRequest = (planId) => request.get("/api/deployPlan/detail/" + planId);
// 发布计划流水线详情
const listDeployPlanStagesRequest = (planId) => request.get("/api/deployPlan/listStages/" + planId);
// 重新执行stage
const redoDeployAgentStageRequest = (stageId) => request.put("/api/deployStage/redoAgent/" + stageId);
// kill stage
const killDeployStageRequest = (planId, index) => request.put(`/api/deployStage/kill/${planId}/${index}`);
// 确认交互阶段
const confirmInteractStageRequest = (data) => request.post("/api/deployStage/confirm", data);
// 强制重新执行
const forceRedoStageRequest = (data) => request.post("/api/deployStage/forceRedoStage", data);
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