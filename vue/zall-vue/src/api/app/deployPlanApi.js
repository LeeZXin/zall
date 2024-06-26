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
// 服务配置列表
const listServiceWhenCreateDeployPlanRequest = (data, env) => request.get("/api/deployPlan/listService", {
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
export {
    listDeployPlanRequest,
    createDeployPlanRequest,
    listServiceWhenCreateDeployPlanRequest,
    closeDeployPlanRequest,
}